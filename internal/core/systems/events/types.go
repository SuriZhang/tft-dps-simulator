package eventsys

import "tft-dps-simulator/internal/core/ecs"

// ChampionActionEvent signals that a champion should attempt to perform an action (attack or spell).
// This is typically enqueued at t=0 or when a previous action cycle completes.
type ChampionActionEvent struct {
    Entity    ecs.Entity
    Timestamp float64 // Time the action check should occur
}

// AttackStartupEvent signals the beginning of an auto-attack wind-up.
// The AutoAttackSystem should handle this to schedule the AttackFiredEvent.
type AttackStartupEvent struct {
    Entity    ecs.Entity
    Timestamp float64 // Time the attack startup begins
}

// AttackFiredEvent signals the point in an attack animation where the projectile is launched
// or the hit connects (before recovery). Damage calculation might happen later.
// Enqueued by the AutoAttackSystem after startup duration.
type AttackFiredEvent struct {
    Source    ecs.Entity
    Target    ecs.Entity // Target chosen at the start of the attack
    Timestamp float64    // Time the attack fires
}

// AttackLandedEvent is triggered when an auto attack successfully lands (after potential travel time/checks).
// This might be merged with AttackFiredEvent if projectile travel isn't simulated.
// Enqueued by DamageSystem or a dedicated ProjectileSystem.
type AttackLandedEvent struct {
    Source     ecs.Entity
    Target     ecs.Entity
    BaseDamage float64 // Base AD at the time of landing
    Timestamp  float64
}

// AttackRecoveryEndEvent signals the end of the attack recovery period.
// The ChampionActionSystem might listen for this to potentially start the next action.
type AttackRecoveryEndEvent struct {
    Entity    ecs.Entity
    Timestamp float64
}

// AttackCooldownStartEvent signals that an entity should begin its attack cooldown period.
// Enqueued by ChampionActionSystem after recovery if mana is not full.
// Handled by AutoAttackSystem to calculate duration and schedule AttackCooldownEndEvent.
type AttackCooldownStartEvent struct {
    Entity    ecs.Entity
    Timestamp float64 // Time the cooldown should start
}

// AttackCooldownEndEvent signals the end of the attack cooldown (based on AS).
// The ChampionActionSystem might listen for this to potentially start the next action.
type AttackCooldownEndEvent struct {
    Entity    ecs.Entity
    Timestamp float64
}

// DamageAppliedEvent is triggered after damage calculation is complete.
type DamageAppliedEvent struct {
    Source           ecs.Entity
    Target           ecs.Entity
    DamageType       string // "AD", "AP", "True"
    DamageSource     string // "Attack", "Spell", "Item", "Trait", "Burn"
    RawDamage        float64
    PreMitigationDamage float64
    MitigatedDamage  float64
    FinalTotalDamage float64
    IsCrit           bool
    IsAbilityCrit    bool
    Timestamp        float64
}

// DeathEvent signifies an entity's HP reached zero or below.
type DeathEvent struct {
	Target    ecs.Entity
	Timestamp float64 // Time of death
}

// KillEvent signifies that an entity caused another entity's death.
// This is typically triggered alongside a DeathEvent.
type KillEvent struct {
    Killer    ecs.Entity // The entity that dealt the killing blow
    Victim    ecs.Entity // The entity that died
    Timestamp float64    // Time the kill occurred (same as DeathEvent)
}

// AsistEvent signifies that an entity assisted in a kill.
// Dealing Damage: If a unit inflicts damage on an enemy before that enemy is eliminated, it's considered to have participated in the kill.
// TODO: Applying Crowd Control or Debuffs: Units that apply effects like stuns, slows, or other debuffs to an enemy before its death are also considered participants.
// TODO: Providing Support: Units that heal, shield, or buff allies who then secure a kill may be deemed to have participated, depending on the specific trait or augment mechanics
type AssistEvent struct {
    Assistor   ecs.Entity // The entity that assisted
    Victim  ecs.Entity // The entity that was killed
    Timestamp float64 // Time the assist occurred
}

// SpellCastCycleStartEvent signals the beginning of a spell cast animation/channel.
// The SpellCastSystem should handle this to schedule the SpellLandedEvent.
type SpellCastCycleStartEvent struct {
    Entity    ecs.Entity
    Timestamp float64 // Time the cast begins
}

// SpellLandedEvent is triggered when a spell effect should be applied.
// Enqueued by the SpellCastSystem after cast duration.
type SpellLandedEvent struct {
    Source    ecs.Entity
    Target    ecs.Entity // Can be self or other
    SpellName string
    Timestamp float64
    // Add spell-specific payload if needed, or use separate events per spell type
}

// SpellRecoveryEndEvent signals the end of the spell recovery/lockout period.
// The ChampionActionSystem might listen for this to potentially start the next action.
type SpellRecoveryEndEvent struct {
    Entity    ecs.Entity
    Timestamp float64
}

// --- Item Specific Events (Examples) ---

// ArchangelsTickEvent signals a time-based tick for Archangel's Staff.
type ArchangelsTickEvent struct {
    Entity    ecs.Entity
    Timestamp float64
}

// GuinsoosRagebladeTickEvent signals a time-based tick for Guinsoo's Rageblade.
type GuinsoosRagebladeTickEvent struct {
    Entity    ecs.Entity
    Timestamp float64
}

// QuicksilverProcEvent signals a time-based proc for Quicksilver's AS bonus.
type QuicksilverProcEvent struct {
    Entity    ecs.Entity
    Timestamp float64
}

// QuicksilverEndEvent signals the end of Quicksilver's active duration.
type QuicksilverEndEvent struct {
    Entity    ecs.Entity
    Timestamp float64 // Time the effect ends
}

type SpiritVisageHealTickEvent struct { // New Event
    Entity    ecs.Entity
    Timestamp float64
}

// RecalculateStatsEvent signals that an entity's stats need recalculation due to a change.
type RecalculateStatsEvent struct {
    Entity    ecs.Entity
    Timestamp float64 // Time the recalculation is requested (can be same as triggering event)
}