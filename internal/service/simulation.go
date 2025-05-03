package service

import (
	"fmt"
	"log"
	"time" // Import time package

	"tft-dps-simulator/internal/core/components"
	"tft-dps-simulator/internal/core/data"
	"tft-dps-simulator/internal/core/ecs"
	"tft-dps-simulator/internal/core/factory"
	"tft-dps-simulator/internal/core/managers"
	"tft-dps-simulator/internal/core/simulation"
	// "tft-dps-simulator/internal/core/systems" // Systems are managed internally by simulation now
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
	// traitManager := managers.NewTraitManager(world) // Trait manager might be needed if not handled internally by simulation setup

	// Map to link request champion ID (ApiName) to ECS entity ID
	entityMap := make(map[string]ecs.Entity)
	var championEntities []ecs.Entity

	// 3. Create Champion Entities from Request
	log.Printf("Processing %d requested champions...", len(requestChampions))
	for _, reqChamp := range requestChampions {
		// Use CreatePlayerChampion based on tests
		entityID, err := championFactory.CreatePlayerChampion(reqChamp.ApiName, reqChamp.Stars)
		if err != nil {
			log.Printf("Error creating champion entity %s: %v. Skipping.", reqChamp.ApiName, err)
			continue // Or return error
		}
		entityMap[reqChamp.ApiName] = entityID
		championEntities = append(championEntities, entityID)

		// Add items using AddItemToChampion based on tests
		for _, itemReq := range reqChamp.Items {
			// Need the item's API name (string) for AddItemToChampion
			err := equipmentManager.AddItemToChampion(entityID, itemReq.ApiName)
			if err != nil {
				log.Printf("Error equipping item %s to champion %s: %v. Skipping item.", itemReq.ApiName, reqChamp.ApiName, err)
				// Decide if this should be a fatal error for the request
			}
		}

		// Set position and team (Team 0 for player champions)
		// Use AddComponent with component instance based on tests
		world.AddComponent(entityID, components.NewPosition(reqChamp.Position.Row, reqChamp.Position.Col))
		world.AddComponent(entityID, components.NewTeam(0)) // Team 0 for player

		log.Printf("Created entity %d for champion %s at (%d, %d) with %d items", entityID, reqChamp.ApiName, reqChamp.Position.Row, reqChamp.Position.Col, len(reqChamp.Items))
	}

	// 4. Add Target Dummies (Team 1)
	log.Println("Adding enemy training dummies")
	targetDummy, err := championFactory.CreateEnemyChampion("TFT_TrainingDummy", 3)
	if err != nil {
		log.Printf("Error creating target dummy: %v", err)
		return nil, fmt.Errorf("error creating target dummy: %w", err)
	}
	world.AddComponent(targetDummy, components.NewPosition(0, 0)) // Dummy position
	targetDummyHealth, _ := world.GetHealth(targetDummy)
	targetDummyHealth.SetBaseMaxHP(10000) // Set dummy health (example value)

	// 5. Configure and Run Simulation
	log.Println("Configuring simulation...")
	simulationDuration := 30.0 // seconds
	// Use DefaultConfig and builder methods based on config.go and tests
	cfg := simulation.DefaultConfig().
		WithMaxTime(simulationDuration).
		WithDebugMode(false) // Set debug mode as needed

	// Validate config
	if err := cfg.Validate(); err != nil {
		log.Printf("Invalid simulation config: %v", err)
		return nil, fmt.Errorf("invalid simulation config: %w", err)
	}

	// Instantiate simulation using NewSimulationWithConfig based on tests
	sim := simulation.NewSimulationWithConfig(world, cfg)

	log.Println("Running simulation...")
	// Call RunSimulation based on tests
	sim.RunSimulation() // RunSimulation doesn't return error or stats directly
	// No error check needed here unless RunSimulation signature changes

	log.Println("Simulation finished.")

	// 6. Process Results (Placeholder - requires fetching data from world components/events)
	log.Println("Processing simulation results (using placeholders)...")
	// Use service types instead of server types
	results := []ChampionSimulationResult{}

	// TODO: Replace this placeholder loop with processing of actual simulation results.
	// This involves querying components (like Attack, Spell, Health, potentially new stats components)
	// from the world *after* the simulation has run.
	// Example: Get total damage dealt by querying a DamageDealt component or summing DamageAppliedEvents.
	// Example: Get attack count from the Attack component.
	// Example: Get spell casts from the Spell component.

	for apiName, entityID := range entityMap {
		reqChamp := findRequestChampion(requestChampions, apiName) // Helper to get stars etc.
		stars := 1
		if reqChamp != nil {
			stars = reqChamp.Stars
		}

		// Fetch final health to check if alive (example of reading state post-simulation)
		// Use helper functions like in tests if available, otherwise direct component access
		healthComp, healthOk := world.GetHealth(entityID)
		isAlive := false
		if healthOk {
			isAlive = healthComp.GetCurrentHP() > 0
		}
		log.Printf("Post-simulation check: Champion %s (Entity %d) Alive: %t", apiName, entityID, isAlive)

		// *** Replace with actual stats fetched from world components/events ***
		// Example: Fetch attack count
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
		// Example: Fetch damage stats (requires a damage tracking mechanism)
		// This part is highly dependent on how damage is tracked (e.g., events, stats component)
		// Placeholder values remain for now.
		placeholderStats := DamageStats{
			TotalDamage:           1000 * float64(stars) * (float64(entityID%5 + 1)), // Fake variation
			DamagePerSecond:       (1000 * float64(stars) * (float64(entityID%5 + 1))) / simulationDuration,
			TotalADDamage:         600 * float64(stars),
			TotalAPDamage:         300 * float64(stars),
			TotalTrueDamage:       100 * float64(stars),
			TotalAutoAttackCounts: attackCount, // Use fetched value
			TotalSpellCastCounts:  spellCastCount, // Use fetched value
		}
		// *** End Placeholder Stats ***

		// Use service types
		results = append(results, ChampionSimulationResult{
			ChampionApiName:  apiName,
			DamageStats: placeholderStats,
		})
	}

	// 7. Build and Return Response
	// Use service types
	response := &RunSimulationResponse{
		Results: results,
	}

	elapsed := time.Since(startTime)
	log.Printf("Simulation request processed successfully in %s.", elapsed)
	return response, nil
}

// Helper function to find the original request data for a champion by ApiName
// Use service types
func findRequestChampion(requestChampions []BoardChampion, apiName string) *BoardChampion {
	for i := range requestChampions {
		if requestChampions[i].ApiName == apiName {
			return &requestChampions[i]
		}
	}
	return nil
}
