package main

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"

	"github.com/suriz/tft-dps-simulator/components"
	"github.com/suriz/tft-dps-simulator/data"
	"github.com/suriz/tft-dps-simulator/ecs"
	"github.com/suriz/tft-dps-simulator/factory"
)

func main() {
    // Get the path to your data files
    dataDir := "./data/data_files" // Or use a configurable path
    filePath := filepath.Join(dataDir, "en_us_14.1b.json")
    
	fmt.Println("------------Loading Set Data---------------")
    // Load TFT set data with specified mutator
    setData, err := data.LoadSetDataFromFile(filePath, "TFTSet14")
    if err != nil {
        fmt.Printf("Error loading set data: %v\n", err)
        os.Exit(1)
    }
    
    fmt.Printf("Successfully loaded TFT set: %s\n", setData.SetData[0].Name)
    fmt.Printf("Found %d champions in the set\n", len(setData.SetData[0].Champions))
	fmt.Printf("Found %d traits in the set\n", len(setData.SetData[0].Traits))
	fmt.Printf("Found %d items in the set\n", len(setData.SetData[0].Items))
	fmt.Printf("Found %d augments in the set\n", len(setData.SetData[0].Augments))
	fmt.Println("------------Finished Loading Set Data---------------")

	
    
    // Create ECS world
    world := ecs.NewWorld()
    
    // Create champion factory
    championBuilder := factory.NewChampionBuilder(world)
    
    // Create a team of champions for simulation
    teamEntities := make([]ecs.Entity, 0)
    
    // Add some champions to the team (example: first 5 champions at 2 stars)
    for i, champData := range setData.SetData[0].Champions {
        if i >= 5 {
            break
        }
        
        // Create champion entity at 2 stars
        entity := championBuilder.CreateChampion(champData, 2)
        teamEntities = append(teamEntities, entity)
        
        // Get Identity component to print info
        idComponent, _ := world.GetComponent(entity, reflect.TypeOf(components.ChampionInfo{}))
        id := idComponent.(components.ChampionInfo)
        
        // Get Health component to print info
        healthComponent, _ := world.GetComponent(entity, reflect.TypeOf(components.Health{}))
        health := healthComponent.(components.Health)
        
        fmt.Printf("Created champion: %s (★%d) - HP: %.0f\n", 
            id.Name, id.StarLevel, health.Max)
    }
    
    fmt.Printf("Created %d champion entities\n", len(teamEntities))
	// printing the team entities
	for _, entity := range teamEntities {
		idComponent, _ := world.GetComponent(entity, reflect.TypeOf(components.ChampionInfo{}))
		id := idComponent.(components.ChampionInfo)	
		fmt.Printf("Entity ID: %d, Champion Name: %s\n", entity, id.Name)
	}
    
    // Now we can run simulations with these entities
    // TODO: Add simulation logic
}