FROM python:3.11-slim

ENV PYTHONDONTWRITEBYTECODE=1 \
    PYTHONUNBUFFERED=1 \
    PYTHONPATH=/opt/helpers

RUN apt-get update \
    && apt-get install -y --no-install-recommends fonts-dejavu-core fonts-dejavu-extra \
    && rm -rf /var/lib/apt/lists/*

RUN pip install --no-cache-dir \
    python-docx==1.1.2 \
    lxml==6.1.3 \
    reportlab==4.5.1 \
    openpyxl==3.1.5 \
    matplotlib

COPY helper/pdf_helpers.py /opt/helpers/pdf_helpers.py

RUN mkdir -p /workspace/input /workspace/output
WORKDIR /workspace

RUN useradd -m sandboxuser && chown -R sandboxuser:sandboxuser /workspace
USER sandboxuser