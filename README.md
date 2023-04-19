# URL Shortener API

A simple URL shortener API built using Go, Gin, and MongoDB.

## Tech Stack

- **Go**: The backend server is written in Go.
- **Gin**: Gin is used as the web framework for the API.
- **MongoDB**: MongoDB is used as the database to store the shortened URLs and their corresponding original URLs.

## Features

- Create a short URL from a given URL.
- Redirect users to the original URL when accessing the short URL.

## Getting Started

### Prerequisites

- Go (version 1.16 or later)
- MongoDB
- An AWS account (for deployment)

### Running the API locally

1. Clone the repository:

```bash
git clone https://github.com/lewismunday/url-shortener.git
cd url-shortener
```


2. Install dependencies:

```bash
go mod download
```


3. Start the local MongoDB instance or configure the `.env` file with the MongoDB Atlas connection string.

4. Run the API:

```bash
go run main.go
```


The API should now be accessible at `http://localhost:8080`.

## API Usage

1. To create a short URL, send a POST request to the `/shorten` endpoint with a JSON payload containing the URL:

```json
{
  "url": "https://www.example.com"
}
```


The server will return the short URL:

```json
{
    "message": "URL inserted successfully",
    "result": {
        "InsertedID": "64404fa2b3eab94a423c4d72"
    },
    "shortUrl": "wVSBL"
}
```

2. To access the original URL, visit the short URL in a web browser or send a GET request to the short URL. The server will redirect to the original URL:

Accessing via `http://localhost:8080/<shortURL>`

For example, going to `http://localhost:8080/wVSBL` will direct you to `https://www.example.com`

