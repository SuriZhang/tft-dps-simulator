package itemsys

import (
	"log"
	"reflect"

	"tft-dps-simulator/internal/core/components"
	"tft-dps-simulator/internal/core/ecs"
)

// AbilityCritSystem checks for items granting the ability crit effect (JG, IE)
// and adds/removes the CanAbilitiesCrit marker component accordingly.
// It should run after equipment changes have been processed.
type AbilityCritSystem struct {
	world *ecs.World
}

// NewAbilityCritSystem creates a new AbilityCritSystem.
func NewAbilityCritSystem(world *ecs.World) *AbilityCritSystem {
	return &AbilityCritSystem{world: world}
}

// Update checks equipped items and applies the CanAbilitiesCrit marker.
func (s *AbilityCritSystem) Update() {
	// We need entities that have equipment to check.
	equipmentType := reflect.TypeOf(components.Equipment{})
	entities := s.world.GetEntitiesWithComponents(equipmentType)
	itemMarkerType := reflect.TypeOf(components.CanAbilityCritFromItems{})

	for _, entity := range entities {
		equipment, _ := s.world.GetEquipment(entity)
		hasInifityEdge := equipment.HasItem("TFT_Item_InfinityEdge")
		hasJeweledGauntlet := equipment.HasItem("TFT_Item_JeweledGauntlet")
		_, hasAbilityCritFromItemsMarker := s.world.GetComponent(entity, itemMarkerType)

		if hasInifityEdge || hasJeweledGauntlet && !hasAbilityCritFromItemsMarker {
			log.Printf("(AbilityCritSystem) Entity %d: Adding CanAbilityCritFromItems component .", entity)
			err := s.world.AddComponent(entity, &components.CanAbilityCritFromItems{})
			if err != nil {
				log.Printf("ERROR (AbilityCritSystem): Failed to add CanAbilityCritFromItems component to entity %d: %v", entity, err)
			}
		} else if !(hasInifityEdge || hasJeweledGauntlet) && hasAbilityCritFromItemsMarker {
			// Use Case: When IE/JG is removed from a champion, we need to remove the marker.
			// Remove the marker component if no relevant items are equipped
			log.Printf("(AbilityCritSystem) Entity %d: Removing CanAbilityCritFromItems component.", entity)
			s.world.RemoveComponent(entity, itemMarkerType)
		}

	}
}
