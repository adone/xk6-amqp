package amqp

import (
	"github.com/streadway/amqp"
)

type Queue struct {
	*Instance

	Name             string
	Durable          bool
	DeleteWhenUnused bool
	Exclusive        bool
	Args             amqp.Table
}

func (queue *Queue) Inspect() (amqp.Queue, error) {
	if err := queue.setup(); err != nil {
		return amqp.Queue{}, err
	}

	return queue.channel.QueueInspect(queue.Name)
}

func (queue *Queue) Delete() error {
	if err := queue.setup(); err != nil {
		return err
	}

	_, err := queue.channel.QueueDelete(
		queue.Name,
		false, // ifUnused
		false, // ifEmpty
		false, // noWait
	)

	return err
}

func (queue *Queue) Bind(exchange Exchange, key string) error {
	if err := queue.setup(); err != nil {
		return err
	}

	return queue.channel.QueueBind(
		queue.Name,
		key,
		exchange.Name,
		false, // noWait
		amqp.Table{},
	)
}

func (queue *Queue) Unbind(exchange Exchange, key string) error {
	if err := queue.setup(); err != nil {
		return err
	}

	return queue.channel.QueueUnbind(
		queue.Name,
		key,
		exchange.Name,
		amqp.Table{},
	)
}

func (queue *Queue) Purge() (int, error) {
	if err := queue.setup(); err != nil {
		return 0, err
	}

	return queue.channel.QueuePurge(queue.Name, false)
}

func (queue *Queue) Listen(options ListenOptions) error {
	if err := queue.setup(); err != nil {
		return err
	}

	msgs, err := queue.channel.Consume(
		queue.Name,
		options.Consumer,
		options.AutoAck,
		options.Exclusive,
		options.NoLocal,
		options.NoWait,
		options.Args,
	)
	if err != nil {
		return err
	}

	go func() {
		for message := range msgs {
			options.Listener(string(message.Body))
		}
	}()

	return nil
}
