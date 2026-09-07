import os
from dotenv import load_dotenv
from qdrant_client import QdrantClient, models

load_dotenv()

qdrant = QdrantClient(
    url=os.getenv("QDRANT_URL"),
    api_key=os.getenv("QDRANT_API_KEY")
)

COLLECTION_NAME = "pdf_pages"

qdrant.create_payload_index(
    collection_name=COLLECTION_NAME,
    field_name="document_id",
    field_schema="keyword",
)
qdrant.create_payload_index(
    collection_name=COLLECTION_NAME,
    field_name="page_number",
    field_schema="integer",
)

print("Index created successfully.")