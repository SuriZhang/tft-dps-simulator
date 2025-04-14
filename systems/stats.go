package systems

import (
	"log"
	"reflect"

	"github.com/suriz/tft-dps-simulator/components"
	"github.com/suriz/tft-dps-simulator/ecs"
)

// StatCalculationSystem calculates the final derived stats for entities
// based on their base stats and accumulated bonuses from items, traits, etc.
// It should run AFTER systems that apply bonuses (like BaseStaticItemSystem, TraitSystems).
type StatCalculationSystem struct {
	world *ecs.World
}

// NewStatCalculationSystem creates a new StatCalculationSystem.
func NewStatCalculationSystem(world *ecs.World) *StatCalculationSystem {
	return &StatCalculationSystem{world: world}
}

// Update calculates and sets the final stats for all relevant entities.
func (s *StatCalculationSystem) Update() {
	// Define component types needed. We need entities that have stats to calculate.
	// Querying for just one core stat component like Health might be sufficient,
	// as entities with stats usually have multiple stat components.
	healthType := reflect.TypeOf(components.Health{})
	attackType := reflect.TypeOf(components.Attack{})
	manaType := reflect.TypeOf(components.Mana{})

	entities := s.world.GetEntitiesWithComponents(healthType, manaType, attackType)

	for _, entity := range entities {
		// --- Calculate Final Stats ---
		s.calculateHealthStats(entity)
		s.calculateAttackStats(entity)
		s.calculateManaStats(entity)

		// --- Apply Consequential Logic ---
		s.applyHealthConsequences(entity)
		// Add other consequences if needed (e.g., mana adjustments)
	}
}

// calculateHealthStats calculates FinalMaxHP, FinalArmor, FinalMR, FinalDurability.
func (s *StatCalculationSystem) calculateHealthStats(entity ecs.Entity) {
	health, ok := s.world.GetHealth(entity)
	if !ok {
		return // No health component, nothing to calculate
	}

	// Max HP: (Base + FlatBonus) * (1 + PercentBonus)
	calculatedMaxHp := (health.GetBaseMaxHp() + health.GetBonusMaxHP()) * (1 + health.GetBonusPercentHp())
	health.SetFinalMaxHealth(calculatedMaxHp)

	// Armor: Base + FlatBonus (Add % bonus calculation if needed)
	calculatedArmor := health.GetBaseArmor() + health.GetBonusArmor()
	health.SetFinalArmor(calculatedArmor)

	// MR: Base + FlatBonus (Add % bonus calculation if needed)
	calculatedMR := health.GetBaseMR() + health.GetBonusMR()
	health.SetFinalMR(calculatedMR)

	// Durability: Base + FlatBonus (Assuming simple addition for now)
	// You might need a BaseDurability field if it exists.
	calculatedDurability := health.GetBonusDurability() // Or Base + Bonus
	health.SetFinalDurability(calculatedDurability)
}

// calculateAttackStats calculates FinalAD, FinalAS, FinalCritChance, FinalCritMultiplier, FinalDamageAmp.
func (s *StatCalculationSystem) calculateAttackStats(entity ecs.Entity) {
	attack, ok := s.world.GetAttack(entity)
	if !ok {
		return // No attack component
	}

	// AD: (Base + FlatBonus) * (1 + PercentBonus)
	// normally there's no flat bonus AD in TFT, but there are exceptions
	calculatedAD := (attack.GetBaseAD() + attack.GetBonusAD()) * (1 + attack.GetBonusPercentAD())
	attack.SetFinalAD(calculatedAD)

	// Attack Speed: BaseAS * (1 + TotalBonusAS%)
	calculatedAS := attack.GetBaseAttackSpeed() * (1 + attack.GetBonusPercentAttackSpeed())
	attack.SetFinalAttackSpeed(calculatedAS)

	// Crit Chance: Base + Bonus (Capped at 1.0)
	calculatedCritChance := attack.GetBaseCritChance() + attack.GetBonusCritChance()
	calculatedBonusCritDamange := 0.0
	if calculatedCritChance > 1.0 {
		attack.SetFinalCritChance(1.0)
		calculatedBonusCritDamange = (calculatedCritChance - 1.0) / 2
		log.Printf("Entity %d: Crit chance capped at 1.0 (was %.2f), extra crit chance(%.2f) is contributed to crit damange at 50%%.", entity, calculatedCritChance, calculatedBonusCritDamange)
	}
	attack.SetFinalCritChance(calculatedCritChance)

	// --- IE/JG Bonus Crit Damage Logic ---
	equipment, okEq := s.world.GetEquipment(entity)

	hasIE := false
	hasJG := false
	if okEq {
		_, hasIE = equipment.GetItem("TFT_Item_InfinityEdge")
		_, hasJG = equipment.GetItem("TFT_Item_JeweledGauntlet")
	}

	// TODO: check the case where a champion wears both IE and JG, it's a rare case but possible
	critDamageBonusFromItem := 0.0
	if hasIE || hasJG {
		traitCritMarkerType := reflect.TypeOf(components.CanAbilityCritFromTraits{})
		_, canAbilityCritFromTraits := s.world.GetComponent(entity, traitCritMarkerType)

		if canAbilityCritFromTraits {
			critDamageBonusFromItem = attack.GetBonusCritDamageToGive()
			log.Printf("Entity %d: Applying IE bonus crit damage (0.10) because abilities can already crit.", entity)
		}
	}

	// Crit Multiplier: Base + Bonus
	calculatedCritMultiplier := attack.GetBaseCritMultiplier() + attack.GetBonusCritMultiplier() + critDamageBonusFromItem + calculatedBonusCritDamange
	attack.SetFinalCritMultiplier(calculatedCritMultiplier)

	// Damage Amp: Base + Bonus (Assuming additive for now)
	calculatedDamageAmp := attack.GetBaseDamageAmp() + attack.GetBonusDamageAmp()
	attack.SetFinalDamageAmp(calculatedDamageAmp)

	// Range: Base + Bonus (Assuming simple addition)
	attack.SetFinalRange(attack.GetBaseRange() + attack.GetBonusRange())
}

// calculateManaStats calculates FinalInitialMana.
func (s *StatCalculationSystem) calculateManaStats(entity ecs.Entity) {
	mana, ok := s.world.GetMana(entity)
	if !ok {
		return // No mana component
	}

	// Initial Mana: Base + Bonus
	calculatedInitialMana := mana.GetBaseInitialMana() + mana.GetBonusInitialMana()
	mana.SetFinalInitialMana(calculatedInitialMana)

	// Potentially adjust CurrentMana based on InitialMana at combat start?
	// This might belong in a combat initialization system.
}

// applyHealthConsequences handles adjustments needed after health stats change.
func (s *StatCalculationSystem) applyHealthConsequences(entity ecs.Entity) {
	health, ok := s.world.GetHealth(entity)
	if !ok {
		return
	}

	// Adjust CurrentHP if it exceeds the new FinalMaxHP
	if health.GetCurrentHealth() > health.GetFinalMaxHP() {
		log.Printf("Entity %d: CurrentHP (%.2f) exceeds new FinalMaxHP (%.2f). Clamping.", entity, health.GetCurrentHealth(), health.GetFinalMaxHP())
		health.SetCurrentHealth(health.GetFinalMaxHP())
	}
	// Note: TFT generally doesn't automatically heal champions if their MaxHP increases mid-combat.
	// Resetting CurrentHP to MaxHP usually happens only at the start of combat.
}
