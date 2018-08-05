package main

import(
   "log"
   "fmt"
   "os"
   "strings"
   "github.com/streadway/amqp"

   
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


// forever := make(chan bool)
seen := make(map[string] bool)
for{
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

for list := range urls{
	for _,link := range list{
		if !seen[link] {
			seen[link] = true
			go func(link string){

			url_list = crawl(link)

			for _, url:= range url_list{
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
			}
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
}

func crawl(url string) []string {
	fmt.Println(url)
	list, err := Extract(url)
	if err != nil {
		log.Print(err)
	}
	return list
}

	
func Extract(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, fmt.Errorf("getting %s: %s", url, resp.Status)
	}

	doc, err := html.Parse(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, fmt.Errorf("parsing %s as HTML: %v", url, err)
	}

	var links []string
	visitNode := func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "a" {
			for _, a := range n.Attr {
				if a.Key != "href" {
					continue
				}
				link, err := resp.Request.URL.Parse(a.Val)
				if err != nil {
					continue // ignore bad URLs
				}
				links = append(links, link.String())
			}
		}
	}
	forEachNode(doc, visitNode, nil)
	return links, nil
}

//!-Extract

// Copied from gopl.io/ch5/outline2.
func forEachNode(n *html.Node, pre, post func(n *html.Node)) {
	if pre != nil {
		pre(n)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		forEachNode(c, pre, post)
	}
	if post != nil {
		post(n)
	}
}