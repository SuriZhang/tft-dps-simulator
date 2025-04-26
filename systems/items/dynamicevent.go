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

func (s *DynamicEventItemSystem) CanHandle(evt interface{}) bool {
    // Check if the event is one we can handle
    switch evt.(type) {
    case eventsys.AttackLandedEvent, eventsys.DamageAppliedEvent:
        return true
    default:
        return false
    }
}

// HandleEvent processes incoming events.
func (s *DynamicEventItemSystem) HandleEvent(event interface{}) {
    switch evt := event.(type) {
    // Use AttackLandedEvent
    case eventsys.AttackLandedEvent: 
        // Triggered when the entity *attacks* and hits
        s.handleTitansTrigger(evt.Source)
		s.handleRagebladeTrigger(evt.Source) // Handle Rageblade trigger
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
    stackAdded, reachedMax := effect.IncrementStacks()

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
        if reachedMax && !effect.IsBonusResistsApplied() {
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

// handleRagebladeTrigger checks if an entity has Guinsoo's Rageblade and processes stack gain.
func (s *DynamicEventItemSystem) handleRagebladeTrigger(entity ecs.Entity) {
    // Check if the attacker has the GuinsoosRagebladeEffect component
    guinsosEffect, hasGuinsos := s.world.GetGuinsoosRagebladeEffect(entity)
    if !hasGuinsos {
        return // Entity doesn't have the effect component
    }

    // Check if the attacker has an Attack component
    attackComp, hasAttack := s.world.GetAttack(entity)
    if !hasAttack {
        log.Printf("Warning: Entity %d in handleRagebladeTrigger lacks Attack component", entity)
        return
    }

    // Check if the attacker has an Equipment component to count Rageblades
    equipment, hasEquipment := s.world.GetEquipment(entity)
    if !hasEquipment {
        log.Printf("Warning: Entity %d in handleRagebladeTrigger lacks Equipment component", entity)
        return // Cannot determine number of Rageblades
    }

    // --- Stacking Logic ---
    numRageblades := equipment.GetItemCount(data.TFT_Item_GuinsoosRageblade)
    if numRageblades == 0 {
        // Safeguard: Component exists but no items? Should be removed by EquipmentManager.
        log.Printf("Warning: Entity %d has GuinsoosRagebladeEffect but 0 Rageblades in equipment?", entity)
        // Optionally remove the component here if EquipmentManager failed?
        // s.world.RemoveComponent(entity, reflect.TypeOf(effects.GuinsoosRagebladeEffect{}))
        return
    }

    // Get the bonus AS *before* adding stacks for this hit
    bonusASBefore := guinsosEffect.GetCurrentBonusAS()

    // Increment stacks (potentially multiple times if multiple Rageblades)
    stacksBefore := guinsosEffect.GetCurrentStacks() // Store stacks before incrementing
    for i := 0; i < numRageblades; i++ {
        guinsosEffect.IncrementStacks() // Use the method from the effect component
		log.Printf("Entity %d Guinsoo's Rageblade: stack incremented. Current stacks: %d", entity, guinsosEffect.GetCurrentStacks())
    }

    // Get the bonus AS *after* adding stacks
    bonusASAfter := guinsosEffect.GetCurrentBonusAS()

    // Calculate the *difference* in bonus AS to apply to the Attack component
    deltaBonusAS := bonusASAfter - bonusASBefore

    if deltaBonusAS > 0 {
        // Apply the *change* in bonus AS to the Attack component's bonus field
        attackComp.AddBonusPercentAttackSpeed(deltaBonusAS)

        // Log the stacking
        log.Printf("Guinsoo's Rageblade: Entity %d (%d Rageblades) attacked. Stacks: %d -> %d. Bonus AS: +%.2f%% (Total Bonus: %.2f%%)",
            entity,
            numRageblades,
            stacksBefore,                         // Stacks before this hit
            guinsosEffect.GetCurrentStacks(),     // Stacks after this hit
            deltaBonusAS*100,                     // Delta percentage
            attackComp.GetBonusPercentAttackSpeed()*100) // Total bonus AS on Attack comp

        // IMPORTANT: Ensure StatCalculationSystem runs after this event handler
        // to update FinalAttackSpeed based on the new BonusPercentAttackSpeed.
        // Example: s.world.MarkForStatUpdate(entity)
    }
}

// Update function - currently empty as logic is event-driven.
func (s *DynamicEventItemSystem) Update(dt float64) {
    // No periodic updates needed for Titan's Resolve itself
}