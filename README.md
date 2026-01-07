# Go Clean Architecture Boilerplate 🚀

A production-ready Go backend API boilerplate following **Clean Architecture** principles. Built to be as easy to use as Laravel - focus on your business logic while the infrastructure is handled for you.

## ✨ Features

- 🏗️ **Clean Architecture** - Separation of concerns with clear layers (Domain, Repository, Service, Handler)
- 🚀 **Gin Framework** - Fast HTTP web framework
- 🗄️ **GORM** - Powerful ORM with PostgreSQL support
- 🔧 **Viper** - Configuration management (YAML + Environment variables)
- 🔄 **golang-migrate** - Database migration management
- 🔐 **JWT Authentication** - Secure authentication with golang-jwt
- 📝 **Zap Logger** - High-performance structured logging
- 🔥 **Air** - Live reload for development
- 🐳 **Docker** - Fully dockerized with docker-compose
- ✅ **Request Validation** - Built-in validation with go-playground/validator
- 📦 **Standardized Response** - Consistent API response format
- 🔒 **Security** - Password hashing with bcrypt, CORS, and more

## 📁 Project Structure

```
go-clean-boiler/
├── cmd/
│   └── api/
│       └── main.go                 # Application entry point
├── internal/
│   ├── domain/                     # Entities/Models
│   │   └── user.go
│   ├── repository/                 # Data access layer
│   │   ├── user_repository.go      # Interface
│   │   └── postgres/
│   │       └── user_repository.go  # Implementation
│   ├── service/                    # Business logic
│   │   ├── user_service.go
│   │   └── auth_service.go
│   ├── handler/                    # HTTP handlers/controllers
│   │   ├── user_handler.go
│   │   └── auth_handler.go
│   ├── middleware/                 # HTTP middlewares
│   │   ├── auth.go
│   │   ├── logger.go
│   │   ├── error.go
│   │   └── cors.go
│   ├── dto/                        # Data Transfer Objects
│   │   ├── request/
│   │   └── response/
│   └── router/                     # Route definitions
│       └── router.go
├── pkg/                            # Shared utilities
│   ├── config/                     # Configuration
│   ├── database/                   # Database setup
│   ├── logger/                     # Logger setup
│   ├── jwt/                        # JWT utilities
│   ├── response/                   # Response format
│   └── validator/                  # Validation
├── migrations/                     # Database migrations
├── config/                         # Config files
├── .air.toml                       # Air config
├── Dockerfile
├── docker-compose.yml
├── Makefile
└── README.md
```

## 🚀 Quick Start

### Prerequisites

- Go 1.21 or higher
- PostgreSQL 15 or higher (or use Docker)
- golang-migrate CLI (optional, for manual migrations)

### Option 1: Using Docker (Recommended)

1. **Clone the repository**
```bash
git clone <repository-url>
cd go-clean-boiler
```

2. **Copy environment file**
```bash
cp .env.example .env
```

3. **Start with Docker Compose**
```bash
make docker-up
```

The API will be available at `http://localhost:8080`

### Option 2: Local Development

1. **Clone the repository**
```bash
git clone <repository-url>
cd go-clean-boiler
```

2. **Install dependencies**
```bash
make deps
```

3. **Set up PostgreSQL**
```bash
# Create database
createdb go_clean_boiler
```

4. **Copy and configure environment**
```bash
cp .env.example .env
# Edit .env with your database credentials
```

5. **Run migrations** (optional, auto-migration is enabled)
```bash
make migrate-up
```

6. **Run with hot reload**
```bash
make dev
```

Or run without hot reload:
```bash
make run
```

## 🛠️ Available Commands

```bash
make help          # Display all available commands
make dev           # Run with hot reload (Air)
make run           # Run without hot reload
make build         # Build the application
make test          # Run tests
make clean         # Clean build files
make docker-up     # Start Docker containers
make docker-down   # Stop Docker containers
make docker-logs   # View Docker logs
make migrate-up    # Run database migrations
make migrate-down  # Rollback database migrations
make migrate-create name=migration_name  # Create new migration
make tidy          # Tidy go modules
make deps          # Download dependencies
```

## 📝 API Endpoints

### Authentication

```bash
# Register
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123",
  "name": "John Doe"
}

# Login
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "password123"
}
```

