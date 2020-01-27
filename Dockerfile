# Start from the alpine golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Abdullahi Innocent <deewai48@gmail.com>"

# Set the Current Working Directory inside the container
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source from the current directory to the Working Directory inside the container
COPY . .

# Run the tests
RUN go test ./app/

# Build the Go app
RUN go build -o main .

# Expose port 8000
EXPOSE 3000

# Command to run the executable
CMD ["./main"]