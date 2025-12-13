# Multi-stage build
FROM golang:1.21 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o github.com/dhanarrizky/Golang-template cmd/app/main.go

FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/github.com/dhanarrizky/Golang-template .
CMD ["./github.com/dhanarrizky/Golang-template"]