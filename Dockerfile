FROM golang:1.23-alpine AS builder

# Установка необходимых зависимостей для сборки
RUN apk add --no-cache git make

# Создание рабочей директории
WORKDIR /app

# Копирование go.mod и go.sum для предварительной загрузки зависимостей
COPY go.mod go.sum ./
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения с оптимизациями
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags="-s -w" -o job_solution ./cmd/api

# Второй этап: минимальный образ для запуска
FROM alpine:latest

# Установка необходимых пакетов для работы приложения
RUN apk --no-cache add ca-certificates tzdata

# Создание директории для миграций
RUN mkdir -p /app/internal/db/migrations

# Настройка рабочей директории
WORKDIR /app

# Копирование бинарного файла из первого этапа
COPY --from=builder /app/job_solution .

# Создание пользователя без привилегий для безопасности
RUN adduser -D -H -h /app appuser
RUN chown -R appuser:appuser /app
USER appuser

# Открываем порт приложения
EXPOSE 8080

# Запуск приложения
CMD ["./job_solution"] 