FROM golang:1.24-alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY internal/ ./internal/
COPY pkg/ ./pkg/


RUN CGO_ENABLED=0 GOOS=linux go build -o /agent-service ./internal/services/agent/main.go


FROM alpine:latest

WORKDIR /app

COPY --from=builder /agent-service .

EXPOSE 50052

CMD ["./agent-service"]