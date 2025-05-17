package eventsys

import (
    "container/heap"
    "math/rand"
    "time"
)

// EventItem holds an event and its priority (EnqueueTimestamp) for the queue.
type EventItem struct {
    Event            interface{} // The actual event data
    Timestamp        float64     // The logical time the event occurs
    EnqueueTimestamp float64     // Timestamp + jitter, used for priority sorting
    index            int         // Index of the item in the heap
}

// EventQueue implements heap.Interface and holds EventItems.
type EventQueue []*EventItem

func (eq EventQueue) Len() int { return len(eq) }

func (eq EventQueue) Less(i, j int) bool {
    // Min-heap based on EnqueueTimestamp
    return eq[i].EnqueueTimestamp < eq[j].EnqueueTimestamp
}

func (eq EventQueue) Swap(i, j int) {
    eq[i], eq[j] = eq[j], eq[i]
    eq[i].index = i
    eq[j].index = j
}

// Push adds an item to the queue.
func (eq *EventQueue) Push(x interface{}) {
    n := len(*eq)
    item := x.(*EventItem)
    item.index = n
    *eq = append(*eq, item)
}

// Pop removes and returns the item with the highest priority (lowest EnqueueTimestamp).
func (eq *EventQueue) Pop() interface{} {
    old := *eq
    n := len(old)
    item := old[n-1]
    old[n-1] = nil  // avoid memory leak
    item.index = -1 // for safety
    *eq = old[0 : n-1]
    return item
}

// EventEnqueuer provides an interface for adding events to the queue.
type EventEnqueuer interface {

    Enqueue(evt interface{}, timestamp float64)
}

// PriorityQueue wraps the EventQueue and provides Enqueue method.
type PriorityQueue struct {
    queue *EventQueue
    rng   *rand.Rand
}

// NewPriorityQueue creates a new event priority queue.
func NewPriorityQueue() *PriorityQueue {
    eq := make(EventQueue, 0)
    heap.Init(&eq)
    // Seed random number generator for jitter
    source := rand.NewSource(time.Now().UnixNano())
    return &PriorityQueue{
        queue: &eq,
        rng:   rand.New(source),
    }
}

// Enqueue adds an event to the priority queue with jitter.
func (pq *PriorityQueue) Enqueue(evt interface{}, timestamp float64) {
    // Add small random jitter based on design doc U(-10^-5, +10^-5)
    jitter := (pq.rng.Float64()*2 - 1) * 1e-5
    enqueueTimestamp := timestamp + jitter
    item := &EventItem{
        Event:            evt,
        Timestamp:        timestamp,
        EnqueueTimestamp: enqueueTimestamp,
    }
    heap.Push(pq.queue, item)
}

// Dequeue removes and returns the next event item from the queue.
// Returns nil if the queue is empty.
func (pq *PriorityQueue) Dequeue() *EventItem {
    if pq.queue.Len() == 0 {
        return nil
    }
    return heap.Pop(pq.queue).(*EventItem)
}

// Len returns the number of items in the queue.
func (pq *PriorityQueue) Len() int {
    return pq.queue.Len()
}