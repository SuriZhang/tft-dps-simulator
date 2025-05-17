package itemhandlers

import (
	"log"

	"tft-dps-simulator/internal/core/data"
	"tft-dps-simulator/internal/core/ecs"
	eventsys "tft-dps-simulator/internal/core/systems/events"
	itemsys "tft-dps-simulator/internal/core/systems/items"
)

type ArchangelsStaffHandler struct{}

func init() {
	itemsys.RegisterItemHandler(data.TFT_Item_ArchangelsStaff, &ArchangelsStaffHandler{})
}

// OnEquip implements itemsys.ItemHandler.
func (h *ArchangelsStaffHandler) OnEquip(entity ecs.Entity, world *ecs.World, eventBus eventsys.EventBus) {
	log.Printf("ArchangelsStaffHandler: OnEquip for entity %d", entity)

	effect, exists := world.GetArchangelsStaffEffect(entity)

	if exists && effect.GetInterval() > 0 {
		effect.ResetEffects()
		firstTickTime := effect.GetInterval() // First tick happens after interval
		tickEvent := eventsys.ArchangelsTickEvent{Entity: entity, Timestamp: firstTickTime}
		eventBus.Enqueue(tickEvent, firstTickTime)
		log.Printf("  Enqueued initial ArchangelsTickEvent for entity %d at t=%.3fs", entity, firstTickTime)
	} else {
        log.Printf("  ArchangelsStaffHandler: OnEquip for entity %d - effect not found or interval is 0.", entity)
    }
}

// ProcessEvent implements itemsys.ItemHandler.
func (h *ArchangelsStaffHandler) ProcessEvent(event interface{}, entity ecs.Entity, world *ecs.World, eventBus eventsys.EventBus) {
	archangelsTickEvent, ok := event.(eventsys.ArchangelsTickEvent)
	if !ok || archangelsTickEvent.Entity != entity {
		return // Not an ArchangelsTickEvent for this entity
	}

	currentTime := archangelsTickEvent.Timestamp
	log.Printf("ArchangelsStaffHandler: Processing ArchangelsTickEvent for entity %d at t=%.3fs", entity, currentTime)

	equipment, equipOk := world.GetEquipment(entity)
	effect, effectOk := world.GetArchangelsStaffEffect(entity)
	spell, spellOk := world.GetSpell(entity)
	health, healthOk := world.GetHealth(entity)

	if !equipOk || !effectOk || !spellOk || !healthOk || health.GetCurrentHP() <= 0 || !equipment.HasItem(data.TFT_Item_ArchangelsStaff) {
		log.Printf("ArchangelsStaffHandler (Tick): Entity %d no longer valid or item removed at %.3fs. Stopping ticks.", entity, currentTime)
		return
	}

	archangelsCount := equipment.GetItemCount(data.TFT_Item_ArchangelsStaff)
	apGain := effect.GetAPPerInterval() * float64(archangelsCount)
	spell.AddBonusAP(apGain) // This directly modifies the Spell component
	effect.AddStacks(1)

	log.Printf("ArchangelsStaffHandler (Tick): Entity %d gained %.1f AP (Stacks: %d, Count: %d). Total Bonus AP: %.1f",
		entity, apGain, effect.GetStacks(), archangelsCount, spell.GetBonusAP())

	// Enqueue event to recalculate stats
	recalcEvent := eventsys.RecalculateStatsEvent{Entity: entity, Timestamp: currentTime}
	eventBus.Enqueue(recalcEvent, currentTime)

	// Enqueue the next tick event
	nextTickTime := currentTime + effect.GetInterval()
	nextEvent := eventsys.ArchangelsTickEvent{Entity: entity, Timestamp: nextTickTime}
	eventBus.Enqueue(nextEvent, nextTickTime)
	log.Printf("ArchangelsStaffHandler (Tick): Enqueued next ArchangelsTickEvent for entity %d at t=%.3fs", entity, nextTickTime)
}
