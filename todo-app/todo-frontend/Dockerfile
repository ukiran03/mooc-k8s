# --- Build Stage ---
FROM golang:1.26 AS builder

WORKDIR /frontend

# Copy the rest of the source code
COPY . .

# Build the static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o todo-frontend .

# --- Final Stage ---
FROM alpine:3.20

WORKDIR /app

# Create a non-root user for security
RUN adduser -D appuser && chown -R appuser:appuser /app
USER appuser

# Copy the binary and Image files from the builder
COPY --from=builder --chown=appuser:appuser /frontend/images ./images
COPY --from=builder --chown=appuser:appuser /frontend/todo-frontend .

CMD ["./todo-frontend"]
