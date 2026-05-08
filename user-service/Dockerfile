FROM golang:1.25-alpine AS builder
WORKDIR /app
COPY go.work go.work.sum* ./
COPY gen/ gen/
COPY api-gateway/ api-gateway/
COPY user-service/ user-service/
COPY order-service/ order-service/
COPY booking-service/ booking-service/
COPY notification-service/ notification-service/
RUN go build -o /user-service ./user-service/cmd/main.go

FROM golang:1.25-alpine
COPY --from=builder /user-service /user-service
ENTRYPOINT ["/user-service"]
