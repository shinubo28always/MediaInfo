# Unrated Coder t.me/Unrated_Coder

# Build stage
FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o bot ./cmd/bot

# Run stage
FROM alpine:latest
RUN apk add --no-cache mediainfo ffmpeg
WORKDIR /app
COPY --from=builder /app/bot .
CMD ["./bot"]
