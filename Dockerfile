FROM golang:latest

RUN go env

# Bot Token
ENV SPLAT_BOT_TOKEN
# 打包为release
ENV SPLAT_ENV release

ENV GO111MODULE auto

ENV GOPATH /go:/SplatBot

RUN mkdir -p /SplatBot

COPY . /SplatBot
WORKDIR /SplatBot

RUN go build .

ENTRYPOINT ["./SplatBot"]
