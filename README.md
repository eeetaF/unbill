# We make splitting bills easy!

![Unbill logo[]](logo.jpg)

---

## What makes UnBill different

- **Fast and intuitive** - Telegram bot provides a smooth and interactive bill-splitting experience without the need to install extra apps  
- **Smart receipt processing** - Integrated OCR and NLP modules analyze complex restaurant bills with high accuracy
- **Flexible splitting options** - Supports equal, item-based, and custom logic splitting tailored to real-world scenarios
- **Collaborative editing** - All participants can review and adjust the bill in real time
- **Secure and extendable** - Dockerized backend, clean API, optional third-party NLP integrations  

 ### Built for real-life messiness - from shared desserts to chaotic receipts


## Stack
Every microservice runs in **docker** container.
- ### Frontend
**Python**, telegram library
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
- ### `PUT /api/products/{filename}`
**Updates info about products in filename**

Request body:
```json
{
  "product_units": [
    {
      "id": 0,
      "name": "product",
      "quantity": 1,
      "price": 400
    }
  ],
}
```
Response body:
```json
{
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

## How to use

TODO...
for now, use docker-composes.yml )

## Project structure
- ### backend /
| app.ini (configure application)

| main.go (entry point)

| docker-compose.yml, Dockerfile

other source directories...

- ### OCR_service /
| app/ocr_server.py, app/openrouter_nlp.py (source files)

| docker-compose.yml, Dockerfile

- ### data_shared /
shared directory among all containers (pictures of bills for backend and OCR_service communication)

- ### python /
telebot source

## Why UnBill is safe to use

Even though all data about bills (including pictures and analyzed data) is located in the same directory, which means any user can make a request on any other's bill, it's impossible to guess the name of the file. 

We encode the name based on current time and a random key that is generated each time the program starts.

We also use an autocleaner (check backend/vacuum_cleaner). It clears files older than (backend/app.ini:maxFileAge) each (backend/app.ini:maxFileAge)/24 (by default, clears files older than 24 hours each hour)

This makes using shared directory practically safe and easy.
