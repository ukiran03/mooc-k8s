FROM golang:1.26 AS builder

WORKDIR /usr/src/app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o todo-app .

# Final Stage
FROM alpine:latest

WORKDIR /app

# Copy the binary and UI files
COPY --from=builder /usr/src/app/ui ./ui
COPY --from=builder /usr/src/app/todo-app .

CMD ["./todo-app"]
