package itemhandlers

import (
	"log"
	"tft-dps-simulator/internal/core/data"
	"tft-dps-simulator/internal/core/ecs"

	eventsys "tft-dps-simulator/internal/core/systems/events"
	itemsys "tft-dps-simulator/internal/core/systems/items"
)

type KrankensFuryHandler struct{}

func init() {
	itemsys.RegisterItemHandler(data.TFT_Item_KrakensFury, &KrankensFuryHandler{})
}

func (h *KrankensFuryHandler) OnEquip(entity ecs.Entity, world *ecs.World, eventBus eventsys.EventBus) {
	log.Printf("KrankensFuryHandler: OnEquip for entity %d. No initial events to enqueue.", entity)

	effect, exists := world.GetKrakensFuryEffect(entity)
	if exists {
		effect.ResetEffects() // Reset stacks on equip
		log.Printf("KrankensFuryHandler: Reset Krankens Fury stacks for entity %d on equip.", entity)
	} else {
		log.Printf("KrankensFuryHandler: WARNING - KrankensFuryEffect component not found for entity %d on equip.", entity)
	}
}

func (h *KrankensFuryHandler) ProcessEvent(event interface{}, entity ecs.Entity, world *ecs.World, eventBus eventsys.EventBus) {

	// Generic check for entity validity
	health, healthOk := world.GetHealth(entity)
	equipment, equipOk := world.GetEquipment(entity)
	if !healthOk || !equipOk || health.GetCurrentHP() <= 0 || !equipment.HasItem(data.TFT_Item_KrakensFury) {
		log.Printf("KrankensFuryHandler: Entity %d no longer valid or item removed. Skipping event processing.", entity)
		return
	}

	_, effectOk := world.GetKrakensFuryEffect(entity)
	if !effectOk {
		log.Printf("KrankensFuryHandler: ERROR - KrankensFuryEffect component not found for entity %d during event processing.", entity)
		return
	}

	switch e := event.(type) {
	case eventsys.AttackLandedEvent:
		if e.Source != entity {
			return // Event not for this entity as attacker
		}
		h.handleKrankensTrigger(entity, e.Timestamp, world, eventBus)

	default:
		return
	}

}

func (h *KrankensFuryHandler) handleKrankensTrigger(entity ecs.Entity, evtTimestamp float64, world *ecs.World, eventBus eventsys.EventBus) {
	equipment, ok := world.GetEquipment(entity)
	if !ok || !equipment.HasItem(data.TFT_Item_KrakensFury) {
		return // Entity doesn't have the item
	}

	effect, ok := world.GetKrakensFuryEffect(entity)
	if !ok {
		log.Printf("KrankensFuryHandler: ERROR - KrankensFuryEffect component not found for entity %d during event processing.", entity)
		return
	}

	effect.IncrementStacks()
	log.Printf("KrankensFuryHandler: Incremented stacks for entity %d to %d.", entity, effect.GetCurrentStacks())

	attackComp, ok := world.GetAttack(entity)
	if !ok {
		log.Printf("KrakensFuryHandler: Entity %d with Kraken's Fury has no Attack component.", entity)
		return
	}

	attackComp.AddBonusAD(effect.GetADPerStack())
	log.Printf("KrankensFuryHandler: Added delta %f bonus AD for entity %d. Total stacks: %d.", effect.GetADPerStack(), entity, effect.GetCurrentStacks())

	// Enqueue event to notify other systems of the change
	recalcEvent := eventsys.RecalculateStatsEvent{Entity: entity, Timestamp: evtTimestamp}
	eventBus.Enqueue(recalcEvent, evtTimestamp)
	log.Printf("KrankensFuryHandler: Enqueued RecalculateStatsEvent for entity %d.", entity)
}