### Users (Protected - Requires JWT Token)

```bash
# Get all users (with pagination)
GET /api/v1/users?page=1&per_page=10
Authorization: Bearer <your-jwt-token>

# Get user by ID
GET /api/v1/users/:id
Authorization: Bearer <your-jwt-token>

# Create user
POST /api/v1/users
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "email": "newuser@example.com",
  "password": "password123",
  "name": "Jane Doe"
}

# Update user
PUT /api/v1/users/:id
Authorization: Bearer <your-jwt-token>
Content-Type: application/json

{
  "email": "updated@example.com",
  "name": "Updated Name"
}

# Delete user
DELETE /api/v1/users/:id
Authorization: Bearer <your-jwt-token>
```

### Health Check

```bash
GET /health
```

## 🎯 How to Add New Features

This boilerplate makes it easy to add new features. Here's a step-by-step guide:

### 1. Create Migration

```bash
make migrate-create name=create_products_table
```

Edit the generated migration file:
```sql
-- migrations/000002_create_products_table.up.sql
CREATE TABLE IF NOT EXISTS products (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    price DECIMAL(10, 2) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);
```

Run migration:
```bash
make migrate-up
```

### 2. Create Entity (Domain)

Create `internal/domain/product.go`:
```go
package domain

import (
    "time"
    "gorm.io/gorm"
)

type Product struct {
    ID        uint           `gorm:"primarykey" json:"id"`
    Name      string         `gorm:"not null" json:"name"`
    Price     float64        `gorm:"not null" json:"price"`
    CreatedAt time.Time      `json:"created_at"`
    UpdatedAt time.Time      `json:"updated_at"`
    DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
```

### 3. Create DTOs

Create `internal/dto/request/product_request.go`:
```go
package request

type CreateProductRequest struct {
    Name  string  `json:"name" validate:"required"`
    Price float64 `json:"price" validate:"required,gt=0"`
}

type UpdateProductRequest struct {
    Name  string  `json:"name" validate:"omitempty"`
    Price float64 `json:"price" validate:"omitempty,gt=0"`
}
```

Create `internal/dto/response/product_response.go`:
```go
package response

import "time"

type ProductResponse struct {
    ID        uint      `json:"id"`
    Name      string    `json:"name"`
    Price     float64   `json:"price"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

### 4. Create Repository

Create interface in `internal/repository/product_repository.go`:
```go
package repository

import "github.com/firdanbash/go-clean-boiler/internal/domain"

type ProductRepository interface {
    Create(product *domain.Product) error
    FindByID(id uint) (*domain.Product, error)
    FindAll(limit, offset int) ([]domain.Product, int64, error)
    Update(product *domain.Product) error
    Delete(id uint) error
}
```

Create implementation in `internal/repository/postgres/product_repository.go`:
```go
package postgres

import (
    "github.com/firdanbash/go-clean-boiler/internal/domain"
    "github.com/firdanbash/go-clean-boiler/internal/repository"
    "gorm.io/gorm"
)

type productRepository struct {
    db *gorm.DB
}

func NewProductRepository(db *gorm.DB) repository.ProductRepository {
    return &productRepository{db: db}
}

// Implement all interface methods...
```

### 5. Create Service

Create `internal/service/product_service.go`:
```go
package service

import (
    "github.com/firdanbash/go-clean-boiler/internal/dto/request"
    "github.com/firdanbash/go-clean-boiler/internal/dto/response"
    "github.com/firdanbash/go-clean-boiler/internal/repository"
)

type ProductService interface {
    Create(req *request.CreateProductRequest) (*response.ProductResponse, error)
    GetByID(id uint) (*response.ProductResponse, error)
    GetAll(page, perPage int) ([]response.ProductResponse, int64, error)
    Update(id uint, req *request.UpdateProductRequest) (*response.ProductResponse, error)
    Delete(id uint) error
}

type productService struct {
    repo repository.ProductRepository
}

func NewProductService(repo repository.ProductRepository) ProductService {
    return &productService{repo: repo}
}

// Implement all business logic methods...
```

### 6. Create Handler

