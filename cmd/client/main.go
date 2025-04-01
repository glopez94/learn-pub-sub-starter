package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// Obtener el nombre de usuario
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Failed to get username: %v", err)
	}

	connStr := "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connStr)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	fmt.Println("Successfully connected to RabbitMQ")

	// Declarar y vincular la cola
	ch, q, err := pubsub.DeclareAndBind(conn, routing.ExchangePerilDirect, "pause."+username, routing.PauseKey, pubsub.Transient)
	if err != nil {
		log.Fatalf("Failed to declare and bind queue: %v", err)
	}
	defer ch.Close()

	fmt.Printf("Queue %s declared and bound to exchange %s\n", q.Name, routing.ExchangePerilDirect)

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Println("Shutting down...")
}
