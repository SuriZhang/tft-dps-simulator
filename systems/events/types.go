package eventsys

import (
	"github.com/suriz/tft-dps-simulator/ecs"
)

// AttackLandedEvent signifies an attack attempt after range/timing checks.
type AttackLandedEvent struct {
	Source     ecs.Entity
	Target     ecs.Entity
	BaseDamage float64
	Timestamp  float64 // Time the attack was initiated
}

// DamageAppliedEvent signifies final damage being dealt after calculations.
type DamageAppliedEvent struct {
	Source      ecs.Entity
	Target      ecs.Entity
	FinalDamage float64
	Timestamp   float64 // Time the damage was applied
}

// DeathEvent signifies an entity's HP reached zero or below.
type DeathEvent struct {
	Target    ecs.Entity
	Timestamp float64 // Time of death
}
