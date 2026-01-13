# Stage 1: Build
FROM golang:1.25-alpine AS builder

# Install git and certificates (needed for some go modules)
RUN apk add --no-cache git

WORKDIR /app

# 1. Copy go.mod and go.sum first to leverage Docker cache
# If you haven't run 'go mod init' yet, the next step handles it
COPY go.mod* go.sum* ./
RUN if [ ! -f go.mod ]; then go mod init smartload; fi
RUN go mod tidy

# 2. Copy the rest of the source code
COPY . .

# 3. Build - we use "." because main.go is in the root
RUN CGO_ENABLED=0 GOOS=linux go build -o /smartload .

# Stage 2: Final
FROM alpine:latest
RUN apk --no-cache add ca-certificates

# Create a non-privileged user for security (Best Practice)
RUN adduser -D appuser
USER appuser

WORKDIR /home/appuser/
COPY --from=builder /smartload .

# Requirement: Must listen on port 8080
EXPOSE 8080

CMD ["./smartload"]