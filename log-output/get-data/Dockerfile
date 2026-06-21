FROM golang:1.26 AS builder

WORKDIR /usr/src/app

COPY . .

RUN go mod download && go mod verify

RUN CGO_ENABLED=0 GOOS=linux go build -o get-data .

FROM alpine:latest

WORKDIR /app

COPY --from=builder /usr/src/app/get-data .

CMD ["./get-data"]
