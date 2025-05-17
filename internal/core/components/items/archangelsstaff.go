package items

// ArchangelsEffect holds the state for Archangel's Staff.
/**
 * "desc": "Combat start: Gain @APPerInterval@ Ability Power every @IntervalSeconds@ seconds in combat.",
 */

type ArchangelsStaffEffect struct {
	stacks        int     // Number of AP stacks gained.
	interval      float64 // Time interval to gain a stack (e.g., 5.0 seconds).
	apPerInterval float64 // AP gained per stack.
}

// NewArchangelsEffect creates a new ArchangelsEffect component using values
// typically derived from item data.
func NewArchangelsEffect(intervalSeconds, apPerInterval float64) *ArchangelsStaffEffect {
	return &ArchangelsStaffEffect{
		stacks:        0,
		interval:      intervalSeconds, // Use parameter
		apPerInterval: apPerInterval,   // Use parameter
	}
}

// Stacks returns the number of AP stacks gained.
func (a *ArchangelsStaffEffect) GetStacks() int {
	return a.stacks
}

// SetStacks sets the number of AP stacks gained.
func (a *ArchangelsStaffEffect) SetStacks(stacks int) {
	a.stacks = stacks
}

// Interval returns the time interval to gain a stack.
func (a *ArchangelsStaffEffect) GetInterval() float64 {
	return a.interval
}

// SetInterval sets the time interval to gain a stack.
func (a *ArchangelsStaffEffect) SetInterval(interval float64) {
	a.interval = interval
}

// AddStack increments the number of stacks by deltaStacks.
func (a *ArchangelsStaffEffect) AddStacks(deltaStacks int) {
	a.stacks += deltaStacks
}

func (a *ArchangelsStaffEffect) ResetEffects() {
	a.stacks = 0
}

// APPerStack returns the AP gained per stack.
func (a *ArchangelsStaffEffect) GetAPPerInterval() float64 {
	return a.apPerInterval
}

// SetAPPerInterval sets the AP gained per stack.
func (a *ArchangelsStaffEffect) SetAPPerInterval(apPerStack float64) {
	a.apPerInterval = apPerStack
}
