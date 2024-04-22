import argparse
import os

import jieba
from langchain_community.document_loaders import (
    TextLoader,
    JSONLoader

)
from langchain_text_splitters import RecursiveCharacterTextSplitter
from rank_bm25 import BM25Okapi
from sentence_transformers import SentenceTransformer, util


def parse_args():
    parser = argparse.ArgumentParser()
    parser.add_argument("--train_path", type=str, default='enhancement_data.jsonl', help="Retrieved documents")
    parser.add_argument("--retrieval_method", type=str, default="bm25", choices=["bm25", "sentence_transformers"],
                        help="Method for document retrieval")
    parser.add_argument("--threshold", type=float, default=0.68, help="Threshold for similarity to refuse answering")
    parser.add_argument("--sentence_asymmetrical_path", type=str, default='shibing624/text2vec-base-chinese')

    return parser.parse_args()


def is_supported_format(filename):
    return filename.endswith((".json", ".jsonl"))


def load_json_from_dir(directory_path):
    data = []
    for filename in os.listdir(directory_path):
        if is_supported_format(filename):
            loader = JSONLoader(
                file_path=f'{directory_path}/{filename}',
                jq_schema='.document',
                text_content=False,
                json_lines=True if filename.endswith(".jsonl") else False
            )
            data.append(loader.load())
    if len(data) == 0:
        print("目录中不包含json或jsonl文件")
        return None
    merged_data = merge_loaded_data(data)
    return merged_data


def merge_loaded_data(data):
    merged_data = []
    for sublist in data:
        merged_data.extend(sublist)
    return merged_data


def load_json_from_one(filename):
    if is_supported_format(filename):
        loader = JSONLoader(
            file_path=filename,
            jq_schema='.document',
            text_content=False,
            json_lines=True if filename.endswith(".jsonl") else False
        )
        data = loader.load()
    else:
        print("文件格式错误，请传入json或jsonl文件")
        return None
    return data


def load_txt_from_dir(directory_path):
    data = []
    for filename in os.listdir(directory_path):
        if filename.endswith(".txt"):
            loader = TextLoader(f'{directory_path}/{filename}')
            data.append(loader.load())
            print(filename, "加载完毕")
    merged_data = merge_loaded_data(data)
    return merged_data


def load_text_from_one(filename):
    data = ''
    if filename.endswith(".txt"):
        loader = TextLoader(f'{filename}')
        data = loader.load()
    return data


def ret_corpus_meta(data):
    corpus = [doc.page_content for doc in data]
    meta = [doc.metadata for doc in data]
    return corpus, meta


def split(corpus, meta):
    text_splitter = RecursiveCharacterTextSplitter(
        chunk_size=256,
        chunk_overlap=20,
        length_function=len,
        separators=[
            "\n\n",
            "\n",
            " ",
            ".",
            ",",
            "\u200B",  # Zero-width space
            "\uff0c",  # Fullwidth comma
            "\u3001",  # Ideographic comma
            "\uff0e",  # Fullwidth full stop
            "\u3002",  # Ideographic full stop
            "",
        ],
        is_separator_regex=False,
    )
    if not meta:
        meta = [{'source': '无', 'seq_num': i} for i in range(len(corpus))]
    texts = text_splitter.create_documents(corpus, metadatas=meta)

    return texts


def retrieve_documents(query, data, meta, embeddings, ret_model, retrieval_method, top_k=1):
    if retrieval_method == "bm25":
        scores = embeddings.get_scores(list(jieba.cut(query)))
        top_docs = [{'doc': data[idx], 'meta': meta[idx], 'score': scores[idx] / (scores[idx] + 10)} for idx in
                    sorted(range(len(scores)), key=lambda i: scores[i], reverse=True)[:top_k]]

        # tokenized_query = list(jieba.cut(query))
        # top_docs = bm25.get_top_n(tokenized_query, documents, n=top_k)
    else:
        instruction = "为这个句子生成表示以用于检索相关文章："
        query_embedding = ret_model.encode([instruction + query], convert_to_tensor=True)
        # scores = util.pytorch_cos_sim(query_embedding, embeddings)[0]
        result = util.semantic_search(query_embedding, embeddings, top_k=top_k)
        # top_docs = [{'doc': data[idx], 'meta': meta[idx], 'score': scores[idx]} for idx in
        #             scores.argsort(descending=True)[:top_k]]
        top_docs = [{'doc': data[res['corpus_id']], 'meta': meta[res['corpus_id']], 'score': res['score']} for res in
                    result[0]]
    print("score:", top_docs[0]['score'], "来源：", top_docs[0]['meta'])

    return top_docs


def retrieval_json(train_path, retrieval_method, sentence_asymmetrical_path):
    if os.path.isdir(train_path):
        data = load_json_from_dir(train_path)
    else:
        data = load_json_from_one(train_path)
    corpus1, meta1 = ret_corpus_meta(data)
    document = split(corpus1, meta1)
    corpus2, meta2 = ret_corpus_meta(document)
    ret_model = None
    if retrieval_method == "bm25":
        tokenized_corpus = [list(jieba.cut(doc)) for doc in corpus2]
        embeddings = BM25Okapi(tokenized_corpus)
    else:
        ret_model = SentenceTransformer(sentence_asymmetrical_path)
        embeddings = ret_model.encode(corpus2, convert_to_tensor=True)

    return corpus2, meta2, ret_model, embeddings


if __name__ == '__main__':
    args = parse_args()
    corpus, meta, ret_model, embeddings = retrieval_json(args.train_path, args.retrieval_method,
                                                         args.sentence_asymmetrical_path)
    doc = retrieve_documents("你是谁？", corpus, meta, embeddings, ret_model, args.retrieval_method)
    doc = doc[0]['doc']
    print(doc)
