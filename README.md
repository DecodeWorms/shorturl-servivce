Go URL Shortener Service

A lightweight URL shortening service built in Go, using MongoDB for storage, Redis for caching, and Gin for routing.

🛠️ Tech Stack

Go 1.24+ – Backend implementation

MongoDB – Primary database

Redis – Cache layer

Gin – Web framework

Docker – Containerization

GitHub Actions – Continuous integration

🚀 Setup Instructions

Prerequisites

Go installed (>= 1.24)

Docker + Docker Compose

MongoDB & Redis running locally (or use Docker)

Clone the Repository
git clone https://github.com/DecodeWorms/shorturl-service
cd shorturl-service

Configure Environment

cp .env.example .env

Update your .env file with correct MongoDB and Redis URIs

Run Locally

go run main.go

Or Run with Docker

docker-compose up --build

📅 API Endpoints

Shorten URL

Create user
POST /api/v1/user/?id=user_id

{
"user_name": "John Doe",
"email": "john@example.com"
}

Create a short url
POST /api/v1/url/?id=user_id
{
"long_url":"https://github.com/login"
}
Response{
"short_url":"mjsgz7q"
}

Redirect to a long_url
GET /api/v1/url/?short_url=mjsgz7q


🔧 Running Tests

go test -v -cover ./...

🛠️ Deployment Notes

Includes GitHub Actions CI pipeline

Waits for Mongo to be ready

Test coverage report uploaded as artifact