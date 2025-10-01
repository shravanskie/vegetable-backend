# Use official Go image
FROM golang:1.25.1

# Set working directory inside container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .
# Build binary for Linux ARM64
RUN GOOS=linux GOARCH=arm64 go build -o /app/main ./cmd/main.go

RUN ls -lh /app

CMD ["/app/main"]