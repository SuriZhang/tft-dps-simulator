package components

// Attack contains champion attack information
type Attack struct {
    Damage         float64
    AttackSpeed          float64
    Range          float64
    CritChance     float64
    CritMultiplier float64
    LastAttackTime float64 // For tracking attack cooldown
}

// NewAttack creates an Attack component
func NewAttack(damage, attackSpeed, attackRange, critChance, critMulti float64) Attack {
    return Attack{
        Damage:         damage,
        AttackSpeed:          attackSpeed,
        Range:          attackRange,
        CritChance:     critChance,
        CritMultiplier: critMulti,
        LastAttackTime: 0,
    }
}