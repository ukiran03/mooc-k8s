FROM golang:1.26 AS builder

WORKDIR /usr/src/app

# Copy only the dependency files first
COPY go.mod ./

RUN go mod download && go mod verify

# Copy the rest of the source code (this changes often)
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o log-output .

# Final Stage
FROM alpine:latest

WORKDIR /app

COPY --from=builder /usr/src/app/log-output .

CMD ["./log-output"]
