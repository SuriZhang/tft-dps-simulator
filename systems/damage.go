package systems

import (
	"fmt"
	"log"
	"reflect"

	"github.com/suriz/tft-dps-simulator/components"
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
	case eventsys.SpellCastEvent:
		s.onSpellCast(event)
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
	attackerCrit, okCrit := s.world.GetCrit(attacker)
	if !okCrit {
		log.Printf("DamageSystem Error: Attacker %d has no Crit component in onAttackLanded.\n", attacker)
		return
	}
	targetHealth, okHp := s.world.GetHealth(target)
	if !okHp {
		log.Printf("DamageSystem Error: Target %d has no Health component in onAttackLanded.\n", target)
		return
	}

	// --- Damage Calculation ---
	// 1. Base Physical Damage
	basePhysical := evt.BaseDamage // Use damage from event (FinalAD)

	// 2. Crit Multiplier (Expected Value) - From Crit component
	critChance := attackerCrit.GetFinalCritChance()
	critMultiplier := attackerCrit.GetFinalCritMultiplier()
	critMultiplierEV := (1.0 - critChance) + (critChance * critMultiplier)

	// 3. Amplification Multiplier - From Attack component
	ampMultiplier := 1.0 + attackerAttack.GetFinalDamageAmp()

	// 4. Armor Reduction Multiplier
	finalArmor := targetHealth.GetFinalArmor()
	armorMultiplier := 100.0 / (100.0 + finalArmor)

	// 5. Durability Multiplier (Inverse of Damage Reduction)
	durabilityMultiplier := 1.0 - targetHealth.GetFinalDurability()

	// Calculate Damage Stages
	preMitigationPhysical := basePhysical * critMultiplierEV * ampMultiplier
	mitigatedPhysical := preMitigationPhysical * armorMultiplier
	finalPhysical := mitigatedPhysical * durabilityMultiplier
	finalTotalDamage := finalPhysical // Auto-attacks are purely physical

	// Update Attacker's LastAttackTime
	attackerAttack.SetLastAttackTime(evt.Timestamp)

	// Enqueue DamageAppliedEvent
	damageAppliedEvent := eventsys.DamageAppliedEvent{
		Source:                attacker,
		Target:                target,
		Timestamp:             evt.Timestamp,
		IsSpell:               false,
		PreMitigationPhysical: preMitigationPhysical,
		PreMitigationMagic:    0, // No magic damage from basic attack
		FinalPhysicalDamage:   finalPhysical,
		FinalMagicDamage:      0, // No magic damage
		FinalTotalDamage:      finalTotalDamage,
	}
	s.eventBus.Enqueue(damageAppliedEvent)

	log.Printf("DamageSystem (onAttackLanded): Calculated %.1f final physical damage from %d to %d.\n", finalTotalDamage, attacker, target)
}

// onDamageApplied applies the final damage to the target and handles mana gain/death.
func (s *DamageSystem) onDamageApplied(evt eventsys.DamageAppliedEvent) {
	attacker := evt.Source
	target := evt.Target
	// Use the total final damage calculated previously
	finalDamageToApply := evt.FinalTotalDamage

	// Get target health
	targetHealth, okHealth := s.world.GetHealth(target)
	if !okHealth {
		log.Printf("DamageSystem Error: Target %d has no Health component in onDamageApplied.\n", target)
		return
	}

	// Get attacker info
	attackerInfo, _ := s.world.GetChampionInfo(attacker)
	attackerName := fmt.Sprintf("Entity %d", attacker)
	if attackerInfo != nil {
		attackerName = attackerInfo.Name
	}

	// Get target info
	targetInfo, _ := s.world.GetChampionInfo(target)
	targetName := fmt.Sprintf("Entity %d", target)
	if targetInfo != nil {
		targetName = targetInfo.Name
	}

	// Apply damage
	initialHP := targetHealth.CurrentHP
	targetHealth.CurrentHP -= finalDamageToApply // Apply the total damage

	displayHealth := targetHealth.CurrentHP
	if displayHealth < 0 {
		displayHealth = 0
	}

	// Updated Log Message
	damageType := "Attack"
	if evt.IsSpell {
		damageType = "Spell"
	}
	log.Printf("DamageSystem (onDamageApplied): %s hits %s with %s for %.1f total damage (%.1f Phys, %.1f Mag | PreMit: %.1f Phys, %.1f Mag) (%.1f -> %.1f HP)\n",
		attackerName, targetName, damageType,
		evt.FinalTotalDamage, evt.FinalPhysicalDamage, evt.FinalMagicDamage,
		evt.PreMitigationPhysical, evt.PreMitigationMagic,
		initialHP, displayHealth)

	// Check for death
	if targetHealth.CurrentHP <= 0 && initialHP > 0 {
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
		log.Printf("DamageSystem (onDamageApplied): %s gets kill credit for %s.\n", attackerName, targetName)
	}

	// --- Mana Gain for Attacker (Only on Auto-Attacks) ---
	if !evt.IsSpell { // Check the flag
		attackerMana, okMana := s.world.GetMana(attacker)
		if okMana {
			manaGain := 10.0 // Standard TFT mana gain per auto-attack
			attackerMana.AddCurrentMana(manaGain)

			// Clamp mana to max
			if attackerMana.GetCurrentMana() > attackerMana.GetMaxMana() {
				attackerMana.SetCurrentMana(attackerMana.GetMaxMana())
			}
			log.Printf("DamageSystem (onDamageApplied): %s gains %.1f mana (now %.1f / %.1f)\n", attackerName, manaGain, attackerMana.GetCurrentMana(), attackerMana.GetMaxMana())
		} else {
			// This might be expected for units without mana
			// log.Printf("DamageSystem: Warning: Attacker %s has no Mana component, cannot gain mana.\n", attackerName)
		}
	}
	// --- End Mana Gain ---

	// TODO: Target also gains mana on being attacked
}

