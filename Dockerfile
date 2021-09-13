FROM golang:alpine AS builder

COPY . /order_book/cmd
WORKDIR /order_book/cmd

RUN go mod download
RUN go build -o fetcher cmd/main.go

FROM alpine:3.12

COPY --from=0 ./order_book/cmd/fetcher .
COPY --from=0 ./order_book/cmd/configs ./configs

CMD ["./fetcher"]
