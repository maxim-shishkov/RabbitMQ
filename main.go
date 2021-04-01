package main

import (
	"github.com/streadway/amqp"
	"log"
	"time"
)

const (
	URL = "amqp://guest:guest@localhost:5672/"
	TIMEFORMAT = "15:04:05"
)

func main() {
	rst := make(chan struct{})

	go workerRabbitMQ(rst)

	for{
		select {
		case <-rst:
			log.Println("Restarting Rabbit MQ")
			time.Sleep(1 * time.Second)
			go workerRabbitMQ(rst)
		}
	}
}

func workerRabbitMQ(rst chan struct{}) {
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

	errorAmqp := ch.NotifyClose(make(chan *amqp.Error, 1))

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
			case e := <- errorAmqp:
				log.Printf("chan error: %s", e.Error())
				rst <- struct{}{}
			case m := <-msgs:
				log.Printf("Received a message: %s", m.Body)
		}
	}

}


