package systems_test

import (
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/suriz/tft-dps-simulator/data"
	eventsys "github.com/suriz/tft-dps-simulator/systems/events"

	"github.com/onsi/ginkgo/v2"
	"github.com/onsi/gomega"
)

func TestSystems(t *testing.T) {
	gomega.RegisterFailHandler(ginkgo.Fail)
	ginkgo.RunSpecs(t, "Systems Suite")
}

var _ = ginkgo.BeforeSuite(func() {
	// Load item data once for the entire manager suite
	// Adjust the path to your actual item data file
	dataDir := "../assets"
	fileName := "en_us_14.1b.json"
	filePath := filepath.Join(dataDir, fileName)
	tftData, err := data.LoadSetDataFromFile(filePath, "TFTSet14")
	if err != nil {
		log.Printf("Error loading set data: %v\n", err)
		os.Exit(1)
	}
	data.InitializeChampions(tftData)

	data.InitializeTraits(tftData)
	data.InitializeSetActiveAugments(tftData, filePath)

	data.InitializeSetActiveItems(tftData, filePath)

	gomega.Expect(err).NotTo(gomega.HaveOccurred(), "Failed to load item data")
})

// --- Mock Event Bus ---
// Simple mock to capture enqueued events for testing
type MockEventBus struct {
	EnqueuedEvents []interface{}
}

func NewMockEventBus() *MockEventBus {
	return &MockEventBus{EnqueuedEvents: make([]interface{}, 0)}
}

func (m *MockEventBus) Enqueue(evt interface{}) {
	m.EnqueuedEvents = append(m.EnqueuedEvents, evt)
}

// RegisterHandler is a no-op for this mock in AutoAttackSystem tests
func (m *MockEventBus) RegisterHandler(h eventsys.EventHandler) {}

// ProcessAll is a no-op for this mock in AutoAttackSystem tests
func (m *MockEventBus) ProcessAll() {
	// No-op for this mock in AutoAttackSystem tests
}

// GetLastEvent returns the last event enqueued
func (m *MockEventBus) GetLastEvent() interface{} {
	if len(m.EnqueuedEvents) == 0 {
		return nil
	}
	return m.EnqueuedEvents[len(m.EnqueuedEvents)-1]
}

// GetAllEvents returns all enqueued events
func (m *MockEventBus) GetAllEvents() []interface{} {
	if (len(m.EnqueuedEvents)) == 0 {
		return nil
	}
	return m.EnqueuedEvents
}

func (m *MockEventBus) ClearEvents() {
	m.EnqueuedEvents = make([]interface{}, 0)
}
