package amqp

import (
	"sync"

	"github.com/streadway/amqp"
	"go.k6.io/k6/js/modules"
)

var (
	_ modules.Module = &RootModule{}
)

func New() *RootModule {
	var (
		mutex = new(sync.RWMutex)
		cond  = sync.NewCond(mutex.RLocker())
	)

	return &RootModule{mutex, cond, nil}
}

type RootModule struct {
	*sync.RWMutex
	*sync.Cond
	connection *amqp.Connection
}

func (root *RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &Instance{vu: vu, root: root}
}

func (m *RootModule) connect(urls []string) (err error) {
	if m.check() {
		return nil
	}

	defer m.Broadcast()

	m.Lock()
	defer m.Unlock()

	m.connection, err = amqp.Dial(urls[0])

	return
}

func (m *RootModule) check() bool {
	m.RLock()
	defer m.RUnlock()

	return m.connection != nil
}

func (m *RootModule) channel() (*amqp.Channel, error) {
	m.RLock()
	defer m.RUnlock()

	for m.connection == nil {
		m.Wait()
	}

	return m.connection.Channel()
}
