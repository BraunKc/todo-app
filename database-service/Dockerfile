FROM golang:1.25.1-alpine AS builder

WORKDIR /app

COPY go.sum go.mod ./

RUN go mod download

COPY . .

RUN go build -o main ./cmd/todo-db/main.go

FROM alpine:latest

WORKDIR /db-service

COPY --from=builder /app/config/config.yml ./config/
COPY --from=builder /app/.env .

COPY --from=builder /app/main .

EXPOSE 50051

CMD ["./main"]
