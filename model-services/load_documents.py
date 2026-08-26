import requests
from pdf2image import convert_from_path

PDFs = [
    {'title': "Attention Is All You Need", 'file': "https://arxiv.org/pdf/1706.03762"},
    {'title': "Deep Residual Learning", 'file': "https://arxiv.org/pdf/1512.03385"},
    {'title': "BERT", 'file': "https://arxiv.org/pdf/1810.04805"}
]

for pdf in PDFs:
    with open(f"{pdf['title']}.pdf", "wb") as f:
        f.write(requests.get(pdf['file']).content)

for pdf in PDFs:
    pdf['images'] = convert_from_path(f"{pdf['title']}.pdf")