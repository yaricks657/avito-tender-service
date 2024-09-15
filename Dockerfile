# Используем официальный легкий образ с Go
FROM golang:1.22-alpine

# Создаем директорию для приложения
WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /app/main cmd/main.go

# Определяем переменную окружения для слушающего адреса
ENV SERVER_ADDRESS="0.0.0.0:8080"

# Открываем порт
EXPOSE 8080

ENTRYPOINT ["/app/main"]