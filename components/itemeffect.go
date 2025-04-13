package components

// ItemEffect holds the aggregated passive stat bonuses from all items.
type ItemEffect struct {
	bonusHealth             float64
	bonusPercentHp          float64 // Represented as a percentage, e.g., 0.1 for +10%
	bonusInitialMana        float64
	bonusArmor              float64
	bonusMR                 float64
	bonusPercentAD          float64 // Represented as a percentage
	bonusDamgeAmp                float64 // Represented as a percentage, e.g., 0.1 for +10%
	bonusAP                 float64
	bonusPercentAttackSpeed float64 // Represented as a multiplier bonus, e.g., 0.1 for +10%
	bonusCritChance         float64
	bonusCritDamage              float64 // Represented as a multiplier bonus, e.g., 0.1 for +10%
	durability              float64 // percent
	// Add other stats as needed (MoveSpeed, Range, Omnivamp, etc.)

}

func NewItemEffect() *ItemEffect {
	return &ItemEffect{} // Initialize with zero values
}

// Add other item-specific components here later if needed for effects,
// e.g., HasGuinsoosRageblade, HasStatikkShiv, etc.
// type HasTearOfTheGoddess struct{} // Example marker component

func (ie *ItemEffect) ResetStats() {
	ie.bonusHealth = 0
	ie.bonusPercentHp = 0
	ie.bonusInitialMana = 0
	ie.bonusArmor = 0
	ie.bonusMR = 0
	ie.bonusPercentAD = 0
	ie.bonusDamgeAmp = 0
	ie.bonusAP = 0
	ie.bonusPercentAttackSpeed = 0
	ie.bonusCritChance = 0
	ie.bonusCritDamage = 0
	ie.durability = 0
	// Reset other stats as needed...
}

// Methods to modify the stats
func (ie *ItemEffect) AddBonusHealth(amount float64) {
	ie.bonusHealth += amount
}

func (ie *ItemEffect) AddBonusPercentHp(amount float64) {
	ie.bonusPercentHp += amount
}

func (ie *ItemEffect) AddBonusInitialMana(amount float64) {
	ie.bonusInitialMana += amount
}

func (ie *ItemEffect) AddBonusArmor(amount float64) {
	ie.bonusArmor += amount
}

func (ie *ItemEffect) AddBonusMR(amount float64) {
	ie.bonusMR += amount
}

func (ie *ItemEffect) AddBonusPercentAD(amount float64) {
	ie.bonusPercentAD += amount
}

func (ie *ItemEffect) AddBonusDamageAmp(amount float64) {
	ie.bonusDamgeAmp += amount
}

func (ie *ItemEffect) AddBonusAP(amount float64) {
	ie.bonusAP += amount
}

func (ie *ItemEffect) AddBonusPercentAttackSpeed(amount float64) {
	ie.bonusPercentAttackSpeed += amount
}

func (ie *ItemEffect) AddBonusCritChance(amount float64) {
	ie.bonusCritChance += amount
}

func (ie *ItemEffect) AddDurability(amount float64) {
	ie.durability += amount
}

func (ie *ItemEffect) AddCritDamage(amount float64) {
	ie.bonusCritDamage += amount
}

func (ie *ItemEffect) SetBonusHealth(amount float64) {
	ie.bonusHealth = amount
}

func (ie *ItemEffect) SetBonusInitialMana(amount float64) {
	ie.bonusInitialMana = amount
}

func (ie *ItemEffect) SetBonusArmor(amount float64) {
	ie.bonusArmor = amount
}

func (ie *ItemEffect) SetBonusMR(amount float64) {
	ie.bonusMR = amount
}

func (ie *ItemEffect) SetBonusPercentAD(amount float64) {
	ie.bonusPercentAD = amount
}

func (ie *ItemEffect) SetDamageAmp(amount float64) {
	ie.bonusDamgeAmp = amount
}

func (ie *ItemEffect) SetBonusAP(amount float64) {
	ie.bonusAP = amount
}

func (ie *ItemEffect) SetBonusPercentAttackSpeed(amount float64) {
	ie.bonusPercentAttackSpeed = amount
}

func (ie *ItemEffect) SetBonusCritChance(amount float64) {
	ie.bonusCritChance = amount
}

func (ie *ItemEffect) SetCritDamage(amount float64) {
	ie.bonusCritDamage = amount
}

func (ie *ItemEffect) SetDurability(amount float64) {
	ie.durability = amount
}

func (ie *ItemEffect) GetBonusHealth() float64 {
	return ie.bonusHealth
}

func (ie *ItemEffect) GetBonusPercentHp() float64 {
	return ie.bonusPercentHp
}

func (ie *ItemEffect) GetBonusInitialMana() float64 {
	return ie.bonusInitialMana
}

func (ie *ItemEffect) GetBonusArmor() float64 {
	return ie.bonusArmor
}

func (ie *ItemEffect) GetBonusMR() float64 {
	return ie.bonusMR
}

func (ie *ItemEffect) GetBonusPercentAD() float64 {
	return ie.bonusPercentAD
}

func (ie *ItemEffect) GetDamageAmp() float64 {
	return ie.bonusDamgeAmp
}

func (ie *ItemEffect) GetBonusAP() float64 {
	return ie.bonusAP
}

func (ie *ItemEffect) GetBonusPercentAttackSpeed() float64 {
	return ie.bonusPercentAttackSpeed
}

func (ie *ItemEffect) GetBonusCritChance() float64 {
	return ie.bonusCritChance
}

func (ie *ItemEffect) GetBonusCritDamage() float64 {
	return ie.bonusCritDamage
}

func (ie *ItemEffect) GetDurability() float64 {
	return ie.durability
}
