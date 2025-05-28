FROM golang:1.23-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Сборка основного приложения из ./cmd/main.go в бинарь с именем app
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o app ./cmd

FROM alpine:latest

RUN adduser -D appuser

# Копируем бинарь в финальный образ
COPY --from=builder /app/app /app/app

WORKDIR /app
USER appuser

EXPOSE 8080

ENTRYPOINT ["./app"]