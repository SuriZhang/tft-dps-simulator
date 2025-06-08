package items

// EvenshroudEffect holds the state for Evenshroud's sunder aura and temporary resistance bonus.
// "desc": "@ARReductionAmount@% <TFTKeyword>Sunder</TFTKeyword> enemies within @HexRange@ hexes. Gain @BonusResists@ Armor and Magic Resist for the first @BonusResistDuration@ seconds of combat."
type EvenshroudEffect struct {
    arReductionAmount   float64 // Armor reduction percentage (e.g., 30.0 for 30%)
    hexRange            float64 // Range in hexes for the sunder effect (e.g., 2.0)
    bonusResists        float64 // Temporary armor and MR bonus (e.g., 25.0)
    bonusResistDuration float64 // Duration of the resistance bonus in seconds (e.g., 15.0)
    resistBonusActive   bool    // Whether the resistance bonus is currently active
}

// NewEvenshroudEffect creates a new EvenshroudEffect component.
func NewEvenshroudEffect(arReductionAmount, hexRange, bonusResists, bonusResistDuration float64) *EvenshroudEffect {
    return &EvenshroudEffect{
        arReductionAmount:   arReductionAmount / 100.0, // Convert percentage to decimal
        hexRange:            hexRange,
        bonusResists:        bonusResists,
        bonusResistDuration: bonusResistDuration,
        resistBonusActive:   false,
    }
}

// GetARReductionAmount returns the armor reduction percentage as decimal.
func (e *EvenshroudEffect) GetARReductionAmount() float64 {
    return e.arReductionAmount
}

// GetHexRange returns the range of the sunder effect.
func (e *EvenshroudEffect) GetHexRange() float64 {
    return e.hexRange
}

// GetBonusResists returns the temporary resistance bonus amount.
func (e *EvenshroudEffect) GetBonusResists() float64 {
    return e.bonusResists
}

// GetBonusResistDuration returns the duration of the resistance bonus.
func (e *EvenshroudEffect) GetBonusResistDuration() float64 {
    return e.bonusResistDuration
}

// IsResistBonusActive returns whether the resistance bonus is currently active.
func (e *EvenshroudEffect) IsResistBonusActive() bool {
    return e.resistBonusActive
}

// ActivateResistBonus activates the temporary resistance bonus.
func (e *EvenshroudEffect) ActivateResistBonus() {
    e.resistBonusActive = true
}

// DeactivateResistBonus deactivates the temporary resistance bonus.
func (e *EvenshroudEffect) DeactivateResistBonus() {
    e.resistBonusActive = false
}

// ResetEffects resets the effect state.
func (e *EvenshroudEffect) ResetEffects() {
    e.resistBonusActive = false
}