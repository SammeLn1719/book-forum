FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /book-forum ./cmd/server/main.go

FROM alpine:latest

WORKDIR /app
COPY --from=builder /book-forum .
COPY migrations ./migrations
COPY .env .
COPY wait-for-postgres.sh .  # Теперь файл существует

RUN chmod +x wait-for-postgres.sh

EXPOSE 8080
CMD ["./wait-for-postgres.sh", "postgres", "--", "./book-forum"]
