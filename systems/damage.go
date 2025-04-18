package systems

import (
	"fmt"
	"log"

	"github.com/suriz/tft-dps-simulator/ecs"
	eventsys "github.com/suriz/tft-dps-simulator/systems/events"
)

// DamageSystem handles damage calculation and application based on events.
type DamageSystem struct {
	world    *ecs.World
	eventBus eventsys.EventBus
}

// NewDamageSystem creates a new damage system.
func NewDamageSystem(world *ecs.World, bus eventsys.EventBus) *DamageSystem {
	return &DamageSystem{
		world:    world,
		eventBus: bus,
	}
}

// HandleEvent processes incoming game events.
func (s *DamageSystem) HandleEvent(evt interface{}) {
	switch event := evt.(type) {
	case eventsys.AttackLandedEvent:
		s.onAttackLanded(event)
	case eventsys.DamageAppliedEvent:
		s.onDamageApplied(event)
		// Add cases for other event types this system might handle
	}
}

// onAttackLanded calculates final damage from an attack and enqueues DamageAppliedEvent.
func (s *DamageSystem) onAttackLanded(evt eventsys.AttackLandedEvent) {
	attacker := evt.Source
	target := evt.Target

	attackerAttack, okAtk := s.world.GetAttack(attacker)
	if !okAtk {
		log.Printf("DamageSystem Error: Attacker %d has no Attack component in onAttackLanded.\n", attacker)
		return
	}
	targetHealth, okHp := s.world.GetHealth(target)
	if !okHp {
		log.Printf("DamageSystem Error: Target %d has no Health component in onAttackLanded.\n", target)
		return
	}

	// --- Damage Calculation ---
	// 1. Base Damage (usually FinalAD for auto-attacks)
	baseDamage := evt.BaseDamage // Use damage from event

	// 2. Crit Multiplier (Expected Value)
	critChance := attackerAttack.GetFinalCritChance()
	critMultiplier := attackerAttack.GetFinalCritMultiplier()
	critMultiplierEV := (1.0 - critChance) + (critChance * critMultiplier)

	// 3. Amplification Multiplier
	ampMultiplier := 1.0 + attackerAttack.GetFinalDamageAmp()

	// 4. Armor Reduction Multiplier
	finalArmor := targetHealth.GetFinalArmor()
	armorMultiplier := 100.0 / (100.0 + finalArmor)

	// 5. Durability Multiplier (Inverse of Damage Reduction)
	durabilityMultiplier := 1.0 - targetHealth.GetFinalDurability()

	// Calculate Final Damage
	finalDamage := baseDamage * critMultiplierEV * ampMultiplier * armorMultiplier * durabilityMultiplier

	// Update Attacker's LastAttackTime (moved from AutoAttackSystem)
	attackerAttack.SetLastAttackTime(evt.Timestamp)

	// Enqueue DamageAppliedEvent
	damageAppliedEvent := eventsys.DamageAppliedEvent{
		Source:      attacker,
		Target:      target,
		FinalDamage: finalDamage,
		Timestamp:   evt.Timestamp,
	}
	s.eventBus.Enqueue(damageAppliedEvent)

	// Log the calculation result
	log.Printf("DamageSystem (onAttackLanded): Calculated %.1f final damage from %d to %d.\n", finalDamage, attacker, target) 
}

// onDamageApplied applies the final damage to the target and handles mana gain/death.
func (s *DamageSystem) onDamageApplied(evt eventsys.DamageAppliedEvent) {
	attacker := evt.Source
	target := evt.Target
	finalDamage := evt.FinalDamage

	// Get target health
	targetHealth, okHealth := s.world.GetHealth(target)
	if !okHealth {
		log.Printf("DamageSystem: Target %d missing Health component for DamageAppliedEvent.\n", target)
		return
	}

	// Get attacker info for logging
	attackerInfo, _ := s.world.GetChampionInfo(attacker)
	attackerName := fmt.Sprintf("Entity %d", attacker)
	if attackerInfo != nil {
		attackerName = attackerInfo.Name
	}

	// Get target info for logging
	targetInfo, _ := s.world.GetChampionInfo(target)
	targetName := fmt.Sprintf("Entity %d", target)
	if targetInfo != nil {
		targetName = targetInfo.Name
	}

	// Apply damage
	initialHP := targetHealth.CurrentHP
	targetHealth.CurrentHP -= finalDamage

	displayHealth := targetHealth.CurrentHP
	if displayHealth < 0 {
		displayHealth = 0
	}
	log.Printf("DamageSystem (onDamageApplied): %s attacks %s for %.1f damage (%.1f -> %.1f HP)\n", attackerName, targetName, finalDamage, initialHP, displayHealth)

	// Check for death
	if targetHealth.CurrentHP <= 0 && initialHP > 0 { // Check initialHP > 0 to prevent multiple death events
		log.Printf("DamageSystem (onDamageApplied): %s has been defeated!\n", targetName)
		deathEvent := eventsys.DeathEvent{
			Target:    target,
			Timestamp: evt.Timestamp, // Use the timestamp from the attack event
		}
		s.eventBus.Enqueue(deathEvent)
		// Note: We might need a system to handle the DeathEvent (e.g., remove entity, grant kill credit)

		// Enqueue Kill Event
        killEvent := eventsys.KillEvent{
            Killer:    attacker, // The source of the DamageAppliedEvent is the killer
            Victim:    target,
            Timestamp: evt.Timestamp, // Use the same timestamp
        }
        s.eventBus.Enqueue(killEvent)
        log.Printf("DamageSystem: %s gets kill credit for %s.\n", attackerName, targetName)
	}

	// --- Mana Gain for Attacker ---
	attackerMana, okMana := s.world.GetMana(attacker)
	if okMana {
		manaGain := 10.0 // Standard TFT mana gain per auto-attack
		attackerMana.AddCurrentMana(manaGain)

		// Clamp mana to max
		if attackerMana.GetCurrentMana() > attackerMana.GetMaxMana() {
			attackerMana.SetCurrentMana(attackerMana.GetMaxMana())
		}
		log.Printf("DamageSystem: %s gains %.1f mana (now %.1f / %.1f)\n", attackerName, manaGain, attackerMana.GetCurrentMana(), attackerMana.GetMaxMana())
	} else {
		// This might be expected for units without mana
		// log.Printf("DamageSystem: Warning: Attacker %s has no Mana component, cannot gain mana.\n", attackerName)
	}
	// --- End Mana Gain ---

	// TODO: Target also gains mana on being attacked (would likely be another handler or system reacting to DamageAppliedEvent)
}
