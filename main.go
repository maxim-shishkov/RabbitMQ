package main

import (
	"github.com/streadway/amqp"
	"log"
	"time"
)
const URL = "amqp://guest:guest@localhost:5672/"


	func main() {
		RabbitMQ()
	//	send("test Msg")
	}

	func failOnError(err error, msg string) {
		if err != nil {
			log.Printf("%s: %s", msg, err)
		}
	}

	func RabbitMQ() {
		defer func() {
			if err := recover(); err != nil {
				time.Sleep(1 * time.Second)
				RabbitMQ()
			}
		}()

		conn, err := amqp.Dial(URL)
		failOnError(err, "Failed to connect to RabbitMQ")
		defer conn.Close()

		ch, err := conn.Channel()
		failOnError(err, "Failed to open a channel")
		defer ch.Close()

		closeChan := make(chan *amqp.Error, 1)
		notify := ch.NotifyClose(closeChan)

		q, err := ch.QueueDeclare(
			"hello", // name
			false,   // durable
			false,   // delete when unused
			false,   // exclusive
			false,   // no-wait
			nil,     // arguments
		)
		msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer
			true,   // auto-ack
			false,  // exclusive
			false,  // no-local
			false,  // no-wait
			nil,    // args
		)

		for {
			select {
				case e := <-notify:
					log.Printf("chan error: %s", e.Error())
					close(closeChan)
					RabbitMQ()
				case m := <-msgs:
					log.Printf("Received a message: %s", m.Body)
			}
		}
	}


