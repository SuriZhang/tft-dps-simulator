package ecs

import (
	"fmt"
	"log"
	"reflect"
	"sync/atomic"

	"tft-dps-simulator/internal/core/components"
	"tft-dps-simulator/internal/core/components/debuffs"
	"tft-dps-simulator/internal/core/components/items"
	"tft-dps-simulator/internal/core/components/traits"
	"tft-dps-simulator/internal/core/entity"
)

// World contains all entities and their components, stored in type-specific maps.
type World struct {
	nextEntityId uint32
	// --- Component Maps ---
	// Using pointers (*components.Type) allows checking for nil to see if an entity.Entity has the component.
	Health                   map[entity.Entity]*components.Health
	Mana                     map[entity.Entity]*components.Mana
	Attack                   map[entity.Entity]*components.Attack
	Traits                   map[entity.Entity]*components.Traits
	ChampionInfo             map[entity.Entity]*components.ChampionInfo
	Position                 map[entity.Entity]*components.Position
	Team                     map[entity.Entity]*components.Team
	Item                     map[entity.Entity]*items.ItemStaticEffect
	Equipment                map[entity.Entity]*components.Equipment
	CanAbilityCritFromTraits map[entity.Entity]*components.CanAbilityCritFromTraits
	CanAbilityCritFromItems  map[entity.Entity]*components.CanAbilityCritFromItems
	Spell                    map[entity.Entity]*components.Spell
	Crit                     map[entity.Entity]*components.Crit
	State                    map[entity.Entity]*components.State
	DamageStats              map[entity.Entity]*components.DamageStats

	// --- Debuff Components ---
	ShredEffects  map[entity.Entity]*debuffs.ShredEffect
	SunderEffects map[entity.Entity]*debuffs.SunderEffect
	WoundEffects  map[entity.Entity]*debuffs.WoundEffect
	BurnEffects   map[entity.Entity]*debuffs.BurnEffect

	// --- Dynamic Item Effect Components ---
	ArchangelsStaffEffects   map[entity.Entity]*items.ArchangelsStaffEffect
	QuicksilverEffects       map[entity.Entity]*items.QuicksilverEffect
	TitansResolveEffects     map[entity.Entity]*items.TitansResolveEffect
	GuinsoosRagebladeEffects map[entity.Entity]*items.GuinsoosRagebladeEffect
	SpiritVisageEffects      map[entity.Entity]*items.SpiritVisageEffect
	KrakensFuryEffects       map[entity.Entity]*items.KrakensFuryEffect
	SpearOfShojinEffects     map[entity.Entity]*items.SpearOfShojinEffect
	BlueBuffEffects          map[entity.Entity]*items.BlueBuffEffect
	FlickerbladeEffects      map[entity.Entity]*items.FlickerbladeEffect
	NashorsToothEffects      map[entity.Entity]*items.NashorsToothEffect

	// --- Trait Effect Components ---
	RapidfireEffects map[entity.Entity]*traits.RapidfireEffect
}

// NewWorld creates a new empty world, initializing all component maps.
func NewWorld() *World {
	return &World{
		nextEntityId: 0,
		// Initialize maps
		Health:                   make(map[entity.Entity]*components.Health),
		Mana:                     make(map[entity.Entity]*components.Mana),
		Attack:                   make(map[entity.Entity]*components.Attack),
		Traits:                   make(map[entity.Entity]*components.Traits),
		ChampionInfo:             make(map[entity.Entity]*components.ChampionInfo),
		Position:                 make(map[entity.Entity]*components.Position),
		Team:                     make(map[entity.Entity]*components.Team),
		Item:                     make(map[entity.Entity]*items.ItemStaticEffect),
		Equipment:                make(map[entity.Entity]*components.Equipment),
		CanAbilityCritFromTraits: make(map[entity.Entity]*components.CanAbilityCritFromTraits),
		CanAbilityCritFromItems:  make(map[entity.Entity]*components.CanAbilityCritFromItems),
		Spell:                    make(map[entity.Entity]*components.Spell),
		Crit:                     make(map[entity.Entity]*components.Crit),
		State:                    make(map[entity.Entity]*components.State),
		DamageStats:              make(map[entity.Entity]*components.DamageStats),

		// --- Debuff Components ---
		ShredEffects:  make(map[entity.Entity]*debuffs.ShredEffect),
		SunderEffects: make(map[entity.Entity]*debuffs.SunderEffect),
		WoundEffects:  make(map[entity.Entity]*debuffs.WoundEffect),
		BurnEffects:   make(map[entity.Entity]*debuffs.BurnEffect),

		// --- Dynamic Item Effect Components ---
		ArchangelsStaffEffects:   make(map[entity.Entity]*items.ArchangelsStaffEffect),
		QuicksilverEffects:       make(map[entity.Entity]*items.QuicksilverEffect),
		TitansResolveEffects:     make(map[entity.Entity]*items.TitansResolveEffect),
		GuinsoosRagebladeEffects: make(map[entity.Entity]*items.GuinsoosRagebladeEffect),
		SpiritVisageEffects:      make(map[entity.Entity]*items.SpiritVisageEffect),
		KrakensFuryEffects:       make(map[entity.Entity]*items.KrakensFuryEffect),
		SpearOfShojinEffects:     make(map[entity.Entity]*items.SpearOfShojinEffect),
		BlueBuffEffects:          make(map[entity.Entity]*items.BlueBuffEffect),
		FlickerbladeEffects:      make(map[entity.Entity]*items.FlickerbladeEffect),
		NashorsToothEffects:      make(map[entity.Entity]*items.NashorsToothEffect),

		// --- Traits ---
		RapidfireEffects: make(map[entity.Entity]*traits.RapidfireEffect),
	}
}

