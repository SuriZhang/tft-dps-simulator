package itemhandlers

import (
	"log"

	"tft-dps-simulator/internal/core/data"
	"tft-dps-simulator/internal/core/ecs"
	"tft-dps-simulator/internal/core/entity"
	eventsys "tft-dps-simulator/internal/core/systems/events"
	itemsys "tft-dps-simulator/internal/core/systems/items"
)

type NashorsToothHandler struct{}

func init() {
	itemsys.RegisterItemHandler(data.TFT_Item_NashorsTooth, &NashorsToothHandler{})
}

// OnEquip implements itemsys.ItemHandler.
func (h *NashorsToothHandler) OnEquip(entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
    log.Printf("NashorsToothHandler: OnEquip for entity %d", entity)
    
    effect, exists := world.GetNashorsToothEffect(entity)
    if exists {
        effect.ResetEffects()
        log.Printf("NashorsToothHandler: Reset Nashor's Tooth effect for entity %d on equip.", entity)
    } else {
        log.Printf("NashorsToothHandler: WARNING - NashorsToothEffect component not found for entity %d on equip.", entity)
    }
}

// // OnUnequip implements itemsys.ItemHandler.
// func (h *NashorsToothHandler) OnUnequip(entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
//     log.Printf("NashorsToothHandler: OnUnequip for entity %d", entity)
    
//     // Remove any active buff when item is unequipped
//     if effect, exists := world.GetNashorsToothEffect(entity); exists {
//         wasActive := effect.IsActive(0) // We don't have timestamp here, use 0
//         if wasActive {
//             // Remove the attack speed bonus
//             if attack, ok := world.GetAttack(entity); ok {
//                 attack.AddBonusPercentAttackSpeed(-effect.GetBonusAS())
//                 log.Printf("NashorsToothHandler: Removed %.2f%% attack speed from entity %d on unequip", 
//                     effect.GetBonusAS()*100, entity)
//             }
//         }
//         effect.ResetEffects()
//     }
// }

// ProcessEvent implements itemsys.ItemHandler.
func (h *NashorsToothHandler) ProcessEvent(event interface{}, entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
    switch evt := event.(type) {
    case eventsys.SpellLandedEvent:
        if evt.Source == entity {
            h.handleSpellLanded(evt, entity, world, eventBus)
        }
    case eventsys.NashorsToothDeactivateEvent:
        if evt.Entity == entity {
            h.handleBuffDeactivation(evt, entity, world, eventBus)
        }
    }
}

// handleSpellLanded activates the attack speed buff when a spell is cast
func (h *NashorsToothHandler) handleSpellLanded(evt eventsys.SpellLandedEvent, entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
    // Verify entity still has Nashor's Tooth
    equipment, okEq := world.GetEquipment(entity)
    if !okEq || equipment.GetItemCount(data.TFT_Item_NashorsTooth) == 0 {
        log.Printf("NashorsToothHandler (handleSpellLanded): Entity %d no longer has Nashor's Tooth at %.3fs.", entity, evt.Timestamp)
        return
    }

    effect, okNT := world.GetNashorsToothEffect(entity)
    if !okNT {
        log.Printf("NashorsToothHandler (handleSpellLanded): Entity %d missing NashorsToothEffect component at %.3fs.", entity, evt.Timestamp)
        return
    }

    attack, okAttack := world.GetAttack(entity)
    if !okAttack {
        log.Printf("NashorsToothHandler (handleSpellLanded): Entity %d missing Attack component at %.3fs.", entity, evt.Timestamp)
        return
    }

    // Check if already active - this handles timer resets
    wasActive := effect.IsActive(evt.Timestamp)

    // Increment sequence to cancel any pending deactivation
    effect.IncrementCurrentSequence()

    // Always activate/reset the buff
    effect.ActivateBuff(evt.Timestamp)

    // If not already active, add the attack speed bonus
    if !wasActive {
        nashorsCount := equipment.GetItemCount(data.TFT_Item_NashorsTooth)
        asGain := effect.GetBonusAS() * float64(nashorsCount)
        attack.AddBonusPercentAttackSpeed(asGain)
        
        log.Printf("NashorsToothHandler (handleSpellLanded): Entity %d activated attack speed buff at %.3fs (%.1f%% for %.1fs)",
            entity, evt.Timestamp, asGain*100, effect.GetDuration())
    } else {
        log.Printf("NashorsToothHandler (handleSpellLanded): Entity %d refreshed attack speed buff at %.3fs (%.1f%% for %.1fs)",
            entity, evt.Timestamp, effect.GetBonusAS()*100, effect.GetDuration())
    }

    // Schedule deactivation event
    deactivateTime := evt.Timestamp + effect.GetDuration()
    deactivateEvent := eventsys.NashorsToothDeactivateEvent{
        Entity:    entity,
        Timestamp: deactivateTime,
        Sequence:  effect.GetCurrentSequence(),
    }
    eventBus.Enqueue(deactivateEvent, deactivateTime)

    // Enqueue RecalculateStatsEvent to update champion stats
    recalcEvent := eventsys.RecalculateStatsEvent{Entity: entity, Timestamp: evt.Timestamp}
    eventBus.Enqueue(recalcEvent, evt.Timestamp)
}

// handleBuffDeactivation removes the attack speed buff when it expires
func (h *NashorsToothHandler) handleBuffDeactivation(evt eventsys.NashorsToothDeactivateEvent, entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
    effect, exists := world.GetNashorsToothEffect(entity)
    if !exists {
        log.Printf("NashorsToothHandler (handleBuffDeactivation): Entity %d missing NashorsToothEffect component at %.3fs.", entity, evt.Timestamp)
        return
    }

    // Check if this deactivation is for the current activation sequence
    if evt.Sequence != effect.GetCurrentSequence() {
        log.Printf("NashorsToothHandler (handleBuffDeactivation): Entity %d ignoring outdated deactivation event at %.3fs (seq %d vs %d).",
            entity, evt.Timestamp, evt.Sequence, effect.GetCurrentSequence())
        return
    }

    // Only deactivate if still active and not refreshed
    if effect.IsActive(evt.Timestamp) {
        attack, okAttack := world.GetAttack(entity)
        if !okAttack {
            log.Printf("NashorsToothHandler (handleBuffDeactivation): Entity %d missing Attack component at %.3fs.", entity, evt.Timestamp)
            return
        }

        // Remove the attack speed bonus
        equipment, _ := world.GetEquipment(entity)
        nashorsCount := equipment.GetItemCount(data.TFT_Item_NashorsTooth)
        asLoss := effect.GetBonusAS() * float64(nashorsCount)
        attack.AddBonusPercentAttackSpeed(-asLoss)
        
        effect.DeactivateBuff()

        log.Printf("NashorsToothHandler (handleBuffDeactivation): Entity %d deactivated attack speed buff at %.3fs (removed %.1f%% AS)",
            entity, evt.Timestamp, asLoss*100)

        // Enqueue RecalculateStatsEvent to update champion stats
        recalcEvent := eventsys.RecalculateStatsEvent{Entity: entity, Timestamp: evt.Timestamp}
        eventBus.Enqueue(recalcEvent, evt.Timestamp)
    }
}