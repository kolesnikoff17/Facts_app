FROM golang:1.18 AS builder
WORKDIR /build

COPY /src /build

RUN go build . -o main

FROM alpine AS server

COPY --from=builder /build/main /opt/app/

WORKDIR /opt/app

CMD ["./main"]