// NewEntity generates a new unique entity ID using the global NewEntity function.
// Note: This only reserves the ID; components must be added separately.
func (w *World) NewEntity() entity.Entity {
	// Use the world's internal atomic counter
	return entity.Entity(atomic.AddUint32(&w.nextEntityId, 1))
}

// RemoveEntity removes an entity and all its associated components from the world.
func (w *World) RemoveEntity(e entity.Entity) {
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
	delete(w.CanAbilityCritFromTraits, e)
	delete(w.CanAbilityCritFromItems, e)
	delete(w.Spell, e)
	delete(w.Crit, e)
	delete(w.State, e)
	delete(w.DamageStats, e)
	// --- Debuff Components ---
	delete(w.ShredEffects, e)
	delete(w.SunderEffects, e)
	delete(w.WoundEffects, e)
	delete(w.BurnEffects, e)
	// --- Dynamic Item Effect Components ---
	delete(w.ArchangelsStaffEffects, e)
	delete(w.QuicksilverEffects, e)
	delete(w.TitansResolveEffects, e)
	delete(w.GuinsoosRagebladeEffects, e)
	delete(w.SpiritVisageEffects, e)
	delete(w.KrakensFuryEffects, e)
	delete(w.SpearOfShojinEffects, e)
	delete(w.BlueBuffEffects, e)
	delete(w.FlickerbladeEffects, e)
	delete(w.NashorsToothEffects, e)
	// Traits
	delete(w.RapidfireEffects, e)
	// Delete from other maps here...
}

