package effects

// GuinsoosRagebladeEffect tracks the stacking attack speed bonus from Guinsoo's Rageblade.
type GuinsoosRagebladeEffect struct {
    AttackSpeedPerStack float64 // The % AS gained per stack (e.g., 0.05 for 5%)
    CurrentStacks       int     // Current number of stacks
    // MaxStacks           int     // Optional: If there's a cap in the future
    CurrentBonusAS      float64 // The total bonus AS currently provided by stacks
}

// NewGuinsoosRagebaldeEffect creates a new GuinsoosRagebladeEffect component.
func NewGuinsoosRagebaldeEffect(asPerStack float64 /*, maxStacks int*/) *GuinsoosRagebladeEffect {
    return &GuinsoosRagebladeEffect{
        AttackSpeedPerStack: asPerStack,
        CurrentStacks:       0,
        // MaxStacks:           maxStacks,
        CurrentBonusAS:      0.0,
    }
}

// AddStack increments the stack count and updates the bonus AS.
func (e *GuinsoosRagebladeEffect) IncrementStacks() {
    // if e.MaxStacks > 0 && e.CurrentStacks >= e.MaxStacks {
    //     return // Already at max stacks
    // }
    e.CurrentStacks++
    e.CurrentBonusAS = float64(e.CurrentStacks) * e.AttackSpeedPerStack
}

// GetCurrentBonusAS returns the total bonus attack speed from current stacks.
func (e *GuinsoosRagebladeEffect) GetCurrentBonusAS() float64 {
    return e.CurrentBonusAS
}

// GetCurrentStacks returns the current number of stacks.
func (e *GuinsoosRagebladeEffect) GetCurrentStacks() int {
    return e.CurrentStacks
}

// ResetEffects resets the stacks and bonus AS.
func (e *GuinsoosRagebladeEffect) ResetEffects() {
    e.CurrentStacks = 0
    e.CurrentBonusAS = 0.0
}