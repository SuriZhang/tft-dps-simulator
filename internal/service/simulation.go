package service

import (
	"fmt"
	"log"
	"time"

	"tft-dps-simulator/internal/core/data"
	"tft-dps-simulator/internal/core/ecs"
	"tft-dps-simulator/internal/core/entity"
	"tft-dps-simulator/internal/core/factory"
	"tft-dps-simulator/internal/core/managers"
	"tft-dps-simulator/internal/core/simulation"
)

// SimulationService handles the logic for running combat simulations.
type SimulationService struct {
	loadedData *data.TFTSetData // Access to game data (champions, items, etc.)
}

// NewSimulationService creates a new SimulationService.
func NewSimulationService(loadedData *data.TFTSetData) *SimulationService {
	return &SimulationService{
		loadedData: loadedData,
	}
}

// RunSimulation executes a combat simulation based on the provided request.
func (s *SimulationService) RunSimulation(requestChampions []BoardChampion) (*RunSimulationResponse, error) {
	log.Println("Starting simulation run...")
	startTime := time.Now()

	// 1. Initialize ECS world
	world := ecs.NewWorld()

	// 2. Initialize Factory and Managers
	// Pass world instead of loadedData based on test setup
	championFactory := factory.NewChampionFactory(world)
	equipmentManager := managers.NewEquipmentManager(world)

	// Map to link request champion ID (ApiName) to ECS entity ID
	entityMap := make(map[entity.Entity]string)
	// var championEntities []entity.Entity

	// 3. Create Champion Entities from Request
	log.Printf("Processing %d requested champions...", len(requestChampions))
	for _, reqChamp := range requestChampions {
		// Use CreatePlayerChampion based on tests
		entityID, err := championFactory.CreatePlayerChampion(reqChamp.ApiName, reqChamp.Stars)
		if err != nil {
			log.Printf("Error creating champion entity %s: %v. Skipping.", reqChamp.ApiName, err)
			continue // Or return error
		}
		entityMap[entityID] = reqChamp.ApiName // Store the mapping

		log.Printf("Items for champion %s: %v", reqChamp.ApiName, reqChamp.Items)

		// Add items using AddItemToChampion based on tests
		for _, itemReq := range reqChamp.Items {

			// Need the item's API name (string) for AddItemToChampion
			err := equipmentManager.AddItemToChampion(entityID, itemReq.ApiName)
			if err != nil {
				log.Printf("Error equipping item %s to champion %s: %v. Skipping item.", itemReq.ApiName, reqChamp.ApiName, err)
				// Decide if this should be a fatal error for the request
			}
		}

		log.Printf("Created entity %d for champion %s at (%d, %d) with %d items", entityID, reqChamp.ApiName, reqChamp.Position.Row, reqChamp.Position.Col, len(reqChamp.Items))
	}

	// 4. Add Target Dummies (Team 1)
	log.Println("Adding enemy training dummies")
	targetDummy, err := championFactory.CreateEnemyChampion("TFT_TrainingDummy", 3)
	if err != nil {
		log.Printf("Error creating target dummy: %v", err)
		return nil, fmt.Errorf("error creating target dummy: %w", err)
	}

	dummyHealth, ok := world.GetHealth(targetDummy)
	if !ok {
		log.Printf("Error getting health component for target dummy: %v", err)
		return nil, fmt.Errorf("error getting health component for target dummy: %w", err)
	}
	dummyHealth.SetBaseMaxHP(1000000) // Set dummy health to 1000 for testing
	dummyHealth.SetBaseMR(0.0)      // Set dummy MR to 0 for testing
	dummyHealth.SetBaseArmor(0.0)   // Set dummy Armor to 0 for testing
	dummyAttack, ok := world.GetAttack(targetDummy)
	if !ok {
		log.Printf("Error getting attack component for target dummy: %v", err)
		return nil, fmt.Errorf("error getting attack component for target dummy: %w", err)
	}
	dummyAttack.SetBaseAttackSpeed(0.0)

	// 5. Configure and Run Simulation
	log.Println("Configuring simulation...")
	config := simulation.DefaultConfig()

	// Validate config
	if err := config.Validate(); err != nil {
		log.Printf("Invalid simulation config: %v", err)
		return nil, fmt.Errorf("invalid simulation config: %w", err)
	}

	// Instantiate simulation using NewSimulationWithConfig based on tests
	sim := simulation.NewSimulationWithConfig(world, config)

	log.Println("Running simulation...")
	sim.RunSimulation()

	log.Println("Simulation finished.")

	archievedEvents := make([]ArchivedEvent, 0, len(sim.GetArchiveEvents()))

	for _, event := range sim.GetArchiveEvents() {
		archivedEvent := ArchivedEvent{
			EventItem: *event,
			EventType: fmt.Sprintf("%T", event.Event),
		}
		archievedEvents = append(archievedEvents, archivedEvent)
	}

	// for _, event := range archievedEvents {
	// 	log.Printf("%s: %+v", event.EventType, event.EventItem)
	// }

	// 6. Process Results
	log.Println("Processing simulation results")
	// Use service types instead of server types
	results := []ChampionSimulationResult{}

	for entityID, apiName := range entityMap {
		// Fetch final health to check if alive (example of reading state post-simulation)
		// Use helper functions like in tests if available, otherwise direct component access
		healthComp, healthOk := world.GetHealth(entityID)
		isAlive := false
		if healthOk {
			isAlive = healthComp.GetCurrentHP() > 0
		}
		log.Printf("Post-simulation check: Champion %s (Entity %d) Alive: %t", apiName, entityID, isAlive)

		attackComp, attackOk := world.GetAttack(entityID)
		attackCount := 0
		if attackOk {
			attackCount = attackComp.GetAttackCount() // Assuming GetAttackCount exists
		}
		// Example: Fetch spell casts
		spellComp, spellOk := world.GetSpell(entityID)
		spellCastCount := 0
		if spellOk {
			spellCastCount = spellComp.GetCastCount()
		}

		damageStats, dsOK := world.GetDamageStats(entityID)
		if !dsOK {
			log.Printf("Error getting damage stats for champion %s (Entity %d)", apiName, entityID)
		}

		damageStats.TotalAutoAttackCounts = attackCount
		damageStats.TotalSpellCastCounts = spellCastCount
		damageStats.DamagePerSecond = damageStats.TotalDamage / config.MaxTime

		// Use service types
		results = append(results, ChampionSimulationResult{
			ChampionApiName: apiName,
			ChampionEntityID: entityID,
			DamageStats:     *damageStats,
		})
	}

	response := &RunSimulationResponse{
		Results:        results,
		ArchieveEvents: archievedEvents, // Assign the dereferenced slice
	}

	elapsed := time.Since(startTime)
	log.Printf("Simulation request processed successfully in %s.", elapsed)
	return response, nil
}

