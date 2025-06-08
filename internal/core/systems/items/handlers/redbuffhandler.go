package itemhandlers

import (
	"log"

    "tft-dps-simulator/internal/core/components/debuffs"
	"tft-dps-simulator/internal/core/components/items"
    "tft-dps-simulator/internal/core/data"
    "tft-dps-simulator/internal/core/ecs"
    "tft-dps-simulator/internal/core/entity"
    eventsys "tft-dps-simulator/internal/core/systems/events"
    itemsys "tft-dps-simulator/internal/core/systems/items"
)

type RedBuffHandler struct{}

// OnEquip implements itemsys.ItemHandler.
func (h *RedBuffHandler) OnEquip(entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
	log.Printf("RedBuffHandler: OnEquip for entity %d", entity)

	// Red Buff doesn't need initial tick events
	// It's a reactive item that triggers on attacks and abilities
}

// OnUnequip implements itemsys.ItemHandler.
func (h *RedBuffHandler) OnUnequip(entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
	log.Printf("RedBuffHandler: OnUnequip for entity %d", entity)

	// No cleanup needed as Red Buff doesn't maintain persistent state
}

// ProcessEvent implements itemsys.ItemHandler.
func (h *RedBuffHandler) ProcessEvent(event interface{}, entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
	switch evt := event.(type) {
	case eventsys.AttackLandedEvent:
		if evt.Source == entity {
			h.handleAttackLanded(evt, entity, world, eventBus)
		}
	case eventsys.SpellLandedEvent:
		if evt.Source == entity {
			h.handleSpellLanded(evt, entity, world, eventBus)
		}
	}
}

// handleAttackLanded applies burn and wound when an attack lands
func (h *RedBuffHandler) handleAttackLanded(evt eventsys.AttackLandedEvent, entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
	// Verify entity still has Red Buff
	equipment, okEq := world.GetEquipment(entity)
	if !okEq || equipment.GetItemCount(data.TFT_Item_RedBuff) == 0 {
		log.Printf("RedBuffHandler (handleAttackLanded): Entity %d no longer has Red Buff at %.3fs.", entity, evt.Timestamp)
		return
	}

	effect, okRB := world.GetRedBuffEffect(entity)
	if !okRB {
		log.Printf("RedBuffHandler (handleAttackLanded): Entity %d missing RedBuffEffect component at %.3fs.", entity, evt.Timestamp)
		return
	}

	h.applyBurnAndWound(evt.Target, evt.Source, effect, evt.Timestamp, eventBus)
}

// handleSpellLanded applies burn and wound when ability damage is dealt
func (h *RedBuffHandler) handleSpellLanded(evt eventsys.SpellLandedEvent, entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
	// Verify entity still has Red Buff
	equipment, okEq := world.GetEquipment(entity)
	if !okEq || equipment.GetItemCount(data.TFT_Item_RedBuff) == 0 {
		log.Printf("RedBuffHandler (handleSpellLanded): Entity %d no longer has Red Buff at %.3fs.", entity, evt.Timestamp)
		return
	}

	effect, okRB := world.GetRedBuffEffect(entity)
	if !okRB {
		log.Printf("RedBuffHandler (handleSpellLanded): Entity %d missing RedBuffEffect component at %.3fs.", entity, evt.Timestamp)
		return
	}

	h.applyBurnAndWound(evt.Target, evt.Source, effect, evt.Timestamp, eventBus)
}

// applyBurnAndWound applies both burn and wound debuffs to the target
func (h *RedBuffHandler) applyBurnAndWound(target, source entity.Entity, effect *items.RedBuffEffect, timestamp float64, eventBus eventsys.EventBus) {
	burnPercent := effect.GetBurnPercent()
	healingReductionPct := effect.GetHealingReductionPct()
	duration := effect.GetDuration()

	// Apply Burn debuff
	burnEvent := eventsys.ApplyDebuffEvent{
		Target:     target,
		Source:     source,
		DebuffType: debuffs.Burn,
		Value:      burnPercent,
		Duration:   duration,
		Timestamp:  timestamp,
		SourceType: "Item",
		SourceId:   data.TFT_Item_RedBuff,
	}
	eventBus.Enqueue(burnEvent, timestamp)

	// Apply Wound debuff
	woundEvent := eventsys.ApplyDebuffEvent{
		Target:     target,
		Source:     source,
		DebuffType: debuffs.Wound,
		Value:      healingReductionPct / 100.0, // Convert percentage to decimal
		Duration:   duration,
		Timestamp:  timestamp,
		SourceType: "Item",
		SourceId:   data.TFT_Item_RedBuff,
	}
	eventBus.Enqueue(woundEvent, timestamp)

	log.Printf("RedBuffHandler: Entity %d applied %.1f%% burn and %.1f%% wound to entity %d for %.1fs at %.3fs",
		source, burnPercent*100, healingReductionPct, target, duration, timestamp)
}

// Register the handler
func init() {
	itemsys.RegisterItemHandler(data.TFT_Item_RedBuff, &RedBuffHandler{})
}
