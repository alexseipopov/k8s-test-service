# Многоэтапная сборка для оптимизации размера образа
FROM --platform=linux/amd64 golang:1.24-alpine AS builder

# Устанавливаем необходимые пакеты
RUN apk add --no-cache git ca-certificates tzdata

# Создаем рабочую директорию
WORKDIR /app

# Копируем go.mod и go.sum для кэширования зависимостей
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main .

# Финальный образ
FROM --platform=linux/amd64 alpine:latest

# Устанавливаем ca-certificates для HTTPS запросов
RUN apk --no-cache add ca-certificates tzdata

# Создаем рабочую директорию и пользователя
WORKDIR /app
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# Копируем бинарный файл из builder stage
COPY --from=builder /app/main .

# Меняем владельца файла и директории, даем права на выполнение
RUN chown -R appuser:appgroup /app && chmod +x main

# Переключаемся на непривилегированного пользователя
USER appuser

# Открываем порт (если понадобится в будущем)
EXPOSE 8080

# Команда запуска
CMD ["./main"]
