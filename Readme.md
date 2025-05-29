# Rota API

Backend service for Rota application.

## 🚀 Quick Start

### Prerequisites

- Docker & Docker Compose

### Setup

1. **Clone the repository**
   ```bash
   git clone <repository-url>
   cd rota-api
   ```

2. **Configure environment**
   ```bash
   cp .env.example .env
   # Update .env with your configuration
   ```

3. **Start services**
   ```bash
   docker-compose up -d
   ```

4. **Run database migrations**
   ```bash
   make migrate-up
   ```

5. **Start the application**
   ```bash
   # Development with hot reload
   air
   ```

## 🧪 Testing

```bash
# Run tests
go test -v ./...
```

## 🛠 Development Commands

```bash
# Format code
gofmt -w .

# Database migrations
make migrate-up    # Apply migrations
make migrate-down  # Rollback last migration

# Docker commands
docker-compose up -d     # Start services
docker-compose logs -f   # View logs
docker-compose down     # Stop services
```

## 🔐 Authentication

- JWT-based authentication
- Token expiration: 24 hours
- Uses middleware for token verification
- Requires token in header: `Authorization: Bearer <token>`
