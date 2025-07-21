# Build image
FROM golang:1.23-alpine AS builder

# Install necessary tools
RUN apk add --no-cache git make

# Install swagger CLI
RUN go install github.com/go-swagger/go-swagger/cmd/swagger@latest

# Install SQLC cli
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest


# Set the working directory
WORKDIR /app

# Copy the entire project
COPY . .

RUN touch .env
RUN make swag
RUN make sqlc

# Set working directory to src
WORKDIR /app/src

# Download all dependencies
RUN go mod download

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server cmd/server/main.go
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o migrate cmd/migrate/main.go

# Start a new stage from scratch
FROM alpine:latest

WORKDIR /app

# Binary execute file
RUN touch .env
COPY --from=builder /app/src/server .
COPY --from=builder /app/src/static/ ./static

# Migration
COPY --from=builder /app/src/migrate /app/
COPY --from=builder /app/src/infrastructure/migrations /app/migrations

# Expose port 8080 to the outside world
EXPOSE 8000

COPY --from=builder /app/entrypoint.sh /
RUN chmod +x /entrypoint.sh
ENTRYPOINT ["/entrypoint.sh"]

CMD ["/bin/sh"]