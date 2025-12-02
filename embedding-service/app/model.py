from sentence_transformers import SentenceTransformer
import os

MODEL_NAME = os.getenv("SENTENCE_TRANSFORMER_MODEL", "sentence-transformers/all-MiniLM-L6-v2")

def get_model():
    return SentenceTransformer(MODEL_NAME)
