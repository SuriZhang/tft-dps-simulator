package items

// VoidStaffEffect holds the state for Void Staff's shred effect.
// "desc": "Attacks and Ability damage @MRShred@% <TFTKeyword>Shred</TFTKeyword> the target for @MRShredDuration@ seconds. This effect does not stack."
type VoidStaffEffect struct {
    mrShred    float64 // Magic resist reduction percentage (e.g., 30.0 for 30%)
    duration   float64 // Duration of the shred effect in seconds (e.g., 3.0)
}

// NewVoidStaffEffect creates a new VoidStaffEffect component.
func NewVoidStaffEffect(mrShred, duration float64) *VoidStaffEffect {
    return &VoidStaffEffect{
        mrShred:  mrShred,
        duration: duration,
    }
}

// GetMRShred returns the magic resist reduction amount.
func (vs *VoidStaffEffect) GetMRShred() float64 {
    return vs.mrShred
}

// GetDuration returns the shred duration.
func (vs *VoidStaffEffect) GetDuration() float64 {
    return vs.duration
}

// ResetEffects resets the effect state (if needed for future extensions).
func (vs *VoidStaffEffect) ResetEffects() {
    // Currently no internal state to reset, but keeping for consistency
}