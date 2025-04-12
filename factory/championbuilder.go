package factory

import (
	"github.com/suriz/tft-dps-simulator/components"
	"github.com/suriz/tft-dps-simulator/data"
	"github.com/suriz/tft-dps-simulator/ecs"
)

// ChampionBuilder creates champion entities from champion data
type ChampionBuilder struct {
    world *ecs.World
}

// NewChampionBuilder creates a new ChampionBuilder
func NewChampionBuilder(world *ecs.World) *ChampionBuilder {
    return &ChampionBuilder{
        world: world,
    }
}

// StarMultiplier returns the stat multiplier for a given star level
func StarMultiplier(starLevel int) float64 {
    switch starLevel {
    case 1:
        return 1.0
    case 2:
        return 1.8
    case 3:
        return 3.24
    default:
        return 1.0
    }
}

// CreateChampion creates a champion entity with components from data
func (cb *ChampionBuilder) CreateChampion(championData data.Champion, starLevel int) ecs.Entity {
    // Create entity
    entity := cb.world.CreateEntity()
    
    // Apply star level multiplier
    multiplier := StarMultiplier(starLevel)
    
    // Add components
    cb.world.AddComponent(entity, components.NewHealth(
        championData.Stats.HP * multiplier,
        championData.Stats.Armor,
        championData.Stats.MagicResist,
    ))
    
    cb.world.AddComponent(entity, components.NewAttack(
        championData.Stats.Damage * multiplier,
        championData.Stats.AttackSpeed,
        championData.Stats.Range,
        championData.Stats.CritChance,
        championData.Stats.CritMultiplier,
    ))
    
    cb.world.AddComponent(entity, components.NewMana(
        championData.Stats.Mana,
        championData.Stats.InitialMana,
    ))
    
    cb.world.AddComponent(entity, components.NewTraits(
        championData.Traits,
    ))
    
    // Add champion identity component
    cb.world.AddComponent(entity, components.NewChampionInfo(
        championData.Name,
        championData.Cost,
        starLevel,
    ))
    
    return entity
}