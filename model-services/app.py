import base64
import shutil
import tempfile
import torch
import numpy as np
import os
import uuid
from dotenv import load_dotenv
from io import BytesIO

from qdrant_client import QdrantClient, models
from fastapi import FastAPI, UploadFile, File
from pydantic import BaseModel
from sklearn.metrics.pairwise import cosine_similarity
from pdf2image import convert_from_path

from embedder import model, processor

load_dotenv()

app = FastAPI(title="Agentic RAG - Model Services(Internal)")

qdrant = QdrantClient(
    url=os.getenv("QDRANT_URL"),
    api_key=os.getenv("QDRANT_API_KEY")
)

COLLECTION_NAME = "pdf_pages"

EMBEDDING_DIM = 128

def ensure_collection():
    if not qdrant.collection_exists(COLLECTION_NAME):
        qdrant.create_collection(
            collection_name=COLLECTION_NAME,
            vectors_config={
                "colqwen": models.VectorParams(
                    size=EMBEDDING_DIM,
                    distance=models.Distance.COSINE,
                    multivector_config=models.MultiVectorConfig(comparator=models.MultiVectorComparator.MAX_SIM),
                    hnsw_config=models.HnswConfigDiff(m=0),
                )
            },
        )

ensure_collection()
os.makedirs("storage/pages", exist_ok=True)

@app.post('/embed')
async def embed(file: UploadFile = File(...)):
    with tempfile.NamedTemporaryFile(delete=False, suffix=".pdf") as tmp:
        tmp.write(await file.read())
        tmp_path = tmp.name

    try:
        title = os.path.splitext(file.filename)[0]
        images = convert_from_path(tmp_path)

        points = []
        batch_size = 4
        for i in range(0, len(images), batch_size):
            batch_images = images[i:i + batch_size]
            inputs = processor.process_images(batch_images)
            inputs = {k: v.to(model.device) for k, v in inputs.items()}
            with torch.no_grad():
                embeddings = model(**inputs).embeddings
            embeddings = embeddings.cpu().float()

            for j, emb in enumerate(embeddings):
                page_idx = i + j
                if page_idx >= len(images):
                    continue

                image_path = f"storage/pages/{uuid.uuid4()}.png"
                images[page_idx].save(image_path, format="PNG")
                points.append(models.PointStruct(
                    id=str(uuid.uuid4()),
                    vector={"colqwen": emb.tolist()},
                    payload={
                        "title": title,
                        "page_number": page_idx + 1,
                        "image_path": image_path,
                    },
                ))
 
        qdrant.upsert(collection_name=COLLECTION_NAME, points=points)
        
    finally:
        os.remove(tmp_path)

    return {"title": title, "num_pages": len(images)}

class RetrieveRequest(BaseModel):
    query: str
    k: int = 3

@app.post('/retrieve')
async def do_retrieve(req: RetrieveRequest):
    collection_info = qdrant.get_collection(COLLECTION_NAME)
    if collection_info.points_count == 0:
        return {"pages": []}

    query_inputs = processor.process_queries([req.query])
    query_inputs = {k: v.to(model.device) for k, v in query_inputs.items()}

    with torch.no_grad():
        query_embedding = model(**query_inputs).embeddings.cpu().float()[0]

    search_result = qdrant.query_points(
        collection_name=COLLECTION_NAME,
        query=query_embedding.tolist(),
        using="colqwen",
        limit=req.k,
    ).points

    results = []
    for point in search_result:
        payload = point.payload
        with open(payload["image_path"], "rb") as f:
            img_b64 = base64.b64encode(f.read()).decode("utf-8")
        results.append({
            "title": payload["title"],
            "page_number": payload["page_number"],
            "image_base64": img_b64,
        })
    return {"pages": results}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8001)