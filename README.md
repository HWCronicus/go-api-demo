# Go API Demo

A RESTful API built with Go featuring user authentication, comments system, and PostgreSQL database integration.

## Features

- JWT authentication
- User registration and login
- Comments system
- PostgreSQL database with SQLC
- Swagger documentation
- Docker for easy deployment

## Quick Start

### Prerequisites

- Docker & Docker Compose

### Running the Application

```bash
docker-compose up
```

The API will be available at `http://localhost:8080` (or http://127.0.0.1:8080 if your machine is weird like mine..)

## API Endpoints

### Public Routes

- `GET /` - Serve index page
- `GET /resume` - Download resume PDF
- `POST /user` - Register new user
- `POST /login` - User authentication
- `GET /comments` - List all comments

### Protected Routes (Requires JWT)

- `POST /comment` - Create a new comment

### Documentation

- `GET /swagger/` - Interactive API documentation

## Tech Stack

- **Language**: Go 1.25.4
- **Database**: PostgreSQL 16
- **"ORM"**: SQLC
- **Authentication**: JWT with bcrypt
- **Documentation**: Swagger/OpenAPI
- **Containerization**: Docker & Docker Compose

## Project Structure

```
.
├── api/
│   ├── database/         # Generated SQLC code
│   ├── docs/            # Swagger documentation
│   ├── html/            # Static HTML files
│   ├── sql/
│   │   ├── queries/     # SQL queries
│   │   └── schema/      # Database schema
│   ├── main.go          # Main application
│   └── go.mod           # Go dependencies
└── docker-compose.yml   # Docker configuration
```
