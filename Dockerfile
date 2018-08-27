FROM golang:1.7

RUN mkdir -p /crawlServer
WORKDIR /crawlServer
ADD ./crawlServer/main /crawlServer
ENTRYPOINT ./main