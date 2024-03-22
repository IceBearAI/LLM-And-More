import argparse
import json
from collections import defaultdict
from typing import List
import requests

import numpy as np
from pydantic import BaseModel
from sentence_transformers import SentenceTransformer
from sklearn.metrics.pairwise import cosine_similarity


class QuestionIntent(BaseModel):
    """The request model for creating an annotation."""
    input: str
    intent: str
    output: str


class MismatchedIntents(BaseModel):
    """The mismatched intents of the annotation."""
    questionPair: List[str] = []
    """The question pair of the annotation."""
    intent1: str = ""
    """The intent1 of the annotation."""
    intent2: str = ""
    """The intent2 of the annotation."""
    answer1: str = ""
    """The answer1 of the annotation."""
    answer2: str = ""
    """The answer2 of the annotation."""
    lineNumbers: List[int] = []
    """The line numbers of the annotation."""


class SimilarIntents(BaseModel):
    """The similar intents of the annotation."""
    intentPair: List[str] = []
    lineNumbers: List[int] = []
    """The intent pair of the annotation."""


class SimilarQuestionIntent(BaseModel):
    """The response model for a list of data annotations."""
    mismatchedIntents: List[MismatchedIntents] = []
    similarIntents: List[SimilarIntents] = []


class DatasetsModel:
    """Represents a model for handling datasets."""

    def __init__(self, model_path: str = "uer/sbert-base-chinese-nli"):
        """Initializes a new instance of the DatasetsModel class.

        Args:
            model_path (str): The path to the model. Defaults to "uer/sbert-base-chinese-nli".
        """
        self.model_path = model_path

        # self.device = torch.device(device)  # Use torch.device for compatibility
        self.model = SentenceTransformer(model_name_or_path=model_path)

    def analyze_similar_questions_and_intents(self, data: List[QuestionIntent],
                                              similarity_threshold: float = 0.91,
                                              intent_similarity_threshold: float = 0.86) -> SimilarQuestionIntent:
        """Analyze similar questions and intents."""

        all_query = []
        all_intents = []
        all_answers = []
        indices = []
        for i, item in enumerate(data):
            questions = item.question[0:].split(',')
            questions = [q.strip() for q in questions]
            intents = item.intent
            answers = item.output
            for q in questions:
                all_query.append(q)
            for _ in range(len(questions)):
                all_intents.append(intents)
                all_answers.append(answers)
                indices.append(i)
            i += 1
        sentence_embeddings = self.model.encode(all_query)
        cosine_score = cosine_similarity(sentence_embeddings)
        similar_indices = np.argwhere(cosine_score >= similarity_threshold)
        intent_question = defaultdict(list)

        similar_intents: List[SimilarIntents] = []
        mismatched_intents: List[MismatchedIntents] = []

        for row, col in similar_indices:
            if row < col:
                intent1 = all_intents[row]
                intent2 = all_intents[col]
                if intent1 != intent2:  # 不考虑同一个意图的情况
                    question1 = all_query[row]
                    question2 = all_query[col]
                    answer1 = all_answers[row]
                    answer2 = all_answers[col]
                    intent_question[tuple((intent1, intent2))].append({
                        "questionPair": [question1, question2],
                        "intent1": intent1,
                        "intent2": intent2,
                        "answer1": answer1,
                        "answer2": answer2,
                        "lineNumbers": [indices[row], indices[col]],
                    })
        for intent_pair, questions in intent_question.items():
            intents_ = list(intent_pair)
            sentence_embeddings = self.model.encode(intents_)
            intent_sim = cosine_similarity([sentence_embeddings[0]], [sentence_embeddings[1]])[0][0]
            if intent_sim > intent_similarity_threshold:
                for question_info in questions:
                    similar_intents.append(
                        SimilarIntents(intentPair=intents_, lineNumbers=question_info["lineNumbers"]))
            else:
                for question_info in questions:
                    mismatched_intents.append(MismatchedIntents(
                        questionPair=question_info["questionPair"],
                        intent1=question_info["intent1"],
                        intent2=question_info["intent2"],
                        answer1=question_info["answer1"],
                        answer2=question_info["answer2"],
                        lineNumbers=question_info["lineNumbers"]
                    ))

        return SimilarQuestionIntent(similarIntents=similar_intents,
                                     mismatchedIntents=mismatched_intents)


def analyze_similar(params) -> SimilarQuestionIntent:
    """Analyze similar questions and intents."""
    if not params.dataset:
        raise ValueError("dataset file path is required")

    # 判断dataset是否是url
    if params.dataset.startswith("http"):
        response = requests.get(params.dataset)
        if response.status_code != 200:
            raise ValueError("dataset file path is invalid")
        data_list = response.content.decode("utf-8").split("\n")
    else:
        with open(params.dataset, "r", encoding="utf-8") as f:
            data_list = f.readlines()

    if params.dataset_type == "faq":
        data = []
        for line in data_list:
            if not line.strip():
                continue
            item = json.loads(line)
            data.append(QuestionIntent(input=item['question'], intent=item['intent'], output=item['output']))
        return DatasetsModel(params.model_name).analyze_similar_questions_and_intents(
            data, params.similarity_threshold, params.intent_similarity_threshold)

    return SimilarQuestionIntent()


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="dataset analyze similar"
    )

    parser.add_argument("--model_name", type=str, default="uer/sbert-base-chinese-nli", help="model name")
    parser.add_argument("--device", type=str, default="mps", help="device")
    parser.add_argument("--similarity_threshold", type=float, default=0.91, help="similarity threshold")
    parser.add_argument("--intent_similarity_threshold", type=float, default=0.86, help="intent similarity threshold")
    parser.add_argument("--dataset", type=str, default="", help="dataset file path")
    parser.add_argument("--dataset_type", type=str, default="faq", help="dataset type: faq, rag, general.")
    parser.add_argument("--output_file", type=str, default="/tmp/result.json", help="output file path.")

    args = parser.parse_args()

    result = analyze_similar(params=args)
    with open(args.output_file, "w", encoding="utf-8") as f:
        f.write(result.json())
