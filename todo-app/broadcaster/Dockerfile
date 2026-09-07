# --- Build Stage ---
FROM golang:1.26 AS builder

WORKDIR /broadcaster

# Cache Go modules dependency layer
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the source code
COPY . .

# Build the static binary
RUN CGO_ENABLED=0 GOOS=linux go build -o todo-broadcaster .

# --- Final Stage ---
FROM alpine:latest

WORKDIR /app

# Create a non-root user and setup the data directory in one step
RUN adduser -D appuser && \
    mkdir -p /app/data && \
    chown -R appuser:appuser /app

# Switch to the non-root user
USER appuser

# Copy the binary and ensure appuser owns it
COPY --from=builder --chown=appuser:appuser /broadcaster/todo-broadcaster .

CMD ["./todo-broadcaster"]
