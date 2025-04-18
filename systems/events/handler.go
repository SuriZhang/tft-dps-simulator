package eventsys

// EventHandler is implemented by any system that wants to receive events.
type EventHandler interface {
    HandleEvent(evt interface{})
}