import numpy as np
import torch
from sklearn.metrics.pairwise import cosine_similarity
from load_documents import PDFs
from embedder import model, processor

embeddings = np.stack([
    embedding.float().numpy()
    for pdf in PDFs for embedding in pdf['page_embeddings']
])

data = []
page_count = 0
for pdf in PDFs:
    for page_idx in range(len(pdf['images'])):
        data.append({
            "title": pdf["title"],
            "file": pdf["file"],
            "page_number": page_idx + 1,
            "image": pdf["images"][page_idx],
            "id": page_count
        })
        page_count += 1

def retrieve(query: str, k: int =3) -> list:
    query = processor.process_queries([query])
    with torch.no_grad():
        query = {k: v.to(model.device) for k, v in query.items()}
        query_embedding = model(**query).float().cpu().numpy()
    query_embedding = query_embedding / np.linalg.norm(query_embedding)
    cos_sim = cosine_similarity(query_embedding, embeddings)[0]
    idx_sorted_by_cosine_sim = np.argsort(cos_sim)[::-1]
    sorted_data = [data[i] for i in idx_sorted_by_cosine_sim]
    return sorted_data[:k]