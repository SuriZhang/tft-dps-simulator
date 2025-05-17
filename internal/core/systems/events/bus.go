package eventsys

import (
	"log"
)

// EventBus defines the interface for an event dispatch system.
type EventBus interface {
	RegisterHandler(handler EventHandler)
	Dispatch(evt interface{}) // Dispatch a single event to appropriate handlers
	EventEnqueuer             // Embed the enqueuer interface
	GetArchive() []*EventItem // Return the archived events
}

// SimpleBus implements a basic synchronous event bus.
type SimpleBus struct {
	handlers     []EventHandler
	queue        *PriorityQueue // Use the priority queue
	archiveQueue []*EventItem   // Archive of processed events
}

// NewSimpleBus creates a new SimpleBus.
func NewSimpleBus() *SimpleBus {
	return &SimpleBus{
		handlers:     make([]EventHandler, 0),
		queue:        NewPriorityQueue(),    // Initialize the priority queue
		archiveQueue: make([]*EventItem, 0), // Initialize the archive queue
	}
}

// RegisterHandler adds a new event handler.
func (b *SimpleBus) RegisterHandler(handler EventHandler) {
	b.handlers = append(b.handlers, handler)
}

// Enqueue adds an event to the priority queue.
func (b *SimpleBus) Enqueue(evt interface{}, timestamp float64) {
	log.Printf("DEBUG: Enqueueing event (%T): %+v at timestamp: %f", evt, evt, timestamp)
	b.queue.Enqueue(evt, timestamp)
}

// Dequeue removes and returns the next event item from the queue.
// The dequeued event is saved to the archiveQueue for later analysis.
func (b *SimpleBus) Dequeue() *EventItem {
	item := b.queue.Dequeue()
	if item != nil {
		// Archive the event for later analysis
		b.archiveQueue = append(b.archiveQueue, item)
	}
	return item
}

// Len returns the number of events currently in the queue.
func (b *SimpleBus) Len() int {
	return b.queue.Len()
}

// GetArchive returns the archive of processed events.
// This can be used to generate time sequence diagrams.
func (b *SimpleBus) GetArchive() []*EventItem {
	return b.archiveQueue
}

// Dispatch sends a single event to all registered handlers that can handle it.
func (b *SimpleBus) Dispatch(evt interface{}) {
	if evt == nil {
		return
	}
	for _, handler := range b.handlers {
		if handler.CanHandle(evt) {
			log.Printf("DEBUG: Dispatching event (%T): %+v to handler: %T", evt, evt, handler)
			handler.HandleEvent(evt)
		}
	}
}
