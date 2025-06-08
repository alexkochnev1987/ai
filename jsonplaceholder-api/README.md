# JSONPlaceholder API Clone

A full-featured REST API that replicates the functionality of JSONPlaceholder.typicode.com with enhanced features including authentication, user management, and containerized deployment.

## 🚀 Features

- **Full CRUD Operations** - Complete user management (Create, Read, Update, Delete)
- **JWT Authentication** - Secure authentication with access and refresh tokens
- **RESTful Design** - Following REST API best practices
- **Database Integration** - PostgreSQL with GORM ORM
- **Containerized** - Docker and Docker Compose ready
- **Middleware Support** - CORS, logging, authentication, and security headers
- **Input Validation** - Request validation with proper error handling
- **Pagination** - Efficient data pagination for large datasets
- **Health Checks** - Built-in health monitoring endpoints
- **Graceful Shutdown** - Proper server shutdown handling

## 🛠 Tech Stack

- **Language**: Go 1.21+
- **Framework**: Gin HTTP Framework
- **Database**: PostgreSQL 15
- **ORM**: GORM
- **Authentication**: JWT tokens
- **Caching**: Redis (optional)
- **Containerization**: Docker & Docker Compose
- **Testing**: Go testing framework

## 📋 Prerequisites

- Go 1.21 or higher
- PostgreSQL 15+
- Docker & Docker Compose (for containerized deployment)
- Redis (optional, for caching)

## 🚀 Quick Start

### Using Docker Compose (Recommended)

1. **Clone the repository**

   ```bash
   git clone <repository-url>
   cd jsonplaceholder-api
   ```

2. **Start the services**

   ```bash
   docker-compose up -d
   ```

3. **The API will be available at:**
   - API: http://localhost:8080
   - Health Check: http://localhost:8080/health
   - Database: localhost:5432
   - Redis: localhost:6379

### Manual Setup

1. **Clone and setup**

   ```bash
   git clone <repository-url>
   cd jsonplaceholder-api
   go mod download
   ```

2. **Setup PostgreSQL database**

   ```bash
   createdb jsonplaceholder
   ```

3. **Configure environment variables**

   ```bash
   cp env.example .env
   # Edit .env with your configuration
   ```

4. **Run the application**
   ```bash
   go run cmd/server/main.go
   ```

## 📡 API Endpoints

### Authentication

| Method | Endpoint                | Description          | Auth Required |
| ------ | ----------------------- | -------------------- | ------------- |
| POST   | `/api/v1/auth/register` | Register new user    | No            |
| POST   | `/api/v1/auth/login`    | Login user           | No            |
| POST   | `/api/v1/auth/refresh`  | Refresh access token | No            |
| POST   | `/api/v1/auth/logout`   | Logout user          | No            |
| GET    | `/api/v1/auth/me`       | Get current user     | Yes           |

### Users

| Method | Endpoint             | Description     | Auth Required |
| ------ | -------------------- | --------------- | ------------- |
| GET    | `/api/v1/users`      | Get all users   | No            |
| GET    | `/api/v1/users/{id}` | Get user by ID  | No            |
| POST   | `/api/v1/users`      | Create new user | Yes           |
| PUT    | `/api/v1/users/{id}` | Update user     | Yes           |
| DELETE | `/api/v1/users/{id}` | Delete user     | Yes           |

### Health Check

| Method | Endpoint  | Description       |
| ------ | --------- | ----------------- |
| GET    | `/health` | API health status |

## 📝 API Usage Examples

### Register a new user

```bash
curl -X POST http://localhost:8080/api/v1/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "username": "johndoe",
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Login

```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Get all users

```bash
curl http://localhost:8080/api/v1/users?page=1&limit=10
```

### Get user by ID

```bash
curl http://localhost:8080/api/v1/users/1
```

### Create user (authenticated)

```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer YOUR_TOKEN" \
  -d '{
    "name": "Jane Doe",
    "username": "janedoe",
    "email": "jane@example.com",
    "password": "password123"
  }'
```

## 🔧 Configuration

Configuration is handled through environment variables. Copy `env.example` to `.env` and modify as needed:

```bash
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost

# Database Configuration
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=jsonplaceholder

# JWT Configuration
JWT_SECRET_KEY=your-super-secret-jwt-key
JWT_ACCESS_TOKEN_EXPIRY=1h
JWT_REFRESH_TOKEN_EXPIRY=24h
```

## 🧪 Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -cover ./...

# Run tests with detailed output
go test -v ./...
```

## 📦 Project Structure

```
jsonplaceholder-api/
├── cmd/
│   └── server/
│       └── main.go              # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── database/
│   │   └── database.go          # Database connection and migrations
│   ├── handlers/
│   │   ├── auth_handler.go      # Authentication endpoints
│   │   └── user_handler.go      # User CRUD endpoints
│   ├── middleware/
│   │   ├── auth.go              # JWT authentication middleware
│   │   └── middleware.go        # Other middleware
│   ├── models/
│   │   ├── auth.go              # Authentication models
│   │   ├── response.go          # Response models
│   │   └── user.go              # User models
│   ├── repositories/
│   │   ├── auth_repository.go   # Auth data access
│   │   └── user_repository.go   # User data access
│   ├── routes/
│   │   └── routes.go            # Route definitions
│   └── services/
│       ├── auth_service.go      # Authentication business logic
│       └── user_service.go      # User business logic
├── docker/
│   └── init-db.sql              # Database initialization
├── migrations/                   # Database migrations
├── scripts/                      # Utility scripts
├── tests/                        # Test files
├── docker-compose.yml           # Docker Compose configuration
├── Dockerfile                   # Docker image configuration
├── go.mod                       # Go module definition
└── README.md                    # This file
```

## 🔒 Security Features

- **JWT Authentication** - Secure token-based authentication
- **Password Hashing** - Bcrypt password hashing
- **Input Validation** - Request validation and sanitization
- **CORS Protection** - Cross-origin request handling
- **Security Headers** - Security-focused HTTP headers
- **SQL Injection Protection** - GORM ORM prevents SQL injection
- **Rate Limiting** - Request rate limiting (configurable)

## 🐳 Docker Support

### Build and run with Docker

```bash
# Build the image
docker build -t jsonplaceholder-api .

# Run the container
docker run -p 8080:8080 jsonplaceholder-api
```

### Using Docker Compose

```bash
# Start all services
docker-compose up -d

# View logs
docker-compose logs -f

# Stop services
docker-compose down

# Rebuild and start
docker-compose up --build -d
```

## 🔄 Development Workflow

1. **Make changes** to the code
2. **Run tests** to ensure everything works
3. **Build and test** locally with Docker
4. **Deploy** using Docker Compose

## 📊 Monitoring

- **Health Check**: `/health` endpoint for monitoring
- **Structured Logging**: JSON-formatted logs for easy parsing
- **Metrics**: Ready for integration with Prometheus/Grafana

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Ensure all tests pass
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [JSONPlaceholder](https://jsonplaceholder.typicode.com/) for API inspiration
- [Gin Framework](https://gin-gonic.com/) for the excellent HTTP framework
- [GORM](https://gorm.io/) for the fantastic ORM
