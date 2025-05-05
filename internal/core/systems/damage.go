package systems

import (
	"fmt"
	"log"
	"reflect"

	"tft-dps-simulator/internal/core/components"
	"tft-dps-simulator/internal/core/ecs"
	eventsys "tft-dps-simulator/internal/core/systems/events"
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
	case eventsys.SpellLandedEvent:
		s.onSpellLanded(event)
		// Add cases for other event types this system might handle
	}
}

func (s *DamageSystem) CanHandle(evt interface{}) bool {
	switch evt.(type) {
	case eventsys.AttackLandedEvent, eventsys.DamageAppliedEvent, eventsys.SpellLandedEvent:
		return true
	default:
		return false
	}
}


// onAttackLanded calculates final damage from an attack and enqueues DamageAppliedEvent.
func (s *DamageSystem) onAttackLanded(evt eventsys.AttackLandedEvent) {
	attacker := evt.Source
    target := evt.Target
    eventTime := evt.Timestamp // Use timestamp from the event

    // --- Get Components (same as before) ---
    attackerAttack, okAtk := s.world.GetAttack(attacker)
    if !okAtk {
        log.Printf("DamageSystem Error: Attacker %d has no Attack component in onAttackFired.\n", attacker)
        return
    }
    attackerCrit, okCrit := s.world.GetCrit(attacker)
    if !okCrit {
        log.Printf("DamageSystem Error: Attacker %d has no Crit component in onAttackFired.\n", attacker)
        return
    }
    targetHealth, okHp := s.world.GetHealth(target)
    if !okHp {
        log.Printf("DamageSystem Error: Target %d has no Health component in onAttackFired.\n", target)
        return
    }

    // --- Damage Calculation ---
    // 1. Raw Damage (Base AD)
    rawDamage := attackerAttack.GetFinalAD()

    // 2. Crit Check & Multiplier
    // TODO: Implement actual random crit check instead of EV for more accurate simulation?
    // For now, using Expected Value (EV) based on crit chance.
    isCrit := false // Placeholder for actual crit check
    critChance := attackerCrit.GetFinalCritChance()
    critMultiplier := attackerCrit.GetFinalCritMultiplier()

    // Using EV for pre-mitigation calculation for now, consistent with previous logic
    critMultiplierEV := (1.0 - critChance) + (critChance * critMultiplier)

    // 3. Amplification Multiplier
    ampMultiplier := 1.0 + attackerAttack.GetFinalDamageAmp()

    // 4. Pre-Mitigation Damage
    preMitigationDamage := rawDamage * critMultiplierEV * ampMultiplier

    // 5. Armor Reduction Multiplier & Mitigation Amount
    finalArmor := targetHealth.GetFinalArmor()
    armorMultiplier := 100.0 / (100.0 + finalArmor) // Damage multiplier based on armor
    mitigatedByArmor := preMitigationDamage * (1.0 - armorMultiplier)

    // 6. Durability Multiplier & Mitigation Amount
    finalDurability := targetHealth.GetFinalDurability()
    durabilityMultiplier := 1.0 - finalDurability // Damage multiplier based on durability
    mitigatedByDurability := (preMitigationDamage - mitigatedByArmor) * (1.0 - durabilityMultiplier)

    // 7. Final Damage & Total Mitigation
    finalDamage := preMitigationDamage * armorMultiplier * durabilityMultiplier
    totalMitigation := mitigatedByArmor + mitigatedByDurability

    // --- Enqueue DamageAppliedEvent ---
    damageAppliedEvent := eventsys.DamageAppliedEvent{
        Source:           attacker,
        Target:           target,
        Timestamp:        eventTime,
        DamageType:       "AD", // Attacks are physical
        DamageSource:     "Attack",
        RawDamage:        rawDamage,
        PreMitigationDamage: preMitigationDamage,
        MitigatedDamage:  totalMitigation,
        FinalTotalDamage:      finalDamage,
        IsCrit:           isCrit, // Use actual crit result if implemented
        IsAbilityCrit:    false,  // Attacks are not ability crits
    }
    s.eventBus.Enqueue(damageAppliedEvent, eventTime) // Use eventTime for enqueueing

    log.Printf("DamageSystem (onAttackLanded): Calculated %.1f final physical damage from %d to %d. Enqueued DamageAppliedEvent.", finalDamage, attacker, target)
}

