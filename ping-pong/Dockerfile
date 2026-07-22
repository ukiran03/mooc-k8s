FROM golang:1.26 AS builder

WORKDIR /app

# Copy only the dependency files first
COPY go.mod ./

RUN go mod download && go mod verify

# Copy the rest of the source code (this changes often)
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o ping-pong .

# Final Stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/ping-pong .

CMD ["./ping-pong"]
