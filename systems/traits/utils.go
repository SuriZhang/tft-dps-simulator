package traitsys

import (
	"reflect"

	"github.com/suriz/tft-dps-simulator/components"
	"github.com/suriz/tft-dps-simulator/ecs"
)

// hasTrait checks if a champion has a specific trait. (Keep this helper)
func (h *RapidfireHandler) hasTrait(world *ecs.World, entity ecs.Entity, traitApiName string) bool {
    traits, ok := world.GetTraits(entity)
    if !ok {
        return false // No traits component found
    }

    for _, t := range traits.GetTraits() {
        if t == traitApiName {
            return true
        }
    }
    return false
}

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