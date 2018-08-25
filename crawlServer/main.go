package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/streadway/amqp"
	"golang.org/x/net/html"
)

type RbmqConfig struct {
	q       amqp.Queue
	ch      *amqp.Channel
	conn    *amqp.Connection
	rbmqErr error
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

// visit appends to links each link found in n, and returns the result.
func visit(config *RbmqConfig, links []string, n *html.Node, baseUrl string, seen map[string]bool) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {

				if strings.HasPrefix(a.Val, "/") {
					url_ := baseUrl + a.Val
					links = append(links, url_)
					if !seen[url_] {
						seen[url_] = true
						publishMessages(config, url_)
					}

				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = visit(config, links, c, baseUrl, seen)
	}
	return links
}
func initAmqp() *RbmqConfig {
	config := &RbmqConfig{}
	config.conn, config.rbmqErr = amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(config.rbmqErr, "Failed to connect to RabbitMQ")

	log.Printf("got Connection, getting Channel...")

	config.ch, config.rbmqErr = config.conn.Channel()
	failOnError(config.rbmqErr, "Failed to open a channel")

	log.Printf("got Channel, declaring Exchange (%s)", "go-crawl-exchange")

	log.Printf("declared Exchange, declaring Queue (%s)", "go-crawl-queue")

	config.q, config.rbmqErr = config.ch.QueueDeclare(
		"go-crawl-queue", // name, leave empty to generate a unique name
		true,             // durable
		false,            // delete when usused
		false,            // exclusive
		false,            // noWait
		nil,              // arguments
	)
	failOnError(config.rbmqErr, "Error declaring the Queue")

	return config
}

func consumeMessages(config *RbmqConfig, baseUrl string, seen map[string]bool) {
	var err error
	msgs, err := config.ch.Consume(
		"go-crawl-queue", // queue
		"",               // consumer
		true,             // auto-ack
		false,            // exclusive
		false,            // no-local
		false,            // no-wait
		nil,              // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		fmt.Println("************************REAhed here")
		for url := range msgs {
			fmt.Println(url)
			fmt.Println(string(url.Body[:]))
			links, err := findLinks(config, string(url.Body[:]), baseUrl, seen)
			if err != nil {
				fmt.Fprintf(os.Stderr, "findlinks2: %v\n", err)
			}

			fmt.Println(strings.Join(links[:], ","))
			log.Printf("Received a message: %s", url.Body)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func publishMessages(config *RbmqConfig, url string) {
	config.rbmqErr = config.ch.Publish(
		"",               // exchange
		"go-crawl-queue", // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			DeliveryMode: amqp.Transient,
			ContentType:  "text/plain",
			Body:         []byte(url),
			Timestamp:    time.Now(),
		})
	failOnError(config.rbmqErr, "Failed to Publish on RabbitMQ")
}

//!+
func main() {
	baseUrl := os.Args[1]
	config := initAmqp()
	seen := make(map[string]bool)
	seen[baseUrl] = true
	publishMessages(config, baseUrl)
	consumeMessages(config, baseUrl, seen)
	printMap(seen)
	defer config.conn.Close()

}

func printMap(seen map[string]bool) {

	fmt.Println(seen)
}

// findLinks performs an HTTP GET request for url, parses the
// response as HTML, and extracts and returns the links.
func findLinks(config *RbmqConfig, url string, baseUrl string, seen map[string]bool) ([]string, error) {
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
	return visit(config, nil, doc, baseUrl, seen), nil
}
