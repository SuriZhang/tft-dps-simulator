// filepath: internal/core/components/items/krakensfury.go
package items

// KrakensFuryEffect holds the state for Kraken's Fury's stacking AD.
type KrakensFuryEffect struct {
    currentStacks int
    adPerStack    float64 // Bonus %AD per stack (e.g., 0.03 for 3%)
}

// NewKrakensFuryEffect creates a new KrakensFuryEffect component.
func NewKrakensFuryEffect(adPerStack float64) *KrakensFuryEffect {
    return &KrakensFuryEffect{
        currentStacks: 0,
        adPerStack:    adPerStack,
    }
}

// IncrementStacks increases the stack count.
func (kf *KrakensFuryEffect) IncrementStacks() {
    kf.currentStacks++
}

// GetCurrentStacks returns the current number of stacks.
func (kf *KrakensFuryEffect) GetCurrentStacks() int {
    return kf.currentStacks
}

// GetADPerStack returns the bonus %AD gained per stack.
func (kf *KrakensFuryEffect) GetADPerStack() float64 {
    return kf.adPerStack
}

// ResetEffects resets the stacks (e.g., at the start of combat if needed, though typically stacks reset with new simulation).
func (kf *KrakensFuryEffect) ResetEffects() {
    kf.currentStacks = 0
}