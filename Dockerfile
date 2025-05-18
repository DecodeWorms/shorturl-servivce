# Step 1: Use an official Go image as the base image for building the app
FROM golang:1.24-alpine AS builder

# Step 2: Set the working directory inside the builder container
WORKDIR /app

# Step 3: Copy only the Go modules files first to take advantage of Docker layer caching
COPY go.mod go.sum ./

# Step 4: Download and cache dependencies to optimize the build process
RUN go mod download

# Step 5: Copy the remaining application code into the container
COPY . .

# Step 6: Build the Go application with optimizations (static binary for Alpine)
RUN go build -ldflags="-w -s" -o main .

# Step 7: Debugging step: List the contents of the /app directory
RUN ls -la /app

# Step 8: Use a minimal lightweight base image for the final container
FROM alpine:latest

# Step 9: Install necessary certificates for HTTPS (if your app connects to external APIs or uses HTTPS)
RUN apk --no-cache add ca-certificates

# Step 10: Set the working directory for the final container
WORKDIR /app

# Step 11: Copy the compiled binary from the builder stage to the final container
COPY --from=builder /app/main .

# Step 13: Expose the application port
EXPOSE 8001

# Step 14: Set the default command to run the application
CMD ["./main"]
