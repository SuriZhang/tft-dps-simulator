package systems

import (
	"log"
	"reflect"

	"github.com/suriz/tft-dps-simulator/components"
	"github.com/suriz/tft-dps-simulator/ecs"
	eventsys "github.com/suriz/tft-dps-simulator/systems/events"
	"github.com/suriz/tft-dps-simulator/utils"
)

// AutoAttackSystem handles champion auto attacks based on Team ID.
type AutoAttackSystem struct {
	world       *ecs.World
	eventBus    eventsys.EventBus // Add event bus
	currentTime float64
}

// NewAutoAttackSystem creates a new auto attack system.
func NewAutoAttackSystem(world *ecs.World, bus eventsys.EventBus) *AutoAttackSystem { // Accept event bus
	return &AutoAttackSystem{
		world:       world,
		eventBus:    bus, // Store event bus
		currentTime: 0.0,
	}
}

// TriggerAutoAttack processes auto attacks for the current timestep.
func (s *AutoAttackSystem) TriggerAutoAttack(deltaTime float64) {
	s.currentTime += deltaTime

	attackType := reflect.TypeOf(components.Attack{})
	teamType := reflect.TypeOf(components.Team{})
	posType := reflect.TypeOf(components.Position{})
	healthType := reflect.TypeOf(components.Health{})

	potentialAttackers := s.world.GetEntitiesWithComponents(attackType, teamType, posType, healthType)

	for _, attacker := range potentialAttackers {
		attack, okAtk := s.world.GetAttack(attacker)
		team, okTeam := s.world.GetTeam(attacker)
		pos, okPos := s.world.GetPosition(attacker)
		health, okHealth := s.world.GetHealth(attacker)

		if !okAtk || !okTeam || !okPos || !okHealth || health.CurrentHP <= 0 || team.ID != 0 {
			continue // Skip if missing components, dead, or not on player team
		}

		if attack.GetFinalAttackSpeed() <= 0 {
			log.Printf("Attacker %d has 0 AS, skipping.\n", attacker)
			continue
		}

		target, foundTarget := utils.FindNearestEnemy(s.world, attacker, team.ID)
		if !foundTarget {
			log.Printf("Attacker %d found no target.\n", attacker)
			continue // No target found for this attacker, move to the next
		}

		targetPos, okTargetPos := s.world.GetPosition(target)
		if !okTargetPos {
			log.Printf("Attacker %d found target %d with no position.\n", attacker, target)
			continue // Target has no position, cannot calculate range
		}

		dx := targetPos.GetX() - pos.GetX()
		dy := targetPos.GetY() - pos.GetY()
		distSq := dx*dx + dy*dy
		attackRange := attack.GetFinalRange()
		rangeSq := attackRange * attackRange

		if distSq > rangeSq {
			log.Printf("Attacker %d target %d is out of range (DistSq: %.2f, RangeSq: %.2f).\n", attacker, target, distSq, rangeSq)
			continue // Target found, but out of range
		}

		attackDelay := 1.0 / attack.GetFinalAttackSpeed()
		timeSinceLastAttack := s.currentTime - attack.GetLastAttackTime()

		if timeSinceLastAttack >= attackDelay {
			// --- Enqueue AttackLandedEvent ---
			// Calculate base damage (before reductions, crit, etc.)
			baseDamage := attack.GetFinalAD()

			// Create and enqueue the event
			attackEvent := eventsys.AttackLandedEvent{
				Source:     attacker,
				Target:     target,
				BaseDamage: baseDamage,
				Timestamp:  s.currentTime, // Include timestamp for potential future use
			}
			s.eventBus.Enqueue(attackEvent)
			// Note: LastAttackTime is NOT updated here. It will be updated by a handler
			// reacting to AttackLandedEvent to ensure the attack actually resolves.
			// --- End Enqueue ---
		}
	}
}
