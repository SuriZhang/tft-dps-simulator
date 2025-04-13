package systems

import (
	"math"
	"testing"

	"github.com/suriz/tft-dps-simulator/components"
	"github.com/suriz/tft-dps-simulator/ecs"
)

// Helper function to setup a basic world for testing
func setupTestWorldWithEntities() (*ecs.World, *AutoAttackSystem, ecs.Entity, ecs.Entity) {
	world := ecs.NewWorld()
	system := NewAutoAttackSystem(world)

	// Attacker (Team 0)
	attacker := world.CreateEntity()
	world.AddComponent(attacker, components.NewChampionInfo("Attacker", "Attacker", 1, 1))
	world.AddComponent(attacker, components.NewHealth(100, 0, 0))
	world.AddComponent(attacker, components.NewAttack(50, 1.0, 1, 0.25, 1.5)) // 50 dmg, 1.0 AS, 1 range, 25% crit, 1.5x crit mult
	world.AddComponent(attacker, components.NewPosition(1, 1))
	world.AddComponent(attacker, components.NewTeam(0))

	// Target (Team 1)
	target := world.CreateEntity()
	world.AddComponent(target, components.NewChampionInfo("Target", "Target", 1, 1))
	world.AddComponent(target, components.NewHealth(1000, 0, 0)) // High health, 0 armor/mr
	world.AddComponent(target, components.NewPosition(2, 1))     // Within range 1
	world.AddComponent(target, components.NewTeam(1))

	return world, system, attacker, target
}

// Test a single basic attack hitting the target
func TestAutoAttack_BasicHit(t *testing.T) {
	world, system, attacker, target := setupTestWorldWithEntities()
	initialTargetHealth, _ := world.GetHealth(target)
	initialHP := initialTargetHealth.CurrentHP

	// Update slightly more than cooldown (1.0s for 1.0 AS)
	system.Update(1.1)

	finalTargetHealth, ok := world.GetHealth(target)
	if !ok {
		t.Fatalf("Target lost health component")
	}

	// Expected damage = 50 (base) * (100 / (100 + 0 armor)) = 50
	expectedHP := initialHP - 50.0
	// Use a small tolerance for float comparison
	if math.Abs(finalTargetHealth.CurrentHP-expectedHP) > 0.01 {
		t.Errorf("Expected target health %.2f, got %.2f", expectedHP, finalTargetHealth.CurrentHP)
	}

	// Check if attacker's LastAttackTime was updated
	attackComp, _ := world.GetAttack(attacker)
	if math.Abs(attackComp.LastAttackTime-system.currentTime) > 0.01 {
		t.Errorf("Expected attacker LastAttackTime %.2f, got %.2f", system.currentTime, attackComp.LastAttackTime)
	}
}

// Test that the attacker respects the cooldown
func TestAutoAttack_Cooldown(t *testing.T) {
	world, system, _, target := setupTestWorldWithEntities()
	initialTargetHealth, _ := world.GetHealth(target)
	initialHP := initialTargetHealth.CurrentHP

	// Update less than cooldown (1.0s) -> should not attack
	system.Update(0.5)

	targetHealthAfterFirstUpdate, _ := world.GetHealth(target)
	if math.Abs(targetHealthAfterFirstUpdate.CurrentHP-initialHP) > 0.01 {
		t.Errorf("Target health changed before cooldown expired, expected %.2f, got %.2f", initialHP, targetHealthAfterFirstUpdate.CurrentHP)
	}

	// Update past the cooldown -> should attack
	system.Update(0.6) // Total time = 0.5 + 0.6 = 1.1s

	targetHealthAfterSecondUpdate, _ := world.GetHealth(target)
	expectedHP := initialHP - 50.0 // Should have taken one hit
	if math.Abs(targetHealthAfterSecondUpdate.CurrentHP-expectedHP) > 0.01 {
		t.Errorf("Target health incorrect after cooldown, expected %.2f, got %.2f", expectedHP, targetHealthAfterSecondUpdate.CurrentHP)
	}
}

