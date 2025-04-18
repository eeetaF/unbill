# We make splitting bills easy!

![Unbill logo[product-screenshot]](logo.jpg)

---

## Why to use unbill (rename this section pls)

---

## Stack
Every microservice runs in **docker** container.
- ### Frontend
**Python**, telebot
- ### Backend
**Go**, http standard package
- ### OCR and NLP service
**Python**, easyocr, cv2

---

## Endpoints
Backend listens on **localhost:8080** 
- ### `POST /api/upload_and_analyze`
**Upload image and analyze bill on it**

Request body:
```json
{
  "ext": "png",
  "content": "base64encodedImageContent"
}
```
Response body:
```json
{
  "filename": "generatedname.jpeg",
  "product_units": [
    {
      "id": 0,
      "name": "product",
      "quantity": 1,
      "price": 400
    }
  ],
  "total_price": 400
}
```
- ### `GET /api/split_equally/{filename}/{num_people}`
**Split the uploaded bill evenly between num_people**

Response body:
```json
{
  "total_per_person": 1300
}
```


---

## How to use

---

## Project structure

---

## Why unbill is safe to use

---