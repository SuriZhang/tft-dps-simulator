package components

// Spell holds data related to a champion's ability/spell.
type Spell struct {
	Name    string
	icon string 
	// Base stats (potentially loaded from champion data)
	BaseAP   float64
	ManaCost float64
	castStartUp float64 // The time it takes to cast the spell. This is the time before the spell animation starts.
	castRecovery float64 // the time after the spell animation finishes before the next spell can be cast.The period where the champion is locked out of auto-attacking is the cast animation time or cast lockout.
	lockManaDuringCast bool // Whether the champion should gain mana during the cast animation

	// --- Spell Variables (should be read from champion ability variables) ---
	// TODO: Load these dynamically based on the spell/champion
	// VarBaseDamage       float64 
	// VarPercentADDamage  float64 
	// VarAPScaling        float64 

	// Bonus stats accumulated from items, traits, etc.
	BonusAP              float64

	// Final calculated stats used by systems
	FinalAP              float64

	// State
	CurrentCooldown float64
	LastCastTime    float64 // For tracking cooldowns
}

// NewSpell creates a Spell component, potentially initializing from base stats.
func NewSpell(name, icon string, manaCost, castStartUp, castRecovery float64) *Spell {
	// TODO: Initialize VarBaseDamage, VarPercentADDamage, VarAPScaling from data
	return &Spell{
		Name:                 name,
		icon:                 icon,
		BaseAP:               100.0, // default base AP in TFT
		ManaCost:             manaCost,
		lockManaDuringCast:   true, // default to lock mana during cast
		castStartUp: castStartUp,
		castRecovery:             castRecovery,
		BonusAP:              0.0,

		FinalAP:              100.0, // init to base AP
		CurrentCooldown:      0.0,
	}
}

// GetName returns the name of the spell.
func (s *Spell) GetName() string {
	return s.Name
}

// GetIcon returns the icon path of the spell.
func (s *Spell) GetIcon() string {
	return s.icon
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

// GetCastStartUp returns the cast start-up time.
func (s *Spell) GetCastStartUp() float64 {
	return s.castStartUp
}

// GetCastRecovery returns the base cooldown.
func (s *Spell) GetCastRecovery() float64 {
	return s.castRecovery
}

// SetCooldown sets the base cooldown.
func (s *Spell) SetCooldown(value float64) {
	s.castRecovery = value
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

// --- Methods for Bonus Spell Stats ---

// --- Methods to SET FINAL calculated spell stats ---

// --- Methods to GET FINAL calculated spell stats ---

// --- Methods for Spell Variables (Example) ---
// TODO: Replace with a better system if needed

// func (s *Spell) SetVarBaseDamage(value float64) {
// 	s.VarBaseDamage = value
// }

// func (s *Spell) GetVarBaseDamage() float64 {
// 	return s.VarBaseDamage
// }

// func (s *Spell) SetVarPercentADDamage(value float64) {
// 	s.VarPercentADDamage = value
// }

// func (s *Spell) GetVarPercentADDamage() float64 {
// 	return s.VarPercentADDamage
// }

// func (s *Spell) SetVarAPScaling(value float64) {
// 	s.VarAPScaling = value
// }

// func (s *Spell) GetVarAPScaling() float64 {
// 	return s.VarAPScaling
// }

// --- Reset Bonus Stats ---
func (s *Spell) ResetBonuses() {
	s.BonusAP = 0.0
	// Do not reset Vars, they are loaded from data
}
