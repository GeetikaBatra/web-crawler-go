package main

import(
   "log"
   "fmt"
   "os"
   "strings"
   "github.com/streadway/amqp"
   "web-crawler-go/crawlServer/links"
   
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

	
func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
	"crawl", // name
	false,   // durable
	false,   // delete when unused
	false,   // exclusive
	false,   // no-wait
	nil,     // arguments
)
failOnError(err, "Failed to declare a queue")
url := bodyFrom(os.Args)
err = ch.Publish(
  "",     // exchange	
  q.Name, // routing key
  false,  // mandatory
  false,  // immediate
  amqp.Publishing {
    ContentType: "text/plain",
    Body:        []byte(url),
  })
	failOnError(err, "Failed to publish a message")

urls, err := ch.Consume(
  q.Name, // queue
  "",     // consumer
  true,   // auto-ack
  false,  // exclusive
  false,  // no-local
  false,  // no-wait
  nil,    // args
)
failOnError(err, "Failed to register a consumer")

// forever := make(chan bool)
worklist := make(chan[] string)
go func(){ worklist <- string(urls)}()
seen := make(map[string] bool)
for list := range worklist{
	for _,link := range list{
		if !seen[link] {
			seen[link] = true
			go func(link string){
				worklist <- crawl(link)
			}(link)
		}
	}
}


fmt.Println("*******************	working")
 // <-forever		
}
func bodyFrom(args []string) string {
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		fmt.Println("Reached here")
		s = "hello"
	} else {
		s = strings.Join(args[1:], " ")
	}
	return s
}


func crawl(url string) []string {
	fmt.Println(url)
	list, err := links.Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

	
