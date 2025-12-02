FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY . .
RUN apk add --no-cache git
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /api-news ./cmd

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=builder /api-news /api-news
EXPOSE 8081
ENTRYPOINT ["/api-news"]
