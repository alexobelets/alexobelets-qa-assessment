# Stage 1: Build the application
FROM golang:1.22.3-alpine AS application

# Set the working directory in the container to /app
WORKDIR /app
# Copy the files from the current directory to /app in the container
COPY go.mod go.sum ./
# Install the dependencies
RUN go mod download
# Copy the files from the current directory to /app in the container
COPY . .
# Build the application and save it with name after -o flag
RUN go build -o qa-challenge-application .

# Stage 2: Reduce the size of image and Run the application
FROM alpine:latest
# Install necessary bash commands
RUN apk add --no-cache bash
RUN apk add --no-cache make
# Set the working directory in the container to /app
WORKDIR /app
# Copy the files from the current directory to /app in the container
COPY --from=application /app/qa-challenge-application ./app/
COPY config.toml Makefile run_app.sh ./
# Make the run_app.sh script executable (unix line ednings)
RUN chmod +x ./run_app.sh
# Set the entry point for the Docker container
CMD ["./run_app.sh"]
