# web-crawler-go
A distributed web crawler

### The Crawler is based on microservice architecture and has following components:

#### crawl-broker : 
Uses RabbitMq 
#### crawl-server:
Runs the webcrawler code based on golang
#### cassandra-docker:
Runs Cassandra Database
#### Janus:
It is the Graph Database layer on top of cassandra. Janus uses gremlin query language to insert data into database


### Steps to run the service
##### Clone the repository 

`https://github.com/GeetikaBatra/web-crawler-go.git`

##### Go to the directory web-crawler-go

`cd web-crawler-go`

#### Build the service

`sh docker-compose.sh build `

#### Run the service

`sh docker-compose.sh up`


Graph Database is available at port `8182` which is exposed. 
To check the number of nodes available in the graph, following can be done

`GET localhost:8182?gremlin=g.V().count()`
