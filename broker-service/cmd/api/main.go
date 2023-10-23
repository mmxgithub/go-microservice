package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

const webPort = "8080"

type Config struct {
	RabbitConn *amqp.Connection
}

func main() {
	// connect to rabbitmq
	rabbitConn, err := newConnection()
	if err != nil {
		log.Panic(err)
	}
	defer rabbitConn.Close()

	app := Config{
		RabbitConn: rabbitConn,
	}

	log.Printf("Starting broker service on port %s\n", webPort)

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", webPort),
		Handler: app.routes(),
	}

	err = srv.ListenAndServe()
	if err != nil {
		log.Panic(err)
	}
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
