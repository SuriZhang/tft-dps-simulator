package items

// SpearOfShojinEffect holds the state for the Spear of Shojin item.
/**
 * "desc": "Attacks grant @FlatManaRestore@ bonus Mana."
 */
type SpearOfShojinEffect struct {
    flatManaRestore float64 // Mana restored per attack
}

// NewSpearOfShojinEffect creates a new SpearOfShojinEffect component.
func NewSpearOfShojinEffect(flatManaRestore float64) *SpearOfShojinEffect {
    return &SpearOfShojinEffect{
        flatManaRestore: flatManaRestore,
    }
}

// GetFlatManaRestore returns the mana restored per attack.
func (s *SpearOfShojinEffect) GetFlatManaRestore() float64 {
    return s.flatManaRestore
}

// SetFlatManaRestore sets the mana restored per attack.
func (s *SpearOfShojinEffect) SetFlatManaRestore(value float64) {
    s.flatManaRestore = value
}

// ResetEffects resets the effect (no persistent state for this item).
func (s *SpearOfShojinEffect) ResetEffects() {
    // No persistent state to reset for Spear of Shojin
}