// AddComponent adds a component to an entity.Entity. It uses a type switch
// to place the component in the correct map.
// Returns an error if the component type is not recognized by the World struct.
func (w *World) AddComponent(e entity.Entity, component interface{}) error {
	// Ensure component is not nil
	if component == nil {
		return fmt.Errorf("cannot add nil component to entity.Entity %d", e)
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
	case items.ItemStaticEffect:
		w.Item[e] = &c
	case *items.ItemStaticEffect:
		w.Item[e] = c
	case components.Equipment:
		w.Equipment[e] = &c
	case *components.Equipment:
		w.Equipment[e] = c
	case components.CanAbilityCritFromTraits:
		w.CanAbilityCritFromTraits[e] = &c
	case *components.CanAbilityCritFromTraits:
		w.CanAbilityCritFromTraits[e] = c
	case components.CanAbilityCritFromItems:
		w.CanAbilityCritFromItems[e] = &c
	case *components.CanAbilityCritFromItems:
		w.CanAbilityCritFromItems[e] = c
	case components.Spell:
		w.Spell[e] = &c
	case *components.Spell:
		w.Spell[e] = c
	case components.Crit:
		w.Crit[e] = &c
	case *components.Crit:
		w.Crit[e] = c
	case components.State:
		w.State[e] = &c
	case *components.State:
		w.State[e] = c
	case components.DamageStats:
		w.DamageStats[e] = &c
	case *components.DamageStats:
		w.DamageStats[e] = c
	// --- Debuff Components ---
	case debuffs.ShredEffect:
		w.ShredEffects[e] = &c
	case *debuffs.ShredEffect:
		w.ShredEffects[e] = c
	case debuffs.SunderEffect:
		w.SunderEffects[e] = &c
	case *debuffs.SunderEffect:
		w.SunderEffects[e] = c
	case debuffs.WoundEffect:
		w.WoundEffects[e] = &c
	case *debuffs.WoundEffect:
		w.WoundEffects[e] = c
	case debuffs.BurnEffect:
		w.BurnEffects[e] = &c
	case *debuffs.BurnEffect:
		w.BurnEffects[e] = c
	// --- Dynamic Item Effect Components ---
	case items.ArchangelsStaffEffect:
		w.ArchangelsStaffEffects[e] = &c
	case *items.ArchangelsStaffEffect:
		w.ArchangelsStaffEffects[e] = c
	case items.QuicksilverEffect:
		w.QuicksilverEffects[e] = &c
	case *items.QuicksilverEffect:
		w.QuicksilverEffects[e] = c
	case items.TitansResolveEffect:
		w.TitansResolveEffects[e] = &c
	case *items.TitansResolveEffect:
		w.TitansResolveEffects[e] = c
	case items.GuinsoosRagebladeEffect:
		w.GuinsoosRagebladeEffects[e] = &c
	case *items.GuinsoosRagebladeEffect:
		w.GuinsoosRagebladeEffects[e] = c
	case items.SpiritVisageEffect:
		w.SpiritVisageEffects[e] = &c
	case *items.SpiritVisageEffect:
		w.SpiritVisageEffects[e] = c
	case items.KrakensFuryEffect:
		w.KrakensFuryEffects[e] = &c
	case *items.KrakensFuryEffect:
		w.KrakensFuryEffects[e] = c
	case items.SpearOfShojinEffect:
		w.SpearOfShojinEffects[e] = &c
	case *items.SpearOfShojinEffect:
		w.SpearOfShojinEffects[e] = c
	case items.BlueBuffEffect:
		w.BlueBuffEffects[e] = &c
	case *items.BlueBuffEffect:
		w.BlueBuffEffects[e] = c
	case items.FlickerbladeEffect:
		w.FlickerbladeEffects[e] = &c
	case *items.FlickerbladeEffect:
		w.FlickerbladeEffects[e] = c
	case items.NashorsToothEffect:
		w.NashorsToothEffects[e] = &c
	case *items.NashorsToothEffect:
		w.NashorsToothEffects[e] = c
	// Traits
	case traits.RapidfireEffect:
		w.RapidfireEffects[e] = &c
	case *traits.RapidfireEffect:
		w.RapidfireEffects[e] = c
	// Add cases for other component types here...
	default:
		// Use reflection to get the type name for the error message
		return fmt.Errorf("unknown component type: %v", reflect.TypeOf(component))
	}
	return nil
}

// GetComponent retrieves a component of a specific type for an entity.Entity.
// It returns the component (as interface{}) and true if found, otherwise nil and false.
// This generic version is kept for flexibility but type-safe getters are preferred.
func (w *World) GetComponent(e entity.Entity, componentType reflect.Type) (interface{}, bool) {
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
	case reflect.TypeOf(items.ItemStaticEffect{}):
		comp, ok := w.Item[e]
		return comp, ok
	case reflect.TypeOf(components.Equipment{}):
		comp, ok := w.Equipment[e]
		return comp, ok
	case reflect.TypeOf(components.CanAbilityCritFromTraits{}):
		comp, ok := w.CanAbilityCritFromTraits[e]
		return comp, ok
	case reflect.TypeOf(components.CanAbilityCritFromItems{}):
		comp, ok := w.CanAbilityCritFromItems[e]
		return comp, ok
	case reflect.TypeOf(components.Spell{}):
		comp, ok := w.Spell[e]
		return comp, ok
	case reflect.TypeOf(components.Crit{}):
		comp, ok := w.Crit[e]
		return comp, ok
	case reflect.TypeOf(components.State{}):
		comp, ok := w.State[e]
		return comp, ok
	case reflect.TypeOf(components.DamageStats{}):
		comp, ok := w.DamageStats[e]
		return comp, ok
	// --- Debuff Components ---
	case reflect.TypeOf(debuffs.ShredEffect{}):
		comp, ok := w.ShredEffects[e]
		return comp, ok
	case reflect.TypeOf(debuffs.SunderEffect{}):
		comp, ok := w.SunderEffects[e]
		return comp, ok
	case reflect.TypeOf(debuffs.WoundEffect{}):
		comp, ok := w.WoundEffects[e]
		return comp, ok
	case reflect.TypeOf(debuffs.BurnEffect{}):
		comp, ok := w.BurnEffects[e]
		return comp, ok
	// --- Dynamic Item Effect Components ---
	case reflect.TypeOf(items.ArchangelsStaffEffect{}):
		comp, ok := w.ArchangelsStaffEffects[e]
		return comp, ok
	case reflect.TypeOf(items.QuicksilverEffect{}):
		comp, ok := w.QuicksilverEffects[e]
		return comp, ok
	case reflect.TypeOf(items.TitansResolveEffect{}):
		comp, ok := w.TitansResolveEffects[e]
		return comp, ok
	case reflect.TypeOf(items.GuinsoosRagebladeEffect{}):
		comp, ok := w.GuinsoosRagebladeEffects[e]
		return comp, ok
	case reflect.TypeOf(items.SpiritVisageEffect{}):
		comp, ok := w.SpiritVisageEffects[e]
		return comp, ok
	case reflect.TypeOf(items.KrakensFuryEffect{}):
		comp, ok := w.KrakensFuryEffects[e]
		return comp, ok
	case reflect.TypeOf(items.SpearOfShojinEffect{}):
		comp, ok := w.SpearOfShojinEffects[e]
		return comp, ok
	case reflect.TypeOf(items.BlueBuffEffect{}):
		comp, ok := w.BlueBuffEffects[e]
		return comp, ok
	case reflect.TypeOf(items.FlickerbladeEffect{}):
		comp, ok := w.FlickerbladeEffects[e]
		return comp, ok
	case reflect.TypeOf(items.NashorsToothEffect{}):
		comp, ok := w.NashorsToothEffects[e]
		return comp, ok
	// Traits
	case reflect.TypeOf(traits.RapidfireEffect{}):
		comp, ok := w.RapidfireEffects[e]
		return comp, ok
	// Add cases for other component types here...
	default:
		return nil, false
	}
}

