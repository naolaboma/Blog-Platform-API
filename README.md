# Blog API

A robust, scalable RESTful API for a blog platform built with Go, featuring user authentication, blog management, AI-powered content generation, and comprehensive search capabilities.

## Table of Contents

- [Features](#features)
- [Technology Stack](#technology-stack)
- [Prerequisites](#prerequisites)
- [Installation](#installation)
- [Configuration](#configuration)
- [Database Setup](#database-setup)
- [Running the Application](#running-the-application)
- [API Documentation](#api-documentation)
- [Project Structure](#project-structure)
- [Testing](#testing)
- [Deployment](#deployment)
- [Contributing](#contributing)
- [License](#license)

## Features

### Core Functionality

- **User Management**: Registration, authentication, profile management, and role-based access control
- **Blog Management**: Create, read, update, and delete blog posts with embedded comments and reactions
- **Search & Filtering**: Advanced search by title, author, tags, and date with pagination support
- **AI Integration**: Powered by Groq AI for blog content generation, enhancement, and idea suggestions
- **Authentication**: JWT-based authentication with refresh tokens and session management
- **OAuth Integration**: Support for Google and GitHub authentication
- **Email Services**: Email verification and password reset functionality
- **File Upload**: Profile picture upload with validation and storage

### Technical Features

- **High Performance**: Redis caching, database indexing, and optimized queries
- **Scalability**: Worker pools, goroutines, and connection pooling
- **Security**: Password hashing, input validation, and role-based authorization
- **Monitoring**: Graceful shutdown, health checks, and error handling
- **Documentation**: Comprehensive Postman collection and API documentation

## Technology Stack

- **Backend**: Go 1.23+
- **Framework**: Gin (HTTP web framework)
- **Database**: MongoDB with proper indexing
- **Cache**: Redis for performance optimization
- **Authentication**: JWT tokens with bcrypt password hashing
- **AI Integration**: Groq API
- **OAuth**: Google and GitHub integration
- **Email**: SMTP with HTML templates
- **File Storage**: Local filesystem with validation
- **Validation**: Go validator package

## Prerequisites

Before running this application, ensure you have the following installed:

- **Go**: Version 1.23 or higher
- **MongoDB**: Version 6.0 or higher
- **Redis**: Version 7.0 or higher
- **Git**: For cloning the repository

## Installation

1. Clone the repository:

```bash
git clone <repository-url>
cd Blog-API
```

2. Install Go dependencies:

```bash
go mod download
```

3. Create environment configuration file:

```bash
cp .env.example .env
```

## Configuration

Create a `.env` file in the root directory as the example env

## Database Setup

### Option 1: Using the Setup Script

1. Ensure MongoDB is running:

```bash
sudo systemctl start mongod
sudo systemctl enable mongod
```

2. Run the setup script:

```bash
mongosh < mongodb_setup.js
```

## Running the Application

1. Start Redis server:

```bash
sudo systemctl start redis-server
sudo systemctl enable redis-server
```

2. Run the application:

```bash
go run cmd/server/main.go
```

The server will start on `http://localhost:8080` by default.

## API Documentation

### Base URL

```
http://localhost:8080/api/v1
```

### Authentication Endpoints

#### User Registration

```http
POST /auth/register
Content-Type: application/json

{
  "username": "testuser",
  "email": "test@example.com",
  "password": "SecurePass123!"
}
```

#### User Login

```http
POST /auth/login
Content-Type: application/json

{
  "email": "test@example.com",
  "password": "SecurePass123!"
}
```

#### Refresh Token

```http
POST /auth/refresh
Content-Type: application/json

{
  "refresh_token": "your-refresh-token"
}
```

#### Logout

```http
POST /auth/logout
Authorization: Bearer <access-token>
```

#### Email Verification

```http
POST /auth/send-verification
Content-Type: application/json

{
  "email": "test@example.com"
}
```

#### Password Reset

```http
POST /auth/forgot-password
Content-Type: application/json

{
  "email": "test@example.com"
}
```

### Blog Endpoints

#### Get All Blogs

```http
GET /blogs?page=1&limit=10&sort=newest
```

#### Get Blog by ID

```http
GET /blogs/{blog-id}
```

#### Create Blog (Authenticated)

```http
POST /blogs
Authorization: Bearer <access-token>
Content-Type: application/json

{
  "title": "My Blog Post",
  "content": "This is the content of my blog post.",
  "tags": ["technology", "programming"]
}
```

#### Update Blog (Author/Admin Only)

```http
PUT /blogs/{blog-id}
Authorization: Bearer <access-token>
Content-Type: application/json

{
  "title": "Updated Blog Title",
  "content": "Updated content."
}
```

#### Delete Blog (Author/Admin Only)

```http
DELETE /blogs/{blog-id}
Authorization: Bearer <access-token>
```

#### Search Blogs by Title

```http
GET /blogs/search/title?title=search-term&page=1&limit=10
```

#### Search Blogs by Author

```http
GET /blogs/search/author?author=username&page=1&limit=10
```

#### Filter Blogs by Tags

```http
GET /blogs/filter/tags?tags=technology,programming&page=1&limit=10
```

#### Get Popular Blogs

```http
GET /blogs/popular?limit=10
```

### User Management Endpoints

#### Get User Profile (Authenticated)

```http
GET /users/profile
Authorization: Bearer <access-token>
```

#### Update User Profile (Authenticated)

```http
PUT /users/profile
Authorization: Bearer <access-token>
Content-Type: application/json

{
  "username": "newusername",
  "bio": "Updated bio information"
}
```

#### Upload Profile Picture (Authenticated)

```http
POST /users/profile/picture
Authorization: Bearer <access-token>
Content-Type: multipart/form-data

profile_picture: <file>
```

### AI Integration Endpoints

#### Generate Blog Content (Authenticated)

```http
POST /ai/generate-blog
Authorization: Bearer <access-token>
Content-Type: application/json

{
  "topic": "Go Programming Best Practices"
}
```

#### Enhance Blog Content (Authenticated)

```http
POST /ai/enhance-blog
Authorization: Bearer <access-token>
Content-Type: application/json

{
  "content": "Your blog content here..."
}
```

#### Suggest Blog Ideas (Authenticated)

```http
POST /ai/suggest-ideas
Authorization: Bearer <access-token>
Content-Type: application/json

{
  "keywords": ["technology", "programming", "golang"]
}
```

### Admin Endpoints

#### Promote User to Admin

```http
PUT /admin/users/{user-id}/promote
Authorization: Bearer <admin-access-token>
Content-Type: application/json

{
  "role": "admin"
}
```

#### Demote Admin to User

```http
PUT /admin/users/{user-id}/demote
Authorization: Bearer <admin-access-token>
Content-Type: application/json

{
  "role": "user"
}
```

## Project Structure

```
Blog-API/
├── cmd/
│   └── server/
│       └── main.go                 # Application entry point
├── internal/
│   ├── delivery/                   # HTTP layer
│   │   ├── controllers/           # Request handlers
│   │   └── router/                # Route definitions
│   ├── domain/                    # Business logic interfaces
│   ├── infrastructure/            # External services
│   │   ├── ai/                   # AI service integration
│   │   ├── cache/                # Redis caching
│   │   ├── database/             # MongoDB connection
│   │   ├── email/                # Email service
│   │   ├── filesystem/           # File upload handling
│   │   ├── jwt/                  # JWT authentication
│   │   ├── middleware/           # HTTP middleware
│   │   ├── oauth/                # OAuth integration
│   │   ├── password/             # Password utilities
│   │   └── worker/               # Background job processing
│   ├── repository/               # Data access layer
│   └── usecase/                  # Business logic implementation
├── pkg/
│   └── config/                   # Configuration management
├── docs/
│   └── postman/                  # API documentation
├── mongodb_collections.json      # Database schema documentation
├── mongodb_setup.js             # Database setup script
└── README.md                     # This file
```

## Testing

### Manual Testing

Use the provided Postman collection in `docs/postman/Blog-API.postman_collection.json` to test all endpoints.

### Automated Testing

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run specific package tests
go test ./internal/usecase
```

## Deployment

### Production Considerations

1. **Environment Variables**: Ensure all sensitive configuration is properly set
2. **Database**: Use production MongoDB instance with proper authentication
3. **Redis**: Configure Redis with authentication and persistence
4. **HTTPS**: Enable TLS/SSL for production endpoints
5. **Monitoring**: Implement logging, metrics, and health checks
6. **Scaling**: Consider horizontal scaling with load balancers

### Development Guidelines

- Follow Go coding standards and conventions
- Write comprehensive tests for new features
- Update documentation for API changes
- Ensure all tests pass before submitting PR
- Use meaningful commit messages
- Follow pure clean architecture approach
