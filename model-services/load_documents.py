import requests

PDFs = [
    {'title': "Attention Is All You Need", 'file': "https://arxiv.org/pdf/1706.03762"},
    {'title': "Deep Residual Learning", 'file': "https://arxiv.org/pdf/1512.03385"},
    {'title': "BERT", 'file': "https://arxiv.org/pdf/1810.04805"},
    {'title': "GPT-3", 'file': "https://arxiv.org/pdf/2005.14165"},
    {'title': "Adam Optimizer", 'file': "https://arxiv.org/pdf/1412.6980"},
    {'title': "GANs", 'file': "https://arxiv.org/pdf/1406.2661"},
    {'title': "U-Net", 'file': "https://arxiv.org/pdf/1505.04597"},
    {'title': "DALL-E 2", 'file': "https://arxiv.org/pdf/2204.06125"},
    {'title': "Stable Diffusion", 'file': "https://arxiv.org/pdf/2112.10752"}
]

for pdf in PDFs:
    with open(f"{pdf['title']}.pdf", "wb") as f:
        f.write(requests.get(pdf['file']).content)