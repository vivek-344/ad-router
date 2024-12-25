# Build stage
FROM golang:1.23.4-alpine3.21 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Final application image
FROM alpine:3.21
WORKDIR /app
RUN apk add --no-cache bash
COPY --from=builder /app/main .
COPY app.env .
COPY start.sh .

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT [ "/app/start.sh" ]