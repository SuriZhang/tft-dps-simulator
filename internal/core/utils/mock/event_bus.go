package mock

import (
	eventsys "tft-dps-simulator/internal/core/systems/events"
)

type EventBus struct {
	EnqueuedEvents []eventsys.Event // Store enqueued events for inspection
}

func (m *EventBus) Enqueue(event eventsys.Event, timestamp float64) {
	m.EnqueuedEvents = append(m.EnqueuedEvents, event)
}

func (m *EventBus) Dequeue() (eventsys.Event, bool) {
	if len(m.EnqueuedEvents) == 0 {
		return nil, false
	}
	event := m.EnqueuedEvents[0]
	m.EnqueuedEvents = m.EnqueuedEvents[1:]
	return event, true
}

func (m *EventBus) IsEmpty() bool {
	return len(m.EnqueuedEvents) == 0
}

func (m *EventBus) Peek() (eventsys.Event, bool) {
	if len(m.EnqueuedEvents) == 0 {
		return nil, false
	}
	return m.EnqueuedEvents[0], true
}

func (m *EventBus) Subscribe(handler eventsys.EventHandler) {
	// No-op for mock, or add basic subscription tracking if needed
}

func (m *EventBus) Unsubscribe(handler eventsys.EventHandler) {
	// No-op for mock
}

func (m *EventBus) ProcessEvents(currentTime float64) {
	// No-op for mock, or simulate event processing if necessary for specific tests
}
func (m *EventBus) Clear() {
	m.EnqueuedEvents = []eventsys.Event{}
}
