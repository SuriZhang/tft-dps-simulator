package itemsys

import (
	"reflect"

	"tft-dps-simulator/internal/core/ecs"
	"tft-dps-simulator/internal/core/components/items"
)

// BaseStaticItemSystem applies the aggregated stat bonuses from items
// to the champion's base stats.
// It should run early in the update cycle, before systems that depend on final stats.
type BaseStaticItemSystem struct {
	world *ecs.World
}

func NewBaseStaticItemSystem(world *ecs.World) *BaseStaticItemSystem {
	return &BaseStaticItemSystem{world: world}
}

// ApplyStaticItemsBonus reads the ItemEffect component for relevant entities
// and modifies their base stat components (Health, Attack, Mana, etc.)
// based on the aggregated bonuses. This should be called after
// ItemEffect has been calculated/updated (e.g., after adding/removing items).
// Input: None (operates on the world state).
// Output: None (modifies components directly).
func (s *BaseStaticItemSystem) ApplyStaticItemsBonus() {
	// Define the component types needed for this system
	itemEffectType := reflect.TypeOf(items.ItemStaticEffect{})
	entitiesWithItemEffect := s.world.GetEntitiesWithComponents(itemEffectType)

	for _, entity := range entitiesWithItemEffect {
		// // Reset component bonuses BEFORE applying static item bonuses
        // if health, ok := s.world.GetHealth(entity); ok {
        //     health.ResetBonuses()
        // }
        // if attack, ok := s.world.GetAttack(entity); ok {
        //     attack.ResetBonuses()
        // }
        // if spell, ok := s.world.GetSpell(entity); ok {
        //     spell.ResetBonuses()
        // }
        // if crit, ok := s.world.GetCrit(entity); ok {
        //     crit.ResetBonuses()
        // }
        // if mana, ok := s.world.GetMana(entity); ok {
        //     mana.ResetBonuses()
        // }

		itemEffect, _ := s.world.GetItemEffect(entity)

		s.applyHealthBonuses(entity, itemEffect)
		s.applyAttackBonuses(entity, itemEffect)
		s.applyManaBonuses(entity, itemEffect)
		s.applyCritBonuses(entity, itemEffect)
		s.applySpellBonuses(entity, itemEffect)

		// Apply to Defense (when Defense component exists)
		// ...
	}
}

// --- Helper functions (Optional) ---
func (s *BaseStaticItemSystem) applyHealthBonuses(entity ecs.Entity, itemEffect *items.ItemStaticEffect) {
	if health, ok := s.world.GetHealth(entity); ok {
		health.AddBonusMaxHealth(itemEffect.GetBonusHealth())
		health.AddBonusPercentHealth(itemEffect.GetBonusPercentHp())
		health.AddBonusArmor(itemEffect.GetBonusArmor())
		health.AddBonusMR(itemEffect.GetBonusMR())
		health.AddBonusDurability(itemEffect.GetDurability())
	}

}

func (s *BaseStaticItemSystem) applyAttackBonuses(entity ecs.Entity, itemEffect *items.ItemStaticEffect) {
	if attack, ok := s.world.GetAttack(entity); ok {
		attack.AddBonusPercentAD(itemEffect.GetBonusPercentAD())
		attack.AddBonusDamageAmp(itemEffect.GetDamageAmp())
		attack.AddBonusPercentAttackSpeed(itemEffect.GetBonusPercentAttackSpeed())
	}
}

func (s *BaseStaticItemSystem) applyCritBonuses(entity ecs.Entity, itemEffect *items.ItemStaticEffect) {
	if crit, ok := s.world.GetCrit(entity); ok {
		crit.AddBonusCritChance(itemEffect.GetBonusCritChance())
		// crit.AddBonusCritMultiplier(itemEffect.GetBonusCritMultiplier())
		
		// Sepecific to Infinity Edge & Jeweled Gauntlet
		crit.AddBonusCritDamageToGive(itemEffect.GetCritDamageToGive())
	}
}

func (s *BaseStaticItemSystem) applyManaBonuses(entity ecs.Entity, itemEffect *items.ItemStaticEffect) {
	if mana, ok := s.world.GetMana(entity); ok {
		mana.AddBonusInitialMana(itemEffect.GetBonusInitialMana())
	}
}

func (s *BaseStaticItemSystem) applySpellBonuses(entity ecs.Entity, itemEffect *items.ItemStaticEffect) {
    if spell, ok := s.world.GetSpell(entity); ok {
        spell.AddBonusAP(itemEffect.GetBonusAP())
    }
}