package factory

import (
	"fmt"
	"log"

	"github.com/suriz/tft-dps-simulator/components"
	"github.com/suriz/tft-dps-simulator/components/effects"
	"github.com/suriz/tft-dps-simulator/data"
	"github.com/suriz/tft-dps-simulator/ecs"
)

// ChampionFactory creates champion entities from champion data.
type ChampionFactory struct {
	world *ecs.World
}

// NewChampionFactory creates a new ChampionFactory.
func NewChampionFactory(world *ecs.World) *ChampionFactory {
	return &ChampionFactory{
		world: world,
	}
}

// StarMultiplier returns the stat multiplier for a given star level.
func StarMultiplier(starLevel int) float64 {
	switch starLevel {
	case 1:
		return 1.0
	case 2:
		return 1.8
	case 3:
		return 3.24
	default:
		log.Printf("Warning: Invalid star level %d requested, using 1-star multiplier.\n", starLevel)
		return 1.0
	}
}

// CreateChampion creates a champion entity with components from data.
// It now returns an error if adding any component fails.
func (cf *ChampionFactory) CreateChampion(championData data.Champion, starLevel int) (ecs.Entity, error) {
	// Create entity
	entity := cf.world.NewEntity()
	var err error // Variable to hold potential errors

	// Apply star level multiplier
	multiplier := StarMultiplier(starLevel)

	// --- Add components, checking for errors after each AddComponent call ---

	// Health
	healthComp := components.NewHealth(
		championData.Stats.HP*multiplier,
		championData.Stats.Armor,
		championData.Stats.MagicResist,
	)
	err = cf.world.AddComponent(entity, healthComp)
	if err != nil {
		return 0, fmt.Errorf("failed to add Health component to %s: %w", championData.Name, err)
	}

	// Attack
	attackComp := components.NewAttack(
		championData.Stats.Damage*multiplier,
		championData.Stats.AttackSpeed,
		championData.Stats.Range,
		0.0, // startup time, not available in data yet
		0.0, // recovery time, not available in data yet
	)
	err = cf.world.AddComponent(entity, attackComp)
	if err != nil {
		return 0, fmt.Errorf("failed to add Attack component to %s: %w", championData.Name, err)
	}

	// Crit
	critComp := components.NewCrit(	championData.Stats.CritChance,
		championData.Stats.CritMultiplier,)
	err = cf.world.AddComponent(entity, critComp)
	if err != nil {
		return 0, fmt.Errorf("failed to add Crit component to %s: %w", championData.Name, err)
	}

	// Mana
	manaComp := components.NewMana(
		championData.Stats.Mana,
		championData.Stats.InitialMana,
	)
	err = cf.world.AddComponent(entity, manaComp)
	if err != nil {
		return 0, fmt.Errorf("failed to add Mana component to %s: %w", championData.Name, err)
	}

	// Add Spell component
	// TODO 1: fix cooldown later, assume it's 1 for now
	// TODO 2: process spell variables
	spellComp := components.NewSpell(
		championData.Ability.Name,
		championData.Ability.Icon,
		championData.Stats.Mana,
		1, // startup time (not available in data yet)
		1, // recovery time (not available in data yet)
	)
	err = cf.world.AddComponent(entity, spellComp)
	if err != nil {
		return 0, fmt.Errorf("failed to add Spell component to %s: %w", championData.Name, err)
	}

	// ChampionInfo
	infoComp := components.NewChampionInfo(
		championData.ApiName,
		championData.Name,
		championData.Cost,
		starLevel,
	)
	err = cf.world.AddComponent(entity, infoComp)
	if err != nil {
		return 0, fmt.Errorf("failed to add ChampionInfo component to %s: %w", championData.Name, err)
	}

	// Traits (only if they exist)
	if len(championData.Traits) > 0 {
		traitsComp := components.NewTraits(championData.Traits)
		err = cf.world.AddComponent(entity, traitsComp)
		if err != nil {
			return 0, fmt.Errorf("failed to add Traits component to %s: %w", championData.Name, err)
		}
	}

	// create empty inventory
	err = cf.world.AddComponent(entity, components.NewEquipment())
	if err != nil {
		return 0, fmt.Errorf("failed to add Inventory component to %s: %w", championData.Name, err)
	}

	// create empty item effect component
	err = cf.world.AddComponent(entity, effects.NewItemStaticEffect())
	if err != nil {
		return 0, fmt.Errorf("failed to add ItemEffect component to %s: %w", championData.Name, err)
	}

	err = cf.world.AddComponent(entity, components.NewPosition(0, 0))
	if err != nil {
		return 0, fmt.Errorf("failed to add Position component to %s: %w", championData.Name, err)
	}
	// Add other essential components here (Position, Team are usually added later/externally)

	// If we reached here, all essential components were added successfully
	return entity, nil
}

// CreateChampionByApiName creates a champion entity by searching for it by name.
// It now propagates errors from CreateChampion.
func (cf *ChampionFactory) CreateChampionByApiName(apiName string, starLevel int, team int) (ecs.Entity, error) {
	// Find champion data by name using the function from the data package
	championData := data.GetChampionByApiName(apiName)
	if championData == nil {
		return 0, fmt.Errorf("champion data for '%s' not found", apiName)
	}

	// Create entity using the found data, passing any error up
	entity, err := cf.CreateChampion(*championData, starLevel)
	if err != nil {
		return 0, err
	}

	// assign a team
	err = cf.world.AddComponent(entity, components.NewTeam(team))
	if err != nil {
		return 0, fmt.Errorf("failed to add Team component to %s: %w", championData.Name, err)
	}

	return entity, nil
}

// CreatePlayerChampion creates a champion entity for the player team (team ID 0).
func (cf *ChampionFactory) CreatePlayerChampion(apiName string, starLevel int) (ecs.Entity, error) {
	return cf.CreateChampionByApiName(apiName, starLevel, 0)
}

// CreateEnemyChampion creates a champion entity for the enemy team (team ID 1).
func (cf *ChampionFactory) CreateEnemyChampion(apiName string, starLevel int) (ecs.Entity, error) {
	return cf.CreateChampionByApiName(apiName, starLevel, 1)
}
