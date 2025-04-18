# We make splitting bills easy!

![Unbill logo[]](logo.jpg)

---

## Why to use unbill (rename this section pls)

## Stack
Every microservice runs in **docker** container.
- ### Frontend
**Python**, telebot
- ### Backend
**Go**, http standard package
- ### OCR and NLP service
**Python**, easyocr, cv2

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

## How to use

run script (я его еще не сделал)

## Project structure
- ### backend /
| app.ini (configure application)

| main.go (entry point)

other source directories...

- ### OCR_service /
| app/ocr_server.py, app/openrouter_nlp.py (source files)

| docker-compose.yml, Dockerfile

- ### data_shared /
shared directory among all containers (pictures of bills for backend and OCR_service communication)

- ### Telebot / ?

## Why unbill is safe to use

Even though all data about bills (including pictures and analyzed data) is located in the same directory, which means any user can make a request on any other's bill, it's impossible to guess the name of the file. 

We encode the name based on current time and a random key that is generated each time the program starts.

We also use an autocleaner (check backend/vacuum_cleaner). It clears files older than (backend/app.ini:maxFileAge) each (backend/app.ini:maxFileAge)/24 (by default, clears files older than 24 hours each hour)

This makes using shared directory practically safe and easy.
