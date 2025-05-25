package components

import (
	"fmt"
	"math"
	"strings"
)

// Health contains champion health information
type Health struct {
	// --- Base Stats (from champion data * star level) ---
	BaseMaxHP float64
	BaseArmor float64
	BaseMR    float64

	// --- Aggregated Bonus Stats (Sum from Items, Traits, Temp Buffs, etc.) ---
	// These are modified by systems like ItemSystem, TraitSystems, BuffSystems
	BonusMaxHP         float64 // Flat HP bonuses
	BonusPercentHp     float64 // Additive % HP bonuses (e.g., 0.1 + 0.2 = 0.3 for +30%)
	BonusArmor         float64 // Flat Armor bonuses
	BonusMR            float64 // Flat MR bonuses
	BonusDurability    float64 // Flat durability bonuses (e.g., from items or traits)
	// Add BonusPercentArmor/MR if needed by mechanics
	
	// --- Final Calculated Stats (Calculated by StatCalculationSystem) ---
	// These are the values used in combat calculations (damage reduction, checking max health)
	FinalMaxHP      float64
	FinalArmor      float64
	FinalMR         float64
	FinalDurability float64 // This is the durability that can be set, not the current durability
	
	HealReduction float64 // Percentage healing reduction (e.g., from Wound Debuff)
	// --- Current State ---
	CurrentHP float64
}

// NewHealth initializes the component with base stats.
func NewHealth(baseHp, baseArmor, baseMr float64) *Health {
	if baseHp < 0 || math.IsNaN(baseHp) {
		baseHp = 0
	}
	if baseArmor < 0 || math.IsNaN(baseArmor) {
		baseArmor = 0
	}
	if baseMr < 0 || math.IsNaN(baseMr) {
		baseMr = 0
	}

	return &Health{
		BaseMaxHP: baseHp,
		BaseArmor: baseArmor,
		BaseMR:    baseMr,

		// Initialize bonuses to 0
		BonusMaxHP:         0,
		BonusPercentHp:     0,
		BonusArmor:         0,
		BonusMR:            0,
		
		// Initialize final stats to base stats initially
		FinalMaxHP: baseHp,
		FinalArmor: baseArmor,
		FinalMR:    baseMr,
		HealReduction: 0,

		// Initialize current health
		CurrentHP: baseHp,
	}
}

func (h *Health) SetBaseMaxHP(maxHealth float64) {
	h.BaseMaxHP = maxHealth
}

func (h *Health) SetCurrentHP(currentHealth float64) {
	h.CurrentHP = currentHealth
}

// --- Methods to SET FINAL calculated stats (called by StatCalculationSystem) ---
func (h *Health) SetFinalMaxHP(value float64) {
	h.FinalMaxHP = value
}

func (h *Health) SetBaseArmor(armor float64) {
	h.BaseArmor = armor
}

func (h *Health) SetFinalArmor(armor float64) {
	h.FinalArmor = armor
}

func (h *Health) SetBaseMR(mr float64) {
	h.BaseMR = mr
}

func (h *Health) SetFinalMR(mr float64) {
	h.FinalMR = mr
}

func (h *Health) SetFinalDurability(durability float64) {
	h.FinalDurability = durability
}

// --- Methods to ADD to BONUS fields
func (h *Health) AddBonusMaxHealth(amount float64) {
	h.BonusMaxHP += amount
}
func (h *Health) AddBonusPercentHealth(amount float64) {
	h.BonusPercentHp += amount
}
func (h *Health) AddBonusArmor(amount float64) {
	h.BonusArmor += amount
}
func (h *Health) AddBonusMR(amount float64) {
	h.BonusMR += amount
}
func (h *Health) AddBonusDurability(amount float64) {
	h.BonusDurability += amount
}
func (h *Health) SetHealReduction(amount float64) {
	h.HealReduction = amount
}

func (h *Health) ResetBonuses() {
	h.BonusMaxHP = 0
	h.BonusPercentHp = 0
	h.BonusArmor = 0
	h.BonusMR = 0
	h.BonusDurability = 0
	// Reset other bonuses if added
}

