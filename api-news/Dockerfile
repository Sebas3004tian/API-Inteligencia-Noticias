FROM golang:1.24.4-alpine AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o main ./cmd

FROM alpine:3.19
WORKDIR /app

EXPOSE 8080
COPY --from=builder /app/main .

CMD ["./main"]
