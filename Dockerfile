FROM golang:latest

RUN go env

# Bot Token
ENV SPLAT_BOT_TOKEN
# 打包为release
ENV SPLAT_ENV release

RUN mkdir -p /go/src/github.com/PangPangPangPangPang

COPY . /go/src/github.com/PangPangPangPangPang/SplatBot
WORKDIR /go/src/github.com/PangPangPangPangPang/SplatBot

RUN go build .

ENTRYPOINT ["./SplatBot"]
