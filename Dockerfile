FROM golang:1.24-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o app ./cmd

FROM alpine:latest

RUN adduser -D appuser

COPY --from=builder /app/app /app/app

WORKDIR /app
USER appuser

EXPOSE 8080

ENTRYPOINT ["./app"]