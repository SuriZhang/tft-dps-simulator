package utils

import (
	"fmt"
	"math"
	"reflect"

	"github.com/suriz/tft-dps-simulator/components"
	"github.com/suriz/tft-dps-simulator/ecs"
)

// FindNearestEnemy finds the closest entity on an opposing team.
// Uses type-safe getters. Range is ignored.
func FindNearestEnemy(world *ecs.World, source ecs.Entity, sourceTeamID int) (ecs.Entity, bool) {
    // Get source position using type-safe getter
    sourcePos, okPosSource := world.GetPosition(source)
    if !okPosSource {
        fmt.Printf("Warning: Source entity %d for targeting has no Position.\n", source)
        return 0, false // Source has no position
    }

    // Define component types needed for GetEntitiesWithComponents (still uses reflect)
    posType := reflect.TypeOf(components.Position{})
    healthType := reflect.TypeOf(components.Health{})
	teamType := reflect.TypeOf(components.Team{})

	// Get all entities that have Position, Health, and Team components
	entities := world.GetEntitiesWithComponents(posType, healthType, teamType)

	// Filter for entities specifically on Team 1
	var potentialTargets []ecs.Entity
	for _, entity := range entities {
		team, ok := world.GetTeam(entity)
		// Ensure the entity has a Team component and its ID is 1
		if ok && team.ID == 1 {
			potentialTargets = append(potentialTargets, entity)
		}
	}

    var closestEnemy ecs.Entity
    closestDistSq := math.MaxFloat64
    foundTarget := false

    for _, target := range potentialTargets {

        // --- Get Target Components using Type-Safe Getters ---
        targetHealth, okHealth := world.GetHealth(target)
        targetPos, okPos := world.GetPosition(target)

        // Check if components exist (should exist based on query, but good practice)
        if  !okHealth || !okPos {
            continue
        }

        // Check if entity is alive
        if targetHealth.CurrentHP <= 0 {
            continue // Dead, skip
        }

        // Calculate distance squared
        dx := targetPos.X - sourcePos.X
        dy := targetPos.Y - sourcePos.Y
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