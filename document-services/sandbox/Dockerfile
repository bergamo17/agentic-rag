FROM python:3.11-slim

ENV PYTHONDONTWRTIEBYCODE=1 \
    PYTHONUNBUFFERED=1

RUN pip install --no-cache-dir \
    python-docx==1.1.2 \
    lxml==6.1.3 \
    reportlab==4.5.1 \
    openpyxl==3.1.5

RUN mkdir -p /workspace/input /workspace/output
WORKDIR /workspace

RUN useradd -m sandboxuser && chown -R sandboxuser:sandboxuser /workspace
USER sandboxuser