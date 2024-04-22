import argparse
import json
import math
import os
import shutil

import deepspeed
import torch
from peft import LoraConfig, get_peft_model
from torch.utils.data import DataLoader, RandomSampler
from torch.utils.data.distributed import DistributedSampler
from tqdm import tqdm

from data_preparation import data_enhancement
from merge_lora import merge_lora
from model import MODE
from retrieval_side import is_supported_format
from utils import DataCollator
from utils import print_trainable_parameters, print_rank_0, to_device, set_random_seed, save_model

try:
    from torch.utils.tensorboard import SummaryWriter
except ImportError:
    from tensorboard import SummaryWriter


def log_info(rank, epoch, step, loss, learning_rate):
    return {
        'rank': rank,
        'epoch': epoch,
        'step': step,
        'loss': loss,
        'learning_rate': learning_rate
    }


def save_checkpoint(model, tokenizer, epoch, global_step, args, ds_config):
    checkpoint_path = os.path.join(args.output_dir, f"epoch-{epoch}-step-{global_step}")
    if ds_config["zero_optimization"]["stage"] == 3:
        state_dict = model._zero3_consolidated_16bit_state_dict()
    else:
        state_dict = None
    if args.global_rank <= 0:
        save_model(model, tokenizer, checkpoint_path, '', state_dict)
    return checkpoint_path


def manage_checkpoints(checkpoints, max_checkpoints):
    if len(checkpoints) > max_checkpoints:
        # Remove oldest checkpoints
        for checkpoint_to_delete in checkpoints[:len(checkpoints) - max_checkpoints]:
            if os.path.exists(checkpoint_to_delete):
                if os.path.isdir(checkpoint_to_delete):  # 确保是文件夹
                    shutil.rmtree(checkpoint_to_delete)
                else:
                    print(f"Error: {checkpoint_to_delete} is not a directory")
            else:
                print(f"Error: {checkpoint_to_delete} does not exist")
        return checkpoints[len(checkpoints) - max_checkpoints:]
    return checkpoints


def enhance_data(input_path, output_dir):
    if os.path.isdir(input_path):
        enhanced_dir = os.path.join(output_dir, "enhancement")
        os.makedirs(enhanced_dir, exist_ok=True)
        for filename in os.listdir(input_path):
            if is_supported_format(filename):
                input_file = os.path.join(input_path, filename)
                enhancement_data_path = os.path.join(enhanced_dir, filename)
                data_enhancement(input_file, enhancement_data_path)
        input_path = enhanced_dir
    else:
        if is_supported_format(input_path):
            output_dir = os.path.dirname(input_path)
            base_filename, ext = os.path.splitext(os.path.basename(input_path))
            enhancement_data_path = os.path.join(output_dir, f"{base_filename}_enhanced{ext}")
            data_enhancement(input_path, enhancement_data_path)
        else:
            print("不支持的文件格式")
            return None
        input_path = enhancement_data_path

    return input_path


def parse_args():
    parser = argparse.ArgumentParser()
    # Model
    parser.add_argument("--model_name_or_path",
                        type=str, help="", required=True)
    # DataSet
    parser.add_argument("--train_path", default="", type=str, help="")
    parser.add_argument("--enhancement", default="false", type=str, help="")
    parser.add_argument("--retrieval_method", type=str, default="st", choices=["bm25", "st"],
                        help="Method for document retrieval")
    parser.add_argument("--st", type=str, default='BAAI/bge-base-zh-v1.5')
    parser.add_argument("--top_k", type=int, default=1)
    parser.add_argument("--max_len", type=int, default=1024, help="")
    parser.add_argument("--max_src_len", type=int, default=256, help="")
    parser.add_argument("--is_skip", action='store_true', help="")
    # Train
    parser.add_argument("--per_device_train_batch_size",
                        type=int, default=16, help="")
    parser.add_argument("--learning_rate", type=float, default=1e-3, help="")
    parser.add_argument("--weight_decay", type=float, default=0.1, help="")
    parser.add_argument("--num_train_epochs", type=int, default=1, help="")
    parser.add_argument("--gradient_accumulation_steps",
                        type=int, default=1, help="")
    parser.add_argument("--warmup_ratio", type=float, default=0.1, help="")
    parser.add_argument("--output_dir", type=str, default=None, help="")
    parser.add_argument("--mode", type=str, default="glm2", help="")
    parser.add_argument("--train_type", type=str, default="lora", help="")
    parser.add_argument("--seed", type=int, default=1234, help="")
    parser.add_argument("--local_rank", type=int, default=-1, help="")
    parser.add_argument("--show_loss_step", default=10, type=int, help="")
    parser.add_argument("--gradient_checkpointing",
                        action='store_true', help="")
    parser.add_argument("--save_model_step", default=100, type=int, help="")
    # deepspeed features
    parser.add_argument("--ds_file", type=str,
                        default="ds_zero2.json", help="")
    # LoRA
    parser.add_argument("--lora_dim", type=int, default=8, help="")
    parser.add_argument("--lora_alpha", type=int, default=30, help="")
    parser.add_argument("--lora_dropout", type=float, default=0.1, help="")
    parser.add_argument("--lora_module_name", type=str,
                        default="query_key_value", help="")
    parser.add_argument("--merge_lora", default="true", type=str, help="")
    # Freeze
    parser.add_argument("--freeze_module_name", type=str,
                        default="layers.27.", help="")
    # P-tuning
    parser.add_argument('--pre_seq_len', type=int, default=16, help='')
    parser.add_argument('--prefix_projection',
                        type=bool, default=True, help='')
    parser = deepspeed.add_config_arguments(parser)
    return parser.parse_args()


