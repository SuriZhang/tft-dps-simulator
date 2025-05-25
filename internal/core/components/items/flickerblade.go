package items

// FlickerbladeEffect holds the state for Flickerblade's stacking effects.
// "desc": "Attacks grant @ASPerStack*100@% stacking Attack Speed. Every @StacksPerBonus@ attacks also grant @ADPerBonus*100@% Attack Damage and @APPerBonus@% Ability Power."
type FlickerbladeEffect struct {
	asPerStack     float64 // Percent AS gained per attack (e.g., 0.08 for 8%)
	adPerBonus     float64 // Percent AD gained every StacksPerBonus attacks (e.g., 0.04 for 4%)
	apPerBonus     float64 // Flat AP gained every StacksPerBonus attacks
	stacksPerBonus float64 // Number of attacks to trigger AD/AP bonus (e.g., 5.0)

	// Internal state
	attackCounter  int     // Counts attacks towards the AD/AP bonus
	totalASApplied float64 // Total %AS dynamically applied by this item's effect
	totalADApplied float64 // Total %AD dynamically applied by this item's effect
	totalAPApplied float64 // Total flat AP dynamically applied by this item's effect
}

// NewFlickerbladeEffect creates a new FlickerbladeEffect component.
func NewFlickerbladeEffect(asPerStack, adPerBonus, apPerBonus, stacksPerBonus float64) *FlickerbladeEffect {
	return &FlickerbladeEffect{
		asPerStack:     asPerStack,
		adPerBonus:     adPerBonus,
		apPerBonus:     apPerBonus,
		stacksPerBonus: stacksPerBonus,
		attackCounter:  0,
		totalASApplied: 0.0,
		totalADApplied: 0.0,
		totalAPApplied: 0.0,
	}
}

// Getters for parameters
func (fe *FlickerbladeEffect) GetASPerStack() float64 {
	return fe.asPerStack
}
func (fe *FlickerbladeEffect) GetADPerBonus() float64 {
	return fe.adPerBonus
}
func (fe *FlickerbladeEffect) GetAPPerBonus() float64 {
	return fe.apPerBonus
}
func (fe *FlickerbladeEffect) GetStacksPerBonus() float64 {
	return fe.stacksPerBonus
}
func (fe *FlickerbladeEffect) GetAttackCounter() int {
	return fe.attackCounter
}

// Getters for total dynamically applied stats (for removal logic)
func (fe *FlickerbladeEffect) GetTotalASApplied() float64 {
	return fe.totalASApplied
}
func (fe *FlickerbladeEffect) GetTotalADApplied() float64 {
	return fe.totalADApplied
}
func (fe *FlickerbladeEffect) GetTotalAPApplied() float64 {
	return fe.totalAPApplied
}

// Setters for total dynamically applied stats
func (fe *FlickerbladeEffect) SetTotalASApplied(value float64) {
	fe.totalASApplied = value
}

func (fe *FlickerbladeEffect) SetTotalADApplied(value float64) {
	fe.totalADApplied = value
}

func (fe *FlickerbladeEffect) SetTotalAPApplied(value float64) {
	fe.totalAPApplied = value	
}

// IncrementAttackCounter increments the attack counter.
func (fe *FlickerbladeEffect) IncrementAttackCounter() {
	fe.attackCounter++
}

// ResetEffects resets all dynamic bonuses and counters.
func (fe *FlickerbladeEffect) ResetEffects() {
	fe.attackCounter = 0
	fe.totalASApplied = 0.0
	fe.totalADApplied = 0.0
	fe.totalAPApplied = 0.0
}
