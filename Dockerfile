FROM golang:1.24-alpine

# Install build dependencies
RUN apk add --no-cache gcc musl-dev git

WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Expose port
EXPOSE 8080

# Command to run the application
CMD ["go", "run", "main.go"]