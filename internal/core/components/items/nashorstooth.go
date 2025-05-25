package items

// NashorsToothEffect holds the state for Nashor's Tooth.
/**
 * "desc": "After casting an Ability, gain @AttackSpeedToGive@% Attack Speed for @ASDuration@ seconds."
 */
type NashorsToothEffect struct {
    bonusAttackSpeed float64 // Attack speed gained per proc (as decimal, e.g., 0.6 for 60%)
    duration        float64 // Duration of the buff in seconds (e.g., 5.0)
    
    // State tracking for the buff
    isActive    bool    // Whether the buff is currently active
    endTime     float64 // When the buff expires
    currentSequence uint64 // To handle overlapping activations
}

// NewNashorsToothEffect creates a new NashorsToothEffect component.
func NewNashorsToothEffect(attackSpeedToGive, duration float64) *NashorsToothEffect {
    return &NashorsToothEffect{
        bonusAttackSpeed: attackSpeedToGive,
        duration:        duration,
        isActive:        false,
        endTime:         0.0,
        currentSequence: 0,
    }
}

// IsActive returns true if the attack speed buff is currently active.
func (nt *NashorsToothEffect) IsActive(currentTime float64) bool {
    return nt.isActive && currentTime < nt.endTime
}

// ActivateBuff activates the attack speed buff.
func (nt *NashorsToothEffect) ActivateBuff(currentTime float64) {
    nt.isActive = true
    nt.endTime = currentTime + nt.duration
}

// DeactivateBuff deactivates the attack speed buff.
func (nt *NashorsToothEffect) DeactivateBuff() {
    nt.isActive = false
    nt.endTime = 0.0
}

// GetBonusAS returns the attack speed gain amount.
func (nt *NashorsToothEffect) GetBonusAS() float64 {
    return nt.bonusAttackSpeed
}

// GetDuration returns the buff duration.
func (nt *NashorsToothEffect) GetDuration() float64 {
    return nt.duration
}

// IncrementCurrentSequence increments the activation sequence number.
func (nt *NashorsToothEffect) IncrementCurrentSequence() {
    nt.currentSequence++
}

// GetCurrentSequence returns the current activation sequence number.
func (nt *NashorsToothEffect) GetCurrentSequence() uint64 {
    return nt.currentSequence
}

// ResetEffects resets the effect state.
func (nt *NashorsToothEffect) ResetEffects() {
    nt.isActive = false
    nt.endTime = 0.0
    nt.currentSequence = 0
}