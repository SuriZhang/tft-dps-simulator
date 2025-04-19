package components

import (
	"fmt"
	"math"
	"strings"
)

// Attack contains champion attack information
type Attack struct {
	// --- Base Stats (from champion data * star level) ---
	BaseAD          float64
	BaseAttackSpeed float64 // Champion's inherent AS value
	BaseDamageAmp   float64 // Usually 0.0
	BaseRange       float64 // Range is often less modified, keep simple for now

	// --- Aggregated Bonus Stats (Sum from Items, Traits, Temp Buffs, etc.) ---
	BonusAD                 float64 // Flat AD bonuses
	BonusPercentAD          float64 // Additive % AD bonuses (e.g., from Deathblade)
	BonusPercentAttackSpeed float64 // Additive % AS bonuses (e.g., 0.1 + 0.3 = 0.4 for +40%)
	BonusDamageAmp          float64 // Additive Damage Amp bonuses
	BonusRange              float64 // Flat Range bonuses (e.g., from items or traits)

	// --- Final Calculated Stats (Calculated by StatCalculationSystem) ---
	FinalAD          float64
	FinalAttackSpeed float64 // Calculated: BaseAS * (1 + TotalBonusAS%)
	FinalDamageAmp   float64 // Calculated: BaseDamageAmp + BonusDamageAmp (or multiplicative?)
	FinalRange       float64 // Calculated: BaseRange + BonusRange

	// --- Current State ---
	LastAttackTime float64 // For tracking attack cooldown
}

// NewAttack creates an Attack component
func NewAttack(baseAd, baseAs, baseRange float64) *Attack {
	if baseAd < 0 || math.IsNaN(baseAd) {
		baseAd = 0
	}
	if baseAs < 0 || math.IsNaN(baseAs) {
		baseAs = 0
	}
	if baseRange < 0 || math.IsNaN(baseRange) {
		baseRange = 0
	}

	return &Attack{
		// Base Stats
		BaseAD:          baseAd,
		BaseAttackSpeed: baseAs,
		BaseRange:       baseRange,
		BaseDamageAmp:   0.0,

		// Bonus Stats (Initialize to 0)
		BonusAD:                 0.0,
		BonusPercentAD:          0.0,
		BonusPercentAttackSpeed: 0.0,
		BonusDamageAmp:          0.0,
		BonusRange:              0.0,

		// Final Stats (Initialize to Base initially)
		FinalAD:          baseAd,
		FinalAttackSpeed: baseAs,
		FinalDamageAmp:   0.0,
		FinalRange:       baseRange,

		// State
		LastAttackTime: 0.0,
	}
}

// --- Methods to ADD to BONUS fields (called by ItemSystem, TraitSystems, etc.) ---
func (a *Attack) AddBonusAD(amount float64) {
	a.BonusAD += amount
}
func (a *Attack) AddBonusPercentAD(amount float64) {
	a.BonusPercentAD += amount
}
func (a *Attack) AddBonusPercentAttackSpeed(amount float64) {
	a.BonusPercentAttackSpeed += amount
}
func (a *Attack) AddBonusDamageAmp(amount float64) {
	a.BonusDamageAmp += amount
}

func (a *Attack) AddBonusRange(amount float64) {
	a.BonusRange += amount
}

// ResetBonuses resets all bonus stats to 0.0
func (a *Attack) ResetBonuses() {
	a.BonusAD = 0.0
	a.BonusPercentAD = 0.0
	a.BonusPercentAttackSpeed = 0.0
	a.BonusDamageAmp = 0.0
	a.BonusRange = 0.0
}

// --- Methods to SET FINAL calculated stats (called by StatCalculationSystem) ---
// These just perform the assignment.
func (a *Attack) SetFinalAD(value float64) {
	a.FinalAD = value
}
func (a *Attack) SetFinalAttackSpeed(value float64) {
	a.FinalAttackSpeed = value
}
func (a *Attack) SetFinalDamageAmp(value float64) {
	a.FinalDamageAmp = value
}

func (a *Attack) SetFinalRange(value float64) {
	a.FinalRange = value
}
func (a *Attack) SetBaseAttackSpeed(value float64) {
	a.BaseAttackSpeed = value
}
func (a *Attack) SetBaseRange(value float64) {
	a.BaseRange = value
}

func (a *Attack) SetBaseAD(value float64) {
	a.BaseAD = value
}

func (a *Attack) SetBaseDamageAmp(value float64) {
	a.BaseDamageAmp = value
}

func (a *Attack) SetBonusPercentAttackSpeed(value float64) {
	a.BonusPercentAttackSpeed = value
}

// --- Methods to GET FINAL stats (used by combat systems) ---
func (a *Attack) GetFinalAD() float64 {
	return a.FinalAD
}
func (a *Attack) GetFinalAttackSpeed() float64 {
	return a.FinalAttackSpeed
}

func (a *Attack) GetFinalDamageAmp() float64 {
	return a.FinalDamageAmp
}
func (a *Attack) GetFinalRange() float64 {
	return a.FinalRange
} // Keep simple for now

// --- Methods for Current State ---
func (a *Attack) GetLastAttackTime() float64 {
	return a.LastAttackTime
}
func (a *Attack) SetLastAttackTime(value float64) {
	a.LastAttackTime = value
}

// --- Getters for Base/Bonus stats (Optional: for debugging/systems) ---
func (a *Attack) GetBaseAD() float64 {
	return a.BaseAD
}

func (a *Attack) GetBonusPercentAttackSpeed() float64 {
	return a.BonusPercentAttackSpeed
}

func (a *Attack) GetBaseDamageAmp() float64 {
	return a.BaseDamageAmp
}
func (a *Attack) GetBonusDamageAmp() float64 {
	return a.BonusDamageAmp
}

func (a *Attack) GetBonusAD() float64 {
	return a.BonusAD
}
func (a *Attack) GetBonusPercentAD() float64 {
	return a.BonusPercentAD
}
func (a *Attack) GetBonusRange() float64 {
	return a.BonusRange
}

func (a *Attack) GetBaseAttackSpeed() float64 {
	return a.BaseAttackSpeed
}
func (a *Attack) GetBaseRange() float64 {
	return a.BaseRange
}

// String returns a multi-line string representation of the Attack component.
func (a *Attack) String() string {
	var sb strings.Builder // Use strings.Builder for efficiency

	sb.WriteString(fmt.Sprintf("  BaseAD: %.2f, BonusAD: %.2f, BonusPercentAD: %.2f, FinalAD: %.2f\n", a.BaseAD, a.BonusAD, a.BonusPercentAD, a.FinalAD))
	sb.WriteString(fmt.Sprintf("  BaseAS: %.3f, BonusASPercent: %.2f, FinalAS: %.3f\n", a.BaseAttackSpeed, a.BonusPercentAttackSpeed, a.FinalAttackSpeed))
	sb.WriteString(fmt.Sprintf("  BaseDamageAmp: %.2f, BonusDamageAmp: %.2f, FinalDamageAmp: %.2f\n", a.BaseDamageAmp, a.BonusDamageAmp, a.FinalDamageAmp))
	sb.WriteString(fmt.Sprintf("  Range: %.2f, BonusRange: %.2f, FinalRange: %.2f\n", a.BaseRange, a.BonusRange, a.FinalRange))
	sb.WriteString(fmt.Sprintf("  LastAttackTime: %.2f", a.LastAttackTime))

	return sb.String()
}
