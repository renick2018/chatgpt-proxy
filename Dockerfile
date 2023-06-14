FROM golang:alpine AS builder

LABEL stage=gobuilder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOPROXY=https://goproxy.cn,https://goproxy.io,direct

WORKDIR /build

COPY . .
RUN go mod download

RUN go build -o /app/chatgpt-proxy .

FROM alpine

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories

RUN apk update --no-cache && apk add --no-cache ca-certificates tzdata
ENV TZ Asia/Shanghai

WORKDIR /app
COPY --from=builder /app/chatgpt-proxy /app/chatgpt-proxy
COPY --from=builder /build/.conf.yml /app/.conf.yml

VOLUME /dir

EXPOSE 8088

ENTRYPOINT ["./chatgpt-proxy"]
