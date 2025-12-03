from light_embed import TextEmbedding
import os


MODEL_NAME = os.getenv("SENTENCE_TRANSFORMER_MODEL", "sentence-transformers/all-MiniLM-L6-v2")

def get_model():
    model = TextEmbedding(MODEL_NAME)
    return model