// onDamageApplied applies the final damage to the target and handles mana gain/death.
func (s *DamageSystem) onDamageApplied(evt eventsys.DamageAppliedEvent) {
    attacker := evt.Source // Renaming for clarity, though it's evt.Source
    target := evt.Target
    finalDamageToApply := evt.FinalTotalDamage 

    // Get target health
    targetHealth, okHealth := s.world.GetHealth(target)
    if !okHealth {
        log.Printf("DamageSystem Error: Target %d has no Health component in onDamageApplied.\n", target)
        return
    }

    // Get attacker/target names (optional, for logging)
    attackerName := fmt.Sprintf("Entity %d", attacker)
    targetName := fmt.Sprintf("Entity %d", target)

    // Apply damage
    initialHP := targetHealth.GetCurrentHP()
    if initialHP <= 0 {
         log.Printf("DamageSystem (onDamageApplied): Target %s already defeated. Ignoring damage.", targetName)
         return // Don't apply damage or trigger effects if already dead
    }
    targetHealth.SetCurrentHP(initialHP - finalDamageToApply)

    displayHealth := targetHealth.CurrentHP
    if displayHealth < 0 {
        displayHealth = 0
    }

    // Updated Log Message using DamageAppliedEvent fields
    log.Printf("DamageSystem (onDamageApplied): %s hits %s with %s (%s) for %.1f damage (Raw: %.1f, PreMit: %.1f, Mit: %.1f). HP: %.1f -> %.1f",
        attackerName, targetName, evt.DamageSource, evt.DamageType,
        evt.FinalTotalDamage, evt.RawDamage, evt.PreMitigationDamage, evt.MitigatedDamage,
        initialHP, displayHealth)
    // if evt.IsCrit || evt.IsAbilityCrit {
    //      log.Printf("  (Crit!)")
    // }

    // Check for death
    if targetHealth.CurrentHP <= 0 && initialHP > 0 { // Ensure initialHP was > 0 to only trigger once
        log.Printf("DamageSystem (onDamageApplied): %s has been defeated!\n", targetName)
        deathEvent := eventsys.DeathEvent{
            Target:    target,
            Timestamp: evt.Timestamp,
        }
        s.eventBus.Enqueue(deathEvent, evt.Timestamp)

        killEvent := eventsys.KillEvent{
            Killer:    attacker,
            Victim:    target,
            Timestamp: evt.Timestamp,
        }
        s.eventBus.Enqueue(killEvent, evt.Timestamp)
        log.Printf("DamageSystem (onDamageApplied): %s gets kill credit for %s.\n", attackerName, targetName)
    }

    // --- Mana Gain ---
    // Attacker gains mana only if source was "Attack"
    if evt.DamageSource == "Attack" {
        attackerMana, okMana := s.world.GetMana(attacker)
        if okMana {
            manaGain := 10.0 // Standard TFT mana gain per auto-attack
            attackerMana.AddCurrentMana(manaGain)
            log.Printf("DamageSystem (onDamageApplied): Attacker %s gains %.1f mana (now %.1f / %.1f)\n", attackerName, manaGain, attackerMana.GetCurrentMana(), attackerMana.GetMaxMana())
        } else {
             log.Printf("DamageSystem: Warning: Attacker %s has no Mana component, cannot gain mana.\n", attackerName)
        }
    }

    // Update attacker's DamageStats
    attackerDamageStats, okStats := s.world.GetDamageStats(attacker)
    if okStats {
        attackerDamageStats.TotalDamage += finalDamageToApply
        switch evt.DamageSource {
        case "Attack":
            attackerDamageStats.AutoAttackDamage += finalDamageToApply
        case "Spell":
            attackerDamageStats.SpellDamage += finalDamageToApply
        }

        switch evt.DamageType {
        case "AD":
            attackerDamageStats.TotalADDamage += finalDamageToApply
        case "AP":
            attackerDamageStats.TotalAPDamage += finalDamageToApply
        case "True":
            attackerDamageStats.TotalTrueDamage += finalDamageToApply
        }
    }
    // Target gains mana on being hit (regardless of source?)
    // TODO: Revisit mana lock during cast
    targetMana, okMana := s.world.GetMana(target)
    if okMana {
        // TODO: Mana gain on hit might scale with damage taken? Using flat 10 for now.
        manaGainOnHit := 10.0
        targetMana.AddCurrentMana(manaGainOnHit)
        log.Printf("DamageSystem (onDamageApplied): Target %s gains %.1f mana from being hit (now %.1f / %.1f)\n", targetName, manaGainOnHit, targetMana.GetCurrentMana(), targetMana.GetMaxMana())
    }
    // No warning if target has no mana, common for dummies/some units.
}

