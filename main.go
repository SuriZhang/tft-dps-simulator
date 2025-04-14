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

	"github.com/suriz/tft-dps-simulator/managers"
	"github.com/suriz/tft-dps-simulator/systems"
	itemsys "github.com/suriz/tft-dps-simulator/systems/items"
	"github.com/suriz/tft-dps-simulator/utils"
)

// Helper function to handle AddComponent errors
func addComponentOrLog(world *ecs.World, entity ecs.Entity, component interface{}) {
	err := world.AddComponent(entity, component)
	if err != nil {
		// Decide how to handle: log, panic, etc.
		fmt.Printf("Error adding component %T to entity %d: %v\n", component, entity, err)
	}
}

func main() {
	// --- Data Loading ---
	dataDir := "./data/data_files"
	fileName := "en_us_14.1b.json"
	filePath := filepath.Join(dataDir, fileName)
	fmt.Println("------------Loading Set Data---------------")
	tftData, err := data.LoadSetDataFromFile(filePath, "TFTSet14")
	if err != nil {
		fmt.Printf("Error loading set data: %v\n", err)
		os.Exit(1)
	}
	data.InitializeChampions(tftData)
	utils.PrintTftDataLoaded(tftData)

	data.InitializeTraits(tftData)

	fmt.Println("------------Loading Item Data---------------")
	// items, err := data.LoadItemDataFromFile(filepath.Join(dataDir, fileName))
	// if err != nil {
	// 	fmt.Printf("Error loading item data: %v\n", err)
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
		fmt.Printf("Error creating champion: %v\n", err)
		return
	}

	addComponentOrLog(world, kindred, components.CanAbilityCritFromTraits{})

	err = equipmentManager.AddItemToChampion(kindred, "TFT_Item_InfinityEdge")
	if err != nil {
		fmt.Printf("Error adding item to Kindred: %v\n", err)
	}

	err = equipmentManager.AddItemToChampion(kindred, "TFT_Item_Deathblade")
	if err != nil {
		fmt.Printf("Error adding item to Kindred: %v\n", err)
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

	statCalculationSystem.Update()

	// brand, err := championFactory.CreateAllyChampion("Brand", 1)
	// if err != nil {
	// 	fmt.Printf("Error creating Brand: %v\n", err)
	// 	return
	// }
	// addComponentOrLog(world, brand, components.NewTeam(0))

	// Print Team 0 champions
	utils.PrintTeamStats(world)

	// // --- Simulation Setup ---
	fmt.Println("\n------------Setting up Simulation---------------")

	// // Add position using the helper
	// addComponentOrLog(world, voidspawn, components.NewPosition(1, 1))

	// health, ok := world.GetHealth(voidspawn)
	// if ok {
	// 	health.UpdateCurrentHealth(1000)
	// 	health.UpdateMaxHealth(1000)
	// } else {
	// 	fmt.Println("voidspawn Health component not found.")
	// }
	// attack, ok := world.GetAttack(voidspawn)
	// if ok {
	// 	attack.UpdateDamage(100)
	// } else {
	// 	fmt.Println("voidspawn Attack component not found.")
	// }

	// // Create Target Dummy manually
	// targetDummy, err := championFactory.CreateEnemyChampion("Training Dummy", 1)
	// if err != nil {
	// 	fmt.Printf("Error creating Traning Dummy: %v\n", err)
	// 	return
	// }
	// addComponentOrLog(world, targetDummy, components.NewHealth(10000, 0, 0))
	// addComponentOrLog(world, targetDummy, components.NewPosition(5, 1))
	// // Update Attack component as enemy will not attack for MVP1
	// addComponentOrLog(world, targetDummy, components.NewAttack(0, 0, 0, 0, 0))

	// utils.PrintChampionStats(world, targetDummy)

	// fmt.Println("Created Voidspawn and Enemy Dummy")

	// // --- Instantiate Systems ---
	// autoAttackSystem := systems.NewAutoAttackSystem(world)

	// // --- Run Simulation (remains the same logic) ---
	// fmt.Println("\nStarting Auto Attack Simulation (30s)...")
	// const maxTimeSeconds = 30.0
	// const timeStepSeconds = 1.0
	// elapsedTime := 0.0
	// for elapsedTime < maxTimeSeconds {
	// 	autoAttackSystem.Update(timeStepSeconds)
	// 	elapsedTime += timeStepSeconds
	// }
	// fmt.Printf("\nSimulation Ended (Reached %.1fs)...\n", elapsedTime)

	// // --- Final Status ---
	// fmt.Println("\n------------Final Stats---------------")
	// // Assuming it has been refactored:
	// utils.PrintChampionStats(world, voidspawn)

	// // Use type-safe getters for final status check
	// targetHealth, okHealth := world.GetHealth(targetDummy)
	// targetInfo, okInfo := world.GetChampionInfo(targetDummy)
	// if okHealth && okInfo {
	// 	fmt.Printf("  Name: %s, Current Health: %.1f / %.1f\n", targetInfo.Name, targetHealth.Current, targetHealth.Max)
	// } else {
	// 	fmt.Println("  Could not retrieve final dummy stats (missing components?).")
	// }
}
