FROM golang:1.26 AS builder

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o wiki-job .

FROM alpine:latest

RUN adduser -D appuser
USER appuser

WORKDIR /usr/src/app
COPY --from=builder /app/wiki-job .

CMD ["./wiki-job"]
