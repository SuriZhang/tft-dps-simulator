package items

// BlueBuffEffect represents the Blue Buff item component.
// Dynamic effects:
// 1. Gain 10 Mana after casting (on SpellCastCycleStartEvent)
// 2. Deal 5% more damage for 8 seconds after getting a takedown (KillEvent or AssistEvent)
type BlueBuffEffect struct {
	ManaRefund    float64 // 10.0 - Mana gained after casting
	DamageAmp     float64 // 0.05 - Damage amplification (5%)
	TakedownTimer float64 // 8.0 - Duration of damage amplification

	// State tracking for damage amplification buff
	IsAmplified               bool    // Whether damage amp is currently active
	AmplificationEndTime      float64 // When the damage amp expires
	CurrentActivationSequence uint64  // To handle timer resets
}

// NewBlueBuff creates a new BlueBuff component with default values.
func NewBlueBuff() *BlueBuffEffect {
	return &BlueBuffEffect{
		ManaRefund:                10.0,
		DamageAmp:                 0.05, // 5%
		TakedownTimer:             8.0,
		IsAmplified:               false,
		AmplificationEndTime:      0.0,
		CurrentActivationSequence: 0,
	}
}

// IsActive returns true if the damage amplification effect is currently active.
func (bb *BlueBuffEffect) IsActive(currentTime float64) bool {
	return bb.IsAmplified && currentTime < bb.AmplificationEndTime
}

// ActivateAmplification activates the damage amplification effect.
func (bb *BlueBuffEffect) ActivateAmplification(currentTime float64) {
	bb.IsAmplified = true
	bb.AmplificationEndTime = currentTime + bb.TakedownTimer
}

// DeactivateAmplification deactivates the damage amplification effect.
func (bb *BlueBuffEffect) DeactivateAmplification() {
	bb.IsAmplified = false
	bb.AmplificationEndTime = 0.0
}

func (bb *BlueBuffEffect) GetManaRefund() float64 {
	return bb.ManaRefund
}

func (bb *BlueBuffEffect) IncremenetCurrentActivationSequence() {
	bb.CurrentActivationSequence++
}
