package components

// Health contains champion health information
type Health struct {
	Max     float64
	Current float64
	Armor   float64
	MR      float64
}

// NewHealth creates a Health component with current health set to max
func NewHealth(max, armor, mr float64) Health {
	return Health{
		Max:     max,
		Current: max,
		Armor:   armor,
		MR:      mr,
	}
}

func (h *Health) UpdateMaxHealth(maxHealth float64) {
	h.Max = maxHealth
}

func (h *Health) UpdateCurrentHealth(currentHealth float64) {
	h.Current = currentHealth
}

func (h *Health) UpdateArmor(armor float64) {
	h.Armor = armor
}

func (h *Health) UpdateMR(mr float64) {
	h.MR = mr
}
