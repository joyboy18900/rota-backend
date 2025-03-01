# Step 1: Use Go 1.24-alpine as the base image for building the app
FROM golang:1.24-alpine AS builder

# Step 2: Install gcc and musl-dev (needed for cgo dependencies)
RUN apk add --no-cache gcc musl-dev

# Step 3: Set the working directory
WORKDIR /app

# Step 4: Copy go.mod and go.sum first to download dependencies (this avoids downloading dependencies every time)
COPY go.mod go.sum ./
RUN go mod download

# Step 5: Copy the entire source code into the container (volume mount will be better for dev)
COPY . .

# In development, use go run instead of go build to avoid rebuilding every time
CMD ["go", "run", "main.go"]