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

type VoidStaffHandler struct{}

func init() {
	itemsys.RegisterItemHandler(data.TFT_Item_VoidStaff, &VoidStaffHandler{})
}

// OnEquip implements itemsys.ItemHandler.
func (h *VoidStaffHandler) OnEquip(entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
    log.Printf("VoidStaffHandler: OnEquip for entity %d", entity)
    
    // Void Staff doesn't need initial tick events like Archangel's Staff
    // It's a reactive item that triggers on attacks and ability damage
}

// OnUnequip implements itemsys.ItemHandler.
func (h *VoidStaffHandler) OnUnequip(entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
    log.Printf("VoidStaffHandler: OnUnequip for entity %d", entity)
    
    // No cleanup needed as Void Staff doesn't maintain persistent state
}

// ProcessEvent implements itemsys.ItemHandler.
func (h *VoidStaffHandler) ProcessEvent(event interface{}, entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
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

// handleAttackLanded applies shred when an attack lands
func (h *VoidStaffHandler) handleAttackLanded(evt eventsys.AttackLandedEvent, entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
    // Verify entity still has Void Staff
    equipment, okEq := world.GetEquipment(entity)
    if !okEq || equipment.GetItemCount(data.TFT_Item_VoidStaff) == 0 {
        log.Printf("VoidStaffHandler (handleAttackLanded): Entity %d no longer has Void Staff at %.3fs.", entity, evt.Timestamp)
        return
    }

    effect, okVS := world.GetVoidStaffEffect(entity)
    if !okVS {
        log.Printf("VoidStaffHandler (handleAttackLanded): Entity %d missing VoidStaffEffect component at %.3fs.", entity, evt.Timestamp)
        return
    }

    h.applyShred(evt.Target, evt.Source, effect, evt.Timestamp, eventBus)
}

// handleSpellLanded applies shred when ability damage is dealt
func (h *VoidStaffHandler) handleSpellLanded(evt eventsys.SpellLandedEvent, entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
    // Verify entity still has Void Staff
    equipment, okEq := world.GetEquipment(entity)
    if !okEq || equipment.GetItemCount(data.TFT_Item_VoidStaff) == 0 {
        log.Printf("VoidStaffHandler (handleSpellLanded): Entity %d no longer has Void Staff at %.3fs.", entity, evt.Timestamp)
        return
    }

    effect, okVS := world.GetVoidStaffEffect(entity)
    if !okVS {
        log.Printf("VoidStaffHandler (handleSpellLanded): Entity %d missing VoidStaffEffect component at %.3fs.", entity, evt.Timestamp)
        return
    }

    h.applyShred(evt.Target, evt.Source, effect, evt.Timestamp, eventBus)
}

// applyShred applies the shred debuff to the target
func (h *VoidStaffHandler) applyShred(target, source entity.Entity, effect *items.VoidStaffEffect, timestamp float64, eventBus eventsys.EventBus) {
    mrShred := effect.GetMRShred()
    duration := effect.GetDuration()

    // Create and enqueue the apply debuff event
    shredEvent := eventsys.ApplyDebuffEvent{
        Target:     target,
        Source:     source,
        DebuffType: debuffs.Shred,
        Value:      mrShred,
        Duration:   duration,
        Timestamp:  timestamp,
        SourceType: "Item",
        SourceId:   data.TFT_Item_VoidStaff,
    }

    eventBus.Enqueue(shredEvent, timestamp)
    
    log.Printf("VoidStaffHandler (applyShred): Entity %d applied %.1f%% MR shred to entity %d for %.1fs at %.3fs",
        source, mrShred, target, duration, timestamp)
}
