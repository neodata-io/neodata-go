# Stage 1: Build the Go application using an official Go image
FROM golang:1.22-alpine AS builder

# Set environment variables for Go module and proxy
ENV GOPROXY=https://proxy.golang.org,direct
ENV CGO_ENABLED=0

# Set the working directory inside the container for the build process
WORKDIR /app

# Copy go.mod and go.sum files to download dependencies first
COPY go.mod go.sum ./

# Download the Go module dependencies
RUN go mod download

# Copy the rest of the source code from the host machine to the container
COPY . .

# Build the Go application binary
RUN go build -o /neodata-go

# Stage 2: Create a lightweight runtime image using Alpine
FROM alpine:latest

# Set the working directory inside the runtime container
WORKDIR /

# Copy the built binary from the builder stage
COPY --from=builder /neodata-go /neodata-go

# Copy Casbin configuration file to the runtime container
COPY ./config/casbin/rbac_model.conf /app/config/casbin/rbac_model.conf
COPY ./config/casbin/policy.csv /app/config/casbin/policy.csv

# Expose the application port to the host
EXPOSE 8080

# Specify the command to run the Go application when the container starts
CMD ["/neodata-go"]
