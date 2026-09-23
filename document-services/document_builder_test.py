import json
from document_builder import build_document

with open("./docs/sample.json") as f:
    data = json.load(f)

build_document(
    theme_name=data["theme"],
    title=data["title"],
    sections=data["sections"],
    output_path="report.docx"
)