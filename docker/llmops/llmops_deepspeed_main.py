#!/usr/bin/env python
# Copyright (c) Microsoft Corporation.
# SPDX-License-Identifier: Apache-2.0

import argparse
import math
import os
# DeepSpeed Team
import shutil
import subprocess
import sys
import time
import warnings

import torch
from torch.utils.data import DataLoader, RandomSampler, SequentialSampler
from torch.utils.data.distributed import DistributedSampler
from transformers import (
    AutoModelForCausalLM,
    SchedulerType,
    default_data_collator,
    get_scheduler,
)

import deepspeed
from deepspeed.ops.adam import DeepSpeedCPUAdam, FusedAdam
from utils.data.data_utils import create_prompt_dataset
from utils.ds_utils import get_train_ds_config
from utils.model.model_utils import create_hf_model
from utils.module.lora import convert_linear_layer_to_lora, convert_lora_to_linear_layer, only_optimize_lora_parameters, \
    make_model_gradient_checkpointing_compatible
from utils.utils import print_rank_0, to_device, save_hf_format, set_random_seed, get_all_reduce_mean, \
    get_optimizer_grouped_parameters, save_zero_three_model, load_hf_tokenizer

sys.path.append(
    os.path.abspath(os.path.join(os.path.dirname(__file__), os.path.pardir)))

warnings.filterwarnings('ignore')


def log_info(Rank, epoch, step, train_loss, learning_rate, eval_loss=None):
    # 创建基本的日志字典
    log_dict = {
        'rank': Rank,
        'epoch': epoch,
        'step': step,
        'loss': train_loss,
        'learning_rate': learning_rate,
    }

    # 如果提供了eval_loss，则添加到字典中
    if eval_loss is not None:
        log_dict['eval_loss'] = eval_loss

    return log_dict


