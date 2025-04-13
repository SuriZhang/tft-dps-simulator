package ecs

import (
	"fmt"
	"reflect" // Stigithub.com/suriz/tft-dps-simulator/components"

	"github.com/suriz/tft-dps-simulator/components"
)

// World contains all entities and their components, stored in type-specific maps.
type World struct {
	// --- Component Maps ---
	// Using pointers (*components.Type) allows checking for nil to see if an entity has the component.
	Health       map[Entity]*components.Health
	Mana         map[Entity]*components.Mana
	Attack       map[Entity]*components.Attack
	Traits       map[Entity]*components.Traits
	ChampionInfo map[Entity]*components.ChampionInfo
	Position     map[Entity]*components.Position
	Team         map[Entity]*components.Team
	Item         map[Entity]*components.ItemEffect
	Equipment    map[Entity]*components.Equipment // Assuming you have an Inventory component
	// Add maps for other components defined in your components directory as needed:
	// Defense      map[Entity]*components.Defense
	// Spell        map[Entity]*components.Spell
	// Buffs        map[Entity]*components.Buffs
}

// NewWorld creates a new empty world, initializing all component maps.
func NewWorld() *World {
	return &World{
		// Initialize maps
		Health:       make(map[Entity]*components.Health),
		Mana:         make(map[Entity]*components.Mana),
		Attack:       make(map[Entity]*components.Attack),
		Traits:       make(map[Entity]*components.Traits),
		ChampionInfo: make(map[Entity]*components.ChampionInfo),
		Position:     make(map[Entity]*components.Position),
		Team:         make(map[Entity]*components.Team),
		Item:         make(map[Entity]*components.ItemEffect),
		Equipment:    make(map[Entity]*components.Equipment),
		// Initialize other maps here...
	}
}

// CreateEntity generates a new unique entity ID using the global NewEntity function.
// Note: This only reserves the ID; components must be added separately.
func (w *World) CreateEntity() Entity {
	return NewEntity()
}

// RemoveEntity removes an entity and all its associated components from the world.
func (w *World) RemoveEntity(e Entity) {
	// Delete the entity from every component map
	delete(w.Health, e)
	delete(w.Mana, e)
	delete(w.Attack, e)
	delete(w.Traits, e)
	delete(w.ChampionInfo, e)
	delete(w.Position, e)
	delete(w.Team, e)
	delete(w.Item, e)
	delete(w.Equipment, e)
	// Delete from other maps here...
}

// AddComponent adds a component to an entity. It uses a type switch
// to place the component in the correct map.
// Returns an error if the component type is not recognized by the World struct.
func (w *World) AddComponent(e Entity, component interface{}) error {
	// Ensure component is not nil
	if component == nil {
		return fmt.Errorf("cannot add nil component to entity %d", e)
	}

	switch c := component.(type) {
	// Handle both value and pointer types for convenience
	case components.Health:
		w.Health[e] = &c
	case *components.Health:
		w.Health[e] = c
	case components.Mana:
		w.Mana[e] = &c
	case *components.Mana:
		w.Mana[e] = c
	case components.Attack:
		w.Attack[e] = &c
	case *components.Attack:
		w.Attack[e] = c
	case components.Traits:
		w.Traits[e] = &c
	case *components.Traits:
		w.Traits[e] = c
	case components.ChampionInfo:
		w.ChampionInfo[e] = &c
	case *components.ChampionInfo:
		w.ChampionInfo[e] = c
	case components.Position:
		w.Position[e] = &c
	case *components.Position:
		w.Position[e] = c
	case components.Team:
		w.Team[e] = &c
	case *components.Team:
		w.Team[e] = c
	case components.ItemEffect:
		w.Item[e] = &c
	case *components.ItemEffect:
		w.Item[e] = c
	case components.Equipment:
		w.Equipment[e] = &c
	case *components.Equipment:
		w.Equipment[e] = c
	// Add cases for other component types here...
	default:
		// Use reflection to get the type name for the error message
		return fmt.Errorf("unknown component type: %v", reflect.TypeOf(component))
	}
	return nil
}

