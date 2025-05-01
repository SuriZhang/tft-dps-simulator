package utils

import (
	"fmt"
	"reflect"

	"tft-dps-simulator/internal/core/components"
	"tft-dps-simulator/internal/core/data"
	"tft-dps-simulator/internal/core/ecs"
)

// PrintChampionStats prints detailed information about a champion entity using type-safe getters.
func PrintChampionStats(world *ecs.World, entity ecs.Entity) {
	// Use type-safe getters and check the 'ok' value
	info, okInfo := world.GetChampionInfo(entity)
	health, okHealth := world.GetHealth(entity)
	mana, okMana := world.GetMana(entity)
	attack, okAttack := world.GetAttack(entity)
	traits, okTraits := world.GetTraits(entity)
	position, okPos := world.GetPosition(entity)
	team, okTeam := world.GetTeam(entity)

	// Check if essential info is present
	if !okInfo {
		fmt.Printf("\n=== Entity ID: %d (No ChampionInfo) ===\n", entity)
		fmt.Println("------------------------------")
		return
	}

	fmt.Printf("\n=== Champion: %s (★ %d) ===\n", info.Name, info.StarLevel)
	fmt.Printf("Entity ID: %d\n", entity)

	// Print components if they exist
	fmt.Printf("Info Component:\n  %+v\n", *info) // Dereference pointer for printing value

	if okHealth {
		// Use the String() method implicitly via %v
		fmt.Printf("Health Component:\n%v\n", health) // String() method already includes newlines
	} else {
		fmt.Println("Health Component: <Missing>")
	}

	if okMana {
		// Use %+v for detailed struct view
		fmt.Printf("Mana Component:\n  %+v\n", *mana)
	} else {
		fmt.Println("Mana Component: <Missing>")
	}

	if okAttack {
		// Use the String() method implicitly via %v
		fmt.Printf("Attack Component:\n%v\n", attack) // String() method already includes newlines
	} else {
		fmt.Println("Attack Component: <Missing>")
	}

	if okTraits {
		// Use %+v for detailed struct view
		fmt.Printf("Traits Component:\n  %+v\n", *traits)
	} else {
		fmt.Println("Traits Component: <Missing>")
	}

	if okPos {
		// Use %+v for detailed struct view
		fmt.Printf("Position Component:\n  %+v\n", *position)
	} else {
		fmt.Println("Position Component: <Missing>")
	}

	if okTeam {
		fmt.Printf("Team Component:\n%v\n", team)
	} else {
		fmt.Println("Team Component: <Missing>")
	}

	fmt.Printf("------------------------------\n")
}

func PrintTeamStats(world *ecs.World) {
	for _, entity := range world.GetEntitiesWithComponents(reflect.TypeOf(components.Team{})) {
		team, ok := world.GetTeam(entity)
		if ok && team.ID == 0 {
			PrintChampionStats(world, entity)
		}
	}
}

func PrintTftDataLoaded(tftData *data.TFTSetData) {
	if tftData == nil || len(tftData.SetData) == 0 {
		fmt.Println("No TFT set data loaded.")
		return
	}
	setInfo := tftData.SetData[0] // Assuming only one set is loaded
	fmt.Printf("Loaded Set: %s\n", setInfo.Mutator)
	fmt.Printf("  Champions: %d\n", len(setInfo.Champions))
	fmt.Printf("  Traits: %d\n", len(setInfo.Traits))
	fmt.Printf("  Items: %d\n", len(setInfo.SetItems))
	fmt.Printf("  Augments: %d\n", len(setInfo.SetAugments))
	fmt.Println("-------------------------------------------")
}

func PrintAugmentStats(augment *data.Item) {
	if augment == nil {
		fmt.Println("No augment data available.")
		return
	}
	fmt.Printf("Augment: %s\n", augment.Name)
	fmt.Printf("  Description: %s\n", augment.Desc)
	fmt.Printf("  AssociatedTraits: %s\n", augment.AssociatedTraits)
	fmt.Println("-------------------------------------------")
}

func PrintItemStats(item *data.Item) {
	if item == nil {
		fmt.Println("No item data available.")
		return
	}
	fmt.Printf("Item: %s\n", item.Name)
	fmt.Printf("  Description: %s\n", item.Desc)
	fmt.Printf("  AssociatedTraits: %s\n", item.AssociatedTraits)
	fmt.Printf("  Composition: %+v\n", item.Composition)
	// print effects
	for _, effect := range item.Effects {
		fmt.Printf("  Effect: %+v\n", effect)
	}
	fmt.Println("-------------------------------------------")
}

func PrintTraitStats(trait *data.Trait) {
	if trait == nil {
		fmt.Println("No trait data available.")
		return
	}
	fmt.Printf("Trait: %s\n", trait.ApiName)
	fmt.Printf("  Description: %s\n", trait.Desc)
	fmt.Printf("  Name: %s\n", trait.Name)
	fmt.Printf("  Icon: %s\n", trait.Icon)
	fmt.Printf("  Effect: %v\n", trait.Effects)
	fmt.Println("-------------------------------------------")
}
