package itemsys

import (
	"log"

	"github.com/suriz/tft-dps-simulator/data"
	"github.com/suriz/tft-dps-simulator/ecs"
	eventsys "github.com/suriz/tft-dps-simulator/systems/events"
)

// DynamicEventItemSystem handles items triggered by game events.
type DynamicEventItemSystem struct {
    world    *ecs.World
    eventBus eventsys.EventBus // Keep the interface type
}

// NewDynamicEventItemSystem creates a new system instance and registers it as a handler.
func NewDynamicEventItemSystem(world *ecs.World, bus eventsys.EventBus) *DynamicEventItemSystem {
    s := &DynamicEventItemSystem{
        world:    world,
        eventBus: bus, // Store the bus if needed later, otherwise can be removed
    }
    return s
}

// HandleEvent processes incoming events.
func (s *DynamicEventItemSystem) HandleEvent(event interface{}) {
    switch evt := event.(type) {
    // Use AttackLandedEvent
    case eventsys.AttackLandedEvent: // <<< CORRECTED EVENT TYPE
        // Triggered when the entity *attacks* and hits
        s.handleTitansTrigger(evt.Source)
    case eventsys.DamageAppliedEvent:
        // Triggered when the entity *takes damage*
        s.handleTitansTrigger(evt.Target)
        // Add other event cases here if needed for other items
    }
}

// handleTitansTrigger checks if an entity has Titan's and processes a stack gain.
func (s *DynamicEventItemSystem) handleTitansTrigger(entity ecs.Entity) {
    equipment, ok := s.world.GetEquipment(entity)
    if (!ok || !equipment.HasItem(data.TFT_Item_TitansResolve)) {
        return // Entity doesn't have the item
    }

    effect, ok := s.world.GetTitansResolveEffect(entity)
    if !ok {
        log.Printf("Warning: Entity %d has Titan's Resolve item but no TitansResolveEffect component.", entity)
        return
    }

    // Try to add a stack
    stackAdded, reachedMax := effect.AddStack()

    if stackAdded {
        log.Printf("Entity %d Titan's Resolve: Stack added (%d/%d).", entity, effect.GetCurrentStacks(), effect.GetMaxStacks())

        // Apply delta AD bonus
        if attackComp, ok := s.world.GetAttack(entity); ok {
            deltaAD := effect.GetADPerStack() // AD gained this stack
            attackComp.AddBonusPercentAD(deltaAD)
            log.Printf("  Applied delta AD: +%.2f%%. Total Bonus AD: %.2f%%", deltaAD*100, attackComp.GetBonusPercentAD()*100)
        }

        // Apply delta AP bonus
        if spellComp, ok := s.world.GetSpell(entity); ok {
            deltaAP := effect.GetAPPerStack() // AP gained this stack
            spellComp.AddBonusAP(deltaAP)
             log.Printf("  Applied delta AP: +%.1f. Total Bonus AP: %.1f", deltaAP, spellComp.GetBonusAP())
        }

        // Apply resists bonus ONLY if max stacks were reached *this time*
        if reachedMax {
            if healthComp, ok := s.world.GetHealth(entity); ok {
                bonusMR := effect.GetBonusMRAtMax() // MR gained at max stacks
				bonusArmor := effect.GetBonusArmorAtMax() // Armor gained at max stacks
                healthComp.AddBonusArmor(bonusArmor)
                healthComp.AddBonusMR(bonusMR)
                log.Printf("  Reached max stacks! Applied bonus resists: +%.0f Armor, +%.0f MR.", bonusArmor, bonusMR)
            }
        }
    }
}

// Update function - currently empty as logic is event-driven.
func (s *DynamicEventItemSystem) Update(dt float64) {
    // No periodic updates needed for Titan's Resolve itself
}