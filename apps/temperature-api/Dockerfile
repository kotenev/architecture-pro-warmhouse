# Этап сборки
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GO111MODULE=auto go build -o temperature-api .

# Этап выполнения
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/temperature-api .
EXPOSE 8081
CMD ["./temperature-api"]