// HasComponent checks if an entity.Entity possesses a component of the specified type.
func (w *World) HasComponent(e entity.Entity, componentType reflect.Type) bool {
	_, ok := w.GetComponent(e, componentType) // Leverage existing logic
	return ok
}

// RemoveComponent removes a specific component type from an entity.Entity.
func (w *World) RemoveComponent(e entity.Entity, componentType reflect.Type) {
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
	case reflect.TypeOf(items.ItemStaticEffect{}):
		delete(w.Item, e)
	case reflect.TypeOf(components.Equipment{}):
		delete(w.Equipment, e)
	case reflect.TypeOf(components.CanAbilityCritFromTraits{}):
		delete(w.CanAbilityCritFromTraits, e)
	case reflect.TypeOf(components.CanAbilityCritFromItems{}):
		delete(w.CanAbilityCritFromItems, e)
	case reflect.TypeOf(components.Spell{}):
		delete(w.Spell, e)
	case reflect.TypeOf(components.Crit{}):
		delete(w.Crit, e)
	case reflect.TypeOf(components.State{}):
		delete(w.State, e)
	case reflect.TypeOf(components.DamageStats{}):
		delete(w.DamageStats, e)
	// --- Debuff Components ---
	case reflect.TypeOf(debuffs.ShredEffect{}):
		delete(w.ShredEffects, e)
	case reflect.TypeOf(debuffs.SunderEffect{}):
		delete(w.SunderEffects, e)
	case reflect.TypeOf(debuffs.WoundEffect{}):
		delete(w.WoundEffects, e)
	case reflect.TypeOf(debuffs.BurnEffect{}):
		delete(w.BurnEffects, e)
	// --- Dynamic Item Effect Components ---
	case reflect.TypeOf(items.ArchangelsStaffEffect{}):
		delete(w.ArchangelsStaffEffects, e)
	case reflect.TypeOf(items.QuicksilverEffect{}):
		delete(w.QuicksilverEffects, e)
	case reflect.TypeOf(items.TitansResolveEffect{}):
		delete(w.TitansResolveEffects, e)
	case reflect.TypeOf(items.GuinsoosRagebladeEffect{}):
		delete(w.GuinsoosRagebladeEffects, e)
	case reflect.TypeOf(items.SpiritVisageEffect{}):
		delete(w.SpiritVisageEffects, e)
	case reflect.TypeOf(items.KrakensFuryEffect{}):
		delete(w.KrakensFuryEffects, e)
	case reflect.TypeOf(items.SpearOfShojinEffect{}):
		delete(w.SpearOfShojinEffects, e)
	case reflect.TypeOf(items.BlueBuffEffect{}):
		delete(w.BlueBuffEffects, e)
	case reflect.TypeOf(items.FlickerbladeEffect{}):
		delete(w.FlickerbladeEffects, e)
	case reflect.TypeOf(items.NashorsToothEffect{}):
		delete(w.NashorsToothEffects, e)
	// Traits
	case reflect.TypeOf(traits.RapidfireEffect{}):
		delete(w.RapidfireEffects, e)
	// Add cases for other component types here...
	default:
		log.Printf("Warning: Attempted to remove unknown component type %v from entity.Entity %d\n", componentType, e)
	}
}

