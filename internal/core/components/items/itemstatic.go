package items

import (
	"math"
)

// ItemStaticEffect holds the aggregated passive stat bonuses from all items.
type ItemStaticEffect struct {
	bonusHealth             float64
	bonusPercentHp          float64 // Represented as a percentage, e.g., 0.1 for +10%
	bonusInitialMana        float64
	bonusArmor              float64
	bonusMR                 float64
	bonusPercentAD          float64 // Represented as a percentage
	bonusDamgeAmp           float64 // Represented as a percentage, e.g., 0.1 for +10%
	bonusAP                 float64
	bonusPercentAttackSpeed float64 // Represented as a multiplier bonus, e.g., 0.1 for +10%
	bonusCritChance         float64
	bonusCritDamage         float64 // Represented as a multiplier bonus, e.g., 0.1 for +10%
	durability              float64 // percent
	// Add other stats as needed (MoveSpeed, Range, Omnivamp, etc.)
	critDamangeToGive float64 // Specific to Infity Edge and Jeweled Gauntlet
}

func NewItemStaticEffect() *ItemStaticEffect {
	return &ItemStaticEffect{} // Initialize with zero values
}

// Add other item-specific components here later if needed for effects,
// e.g., HasGuinsoosRageblade, HasStatikkShiv, etc.
// type HasTearOfTheGoddess struct{} // Example marker component

func (ie *ItemStaticEffect) ResetStats() {
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
	ie.critDamangeToGive = 0 // Reset to zero
	// Reset other stats as needed...
}

// Methods to modify the stats
func (ie *ItemStaticEffect) AddBonusHealth(amount float64) {
	ie.bonusHealth += amount
}

func (ie *ItemStaticEffect) AddBonusPercentHp(amount float64) {
	ie.bonusPercentHp += amount
}

func (ie *ItemStaticEffect) AddBonusInitialMana(amount float64) {
	ie.bonusInitialMana += amount
}

func (ie *ItemStaticEffect) AddBonusArmor(amount float64) {
	ie.bonusArmor += amount
}

func (ie *ItemStaticEffect) AddBonusMR(amount float64) {
	ie.bonusMR += amount
}

func (ie *ItemStaticEffect) AddBonusPercentAD(amount float64) {
	ie.bonusPercentAD += amount
}

func (ie *ItemStaticEffect) AddBonusDamageAmp(amount float64) {
	ie.bonusDamgeAmp += amount
}

func (ie *ItemStaticEffect) AddBonusAP(amount float64) {
	ie.bonusAP += amount
}

func (ie *ItemStaticEffect) AddBonusPercentAttackSpeed(amount float64) {
	ie.bonusPercentAttackSpeed += amount
}

func (ie *ItemStaticEffect) AddBonusCritChance(amount float64) {
	ie.bonusCritChance += amount
}

func (ie *ItemStaticEffect) AddDurability(amount float64) {
	ie.durability += amount
}

func (ie *ItemStaticEffect) AddCritDamage(amount float64) {
	ie.bonusCritDamage += amount
}

func (ie *ItemStaticEffect) AddCritDamageToGive(amount float64) {
	ie.critDamangeToGive += amount
}

func (ie *ItemStaticEffect) SetBonusHealth(amount float64) {
	ie.bonusHealth = amount
}

func (ie *ItemStaticEffect) SetBonusInitialMana(amount float64) {
	ie.bonusInitialMana = amount
}

func (ie *ItemStaticEffect) SetBonusArmor(amount float64) {
	ie.bonusArmor = amount
}

func (ie *ItemStaticEffect) SetBonusMR(amount float64) {
	ie.bonusMR = amount
}

func (ie *ItemStaticEffect) SetBonusPercentAD(amount float64) {
	ie.bonusPercentAD = amount
}

func (ie *ItemStaticEffect) SetDamageAmp(amount float64) {
	ie.bonusDamgeAmp = amount
}

func (ie *ItemStaticEffect) SetBonusAP(amount float64) {
	ie.bonusAP = amount
}

func (ie *ItemStaticEffect) SetBonusPercentAttackSpeed(amount float64) {
	ie.bonusPercentAttackSpeed = amount
}

func (ie *ItemStaticEffect) SetBonusCritChance(amount float64) {
	ie.bonusCritChance = amount
}

func (ie *ItemStaticEffect) SetCritDamage(amount float64) {
	ie.bonusCritDamage = amount
}

func (ie *ItemStaticEffect) SetDurability(amount float64) {
	ie.durability = amount
}

func (ie *ItemStaticEffect) SetCritDamageToGive(amount float64) {
	ie.critDamangeToGive = amount
}

func (ie *ItemStaticEffect) GetBonusHealth() float64 {
	return ie.bonusHealth
}

func (ie *ItemStaticEffect) GetBonusPercentHp() float64 {
	return ie.bonusPercentHp
}

func (ie *ItemStaticEffect) GetBonusInitialMana() float64 {
	return ie.bonusInitialMana
}

func (ie *ItemStaticEffect) GetBonusArmor() float64 {
	return ie.bonusArmor
}

func (ie *ItemStaticEffect) GetBonusMR() float64 {
	return ie.bonusMR
}

func (ie *ItemStaticEffect) GetBonusPercentAD() float64 {
	return ie.bonusPercentAD
}

func (ie *ItemStaticEffect) GetDamageAmp() float64 {
	return ie.bonusDamgeAmp
}

func (ie *ItemStaticEffect) GetBonusAP() float64 {
	return ie.bonusAP
}

func (ie *ItemStaticEffect) GetBonusPercentAttackSpeed() float64 {
	return ie.bonusPercentAttackSpeed
}

func (ie *ItemStaticEffect) GetBonusCritChance() float64 {
	return ie.bonusCritChance
}

func (ie *ItemStaticEffect) GetBonusCritDamage() float64 {
	return ie.bonusCritDamage
}

func (ie *ItemStaticEffect) GetDurability() float64 {
	return ie.durability
}

func (ie *ItemStaticEffect) GetCritDamageToGive() float64 {
	if math.IsNaN(ie.critDamangeToGive) {
		return 0
	}
	return ie.critDamangeToGive
}
