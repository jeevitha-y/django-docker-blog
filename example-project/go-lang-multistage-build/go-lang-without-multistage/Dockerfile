# Use official Go image
FROM golang:1.23

# Set working directory inside container
WORKDIR /app

# Copy go.mod and go.sum first (better caching)
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy remaining project files
COPY . .

# Build the application
RUN go build -o app

# Expose application port
EXPOSE 8080

# Run the application
CMD ["./app"]