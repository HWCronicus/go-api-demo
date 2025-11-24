# Go API Demo

This was a project I built primarily to showcase a Go RESTful API I built interfacing with a PostgreSQL database over the weekend. But it turned into a full-stack application with a Next.js frontend.. I was told I needed some UI..

## Features

- JWT authentication
- User registration and login
- Comments system
- PostgreSQL database with SQLC "ORM"
- Swagger documentation
- Next.js 16 frontend with Tailwind CSS
- Docker for easy deployment

## Quick Start

### Prerequisites

- Docker & Docker Compose

### Running the Application

```bash
docker-compose up
```

Services will be available at:

- **Go API**: `http://localhost:8080` (or 127.0.0.1:8080 if your PC is different like mine...)
- **Next.js App**: `http://localhost:3000`
- **Swagger Docs**: `http://localhost:8080/swagger/`

## API Endpoints

### Public Routes

- `GET /` - Server index/health page
- `GET /resume` - Download my resume PDF !!!
- `POST /user` - Create a new user
- `POST /login` - User Login / User authentication
- `GET /comments` - List them comments

### Protected Routes (Requires JWT)

- `POST /comment` - Create a new comment
- `DELETE /comments` - Delete a comment, if you created it.

### Documentation

- `GET /swagger/` - Interactive API documentation (Swagger Docs)

## Tech Stack

### Backend

- **Language**: Go 1.25.4
- **Database**: PostgreSQL 16
- **"ORM"**: SQLC
- **Authentication**: JWT with bcrypt
- **Documentation**: Swagger/OpenAPI

### Frontend

- **Framework**: Next.js 16
- **Styling**: Tailwind CSS 4
- **UI Components**: Radix UI
- **Language**: TypeScript

### Infrastructure

- **Containerization**: Docker & Docker Compose

## Project Structure

```text
.
├── api/
│   ├── database/        # Generated SQLC code
│   ├── docs/            # Swagger documentation
│   ├── files/           # Static other files
│   ├── html/            # Static HTML files
│   ├── sql/
│   │   ├── queries/     # SQL queries
│   │   └── schema/      # Database schema
│   ├── main.go          # Main application
│   └── go.mod           # Go dependencies
├── app/
│   ├── app/             # Next.js app directory
│   ├── components/      # React components
│   ├── lib/             # Utility functions
│   └── package.json     # Node dependencies
└── docker-compose.yml   # Docker configuration
```
