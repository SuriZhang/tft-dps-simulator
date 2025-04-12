package components

// Mana contains champion mana information
type Mana struct {
    Max      float64
    Current  float64
    StartMana float64
}

// NewMana creates a Mana component
func NewMana(max, start float64) Mana {
    return Mana{
        Max:      max,
        Current:  start,
        StartMana: start,
    }
}