func (h *Health) ResetHealth() {
	h.CurrentHP = 0
}

func (h *Health) AddHealth(amount float64) {
	h.CurrentHP += amount
	if h.CurrentHP > h.BaseMaxHP {
		h.CurrentHP = h.BaseMaxHP
	}
}

func (h *Health) AddBaseMaxHp(amount float64) float64 {
	h.BaseMaxHP += amount
	if h.CurrentHP > h.BaseMaxHP {
		h.CurrentHP = h.BaseMaxHP
	}
	return h.BaseMaxHP
}

func (h *Health) AddFinalMaxHP(amount float64) {
	h.FinalMaxHP += amount
	if h.CurrentHP > h.FinalMaxHP {
		h.CurrentHP = h.FinalMaxHP
	}
}

func (h *Health) AddBaseArmor(amount float64) {
	h.BaseArmor += amount
}

func (h *Health) AddDurability(amount float64) {
	h.BonusDurability += amount
}

// Heal adjusts current HP, respecting the already calculated FinalMaxHP
func (h *Health) Heal(amount float64) {
	h.CurrentHP += (1 - h.HealReduction) * amount
	if h.CurrentHP > h.FinalMaxHP { // Use FinalMaxHP here
		h.CurrentHP = h.FinalMaxHP
	}
}

// ResetCurrentHealth sets CurrentHP to the calculated FinalMaxHP (e.g., at combat start)
func (h *Health) ResetCurrentHealth() {
	h.CurrentHP = h.FinalMaxHP // Use FinalMaxHP
}

func (h *Health) IsAlive() bool {
	return h.CurrentHP > 0
}

func (h *Health) IsDead() bool {
	return h.CurrentHP <= 0
}

func (h *Health) GetBaseMaxHp() float64 {
	return h.BaseMaxHP
}

func (h *Health) GetBaseArmor() float64 {
	return h.BaseArmor
}

func (h *Health) GetBaseMR() float64 {
	return h.BaseMR
}

func (h *Health) GetFinalArmor() float64 {
	return h.FinalArmor
}

func (h *Health) GetFinalMR() float64 {
	return h.FinalMR
}

func (h *Health) GetFinalMaxHP() float64 {
	return h.FinalMaxHP
}

func (h *Health) GetBonusDurability() float64 {
	return h.BonusDurability
}

func (h *Health) GetFinalDurability() float64 {
	return h.FinalDurability
}

func (h *Health) GetBonusMaxHP() float64 {
	return h.BonusMaxHP
}

func (h *Health) GetBonusPercentHp() float64 {
	return h.BonusPercentHp
}

func (h *Health) GetBonusArmor() float64 {
	return h.BonusArmor
}

func (h *Health) GetBonusMR() float64 {
	return h.BonusMR
}

func (h *Health) GetCurrentHP() float64 {
	return h.CurrentHP
}

// String returns a multi-line string representation of the Health component.
func (h *Health) String() string {
	var sb strings.Builder // Use strings.Builder for efficiency

	// HP Line
	sb.WriteString(fmt.Sprintf("  HP: %.2f / %.2f (Base: %.2f, BonusFlat: %.2f, BonusPercent: %.2f)\n",
		h.CurrentHP, h.FinalMaxHP, h.BaseMaxHP, h.BonusMaxHP, h.BonusPercentHp))

	// Armor Line
	sb.WriteString(fmt.Sprintf("  Armor: %.2f (Base: %.2f, Bonus: %.2f)\n",
		h.FinalArmor, h.BaseArmor, h.BonusArmor))

	// MR Line
	sb.WriteString(fmt.Sprintf("  MR: %.2f (Base: %.2f, Bonus: %.2f)\n",
		h.FinalMR, h.BaseMR, h.BonusMR))

	// Durability Line
	sb.WriteString(fmt.Sprintf("  Durability: %.2f (Bonus: %.2f)",
		h.FinalDurability, h.BonusDurability))

	return sb.String()
}
