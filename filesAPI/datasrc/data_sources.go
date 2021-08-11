package datasrc

import (
	"fmt"
	"log"
	"os"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type DataSources struct {
	RabbitMQConnection *amqp.Connection
}

// InitDS establishes connections to fields in dataSources
func InitDS() (*DataSources, error) {
	log.Printf("Initializing data sources\n")

	amqpHost := os.Getenv("AMQP_HOST")
	amqpPort := os.Getenv("AMQP_PORT")
	ampqUser := os.Getenv("AMQP_USER")
	amqpPassword := os.Getenv("AMQP_PASSWORD")
	amqpNotReady := true
	counter := 1
	var conn *amqp.Connection
	var err error
	for amqpNotReady && counter < 3 {
		log.Printf("Dial to amqp://%s:%s@%s:%s/", ampqUser, amqpPassword, amqpHost, amqpPort)
		conn, err = amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", ampqUser, amqpPassword, amqpHost, amqpPort))
		if err != nil {
			log.Printf("%v", err)
			time.Sleep(3*time.Second)
			counter++
			continue
		}
		amqpNotReady = false
	}

	if err != nil {
		return nil, fmt.Errorf("error creating RabbitMQ connection: %w", err)
	}

	return &DataSources{
		RabbitMQConnection: conn,
	}, nil
}

// Close to be used in graceful server shutdown
func (d *DataSources) Close() error {
	if err := d.RabbitMQConnection.Close(); err != nil {
		return fmt.Errorf("error closing RabbitMQ: %w", err)
	}

	return nil
}
