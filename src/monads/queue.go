/*
 * For interacting with queues
 */
package main

import (
    "context"
    "time"

    amqp "github.com/rabbitmq/amqp091-go"
)

type Queue struct {
    ch   amqp.Channel
    conn amqp.Connection
    q    amqp.Queue
}

func (q *Queue) Close () {
    q.conn.Close()
    q.ch.Close()
}

// Connect to a queue/channel
func connect(qname string) Queue {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:53261/")
	failOnError(err, "Failed to connect to RabbitMQ")

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")

	q, err := ch.QueueDeclare(
		qname, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Failed to declare a queue")

	return Queue{*ch, *conn, q}
}

func (q *Queue) send(msg []byte) error {
	ctx, cancel := context.WithTimeout(context.Background(), 
                                       5 * time.Second)
	defer cancel()

    return q.ch.PublishWithContext(ctx,
        "",       // exchange
        q.q.Name, // routing key
        false,    // mandatory
        false,    // immediate
        amqp.Publishing{
            ContentType: "text/plain",
            Body:        msg,
        })
}
