package amqp

import (
	"github.com/streadway/amqp"
)

type Exchange struct {
	*Instance

	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	Args       amqp.Table
}

func (exchange *Exchange) Delete() error {
	if err := exchange.setup(); err != nil {
		return err
	}

	return exchange.channel.ExchangeDelete(
		exchange.Name,
		false, // ifUnused
		false, // noWait
	)
}

func (source *Exchange) Bind(destination Exchange, key string) error {
	if err := source.setup(); err != nil {
		return err
	}

	return source.channel.ExchangeBind(
		destination.Name,
		key,
		source.Name,
		false, // noWait
		amqp.Table{},
	)
}

func (source *Exchange) Unbind(destination Exchange, key string) error {
	if err := source.setup(); err != nil {
		return err
	}

	return source.channel.ExchangeUnbind(
		destination.Name,
		key,
		source.Name,
		false, // noWait
		amqp.Table{},
	)
}

type Message struct {
	Mandatory   bool
	Immediate   bool
	ContentType string
	Body        string
}

func (exchange *Exchange) Pulbish(key string, message Message) error {
	if err := exchange.setup(); err != nil {
		return err
	}

	return exchange.channel.Publish(
		exchange.Name,
		key,
		message.Mandatory,
		message.Immediate,
		amqp.Publishing{
			ContentType: message.ContentType,
			Body:        []byte(message.Body),
		},
	)
}
