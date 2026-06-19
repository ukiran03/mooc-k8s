FROM golang:1.26 AS builder

WORKDIR /usr/src/app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o todo-app .

# Final Stage
FROM alpine:latest

# Create a system group and user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

# Copy the binary and UI files
COPY --from=builder /usr/src/app/ui ./ui
COPY --from=builder /usr/src/app/todo-app .

# Ensure the appuser owns the directory if your app needs to write files (optional)
# RUN chown -R appuser:appgroup /app

# Switch to the non-root user
USER appuser

CMD ["./todo-app"]
