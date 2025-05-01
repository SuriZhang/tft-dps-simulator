package traitsys

import (
	"reflect"

	"tft-dps-simulator/internal/core/components"
	"tft-dps-simulator/internal/core/ecs"
)

// getChampionsByTeam finds all entities belonging to a specific team. (Keep this helper)
func getChampionsByTeam(world *ecs.World, teamID int) []ecs.Entity {
    var teamChampions []ecs.Entity
    teamType := reflect.TypeOf(components.Team{})
    entities := world.GetEntitiesWithComponents(teamType)
    for _, entity := range entities {
        if teamComp, ok := world.GetTeam(entity); ok && teamComp.ID == teamID {
            teamChampions = append(teamChampions, entity)
        }
    }
    return teamChampions
}