// onSpellLanded calculates final damage from a spell and enqueues DamageAppliedEvent.
// Triggered by SpellLandedEvent.
func (s *DamageSystem) onSpellLanded(evt eventsys.SpellLandedEvent) { // Changed event type
    caster := evt.Source
    target := evt.Target
    eventTime := evt.Timestamp // Use timestamp from the event

    // --- Get Components ---
    casterSpell, okSpell := s.world.GetSpell(caster)
    if !okSpell {
        log.Printf("DamageSystem Error: Caster %d has no Spell component in onSpellLanded.\n", caster)
        return
    }
    casterAttack, okAttack := s.world.GetAttack(caster) // Needed for AD scaling/Amp
    if !okAttack {
        log.Printf("DamageSystem Error: Caster %d has no Attack component in onSpellLanded.\n", caster)
        return
    }
    casterCrit, okCrit := s.world.GetCrit(caster) // Needed for Ability Crit check
    if !okCrit {
        log.Printf("DamageSystem Error: Caster %d has no Crit component in onSpellLanded.\n", caster)
        return
    }
    targetHealth, okHp := s.world.GetHealth(target)
    if !okHp {
        log.Printf("DamageSystem Error: Target %d has no Health component in onSpellLanded.\n", target)
        return
    }

    // --- Base Spell Damage Components ---
    // TODO: Implement proper spell data lookup based on evt.SpellName
    // For now, assume simple AP scaling magic damage as before.
    rawPhysicalDamage := 0.0
    rawMagicDamage := casterSpell.GetFinalAP() // Placeholder: Use actual spell base + scaling
    rawTrueDamage := 0.0                       // Placeholder for true damage spells

    // Determine primary damage type for resistance calculation (simplification)
    damageType := "AP" // Default for placeholder
    rawDamage := rawMagicDamage // Placeholder
    if rawPhysicalDamage > rawMagicDamage && rawPhysicalDamage > rawTrueDamage {
        damageType = "AD"
        rawDamage = rawPhysicalDamage
    } else if rawTrueDamage > rawMagicDamage {
        damageType = "True"
        rawDamage = rawTrueDamage
    }

    // --- Crit Check & Multiplier ---
    itemCritMarkerType := reflect.TypeOf(components.CanAbilityCritFromItems{})
    traitCritMarkerType := reflect.TypeOf(components.CanAbilityCritFromTraits{})
    _, hasItemCritMarker := s.world.GetComponent(caster, itemCritMarkerType)
    _, hasTraitCritMarker := s.world.GetComponent(caster, traitCritMarkerType)
    canAbilitiesCrit := hasItemCritMarker || hasTraitCritMarker

    isAbilityCrit := false // Placeholder for actual crit check
    critChance := 0.0
    critMultiplier := 1.0
    if canAbilitiesCrit {
        critChance = casterCrit.GetFinalCritChance()
        critMultiplier = casterCrit.GetFinalCritMultiplier()
        // TODO: Implement actual random crit check: isAbilityCrit = rand.Float64() < critChance
    }

    // Using EV for pre-mitigation calculation for now
    critMultiplierEV := (1.0 - critChance) + (critChance * critMultiplier)

    // --- Amplification Multiplier ---
    // TODO: Use Spell Amp if available, otherwise fallback to Attack Amp?
    ampMultiplier := 1.0 + casterAttack.GetFinalDamageAmp() // Using Attack Amp for now

    // --- Pre-Mitigation Damage ---
    preMitigationDamage := rawDamage * critMultiplierEV * ampMultiplier

    // --- Resistance Multipliers & Mitigation Amount ---
    resistanceMultiplier := 1.0
    mitigatedByResistance := 0.0
    if damageType == "AD" {
        finalArmor := targetHealth.GetFinalArmor()
        resistanceMultiplier = 100.0 / (100.0 + finalArmor)
        mitigatedByResistance = preMitigationDamage * (1.0 - resistanceMultiplier)
    } else if damageType == "AP" {
        finalMR := targetHealth.GetFinalMR()
        resistanceMultiplier = 100.0 / (100.0 + finalMR)
        mitigatedByResistance = preMitigationDamage * (1.0 - resistanceMultiplier)
    }
    // True damage ignores resistance (multiplier remains 1.0)

    // --- Durability Multiplier & Mitigation Amount ---
    finalDurability := targetHealth.GetFinalDurability()
    durabilityMultiplier := 1.0 - finalDurability
    mitigatedByDurability := (preMitigationDamage - mitigatedByResistance) * (1.0 - durabilityMultiplier)
    // True damage ignores durability? Check TFT rules. Assuming it does for now.
    if damageType == "True" {
         durabilityMultiplier = 1.0
         mitigatedByDurability = 0.0
    }


    // --- Final Damage & Total Mitigation ---
    finalDamage := preMitigationDamage * resistanceMultiplier * durabilityMultiplier
    totalMitigation := mitigatedByResistance + mitigatedByDurability

    // --- Enqueue DamageAppliedEvent ---
    damageAppliedEvent := eventsys.DamageAppliedEvent{
        Source:           caster,
        Target:           target,
        Timestamp:        eventTime,
        DamageType:       damageType,
        DamageSource:     "Spell", //+ evt.SpellName, // Include spell name
        RawDamage:        rawDamage,
        PreMitigationDamage: preMitigationDamage,
        MitigatedDamage:  totalMitigation,
        FinalTotalDamage:      finalDamage,
        IsCrit:           false, // Spells don't trigger basic attack crit flag
        IsAbilityCrit:    isAbilityCrit, // Use actual crit result if implemented
    }
    s.eventBus.Enqueue(damageAppliedEvent, eventTime) // Use eventTime for enqueueing

    log.Printf("DamageSystem (onSpellLanded): Calculated %.1f final %s damage from %s by %d to %d. Enqueued DamageAppliedEvent.", finalDamage, damageType, evt.SpellName, caster, target)


}