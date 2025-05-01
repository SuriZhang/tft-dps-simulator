package traitsys

import (
	"github.com/suriz/tft-dps-simulator/data"
	"github.com/suriz/tft-dps-simulator/ecs"

	// Import event system if needed by handlers
	eventsys "github.com/suriz/tft-dps-simulator/systems/events"
)

// TraitHandler defines the interface for dynamic trait logic.
type TraitHandler interface {
    // OnActivate is called once when the trait becomes active at a specific tier for a team.
    OnActivate(teamID int, effect data.Effect, world *ecs.World)

    // Handle processes game events relevant to the trait. Operates on entities with the trait component.
    // It receives the eventBus to enqueue follow-up events like RecalculateStatsEvent.
    Handle(event interface{}, entity ecs.Entity, world *ecs.World, eventBus eventsys.EventBus) // Added eventBus

    // OnDeactivate is called when the trait becomes inactive or changes tier for a team.
    // Since tiers are static, this might only be relevant if combat could restart with different tiers.
    // For single combat, Reset might be sufficient cleanup.
    OnDeactivate(teamID int, effect data.Effect, world *ecs.World) 

    // Reset is called at the end of combat to clear all state, potentially by removing components.
    Reset(world *ecs.World)
}