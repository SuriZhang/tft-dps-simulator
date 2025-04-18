package eventsys

import (
	"github.com/suriz/tft-dps-simulator/ecs"
)

// AttackLandedEvent fires the moment an auto‑attack lands (before mitigation).
type AttackLandedEvent struct {
    Source     ecs.Entity
    Target     ecs.Entity
    BaseDamage int
}

// DamageAppliedEvent fires after all resistances/crit/etc have been calculated.
type DamageAppliedEvent struct {
    Source      ecs.Entity
    Target      ecs.Entity
    FinalDamage int
}