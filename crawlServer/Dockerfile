FROM golang

RUN go get github.com/streadway/amqp

WORKDIR /crawlServer
RUN ls /crawlServer
RUN go build -o crawl1 . 