func (s *DamageSystem) onSpellCast(evt eventsys.SpellCastEvent) {
	caster := evt.Source
	target := evt.Target

	casterSpell, okSpell := s.world.GetSpell(caster)
	if !okSpell {
		log.Printf("DamageSystem Error: Caster %d has no Spell component in onSpellCast.\n", caster)
		return
	}
	casterAttack, okAttack := s.world.GetAttack(caster) // Needed for AD scaling
	if !okAttack {
		log.Printf("DamageSystem Error: Caster %d has no Attack component in onSpellCast.\n", caster)
		return
	}
	casterCrit, okCrit := s.world.GetCrit(caster) // Get Crit component
	if !okCrit {
		log.Printf("DamageSystem Error: Caster %d has no Crit component in onSpellCast.\n", caster)
		return
	}
	targetHealth, okHp := s.world.GetHealth(target)
	if !okHp {
		log.Printf("DamageSystem Error: Target %d has no Health component in onSpellCast.\n", target)
		return
	}

	// --- Base Spell Damage Components ---
	rawPhysicalDamage := 0.0
	// rawPhysicalDamage += casterSpell.GetVarPercentADDamage() * casterAttack.GetFinalAD()
	rawMagicDamage := casterSpell.GetFinalAP()
	// rawMagicDamage += casterSpell.GetVarBaseDamage() + (casterSpell.GetVarAPScaling() * casterSpell.GetFinalAP())

	// --- Crit Multiplier (Expected Value) ---
	// Check if abilities can crit (requires CanAbilityCrit flag, set by items/traits)
	itemCritMarkerType := reflect.TypeOf(components.CanAbilityCritFromItems{})
	traitCritMarkerType := reflect.TypeOf(components.CanAbilityCritFromTraits{})
	_, hasItemCritMarker := s.world.GetComponent(caster, itemCritMarkerType)
	_, hasTraitCritMarker := s.world.GetComponent(caster, traitCritMarkerType)
	canAbilitiesCrit := hasItemCritMarker || hasTraitCritMarker

	critChance := 0.0
	critMultiplier := 1.0 // Default to 1.0 if cannot crit
	if canAbilitiesCrit {
		critChance = casterCrit.GetFinalCritChance()
		critMultiplier = casterCrit.GetFinalCritMultiplier()
	}
	critMultiplierEV := (1.0 - critChance) + (critChance * critMultiplier)

	// --- Amplification Multiplier ---
	// TODO: Distinguish between Attack Amp and Spell Amp if necessary. Using Attack Amp for now.
	ampMultiplier := 1.0 + casterAttack.GetFinalDamageAmp()

	// --- Resistance Multipliers ---
	finalArmor := targetHealth.GetFinalArmor()
	armorMultiplier := 100.0 / (100.0 + finalArmor)
	finalMR := targetHealth.GetFinalMR()
	mrMultiplier := 100.0 / (100.0 + finalMR)

	// --- Durability Multiplier ---
	durabilityMultiplier := 1.0 - targetHealth.GetFinalDurability()

	// --- Calculate Damage Stages ---
	// Apply Crit and Amp
	preMitigationPhysical := rawPhysicalDamage * critMultiplierEV * ampMultiplier
	preMitigationMagic := rawMagicDamage * critMultiplierEV * ampMultiplier

	// Apply Resistance
	mitigatedPhysical := preMitigationPhysical * armorMultiplier
	mitigatedMagic := preMitigationMagic * mrMultiplier

	// Apply Durability
	finalPhysical := mitigatedPhysical * durabilityMultiplier
	finalMagic := mitigatedMagic * durabilityMultiplier

	// Total Final Damage
	finalTotalDamage := finalPhysical + finalMagic

	// --- Enqueue DamageAppliedEvent ---
	damageAppliedEvent := eventsys.DamageAppliedEvent{
		Source:                caster,
		Target:                target,
		Timestamp:             evt.Timestamp,
		IsSpell:               true,
		PreMitigationPhysical: preMitigationPhysical,
		PreMitigationMagic:    preMitigationMagic,
		FinalPhysicalDamage:   finalPhysical,
		FinalMagicDamage:      finalMagic,
		FinalTotalDamage:      finalTotalDamage,
	}
	s.eventBus.Enqueue(damageAppliedEvent)

	log.Printf("DamageSystem (onSpellCast): Calculated %.1f total final damage (%.1f Phys, %.1f Mag) from %d to %d.\n",
		finalTotalDamage, finalPhysical, finalMagic, caster, target)
}
