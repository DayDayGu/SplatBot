FROM golang:latest

RUN go env

# Bot Token
ENV SPLAT_BOT_TOKEN 828734497:AAGczERxUIf2Su5hvcrVi7KBg7qd_MzKUXI
# 打包为release
ENV SPLAT_ENV release 

RUN go get -u github.com/golang/dep/cmd/dep
RUN mkdir /go/src/github.com/PangPangPangPangPang

COPY . /go/src/github.com/PangPangPangPangPang/SplatBot
WORKDIR /go/src/github.com/PangPangPangPangPang/SplatBot

RUN dep ensure -v
RUN go build .

ENTRYPOINT ["./SplatBot"]
