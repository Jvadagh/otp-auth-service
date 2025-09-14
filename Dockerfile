FROM golang:1.24.1-alpine AS builder
WORKDIR /app

RUN apk add --no-cache git ca-certificates

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o otp-auth-service cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/otp-auth-service .

COPY migrations ./migrations
COPY internal/docs ./internal/docs
COPY .env .env

RUN apk add --no-cache ca-certificates

EXPOSE 8080

CMD ["./otp-auth-service"]
