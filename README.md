# 📚 Online Bookstore API

RESTful API untuk sistem toko buku online yang dibangun dengan Go, Fiber, dan GORM. API ini menyediakan fitur lengkap untuk manajemen buku, autentikasi user, dan pemrosesan pesanan.

[![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go)](https://golang.org/) [![Fiber](https://img.shields.io/badge/Fiber-v2-00ADD8?style=for-the-badge&logo=go)](https://gofiber.io/) [![MySQL](https://img.shields.io/badge/MySQL-8+-4479A1?style=for-the-badge&logo=mysql)](https://www.mysql.com/)

## 🚀 Fitur

- **User Authentication & Authorization** - Registrasi dan login dengan JWT
- **Book Management** - CRUD buku dengan kategorisasi
- **Order Management** - Sistem pemesanan lengkap
- **Category Management** - Organisasi buku berdasarkan kategori
- **Clean Architecture** - Repository pattern dan dependency injection

## 🛠️ Tech Stack

- **Framework**: Fiber v2
- **Database**: MySQL 8+
- **ORM**: GORM
- **Auth**: JWT
- **Config**: Viper
- **Logging**: Logrus

## 💻 Installation

```bash
git clone https://github.com/fathirarya/online-bookstore-api.git
cd online-bookstore-api
go mod download
```

## 🔧 Environment Variables

Buat file `.env` di root directory:

```env
# Server Configuration
SERVER_HOST=localhost
SERVER_PORT=8080

# Database Configuration
DB_HOST=localhost
DB_PORT=3306
DB_USER=your_user
DB_PASSWORD=your_password
DB_NAME=bookstore_db

# JWT Configuration
JWT_SECRET_KEY=your-super-secret-jwt-key
JWT_ISSUER=online-bookstore-api
JWT_EXPIRE_DURATION=24h

# Log Configuration
LOG_LEVEL=info
```

## 🚀 Run Application

```bash
go run cmd/main.go
```

Server akan berjalan di `http://localhost:8080`

## 📖 API Documentation

**Dokumentasi Lengkap**: [Postman Documentation](https://documenter.getpostman.com/view/30637751/2sB3HgPNyp)

### Base URL
```
http://localhost:8080/api
```

## 🛠 Main Endpoints

### Authentication
- `POST /auth/register` - Registrasi user baru
- `POST /auth/login` - Login user
- `GET /user/profile` - Profile user (Protected)

### Books
- `GET /books` - Get semua buku
- `GET /books/:id` - Get buku by ID
- `POST /books` - Buat buku baru (Protected)
- `PUT /books/:id` - Update buku (Protected)
- `DELETE /books/:id` - Hapus buku (Protected)

### Categories
- `GET /categories` - Get semua kategori
- `POST /categories` - Buat kategori baru (Protected)

### Orders
- `GET /orders` - Get pesanan user (Protected)
- `POST /orders` - Buat pesanan baru (Protected)

## 🔐 Authentication

Gunakan JWT token di header:
```http
Authorization: Bearer <your-jwt-token>
```

## 📝 Contoh Usage

### Register
```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "testuser",
    "email": "test@example.com",
    "password": "password123"
  }'
```

### Login
```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "test@example.com",
    "password": "password123"
  }'
```

### Create Book
```bash
curl -X POST http://localhost:8080/api/books \
  -H "Authorization: Bearer <your-token>" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Book Title",
    "author": "Author Name",
    "price": 29.99,
    "category_id": 1,
    "stock": 100
  }'
```

## 🏗️ Project Structure

```
├── cmd/
│   └── main.go                # Application entry point
├── internal/
│   ├── auth/                  # JWT service
│   ├── config/                # Configuration
│   ├── delivery/http/         # HTTP handlers & routes
│   ├── entity/                # Domain entities
│   ├── repository/            # Data access layer
│   └── usecase/               # Business logic
├── db/migrations/             # Database migrations
└── README.md
```

## 🔒 Security Features

- Password hashing dengan bcrypt
- JWT token authentication
- Protected routes dengan middleware
- Input validation
- GORM SQL injection protection

## 📄 License

MIT License

---

**Made with ❤️ by [fathirarya](https://github.com/fathirarya)**
