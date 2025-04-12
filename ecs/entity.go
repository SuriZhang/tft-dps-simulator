package ecs

import (
	"sync/atomic"
)

// Entity is just an ID to reference components
type Entity uint32

var nextEntityID uint32 = 0

// NewEntity creates a new entity with a unique ID
func NewEntity() Entity {
    return Entity(atomic.AddUint32(&nextEntityID, 1))
}

