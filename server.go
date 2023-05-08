package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type handler func(http.ResponseWriter, *http.Request)

// Util
// ---------------------
func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

// Handlers
// ---------------------
func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func makeSendHandler(q amqp.Queue, ch *amqp.Channel, ctx context.Context) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("Recieved Send Request")

		if r.Method != "POST" {
			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprintf(w, "Incorrect method, this endpoint needs a body")
		}

		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)

		reqBody, err := ioutil.ReadAll(r.Body)
		failOnError(err, "Failed to publish a message")

		err = ch.PublishWithContext(ctx,
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(reqBody),
			})
		failOnError(err, "Failed to publish a message")
		fmt.Fprintf(w, "Message Sent")
		log.Print("Finished Sending Message")
	}
}

// Setup
// ---------------------
func connect() (*amqp.Connection, *amqp.Channel, amqp.Queue) {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:53098/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	return conn, ch, q
}
func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, ch, q := connect()
	defer conn.Close()
	defer ch.Close()

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/send", makeSendHandler(q, ch, ctx))
	log.Print("Starting Server")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
