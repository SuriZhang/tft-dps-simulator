package eventsys

// EventHandler defines the interface for types that can handle specific events.
type EventHandler interface {
    HandleEvent(evt interface{})
    // CanHandle returns true if the handler can process the given event type.
    CanHandle(evt interface{}) bool
}
