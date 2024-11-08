# Use the official Go image
FROM golang:1.22

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

# Expose the application port
EXPOSE 8080

# Run the executable
CMD ["./main"]