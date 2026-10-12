# Этап 1. Сборка приложения
FROM golang:1.26 AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o task-manager ./cmd/app


# Этап 2. Запуск приложения
FROM alpine:3.22

WORKDIR /app

COPY --from=builder /app/task-manager ./task-manager

USER 10001:10001

EXPOSE 8080

CMD ["./task-manager"]