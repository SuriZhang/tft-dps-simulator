package components

type DamageStats struct {
	TotalDamage           float64 `json:"totalDamage"`
	DamagePerSecond       float64 `json:"dps"`
	TotalADDamage         float64 `json:"totalADDamage"`
	TotalAPDamage         float64 `json:"totalAPDamage"`
	TotalTrueDamage       float64 `json:"totalTrueDamage"`
	TotalAutoAttackCounts int     `json:"totalAutoAttackCounts"`
	TotalSpellCastCounts  int     `json:"totalSpellCastCounts"`
	AutoAttackDamage float64 `json:"autoAttackDamage"`
	SpellDamage      float64 `json:"spellDamage"`
}

func NewDamageStats() DamageStats {
	return DamageStats{
		TotalDamage:           0.0,
		DamagePerSecond:       0.0,
		TotalADDamage:         0.0,
		TotalAPDamage:         0.0,
		TotalTrueDamage:       0.0,
		TotalAutoAttackCounts: 0,
		TotalSpellCastCounts:  0, 
		AutoAttackDamage:      0.0,
		SpellDamage:           0.0,
	}
}

func (ds *DamageStats) SetTotalAutoAttackCounts(count int) {
	ds.TotalAutoAttackCounts = count
}

func (ds *DamageStats) SetTotalSpellCastCounts(count int) {
	ds.TotalSpellCastCounts = count
}