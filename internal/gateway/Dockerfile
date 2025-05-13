FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY internal/ ./internal/
COPY pkg/ ./pkg/


RUN CGO_ENABLED=0 GOOS=linux go build -o /gateway ./internal/gateway/main.go


FROM alpine:latest

WORKDIR /app

COPY --from=builder /gateway .

EXPOSE 8080

CMD ["./gateway"]