from sentence_transformers import SentenceTransformer
from fastapi import FastAPI
from pydantic import BaseModel


class Message(BaseModel):
    book_title: str

model = SentenceTransformer('Mihaiii/gte-micro-v2')


app = FastAPI()

@app.get("/")
async def root():
    return {"message": {"This is a sentence model transformer API"}}


@app.post("/model/")
async def create_embedding(message: Message):
    return {"embedding": model.encode(message.book_title).tolist()}