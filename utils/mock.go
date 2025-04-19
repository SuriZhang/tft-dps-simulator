package utils

import (
	eventsys "github.com/suriz/tft-dps-simulator/systems/events"
)

// --- Mock Event Bus ---
// Simple mock to capture enqueued events for testing
type MockEventBus struct {
	enqueuedEvents []interface{}
	handlers       []eventsys.EventHandler
}

func NewMockEventBus() *MockEventBus {
	return &MockEventBus{enqueuedEvents: make([]interface{}, 0)}
}

func (m *MockEventBus) Enqueue(evt interface{}) {
	m.enqueuedEvents = append(m.enqueuedEvents, evt)
}

// RegisterHandler is a no-op for this mock in AutoAttackSystem tests
func (m *MockEventBus) RegisterHandler(h eventsys.EventHandler) {
	m.handlers = append(m.handlers, h)
}

// ProcessAll is a no-op for this mock in AutoAttackSystem tests
func (m *MockEventBus) ProcessAll() {
	for len(m.enqueuedEvents) > 0 {
		evt := m.enqueuedEvents[0]
		m.enqueuedEvents = m.enqueuedEvents[1:]
		for _, h := range m.handlers {
			h.HandleEvent(evt)
		}
	}
}

// GetLastEvent returns the last event enqueued
func (m *MockEventBus) GetLastEvent() interface{} {
	if len(m.enqueuedEvents) == 0 {
		return nil
	}
	return m.enqueuedEvents[len(m.enqueuedEvents)-1]
}

// GetAllEvents returns all enqueued events
func (m *MockEventBus) GetAllEvents() []interface{} {
	if (len(m.enqueuedEvents)) == 0 {
		return nil
	}
	return m.enqueuedEvents
}

func (m *MockEventBus) ClearEvents() {
	m.enqueuedEvents = make([]interface{}, 0)
}

// Helper function to simulate events and process them
func (m *MockEventBus) SimulateAndProcessEvent(event interface{}) {
	m.Enqueue(event)
	m.ProcessAll() // Directly call handler for testing
}
