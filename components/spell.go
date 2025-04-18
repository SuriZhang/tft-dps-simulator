package components

// Spell holds data related to a champion's ability/spell.
type Spell struct {
	// Base stats (potentially loaded from champion data)
	BaseAP   float64
	ManaCost float64
	Cooldown float64 // Base cooldown

	// Bonus stats accumulated from items, traits, etc.
	BonusAP float64
	// Add other potential bonus spell stats (e.g., BonusSpellCritChance, BonusSpellCritDamage) if needed

	// Final calculated stats used by systems
	FinalAP float64

	// State
	CurrentCooldown float64
	// CanCast bool // Could be managed by SpellSystem
}

// NewSpell creates a Spell component, potentially initializing from base stats.
func NewSpell(baseAP, manaCost, cooldown float64) *Spell {
	return &Spell{
		BaseAP:          baseAP,
		ManaCost:        manaCost,
		Cooldown:        cooldown,
		BonusAP:         0,      // Start with no bonus
		FinalAP:         baseAP, // Initial final value
		CurrentCooldown: 0,
	}
}

// GetBaseAP returns the base ability power.
func (s *Spell) GetBaseAP() float64 {
	return s.BaseAP
}

// SetBaseAP sets the base ability power.
func (s *Spell) SetBaseAP(value float64) {
	s.BaseAP = value
	// Potentially recalculate FinalAP here if needed
}

// GetManaCost returns the mana cost.
func (s *Spell) GetManaCost() float64 {
	return s.ManaCost
}

// SetManaCost sets the mana cost.
func (s *Spell) SetManaCost(value float64) {
	s.ManaCost = value
}

// GetCooldown returns the base cooldown.
func (s *Spell) GetCooldown() float64 {
	return s.Cooldown
}

// SetCooldown sets the base cooldown.
func (s *Spell) SetCooldown(value float64) {
	s.Cooldown = value
}

// GetBonusAP returns the bonus ability power.
func (s *Spell) GetBonusAP() float64 {
	return s.BonusAP
}

// SetBonusAP sets the bonus ability power.
func (s *Spell) SetBonusAP(value float64) {
	s.BonusAP = value
	// Potentially recalculate FinalAP here if needed
}

// GetFinalAP returns the final calculated ability power.
func (s *Spell) GetFinalAP() float64 {
	return s.FinalAP
}

// SetFinalAP sets the final calculated ability power.
// Note: This might be better handled by a dedicated calculation function.
func (s *Spell) SetFinalAP(value float64) {
	s.FinalAP = value
}

// GetCurrentCooldown returns the current cooldown remaining.
func (s *Spell) GetCurrentCooldown() float64 {
	return s.CurrentCooldown
}

// SetCurrentCooldown sets the current cooldown remaining.
func (s *Spell) SetCurrentCooldown(value float64) {
	s.CurrentCooldown = value
}

func (s *Spell) AddBonusAP(value float64) {
	s.BonusAP += value
}

func (s *Spell) ResetBonuses() {
	s.BonusAP = 0
}