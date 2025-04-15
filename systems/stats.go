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

	// Crit Chance: Base + Bonus (Capped at 1.0), convert excess to Crit Damage
    calculatedCritChance := attack.GetBaseCritChance() + attack.GetBonusCritChance()
    excessCritDamageBonus := 0.0
    if calculatedCritChance > 1.0 {
        excessCritDamageBonus = (calculatedCritChance - 1.0) / 2.0 // 50% conversion rate
        log.Printf("Entity %d: Crit chance %.2f exceeds 1.0. Adding %.2f bonus crit damage from excess.", entity, calculatedCritChance, excessCritDamageBonus)
        calculatedCritChance = 1.0 // Cap final crit chance
    }
    attack.SetFinalCritChance(calculatedCritChance)

    // --- IE / JG Conditional Bonus Crit Damage---
	// handle logic for: "If the holder's abilities can already critically strike, gain 10% Critical Strike Damage instead."
    equipment, okEq := s.world.GetEquipment(entity)
    numIE := 0
    numJG := 0
    if okEq {
        // Replace the loop:
        numIE = equipment.GetItemCount("TFT_Item_InfinityEdge")
        numJG = equipment.GetItemCount("TFT_Item_JeweledGauntlet")
    }
    totalCritItems := numIE + numJG

    // Check for trait source of ability crit
    traitCritMarkerType := reflect.TypeOf(components.CanAbilityCritFromTraits{})
    _, hasTraitCritMarker := s.world.GetComponent(entity, traitCritMarkerType)

    // Determine how many IE/JG grant the bonus damage
    numBonusGrantingItems := 0
    if hasTraitCritMarker {
        // If crit comes from trait, ALL IE/JG grant the bonus damage
        numBonusGrantingItems = totalCritItems
    } else {
        // If no trait crit, the first IE/JG enables the flag (via AbilityCritSystem),
        // and any subsequent ones grant the bonus damage.
        if totalCritItems > 1 {
            numBonusGrantingItems = totalCritItems - 1
        }
        // If totalCritItems is 0 or 1, and no trait source, no items grant the bonus damage.
    }

    // Calculate the bonus crit damage from the "instead" effect, because GetBonusCritDamageToGive() is cumulated from all IE/JE items, we'll need to remove the first one if it is not contributing to the bonus damage.
    conditionalCritDamageBonus := float64(numBonusGrantingItems) / float64(totalCritItems) * attack.GetBonusCritDamageToGive() 

    if conditionalCritDamageBonus > 0 {
        log.Printf("Entity %d: Applying +%.2f Crit Damage from %d IE/JG instance(s) due to 'already crit' condition.", entity, conditionalCritDamageBonus, numBonusGrantingItems)
    }

    // Crit Multiplier: Base + Bonus (from ItemEffect/Traits) + Conditional Bonus + Excess CritChance Bonus
    calculatedCritMultiplier := attack.GetBaseCritMultiplier() + attack.GetBonusCritMultiplier() + conditionalCritDamageBonus + excessCritDamageBonus
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
}
