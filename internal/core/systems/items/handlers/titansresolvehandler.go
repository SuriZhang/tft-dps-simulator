package itemhandlers

import (
	"log"
	"tft-dps-simulator/internal/core/data"
	"tft-dps-simulator/internal/core/ecs"

	eventsys "tft-dps-simulator/internal/core/systems/events"
	itemsys "tft-dps-simulator/internal/core/systems/items"
)

type TitansResolveHandler struct{}

func init() {
    itemsys.RegisterItemHandler(data.TFT_Item_TitansResolve, &TitansResolveHandler{})
}

// OnEquip implements itemsys.ItemHandler.
// For Titan's Resolve, there's no immediate timed event to schedule on equip.
// Stacks are gained through AttackLandedEvent or DamageAppliedEvent.
func (h *TitansResolveHandler) OnEquip(entity ecs.Entity, world *ecs.World, eventBus eventsys.EventBus) {
    log.Printf("TitansResolveHandler: OnEquip for entity %d. No initial events to enqueue.", entity)

    effect, exists := world.GetTitansResolveEffect(entity)
    if exists {
        effect.ResetStacks() // Reset stacks on equip
        log.Printf("TitansResolveHandler: Reset Titan's Resolve stacks for entity %d on equip.", entity)
    } else {
        log.Printf("TitansResolveHandler: WARNING - TitansResolveEffect component not found for entity %d on equip.", entity)
    }
}

// ProcessEvent implements itemsys.ItemHandler.
// This function will handle AttackLandedEvent and DamageAppliedEvent to grant stacks.
func (h *TitansResolveHandler) ProcessEvent(event interface{}, entity ecs.Entity, world *ecs.World, eventBus eventsys.EventBus) {
    // Generic check for entity validity
    health, healthOk := world.GetHealth(entity)
    equipment, equipOk := world.GetEquipment(entity)
    if !healthOk || !equipOk || health.GetCurrentHP() <= 0 || !equipment.HasItem(data.TFT_Item_TitansResolve) {
        log.Printf("TitansResolveHandler: Entity %d no longer valid or item removed. Skipping event processing.", entity)
        return
    }

    _, effectOk := world.GetTitansResolveEffect(entity)
    if !effectOk {
        log.Printf("TitansResolveHandler: ERROR - TitansResolveEffect component not found for entity %d during event processing.", entity)
        return
    }

    switch e := event.(type) {
    case eventsys.AttackLandedEvent:
        if e.Source != entity {
            return // Event not for this entity as attacker
        }
        h.handleTitansTrigger(entity, e.Timestamp, world, eventBus)

    case eventsys.DamageAppliedEvent:
        if e.Target != entity { // Titan's stacks when the *wearer* takes damage
            return
        }
       	h.handleTitansTrigger(entity, e.Timestamp, world, eventBus)

    default:
        // Not an event this handler is interested in
        return
    }
}

func (h *TitansResolveHandler) handleTitansTrigger(entity ecs.Entity, evtTimestamp float64, world *ecs.World, eventBus eventsys.EventBus) {
    equipment, ok := world.GetEquipment(entity)
    if (!ok || !equipment.HasItem(data.TFT_Item_TitansResolve)) {
        return // Entity doesn't have the item
    }

    effect, ok := world.GetTitansResolveEffect(entity)
    if !ok {
        log.Printf("TitansResolveHandler (handleTitansTrigger): Warning - Entity %d has Titan's Resolve item but no TitansResolveEffect component.", entity)
        return
    }

    // Try to add a stack
    stackAdded, reachedMax := effect.IncrementStacks()

    if stackAdded {
        log.Printf("TitansResolveHandler (handleTitansTrigger): Entity %d Titan's Resolve, Stack added (%d/%d).", entity, effect.GetCurrentStacks(), effect.GetMaxStacks())

        // Apply delta AD bonus
        if attackComp, ok := world.GetAttack(entity); ok {
            deltaAD := effect.GetADPerStack() // AD gained this stack
            attackComp.AddBonusPercentAD(deltaAD)
            log.Printf("  Applied delta AD: +%.2f%%. Total Bonus AD: %.2f%%", deltaAD*100, attackComp.GetBonusPercentAD()*100)
        }

        // Apply delta AP bonus
        if spellComp, ok := world.GetSpell(entity); ok {
            deltaAP := effect.GetAPPerStack() // AP gained this stack
            spellComp.AddBonusAP(deltaAP)
             log.Printf("  Applied delta AP: +%.1f. Total Bonus AP: %.1f", deltaAP, spellComp.GetBonusAP())
        }

        // Apply resists bonus ONLY if max stacks were reached *this time*
        if reachedMax && !effect.IsBonusResistsApplied() {
            if healthComp, ok := world.GetHealth(entity); ok {
                bonusMR := effect.GetBonusMRAtMax() // MR gained at max stacks
				bonusArmor := effect.GetBonusArmorAtMax() // Armor gained at max stacks
                healthComp.AddBonusArmor(bonusArmor)
                healthComp.AddBonusMR(bonusMR)
                log.Printf("  Reached max stacks! Applied bonus resists: +%.0f Armor, +%.0f MR.", bonusArmor, bonusMR)
            }
        }
        recalcEvent := eventsys.RecalculateStatsEvent{Entity: entity, Timestamp: evtTimestamp}
        eventBus.Enqueue(recalcEvent, evtTimestamp)
        log.Printf("TitansResolveHandler (handleTitansTrigger): Enqueued RecalculateStatsEvent for entity %d at %.3fs", entity, evtTimestamp)
    }

}
