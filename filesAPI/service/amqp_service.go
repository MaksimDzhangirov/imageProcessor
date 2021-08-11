// https://www.rabbitmq.com/tutorials/tutorial-two-go.html
// https://x-team.com/blog/set-up-rabbitmq-with-docker-compose/
package service

import (
	"github.com/MaksimDzhangirov/imageProcessor/filesAPI/model"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

type amqpService struct {
	RabbitMQConnection *amqp.Connection
	ch                 *amqp.Channel
	queue              amqp.Queue
}

func NewAmqpService(conn *amqp.Connection) model.AmqpService {
	return &amqpService{
		RabbitMQConnection: conn,
		ch:                 nil,
		queue:              amqp.Queue{Name: ""},
	}
}

func (s *amqpService) openChannel() error {
	ch, err := s.RabbitMQConnection.Channel()
	if err != nil {
		return err
	}
	s.ch = ch

	return nil
}

func (s *amqpService) queueDeclare(queueConfig model.QueueConfig) error {
	q, err := s.ch.QueueDeclare(
		queueConfig.Name,
		queueConfig.Durable,
		queueConfig.DeleteWhenUnused,
		queueConfig.Exclusive,
		queueConfig.NoWait,
		queueConfig.Arguments,
	)
	if err != nil {
		return err
	}
	s.queue = q

	return nil
}

func (s *amqpService) publish(publishConfig model.PublishConfig) error {
	err := s.ch.Publish(
		publishConfig.Exchange,
		s.queue.Name,
		publishConfig.Mandatory,
		publishConfig.Immediate,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         []byte(publishConfig.StringToPublish),
		})
	if err != nil {
		return err
	}

	return nil
}

func (s *amqpService) Send(queueConfig model.QueueConfig, publishConfig model.PublishConfig) error {
	err := s.openChannel()
	if err != nil {
		log.Printf("Failed to open a channel: %v", err)
		return err
	}
	err = s.queueDeclare(queueConfig)
	if err != nil {
		log.Printf("Failed to declare a queue: %v", err)
		return err
	}
	err = s.publish(publishConfig)
	if err != nil {
		log.Printf("Failed to publish a message: %v", err)
		return err
	}

	return s.ch.Close()
}

func (s *amqpService) Consume(queueConfig model.QueueConfig, consumeConfig model.ConsumeConfig) error {
	err := s.openChannel()
	if err != nil {
		log.Printf("Failed to open a channel: %v", err)
		return err
	}
	err = s.queueDeclare(queueConfig)
	if err != nil {
		log.Printf("Failed to declare a queue: %v", err)
		return err
	}

	qosConfig := model.QosConfig{
		PrefetchCount: 1,
		PrefetchSize: 0,
		Global: false,
	}

	err = s.ch.Qos(
		qosConfig.PrefetchCount,
		qosConfig.PrefetchSize,
		qosConfig.Global,
	)
	if err != nil {
		log.Printf("Failed to set QoS: %v", err)
		return err
	}

	msgs, err := s.ch.Consume(
		s.queue.Name,
		consumeConfig.Consumer,
		consumeConfig.AutoAck,
		consumeConfig.Exclusive,
		consumeConfig.NoLocal,
		consumeConfig.NoWait,
		consumeConfig.Arguments,
	)
	if err != nil {
		log.Printf("Failed to register a consumer: %v", err)
		return err
	}

	forever := make(chan bool)
	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			err := consumeConfig.FuncToPerform(string(d.Body))
			if err != nil {
				log.Printf("Cannot perform oparation, %v", err)
			}
			log.Printf("Done")
			d.Ack(false)
		}
	}()
	<-forever

	return s.ch.Close()
}