// Test armor reduction (simplified, doesn't test crit interaction here)
func TestAutoAttack_ArmorReduction(t *testing.T) {
	world, system, _, target := setupTestWorldWithEntities()

	// Give target armor
	targetHealth, _ := world.GetHealth(target)
	targetHealth.BaseArmor = 100 // Should reduce physical damage by 50%
	initialHP := targetHealth.CurrentHP

	// Update past cooldown
	system.Update(1.1)

	finalTargetHealth, _ := world.GetHealth(target)

	// Expected damage = 50 (base) * (100 / (100 + 100 armor)) = 25
	expectedHP := initialHP - 25.0
	if math.Abs(finalTargetHealth.CurrentHP-expectedHP) > 0.01 {
		t.Errorf("Armor reduction failed: Expected target health %.2f, got %.2f", expectedHP, finalTargetHealth.CurrentHP)
	}
}

// Test that friendly fire is prevented
func TestAutoAttack_NoFriendlyFire(t *testing.T) {
	world, system, attacker, target := setupTestWorldWithEntities()

	// Change target's team to be the same as attacker
	targetTeam, _ := world.GetTeam(target)
	targetTeam.ID = 0 // Now both are Team 0
	initialTargetHealth, _ := world.GetHealth(target)
	initialHP := initialTargetHealth.CurrentHP

	// Update past cooldown
	system.Update(1.1)

	finalTargetHealth, _ := world.GetHealth(target)
	// Health should NOT change
	if math.Abs(finalTargetHealth.CurrentHP-initialHP) > 0.01 {
		t.Errorf("Friendly fire occurred! Target health changed, expected %.2f, got %.2f", initialHP, finalTargetHealth.CurrentHP)
	}

	// Check attacker's LastAttackTime - should NOT have been updated as no attack occurred
	attackComp, _ := world.GetAttack(attacker)
	if attackComp.LastAttackTime != 0.0 {
		t.Errorf("Attacker LastAttackTime updated (%.2f) despite no valid target", attackComp.LastAttackTime)
	}
}

// Test that a dead attacker does not attack
func TestAutoAttack_DeadAttacker(t *testing.T) {
	world, system, attacker, target := setupTestWorldWithEntities()

	// Kill the attacker
	attackerHealth, _ := world.GetHealth(attacker)
	attackerHealth.CurrentHP = 0
	initialTargetHealth, _ := world.GetHealth(target)
	initialHP := initialTargetHealth.CurrentHP

	// Update past cooldown
	system.Update(1.1)

	finalTargetHealth, _ := world.GetHealth(target)
	// Health should NOT change
	if math.Abs(finalTargetHealth.CurrentHP-initialHP) > 0.01 {
		t.Errorf("Dead attacker attacked! Target health changed, expected %.2f, got %.2f", initialHP, finalTargetHealth.CurrentHP)
	}
}

// Note: Testing Crit requires either controlling the RNG or running many times.
// For simplicity, we won't test the exact probability here, but you could
// temporarily set crit chance to 1.0 or 0.0 to test the damage calculation.

// Example: Test guaranteed crit damage
func TestAutoAttack_GuaranteedCrit(t *testing.T) {
	world, system, attacker, target := setupTestWorldWithEntities()

	// Force crit
	attackComp, _ := world.GetAttack(attacker)
	attackComp.BaseCritChance = 1.0                 // 100% crit chance
	critMultiplier := attackComp.BaseCritMultiplier // Should be 1.5x
	baseDamage := attackComp.BaseAD             // Should be 50

	initialTargetHealth, _ := world.GetHealth(target)
	initialHP := initialTargetHealth.CurrentHP

	// Update past cooldown
	system.Update(1.1)

	finalTargetHealth, _ := world.GetHealth(target)

	// Expected damage = 50 (base) * 1.5 (crit) * (100 / (100 + 0 armor)) = 75
	expectedHP := initialHP - (baseDamage * critMultiplier)
	if math.Abs(finalTargetHealth.CurrentHP-expectedHP) > 0.01 {
		t.Errorf("Crit damage incorrect: Expected target health %.2f, got %.2f", expectedHP, finalTargetHealth.CurrentHP)
	}
}
