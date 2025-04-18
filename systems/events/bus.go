package eventsys

// EventBus manages registration and dispatch of events.
type EventBus interface {
	Enqueue(evt interface{})
	RegisterHandler(h EventHandler)
	ProcessAll()
}

// SimpleBus is a basic FIFO bus with pub/sub.
type SimpleBus struct {
	handlers []EventHandler
	queue    []interface{}
}

func NewSimpleBus() *SimpleBus {
	return &SimpleBus{
		handlers: make([]EventHandler, 0),
		queue:    make([]interface{}, 0),
	}
}

func (b *SimpleBus) RegisterHandler(h EventHandler) {
	b.handlers = append(b.handlers, h)
}

func (b *SimpleBus) Enqueue(evt interface{}) {
	b.queue = append(b.queue, evt)
}

func (b *SimpleBus) ProcessAll() {
	for len(b.queue) > 0 {
		evt := b.queue[0]
		b.queue = b.queue[1:]
		for _, h := range b.handlers {
			h.HandleEvent(evt)
		}
	}
}
