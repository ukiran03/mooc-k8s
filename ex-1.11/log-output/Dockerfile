FROM golang:1.26 AS builder

WORKDIR /usr/src/app

COPY . .

RUN go mod download && go mod verify

RUN CGO_ENABLED=0 GOOS=linux go build -o log-output .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /usr/src/app/log-output .

CMD ["./log-output"]
