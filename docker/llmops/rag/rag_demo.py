import argparse
import json

import jieba
import torch
from rank_bm25 import BM25Okapi
from sentence_transformers import SentenceTransformer

from retrieval_side import retrieval_json, retrieve_documents, split, ret_corpus_meta
from model import MODE


def parse_args():
    parser = argparse.ArgumentParser()
    # Model
    parser.add_argument("--device", type=str, default="0", help="")
    parser.add_argument("--mode", type=str, default="glm3", help="")
    parser.add_argument("--model_path", type=str, default="./output_model/", help="")
    parser.add_argument("--max_length", type=int, default=8192, help="")
    parser.add_argument("--do_sample", type=bool, default=True, help="")
    parser.add_argument("--top_p", type=float, default=0.8, help="")
    parser.add_argument("--temperature", type=float, default=0.0, help="")
    parser.add_argument("--doc_path", type=str, default='./data/output.jsonl', help="Retrieved documents")
    parser.add_argument("--rag_history_path", type=str, default='./data/rag_history.txt', help="Dialogue with History")
    parser.add_argument("--retrieval_method", type=str, default="sentence_transformer", choices=["bm25", "sentence_transformers"],
                        help="Method for document retrieval")
    parser.add_argument("--threshold", type=float, default=0.68, help="Threshold for similarity to refuse answering")
    parser.add_argument("--top_k", type=int, default=1)
    # parser.add_argument("--sentence_asymmetrical_path", type=str, default='shibing624/text2vec-base-chinese')
    parser.add_argument("--sentence_unsymmetrical_path", type=str,
                        default='BAAI/bge-base-zh-v1.5')
    return parser.parse_args()


def should_refuse_to_answer(retrieved_docs_scores, threshold):
    return retrieved_docs_scores < threshold


def load_doc(data):
    meta = [doc['question'] for doc in data]
    corpus = [doc['document'] for doc in data]
    return corpus,meta


def predict(instruction, document, question, model, tokenizer, args, mode):
    if "glm" in mode:
        result, _ = model.chat(tokenizer, instruction + document + question, max_length=args.max_length,
                               do_sample=args.do_sample,
                               top_p=args.top_p, temperature=args.temperature)
    else:
        input_ids = []
        max_doc_len = args.max_length - len(instruction + question)
        if len(document) > max_doc_len:
            document = document[:max_doc_len]
        value = instruction + document + question
        value_ids = tokenizer.encode(value)
        input_ids += value_ids
        input_ids.append(tokenizer.eos_token_id)
        input_ids = torch.tensor([input_ids], dtype=torch.long).to(
            torch.device("cuda:{}".format(args.device)))
        outputs = model.generate(input_ids, max_length=args.max_length, do_sample=args.do_sample,
                                 top_p=args.top_p, temperature=args.temperature)
        result = tokenizer.decode(
            outputs[0][len(input_ids[0]):], skip_special_tokens=True)

    return result


if __name__ == '__main__':
    args = parse_args()
    model = MODE[args.mode]["model"].from_pretrained(args.model_path, device_map="cuda:{}".format(
        args.device), torch_dtype=torch.float16, trust_remote_code=True)
    tokenizer = MODE[args.mode]["tokenizer"].from_pretrained(
        args.model_path, trust_remote_code=True)
    print('finished model and tokenizer loading')

    # data = []
    # with open(args.doc_path, 'r', encoding='utf-8') as f:
    #     for line in f:
    #         try:
    #             json_data = json.loads(line)
    #             data.append(json_data)
    #         except json.JSONDecodeError as e:
    #             print(f"Error decoding JSON: {e}")

    ret_model = None
    meta= None
    # if 'question' in data[0]:
    #     if args.retrieval_method == "bm25":
    #         corpus1, meta1= load_doc(data)
    #         document1 = split(corpus1, meta1)
    #         corpus, meta = ret_corpus_meta(document1)
    #         tokenized_corpus = [list(jieba.cut(doc)) for doc in corpus]
    #         embeddings = BM25Okapi(tokenized_corpus)
    #     else:
    #         ret_model = SentenceTransformer(args.sentence_asymmetrical_path)
    #         corpus1, meta1= load_doc(data)
    #         document1 = split(corpus1, meta1)
    #         corpus, meta = ret_corpus_meta(document1)
    #         embeddings = ret_model.encode(corpus, convert_to_tensor=True)
    # else:
    corpus, meta, ret_model, embeddings = retrieval_json(args.doc_path, args.retrieval_method,
                                                         args.sentence_unsymmetrical_path)

    instruction = "你是一个专业的客服机器人。你需要使用提供的背景知识来回答问题，请严格根据背景知识的内容回答，对于没有背景知识的信息或与问题不匹配的背景知识，直接回答“抱歉，我不知道”："

    while True:
        question = input("('exit'退出): ")
        if question.lower() == 'exit':
            break

        doc = retrieve_documents(question, corpus, meta, embeddings, ret_model, args.retrieval_method,args.top_k)
        if args.top_k > 1:
            docs = [f"片段{i}:{s['doc']}" for i, s in enumerate(doc)]
            document = "\n".join(docs)
        else:
            document = doc[0]['doc']
            # meta=doc[0]['meta']
            # document=f"{document}\n数据来源：{meta}"
        scores = doc[0]['score']
        document_context = f'背景知识：{document}\n'
        questions = f'问题：{question}\n'
        if should_refuse_to_answer(scores, threshold=args.threshold):
            answer = '抱歉，我不知道'
        else:
            answer = predict(instruction, document_context, questions, model, tokenizer, args, args.mode)
        response_data = {
            "code": 0,
            "msg": "Success",
            "data": {
                "question": question,
                "docs": document_context,
                "answer": answer
            }
        }
        print(json.dumps(response_data, ensure_ascii=False, indent=4))
        answer_str = json.dumps({'document': document, "question": question, 'answer': answer},
                                ensure_ascii=False) + '\n'
        with open(args.rag_history_path, 'a', encoding='utf-8') as f:
            f.write(answer_str)
