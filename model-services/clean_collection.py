from qdrant_client import QdrantClient
from dotenv import load_dotenv
from app import COLLECTION_NAME
import os

load_dotenv()

qdrant = QdrantClient(
    url=os.getenv("QDRANT_URL"),
    api_key=os.getenv("QDRANT_API_KEY")
)

qdrant.delete_collection(COLLECTION_NAME)
print(f"Colecction '{COLLECTION_NAME}' has been deleted.")