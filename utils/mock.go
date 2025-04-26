// filepath: utils/mock.go
package utils

import (
	"container/heap"
	"math/rand"
	"reflect"

	eventsys "github.com/suriz/tft-dps-simulator/systems/events"
)

// MockEventBus simulates the event bus with timestamp ordering for testing.
type MockEventBus struct {
    // Use EventItems similar to the real queue
    queue    eventsys.EventQueue
    handlers []eventsys.EventHandler
    rng      *rand.Rand
    // Store processed events for inspection if needed
    processedEvents []*eventsys.EventItem
}

// NewMockEventBus creates a new MockEventBus.
func NewMockEventBus() *MockEventBus {
    // Seed random number generator for jitter consistency in tests if needed
    source := rand.NewSource(1) // Use a fixed seed for deterministic tests
    eq := make(eventsys.EventQueue, 0)
    heap.Init(&eq) // Initialize the heap
    return &MockEventBus{
        queue:           eq,
        handlers:        make([]eventsys.EventHandler, 0),
        rng:             rand.New(source),
        processedEvents: make([]*eventsys.EventItem, 0),
    }
}

// Enqueue adds an event to the mock queue with timestamp and jitter.
// Matches the EventEnqueuer interface.
func (m *MockEventBus) Enqueue(evt interface{}, timestamp float64) {
    // Add small random jitter based on design doc U(-10^-5, +10^-5)
    jitter := (m.rng.Float64()*2 - 1) * 1e-5
    enqueueTimestamp := timestamp + jitter
    item := &eventsys.EventItem{
        Event:            evt,
        Timestamp:        timestamp,
        EnqueueTimestamp: enqueueTimestamp,
        // index is managed by heap
    }
    heap.Push(&m.queue, item)
}

// RegisterHandler adds a handler.
func (m *MockEventBus) RegisterHandler(h eventsys.EventHandler) {
    m.handlers = append(m.handlers, h)
}

// Dequeue removes and returns the next event based on EnqueueTimestamp.
// Returns nil if the queue is empty.
func (m *MockEventBus) Dequeue() *eventsys.EventItem {
    if m.queue.Len() == 0 {
        return nil
    }
    return heap.Pop(&m.queue).(*eventsys.EventItem)
}

// ProcessNext processes the single next event in the queue.
// Returns the processed event item or nil if queue was empty.
// Simulates the core loop step: Dequeue -> Dispatch
func (m *MockEventBus) ProcessNext() *eventsys.EventItem {
    item := m.Dequeue()
    if item == nil {
        return nil
    }
    m.processedEvents = append(m.processedEvents, item) // Record processed event
    m.Dispatch(item.Event)                              // Dispatch the actual event data
    return item
}

func (m *MockEventBus) Dispatch(evt interface{}) {
	if evt == nil {
		return
	}
	for _, h := range m.handlers {
		if h.CanHandle(evt) {
			h.HandleEvent(evt)
		}
	}
}

// ProcessUntilTime processes events until the next event's timestamp exceeds the given time.
// Returns the timestamp of the last processed event, or the starting time if no events were processed.
func (m *MockEventBus) ProcessUntilTime(targetTime float64) float64 {
    lastProcessedTime := -1.0 // Indicate no events processed yet
    for m.queue.Len() > 0 {
        // Peek at the next event's time without removing it
        nextEventTime := m.queue[0].Timestamp
        if nextEventTime > targetTime {
            break // Stop if the next event is past the target time
        }
        processedItem := m.ProcessNext()
        if processedItem != nil {
             lastProcessedTime = processedItem.Timestamp
        }
    }
    return lastProcessedTime
}


// ProcessUntilEmpty processes all events currently in the queue in order.
func (m *MockEventBus) ProcessUntilEmpty() {
    for m.queue.Len() > 0 {
        m.ProcessNext()
    }
}

// GetLastEvent is less meaningful with timestamp ordering.
// Use GetProcessedEvents() or GetEventsOfType() instead.
func (m *MockEventBus) GetLastEvent() interface{} {
    if len(m.processedEvents) == 0 {
        return nil
    }
    // Returns the last *processed* event's data
    return m.processedEvents[len(m.processedEvents)-1].Event
}

// GetAllEvents returns a copy of the raw events *currently* in the queue (unsorted).
func (m *MockEventBus) GetAllEvents() []interface{} {
    events := make([]interface{}, 0, m.queue.Len())
    // Create a temporary copy to iterate over, as heap order isn't guaranteed stable
    tempQueue := make(eventsys.EventQueue, m.queue.Len())
    copy(tempQueue, m.queue)
    for _, item := range tempQueue {
        events = append(events, item.Event)
    }
    return events
}

// GetProcessedEvents returns a slice of the EventItems processed so far, in order.
func (m *MockEventBus) GetProcessedEvents() []*eventsys.EventItem {
    // Return a copy to prevent external modification
    result := make([]*eventsys.EventItem, len(m.processedEvents))
    copy(result, m.processedEvents)
    return result
}


// GetQueueItems returns a copy of the EventItems currently in the queue (unsorted).
func (m *MockEventBus) GetQueueItems() []*eventsys.EventItem {
    items := make([]*eventsys.EventItem, m.queue.Len())
    // Create a temporary copy to iterate over
    tempQueue := make(eventsys.EventQueue, m.queue.Len())
    copy(tempQueue, m.queue)
    copy(items, tempQueue)
    return items
}


// ClearEvents removes all events from the queue and clears processed history.
func (m *MockEventBus) ClearEvents() {
    m.queue = make(eventsys.EventQueue, 0)
    heap.Init(&m.queue) // Re-initialize heap
    m.processedEvents = make([]*eventsys.EventItem, 0)
}

// Len returns the number of events currently in the queue.
func (m *MockEventBus) Len() int {
    return m.queue.Len()
}

// FindFirstEventOfType finds the first processed event matching the type.
func (m *MockEventBus) FindFirstProcessedEventOfType(eventType reflect.Type) (interface{}, bool) {
    for _, item := range m.processedEvents {
        if reflect.TypeOf(item.Event) == eventType {
            return item.Event, true
        }
    }
    return nil, false
}

// FindAllEventsOfType finds all processed events matching the type.
func (m *MockEventBus) FindAllProcessedEventsOfType(eventType reflect.Type) []interface{} {
    var found []interface{}
    for _, item := range m.processedEvents {
        if reflect.TypeOf(item.Event) == eventType {
            found = append(found, item.Event)
        }
    }
    return found
}
