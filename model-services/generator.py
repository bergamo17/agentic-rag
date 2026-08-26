from PIL.Image import Image
from io import BytesIO
import torch
import openai
import base64

model_name = "gpt-4o"
system_prompt = "You are an expert professional PDF analyst who gives rigorous in-depth answers."

_client = None

def _get_client():
    global _client
    if _client is None:
        _client = openai.OpenAI()
    return _client

def _image_to_data_url(image: Image) -> str:
    buffer = BytesIO()
    image.save(buffer, format="PNG")
    img_b64 = base64.b64encode(buffer.getvalue()).decode("utf-8")
    return f"data:image/png;base64,{img_b64}"

def query_vlm(query: str, images: list[Image]) -> str:
    client = _get_client()
    
    message_content = [
        {"type": "image_url", "image_url": {"url": _image_to_data_url(image)}}
        for image in images
    ] + [{"type": "text", "text": query}]

    messages = [
        {
            "role": "system",
            "content": system_prompt
        },
        {
            "role": "user",
            "content": message_content
        }
    ]

    response = client.chat.completions.create(
        model=model_name,
        messages=messages,
        max_tokens=1000,
    )
    return response.choices[0].message.content