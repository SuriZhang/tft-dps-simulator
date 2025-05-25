package itemhandlers

import (
	"log"

	"tft-dps-simulator/internal/core/data"
	"tft-dps-simulator/internal/core/ecs"
	eventsys "tft-dps-simulator/internal/core/systems/events"
	itemsys "tft-dps-simulator/internal/core/systems/items"
)

type FlickerbladeHandler struct{}

func init() {
	itemsys.RegisterItemHandler(data.TFT_Item_Artifact_NavoriFlickerblades, &FlickerbladeHandler{})
}

func (h *FlickerbladeHandler) OnEquip(entity ecs.Entity, world *ecs.World, eventBus eventsys.EventBus) {
	log.Printf("FlickerbladeHandler: OnEquip for entity %d", entity)
	if effect, exists := world.GetFlickerbladeEffect(entity); exists {
		effect.ResetEffects()
		log.Printf("FlickerbladeHandler: Reset FlickerbladeEffect for entity %d on equip.", entity)
	} else {
		log.Printf("FlickerbladeHandler: WARNING - FlickerbladeEffect component not found for entity %d on equip.", entity)
	}
}

func (h *FlickerbladeHandler) ProcessEvent(event interface{}, entity ecs.Entity, world *ecs.World, eventBus eventsys.EventBus) {
	attackEvent, ok := event.(eventsys.AttackFiredEvent)
	if !ok || attackEvent.Source != entity {
		return
	}

	currentTime := attackEvent.Timestamp

	equipment, equipOk := world.GetEquipment(entity)
	flickerEffect, effectOk := world.GetFlickerbladeEffect(entity)
	attackComp, attackOk := world.GetAttack(entity)
	spellComp, spellOk := world.GetSpell(entity)
	health, healthOk := world.GetHealth(entity)

	if !equipOk || !effectOk || !attackOk || !spellOk || !healthOk || health.GetCurrentHP() <= 0 {
		return // Entity or components not in a valid state
	}

	if !equipment.HasItem(data.TFT_Item_Artifact_NavoriFlickerblades) {
		return // Item no longer equipped
	}

	flickerbladeCount := equipment.GetItemCount(data.TFT_Item_Artifact_NavoriFlickerblades)

	flickerEffect.IncrementAttackCounter()
	asGain := flickerEffect.GetASPerStack() * float64(flickerbladeCount)
	attackComp.AddBonusPercentAttackSpeed(asGain)
	currentAttackCountForBonus := flickerEffect.GetAttackCounter()
	log.Printf("FlickerbladeHandler: Entity %d attacked. Current attack count: %d, AS gain: %.2f%%.",
		entity, currentAttackCountForBonus, asGain*100)

	if (currentAttackCountForBonus % int(flickerEffect.GetStacksPerBonus())) == 0 {
		// Trigger AD/AP bonus application
		bonusAD := flickerEffect.GetADPerBonus() * float64(flickerbladeCount)
		bonusAP := flickerEffect.GetAPPerBonus() * float64(flickerbladeCount)

		attackComp.AddBonusAD(bonusAD)
		spellComp.AddBonusAP(bonusAP)
		log.Printf("FlickerbladeHandler: Entity %d triggered AD/AP bonus. Applied +%.2f%% AD, +%.1f AP.",
			entity, bonusAD*100, bonusAP)
	}

	recalcEvent := eventsys.RecalculateStatsEvent{Entity: entity, Timestamp: currentTime}
	eventBus.Enqueue(recalcEvent, currentTime)
}