// GetComponent retrieves a component of a specific type for an entity.
// It returns the component (as interface{}) and true if found, otherwise nil and false.
// This generic version is kept for flexibility but type-safe getters are preferred.
func (w *World) GetComponent(e Entity, componentType reflect.Type) (interface{}, bool) {
	switch componentType {
	case reflect.TypeOf(components.Health{}):
		comp, ok := w.Health[e] // Use type-safe getter internally
		return comp, ok
	case reflect.TypeOf(components.Mana{}):
		comp, ok := w.Mana[e]
		return comp, ok
	case reflect.TypeOf(components.Attack{}):
		comp, ok := w.Attack[e]
		return comp, ok
	case reflect.TypeOf(components.Traits{}):
		comp, ok := w.Traits[e]
		return comp, ok
	case reflect.TypeOf(components.ChampionInfo{}):
		comp, ok := w.ChampionInfo[e]
		return comp, ok
	case reflect.TypeOf(components.Position{}):
		comp, ok := w.Position[e]
		return comp, ok
	case reflect.TypeOf(components.Team{}):
		comp, ok := w.Team[e]
		return comp, ok
	case reflect.TypeOf(components.ItemEffect{}):
		comp, ok := w.Item[e]
		return comp, ok
	case reflect.TypeOf(components.Equipment{}):
		comp, ok := w.Equipment[e]
		return comp, ok

	// Add cases for other component types here...
	default:
		return nil, false
	}
}

// HasComponent checks if an entity possesses a component of the specified type.
func (w *World) HasComponent(e Entity, componentType reflect.Type) bool {
	_, ok := w.GetComponent(e, componentType) // Leverage existing logic
	return ok
}

// RemoveComponent removes a specific component type from an entity.
func (w *World) RemoveComponent(e Entity, componentType reflect.Type) {
	switch componentType {
	case reflect.TypeOf(components.Health{}):
		delete(w.Health, e)
	case reflect.TypeOf(components.Mana{}):
		delete(w.Mana, e)
	case reflect.TypeOf(components.Attack{}):
		delete(w.Attack, e)
	case reflect.TypeOf(components.Traits{}):
		delete(w.Traits, e)
	case reflect.TypeOf(components.ChampionInfo{}):
		delete(w.ChampionInfo, e)
	case reflect.TypeOf(components.Position{}):
		delete(w.Position, e)
	case reflect.TypeOf(components.Team{}):
		delete(w.Team, e)
	case reflect.TypeOf(components.ItemEffect{}):
		delete(w.Item, e)
	case reflect.TypeOf(components.Equipment{}):
		delete(w.Equipment, e)
	// Add cases for other component types here...
	default:
		fmt.Printf("Warning: Attempted to remove unknown component type %v from entity %d\n", componentType, e)
	}
}

// GetEntitiesWithComponents returns a slice of entities that possess ALL the specified component types.
func (w *World) GetEntitiesWithComponents(componentTypes ...reflect.Type) []Entity {
	if len(componentTypes) == 0 {
		return []Entity{}
	}

	// Optimization: Find the component type with the fewest entities first.
	var smallestMapSize = -1
	var firstType reflect.Type
	for _, ct := range componentTypes {
		currentSize := w.getMapSizeForType(ct)
		if smallestMapSize == -1 || currentSize < smallestMapSize {
			smallestMapSize = currentSize
			firstType = ct
		}
	}

	if smallestMapSize == 0 {
		return []Entity{} // If the smallest map is empty, no entities can have all components.
	}

	// Initialize candidates with entities from the smallest map
	initialSet := w.getEntitiesForType(firstType)
	candidates := make(map[Entity]bool, len(initialSet))
	for _, e := range initialSet {
		candidates[e] = true
	}

	// Filter by the remaining component types
	for _, ct := range componentTypes {
		if ct == firstType { // Skip the type we already used
			continue
		}

		// Check entities for the current component type
		currentTypeSet := make(map[Entity]bool)
		for _, e := range w.getEntitiesForType(ct) {
			currentTypeSet[e] = true
		}

		// Filter candidates: keep only those present in currentTypeSet
		for e := range candidates {
			if !currentTypeSet[e] {
				delete(candidates, e) // Remove if entity doesn't have the current component type
			}
		}

		// Early exit if no candidates remain
		if len(candidates) == 0 {
			return []Entity{}
		}
	}

	// Convert the final candidate map keys to a slice
	result := make([]Entity, 0, len(candidates))
	for e := range candidates {
		result = append(result, e)
	}
	return result
}

// Helper to get map size for a type
func (w *World) getMapSizeForType(componentType reflect.Type) int {
	switch componentType {
	case reflect.TypeOf(components.Health{}):
		return len(w.Health)
	case reflect.TypeOf(components.Mana{}):
		return len(w.Mana)
	case reflect.TypeOf(components.Attack{}):
		return len(w.Attack)
	case reflect.TypeOf(components.Traits{}):
		return len(w.Traits)
	case reflect.TypeOf(components.ChampionInfo{}):
		return len(w.ChampionInfo)
	case reflect.TypeOf(components.Position{}):
		return len(w.Position)
	case reflect.TypeOf(components.Team{}):
		return len(w.Team)
	case reflect.TypeOf(components.ItemEffect{}):
		return len(w.Item)
	case reflect.TypeOf(components.Equipment{}):
		return len(w.Equipment)
	// Add cases for other component types...
	default:
		return 0
	}
}

