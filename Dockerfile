# Build stage
FROM golang:1.23.4-alpine3.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN wget -O migrate.tar.gz https://github.com/golang-migrate/migrate/releases/latest/download/migrate.linux-amd64.tar.gz
RUN tar -xzf migrate.tar.gz -C /app && rm migrate.tar.gz

# Final stage
FROM alpine:3.21 AS final
WORKDIR /app
RUN apk add --no-cache bash
COPY --from=builder /app/main .
COPY --from=builder /app/migrate .
COPY app.env .
COPY db/migration ./db/migration
COPY start.sh .
RUN chmod +x /app/start.sh
RUN chmod +x /app/migrate

EXPOSE 8080
ENTRYPOINT ["/app/start.sh"]
CMD ["/app/main"]