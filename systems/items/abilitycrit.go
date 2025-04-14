package itemsys

import (
	"log"
	"reflect"

	"github.com/suriz/tft-dps-simulator/components"
	"github.com/suriz/tft-dps-simulator/ecs"
	// Assuming data package provides item definitions if needed, but not strictly required here
	// "github.com/suriz/tft-dps-simulator/data"
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

	for _, entity := range entities {
		equipment, _ := s.world.GetEquipment(entity) // We know it exists from the query
		hasInifityEdge := equipment.HasItem("TFT_Item_InfinityEdge")
		hasJeweledGauntlet := equipment.HasItem("TFT_Item_JeweledGauntlet")
		if hasInifityEdge || hasJeweledGauntlet {
			log.Printf("(AbilityCritSystem) Entity %d: Adding CanAbilityCritFromItems component .", entity)
			err := s.world.AddComponent(entity, &components.CanAbilityCritFromItems{})
			if err != nil {
				log.Printf("ERROR (AbilityCritSystem): Failed to add CanAbilityCritFromItems component to entity %d: %v", entity, err)
			}
		}

	}
}
