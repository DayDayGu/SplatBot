FROM golang:latest

# Bot Token
ENV SPLAT_BOT_TOKEN 
# 打包为release
ENV SPLAT_ENV release 

WORKDIR /root/SplatBot
COPY . /root/SplatBot

RUN go get gopkg.in/tucnak/telebot.v2
RUN go get github.com/mattn/go-sqlite3

RUN go build .

EXPOSE 8090

ENTRYPOINT ["./SplatBot"]