// Helper function for GetEntitiesWithComponents to get entities for a single type
func (w *World) getEntitiesForType(componentType reflect.Type) []Entity {
	var entities []Entity
	switch componentType {
	case reflect.TypeOf(components.Health{}):
		entities = make([]Entity, 0, len(w.Health))
		for e := range w.Health {
			entities = append(entities, e)
		}
	case reflect.TypeOf(components.Mana{}):
		entities = make([]Entity, 0, len(w.Mana))
		for e := range w.Mana {
			entities = append(entities, e)
		}
	case reflect.TypeOf(components.Attack{}):
		entities = make([]Entity, 0, len(w.Attack))
		for e := range w.Attack {
			entities = append(entities, e)
		}
	case reflect.TypeOf(components.Traits{}):
		entities = make([]Entity, 0, len(w.Traits))
		for e := range w.Traits {
			entities = append(entities, e)
		}
	case reflect.TypeOf(components.ChampionInfo{}):
		entities = make([]Entity, 0, len(w.ChampionInfo))
		for e := range w.ChampionInfo {
			entities = append(entities, e)
		}
	case reflect.TypeOf(components.Position{}):
		entities = make([]Entity, 0, len(w.Position))
		for e := range w.Position {
			entities = append(entities, e)
		}
	case reflect.TypeOf(components.Team{}):
		entities = make([]Entity, 0, len(w.Team))
		for e := range w.Team {
			entities = append(entities, e)
		}
	case reflect.TypeOf(components.ItemEffect{}):
		entities = make([]Entity, 0, len(w.Item))
		for e := range w.Item {
			entities = append(entities, e)
		}
	case reflect.TypeOf(components.Equipment{}):
		entities = make([]Entity, 0, len(w.Equipment))
		for e := range w.Equipment {
			entities = append(entities, e)
		}
	// Add cases for other component types...
	default:
		return []Entity{} // Return empty slice for unknown types
	}
	return entities
}

// --- Type-Safe Getters (Recommended) ---

// GetHealth returns the Health component for an entity, type-safe.
func (w *World) GetHealth(e Entity) (*components.Health, bool) {
	comp, ok := w.Health[e]
	return comp, ok
}

// GetMana returns the Mana component for an entity, type-safe.
func (w *World) GetMana(e Entity) (*components.Mana, bool) {
	comp, ok := w.Mana[e]
	return comp, ok
}

// GetAttack returns the Attack component for an entity, type-safe.
func (w *World) GetAttack(e Entity) (*components.Attack, bool) {
	comp, ok := w.Attack[e]
	return comp, ok
}

// GetTraits returns the Traits component for an entity, type-safe.
func (w *World) GetTraits(e Entity) (*components.Traits, bool) {
	comp, ok := w.Traits[e]
	return comp, ok
}

// GetChampionInfo returns the ChampionInfo component for an entity, type-safe.
func (w *World) GetChampionInfo(e Entity) (*components.ChampionInfo, bool) {
	comp, ok := w.ChampionInfo[e]
	return comp, ok
}

// GetPosition returns the Position component for an entity, type-safe.
func (w *World) GetPosition(e Entity) (*components.Position, bool) {
	comp, ok := w.Position[e]
	return comp, ok
}

// GetTeam returns the Team component for an entity, type-safe.
func (w *World) GetTeam(e Entity) (*components.Team, bool) {
	comp, ok := w.Team[e]
	return comp, ok
}

// GetChampionByName returns the first entity with the specified champion name.
func (w *World) GetChampionByName(name string) (Entity, bool) {
	for e, info := range w.ChampionInfo {
		if info.Name == name {
			return e, true
		}
	}
	return 0, false // Not found
}

// GetItemEffect returns the ItemEffect component for an entity, type-safe.
func (w *World) GetItemEffect(e Entity) (*components.ItemEffect, bool) {
	comp, ok := w.Item[e]
	return comp, ok
}

// GetEquipment returns the Equipment component for an entity, type-safe.
func (w *World) GetEquipment(e Entity) (*components.Equipment, bool) {
	comp, ok := w.Equipment[e]
	return comp, ok
}
