FROM golang:1.24.1-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git ca-certificates bash postgresql-client
RUN go install github.com/swaggo/swag/cmd/swag@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN swag init -g cmd/main.go -o internal/docs
RUN go build -o otp-auth-service cmd/main.go
RUN go build -o migrate cmd/migrate/main.go

FROM alpine:latest
WORKDIR /app

RUN apk add --no-cache bash postgresql-client ca-certificates

COPY --from=builder /app/otp-auth-service .
COPY --from=builder /app/migrate .
COPY --from=builder /app/migrations ./migrations
COPY --from=builder /app/internal/docs ./internal/docs
COPY --from=builder /app/.env .env

EXPOSE 8080

CMD bash -c "./migrate up && ./otp-auth-service"