def parse_args():
    parser = argparse.ArgumentParser(
        description="Finetune a transformers model on a causal language modeling task")
    parser.add_argument('--data_path',
                        nargs='*',
                        default="formatted_JCQA",
                        help='Path to the training dataset. Accepted format:'
                             '1) a single data path, 2) multiple datasets in the'
                             'form: dataset1-path dataset2-path ...')
    parser.add_argument('--data_split',
                        type=str,
                        default='10,0,0',
                        help='Comma-separated list of proportions for training'
                             'phase 1, 2, and 3 data. For example the split `6,2,2`'
                             'will use 60%% of data for phase 1, 20%% for phase 2'
                             'and 20%% for phase 3.')
    parser.add_argument(
        '--sft_only_data_path',
        nargs='*',
        default=[],
        help='Path to the dataset for only using in SFT phase.')
    parser.add_argument(
        '--data_output_path',
        type=str,
        default='/home/calf/ssd/data_files/',
        help='Where to store the data-related files such as shuffle index. This needs to be on a local storage of a node (not on a shared storage)'
    )
    parser.add_argument(
        "--model_name_or_path",
        type=str,
        default="/home/calf/ssd/outputs/llama-2-sft-new-600w4tasks",
        help="Path to pretrained model or model identifier from huggingface.co/models.",
    )
    parser.add_argument(
        "--per_device_train_batch_size",
        type=int,
        default=16,
        help="Batch size (per device) for the training dataloader.",
    )
    parser.add_argument(
        "--per_device_eval_batch_size",
        type=int,
        default=16,
        help="Batch size (per device) for the evaluation dataloader.",
    )
    parser.add_argument(
        "--max_seq_len",
        type=int,
        default=512,
        help="The maximum sequence length.",
    )
    parser.add_argument(
        "--learning_rate",
        type=float,
        default=1e-5,
        help="Initial learning rate (after the potential warmup period) to use.",
    )
    parser.add_argument("--weight_decay",
                        type=float,
                        default=0.,
                        help="Weight decay to use.")
    parser.add_argument("--num_train_epochs",
                        type=int,
                        default=1,
                        help="Total number of training epochs to perform.")
    parser.add_argument(
        "--gradient_accumulation_steps",
        type=int,
        default=1,
        help="Number of updates steps to accumulate before performing a backward/update pass.",
    )
    parser.add_argument(
        "--lr_scheduler_type",
        type=SchedulerType,
        default="cosine",
        help="The scheduler type to use.",
        choices=[
            "linear", "cosine", "cosine_with_restarts", "polynomial",
            "constant", "constant_with_warmup"
        ],
    )
    parser.add_argument(
        "--num_warmup_steps",
        type=int,
        default=0,
        help="Number of steps for the warmup in the lr scheduler.")
    parser.add_argument("--output_dir",
                        type=str,
                        default='/home/calf/ssd/outputs/llama-2-rope-experimental',
                        help="Where to store the model.")
    parser.add_argument("--seed",
                        type=int,
                        default=1234,
                        help="A seed for reproducible training.")
    parser.add_argument("--local_rank",
                        type=int,
                        default=-1,
                        help="local_rank for distributed training on gpus")
    parser.add_argument('--gradient_checkpointing',
                        action='store_true',
                        help='Enable HF gradient checkpointing for model.')
    parser.add_argument('--disable_dropout',
                        action='store_true',
                        help='Disable the dropout of the model.')
    # deepspeed features
    parser.add_argument('--offload',
                        action='store_true',
                        help='Enable ZeRO Offload techniques.')
    parser.add_argument(
        '--zero_stage',
        type=int,
        default=2,
        help='ZeRO optimization stage for Actor model (and clones).')
    # ROPE
    parser.add_argument(
        "--rope_pi",
        action='store_true',
        help="use rope_pi"
    )
    parser.add_argument("--rope_ratio",
                        type=int,
                        default=8,
                        help="If > 1, extend context window.")
    # LoRA for efficient training setting
    parser.add_argument("--train_type", type=str, default="lora", help="")
    parser.add_argument("--lora_dim",
                        type=int,
                        default=0,
                        help="If > 0, use LoRA for efficient training.")
    parser.add_argument("--lora_module_name",
                        type=str,
                        default="decoder.layers.",
                        help="The scope of LoRA.")
    parser.add_argument('--only_optimize_lora',
                        action='store_true',
                        help='Only optimize the LoRA parameters.')
    parser.add_argument(
        "--lora_learning_rate",
        type=float,
        default=2e-5,
        help="Initial LoRA learning rate (after the potential warmup period) to use."
    )
    # save per steps
    parser.add_argument('--save_per_steps', type=int,
                        help='save per x steps', default=200)
    parser.add_argument('--save_total_limit', type=int, default=1)
    # start from step
    parser.add_argument('--start_from_step', type=int,
                        help='skip first x steps', default=-1)
    # Tensorboard logging
    parser.add_argument('--enable_tensorboard',
                        action='store_true',
                        help='Enable tensorboard logging')
    parser.add_argument('--tensorboard_path',
                        type=str,
                        default="step1_tensorboard")
    parser.add_argument('--tensorboard_port',
                        type=int,
                        default=6006,
                        help='Port for tensorboard')
    parser = deepspeed.add_config_arguments(parser)
    args = parser.parse_args()

    return args


def start_tensorboard(args):
    # 启动TensorBoard服务，并返回URL
    tensorboard_command = f"tensorboard --logdir={args.tensorboard_path} --port={args.tensorboard_port} --bind_all"
    print_rank_0(f"Starting TensorBoard with command: {tensorboard_command}")
    subprocess.Popen(tensorboard_command, shell=True)
    print_rank_0("==" * 30)
    print_rank_0("==" * 30)
    print_rank_0("==" * 30)
    tensorboard_url = f"http://localhost:{args.tensorboard_port}/"
    print_rank_0(f"TensorBoard URL: {tensorboard_url}")
    print_rank_0("==" * 30)
    print_rank_0("==" * 30)
    print_rank_0("==" * 30)
    return tensorboard_url


