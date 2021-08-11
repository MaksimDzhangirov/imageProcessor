package model

import amqp "github.com/rabbitmq/amqp091-go"

type QueueConfig struct {
	Name             string
	Durable          bool
	DeleteWhenUnused bool
	Exclusive        bool
	NoWait           bool
	Arguments        amqp.Table
}

type PublishConfig struct {
	Exchange        string
	Mandatory       bool
	Immediate       bool
	StringToPublish string
}

type ConsumeConfig struct {
	Consumer  string
	AutoAck   bool
	Exclusive bool
	NoLocal   bool
	NoWait    bool
	Arguments amqp.Table
	FuncToPerform func(uid string) error
}

type QosConfig struct {
	PrefetchCount int
	PrefetchSize  int
	Global        bool
}
