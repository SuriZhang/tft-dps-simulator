package items

// SpiritVisageEffect holds the state for the Spirit Visage item's healing.
type SpiritVisageEffect struct {
    missingHealthHealRate float64
    tickInterval          float64 
	maxHeal 		  	float64 
}

// NewSpiritVisageEffect creates a new SpiritVisageEffect component.
func NewSpiritVisageEffect(missingHealthHealRate, tickInterval, maxHeal float64) *SpiritVisageEffect {
    return &SpiritVisageEffect{
        missingHealthHealRate: missingHealthHealRate,
        tickInterval:          tickInterval,
		maxHeal: 			maxHeal,
    }
}

func (sve *SpiritVisageEffect) GetMissingHealthHealRate() float64 {
    return sve.missingHealthHealRate
}

func (sve *SpiritVisageEffect) GetTickInterval() float64 {
    return sve.tickInterval
}

func (sve *SpiritVisageEffect) GetMaxHeal() float64 {
	return sve.maxHeal
}