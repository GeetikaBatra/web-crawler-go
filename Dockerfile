FROM golang:1.7

RUN apt-get -y update && apt-get -y install netcat
RUN mkdir -p /crawlServer
WORKDIR crawlServer/
ADD ./main .
ADD ./scripts/server-entryscript.sh .
CMD ["./server-entryscript.sh"]