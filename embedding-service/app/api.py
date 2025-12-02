from fastapi import APIRouter, Depends
from .model import get_model
from .schemas import EmbedRequest

router = APIRouter()

@router.post("/embed")
def embed(req: EmbedRequest, model = Depends(get_model)):
    vector = model.encode(req.text).tolist()
    return {"embedding": vector}
