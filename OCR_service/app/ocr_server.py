import cv2
import easyocr
import socket
from openrouter_nlp import nlp_generate

reader = easyocr.Reader(['en', 'ru'], gpu=True)

HOST = '0.0.0.0'
PORT = 8082

with socket.socket(socket.AF_INET, socket.SOCK_STREAM) as s:
    s.bind((HOST, PORT))
    s.listen(1)
    print(f"OCR server listening on {HOST}:{PORT}...")
    while True:
        conn, addr = s.accept()
        with conn:
            filename = conn.recv(1024).decode().strip()
            image_path = f"/data/{filename}"
            image = cv2.imread(image_path)
            if image is None:
                conn.sendall(b"ERROR: Image not found.\n")
                continue

            results = reader.readtext(image, detail=1, paragraph=False)
            full_text = ' '.join([text for (_, text, _) in results])

            try:
                nlp_answer = nlp_generate(full_text)
                conn.sendall(nlp_answer.encode('utf-8'))
            except Exception as e:
                conn.sendall(f"ERROR: NLP failed: {e}\n".encode('utf-8'))
