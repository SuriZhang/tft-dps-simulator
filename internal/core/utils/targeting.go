package utils

import (
	"log"
	"math"
	"reflect"

	"tft-dps-simulator/internal/core/components"
	"tft-dps-simulator/internal/core/ecs"
	"tft-dps-simulator/internal/core/entity"
)

// FindNearestEnemy finds the closest entity on an opposing team.
// Uses type-safe getters. Range is ignored.
func FindNearestEnemy(world *ecs.World, source entity.Entity, sourceTeamID int) (entity.Entity, bool) {
	// Get source position using type-safe getter
	sourcePos, okPosSource := world.GetPosition(source)
	if !okPosSource {
		log.Printf("Warning: Source entity %d for targeting has no Position.\n", source)
		return 0, false // Source has no position
	}

	// Define component types needed for GetEntitiesWithComponents (still uses reflect)
	posType := reflect.TypeOf(components.Position{})
	healthType := reflect.TypeOf(components.Health{})
	teamType := reflect.TypeOf(components.Team{})

	// Get all entities that have Position, Health, and Team components
	entities := world.GetEntitiesWithComponents(posType, healthType, teamType)

	// Filter for entities specifically on Team 1
	var potentialTargets []entity.Entity
	for _, entity := range entities {
		team, ok := world.GetTeam(entity)
		// Ensure the entity is on a different team and has a valid Team component
		if ok && team.ID != sourceTeamID {
			potentialTargets = append(potentialTargets, entity)
		}
	}

	var closestEnemy entity.Entity
	closestDistSq := math.MaxInt32 
	foundTarget := false

	for _, target := range potentialTargets {

		// --- Get Target Components using Type-Safe Getters ---
		targetHealth, okHealth := world.GetHealth(target)
		targetPos, okPos := world.GetPosition(target)

		// Check if components exist (should exist based on query, but good practice)
		if !okHealth || !okPos {
			continue
		}

		// Check if entity is alive
		if targetHealth.CurrentHP <= 0 {
			continue // Dead, skip
		}

		// Calculate distance squared
		dx := targetPos.GetX() - sourcePos.GetX()
		dy := targetPos.GetY() - sourcePos.GetY()
		distSq := dx*dx + dy*dy

		// Update closest if this one is closer
		if distSq < closestDistSq {
			closestDistSq = distSq
			closestEnemy = target
			foundTarget = true
		}
	}

	return closestEnemy, foundTarget
}

func DistSq(x1, y1, x2, y2 int) int {
	dx := x2 - x1
	dy := y2 - y1
	return dx*dx + dy*dy
}