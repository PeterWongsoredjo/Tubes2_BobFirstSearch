# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire backend source code into the container
COPY . .

# Build the binaries for cmd (server) and main
RUN CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server.go ./cmd/mainScrape.go
RUN CGO_ENABLED=0 GOOS=linux go build -o /main ./main/TreeHandler.go ./main/BFS.go ./main/DFS.go ./main/BFSplit.go ./main/utils.go

# Run stage (production)
FROM alpine:latest
WORKDIR /

# Copy the built binaries from the builder stage
COPY --from=builder /server /server
COPY --from=builder /main /main

# Copy configs directory
COPY configs /configs/

# Expose ports (8080 for main, 8081 for cmd)
EXPOSE 8080
EXPOSE 8081

# Run both binaries in the background
CMD /server & /main

