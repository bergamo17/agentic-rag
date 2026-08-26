from transformers import ColQwen2ForRetrieval, ColQwen2Processor
import torch

device = "mps" if torch.backends.mps.is_available() else "cpu"
print(f"Using device: {device}")

model_name = "vidore/colqwen2-v1.0-hf"

model = ColQwen2ForRetrieval.from_pretrained(
    model_name,
    dtype=torch.float32,
    device_map=device,
).eval()

processor = ColQwen2Processor.from_pretrained(model_name)

if __name__ == "__main__":
    from load_documents import PDFs

    image_counter = 0
    for pdf_idx, pdf in enumerate(PDFs):
        print(f"Generating embeddings for {len(pdf['images'])} pages in {pdf['title']}")
        pdf['page_embeddings'] = []
        batch_size = 4
        for i in range(0, len(pdf['images']), batch_size):
            batch_images = pdf['images'][i:i+batch_size]
            inputs = processor.process_images(batch_images)
            inputs = {k: v.to(model.device) for k, v in inputs.items()}
            with torch.no_grad():
                embeddings = model(**inputs).embeddings
            embeddings = embeddings.cpu()
            embeddings = embeddings / torch.norm(embeddings, dim=1, keepdim=True)
            for j, emb in enumerate(embeddings):
                if i+j < len(pdf["images"]):
                    if 'page_embeddings' not in pdf:
                        pdf['page_embeddings'] = []
                    pdf['page_embeddings'].append(emb)
                    image_counter += 1

    print(f"Generated embeddings for {image_counter} PDF pages")