// GetEntitiesWithComponents returns a slice of entities that possess ALL the specified component types.
func (w *World) GetEntitiesWithComponents(componentTypes ...reflect.Type) []entity.Entity {
	if len(componentTypes) == 0 {
		return []entity.Entity{}
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
		return []entity.Entity{} // If the smallest map is empty, no entities can have all components.
	}

	// Initialize candidates with entities from the smallest map
	initialSet := w.getEntitiesForType(firstType)
	candidates := make(map[entity.Entity]bool, len(initialSet))
	for _, e := range initialSet {
		candidates[e] = true
	}

	// Filter by the remaining component types
	for _, ct := range componentTypes {
		if ct == firstType { // Skip the type we already used
			continue
		}

		// Check entities for the current component type
		currentTypeSet := make(map[entity.Entity]bool)
		for _, e := range w.getEntitiesForType(ct) {
			currentTypeSet[e] = true
		}

		// Filter candidates: keep only those present in currentTypeSet
		for e := range candidates {
			if !currentTypeSet[e] {
				delete(candidates, e) // Remove if entity.Entity doesn't have the current component type
			}
		}

		// Early exit if no candidates remain
		if len(candidates) == 0 {
			return []entity.Entity{}
		}
	}

	// Convert the final candidate map keys to a slice
	result := make([]entity.Entity, 0, len(candidates))
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
	case reflect.TypeOf(items.ItemStaticEffect{}):
		return len(w.Item)
	case reflect.TypeOf(components.Equipment{}):
		return len(w.Equipment)
	case reflect.TypeOf(components.CanAbilityCritFromTraits{}):
		return len(w.CanAbilityCritFromTraits)
	case reflect.TypeOf(components.CanAbilityCritFromItems{}):
		return len(w.CanAbilityCritFromItems)
	case reflect.TypeOf(components.Spell{}):
		return len(w.Spell)
	case reflect.TypeOf(components.Crit{}):
		return len(w.Crit)
	case reflect.TypeOf(components.State{}):
		return len(w.State)
	case reflect.TypeOf(components.DamageStats{}):
		return len(w.DamageStats)
	// --- Debuff Components ---
	case reflect.TypeOf(debuffs.ShredEffect{}):
		return len(w.ShredEffects)
	case reflect.TypeOf(debuffs.SunderEffect{}):
		return len(w.SunderEffects)
	case reflect.TypeOf(debuffs.WoundEffect{}):
		return len(w.WoundEffects)
	case reflect.TypeOf(debuffs.BurnEffect{}):
		return len(w.BurnEffects)
	// --- Dynamic Item Effect Components ---
	case reflect.TypeOf(items.ArchangelsStaffEffect{}):
		return len(w.ArchangelsStaffEffects)
	case reflect.TypeOf(items.QuicksilverEffect{}):
		return len(w.QuicksilverEffects)
	case reflect.TypeOf(items.TitansResolveEffect{}):
		return len(w.TitansResolveEffects)
	case reflect.TypeOf(items.GuinsoosRagebladeEffect{}):
		return len(w.GuinsoosRagebladeEffects)
	case reflect.TypeOf(items.SpiritVisageEffect{}):
		return len(w.SpiritVisageEffects)
	case reflect.TypeOf(items.KrakensFuryEffect{}):
		return len(w.KrakensFuryEffects)
	case reflect.TypeOf(items.SpearOfShojinEffect{}):
		return len(w.SpearOfShojinEffects)
	case reflect.TypeOf(items.BlueBuffEffect{}):
		return len(w.BlueBuffEffects)
	case reflect.TypeOf(items.FlickerbladeEffect{}):
		return len(w.FlickerbladeEffects)
	case reflect.TypeOf(items.NashorsToothEffect{}):
		return len(w.NashorsToothEffects)
	// Traits
	case reflect.TypeOf(traits.RapidfireEffect{}):
		return len(w.RapidfireEffects)
	// Add cases for other component types...
	default:
		return 0
	}
}

