package items

// RedBuffEffect holds the state for Red Buff's burn and wound effects.
// "desc": "Attacks and Abilities @BurnPercent@% <TFTKeyword>Burn</TFTKeyword> and @HealingReductionPct@% <TFTKeyword>Wound</TFTKeyword> enemies for @Duration@ seconds."
type RedBuffEffect struct {
    burnPercent         float64 // Burn percentage of max HP per second (e.g., 0.01 for 1%)
    healingReductionPct float64 // Wound healing reduction percentage (e.g., 33.0 for 33%)
    duration            float64 // Duration of both effects in seconds (e.g., 5.0)
}

// NewRedBuffEffect creates a new RedBuffEffect component.
func NewRedBuffEffect(burnPercent, healingReductionPct, duration float64) *RedBuffEffect {
    return &RedBuffEffect{
        burnPercent:         burnPercent / 100.0, // Convert percentage to decimal for burn
        healingReductionPct: healingReductionPct,
        duration:            duration,
    }
}

// GetBurnPercent returns the burn damage percentage.
func (rb *RedBuffEffect) GetBurnPercent() float64 {
    return rb.burnPercent
}

// GetHealingReductionPct returns the wound healing reduction percentage.
func (rb *RedBuffEffect) GetHealingReductionPct() float64 {
    return rb.healingReductionPct
}

// GetDuration returns the duration of both effects.
func (rb *RedBuffEffect) GetDuration() float64 {
    return rb.duration
}

// ResetEffects resets the effect state (if needed for future extensions).
func (rb *RedBuffEffect) ResetEffects() {
    // Currently no internal state to reset, but keeping for consistency
}