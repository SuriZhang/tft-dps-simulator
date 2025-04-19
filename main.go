package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"reflect"

	"github.com/suriz/tft-dps-simulator/components"
	"github.com/suriz/tft-dps-simulator/data"
	"github.com/suriz/tft-dps-simulator/ecs"
	"github.com/suriz/tft-dps-simulator/factory"
	"github.com/suriz/tft-dps-simulator/managers"
	"github.com/suriz/tft-dps-simulator/simulation"
	"github.com/suriz/tft-dps-simulator/systems"
	eventsys "github.com/suriz/tft-dps-simulator/systems/events"
	itemsys "github.com/suriz/tft-dps-simulator/systems/items"
	"github.com/suriz/tft-dps-simulator/utils"
)

// Helper function to handle AddComponent errors
func addComponentOrLog(world *ecs.World, entity ecs.Entity, component interface{}) {
	err := world.AddComponent(entity, component)
	if err != nil {
		// Decide how to handle: log, panic, etc.
		log.Printf("Error adding component %T to entity %d: %v\n", component, entity, err)
	}
}

func main() {
	// --- Data Loading ---
	dataDir := "./assets"
	fileName := "en_us_14.1b.json"
	filePath := filepath.Join(dataDir, fileName)
	fmt.Println("------------Loading Set Data---------------")
	tftData, err := data.LoadSetDataFromFile(filePath, "TFTSet14")
	if err != nil {
		log.Printf("Error loading set data: %v\n", err)
		os.Exit(1)
	}
	data.InitializeChampions(tftData)
	utils.PrintTftDataLoaded(tftData)

	data.InitializeTraits(tftData)

	fmt.Println("------------Loading Item Data---------------")
	// items, err := data.LoadItemDataFromFile(filepath.Join(dataDir, fileName))
	// if err != nil {
	// 	log.Printf("Error loading item data: %v\n", err)
	// 	os.Exit(1)
	// }

	data.InitializeSetActiveAugments(tftData, filePath)

	data.InitializeSetActiveItems(tftData, filePath)

	data.GetItemByApiName("")

	// --- ECS Setup ---
	world := ecs.NewWorld()
	championFactory := factory.NewChampionFactory(world)
	equipmentManager := managers.NewEquipmentManager(world)
	statCalculationSystem := systems.NewStatCalculationSystem(world)
	abilityCritSystem := itemsys.NewAbilityCritSystem(world)
	baseStaticItemSystem := itemsys.NewBaseStaticItemSystem(world)
	// --- Create Initial Entities (Example) ---
	// This part remains conceptually similar, but uses the helper
	fmt.Println("\n------------Creating Initial Entities---------------")
	kindred, err := championFactory.CreatePlayerChampion("TFT14_Kindred", 1)
	if err != nil {
		log.Printf("Error creating champion: %v\n", err)
		return
	}

	// addComponentOrLog(world, kindred, components.CanAbilityCritFromTraits{})

	err = equipmentManager.AddItemToChampion(kindred, "TFT_Item_InfinityEdge")
	if err != nil {
		log.Printf("Error adding item to Kindred: %v\n", err)
	}
	err = equipmentManager.AddItemToChampion(kindred, "TFT_Item_InfinityEdge")
	if err != nil {
		log.Printf("Error adding item to Kindred: %v\n", err)
	}
	err = equipmentManager.AddItemToChampion(kindred, "TFT_Item_InfinityEdge")
	if err != nil {
		log.Printf("Error adding item to Kindred: %v\n", err)
	}

	// 1. RESET all bonus stats for all relevant entities
	//    (This could be its own small system or done explicitly here)
	healthType := reflect.TypeOf(components.Health{})
	attackType := reflect.TypeOf(components.Attack{})
	manaType := reflect.TypeOf(components.Mana{})
	entitiesToReset := world.GetEntitiesWithComponents(healthType, attackType, manaType) // Or query individually
	for _, entity := range entitiesToReset {
		if health, ok := world.GetHealth(entity); ok {
			health.ResetBonuses()
		}
		if attack, ok := world.GetAttack(entity); ok {
			attack.ResetBonuses()
		}
		if mana, ok := world.GetMana(entity); ok {
			mana.ResetBonuses()
		}
	}

	// process Inifity Edge and Jeweled Gauntlet
	abilityCritSystem.Update()
	// apply item effects
	baseStaticItemSystem.ApplyStats()

	statCalculationSystem.ApplyStaticBonusStats()

	// brand, err := championFactory.CreateAllyChampion("Brand", 1)
	// if err != nil {
	// 	log.Printf("Error creating Brand: %v\n", err)
	// 	return
	// }
	// addComponentOrLog(world, brand, components.NewTeam(0))

	// Print Team 0 champions
	utils.PrintTeamStats(world)

	// // --- Simulation Setup ---
	fmt.Println("\n------------Setting up Simulation---------------")

	// // Create Target Dummy manually
	targetDummy, err := championFactory.CreateEnemyChampion("TFT_TrainingDummy", 1)
	if err != nil {
		log.Printf("Error creating Traning Dummy: %v\n", err)
		return
	}
	addComponentOrLog(world, targetDummy, components.NewHealth(10000, 0, 0))
	addComponentOrLog(world, targetDummy, components.NewPosition(5, 1))
	// Update Attack component as enemy will not attack for MVP1
	addComponentOrLog(world, targetDummy, components.NewAttack(0, 0, 0))

	utils.PrintChampionStats(world, targetDummy)

	// --- Setup Simulation ---
	config := simulation.DefaultConfig().
		WithMaxTime(30.0).
		WithTimeStep(1.0).
		WithDebugMode(true)

	sim := simulation.NewSimulationWithConfig(world, config)

	// Create Event Bus
	eventBus := eventsys.NewSimpleBus()

	// Create Systems
	autoAttackSystem := systems.NewAutoAttackSystem(world, eventBus)
	damageSystem := systems.NewDamageSystem(world, eventBus)
	// ... create other systems ...

	// Register Event Handlers
	eventBus.RegisterHandler(damageSystem)
	// ... register other handlers ...

	// Simulation Loop
	var currentTime float64 = 0.0
	var deltaTime float64 = 0.1 // Example timestep

	for currentTime < 10.0 { // Example duration
		// Update systems that generate events or act based on time
		autoAttackSystem.TriggerAutoAttack(deltaTime)
		// ... update other systems like spell casting, movement etc. ...

		// Process all queued events
		eventBus.ProcessAll()

		// Increment time
		currentTime += deltaTime

		// Check for simulation end conditions (e.g., one team wiped out)
		// ...
	}

	// --- Run Simulation ---
	sim.RunSimulation()

	// --- Print Results ---
	sim.PrintResults(kindred, targetDummy)
}
