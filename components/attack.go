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

// UpdateDamage updates the attack damage
func (a *Attack) UpdateDamage(damage float64) {
	a.Damage = damage
}

// UpdateAttackSpeed updates the attack speed
func (a *Attack) UpdateAttackSpeed(attackSpeed float64) {
	a.AttackSpeed = attackSpeed
}	

// UpdateRange updates the attack range
func (a *Attack) UpdateRange(rangeVal float64) {
	a.Range = rangeVal
}

// UpdateCritChance updates the critical chance	
func (a *Attack) UpdateCritChance(critChance float64) {
	a.CritChance = critChance
}

// UpdateCritMultiplier updates the critical damage multiplier
func (a *Attack) UpdateCritMultiplier(critMulti float64) {
	a.CritMultiplier = critMulti
}

// UpdateLastAttackTime updates the last attack time
func (a *Attack) UpdateLastAttackTime(lastAttackTime float64) {
	a.LastAttackTime = lastAttackTime
}
