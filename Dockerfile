FROM golang:alpine AS builder

COPY . /order_book/cmd
WORKDIR /order_book/cmd

RUN go mod download
RUN go build -o service cmd/main.go

FROM alpine:3.12

COPY --from=0 ./order_book/cmd/service .
COPY --from=0 ./order_book/cmd/configs ./configs

CMD ["./service"]