Create `internal/handler/product_handler.go`:
```go
package handler

import (
    "github.com/firdanbash/go-clean-boiler/internal/service"
    "github.com/gin-gonic/gin"
)

type ProductHandler struct {
    productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
    return &ProductHandler{productService: productService}
}

func (h *ProductHandler) Create(c *gin.Context) {
    // Implementation...
}

// Implement all handler methods...
```

### 7. Register Routes

Edit `internal/router/router.go`:
```go
// Add to SetupRouter function
products := v1.Group("/products")
products.Use(middleware.AuthMiddleware(jwtSecret))
{
    products.GET("", productHandler.GetAll)
    products.GET("/:id", productHandler.GetByID)
    products.POST("", productHandler.Create)
    products.PUT("/:id", productHandler.Update)
    products.DELETE("/:id", productHandler.Delete)
}
```

### 8. Wire Dependencies

Edit `cmd/api/main.go`:
```go
// Initialize repositories
productRepo := postgres.NewProductRepository(database.DB)

// Initialize services
productService := service.NewProductService(productRepo)

// Initialize handlers
productHandler := handler.NewProductHandler(productService)

// Pass to router
r := router.SetupRouter(authHandler, userHandler, productHandler, cfg.JWT.Secret)
```

## ⚙️ Configuration

Configuration is managed via Viper and supports both YAML files and environment variables.

### Config File (`config/config.yaml`)

```yaml
app:
  name: go-clean-boiler
  env: development
  port: 8080

database:
  host: localhost
  port: 5432
  user: postgres
  password: postgres
  name: go_clean_boiler
  sslmode: disable

jwt:
  secret: your-secret-key
  expiration: 24h

log:
  level: debug
  encoding: console
```

### Environment Variables

Environment variables override config file values:
- `APP_PORT`
- `DB_HOST`, `DB_PORT`, `DB_USER`, `DB_PASSWORD`, `DB_NAME`
- `JWT_SECRET`, `JWT_EXPIRATION`
- `LOG_LEVEL`

## 🧪 Testing

Create tests following Go conventions:

```bash
# Run all tests
make test

# Run specific package tests
go test -v ./internal/service/...

# Run with coverage
go test -cover ./...
```

## 📦 Deployment

### Docker Deployment

```bash
# Build and run with Docker Compose
make docker-up

# View logs
make docker-logs

# Stop containers
make docker-down
```

### Manual Deployment

```bash
# Build binary
make build

# Run binary
./bin/main
```

## 🔒 Security Best Practices

- ✅ Passwords are hashed with bcrypt
- ✅ JWT tokens for authentication
- ✅ CORS middleware included
- ✅ SQL injection protection via GORM
- ✅ Input validation on all requests
- ⚠️ **Change JWT_SECRET in production!**
- ⚠️ Use strong database passwords
- ⚠️ Enable SSL in production (sslmode=require)

## 📚 Tech Stack

- **Framework**: [Gin](https://github.com/gin-gonic/gin)
- **ORM**: [GORM](https://gorm.io/)
- **Config**: [Viper](https://github.com/spf13/viper)
- **Logger**: [Zap](https://github.com/uber-go/zap)
- **JWT**: [golang-jwt](https://github.com/golang-jwt/jwt)
- **Validation**: [validator](https://github.com/go-playground/validator)
- **Migration**: [golang-migrate](https://github.com/golang-migrate/migrate)
- **Hot Reload**: [Air](https://github.com/cosmtrek/air)

## 🤝 Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## 📄 License

This project is licensed under the MIT License.

## 💡 Tips

- Use `make dev` for development with hot reload
- Keep your `.env` file secure and never commit it
- Follow the clean architecture pattern when adding features
- Write tests for your business logic
- Use migrations for database changes
- Check logs with `make docker-logs` when using Docker

## 🆘 Troubleshooting

### Database connection failed
- Make sure PostgreSQL is running
- Check database credentials in `.env`
- If using Docker, ensure containers are running: `docker-compose ps`

### Air not found
- Run `go install github.com/cosmtrek/air@latest`
- Make sure `$GOPATH/bin` is in your PATH

### Migrations not running
- Install migrate CLI: `brew install golang-migrate` (macOS) or download from [releases](https://github.com/golang-migrate/migrate/releases)
- Alternatively, use auto-migration (enabled by default)

## 📬 Contact

For questions or support, please open an issue on GitHub.

---

Made with ❤️ using Go Clean Architecture
