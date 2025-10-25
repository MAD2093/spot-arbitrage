# Этап сборки
FROM golang:1.24.3-alpine AS builder

# Аргументы для production сборки
ARG BUILD_ENV=production
ARG APP_PATH=./cmd/app/main.go

# Установка необходимых зависимостей для сборки
RUN apk add --no-cache git openssh

# Установка рабочей директории
WORKDIR /app

# Копирование файлов зависимостей
COPY go.mod go.sum ./

# Копируем SSH-ключ и настраиваем SSH
RUN mkdir -p /root/.ssh && \
    chmod 700 /root/.ssh
COPY id_rsa /root/.ssh/id_rsa
RUN chmod 600 /root/.ssh/id_rsa && \
    ssh-keyscan github.com >> /root/.ssh/known_hosts

# Настраиваем git для использования SSH вместо HTTPS
RUN git config --global url."git@github.com:".insteadOf "https://github.com/"

# Скачиваем зависимости
RUN go mod download

# Копирование исходного кода
COPY . .

# Сборка приложения с флагами для разных окружений
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo \
    -ldflags "-X main.buildEnv=${BUILD_ENV}" \
    -o main ${APP_PATH}

# Финальный этап
FROM alpine:latest

# Аргумент для runtime
ARG BUILD_ENV=production
ENV APP_ENV=${BUILD_ENV}

# Установка необходимых пакетов для работы
RUN apk --no-cache add ca-certificates tzdata curl

WORKDIR /root/

# Копирование бинарного файла из этапа сборки
COPY --from=builder /app/main .

# Открытие порта (замените на нужный вам порт)
EXPOSE 8081

# Healthcheck (может отличаться для test/prod)
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD curl -f http://localhost:8081/metrics || exit 1

# Запуск приложения
CMD ["./main"]