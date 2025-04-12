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
    fmt.Printf("Current time: %.2f\n", s.currentTime)

	// Define component types needed for GetEntitiesWithComponents (still uses reflect)
	attackType := reflect.TypeOf(components.Attack{})
	healthType := reflect.TypeOf(components.Health{})
	teamType := reflect.TypeOf(components.Team{})
	posType := reflect.TypeOf(components.Position{})

	// Find all entities that *could* potentially attack
	potentialAttackers := s.world.GetEntitiesWithComponents(attackType, healthType, teamType, posType)
    
    // Filter for entities specifically on Team 0 (Player)
    // This is a bit redundant since we check again in the loop, but keeps the logic clear.
    var playerAttackers []ecs.Entity
    for _, entity := range potentialAttackers {
        team, ok := s.world.GetTeam(entity)
        // Check if the entity has a Team component and if its ID is 0
        if ok && team.ID == 0 {
            playerAttackers = append(playerAttackers, entity)
        }
    }

	// Process each potential attacker
	for _, attacker := range playerAttackers {
		// --- Get Components using Type-Safe Getters ---
		team, okTeam := s.world.GetTeam(attacker)
		attack, okAttack := s.world.GetAttack(attacker)
		health, okHealth := s.world.GetHealth(attacker)
        info, okInfo := s.world.GetChampionInfo(attacker)

		// Skip if essential components are missing (shouldn't happen with GetEntitiesWithComponents)
		if !okTeam || !okAttack || !okHealth || !okInfo {
			fmt.Printf("Warning: Attacker entity %d missing core components (Team/Attack/Health/ChampionInfo).\n", attacker)
			continue
		}

		// Skip dead attackers
		if health.Current <= 0 {
            fmt.Printf("Attacker %d is dead, skipping...\n", attacker)
			continue
		}

		// Check if attack is off cooldown
		attackCooldown := 1.0 / attack.AttackSpeed
		if s.currentTime-attack.LastAttackTime < attackCooldown {
            fmt.Printf("Attacker %s is not ready to attack yet (%.2f seconds remaining).\n", info.Name, attackCooldown-(s.currentTime-attack.LastAttackTime))
			continue // Not ready to attack yet
		}

		// Find the nearest ENEMY target using the utility function
		// Assumes FindNearestEnemy is also refactored to use type-safe getters
		target, found := utils.FindNearestEnemy(s.world, attacker, team.ID)
        targetInfo, _ := s.world.GetChampionInfo(target)
        fmt.Printf("Attacker %s found target %s\n", info.Name, targetInfo.Name)

		if !found {
            fmt.Printf("Attacker %s found no valid target.\n", info.Name)
			continue 
		}

		// Perform the attack (pass the pointer to attack component)
		err := s.performAttack(attacker, target, attack) // Pass the pointer
        if err != nil {
            fmt.Printf("Error performing attack: %v\n", err)
            continue // Skip to next attacker
        }
		// Update last attack time *after* a successful attack attempt
		// The 'attack' variable is already a pointer, so we can modify it directly.
		// No need to call AddComponent again unless performAttack replaces the component pointer.
		// Assuming performAttack modifies the pointed-to struct:
		attack.LastAttackTime = s.currentTime
        fmt.Printf("Updated last attack time for entity %d to %.2f\n", attacker, attack.LastAttackTime)
	}
}

// performAttack handles the logic of an entity attacking another.
// Now accepts a pointer to the Attack component.
func (s *AutoAttackSystem) performAttack(attacker, target ecs.Entity, attack *components.Attack) error { 
    fmt.Println("Performing attack...")

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
		return fmt.Errorf("target %s (Entity %d) has no Health component", targetName, target)
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

    fmt.Printf("Damage before reduction: %.1f, after reduction: %.1f\n", damage, finalDamage)

	// Apply damage to target (modify the struct pointed to by targetHealth)
	targetHealth.Current -= finalDamage

	// --- End Damage Calculation ---

	// Debug output
	critText := ""
	if isCrit {
		critText = "***CRIT! "
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

