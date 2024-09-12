# Используем официальный легкий образ с Go
FROM golang:1.20-alpine

# Создаем директорию для приложения
WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o /app/main cmd/main.go

# Открываем порт
EXPOSE ${TODO_PORT}

ENTRYPOINT ["/app/main"]