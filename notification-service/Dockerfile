FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY notification-service/ notification-service/
WORKDIR /app/notification-service
RUN go build -o /notification-service ./cmd/main.go

FROM golang:1.25-alpine
COPY --from=builder /notification-service /notification-service
ENTRYPOINT ["/notification-service"]
