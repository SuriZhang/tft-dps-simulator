package managers

import (
	"log"
	"reflect"

	"tft-dps-simulator/internal/core/components"
	"tft-dps-simulator/internal/core/ecs"
	"tft-dps-simulator/internal/core/entity"
	eventsys "tft-dps-simulator/internal/core/systems/events"
	itemsys "tft-dps-simulator/internal/core/systems/items"
)

// ItemManager manages dynamic item effects by dispatching events to registered ItemHandlers.
type ItemManager struct {
	world    *ecs.World
	eventBus eventsys.EventBus
}

// NewItemManager creates a new ItemSystem.
func NewItemManager(world *ecs.World, bus eventsys.EventBus) *ItemManager {
	im := &ItemManager{
		world:    world,
		eventBus: bus,
	}
	return im
}

// EnqueueInitialEvents checks equipped items and schedules the first timer events.
func (im *ItemManager) EnqueueInitialEvents() {
	log.Printf("DEBUG: EnqueueInitialEvents called in ItemManager")
	entities := im.world.GetEntitiesWithComponents(reflect.TypeOf(components.Equipment{}))

	for _, entity := range entities {
		equipment, ok := im.world.GetEquipment(entity)
		if !ok {
			log.Printf("ItemManager (EnqueueInitialEvents): Entity %d does not have Equipment component", entity)
			continue
		}

		for _, itemSlot := range equipment.GetAllItems() {
			if itemSlot == nil {
				log.Printf("ItemManager (EnqueueInitialEvents): Item slot is nil for entity %d", entity)
				continue
			}
			itemApiName := itemSlot.ApiName
			if handler, exists := itemsys.GetItemHandler(itemApiName); exists {
				log.Printf("ItemManager (EnqueueInitialEvents): Calling OnEquip for item %s on entity %d", itemApiName, entity)
				// Pass itemData, world, and eventBus to OnEquip
				handler.OnEquip(entity, im.world, im.eventBus)
			}
		}
	}
}

// CanHandle checks if the ItemSystem should process this event.
// It should generally try to handle events that item handlers might be interested in.
func (im *ItemManager) CanHandle(event interface{}) bool {
	// ItemSystem should attempt to process any event, and then delegate.
	// Specific item handlers will decide if they actually act on the event.
	// This is a broad catch-all; specific event types can be checked if performance is a concern.
	switch event.(type) {
	// Item-specific tick/proc events (examples)
	case eventsys.ArchangelsTickEvent, eventsys.QuicksilverProcEvent, eventsys.QuicksilverEndEvent, eventsys.GuinsoosRagebladeTickEvent, eventsys.EvenshroudResistActivateEvent, eventsys.EvenshroudResistDeactivateEvent:
		return true
		
	// General game events that items might react to
	case eventsys.AttackLandedEvent, eventsys.DamageAppliedEvent, eventsys.SpellLandedEvent, eventsys.KillEvent, eventsys.AssistEvent:
		return true
	default:
		return false
	}
}

// HandleEvent dispatches the event to relevant item handlers.
func (im *ItemManager) HandleEvent(event any) { // Changed interface{} to any
    uniqueInvolvedEntities := make(map[entity.Entity]struct{})

    // Determine involved entity/entities from the event
    switch evt := event.(type) {
    case eventsys.ArchangelsTickEvent:
        uniqueInvolvedEntities[evt.Entity] = struct{}{}
    case eventsys.QuicksilverProcEvent:
        uniqueInvolvedEntities[evt.Entity] = struct{}{}
    case eventsys.QuicksilverEndEvent:
        uniqueInvolvedEntities[evt.Entity] = struct{}{}
    case eventsys.GuinsoosRagebladeTickEvent:
        uniqueInvolvedEntities[evt.Entity] = struct{}{}
    case eventsys.BlueBuffDamageAmpActivateEvent:
        uniqueInvolvedEntities[evt.Entity] = struct{}{}
    case eventsys.BlueBuffDamageAmpDeactivateEvent:
        uniqueInvolvedEntities[evt.Entity] = struct{}{}
    case eventsys.EvenshroudResistActivateEvent:
        uniqueInvolvedEntities[evt.Entity] = struct{}{}
    case eventsys.EvenshroudResistDeactivateEvent:
        uniqueInvolvedEntities[evt.Entity] = struct{}{}
    // General game events that items might react to
    case eventsys.AttackFiredEvent:
        uniqueInvolvedEntities[evt.Source] = struct{}{} // Item on attacker
    case eventsys.AttackLandedEvent:
        uniqueInvolvedEntities[evt.Source] = struct{}{} // Item on attacker
        // If items on the target can react to an attack landing (e.g., Bramble Vest), add evt.Target:
        // uniqueInvolvedEntities[evt.Target] = struct{}{}
    case eventsys.DamageAppliedEvent:
        uniqueInvolvedEntities[evt.Source] = struct{}{} // For items on the dealer of damage
        uniqueInvolvedEntities[evt.Target] = struct{}{} // For items on the receiver of damage
    case eventsys.SpellLandedEvent:
        uniqueInvolvedEntities[evt.Source] = struct{}{} // For items on the caster
    case eventsys.KillEvent:
        uniqueInvolvedEntities[evt.Killer] = struct{}{} // For items on the killer
    case eventsys.AssistEvent:
        uniqueInvolvedEntities[evt.Assistor] = struct{}{} // For items on the assister
    default:
        // Log for unhandled event types, potentially at a debug level if noisy
        // log.Printf("ItemManager: Unhandled event type: %T for item processing", event)
        return // No specific entities identified for item processing
    }

    if len(uniqueInvolvedEntities) == 0 {
        return
    }

    // This map ensures a specific item handler on a specific entity is called at most once per event.
    processedHandlersThisEvent := make(map[entity.Entity]map[string]struct{}) // entity -> itemApiName -> processed

    for entity := range uniqueInvolvedEntities {
        equipment, ok := im.world.GetEquipment(entity)
        if !ok {
            continue
        }

        if processedHandlersThisEvent[entity] == nil {
            processedHandlersThisEvent[entity] = make(map[string]struct{})
        }

        for _, itemSlot := range equipment.Items {
            if itemSlot == nil {
                continue
            }
            itemApiName := itemSlot.ApiName

            // Check if this item type on this entity has already been processed for this event
            if _, alreadyProcessed := processedHandlersThisEvent[entity][itemApiName]; alreadyProcessed {
                continue
            }

            // itemData := data.GetItemByApiName(itemApiName) // Get full item data
            // if itemData == nil {
            // 	log.Printf("ItemManager: Could not find item data for API name %s on entity %d", itemApiName, entity)
            // 	continue
            // }

            if handler, exists := itemsys.GetItemHandler(itemApiName); exists {
                handler.ProcessEvent(event, entity, im.world, im.eventBus)
                processedHandlersThisEvent[entity][itemApiName] = struct{}{} // Mark as processed for this event
            }
        }
    }
}
