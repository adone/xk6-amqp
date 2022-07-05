package amqp

import (
	amqpDriver "github.com/streadway/amqp"
)

type Exchanges struct {
	Version    string
	Connection *amqpDriver.Connection
}

type ExchangeOptions struct {
	ConnectionUrl string
}

type ExchangeDeclareOptions struct {
	Name       string
	Kind       string
	Durable    bool
	AutoDelete bool
	Internal   bool
	NoWait     bool
	Args       amqpDriver.Table
}

type ExchangeBindOptions struct {
	DestinationExchangeName string
	SourceExchangeName      string
	RoutingKey              string
	NoWait                  bool
	Args                    amqpDriver.Table
}

type ExchangeUnindOptions struct {
	DestinationExchangeName string
	SourceExchangeName      string
	RoutingKey              string
	NoWait                  bool
	Args                    amqpDriver.Table
}

func (exchange *Exchanges) Declare(options ExchangeDeclareOptions) error {
	ch, err := exchange.Connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return ch.ExchangeDeclare(
		options.Name,
		options.Kind,
		options.Durable,
		options.AutoDelete,
		options.Internal,
		options.NoWait,
		options.Args,
	)
}

func (exchange *Exchanges) Delete(name string) error {
	ch, err := exchange.Connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return ch.ExchangeDelete(
		name,
		false, // ifUnused
		false, // noWait
	)
}

func (exchange *Exchanges) Bind(options ExchangeBindOptions) error {
	ch, err := exchange.Connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return ch.ExchangeBind(
		options.DestinationExchangeName,
		options.RoutingKey,
		options.SourceExchangeName,
		options.NoWait,
		options.Args,
	)
}

func (exchange *Exchanges) Unbind(options ExchangeUnindOptions) error {
	ch, err := exchange.Connection.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()
	return ch.ExchangeUnbind(
		options.DestinationExchangeName,
		options.RoutingKey,
		options.SourceExchangeName,
		options.NoWait,
		options.Args,
	)
}
