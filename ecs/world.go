package ecs

import (
	"fmt"
	"reflect"
)

// World contains all entities and their components
type World struct {
    // Maps component type to a map of entity->component
    components map[reflect.Type]map[Entity]interface{}
    // All active entities
    entities map[Entity]bool
}

// NewWorld creates a new empty world
func NewWorld() *World {
    return &World{
        components: make(map[reflect.Type]map[Entity]interface{}),
        entities:   make(map[Entity]bool),
    }
}

// CreateEntity creates a new entity in the world
func (w *World) CreateEntity() Entity {
    e := NewEntity()
    w.entities[e] = true
    return e
}

// RemoveEntity removes an entity and all its components
func (w *World) RemoveEntity(e Entity) {
    delete(w.entities, e)
    // Remove all components for this entity
    for _, componentMap := range w.components {
        delete(componentMap, e)
    }
}

// AddComponent adds a component to an entity
func (w *World) AddComponent(e Entity, component interface{}) {
    componentType := reflect.TypeOf(component)
    
    // Create map for this component type if it doesn't exist
    if _, exists := w.components[componentType]; !exists {
        w.components[componentType] = make(map[Entity]interface{})
    }
    
    // Add component to entity
    w.components[componentType][e] = component
}

// GetComponent returns a component for an entity
func (w *World) GetComponent(e Entity, componentType reflect.Type) (interface{}, bool) {
    if componentMap, exists := w.components[componentType]; exists {
        if component, hasComponent := componentMap[e]; hasComponent {
            return component, true
        }
    }
    return nil, false
}

// HasComponent checks if an entity has a specific component
func (w *World) HasComponent(e Entity, componentType reflect.Type) bool {
    if componentMap, exists := w.components[componentType]; exists {
        _, hasComponent := componentMap[e]
        return hasComponent
    }
    return false
}

// GetEntitiesWithComponents returns all entities that have all the specified component types
func (w *World) GetEntitiesWithComponents(componentTypes ...reflect.Type) []Entity {
    if len(componentTypes) == 0 {
        return []Entity{}
    }
    
    // Start with entities from first component type
    candidateEntities := make(map[Entity]bool)
    if componentMap, exists := w.components[componentTypes[0]]; exists {
        for e := range componentMap {
            candidateEntities[e] = true
        }
    }
    
    // Filter by remaining component types
    for _, ct := range componentTypes[1:] {
        componentMap, exists := w.components[ct]
        if !exists {
            return []Entity{} // No entities have this component
        }
        
        // Keep only entities that have this component
        for e := range candidateEntities {
            if _, hasComponent := componentMap[e]; !hasComponent {
                delete(candidateEntities, e)
            }
        }
    }
    
    // Convert map keys to slice
    result := make([]Entity, 0, len(candidateEntities))
    for e := range candidateEntities {
        result = append(result, e)
    }
    return result
}

// String returns a string representation of the world
func (w *World) String() string {
    return fmt.Sprintf("World with %d entities", len(w.entities))
}

