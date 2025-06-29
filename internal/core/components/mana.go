package components

import (
	"math"
)

// Mana contains champion mana information
type Mana struct {
	Max              float64
	Current          float64
	BaseInitialMana  float64
	BonusInitialMana float64
	FinalInitialMana float64
}

// NewMana creates a Mana component
func NewMana(max, start float64) Mana {
	if max < 0 || math.IsNaN(max) {
		max = 0
	}
	if start < 0 || math.IsNaN(start) {
		start = 0
	}

	return Mana{
		Max:              max,
		Current:          start,
		BaseInitialMana:  start,
		BonusInitialMana: 0,
		FinalInitialMana: start, // when creating a new champion, finalInitialMana is set to the static max mana from data
	}
}

func (m *Mana) SetCurrentMana(currentMana float64) {
	m.Current = currentMana
}

func (m *Mana) SetBaseInitialMana(baseMana float64) {
	m.BaseInitialMana = baseMana
}

func (m *Mana) SetBonusInitialMana(bonusMana float64) {
	m.BonusInitialMana = bonusMana
}

func (m *Mana) SetFinalInitialMana(finalMana float64) {
	m.FinalInitialMana = finalMana
}

func (m *Mana) SetMaxMana(maxMana float64) {
	m.Max = maxMana
}

func (m *Mana) ResetBonuses() {
	m.BonusInitialMana = 0
}

func (m *Mana) ResetCurrentMana() {
	m.Current = 0
}

// AddCurrentMana allows exceeding max mana temporarily, useful for effects that grant mana
func (m *Mana) AddCurrentMana(amount float64) {
	m.Current += amount
}

func (m *Mana) AddBonusInitialMana(amount float64) {
	m.BonusInitialMana += amount
	if m.Current > m.Max {
		m.Current = m.Max
	}
}

func (m *Mana) GetCurrentMana() float64 {
	return m.Current
}

func (m *Mana) GetMaxMana() float64 {
	return m.Max
}

func (m *Mana) GetBaseInitialMana() float64 {
	return m.BaseInitialMana
}

func (m *Mana) GetBonusInitialMana() float64 {
	return m.BonusInitialMana
}

func (m *Mana) GetFinalInitialMana() float64 {
	return m.FinalInitialMana
}

func (m *Mana) IsFull() bool {
	return m.Current >= m.Max
}

func (m *Mana) CanCastSpell() bool {
	return m.Current >= m.Max && m.Max != 0.0
}