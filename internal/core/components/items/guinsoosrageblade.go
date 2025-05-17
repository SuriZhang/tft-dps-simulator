package items

// ArchangelsEffect holds the state for Archangel's Staff.
/**
 *  "desc": "Gain @AttackSpeedPerStack@% stacking Attack Speed every second.",
 */

type GuinsoosRagebladeEffect struct {
	stacks        int     // Number of AP stacks gained.
	interval      float64 // Time interval to gain a stack (e.g., 5.0 seconds).
	attackSpeedPerStack float64 // AP gained per stack.
}

func NewGuinsoosRagebladeEffect(intervalSeconds, attackSpeedPerStack float64) *GuinsoosRagebladeEffect {
	return &GuinsoosRagebladeEffect{
		stacks:        0,
		interval:      intervalSeconds,
		attackSpeedPerStack: attackSpeedPerStack, 
	}
}

// Stacks returns the number of AP stacks gained.
func (g *GuinsoosRagebladeEffect) GetStacks() int {
	return g.stacks
}

// SetStacks sets the number of AP stacks gained.
func (g *GuinsoosRagebladeEffect) SetStacks(stacks int) {
	g.stacks = stacks
}

// Interval returns the time interval to gain a stack.
func (g *GuinsoosRagebladeEffect) GetInterval() float64 {
	return g.interval
}

// SetInterval sets the time interval to gain a stack.
func (g *GuinsoosRagebladeEffect) SetInterval(interval float64) {
	g.interval = interval
}

// IncrementStacks increments the number of stacks by 1.
func (g *GuinsoosRagebladeEffect) IncrementStacks() {
	g.stacks += 1
}

func (g *GuinsoosRagebladeEffect) ResetEffects() {
	g.stacks = 0
}

// APPerStack returns the AP gained per stack.
func (g *GuinsoosRagebladeEffect) GetAttackSpeedPerStack() float64 {
	return g.attackSpeedPerStack
}

// SetAttackSpeedPerStack sets the AP gained per stack.
func (g *GuinsoosRagebladeEffect) SetAttackSpeedPerStack(attackSpeedPerStack float64) {
	g.attackSpeedPerStack = attackSpeedPerStack
}

// GetCurrentBonusAS returns the total bonus attack speed from current stacks.
func (g *GuinsoosRagebladeEffect) GetCurrentBonusAS() float64 {
	return float64(g.stacks) * g.attackSpeedPerStack
}

// GetCurrentStacks returns the current number of stacks.
func (g *GuinsoosRagebladeEffect) GetCurrentStacks() int {
	return g.stacks
}