// Helper function for GetEntitiesWithComponents to get entities for a single type
func (w *World) getEntitiesForType(componentType reflect.Type) []entity.Entity {
	var entities []entity.Entity
	switch componentType {
	case reflect.TypeOf(components.Health{}):
		entities = make([]entity.Entity, 0, len(w.Health))
		for e := range w.Health {
			entities = append(entities, e)
		}
	case reflect.TypeOf(components.Mana{}):
		entities = make([]entity.Entity, 0, len(w.Mana))
		for e := range w.Mana {
			entities = append(entities, e)
		}
	case reflect.TypeOf(components.Attack{}):
		entities = make([]entity.Entity, 0, len(w.Attack))
		for e := range w.Attack {
			entities = append(entities, e)
		}
	case reflect.TypeOf(components.Traits{}):
		entities = make([]entity.Entity, 0, len(w.Traits))
		for e := range w.Traits {
			entities = append(entities, e)
		}
	case reflect.TypeOf(components.ChampionInfo{}):
		entities = make([]entity.Entity, 0, len(w.ChampionInfo))
		for e := range w.ChampionInfo {
			entities = append(entities, e)
		}
	case reflect.TypeOf(components.Position{}):
		entities = make([]entity.Entity, 0, len(w.Position))
		for e := range w.Position {
			entities = append(entities, e)
		}
	case reflect.TypeOf(components.Team{}):
		entities = make([]entity.Entity, 0, len(w.Team))
		for e := range w.Team {
			entities = append(entities, e)
		}
	case reflect.TypeOf(items.ItemStaticEffect{}):
		entities = make([]entity.Entity, 0, len(w.Item))
		for e := range w.Item {
			entities = append(entities, e)
		}
	case reflect.TypeOf(components.Equipment{}):
		entities = make([]entity.Entity, 0, len(w.Equipment))
		for e := range w.Equipment {
			entities = append(entities, e)
		}
	case reflect.TypeOf(components.CanAbilityCritFromTraits{}):
		entities = make([]entity.Entity, 0, len(w.CanAbilityCritFromTraits))
		for e := range w.CanAbilityCritFromTraits {
			entities = append(entities, e)
		}
	case reflect.TypeOf(components.CanAbilityCritFromItems{}):
		entities = make([]entity.Entity, 0, len(w.CanAbilityCritFromItems))
		for e := range w.CanAbilityCritFromItems {
			entities = append(entities, e)
		}
	case reflect.TypeOf(components.Spell{}):
		entities = make([]entity.Entity, 0, len(w.Spell))
		for e := range w.Spell {
			entities = append(entities, e)
		}
	case reflect.TypeOf(components.Crit{}):
		entities = make([]entity.Entity, 0, len(w.Crit))
		for e := range w.Crit {
			entities = append(entities, e)
		}
	case reflect.TypeOf(components.State{}):
		entities = make([]entity.Entity, 0, len(w.State))
		for e := range w.State {
			entities = append(entities, e)
		}
	case reflect.TypeOf(components.DamageStats{}):
		entities = make([]entity.Entity, 0, len(w.DamageStats))
		for e := range w.DamageStats {
			entities = append(entities, e)
		}
	// --- Debuff Components ---
	case reflect.TypeOf(debuffs.ShredEffect{}):
		entities = make([]entity.Entity, 0, len(w.ShredEffects))
		for e := range w.ShredEffects {
			entities = append(entities, e)
		}
	case reflect.TypeOf(debuffs.SunderEffect{}):
		entities = make([]entity.Entity, 0, len(w.SunderEffects))
		for e := range w.SunderEffects {
			entities = append(entities, e)
		}
	case reflect.TypeOf(debuffs.WoundEffect{}):
		entities = make([]entity.Entity, 0, len(w.WoundEffects))
		for e := range w.WoundEffects {
			entities = append(entities, e)
		}
	case reflect.TypeOf(debuffs.BurnEffect{}):
		entities = make([]entity.Entity, 0, len(w.BurnEffects))
		for e := range w.BurnEffects {
			entities = append(entities, e)
		}
	// --- Dynamic Item Effect Components ---
	case reflect.TypeOf(items.ArchangelsStaffEffect{}):
		entities = make([]entity.Entity, 0, len(w.ArchangelsStaffEffects))
		for e := range w.ArchangelsStaffEffects {
			entities = append(entities, e)
		}
	case reflect.TypeOf(items.QuicksilverEffect{}):
		entities = make([]entity.Entity, 0, len(w.QuicksilverEffects))
		for e := range w.QuicksilverEffects {
			entities = append(entities, e)
		}
	case reflect.TypeOf(items.TitansResolveEffect{}):
		entities = make([]entity.Entity, 0, len(w.TitansResolveEffects))
		for e := range w.TitansResolveEffects {
			entities = append(entities, e)
		}
	case reflect.TypeOf(items.GuinsoosRagebladeEffect{}):
		entities = make([]entity.Entity, 0, len(w.GuinsoosRagebladeEffects))
		for e := range w.GuinsoosRagebladeEffects {
			entities = append(entities, e)
		}
	case reflect.TypeOf(items.SpiritVisageEffect{}):
		entities = make([]entity.Entity, 0, len(w.SpiritVisageEffects))
		for e := range w.SpiritVisageEffects {
			entities = append(entities, e)
		}
	case reflect.TypeOf(items.KrakensFuryEffect{}):
		entities = make([]entity.Entity, 0, len(w.KrakensFuryEffects))
		for e := range w.KrakensFuryEffects {
			entities = append(entities, e)
		}
	case reflect.TypeOf(items.SpearOfShojinEffect{}):
		entities = make([]entity.Entity, 0, len(w.SpearOfShojinEffects))
		for e := range w.SpearOfShojinEffects {
			entities = append(entities, e)
		}
	case reflect.TypeOf(items.BlueBuffEffect{}):
		entities = make([]entity.Entity, 0, len(w.BlueBuffEffects))
		for e := range w.BlueBuffEffects {
			entities = append(entities, e)
		}
	case reflect.TypeOf(items.FlickerbladeEffect{}):
		entities = make([]entity.Entity, 0, len(w.FlickerbladeEffects))
		for e := range w.FlickerbladeEffects {
			entities = append(entities, e)
		}
	case reflect.TypeOf(items.NashorsToothEffect{}):
		entities = make([]entity.Entity, 0, len(w.NashorsToothEffects))
		for e := range w.NashorsToothEffects {
			entities = append(entities, e)
		}
	// Traits
	case reflect.TypeOf(traits.RapidfireEffect{}):
		entities = make([]entity.Entity, 0, len(w.RapidfireEffects))
		for e := range w.RapidfireEffects {
			entities = append(entities, e)
		}
	// Add cases for other component types...
	default:
		return []entity.Entity{} // Return empty slice for unknown types
	}
	return entities
}

