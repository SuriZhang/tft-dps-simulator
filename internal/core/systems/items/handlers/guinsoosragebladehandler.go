package itemhandlers

import (
	"log"

	"tft-dps-simulator/internal/core/data"
	"tft-dps-simulator/internal/core/ecs"
	"tft-dps-simulator/internal/core/entity"
	eventsys "tft-dps-simulator/internal/core/systems/events"
	itemsys "tft-dps-simulator/internal/core/systems/items"
)

type GuinsoosRagebladeHandler struct{}

func init() {
	itemsys.RegisterItemHandler(data.TFT_Item_GuinsoosRageblade, &GuinsoosRagebladeHandler{})
}

// OnEquip implements itemsys.ItemHandler.
func (h *GuinsoosRagebladeHandler) OnEquip(entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
	log.Printf("GuinsoosRagebladeHandler: OnEquip for entity %d", entity)

	effect, exists := world.GetGuinsoosRagebladeEffect(entity)
	if exists {
		effect.ResetEffects() // Reset stacks on equip
		log.Printf("GuinsoosRagebladeHandler: Reset Guinsoo's Rageblade stacks for entity %d on equip.", entity)
	} else {
		log.Printf("GuinsoosRagebladeHandler: WARNING - GuinsoosRagebladeEffect component not found for entity %d on equip.", entity)
	}

	if exists && effect.GetInterval() > 0 {
		effect.ResetEffects()
		firstTickTime := effect.GetInterval() // First tick happens after interval
		tickEvent := eventsys.GuinsoosRagebladeTickEvent{Entity: entity, Timestamp: firstTickTime}
		eventBus.Enqueue(tickEvent, firstTickTime)
		log.Printf("  Enqueued initial GuinsoosRagebladeTickEvent for entity %d at t=%.3fs", entity, firstTickTime)
	}
}

func (h *GuinsoosRagebladeHandler) ProcessEvent(event interface{}, entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
	ragebladeTickEvent, ok := event.(eventsys.GuinsoosRagebladeTickEvent)
	if !ok || ragebladeTickEvent.Entity != entity {
		return // Not a GuinsoosRagebladeTickEvent for this entity
	}
	currentTime := ragebladeTickEvent.Timestamp
	log.Printf("GuinsoosRagebladeHandler: Processing GuinsoosRagebladeTickEvent for entity %d at t=%.3fs", entity, currentTime)

	equipment, equipOk := world.GetEquipment(entity)
	effect, effectOk := world.GetGuinsoosRagebladeEffect(entity)
	attack, attackOk := world.GetAttack(entity)
	health, healthOk := world.GetHealth(entity)
	if !equipOk || !effectOk || !attackOk || !healthOk || health.GetCurrentHP() <= 0 || !equipment.HasItem(data.TFT_Item_GuinsoosRageblade) {
		log.Printf("GuinsoosRagebladeHandler (Tick): Entity %d no longer valid or item removed at %.3fs. Stopping ticks.", entity, currentTime)
		return
	}

	guinsoosCount := equipment.GetItemCount(data.TFT_Item_GuinsoosRageblade)
	attackSpeedGain := effect.GetAttackSpeedPerStack() * float64(guinsoosCount)
	attack.AddBonusPercentAttackSpeed(attackSpeedGain) // This directly modifies the Attack component
	effect.IncrementStacks()

	log.Printf("GuinsoosRagebladeHandler (Tick): Entity %d gained %.1f%% Attack Speed (Stacks: %d, Count: %d). Total Bonus Attack Speed: %.1f%%",
		entity, attackSpeedGain*100, effect.GetStacks(), guinsoosCount, attack.GetBonusPercentAttackSpeed()*100)

	// Enqueue event to recalculate the attack speed
	recalcEvent := eventsys.RecalculateStatsEvent{Entity: entity, Timestamp: currentTime}
	eventBus.Enqueue(recalcEvent, currentTime)
	
	// Schedule the next tick
	nextTickTime := currentTime + effect.GetInterval()
	tickEvent := eventsys.GuinsoosRagebladeTickEvent{Entity: entity, Timestamp: nextTickTime}
	eventBus.Enqueue(tickEvent, nextTickTime)
}