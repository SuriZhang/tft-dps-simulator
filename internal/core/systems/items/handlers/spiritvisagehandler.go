package itemhandlers

import (
	"log"
	"math"

	"tft-dps-simulator/internal/core/data"
	"tft-dps-simulator/internal/core/ecs"
	"tft-dps-simulator/internal/core/entity"
	eventsys "tft-dps-simulator/internal/core/systems/events"
	itemsys "tft-dps-simulator/internal/core/systems/items"
)

type SpiritVisageHandler struct{}

func init() {
	itemsys.RegisterItemHandler(data.TFT_Item_SpiritVisage, &SpiritVisageHandler{})
}

func (h *SpiritVisageHandler) OnEquip(entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
	effect, okEffect := world.GetSpiritVisageEffect(entity)
	if !okEffect {
		log.Printf("SpiritVisageHandler: SpiritVisageEffect component not found for entity %d on equip.", entity)
		return
	}

	firstTickTime := effect.GetTickInterval() // First tick happens after interval
	healTickEvent := eventsys.SpiritVisageHealTickEvent{Entity: entity, Timestamp: firstTickTime}
	eventBus.Enqueue(healTickEvent, firstTickTime)
	log.Printf("SpiritVisageHandler: OnEquip for entity %d. Scheduled first heal tick at %.3fs.", entity, healTickEvent.Timestamp)
}

func (h *SpiritVisageHandler) ProcessEvent(event interface{}, entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
	spiritVisageHealTickEvent, ok := event.(eventsys.SpiritVisageHealTickEvent)
	if !ok || spiritVisageHealTickEvent.Entity != entity {
		return 
	}

	currentTime := spiritVisageHealTickEvent.Timestamp
	log.Printf("SpiritVisageHandler (Tick): Processing SpiritVisageHealTickEvent for entity %d at t=%.3fs", entity, currentTime)

	equipment, equipOk := world.GetEquipment(entity)
	effect, effectOk := world.GetSpiritVisageEffect(entity)
	health, healthOk := world.GetHealth(entity)

	if !equipOk || !effectOk || !healthOk || health.GetCurrentHP() <= 0 || !equipment.HasItem(data.TFT_Item_SpiritVisage) {
		log.Printf("SpiritVisageHandler (Tick): Entity %d no longer valid or item removed at %.3fs. Stopping ticks.", entity, currentTime)
		return
	}

	spiritVisageCount := equipment.GetItemCount(data.TFT_Item_SpiritVisage)
	healthGainRate := effect.GetMissingHealthHealRate()
	missingHealth := health.GetFinalMaxHP() - health.GetCurrentHP()

	if missingHealth > 0 {
		healAmountPerCount := math.Max(missingHealth*healthGainRate, effect.GetMaxHeal())
		previousHP := health.GetCurrentHP()
		health.Heal(healAmountPerCount * float64(spiritVisageCount))
		healedAmount := health.GetCurrentHP() - previousHP

		if healedAmount > 0.01 { // Log only if a meaningful amount was healed
			log.Printf("SpiritVisageHandler (Tick): Entity %d healed for %.2f (%.1f%% of missing HP %.2f). HP: %.2f -> %.2f / %.2f. Timestamp: %.3fs",
				entity, healedAmount, healthGainRate*100, missingHealth, previousHP, health.GetCurrentHP(), health.GetFinalMaxHP(), currentTime)
		}

		// Enqueue event to recalculate stats
		recalcEvent := eventsys.RecalculateStatsEvent{Entity: entity, Timestamp: currentTime}
		eventBus.Enqueue(recalcEvent, currentTime)
	}

	nextTickTime := currentTime + effect.GetTickInterval()
	nextHealEvent := eventsys.SpiritVisageHealTickEvent{
		Entity:    entity,
		Timestamp: nextTickTime,
	}
	eventBus.Enqueue(nextHealEvent, nextTickTime)
	log.Printf("SpiritVisageHandler: Scheduled next heal tick for entity %d at %.3fs.", entity, nextTickTime)
}