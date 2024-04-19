import argparse
import json

import jieba
from rank_bm25 import BM25Okapi
from sentence_transformers import SentenceTransformer, util


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
    parser.add_argument("--doc_path", type=str, default='./data/enhancement_data.json', help="Retrieved documents")
    parser.add_argument("--rag_history_path", type=str, default='./data/rag_history.txt', help="Dialogue with History")
    parser.add_argument("--retrieval_method", type=str, default="sentence_transformers", choices=["bm25", "sentence_transformers"],
                        help="Method for document retrieval")
    parser.add_argument("--threshold", type=float, default=0.69, help="Threshold for similarity to refuse answering")
    parser.add_argument("--sentence_asymmetrical_path", type=str, default='./text2vec-base-chinese')
    parser.add_argument("--sentence_unsymmetrical_path", type=str,
                        default='./text2vec-base-chinese')
    return parser.parse_args()


def should_refuse_to_answer(retrieved_docs_scores, threshold):
    return retrieved_docs_scores < threshold


def load_doc1(data):
    documents = [{'question': doc['question'], 'document': doc['document']} for doc in data]
    return documents


def load_doc2(data):
    documents = [{'document': doc['document']} for doc in data]
    return documents


def retrieve_documents_bm25(query, data, bm25, top_k=1):
    scores = bm25.get_scores(list(jieba.cut(query)))

    top_docs = [data[idx] for idx in sorted(range(len(scores)), key=lambda i: scores[i], reverse=True)[:top_k]]

    # tokenized_query = list(jieba.cut(query))
    # top_docs = bm25.get_top_n(tokenized_query, documents, n=top_k)
    return top_docs, max(scores)


def retrieve_documents_st(query, data, corpus_embeddings, top_k=1):
    instruction = "为这个句子生成表示以用于检索相关文章："

    query_embedding = ret_model.encode([instruction + query], convert_to_tensor=True)
    scores = util.pytorch_cos_sim(query_embedding, corpus_embeddings)[0]
    top_docs = [data[idx] for idx in scores.argsort(descending=True)[:top_k]]
    return top_docs, max(scores)


if __name__ == '__main__':
    args = parse_args()

    data = []
    with open(args.doc_path, 'r', encoding='utf-8') as f:
        for line in f:
            try:
                json_data = json.loads(line)
                data.append(json_data)
            except json.JSONDecodeError as e:
                print(f"Error decoding JSON: {e}")

    if args.retrieval_method == "bm25":
        if 'question' in data[0]:
            documents = load_doc1(data)
            tokenized_corpus = [list(jieba.cut(doc['question'])) for doc in documents]
        else:
            documents = load_doc2(data)
            tokenized_corpus = [list(jieba.cut(doc['document'])) for doc in documents]
        bm25 = BM25Okapi(tokenized_corpus)
    else:
        if 'question' in data[0]:  # 对称语义搜索
            # model = SentenceTransformer('./sbert-base-chinese-nli')
            # model = SentenceTransformer('sentence-transformers/all-MiniLM-L6-v2')
            ret_model = SentenceTransformer(args.sentence_asymmetrical_path)
            documents = load_doc1(data)
            corpus = [doc['question'] for doc in documents]
        else:  # 非对称语义搜索
            # model = SentenceTransformer('sentence-transformers/multi-qa-MiniLM-L6-cos-v1')
            # model = SentenceTransformer(args.sentence_unsymmetrical_path)
            # nghuyong/ernie-3.0-base-zh
            ret_model = SentenceTransformer(args.sentence_unsymmetrical_path)
            documents = load_doc2(data)
            corpus = [doc['document'] for doc in documents]

        corpus_embeddings = ret_model.encode(corpus, convert_to_tensor=True)

    instruction = "你是一个专业的客服机器人。你需要使用提供的背景知识来回答问题，请严格根据背景知识的内容回答，对于没有背景知识的信息或与问题不匹配的背景知识，直接回答“不知道”："

    while True:
        question = input("('exit'退出): ")
        if question.lower() == 'exit':
            break
        try:
            retrieved_docs, scores = retrieve_documents_bm25(question, documents, bm25, top_k=3)
        except:
            retrieved_docs, scores = retrieve_documents_st(question, documents, corpus_embeddings, top_k=3)
        print(scores)
        docs = retrieved_docs[0]['document']

        if should_refuse_to_answer(scores, threshold=args.threshold):
            answer = '抱歉，我不知道'
        else:
            answer = docs

        response_data = {
            "code": 0,
            "msg": "Success",
            "data": {
                "question": question,
                "docs": docs,
                "answer": answer
            }
        }
        print(json.dumps(response_data, ensure_ascii=False, indent=4))
        answer_str = json.dumps({'document': docs, "question": question, 'answer': answer},
                                ensure_ascii=False) + '\n'
        with open(args.rag_history_path, 'a', encoding='utf-8') as f:
            f.write(answer_str)
