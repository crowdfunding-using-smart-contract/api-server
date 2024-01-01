# Stage 1 - Builder: Import the golang container.
FROM golang:1.20-buster AS builder

# Set the work directory.
WORKDIR /app

# Copy go mod and sum files.
COPY go.mod .
COPY go.sum .

# Install the dependencies.
RUN go mod download

# Copy the source code into the container.
COPY . .

# Build the source code
RUN CGO_ENABLED=0 go build -o out/app cmd/api/main.go

# Stage 2 - Runner.
FROM alpine:3.16.2
COPY --from=builder /app/out/app .

EXPOSE 8080 9615
CMD ["/app"]
