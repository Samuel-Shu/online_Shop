FROM golang:1.22-alpine3.20 AS builder


WORKDIR /build

COPY . .

ENV GOPROXY=https://goproxy.cn,direct

RUN go mod tidy

RUN go build -o goods_srv main.go


FROM alpine:latest

WORKDIR /run

COPY --from=builder /build/goods_srv .

CMD ["./goods_srv"]
