FROM golang:1.7

RUN mkdir -p /crawlServer
WORKDIR /crawlServer
ADD ./crawl1 /crawlServer
ENTRYPOINT ./crawl1