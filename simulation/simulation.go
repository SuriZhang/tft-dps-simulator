package simulation

import (
	"fmt"
	"log"
	"time"

	"github.com/suriz/tft-dps-simulator/ecs"
	"github.com/suriz/tft-dps-simulator/systems"
	eventsys "github.com/suriz/tft-dps-simulator/systems/events"
	itemsys "github.com/suriz/tft-dps-simulator/systems/items"
	"github.com/suriz/tft-dps-simulator/utils"
)

// Simulation manages the simulation loop and coordinates system execution
type Simulation struct {
	world                 *ecs.World
	eventBus              eventsys.EventBus
	autoAttackSystem      *systems.AutoAttackSystem
	damageSystem          *systems.DamageSystem
	statCalcSystem        *systems.StatCalculationSystem
	baseStaticItemSystem  *itemsys.BaseStaticItemSystem
	abilityCritSystem     *itemsys.AbilityCritSystem
	dynamicTimeItemSystem *itemsys.DynamicTimeItemSystem
	spellCastSystem 	 *systems.SpellCastSystem
	dynamicEventItemSystem *itemsys.DynamicEventItemSystem
	// Add other systems as needed

	config      SimulationConfig
	currentTime float64
}

// NewSimulation creates a new simulation with the given world and default config
func NewSimulation(world *ecs.World) *Simulation {
	return NewSimulationWithConfig(world, DefaultConfig())
}

// NewSimulationWithConfig creates a new simulation with the given world and config
func NewSimulationWithConfig(world *ecs.World, config SimulationConfig) *Simulation {
	if err := config.Validate(); err != nil {
		panic(fmt.Sprintf("Invalid simulation config: %v", err))
	}

	// Create Event Bus
	eventBus := eventsys.NewSimpleBus()

	// Create Systems, passing event bus where needed
	autoAttackSystem := systems.NewAutoAttackSystem(world, eventBus)
	damageSystem := systems.NewDamageSystem(world, eventBus)
	statCalcSystem := systems.NewStatCalculationSystem(world)
	baseStaticItemSystem := itemsys.NewBaseStaticItemSystem(world)
	abilityCritSystem := itemsys.NewAbilityCritSystem(world)
	dynamicTimeItemSystem := itemsys.NewDynamicTimeItemSystem(world)
	spellCastSystem := systems.NewSpellCastSystem(world, eventBus)
	dynamicEventItemSystem := itemsys.NewDynamicEventItemSystem(world, eventBus)

	// Register Event Handlers
	eventBus.RegisterHandler(damageSystem)
	eventBus.RegisterHandler(dynamicEventItemSystem)
	// Register other handlers here...

	sim := &Simulation{
		world:                 world,
		eventBus:              eventBus,
		autoAttackSystem:      autoAttackSystem,
		damageSystem:          damageSystem,
		statCalcSystem:        statCalcSystem,
		baseStaticItemSystem:  baseStaticItemSystem,
		abilityCritSystem:     abilityCritSystem,
		dynamicTimeItemSystem: dynamicTimeItemSystem,
		spellCastSystem:       spellCastSystem,
		dynamicEventItemSystem: dynamicEventItemSystem,
		config:                config,
		currentTime:           0.0,
	}

	// apply bonus static item stats to champions
	sim.applyInitialUpdates() 
	return sim
}

// applyInitialUpdates runs systems that need to execute once before the main loop.
func (s *Simulation) applyInitialUpdates() {
	// Order might matter here depending on dependencies
	s.abilityCritSystem.Update()
	s.baseStaticItemSystem.ApplyStats()
	s.statCalcSystem.ApplyStaticBonusStats()
	s.eventBus.ProcessAll()
}

// SetMaxTime sets the maximum simulation time in seconds
func (s *Simulation) SetMaxTime(seconds float64) {
	s.config = s.config.WithMaxTime(seconds)
}

// SetTimeStep sets the simulation time step in seconds
func (s *Simulation) SetTimeStep(seconds float64) {
	s.config = s.config.WithTimeStep(seconds)
}

// RunSimulation executes the simulation until completion
func (s *Simulation) RunSimulation() {
	startTime := time.Now()
	fmt.Println("\nStarting Simulation...")

	// Reset time
	s.currentTime = 0.0

	// Main simulation loop
	var nextReportTime float64 = s.config.ReportingInterval

	for s.currentTime < s.config.MaxTime {
		s.Step() // Use the updated Step method

		// Status reporting at intervals if enabled
		if s.config.DebugMode && s.currentTime >= nextReportTime {
			log.Printf("Simulation time: %.1fs\n", s.currentTime)
			// Optionally print specific stats here if needed during the run
			nextReportTime += s.config.ReportingInterval
		}

		// TODO: Add termination conditions (e.g., check if one team is defeated)
		// For now, we just run until MaxTime
	}

	elapsed := time.Since(startTime)
	log.Printf("\nSimulation Ended (Reached %.1fs, took %v real time)\n", s.currentTime, elapsed)
}

// Step advances the simulation by one time step
func (s *Simulation) Step() {
	// 1. Update systems that generate events or act based on time
	if s.config.EnableAutoAttacks {
		s.autoAttackSystem.TriggerAutoAttack(s.config.TimeStep)
	}

	if s.config.EnableSpellCasts {
		// Add spell cast system update when implemented
		// s.spellCastSystem.Update(s.config.TimeStep)
	}

	// Update dynamic time items (e.g., Archangel's stacking, QS expiry check)
	s.dynamicTimeItemSystem.Update(s.config.TimeStep)

	// Update stat calculation system AFTER dynamic items might have changed stats
	// This ensures FinalStats reflect changes from items in the current step
	s.statCalcSystem.Update(s.config.TimeStep)

	// Add other time-dependent system updates here...
	// e.g., systems handling buffs/debuffs duration

	// 2. Process all events queued in this step (e.g., AttackLanded -> DamageApplied)
	s.eventBus.ProcessAll()

	// 3. Increment time (Do this LAST in the step)
	s.currentTime += s.config.TimeStep
}

// PrintResults displays the final simulation results
func (s *Simulation) PrintResults(entities ...ecs.Entity) {
	fmt.Println("\n------------Final Stats---------------")
	for _, entity := range entities {
		utils.PrintChampionStats(s.world, entity)
	}
}

// GetConfig returns a copy of the current simulation configuration
func (s *Simulation) GetConfig() SimulationConfig {
	return s.config
}

// SetConfig sets a new configuration for the simulation
func (s *Simulation) SetConfig(config SimulationConfig) error {
	if err := config.Validate(); err != nil {
		return err
	}
	s.config = config
	return nil
}