// --- Type-Safe Getters (Recommended) ---

// GetHealth returns the Health component for an entity.Entity, type-safe.
func (w *World) GetHealth(e entity.Entity) (*components.Health, bool) {
	comp, ok := w.Health[e]
	return comp, ok
}

// GetMana returns the Mana component for an entity.Entity, type-safe.
func (w *World) GetMana(e entity.Entity) (*components.Mana, bool) {
	comp, ok := w.Mana[e]
	return comp, ok
}

// GetAttack returns the Attack component for an entity.Entity, type-safe.
func (w *World) GetAttack(e entity.Entity) (*components.Attack, bool) {
	comp, ok := w.Attack[e]
	return comp, ok
}

// GetTraits returns the Traits component for an entity.Entity, type-safe.
func (w *World) GetTraits(e entity.Entity) (*components.Traits, bool) {
	comp, ok := w.Traits[e]
	return comp, ok
}

// GetChampionInfo returns the ChampionInfo component for an entity.Entity, type-safe.
func (w *World) GetChampionInfo(e entity.Entity) (*components.ChampionInfo, bool) {
	comp, ok := w.ChampionInfo[e]
	return comp, ok
}

// GetPosition returns the Position component for an entity.Entity, type-safe.
func (w *World) GetPosition(e entity.Entity) (*components.Position, bool) {
	comp, ok := w.Position[e]
	return comp, ok
}

// GetTeam returns the Team component for an entity.Entity, type-safe.
func (w *World) GetTeam(e entity.Entity) (*components.Team, bool) {
	comp, ok := w.Team[e]
	return comp, ok
}

// GetChampionByName returns the first entity.Entity with the specified champion name.
func (w *World) GetChampionByName(name string) (entity.Entity, bool) {
	for e, info := range w.ChampionInfo {
		if info.Name == name {
			return e, true
		}
	}
	return 0, false // Not found
}

// GetItemEffect returns the ItemEffect component for an entity.Entity, type-safe.
func (w *World) GetItemEffect(e entity.Entity) (*items.ItemStaticEffect, bool) {
	comp, ok := w.Item[e]
	return comp, ok
}

// GetEquipment returns the Equipment component for an entity.Entity, type-safe.
func (w *World) GetEquipment(e entity.Entity) (*components.Equipment, bool) {
	comp, ok := w.Equipment[e]
	return comp, ok
}

// GetCanAbilityCritFromTraits returns the CanAbilityCritFromTraits component for an entity.Entity, type-safe.
func (w *World) GetCanAbilityCritFromTraits(e entity.Entity) (*components.CanAbilityCritFromTraits, bool) {
	comp, ok := w.CanAbilityCritFromTraits[e]
	return comp, ok
}

// GetCanAbilityCritFromItems returns the CanAbilityCritFromItems component for an entity.Entity, type-safe.
func (w *World) GetCanAbilityCritFromItems(e entity.Entity) (*components.CanAbilityCritFromItems, bool) {
	comp, ok := w.CanAbilityCritFromItems[e]
	return comp, ok
}

// GetSpell returns the Spell component for an entity.Entity, type-safe.
func (w *World) GetSpell(e entity.Entity) (*components.Spell, bool) {
	comp, ok := w.Spell[e]
	return comp, ok
}

// GetCrit returns the Crit component for an entity.Entity, type-safe.
func (w *World) GetCrit(e entity.Entity) (*components.Crit, bool) {
	comp, ok := w.Crit[e]
	return comp, ok
}

// GetState returns the State component for an entity.Entity, type-safe.
func (w *World) GetState(e entity.Entity) (*components.State, bool) {
	comp, ok := w.State[e]
	return comp, ok
}

