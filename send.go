package main

import (
	"github.com/streadway/amqp"
	"log"
	"time"
)

func failOnErrorSend(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func send(str string) {

		conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
		failOnErrorSend(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		ch, err := conn.Channel()
		failOnErrorSend(err, "Failed to open a channel")
		defer ch.Close()

		q, err := ch.QueueDeclare(
			"hello", // name
			false,   // durable
			false,   // delete when unused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)
		failOnErrorSend(err, "Failed to declare a queue")

		body :=  time.Now().Format("15:04:05") + ": " + str
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		failOnErrorSend(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s", body)

}