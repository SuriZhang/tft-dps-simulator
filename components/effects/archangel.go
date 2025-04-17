package effects

// ArchangelsEffect holds the state for Archangel's Staff.
/**
 * "desc": "Combat start: Gain @APPerInterval@ Ability Power every @IntervalSeconds@ seconds in combat.",
 */

type ArchangelsEffect struct {
    timer      float64 // Time elapsed since the last stack gain.
    stacks     int     // Number of AP stacks gained.
    interval   float64 // Time interval to gain a stack (e.g., 5.0 seconds).
    apPerInterval float64 // AP gained per stack.
}

// NewArchangelsEffect creates a new ArchangelsEffect component using values
// typically derived from item data.
func NewArchangelsEffect(intervalSeconds, apPerInterval float64) *ArchangelsEffect {
    return &ArchangelsEffect{
        timer:      0.0,
        stacks:     0,
        interval:   intervalSeconds, // Use parameter
        apPerInterval: apPerInterval,   // Use parameter
    }
}

// Timer returns the time elapsed since the last stack gain.
func (a *ArchangelsEffect) GetTimer() float64 {
    return a.timer
}

// SetTimer sets the time elapsed since the last stack gain.
func (a *ArchangelsEffect) SetTimer(timer float64) {
    a.timer = timer
}

func (a *ArchangelsEffect) AddTimer(deltaTime float64) {
    a.timer += deltaTime
}

func (a *ArchangelsEffect) MinusInterval() {
    a.timer -= a.interval
    if a.timer < 0 {
        a.timer = 0 // Ensure timer doesn't go negative
    }
}

// Stacks returns the number of AP stacks gained.
func (a *ArchangelsEffect) GetStacks() int {
    return a.stacks
}

// SetStacks sets the number of AP stacks gained.
func (a *ArchangelsEffect) SetStacks(stacks int) {
    a.stacks = stacks
}

// Interval returns the time interval to gain a stack.
func (a *ArchangelsEffect) GetInterval() float64 {
    return a.interval
}

// SetInterval sets the time interval to gain a stack.
func (a *ArchangelsEffect) SetInterval(interval float64) {
    a.interval = interval
}

// AddStack increments the number of stacks by deltaStacks.
func (a *ArchangelsEffect) AddStacks(deltaStacks int) {
    a.stacks += deltaStacks
}

func (a *ArchangelsEffect) ResetEffects() {
    a.stacks = 0
    a.timer = 0.0 // Reset timer as well
}

// APPerStack returns the AP gained per stack.
func (a *ArchangelsEffect) GetAPPerInterval() float64 {
    return a.apPerInterval
}

// SetAPPerInterval sets the AP gained per stack.
func (a *ArchangelsEffect) SetAPPerInterval(apPerStack float64) {
    a.apPerInterval = apPerStack
}
