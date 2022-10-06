FROM golang:1.18 AS builder
WORKDIR /build

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./src/main.go

FROM alpine AS server

WORKDIR /app

COPY --from=builder /build/main .

CMD ["./main"]
