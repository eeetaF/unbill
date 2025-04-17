import os
import requests
from dotenv import load_dotenv

load_dotenv()
API_KEY = os.getenv("OPENROUTER_API_KEY")


def nlp_generate(prompt, min_tokens=50, max_tokens=200, temperature=0.7):
    if not API_KEY:
        raise RuntimeError("OPENROUTER_API_KEY not set in environment.")
    headers = {
        "Authorization": f"Bearer {API_KEY}",
        "Content-Type": "application/json",
        "HTTP-Referer": "http://localhost",
    }
    payload = {
        "model": "openai/gpt-3.5-turbo",
        "messages": [
            {"role": "system", "content": "You will see the result of OCR model. This model has read a restaurant bill. Analyze it and provide as output: name, quantity and price of each product in that bill in the following format: \"<name>|<quantity>|<price>\\n\". Also important: use integer for price (for example, detected 240.00 show as 24000). Avoid adding in the list of products extra information from the bill, such as total cost. Format uppercase and lowercase letters as you think will be correct (because the OCR may wrongly detect the case)"},
            {"role": "user", "content": prompt}
        ],
        "temperature": temperature,
        "max_tokens": max_tokens
    }
    response = requests.post("https://openrouter.ai/api/v1/chat/completions", headers=headers, json=payload)
    response.raise_for_status()
    return response.json()["choices"][0]["message"]["content"]
