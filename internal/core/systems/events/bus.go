package eventsys

// EventBus defines the interface for an event dispatch system.
type EventBus interface {
    RegisterHandler(handler EventHandler)
    Dispatch(evt interface{}) // Dispatch a single event to appropriate handlers
    EventEnqueuer // Embed the enqueuer interface
}

// SimpleBus implements a basic synchronous event bus.
type SimpleBus struct {
    handlers []EventHandler
    queue    *PriorityQueue // Use the priority queue
}

// NewSimpleBus creates a new SimpleBus.
func NewSimpleBus() *SimpleBus {
    return &SimpleBus{
        handlers: make([]EventHandler, 0),
        queue:    NewPriorityQueue(), // Initialize the priority queue
    }
}

// RegisterHandler adds a new event handler.
func (b *SimpleBus) RegisterHandler(handler EventHandler) {
    b.handlers = append(b.handlers, handler)
}

// Enqueue adds an event to the priority queue.
func (b *SimpleBus) Enqueue(evt interface{}, timestamp float64) {
    b.queue.Enqueue(evt, timestamp)
}

// Dequeue removes and returns the next event item from the queue.
func (b *SimpleBus) Dequeue() *EventItem {
   return b.queue.Dequeue()
}

// Len returns the number of events currently in the queue.
func (b *SimpleBus) Len() int {
    return b.queue.Len()
}


// Dispatch sends a single event to all registered handlers that can handle it.
func (b *SimpleBus) Dispatch(evt interface{}) {
    if evt == nil {
        return
    }
    for _, handler := range b.handlers {
        if handler.CanHandle(evt) {
            handler.HandleEvent(evt)
        }
    }
}
