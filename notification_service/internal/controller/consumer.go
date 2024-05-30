package controller

import "github.com/streadway/amqp"

type Consumer struct {
	Channel *amqp.Channel
}

func NewConsumer(channel *amqp.Channel) *Consumer {
	return &Consumer{
		Channel: channel,
	}
}
