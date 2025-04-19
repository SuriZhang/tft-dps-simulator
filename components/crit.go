package components

import (
	"fmt"
	"math"
	"strings"
)

// Crit holds shared critical strike statistics for an entity.
type Crit struct {
	// Base stats (Usually 0.25 chance, 1.4 multiplier for attacks, 0/1.4 for spells)
	// Note: Base spell crit is effectively 0 unless enabled by items/traits.
	// We'll handle the base multiplier application logic in StatCalculationSystem.
	BaseCritChance     float64
	BaseCritMultiplier float64 // e.g., 1.4 for basic attacks

	// Bonus stats from items, traits, etc.
	BonusCritChance       float64
	BonusCritMultiplier   float64 // Flat additions (e.g., JG's +0.4)
	BonusCritDamageToGive float64 // Specific to Infinity Edge and Jeweled Gauntlet 
	// Final calculated stats
	FinalCritChance     float64 // Capped at 1.0
	FinalCritMultiplier float64 // Includes base, bonus, excess chance conversion, IE/JG conditional
}

// NewCrit creates a Crit component with default attack values.
func NewCrit(critChance, critMultiplier float64) *Crit {
	if math.IsNaN(critChance) || critChance < 0 {
		critChance = 0.0
	}
	if math.IsNaN(critMultiplier) || critMultiplier < 0 {
		critMultiplier = 0.0
	}
	return &Crit{
		BaseCritChance:    critChance, 
		BaseCritMultiplier: critMultiplier,
		
		BonusCritChance:    0.0, // No bonus crit chance by default
		BonusCritMultiplier: 0.0, // No bonus crit multiplier by default
		BonusCritDamageToGive: 0.0, // No bonus crit damage by default
		
		// Finals default to base values
		FinalCritChance:    critChance, // Start with base crit chance
		FinalCritMultiplier: critMultiplier, // Start with base crit multiplier

	}
}

// --- Getters for Base ---
func (c *Crit) GetBaseCritChance() float64     { return c.BaseCritChance }
func (c *Crit) GetBaseCritMultiplier() float64 { return c.BaseCritMultiplier }

// --- Setters/Adders for Bonus ---
func (c *Crit) AddBonusCritChance(amount float64) {
	c.BonusCritChance += amount
}
func (c *Crit) AddBonusCritMultiplier(amount float64) {
	c.BonusCritMultiplier += amount
}
func (c *Crit) AddBonusCritDamageToGive(amount float64) {
	if math.IsNaN(amount) || amount < 0 {
		amount = 0.0
	}
	c.BonusCritDamageToGive += amount
}
func (c *Crit) ResetBonuses() {
	c.BonusCritChance = 0.0
	c.BonusCritMultiplier = 0.0
	c.BonusCritDamageToGive = 0.0
}

// --- Getters for Bonus ---
func (c *Crit) GetBonusCritChance() float64 {
	return c.BonusCritChance
}
func (c *Crit) GetBonusCritMultiplier() float64 {
	return c.BonusCritMultiplier
}
func (c *Crit) GetBonusCritDamageToGive() float64 {
	if math.IsNaN(c.BonusCritDamageToGive) || c.BonusCritDamageToGive < 0 {
		c.BonusCritDamageToGive = 0.0
	}
	return c.BonusCritDamageToGive
}

// --- Setters for Final ---
func (c *Crit) SetFinalCritChance(value float64) {
	c.FinalCritChance = value
}
func (c *Crit) SetFinalCritMultiplier(value float64) {
	c.FinalCritMultiplier = value
}

// --- Getters for Final ---
func (c *Crit) GetFinalCritChance() float64 {
	return c.FinalCritChance
}
func (c *Crit) GetFinalCritMultiplier() float64 {
	return c.FinalCritMultiplier
}

func (c *Crit) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("  BaseCritChance: %.2f, BonusCritChance: %.2f, FinalCritChance: %.2f\n", c.BaseCritChance, c.BonusCritChance, c.FinalCritChance))
	sb.WriteString(fmt.Sprintf("  BaseCritMulti: %.2f, BonusCritMulti: %.2f, *BonusCritDamageToGive: %.2f, FinalCritMulti: %.2f\n", c.BaseCritMultiplier, c.BonusCritMultiplier, c.BonusCritDamageToGive, c.FinalCritMultiplier))

	return sb.String()
}