package amqp

import (
	"github.com/dop251/goja"
	"go.k6.io/k6/js/common"
)

type API struct {
	*Instance
}

func (api API) connect(urls []string) error {
	return api.root.connect(urls)
}

func (api *API) queue(call goja.ConstructorCall, rt *goja.Runtime) *goja.Object {
	queue := new(Queue)
	queue.Instance = api.Instance

	rt.ExportTo(call.Argument(0), queue)

	if err := api.setup(); err != nil {
		common.Throw(rt, err)
	}

	if _, err := api.channel.QueueDeclare(
		queue.Name,
		queue.Durable,
		queue.DeleteWhenUnused,
		queue.Exclusive,
		false, // noWait
		queue.Args,
	); err != nil {
		common.Throw(rt, err)
	}

	obj := rt.ToValue(queue).(*goja.Object)
	obj.SetPrototype(call.This.Prototype())

	return obj
}

func (api *API) exchange(call goja.ConstructorCall, rt *goja.Runtime) *goja.Object {
	exchange := new(Exchange)
	exchange.Instance = api.Instance

	rt.ExportTo(call.Argument(0), exchange)

	if err := api.setup(); err != nil {
		common.Throw(rt, err)
	}

	if err := api.channel.ExchangeDeclare(
		exchange.Name,
		exchange.Kind,
		exchange.Durable,
		exchange.AutoDelete,
		exchange.Internal,
		false,
		exchange.Args,
	); err != nil {
		common.Throw(rt, err)
	}

	obj := rt.ToValue(exchange).(*goja.Object)
	obj.SetPrototype(call.This.Prototype())

	return obj
}
