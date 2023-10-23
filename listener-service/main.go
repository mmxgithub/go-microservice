package main

import (
	"fmt"
	"listener/event"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// connect to rabbitmq
	rabbitConn, err := newConnection()
	if err != nil {
		log.Panic(err)
	}
	defer rabbitConn.Close()

	// create consumer
	consumer, err := event.NewConsumer(rabbitConn)
	if err != nil {
		panic(err)
	}

	err = consumer.Listen([]string{"log.INFO", "log.WARNING", "log.ERROR"})
	if err != nil {
		log.Println(err)
	}

	// watch the queue and consume event
}

func newConnection() (*amqp.Connection, error) {
	var counts int64
	var connection *amqp.Connection
	cooldown := 1 * time.Second

	for {
		c, err := amqp.Dial("amqp://guest:guest@rabbitmq")
		if err != nil {
			fmt.Println("RabbitMQ not yet ready...")
			counts++
		} else {
			connection = c
			log.Println("RabbitMQ connected")
			break
		}

		if counts > 20 {
			fmt.Println(err)
			return nil, err
		}

		time.Sleep(cooldown)
		continue
	}
	return connection, nil
}
