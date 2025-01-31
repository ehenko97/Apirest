# Этап сборки
FROM golang:1.23.4 AS builder

# Рабочая директория
WORKDIR /app

# Установка goose (миграции)
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

# Копирование файлов зависимостей
COPY go.mod go.sum ./

# Установка зависимостей
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения
RUN GOARCH=arm64 go build -o Projectapirest ./cmd/main.go

# Этап запуска
FROM alpine:latest

# Рабочая директория
WORKDIR /app

# Установка необходимых пакетов
RUN apk add --no-cache ca-certificates tzdata libc6-compat

# Копирование бинарного файла и goose
COPY --from=builder /app/Projectapirest /app/Projectapirest
COPY --from=builder /go/bin/goose /usr/local/bin/goose

# Копирование конфигураций и миграций
COPY config ./config
COPY migrations ./migrations

# Права на выполнение
RUN chmod +x /app/Projectapirest /usr/local/bin/goose

# Открытие портов
EXPOSE 8080
EXPOSE 50051

# Запуск приложения
CMD ["./Projectapirest"]