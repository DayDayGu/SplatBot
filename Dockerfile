FROM golang:latest

RUN go env

# Bot Token
ENV SPLAT_BOT_TOKEN
# 打包为release
ENV SPLAT_ENV release

ENV GO111MODULE auto

ENV GOPATH /go:/root/SplatBot

RUN mkdir -p /root/SplatBot

COPY . /root/SplatBot
WORKDIR /root/SplatBot

RUN go build .

ENTRYPOINT ["./SplatBot"]
