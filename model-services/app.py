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
from fastapi import FastAPI, UploadFile, File, HTTPException
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
        document_id = str(uuid.uuid4())
        images = convert_from_path(tmp_path)

        total_points = 0
        batch_size = 4
        for i in range(0, len(images), batch_size):
            batch_images = images[i:i + batch_size]
            inputs = processor.process_images(batch_images)
            inputs = {k: v.to(model.device) for k, v in inputs.items()}
            with torch.no_grad():
                embeddings = model(**inputs).embeddings
            embeddings = embeddings.cpu().float()

            batch_points = []

            for j, emb in enumerate(embeddings):
                page_idx = i + j
                if page_idx >= len(images):
                    continue

                image_path = f"storage/pages/{uuid.uuid4()}.png"
                images[page_idx].save(image_path, format="PNG")
                batch_points.append(models.PointStruct(
                    id=str(uuid.uuid4()),
                    vector={"colqwen": emb.tolist()},
                    payload={
                        "document_id": document_id,
                        "title": title,
                        "page_number": page_idx + 1,
                        "image_path": image_path,
                    },
                ))

            if batch_points:
                qdrant.upsert(collection_name=COLLECTION_NAME, points=batch_points)
                total_points += len(batch_points)
        
    finally:
        os.remove(tmp_path)

    return {"document_id": document_id, "title": title, "num_pages": len(images)}

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
            "document_id": payload["document_id"],
            "title": payload["title"],
            "page_number": payload["page_number"],
            "image_base64": img_b64,
        })
    return {"pages": results}

@app.get("/page")
async def get_page(document_id: str, page_number: int):
    scroll_result = qdrant.scroll(
        collection_name=COLLECTION_NAME,
        scroll_filter=models.Filter(
            must=[
                models.FieldCondition(key='document_id', match=models.MatchValue(value=document_id)),
                models.FieldCondition(key='page_number', match=models.MatchValue(value=page_number)),
            ]
        ),
        limit=1,
    )
    points, _ = scroll_result
    if not points:
        raise HTTPException(status_code=404, detail="This page is not found")

    payload = points[0].payload
    with open(payload['image_path'], 'rb') as f:
        img_b64 = base64.b64encode(f.read()).decode("utf-8")

    return {
        "document_id": payload["document_id"],
        "title": payload["title"],
        "page_number": payload["page_number"],
        "image_base64": img_b64,
    }

@app.get("/documents")
async def list_documents():
    all_points = []
    next_offset = None
    while True:
        points, next_offset = qdrant.scroll(
            collection_name=COLLECTION_NAME, 
            limit=200,
            offset=next_offset,
            with_payload=["document_id", "title"],
        )
        all_points.extend(points)
        if next_offset is None:
            break

    docs = {}
    for p in all_points:
        doc_id = p.payload["document_id"]
        title = p.payload["title"]
        if doc_id not in docs:
            docs[doc_id] = {"document_id": doc_id, "title": title, "num_pages": 0}
        docs[doc_id]['num_pages'] += 1

    return {"documents": list(docs.values())}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8001)