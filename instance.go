package amqp

import (
	"github.com/streadway/amqp"
	"go.k6.io/k6/js/modules"
)

var (
	_ modules.Instance = &Instance{}
)

type Instance struct {
	root    *RootModule
	vu      modules.VU
	channel *amqp.Channel
}

func (i *Instance) Exports() modules.Exports {
	api := &API{i}

	return modules.Exports{
		Default: api,
		Named: map[string]interface{}{
			"Connect":  api.connect,
			"Queue":    api.queue,
			"Exchange": api.exchange,
		},
	}
}

func (i *Instance) setup() (err error) {
	if i.channel == nil {
		i.channel, err = i.root.channel()
	}

	return
}

func (i *Instance) shutdown() error {
	if i.channel == nil {
		return i.channel.Close()
	}

	return nil
}
