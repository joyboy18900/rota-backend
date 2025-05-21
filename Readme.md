# Rota API Backend (Go)

Backend service for the Rota application built with Go, Fiber, and PostgreSQL.

## 🚀 Features

- 🔐 JWT Authentication
- 🗄️ PostgreSQL Database
- 🚀 High Performance with Fiber
- 🔄 Redis for Token Management
- 📝 API Documentation (Swagger)
- 🔄 Database Migrations

## 📋 Prerequisites

- Go 1.20+
- PostgreSQL 13+
- Redis 6+
- Make (optional)

## 🛠️ Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/rota-backend.git
   cd rota-backend
   ```

2. Copy the environment file and update the values:
   ```bash
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. Install dependencies:
   ```bash
   go mod download
   ```

## 🚀 Running the Application

1. Start PostgreSQL and Redis services

2. Run database migrations:
   ```bash
   make migrate-up
   # or manually:
   # migrate -path ./migrations -database "postgres://user:pass@localhost:5432/dbname?sslmode=disable" up
   ```

3. Start the application:
   ```bash
   go run main.go
   # or with hot reload (install air first: go install github.com/cosmtrek/air@latest)
   air
   ```

The server will start on `http://localhost:8080` by default.

## 📚 API Documentation

After starting the server, access the Swagger documentation at:
- Swagger UI: http://localhost:8080/swagger/index.html
- OpenAPI JSON: http://localhost:8080/swagger/doc.json

## 🧪 Running Tests

```bash
# Run all tests
go test -v ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out
```

## 🛠 Development

### Code Generation

```bash
# Generate mocks
mockgen -source=path/to/interface.go -destination=mocks/interface_mock.go -package=mocks
```

### Linting

```bash
# Install golangci-lint if needed
# https://golangci-lint.run/usage/install/

golangci-lint run
```

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

Distributed under the MIT License. See `LICENSE` for more information.

## 🙏 Acknowledgments

- [Fiber](https://gofiber.io/) - Web framework
- [GORM](https://gorm.io/) - ORM library
- [Swaggo](https://github.com/swaggo/swag) - API documentation
- [Testify](https://github.com/stretchr/testify) - Testing toolkit

### Contact

Your Name - [your.email@example.com](mailto:your.email@example.com)

Project Link: [https://github.com/yourusername/rota-backend](https://github.com/yourusername/rota-backend)