// GetDamageStats returns the DamageStats component for an entity.Entity, type-safe.
func (w *World) GetDamageStats(e entity.Entity) (*components.DamageStats, bool) {
	comp, ok := w.DamageStats[e]
	return comp, ok
}

// GetArchangelsStaffEffect returns the ArchangelsEffect component for an entity.Entity, type-safe.
func (w *World) GetArchangelsStaffEffect(e entity.Entity) (*items.ArchangelsStaffEffect, bool) {
	comp, ok := w.ArchangelsStaffEffects[e]
	return comp, ok
}

// GetQuicksilverEffect returns the QuicksilverEffect component for an entity.Entity, type-safe.
func (w *World) GetQuicksilverEffect(e entity.Entity) (*items.QuicksilverEffect, bool) {
	comp, ok := w.QuicksilverEffects[e]
	return comp, ok
}

// GetTitansResolveEffect returns the TitansResolveEffect component for an entity.Entity, type-safe.
func (w *World) GetTitansResolveEffect(e entity.Entity) (*items.TitansResolveEffect, bool) {
	comp, ok := w.TitansResolveEffects[e]
	return comp, ok
}

// GetGuinsoosRagebladeEffect returns the GuinsoosRagebladeEffect component for an entity.Entity, type-safe.
func (w *World) GetGuinsoosRagebladeEffect(e entity.Entity) (*items.GuinsoosRagebladeEffect, bool) {
	comp, ok := w.GuinsoosRagebladeEffects[e]
	return comp, ok
}

// GetSpiritVisageEffect returns the SpiritVisageEffect component for an entity.Entity, type-safe.
func (w *World) GetSpiritVisageEffect(e entity.Entity) (*items.SpiritVisageEffect, bool) {
	comp, ok := w.SpiritVisageEffects[e]
	return comp, ok
}

// GetKrakensFuryEffect returns the KrakensFuryEffect component for an entity.Entity, type-safe.
func (w *World) GetKrakensFuryEffect(e entity.Entity) (*items.
	KrakensFuryEffect, bool) {
	comp, ok := w.KrakensFuryEffects[e]
	return comp, ok
}

// GetSpearOfShojinEffect returns the SpearOfShojinEffect component for an entity.Entity, type-safe.
func (w *World) GetSpearOfShojinEffect(e entity.Entity) (*items.SpearOfShojinEffect, bool) {
	comp, exists := w.SpearOfShojinEffects[e]
	return comp, exists
}

// GetBlueBuffEffect returns the BlueBuffEffect component for an entity.Entity, type-safe.
func (w *World) GetBlueBuffEffect(e entity.Entity) (*items.BlueBuffEffect, bool) {
	comp, ok := w.BlueBuffEffects[e]
	return comp, ok
}

// GetFlickerbladeEffect returns the FlickerbladeEffect component for an entity.Entity, type-safe.
func (w *World) GetFlickerbladeEffect(e entity.Entity) (*items.FlickerbladeEffect, bool) {
	comp, ok := w.FlickerbladeEffects[e]
	return comp, ok
}

// GetNashorsToothEffect returns the NashorsToothEffect component for an entity, type-safe.
func (w *World) GetNashorsToothEffect(e entity.Entity) (*items.NashorsToothEffect, bool) {
	effect, ok := w.NashorsToothEffects[e]
	return effect, ok
}

// Traits
// GetRapidfireEffect returns the RapidfireEffect component for an entity.Entity, type-safe.
func (w *World) GetRapidfireEffect(e entity.Entity) (*traits.RapidfireEffect, bool) {
	comp, ok := w.RapidfireEffects[e]
	return comp, ok
}

// Debuffs
// GetShredEffect returns the ShredEffect component for an entity, type-safe.
func (w *World) GetShredEffect(e entity.Entity) (*debuffs.ShredEffect, bool) {
	comp, ok := w.ShredEffects[e]
	return comp, ok
}

// GetSunderEffect returns the SunderEffect component for an entity, type-safe.
func (w *World) GetSunderEffect(e entity.Entity) (*debuffs.SunderEffect, bool) {
	comp, ok := w.SunderEffects[e]
	return comp, ok
}

// GetWoundEffect returns the WoundEffect component for an entity, type-safe.
func (w *World) GetWoundEffect(e entity.Entity) (*debuffs.WoundEffect, bool) {
	comp, ok := w.WoundEffects[e]
	return comp, ok
}

// GetBurnEffect returns the BurnEffect component for an entity, type-safe.
func (w *World) GetBurnEffect(e entity.Entity) (*debuffs.BurnEffect, bool) {
	comp, ok := w.BurnEffects[e]
	return comp, ok
}
