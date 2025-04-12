package systems

import (
	"fmt"
	"math/rand"
	"reflect"

	"github.com/suriz/tft-dps-simulator/components"
	"github.com/suriz/tft-dps-simulator/ecs"
	"github.com/suriz/tft-dps-simulator/utils"
)

// AutoAttackSystem handles champion auto attacks based on Team ID.
type AutoAttackSystem struct {
	world       *ecs.World
	currentTime float64
}

// NewAutoAttackSystem creates a new auto attack system.
func NewAutoAttackSystem(world *ecs.World) *AutoAttackSystem {
	return &AutoAttackSystem{
		world:       world,
		currentTime: 0.0,
	}
}

// Update processes auto attacks for the current timestep.
func (s *AutoAttackSystem) Update(deltaTime float64) {
	s.currentTime += deltaTime

	// Define component types needed for GetEntitiesWithComponents (still uses reflect)
	attackType := reflect.TypeOf(components.Attack{})
	healthType := reflect.TypeOf(components.Health{})
	teamType := reflect.TypeOf(components.Team{})
	posType := reflect.TypeOf(components.Position{})

	// Find all entities that *could* potentially attack
	potentialAttackers := s.world.GetEntitiesWithComponents(attackType, healthType, teamType, posType)

	// Process each potential attacker
	for _, attacker := range potentialAttackers {
		// --- Get Components using Type-Safe Getters ---
		team, okTeam := s.world.GetTeam(attacker)
		attack, okAttack := s.world.GetAttack(attacker)
		health, okHealth := s.world.GetHealth(attacker)

		// Skip if essential components are missing (shouldn't happen with GetEntitiesWithComponents)
		if !okTeam || !okAttack || !okHealth {
			fmt.Printf("Warning: Attacker entity %d missing core components (Team/Attack/Health/Position).\n", attacker)
			continue
		}

		// --- Rule #4: Only Team 0 (Player) attacks ---
		if team.ID != 0 {
			continue // Skip if not on the player team (Team 0)
		}

		// Skip dead attackers
		if health.Current <= 0 {
			continue
		}

		// Check if attack is off cooldown
		attackCooldown := 1.0 / attack.AttackSpeed
		if s.currentTime-attack.LastAttackTime < attackCooldown {
			continue // Not ready to attack yet
		}

		// Find the nearest ENEMY target using the utility function
		// Assumes FindNearestEnemy is also refactored to use type-safe getters
		target, found := utils.FindNearestEnemy(s.world, attacker, team.ID)
		if !found {
			continue // No valid targets found
		}

		// Perform the attack (pass the pointer to attack component)
		s.performAttack(attacker, target, attack) // Pass the pointer

		// Update last attack time *after* a successful attack attempt
		// The 'attack' variable is already a pointer, so we can modify it directly.
		// No need to call AddComponent again unless performAttack replaces the component pointer.
		// Assuming performAttack modifies the pointed-to struct:
		attack.LastAttackTime = s.currentTime
		// If performAttack *could* replace the component (unlikely), you'd need:
		// updatedAttack, _ := s.world.GetAttack(attacker)
		// updatedAttack.LastAttackTime = s.currentTime
	}
}

// performAttack handles the logic of an entity attacking another.
// Now accepts a pointer to the Attack component.
func (s *AutoAttackSystem) performAttack(attacker, target ecs.Entity, attack *components.Attack) error { // Takes *components.Attack
	// Get attacker info for logging (using type-safe getter)
	attackerInfo, okInfoAttacker := s.world.GetChampionInfo(attacker)
	attackerName := fmt.Sprintf("Entity %d", attacker) // Default name
	if okInfoAttacker {
		attackerName = attackerInfo.Name
	}

	// Get target info for logging (using type-safe getter)
	targetInfo, okInfoTarget := s.world.GetChampionInfo(target)
	targetName := fmt.Sprintf("Entity %d", target) // Default name
	if okInfoTarget {
		targetName = targetInfo.Name
	}

	// Get target health (using type-safe getter)
	targetHealth, okHealth := s.world.GetHealth(target)
	if !okHealth {
		return fmt.Errorf("Error: Target %s (Entity %d) has no Health component.", targetName, target)
		 // Cannot apply damage
	}

	// --- Damage Calculation ---
	damage := attack.Damage // Access fields directly from the pointer

	// Check for critical hit
	isCrit := rand.Float64() < attack.CritChance
	if isCrit {
		damage = damage * attack.CritMultiplier
	}

	// Apply armor formula: damage * (100 / (100 + armor))
	// Ensure targetHealth.Armor is accessed correctly
	damageReduction := 100.0 / (100.0 + targetHealth.Armor)
	finalDamage := damage * damageReduction

	// Apply damage to target (modify the struct pointed to by targetHealth)
	targetHealth.Current -= finalDamage

	// NO need to call AddComponent(target, targetHealth) here,
	// because targetHealth is a pointer to the component in the world's map.
	// We modified the data *in place*.

	// --- End Damage Calculation ---

	// Debug output
	critText := ""
	if isCrit {
		critText = "CRIT! "
	}
	displayHealth := targetHealth.Current
	if displayHealth < 0 {
		displayHealth = 0
	}
	fmt.Printf("%s attacks %s for %.1f damage %s(%.1f HP remaining)\n",
		attackerName, targetName, finalDamage, critText, displayHealth)

	// Check if target died
	if targetHealth.Current <= 0 {
		fmt.Printf("%s has been defeated!\n", targetName)
	}
	return nil // Indicate successful execution
}

