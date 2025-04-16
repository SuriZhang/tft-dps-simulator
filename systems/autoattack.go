package systems

import (
	"fmt"
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
	// fmt.Printf("Current time: %.2f\n", s.currentTime) 

	// Define component types needed (using reflect for now)
	attackType := reflect.TypeOf(components.Attack{})
	teamType := reflect.TypeOf(components.Team{})
	posType := reflect.TypeOf(components.Position{})
	healthType := reflect.TypeOf(components.Health{}) // Needed to check if attacker is alive

	// Get all entities that can potentially attack (have Attack, Team, Position, Health)
	potentialAttackers := s.world.GetEntitiesWithComponents(attackType, teamType, posType, healthType)

	// Process each potential attacker
	for _, attacker := range potentialAttackers {
		// --- Get Attacker Components ---
		attack, okAtk := s.world.GetAttack(attacker)
		team, okTeam := s.world.GetTeam(attacker)
		pos, okPos := s.world.GetPosition(attacker)
		health, okHealth := s.world.GetHealth(attacker) // Get health component

		// Basic validation: Ensure components exist and attacker is alive and belongs to the player team (Team 0)
		if !okAtk || !okTeam || !okPos || !okHealth || health.CurrentHP <= 0 || team.ID != 0 {
			continue // Skip if missing components, dead, or not on player team
		}

		// Check if attacker can attack (Attack Speed > 0)
		if attack.GetFinalAttackSpeed() <= 0 {
			fmt.Printf("Attacker %d has 0 AS, skipping.\n", attacker) // Optional log
			continue
		}

		// --- Targeting ---
		// Find the nearest enemy (Team 1)
		target, foundTarget := utils.FindNearestEnemy(s.world, attacker, team.ID) // Pass attacker's team ID

		if !foundTarget {
			fmt.Printf("Attacker %d found no target.\n", attacker) // Optional log
			continue // No target found for this attacker, move to the next
		}

		// --- Range Check ---
		targetPos, okTargetPos := s.world.GetPosition(target)
		if !okTargetPos {
			fmt.Printf("Attacker %d found target %d with no position.\n", attacker, target) // Optional log
			continue // Target has no position, cannot calculate range
		}

		dx := targetPos.X - pos.X
		dy := targetPos.Y - pos.Y
		distSq := dx*dx + dy*dy
		attackRange := attack.GetFinalRange()
		rangeSq := attackRange * attackRange

		if distSq > rangeSq {
			fmt.Printf("Attacker %d target %d is out of range (DistSq: %.2f, RangeSq: %.2f).\n", attacker, target, distSq, rangeSq) // Optional log
			continue // Target found, but out of range
		}

		// --- Attack Timing ---
		attackDelay := 1.0 / attack.GetFinalAttackSpeed()
		timeSinceLastAttack := s.currentTime - attack.GetLastAttackTime()

		// Check if enough time has passed to attack again
		if timeSinceLastAttack >= attackDelay {
			// Perform the attack
			err := s.performAttack(attacker, target, attack)
			if err != nil {
				fmt.Printf("Error performing attack from %d to %d: %v\n", attacker, target, err)
				// Decide how to handle attack errors (e.g., skip, log, panic)
			} else {
				// Update last attack time ONLY on successful attack attempt
				attack.SetLastAttackTime(s.currentTime)
			}
		}
		// --- End Attack Timing ---
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
	damage := attack.GetFinalAD() 

	damageAmp := attack.GetFinalDamageAmp()
	if damageAmp != 0 {	
		damage = damage * (1 + damageAmp) // Apply damage amplification
	}

	// Check for critical hit
	// isCrit := rand.Float64() < attack.GetFinalCritChance()
	// if isCrit {
	// 	damage = damage * attack.GetFinalCritMultiplier()
	// }

	// Use Expected Value for Crit Damange calculation
	expectedCritMultiplier := (1.0 - attack.GetFinalCritChance()) + (attack.GetFinalCritChance() * attack.GetFinalCritMultiplier())

	damage = damage * expectedCritMultiplier

	// Apply armor formula: damage * (100 / (100 + armor))* (1 - durability)
	// Ensure targetHealth.Armor is accessed correctly
	damageReduction := (100.0 / (100.0 + targetHealth.GetFinalArmor())) * (1.0 - targetHealth.GetFinalDurability()) 
	finalDamage := damage * damageReduction

    fmt.Printf("Damage before reduction: %.1f, after reduction: %.1f\n", damage, finalDamage)

	// Apply damage to target (modify the struct pointed to by targetHealth)
	targetHealth.CurrentHP -= finalDamage

	// TODO: maybe move the logic to a separate mana system?
	// --- Mana Gain for Attacker ---
    attackerMana, okMana := s.world.GetMana(attacker)
    if okMana {
        // Standard TFT mana gain per auto-attack
        manaGain := 10.0
        attackerMana.AddCurrentMana(manaGain)

        if attackerMana.GetCurrentMana() > attackerMana.GetMaxMana() {
            attackerMana.SetCurrentMana(attackerMana.GetMaxMana())
        }
        fmt.Printf("%s gains %.1f mana (now %.1f / %.1f)\n", attackerName, manaGain, attackerMana.GetCurrentMana(), attackerMana.GetMaxMana())
    } else {
        fmt.Printf("Warning: Attacker %s has no Mana component, cannot gain mana.\n", attackerName)
    }
    // --- End Mana Gain ---

	// --- End Damage Calculation ---

	// Debug output
	// critText := ""
	// if isCrit {
	// 	critText = "***CRIT! "
	// }
	displayHealth := targetHealth.CurrentHP
	if displayHealth < 0 {
		displayHealth = 0
	}
	fmt.Printf("%s attacks %s for %.1f damage (%.1f HP remaining)\n",
		attackerName, targetName, finalDamage, displayHealth)

	// Check if target died
	if targetHealth.CurrentHP <= 0 {
		fmt.Printf("%s has been defeated!\n", targetName)
	}

	// TODO: Target also gains mana on being attacked

	return nil // Indicate successful execution
}

