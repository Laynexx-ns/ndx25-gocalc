FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY internal/ ./internal/
COPY pkg/ ./pkg/


RUN CGO_ENABLED=0 GOOS=linux go build -o /user-service ./internal/services/user-service/main.go


FROM alpine:latest

WORKDIR /app

COPY --from=builder /user-service .
COPY config/ ./config/
COPY db/migrations/ ./db/migrations/

EXPOSE 50053

CMD ["./user-service"]