def notexists_mkdir(args):
    # 不存在即创建输出目录
    if not os.path.exists(args.output_dir):
        os.makedirs(args.output_dir)
    # 不存在即创建tensorboard目录
    try:
        os.rmdir(args.tensorboard_path)
    except:
        print("tensorboard_path not exists.")
    try:
        os.makedirs(args.tensorboard_path)
    except:
        print("tensorboard_path exists.")
    # 不存在即创建data_output目录
    if not os.path.exists(args.data_output_path):
        os.makedirs(args.data_output_path)


def main():
    args = parse_args()
    # 不存在则创建文件夹函数
    notexists_mkdir(args)

    if args.local_rank == -1:
        device = torch.device("cuda")
    else:
        torch.cuda.set_device(args.local_rank)
        device = torch.device("cuda", args.local_rank)
        # Initializes the distributed backend which will take care of sychronizing nodes/GPUs
        # torch.distributed.init_process_group(backend='nccl')
        deepspeed.init_distributed()

    args.global_rank = torch.distributed.get_rank()

    ds_config = get_train_ds_config(offload=args.offload,
                                    stage=args.zero_stage,
                                    enable_tensorboard=args.enable_tensorboard,
                                    tb_path=args.tensorboard_path,
                                    tb_name="step1_model")
    ds_config[
        'train_micro_batch_size_per_gpu'] = args.per_device_train_batch_size
    ds_config[
        'train_batch_size'] = args.per_device_train_batch_size * torch.distributed.get_world_size(
    ) * args.gradient_accumulation_steps

    # If passed along, set the training seed now.
    set_random_seed(args.seed)

    torch.distributed.barrier()
    print_rank_0('finished setting torch.', args.global_rank)

    tokenizer = load_hf_tokenizer(args.model_name_or_path, fast_tokenizer=True)
    model = create_hf_model(AutoModelForCausalLM,
                            args.model_name_or_path,
                            tokenizer,
                            args.rope_pi,
                            args.rope_ratio,
                            ds_config,
                            disable_dropout=args.disable_dropout)

    if args.train_type == "lora":
        lora_module_name = args.lora_module_name.split(",")

        model = convert_linear_layer_to_lora(model, lora_module_name,
                                             args.lora_dim)
        print_rank_0('using LoRA.', args.global_rank)
        if args.only_optimize_lora:
            model = only_optimize_lora_parameters(model)
            model = make_model_gradient_checkpointing_compatible(model)
            print_rank_0('only optimizing LoRA parameters.', args.global_rank)
    else:
        print_rank_0('not using LoRA.', args.global_rank)
    # model = model.to(device)
    # Prepare the data
    print_rank_0('start preparing dataset...', args.global_rank)
    train_phase = 1
    train_dataset, eval_dataset = create_prompt_dataset(
        args.local_rank,
        args.data_path,
        args.data_split,
        args.data_output_path,
        train_phase,
        args.seed,
        tokenizer,
        args.max_seq_len,
        sft_only_data_path=args.sft_only_data_path)
    print_rank_0('dataset prepared.', args.global_rank)
    # DataLoaders creation:
    print_rank_0('preparing dataloader...', args.global_rank)
    if args.local_rank == -1:
        train_sampler = RandomSampler(train_dataset)
        eval_sampler = SequentialSampler(eval_dataset)
    else:
        train_sampler = DistributedSampler(train_dataset)
        eval_sampler = DistributedSampler(eval_dataset)
    train_dataloader = DataLoader(train_dataset,
                                  collate_fn=default_data_collator,
                                  sampler=train_sampler,
                                  batch_size=args.per_device_train_batch_size)
    eval_dataloader = DataLoader(eval_dataset,
                                 collate_fn=default_data_collator,
                                 sampler=eval_sampler,
                                 batch_size=args.per_device_eval_batch_size)

    def evaluation_loss_perplexity(model, eval_dataloader, device):
        model.eval()
        losses = 0
        for step, batch in enumerate(eval_dataloader):
            batch = to_device(batch, device)
            with torch.no_grad():
                outputs = model(**batch)

            loss = outputs.loss
            losses += loss.float()
        average_loss = losses / (step + 1)  # 计算平均损失
        try:
            perplexity = torch.exp(average_loss)
        except OverflowError:
            perplexity = float("inf")
        try:
            perplexity = get_all_reduce_mean(perplexity).item()
        except:
            pass

        # 转换平均损失和困惑度为Python原生数值，方便打印
        average_loss = average_loss.item()
        perplexity = perplexity.item() if not isinstance(
            perplexity, float) else perplexity
        return average_loss, perplexity

    print_rank_0('finished defining evaluation function.', args.global_rank)
    # Split weights in two groups, one with weight decay and the other not.
    optimizer_grouped_parameters = get_optimizer_grouped_parameters(
        model, args.weight_decay, args.lora_learning_rate)
    print_rank_0('finished weights splitting.', args.global_rank)
    AdamOptimizer = DeepSpeedCPUAdam if args.offload else FusedAdam
    optimizer = AdamOptimizer(optimizer_grouped_parameters,
                              lr=args.learning_rate,
                              betas=(0.9, 0.95))
    print_rank_0('finished optimizer init.', args.global_rank)
    num_update_steps_per_epoch = math.ceil(
        len(train_dataloader) / args.gradient_accumulation_steps)
    lr_scheduler = get_scheduler(
        name=args.lr_scheduler_type,
        optimizer=optimizer,
        num_warmup_steps=args.num_warmup_steps,
        num_training_steps=args.num_train_epochs * num_update_steps_per_epoch,
    )
    print_rank_0('finished scheduler init.', args.global_rank)
    model, optimizer, _, lr_scheduler = deepspeed.initialize(
        model=model,
        optimizer=optimizer,
        args=args,
        config=ds_config,
        lr_scheduler=lr_scheduler,
        dist_init_required=True)
    print_rank_0('finished deepspeed init.', args.global_rank)
    if args.gradient_checkpointing:
        model.gradient_checkpointing_enable()
    print_rank_0('start training...', args.global_rank)
    # Train!
    print_rank_0("***** Running training *****", args.global_rank)
    rounds = len(train_dataloader)
    saved_checkpoints = []
    checkpoints_path = os.path.join(args.output_dir, f"checkpoint")
    saved_models = []
    best_perplexity_eval = 9999999.9
    best_epoch = args.num_train_epochs
    for epoch in range(args.num_train_epochs):
        print_rank_0(
            f"Beginning of Epoch {epoch + 1}/{args.num_train_epochs}, Total Micro Batches {rounds}",
            args.global_rank)
        model.train()
        total_train_loss = 0
        num_train_steps = 0
        last_learning_rate = 0 

        for step, batch in enumerate(train_dataloader):
            if step < args.start_from_step:
                print_rank_0(f'skipping {step + 1}-th step of {rounds}.', args.global_rank)
                continue
            print_rank_0(f'training {step + 1}-th step of {rounds}.', args.global_rank)
            start = time.time()
            batch = to_device(batch, device)
            outputs = model(**batch, use_cache=False)
            loss = outputs.loss
            learning_rate = model.optimizer.param_groups[0]['lr']
            progress = epoch + (step + 1) / rounds
            log_data = log_info(args.global_rank, progress, step + 1, loss.item(), learning_rate)
            print_rank_0(log_data, args.global_rank)
            model.backward(loss)
            total_train_loss += loss.item()
            num_train_steps += 1
            model.step()
            last_learning_rate = model.optimizer.param_groups[0]['lr']
            end = time.time()
            print_rank_0(f'finished step {step + 1}, used {end - start} seconds.', args.global_rank)
            # 保存检测点
            if step > 0 and (step + 1) % args.save_per_steps == 0:
                print_rank_0(f'saving checkpoint at step {step + 1}...', args.global_rank)
                # if args.train_type == "lora":
                #     save_dir = os.path.join(checkpoints_path, f"lora-step-{step+1}")
                #     if args.global_rank <= 0:
                #         model.save_pretrained(
                #             save_dir,
                #             state_dict=model._zero3_consolidated_16bit_state_dict() if args.zero_stage == 3 else None)
                #         tokenizer.save_pretrained(save_dir)
                #     saved_checkpoints.append(save_dir)
                # else:
                model.save_checkpoint(checkpoints_path, client_state={
                    "step": step + 1,
                    "epoch": epoch,
                    "learning_rate": learning_rate
                })

                if args.global_rank == 0:
                    saved_checkpoints.append(os.path.join(checkpoints_path, f"global_step{model.global_steps}"))
                    if len(saved_checkpoints) > args.save_total_limit:
                        # 删除最旧的检测点
                        old_checkpoint_path = saved_checkpoints.pop(0)
                        if os.path.exists(old_checkpoint_path):
                            shutil.rmtree(old_checkpoint_path)
        # 每个epoch保存一次模型
        base_model = convert_lora_to_linear_layer(model)
        sub_folder = f"epoch-{epoch+1}"
        output_dir = os.path.join(args.output_dir, sub_folder)
        if args.global_rank == 0:
            save_hf_format(base_model, tokenizer, args, sub_folder=sub_folder)
            saved_models.append(output_dir)
        if args.zero_stage == 3:
            save_zero_three_model(base_model,
                                  args.global_rank,
                                  output_dir,
                                  zero_stage=args.zero_stage)
        print_rank_0(f"Saving to {output_dir}", args.global_rank)

        if args.global_rank == 0:
            while len(saved_models) > 2:
                old_model_dir = saved_models.pop(0)
                if os.path.exists(old_model_dir) and f"epoch-{best_epoch}" not in old_model_dir:
                    shutil.rmtree(old_model_dir)
                    print(f"Deleted {old_model_dir}")
                else:
                    saved_models.append(old_model_dir)

        # 计算平均训练损失
        average_train_loss = total_train_loss / num_train_steps
        # Evaluate perplexity on the validation set.
        print_rank_0(
            f"***** Evaluating perplexity, Epoch {epoch + 1}/{args.num_train_epochs} *****",
            args.global_rank)
        try:
            print_rank_0("evaluating on sft dataset...", args.global_rank)
            eval_loss, perplexity = evaluation_loss_perplexity(
                model, eval_dataloader, device)
            # 分别记录训练和评估的损失
            # 记录日志
            log_data = log_info(
                args.global_rank,
                epoch + 1,
                num_train_steps,  # 传递最后一步的步数
                average_train_loss,
                last_learning_rate,  # 传递最后一步的学习率
                eval_loss
            )
            print_rank_0(log_data, args.global_rank)
            print_rank_0(
                f"Epoch {epoch + 1} finished, Training Loss: {average_train_loss}, Evaluation Loss: {eval_loss}, Perplexity: {perplexity}",
                args.global_rank)
            perplexity_eval=(perplexity+eval_loss)/2
            print(f"perplexity_eval_{args.global_rank}:", perplexity_eval)
            print(f"eval_{args.global_rank}:",eval_loss)
            print(f"perplexity_:{args.global_rank}:", perplexity)
        except Exception as e:
            print_rank_0(f"Evaluation failed: {e}", args.global_rank)
        if not perplexity_eval:
            perplexity_eval= best_perplexity_eval - 0.001
            best_epoch = epoch + 1
            print_rank_0(f"Best model found at epoch:{best_epoch}", args.global_rank)
        elif perplexity_eval <= best_perplexity_eval:
            best_perplexity_eval = perplexity_eval
            best_epoch = epoch + 1
            print_rank_0(f"Best model found at epoch:{best_epoch}")
        model.tput_timer.update_epoch_count()
    best_model_dir = os.path.join(args.output_dir, f"epoch-{best_epoch}")
    if args.global_rank == 0:
        shutil.copytree(best_model_dir, args.output_dir, dirs_exist_ok=True)
        # 删除之前保存的模型
        for model_dir in saved_models:
            if os.path.exists(model_dir):
                shutil.rmtree(model_dir)
    print_rank_0('finished training.', args.global_rank)


if __name__ == "__main__":
    print_rank_0('start running...')
    main()
    print_rank_0('finished running.')
