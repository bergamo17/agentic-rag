import base64
import shutil
import tempfile
import torch
import numpy as np
import os
from io import BytesIO

from fastapi import FastAPI, UploadFile, File
from pydantic import BaseModel
from sklearn.metrics.pairwise import cosine_similarity
from pdf2image import convert_from_path

from embedder import model, processor

app = FastAPI(title="Agentic RAG - Model Services(Internal)")

PDFs = []

def _rebuild_index():
    global _embeddings, _data
    all_embeddings = [
        emb.float().numpy() for pdf in PDFs for emb in pdf["page_embeddings"]
    ]
    _embeddings = np.stack(all_embeddings) if all_embeddings else None
    _data = []
    for pdf in PDFs:
        for page_idx in range(len(pdf["images"])):
            _data.append({
                "title": pdf["title"],
                "page_number": page_idx + 1,
                "image": pdf["images"][page_idx],
            })

_embeddings = None
_data = None
_rebuild_index()

@app.post('/embed')
async def embed(file: UploadFile = File(...)):
    with tempfile.NamedTemporaryFile(delete=False, suffix=".pdf") as tmp:
        tmp.write(await file.read())
        tmp_path = tmp.name

    try:
        title = os.path.splitext(file.filename)[0]
        images = convert_from_path(tmp_path)

        page_embeddings = []
        batch_size = 4
        for i in range(0, len(images), batch_size):
            batch_images = images[i:i + batch_size]
            inputs = processor.process_images(batch_images)
            inputs = {k: v.to(model.device) for k, v in inputs.items()}
            with torch.no_grad():
                embeddings = model(**inputs).embeddings
            embeddings = embeddings.cpu()
            embeddings = embeddings / torch.norm(embeddings, dim=1, keepdim=True)
            page_embeddings.extend(embeddings)
 
        PDFs.append({
            "title": title,
            "file": tmp_path,
            "images": images,
            "page_embeddings": page_embeddings,
        })

        _rebuild_index()
    finally:
        os.remove(tmp_path)

    return {"title": title, "num_pages": len(images)}

class RetrieveRequest(BaseModel):
    query: str
    k: int = 3

@app.post('/retrieve')
async def do_retrieve(req: RetrieveRequest):
    if _embeddings is None or len(_data) == 0:
        return {"pages": []}

    query_inputs = processor.process_queries([req.query])
    query_inputs = {k: v.to(model.device) for k,v in query_inputs.items()}

    with torch.no_grad():
        query_embedding = model(**query_inputs).embeddings().float().cpu().numpy()

    query_embedding = query_embedding / np.linalg.norm(query_embedding)

    cos_sim = cosine_similarity(query_embedding, _embeddings)[0]
    top_idx = np.argsort(cos_sim)[::-1][:req.k]

    results = []
    for i in top_idx:
        page = _data[i]
        buffer = BytesIO()
        page['image'].save(buffer, format="PNG")
        img_b64 = base64.b64encode(buffer.getvalue()).decode("utf-8")
        results.append({
            "title": page['title'],
            "page_number": page['page_number'],
            "image_base64": img_b64,
        })
    return {"pages": results}

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="0.0.0.0", port=8001)