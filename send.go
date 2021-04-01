package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func send(str string) {

		conn, err := amqp.Dial(URL)
		if err != nil {
			log.Printf("Failed to connect to RabbitMQ: %s", err)
		}
		defer conn.Close()

		ch, err := conn.Channel()
		if err != nil {
			log.Printf("Failed to open a channel: %s", err)
		}
		defer ch.Close()

		q, err := ch.QueueDeclare(
			"hello", // name
			false,   // durable
			false,   // delete when unused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)
		if err != nil {
			log.Printf("Failed to declare a queue: %s",  err)
		}

		body :=  fmt.Sprint( time.Now().Format(TIMEFORMAT), ": ", str )
		err = ch.Publish(
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})

		if err != nil {
			log.Printf("Failed to publish a message: %s",  err)
		}

		log.Printf(" [x] Sent %s", body)
}