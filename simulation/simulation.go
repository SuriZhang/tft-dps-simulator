package simulation

import (
	"fmt"
	"time"

	"github.com/suriz/tft-dps-simulator/ecs"
	"github.com/suriz/tft-dps-simulator/systems"
	itemsys "github.com/suriz/tft-dps-simulator/systems/items"
	"github.com/suriz/tft-dps-simulator/utils"
)

// Simulation manages the simulation loop and coordinates system execution
type Simulation struct {
    world               *ecs.World
    autoAttackSystem    *systems.AutoAttackSystem
    statCalcSystem      *systems.StatCalculationSystem
    baseStaticItemSystem *itemsys.BaseStaticItemSystem
    abilityCritSystem   *itemsys.AbilityCritSystem
    // Add other systems as needed
    
    config              SimulationConfig
    currentTime         float64
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

    sim := &Simulation{
        world:               world,
        autoAttackSystem:    systems.NewAutoAttackSystem(world),
        statCalcSystem:      systems.NewStatCalculationSystem(world),
        baseStaticItemSystem: itemsys.NewBaseStaticItemSystem(world),
        abilityCritSystem:   itemsys.NewAbilityCritSystem(world),
        config:              config,
        currentTime:         0.0,
    }
    
    return sim
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
    
    // Apply initial static calculations
    s.abilityCritSystem.Update()
    s.baseStaticItemSystem.ApplyStats()
    s.statCalcSystem.Update()
    
    // Main simulation loop
    var nextReportTime float64 = s.config.ReportingInterval
    
    for s.currentTime < s.config.MaxTime {
        s.Step()
        
        // Status reporting at intervals if enabled
        if s.config.DebugMode && s.currentTime >= nextReportTime {
            fmt.Printf("Simulation time: %.1fs\n", s.currentTime)
            nextReportTime += s.config.ReportingInterval
        }
    }
    
    elapsed := time.Since(startTime)
    fmt.Printf("\nSimulation Ended (Reached %.1fs, took %v real time)\n", s.currentTime, elapsed)
}

// Step advances the simulation by one time step
func (s *Simulation) Step() {
    // Only run systems that are enabled in the config
    if s.config.EnableAutoAttacks {
        s.autoAttackSystem.Update(s.config.TimeStep)
    }
    
    if s.config.EnableSpellCasts {
        // Add spell cast system when implemented
        // s.spellCastSystem.Update(s.config.TimeStep)
    }
    
    // Add other conditional system updates based on config flags
    
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