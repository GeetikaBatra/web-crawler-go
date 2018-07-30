FROM golang

RUN go get github.com/streadway/amqp

WORKDIR /crawlserver 
RUN go build -o crawl1 . 
CMD ["/app/main"]