def main():
    args = parse_args()

    if args.local_rank == -1:
        device = torch.device("cuda")
    else:
        torch.cuda.set_device(args.local_rank)
        device = torch.device("cuda", args.local_rank)
        deepspeed.init_distributed()
    args.global_rank = torch.distributed.get_rank()

    with open(args.ds_file, "r", encoding="utf-8") as fh:
        ds_config = json.load(fh)

    ds_config['train_micro_batch_size_per_gpu'] = args.per_device_train_batch_size
    ds_config[
        'train_batch_size'] = args.per_device_train_batch_size * torch.distributed.get_world_size() * args.gradient_accumulation_steps
    ds_config['gradient_accumulation_steps'] = args.gradient_accumulation_steps

    if args.global_rank <= 0:
        tb_write = SummaryWriter()

    set_random_seed(args.seed)
    torch.distributed.barrier()
    # load tokenizer
    tokenizer = MODE[args.mode]["tokenizer"].from_pretrained(
        args.model_name_or_path, trust_remote_code=True)

    # load model
    if args.train_type == "lora":
        print_rank_0("use LoRA", args.global_rank)
        model = MODE[args.mode]["model"].from_pretrained(
            args.model_name_or_path, trust_remote_code=True)
        lora_module_name = args.lora_module_name.split(",")
        config = LoraConfig(r=args.lora_dim,
                            lora_alpha=args.lora_alpha,
                            target_modules=lora_module_name,
                            lora_dropout=args.lora_dropout,
                            bias="none",
                            task_type="CAUSAL_LM",
                            inference_mode=False,
                            )
        model = get_peft_model(model, config)
        model.config.torch_dtype = torch.float32
    elif args.train_type == "freeze":
        model = MODE[args.mode]["model"].from_pretrained(
            args.model_name_or_path)
        freeze_module_name = args.freeze_module_name.split(",")
        for name, param in model.named_parameters():
            if not any(nd in name for nd in freeze_module_name):
                param.requires_grad = False
    elif args.train_type == "ptuning":
        config = MODE[args.mode]["config"].from_pretrained(
            args.model_name_or_path)
        config.pre_seq_len = args.pre_seq_len
        config.prefix_projection = args.prefix_projection
        model = MODE[args.mode]["model"].from_pretrained(
            args.model_name_or_path, config=config)
        for name, param in model.named_parameters():
            if not any(nd in name for nd in ["prefix_encoder"]):
                param.requires_grad = False
    elif args.train_type == "all":
        print_rank_0("train all", args.global_rank)
        model = MODE[args.mode]["model"].from_pretrained(
            args.model_name_or_path)
    else:
        raise Exception("train_type无效")
    # load data
    if args.enhancement.lower() == "true":
        print("Using enhanced data")
        args.train_path = enhance_data(args.train_path, os.path.dirname(args.train_path))

    train_dataset = MODE[args.mode]["dataset"](
        args.train_path, tokenizer, args.max_len, args.max_src_len, args.is_skip, args.retrieval_method, args.st,
        args.top_k)
    if args.local_rank == -1:
        train_sampler = RandomSampler(train_dataset)
    else:
        train_sampler = DistributedSampler(train_dataset)

    data_collator = DataCollator(tokenizer)
    train_dataloader = DataLoader(train_dataset, collate_fn=data_collator, sampler=train_sampler,
                                  batch_size=args.per_device_train_batch_size)

    # load optimizer
    ds_config["optimizer"]["params"]["lr"] = args.learning_rate
    ds_config["optimizer"]["params"]["betas"] = (0.9, 0.95)
    ds_config["optimizer"]["params"]["eps"] = 1e-8
    ds_config["optimizer"]["params"]["weight_decay"] = 0.1
    num_training_steps = args.num_train_epochs * \
                         math.ceil(len(train_dataloader) / args.gradient_accumulation_steps)
    print_rank_0("num_training_steps = {}".format(
        num_training_steps), args.global_rank)
    num_warmup_steps = int(args.warmup_ratio * num_training_steps)
    print_rank_0("num_warmup_steps = {}".format(
        num_warmup_steps), args.global_rank)
    ds_config["scheduler"]["params"]["total_num_steps"] = num_training_steps
    ds_config["scheduler"]["params"]["warmup_num_steps"] = num_warmup_steps
    ds_config["scheduler"]["params"]["warmup_max_lr"] = args.learning_rate
    ds_config["scheduler"]["params"]["warmup_min_lr"] = args.learning_rate * 0.1

    # print parameters
    for name, param in model.named_parameters():
        if param.requires_grad == True:
            print_rank_0(name, 0)
    print_trainable_parameters(model)

    # gradient_checkpointing
    if args.gradient_checkpointing:
        model.gradient_checkpointing_enable()
        if hasattr(model, "enable_input_require_grads"):
            model.enable_input_require_grads()
        else:
            def make_inputs_require_grad(module, input, output):
                output.requires_grad_(True)

            model.get_input_embeddings().register_forward_hook(make_inputs_require_grad)

    model, optimizer, _, lr_scheduler = deepspeed.initialize(model=model, args=args, config=ds_config,
                                                             dist_init_required=True)
    model.train()
    tr_loss, logging_loss, min_loss = 0.0, 0.0, 0.0
    global_step = 0
    # train
    checkpoints = []  # List to store checkpoint paths
    max_checkpoints = 1
    for epoch in range(args.num_train_epochs):
        print_rank_0("Beginning of Epoch {}/{}, Total Micro Batches {}".format(
            epoch + 1, args.num_train_epochs, len(train_dataloader)), args.global_rank)
        model.train()
        for step, batch in tqdm(enumerate(train_dataloader), total=len(train_dataloader), unit="batch"):
            batch = to_device(batch, device)
            outputs = model(**batch, use_cache=False)
            loss = outputs.loss
            tr_loss += loss.item()
            model.backward(loss)
            torch.nn.utils.clip_grad_norm_(model.parameters(), 1.0)
            model.step()
            # 更新学习率调度器
            lr_scheduler.step()
            if (step + 1) % args.gradient_accumulation_steps == 0:
                log_data = log_info(args.global_rank, epoch + (step + 1) / len(train_dataloader), step + 1,
                                    (tr_loss - logging_loss) / (args.show_loss_step * args.gradient_accumulation_steps),
                                    lr_scheduler.get_last_lr()[0])
                print_rank_0(log_data, args.global_rank)
                global_step += 1
                # write loss
                if global_step % args.show_loss_step == 0:
                    print_rank_0("step: {}-{}-{}".format(step + 1,
                                                         global_step, model.global_steps), args.global_rank)
                    if args.global_rank <= 0:
                        tb_write.add_scalar("train_loss", (tr_loss - logging_loss) /
                                            (args.show_loss_step * args.gradient_accumulation_steps), global_step)
                        logging_loss = tr_loss
                # save model
                if args.save_model_step is not None and global_step % args.save_model_step == 0:
                    checkpoint_path = save_checkpoint(model, tokenizer, epoch + 1, global_step, args, ds_config)
                    checkpoints.append(checkpoint_path)
                    checkpoints = manage_checkpoints(checkpoints, max_checkpoints)

    if args.global_rank <= 0:
        checkpoint_path = save_checkpoint(model, tokenizer, epoch + 1, global_step, args, ds_config)
        checkpoints.append(checkpoint_path)
        checkpoints = manage_checkpoints(checkpoints, max_checkpoints)
    if args.global_rank <= 0 and args.merge_lora.lower() == "true":
        merge_lora(args, checkpoint_path)

    print_rank_0("train end", args.global_rank)


if __name__ == "__main__":
    main()
