# Этап 1: сборка приложения
FROM golang:1.25-alpine AS builder

WORKDIR /src

# Сначала копируем файлы зависимостей.
# Этот слой будет переиспользован, пока go.mod/go.sum не изменились.
COPY go.mod go.sum ./
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем backend-приложение
RUN go build -o /out/backend ./internal/cmd


# Этап 2: финальный runtime-образ
FROM alpine:3.22

# Создаём отдельного пользователя для запуска backend
RUN addgroup -S appgroup && \
    adduser -S -G appgroup appuser

# Рабочая директория приложения
WORKDIR /app


COPY --from=builder --chown=appuser:appgroup /out/backend /app/backend
COPY --from=builder --chown=appuser:appgroup /src/internal/config/docker.yaml /app/internal/config/docker.yaml
COPY --from=builder --chown=appuser:appgroup /src/static /app/static

ENV CONFIG_PATH=/app/internal/config/docker.yaml


# Backend не должен работать от root
USER appuser

ENTRYPOINT ["/app/backend"]