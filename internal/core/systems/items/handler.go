package itemsys

import (
	"tft-dps-simulator/internal/core/ecs"
	eventsys "tft-dps-simulator/internal/core/systems/events"
)

// ItemHandler defines the interface for item-specific logic.
// Each dynamic item should have an implementation of this interface.
type ItemHandler interface {
    // OnEquip is called by the EquipmentManager when the item is added to an entity.
    // It should set up item-specific components (e.g., ArchangelsEffect) and
    // enqueue initial events if necessary (e.g., first ArchangelsTickEvent).
    OnEquip(entity ecs.Entity, world *ecs.World, eventBus eventsys.EventBus)

    // ProcessEvent is called by the main ItemSystem when a relevant game event occurs.
    // The handler should determine if and how to react to the event.
    // This method will handle both general game events (like AttackLandedEvent for Guinsoo's)
    // and item-specific tick/proc events (like ArchangelsTickEvent for Archangel's Staff).
    // The handler is responsible for type-asserting the event and acting accordingly.
    ProcessEvent(event interface{}, entity ecs.Entity, world *ecs.World, eventBus eventsys.EventBus)
}