# Online Bookstore API

A comprehensive RESTful API for managing an online bookstore built with Go, Fiber, and GORM. This API provides complete functionality for book management, user authentication, order processing, and category management.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/)
[![Fiber](https://img.shields.io/badge/Fiber-v2-00ADD8?style=for-the-badge&logo=go)](https://gofiber.io/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-13+-316192?style=for-the-badge&logo=postgresql)](https://www.postgresql.org/)
[![JWT](https://img.shields.io/badge/JWT-Authentication-000000?style=for-the-badge&logo=JSON%20web%20tokens)](https://jwt.io/)

## 🚀 Features

- **User Authentication & Authorization**
  - User registration and login
  - JWT token-based authentication
  - Protected routes with middleware
  - Password hashing with bcrypt

- **Book Management**
  - CRUD operations for books
  - Book categorization
  - Search and filtering capabilities
  - Inventory management

- **Order Management**
  - Create and manage orders
  - Order status tracking
  - Order history for users

- **Category Management**
  - Organize books by categories
  - Category-based filtering

- **Clean Architecture**
  - Repository pattern
  - Use case layer
  - Dependency injection
  - Structured logging

## 🏗️ Architecture

This project follows Clean Architecture principles:

```
├── cmd/
│   └── main.go                    # Application entry point
├── db/
│   └── migrations/
│       └── migrate.go             # Database migrations
├── internal/
│   ├── auth/
│   │   └── jwt.go                 # JWT service
│   ├── config/
│   │   ├── app.go                 # Application config
│   │   ├── fiber.go               # Fiber config
│   │   ├── gorm.go                # GORM config
│   │   ├── jwt.go                 # JWT config
│   │   ├── logrus.go              # Logging config
│   │   └── viper.go               # Configuration management
│   ├── delivery/
│   │   └── http/
│   │       ├── handler/           # HTTP handlers
│   │       ├── middleware/        # Custom middleware
│   │       └── routes/            # Route definitions
│   ├── entity/                    # Domain entities
│   ├── model/                     # Request/Response models
│   ├── repository/                # Data access layer
│   └── usecase/                   # Business logic layer
├── go.mod
├── go.sum
└── README.md
```

## 🛠️ Tech Stack

- **Framework**: [Fiber v2](https://gofiber.io/) - Fast HTTP web framework
- **Database**: [PostgreSQL](https://www.postgresql.org/) - Relational database
- **ORM**: [GORM](https://gorm.io/) - Go ORM library
- **Authentication**: [JWT](https://jwt.io/) - JSON Web Tokens
- **Validation**: [go-playground/validator](https://github.com/go-playground/validator)
- **Logging**: [Logrus](https://github.com/sirupsen/logrus)
- **Configuration**: [Viper](https://github.com/spf13/viper)
- **Password Hashing**: [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)

## 🚦 Getting Started

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 13+
- Git

### Installation

1. **Clone the repository**
   ```bash
   git clone https://github.com/fathirarya/online-bookstore-api.git
   cd online-bookstore-api
   ```

2. **Install dependencies**
   ```bash
   go mod download
   ```

3. **Set up environment variables**
   
   Create a `.env` file in the root directory:
   ```env
   # Server Configuration
   SERVER_HOST=localhost
   SERVER_PORT=8080
   
   # Database Configuration
   DB_HOST=localhost
   DB_PORT=
   DB_USER=
   DB_PASSWORD=your_password
   DB_NAME=bookstore_db
   DB_SSLMODE=disable
   
   # JWT Configuration
   JWT_SECRET_KEY=your-super-secret-jwt-key-change-in-production
   JWT_ISSUER=online-bookstore-api
   JWT_EXPIRE_DURATION=24h
   
   # Log Configuration
   LOG_LEVEL=info
   ```

4. **Set up PostgreSQL database**
   ```sql
   CREATE DATABASE bookstore_db;
   ```

5. **Run the application**
   ```bash
   go run cmd/main.go
   ```

The server will start on `http://localhost:8080`

## 📚 API Documentation

### Authentication Endpoints

#### Register User
```http
POST /api/auth/register
Content-Type: application/json

{
  "username": "johndoe",
  "email": "john@example.com",
  "password": "password123"
}
```

#### Login User
```http
POST /api/auth/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "password123"
}
```

#### Get User Profile (Protected)
```http
GET /api/user/profile
Authorization: Bearer <your-jwt-token>
```

### Book Endpoints

#### Get All Books
```http
GET /api/books
```

#### Get Book by ID
```http
GET /api/books/:id
```

#### Create Book (Protected)
```http
POST /api/books
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "title": "Book Title",
  "author": "Author Name",
  "price": 29.99,
  "category_id": 1,
  "stock": 100,
  "description": "Book description"
}
```

#### Update Book (Protected)
```http
PUT /api/books/:id
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "title": "Updated Title",
  "price": 34.99
}
```

#### Delete Book (Protected)
```http
DELETE /api/books/:id
Authorization: Bearer <your-jwt-token>
```

### Category Endpoints

#### Get All Categories
```http
GET /api/categories
```

#### Create Category (Protected)
```http
POST /api/categories
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "name": "Fiction",
  "description": "Fiction books"
}
```

### Order Endpoints

#### Create Order (Protected)
```http
POST /api/orders
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "books": [
    {
      "book_id": 1,
      "quantity": 2
    }
  ]
}
```

#### Get User Orders (Protected)
```http
GET /api/orders
Authorization: Bearer <your-jwt-token>
```


### API Testing with curl

**Register a new user:**
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

**Login:**
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

**Get profile (replace TOKEN with actual JWT token):**
```bash
curl -X GET http://localhost:8080/api/user/profile \
  -H "Authorization: Bearer TOKEN"
```

## 📁 Project Structure Details

### Entities
- `User` - User account information
- `BookEntity` - Book catalog data
- `CategoryEntity` - Book categories
- `OrderEntity` - Customer orders
- `BookOrdersEntity` - Order items (many-to-many relation)

### Key Components

- **Handlers**: Process HTTP requests and responses
- **Use Cases**: Contain business logic
- **Repositories**: Handle data persistence
- **Middleware**: Authentication, logging, CORS, etc.
- **Models**: Request/response data structures

## 🔒 Security Features

- **Password Security**: Passwords are hashed using bcrypt
- **JWT Authentication**: Secure token-based authentication
- **Protected Routes**: Middleware protection for sensitive endpoints
- **Input Validation**: Request data validation using struct tags
- **SQL Injection Prevention**: GORM provides built-in protection


## 🙏 Acknowledgments

- [Fiber](https://gofiber.io/) for the amazing web framework
- [GORM](https://gorm.io/) for the powerful ORM
- [JWT](https://jwt.io/) for secure authentication
- All contributors and the Go community

---

⭐ Don't forget to star this repository if you found it helpful!
