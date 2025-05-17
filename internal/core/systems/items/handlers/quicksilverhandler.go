package itemhandlers

import (
	"log"
	"tft-dps-simulator/internal/core/data"
	"tft-dps-simulator/internal/core/ecs"
	eventsys "tft-dps-simulator/internal/core/systems/events"
	itemsys "tft-dps-simulator/internal/core/systems/items"
)

type QuicksilverHandler struct{}

func init() {
	itemsys.RegisterItemHandler(data.TFT_Item_Quicksilver, &QuicksilverHandler{})
}

// OnEquip implements itemsys.ItemHandler.
// This function is called when Quicksilver is equipped to an entity.
// It should schedule the initial QuicksilverProcEvent.
func (h *QuicksilverHandler) OnEquip(entity ecs.Entity, world *ecs.World, eventBus eventsys.EventBus) {
	log.Printf("QuicksilverHandler: OnEquip for entity %d", entity)

	effect, exists := world.GetQuicksilverEffect(entity)
	if exists && effect.GetProcInterval() > 0 && effect.GetSpellShieldDuration() > 0 {
		effect.ResetEffects()
		// Enqueue first proc event
		firstProcTime := effect.GetProcInterval()
		if firstProcTime <= effect.GetSpellShieldDuration() { // Only if first proc happens before expiry
			procEvent := eventsys.QuicksilverProcEvent{Entity: entity, Timestamp: firstProcTime}
			eventBus.Enqueue(procEvent, firstProcTime)
			log.Printf("  Enqueued initial QuicksilverProcEvent for entity %d at t=%.3fs", entity, firstProcTime)
		}

		// Enqueue the end event
		endTime := effect.GetSpellShieldDuration() // Duration starts at t=0
		endEvent := eventsys.QuicksilverEndEvent{Entity: entity, Timestamp: endTime}
		eventBus.Enqueue(endEvent, endTime)
		log.Printf("  Enqueued QuicksilverEndEvent for entity %d at t=%.3fs", entity, endTime)
	}
}

// ProcessEvent implements itemsys.ItemHandler.
// This function will handle QuicksilverProcEvent and QuicksilverEndEvent.
func (h *QuicksilverHandler) ProcessEvent(event interface{}, entity ecs.Entity, world *ecs.World, eventBus eventsys.EventBus) {
	// Generic check for entity validity
	health, healthOk := world.GetHealth(entity)
	equipment, equipOk := world.GetEquipment(entity)
	
	if !healthOk || !equipOk || health.GetCurrentHP() <= 0 || !equipment.HasItem(data.TFT_Item_Quicksilver) {
		log.Printf("QuicksilverHandler: Entity %d no longer valid or item removed. Skipping event processing.", entity)
		return
	}

	_, effectOk := world.GetQuicksilverEffect(entity)
	if !effectOk {
		log.Printf("QuicksilverHandler: ERROR - QuicksilverEffect component not found for entity %d during event processing.", entity)
		return
	}

	switch e := event.(type) {
	case eventsys.QuicksilverProcEvent:
		if e.Entity != entity {
			return // Event not for this entity
		}
		currentTime := e.Timestamp
		log.Printf("QuicksilverHandler: Processing QuicksilverProcEvent for entity %d at t=%.3fs", entity, currentTime)
		h.handleQuicksilverProc(e, world, eventBus)
		return

	case eventsys.QuicksilverEndEvent:
		if e.Entity != entity {
			return // Event not for this entity
		}
		currentTime := e.Timestamp
		log.Printf("QuicksilverHandler: Processing QuicksilverEndEvent for entity %d at t=%.3fs", entity, currentTime)

		h.handleQuicksilverEnd(e, world)
		return
		
	default:
		// Not an event this handler is interested in
		return
	}
}

func (h *QuicksilverHandler) handleQuicksilverProc(evt eventsys.QuicksilverProcEvent, world *ecs.World, eventBus eventsys.EventBus) {
	entity := evt.Entity
	currentTime := evt.Timestamp

	// Check if entity still exists and has the item/effect and is active
	equipment, equipOk := world.GetEquipment(entity)
	effect, effectOk := world.GetQuicksilverEffect(entity)
	attackComp, attackOk := world.GetAttack(entity)
	health, healthOk := world.GetHealth(entity) // Check if alive

	// Crucially, also check if the effect is still active (duration hasn't ended)
	if !equipOk || !effectOk || !attackOk || !healthOk || health.GetCurrentHP() <= 0 || !equipment.HasItem(data.TFT_Item_Quicksilver) || !effect.IsActive() {
		log.Printf("QuicksilverHandler (handleQuicksilverProc): Entity %d no longer valid, item removed, or effect inactive at %.3fs. Stopping procs.", entity, currentTime)
		return
	}

	// Apply effect (add bonus AS)
	quicksilverCount := equipment.GetItemCount(data.TFT_Item_Quicksilver) // Should always be 1 due to uniqueness, but check anyway
	asGain := effect.GetProcAttackSpeed() * float64(quicksilverCount)
	attackComp.AddBonusPercentAttackSpeed(asGain)
	effect.AddStacks(1) // Increment internal stack count

	log.Printf("QuicksilverHandler (handleQuicksilverProc): Entity %d gained %.2f%% AS at %.3fs (Stacks: %d). Total Bonus AS: %.2f%%",
		entity, asGain*100, currentTime, effect.GetStacks(), attackComp.GetBonusPercentAttackSpeed()*100)

	// Enqueue event to recalculate stats
	recalcEvent := eventsys.RecalculateStatsEvent{Entity: entity, Timestamp: currentTime}
	eventBus.Enqueue(recalcEvent, currentTime)

	// Enqueue the next proc event ONLY if it happens before the effect expires
	nextProcTime := currentTime + effect.GetProcInterval()
	// Use the original duration stored in the effect, don't rely on RemainingDuration which might change
	expiryTime := effect.GetSpellShieldDuration()
	if nextProcTime <= expiryTime {
		nextProcEvent := eventsys.QuicksilverProcEvent{Entity: entity, Timestamp: nextProcTime}
		eventBus.Enqueue(nextProcEvent, nextProcTime)
	} else {
		log.Printf("DynamicTimeItemSystem (QuicksilverProc): Next proc for entity %d at %.3fs would be after expiry (%.3fs). Not enqueueing.", entity, nextProcTime, expiryTime)
	}
}

// handleQuicksilverEnd marks the effect as inactive.
func (h *QuicksilverHandler) handleQuicksilverEnd(evt eventsys.QuicksilverEndEvent, world *ecs.World) {
	entity := evt.Entity
	currentTime := evt.Timestamp

	effect, effectOk := world.GetQuicksilverEffect(entity)
	if !effectOk {
		// Effect might have been removed if item was removed earlier
		return
	}

	if effect.IsActive() {
		log.Printf("DynamicTimeItemSystem (QuicksilverEnd): Entity %d Quicksilver duration ended at %.3fs. Marking inactive.", entity, currentTime)
		effect.SetIsActive(false)
		// Note: Bonus AS is NOT removed here. It persists but stops stacking.
		// If removal is desired, EquipmentManager should handle it on item removal.
	}
}
