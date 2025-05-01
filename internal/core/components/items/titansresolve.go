package items

/*
 * "desc": "Gain @StackingAD*100@% Attack Damage and @StackingSP@ Ability Power when attacking or taking damage, stacking up to @StackCap@ times.  <br><br>At full stacks, gain @BonusResistsAtStackCap@ Armor and @BonusResistsAtStackCap@ Magic Resist."
 */
// TitansResolveEffect holds the state for the Titan's Resolve item.
type TitansResolveEffect struct {
    currentStacks int
    maxStacks     int
    adPerStack    float64 // Stored as decimal, e.g., 0.02 for 2%
    apPerStack    float64
    bonusArmor  float64 // Armor and MR granted at max stacks
	bonusMR    float64 // MR granted at max stacks (if applicable)
    isMaxStacks   bool    // Flag to track if max stacks reached
	bonusResistsApplied bool // Flag to track if bonus resists were applied
}

// NewTitansResolveEffect creates a new TitansResolveEffect component.
func NewTitansResolveEffect(maxStacks float64, adPerStack, apPerStack, bonusResists float64) *TitansResolveEffect {
    return &TitansResolveEffect{
        currentStacks: 0,
        maxStacks:     int(maxStacks),
        adPerStack:    adPerStack,
        apPerStack:    apPerStack,
        bonusArmor:  bonusResists,
		bonusMR:   bonusResists, 
        isMaxStacks:   false,
		bonusResistsApplied: false,
    }
}

// IncrementStacks increments the stack count by 1 each time, up to the maximum.
// Returns true if the stack count changed, false otherwise.
// Also returns true if max stacks were reached *this time*.
func (t *TitansResolveEffect) IncrementStacks() (stackAdded bool, reachedMax bool) {
    if t.currentStacks < t.maxStacks {
        t.currentStacks ++
        stackAdded = true
        if t.currentStacks == t.maxStacks && !t.isMaxStacks {
            t.isMaxStacks = true
            reachedMax = true
        }
        return stackAdded, reachedMax
    }
    return false, false
}

// GetCurrentBonusAD returns the bonus AD based on current stacks.
func (t *TitansResolveEffect) GetCurrentBonusAD() float64 {
    return float64(t.currentStacks) * t.adPerStack
}

// GetCurrentBonusAP returns the bonus AP based on current stacks.
func (t *TitansResolveEffect) GetCurrentBonusAP() float64 {
    return float64(t.currentStacks) * t.apPerStack
}

// GetCurrentBonusArmor returns the bonus Armor if at max stacks.
func (t *TitansResolveEffect) GetBonusArmorAtMax() float64 {
    if t.isMaxStacks {
        return t.bonusArmor
    }
    return 0.0
}

// GetCurrentBonusMR returns the bonus Magic Resist if at max stacks.
func (t *TitansResolveEffect) GetBonusMRAtMax() float64 {
	if t.isMaxStacks {
		return t.bonusMR
	}
	return 0.0
}

// Getters for config values if needed...
func (t *TitansResolveEffect) GetADPerStack() float64 {
    return t.adPerStack
}

func (t *TitansResolveEffect) GetAPPerStack() float64 {
    return t.apPerStack
}

func (t *TitansResolveEffect) GetCurrentStacks() int {
	return t.currentStacks
}

func (t *TitansResolveEffect) GetMaxStacks() int {
	return t.maxStacks
}

func (t *TitansResolveEffect) IsMaxStacksReached() bool {
	return t.isMaxStacks
}

func (t *TitansResolveEffect) IsBonusResistsApplied() bool {
	return t.bonusResistsApplied
}

func (t *TitansResolveEffect) ResetStacks() {
	t.currentStacks = 0
	t.isMaxStacks = false
	t.bonusResistsApplied = false
}