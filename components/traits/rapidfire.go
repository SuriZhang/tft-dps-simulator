package traits

// RapidfireEffect holds the dynamic stacking state for the Rapidfire trait.
// Only Rapidfire champions will have this component, and it will be added/removed as needed.
type RapidfireEffect struct {
    currentStacks       int
    maxStacks           int
	isMaxStacks		bool
    AttackSpeedPerStack float64
}

// NewRapidfireEffect creates a new RapidfireState component.
func NewRapidfireEffect(maxStacks int, asPerStack float64) *RapidfireEffect {
    return &RapidfireEffect{
        currentStacks:       0,
        maxStacks:           maxStacks,
        AttackSpeedPerStack: asPerStack,
    }
}

// IncrementStacks increases stack count, returns true if stacks changed.
func (rs *RapidfireEffect) IncrementStacks() (stackAdded, reachedMax bool) {
	if rs.currentStacks < rs.maxStacks {
        rs.currentStacks ++
        stackAdded = true
        if rs.currentStacks == rs. maxStacks && !rs.isMaxStacks {
            rs.isMaxStacks = true
            reachedMax = true
        }
        return stackAdded, reachedMax
    }
    return false, false
}

// GetCurrentBonusAS calculates the bonus AS from stacks.
func (rs *RapidfireEffect) GetCurrentBonusAS() float64 {
    return float64(rs.currentStacks) * rs.AttackSpeedPerStack
}

// GetcurrentStacks returns the current stack count.
func (rs *RapidfireEffect) GetCurrentStacks() int {
	return rs.currentStacks
}

// GetMaxStacks returns the maximum stack count.
func (rs *RapidfireEffect) GetMaxStacks() int {
	return rs.maxStacks
}

// GetAttackSpeedPerStack returns the attack speed per stack.
func (rs *RapidfireEffect) GetAttackSpeedPerStack() float64 {
	return rs.AttackSpeedPerStack
}

// Reset clears the stacks.
func (rs *RapidfireEffect) Reset() {
    rs.currentStacks = 0
}