package traits

// RapidfireEffect holds the dynamic stacking state for the Rapidfire trait.
// Only Rapidfire champions will have this component, and it will be added/removed as needed.
type RapidfireEffect struct {
    CurrentStacks       int
    MaxStacks           int
    AttackSpeedPerStack float64
    // Store the applied bonus separately if needed for removal,
    // or recalculate based on stacks when stats are computed.
}

// NewRapidfireEffect creates a new RapidfireState component.
func NewRapidfireEffect(maxStacks int, asPerStack float64) *RapidfireEffect {
    return &RapidfireEffect{
        CurrentStacks:       0,
        MaxStacks:           maxStacks,
        AttackSpeedPerStack: asPerStack,
    }
}

// IncrementStacks increases stack count, returns true if stacks changed.
func (rs *RapidfireEffect) IncrementStacks() bool {
    if rs.CurrentStacks < rs.MaxStacks {
        rs.CurrentStacks++
        return true
    }
    return false
}

// GetCurrentBonusAS calculates the bonus AS from stacks.
func (rs *RapidfireEffect) GetCurrentBonusAS() float64 {
    return float64(rs.CurrentStacks) * rs.AttackSpeedPerStack
}

// Reset clears the stacks.
func (rs *RapidfireEffect) Reset() {
    rs.CurrentStacks = 0
}