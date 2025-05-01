package items

// QuicksilverEffect holds the state for Quicksilver.
/**
 * "desc": "Combat start: Gain immunity to crowd control for @SpellShieldDuration@ seconds. During this time, gain @ProcAttackSpeed*100@% Attack Speed every @ProcInterval@ seconds.<br><br><tftitemrules>[Unique - only 1 per champion]</tftitemrules>"
 */
type QuicksilverEffect struct {
	spellShieldDuration float64 // Duration of CC immunity (in seconds).
	// CC Immunity State
	isActive          bool    // Flag to indicate if the CC immunity is currently active.

	// Attack Speed Proc State (during immunity)
	procInterval    float64 // Time interval to gain AS stack (e.g., 2.0 seconds).
	stacks          int
	procAttackSpeed float64 // Bonus AS gained per proc (as a decimal, e.g., 0.03 for 3%).
	currentBonusAS  float64 // Total accumulated bonus AS from procs.
}

// NewQuicksilverEffect creates a new QuicksilverEffect component using default values
// based on the typical Quicksilver item data.
// TODO: Ideally, fetch Duration, Interval, and ProcAS from item data via EquipmentManager.
func NewQuicksilverEffect(spellShieldDuration, procAttackSpeed, procInterval float64) *QuicksilverEffect {
	return &QuicksilverEffect{
		spellShieldDuration: spellShieldDuration,
		isActive:          true,

		// Attack Speed Proc Defaults (from JSON)
		stacks:          0, // Starts with no stacks
		procInterval:    procInterval,    // ProcInterval
		procAttackSpeed: procAttackSpeed, // ProcAttackSpeed (approx)
		currentBonusAS:  0.0,             // Starts with no bonus AS
	}
}

// --- CC Immunity Methods ---
// IsActive returns whether the CC immunity effect is currently active.
func (q *QuicksilverEffect) IsActive() bool {
	return q.isActive
}

// SetIsActive explicitly sets the active state of the CC immunity.
// Note: Usually managed internally by DecreaseRemainingDuration.
func (q *QuicksilverEffect) SetIsActive(active bool) {
	q.isActive = active
}

// --- Attack Speed Proc Methods ---
// GetSpellShieldDuration returns the duration of CC immunity.
func (q *QuicksilverEffect) GetSpellShieldDuration() float64 {
	return q.spellShieldDuration
}

// GetProcInterval returns the interval between AS procs.
func (q *QuicksilverEffect) GetProcInterval() float64 {
	return q.procInterval
}

// GetProcAttackSpeed returns the bonus AS gained per proc.
func (q *QuicksilverEffect) GetProcAttackSpeed() float64 {
	return q.procAttackSpeed
}

// GetCurrentBonusAS returns the total accumulated bonus AS from procs.
func (q *QuicksilverEffect) GetCurrentBonusAS() float64 {
	return q.currentBonusAS
}

// AddBonusAS increases the accumulated bonus AS.
func (q *QuicksilverEffect) AddBonusAS(deltaAS float64) {
	q.currentBonusAS += deltaAS
}

// ResetBonusAS resets the accumulated bonus AS (e.g., if effect ends).
func (q *QuicksilverEffect) ResetBonusAS() {
	q.currentBonusAS = 0.0
}

func (q *QuicksilverEffect) GetStacks() int {
	return q.stacks
}

func (q *QuicksilverEffect) SetStacks(stacks int) {
	q.stacks = stacks
}

func (q *QuicksilverEffect) AddStacks(deltaStacks int) {
	q.stacks += deltaStacks
}

func (q *QuicksilverEffect) ResetEffects() {
	q.isActive = true

	// Attack Speed Proc State (during immunity)
	q.stacks = 0
	q.currentBonusAS = 0.0
}

// TODO: Consider adding an IsImmuneToCC marker component if CC mechanics are implemented.
// type IsImmuneToCC struct{}
