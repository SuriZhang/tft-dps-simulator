package components

// Attack contains champion attack information
type Attack struct {
    Damage         float64
    Speed          float64
    Range          float64
    CritChance     float64
    CritMultiplier float64
    LastAttackTime float64 // For tracking attack cooldown
}

// NewAttack creates an Attack component
func NewAttack(damage, speed, attackRange, critChance, critMulti float64) Attack {
    return Attack{
        Damage:         damage,
        Speed:          speed,
        Range:          attackRange,
        CritChance:     critChance,
        CritMultiplier: critMulti,
        LastAttackTime: 0,
    }
}