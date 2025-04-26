package simulation

import (
	"fmt"
	"log"
	"reflect"
	"time"

	"github.com/suriz/tft-dps-simulator/components"
	"github.com/suriz/tft-dps-simulator/ecs"
	"github.com/suriz/tft-dps-simulator/systems"
	eventsys "github.com/suriz/tft-dps-simulator/systems/events"
	itemsys "github.com/suriz/tft-dps-simulator/systems/items"
	"github.com/suriz/tft-dps-simulator/utils"
)

// Simulation manages the simulation loop and coordinates system execution
type Simulation struct {
	world    *ecs.World
	eventBus eventsys.EventBus // Interface remains the same
	// Keep system references if needed for setup or direct calls (less common now)
	autoAttackSystem       *systems.AutoAttackSystem
	damageSystem           *systems.DamageSystem
	statCalcSystem         *systems.StatCalculationSystem
	baseStaticItemSystem   *itemsys.BaseStaticItemSystem
	abilityCritSystem      *itemsys.AbilityCritSystem
	dynamicTimeItemSystem  *itemsys.DynamicTimeItemSystem
	spellCastSystem        *systems.SpellCastSystem
	dynamicEventItemSystem *itemsys.DynamicEventItemSystem
	// Add other systems as needed

	config      SimulationConfig
	currentTime float64
	recordQueue []*eventsys.EventItem // No longer needed, managed by eventBus
	// No longer need eventQueue here if it's managed within the eventBus
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

	// Create Event Bus (which now includes the PriorityQueue)
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
	championActionSystem := systems.NewChampionActionSystem(world, eventBus)

	// Register Event Handlers
	eventBus.RegisterHandler(damageSystem)
	eventBus.RegisterHandler(dynamicEventItemSystem)
	eventBus.RegisterHandler(championActionSystem)
	eventBus.RegisterHandler(autoAttackSystem)
	eventBus.RegisterHandler(spellCastSystem)
	// TODO: Register handlers for other systems that react to events (e.g., SpellCastSystem, AutoAttackSystem for recovery events, DynamicTimeItemSystem for timer events)

	sim := &Simulation{
		world:                  world,
		eventBus:               eventBus,         // Store the bus instance
		autoAttackSystem:       autoAttackSystem, // Keep refs for now if needed during transition
		damageSystem:           damageSystem,
		statCalcSystem:         statCalcSystem,
		baseStaticItemSystem:   baseStaticItemSystem,
		abilityCritSystem:      abilityCritSystem,
		dynamicTimeItemSystem:  dynamicTimeItemSystem,
		spellCastSystem:        spellCastSystem,
		dynamicEventItemSystem: dynamicEventItemSystem,
		config:                 config,
		currentTime:            0.0,
		recordQueue: 		  make([]*eventsys.EventItem, 0),
	}

	// apply bonus static item stats to champions AND enqueue initial events
	sim.setupCombat() 
	return sim
}

// setupCombat runs initial setup and enqueues starting events.
func (s *Simulation) setupCombat() {
	fmt.Println("--- Running Initial Combat Setup ---")
	// Apply static bonuses first (devlog.md L281.2)
	s.abilityCritSystem.Update() // For IE/JG check
	s.baseStaticItemSystem.ApplyStats()
	s.statCalcSystem.ApplyStaticBonusStats() // Calculate final stats based on static bonuses

	// TODO: Implement other "before combat" steps from devlog.md (L279)
	// 1. Resolve start-of-combat effects (Items like Thief's Gloves - requires item implementation) (devlog.md L280)
	// 3. Enqueue time effects (e.g., Archangel's) - Requires DynamicTimeItemSystem refactor (devlog.md L282)
	//    Example (conceptual - needs system changes):
	//    s.dynamicTimeItemSystem.EnqueueInitialTimerEvents(s.eventBus)
	// 4. Other special handlings (e.g., Overlord - requires trait implementation) (devlog.md L283)

	// Enqueue first actions for all champions at t=0 (devlog.md L284)
	championInfoType := reflect.TypeOf(components.ChampionInfo{})
	healthType := reflect.TypeOf(components.Health{}) // Ensure we only enqueue for living champions
	champions := s.world.GetEntitiesWithComponents(championInfoType, healthType)
	if s.config.DebugMode {
		log.Printf("Found %d champions to enqueue initial action.", len(champions))
	}
	for _, champ := range champions {
		// Ensure champion is alive before enqueueing action
		if health, ok := s.world.GetHealth(champ); ok && health.GetCurrentHP() > 0 {
			// Enqueue the event for the Action Handler system to process at t=0
			// The handler will decide whether to attack or cast based on state (mana, etc.)
			initialActionEvent := eventsys.ChampionActionEvent{Entity: champ}
			s.eventBus.Enqueue(initialActionEvent, 0.0)
			if s.config.DebugMode {
				log.Printf("Enqueued initial ChampionActionEvent for entity %d at t=0.0", champ)
			}
		} else {
			if s.config.DebugMode {
				log.Printf("Skipping initial ChampionActionEvent for entity %d (not alive or no health component)", champ)
			}
		}
	}
	if s.config.DebugMode {
		log.Println("Initial combat setup complete. Events enqueued:", s.eventBus.(*eventsys.SimpleBus).Len()) // Check queue length
	}
}

// SetMaxTime sets the maximum simulation time in seconds
func (s *Simulation) SetMaxTime(seconds float64) {
	s.config = s.config.WithMaxTime(seconds)
}

// SetTimeStep sets the simulation time step in seconds
func (s *Simulation) SetTimeStep(seconds float64) {
	s.config = s.config.WithTimeStep(seconds)
}

// RunSimulation executes the event-driven simulation loop.
func (s *Simulation) RunSimulation() {
	startTime := time.Now()
	fmt.Println("\nStarting Event-Driven Simulation...")

	s.currentTime = 0.0 // Ensure time starts at 0

	// Get the specific SimpleBus implementation to access Dequeue and Len
	simpleBus, ok := s.eventBus.(*eventsys.SimpleBus)
	if !ok {
		log.Fatal("EventBus is not a *SimpleBus, cannot run simulation")
		return
	}

	// Main event loop
	for simpleBus.Len() > 0 {
		// 1. Dequeue the next event
		eventItem := simpleBus.Dequeue()
		if eventItem == nil { // Should not happen if Len() > 0, but safety check
			break
		}

		// Check for simulation end conditions BEFORE processing
		// Condition 1: Time exceeds MaxTime
		if eventItem.Timestamp > s.config.MaxTime {
			log.Printf("Simulation time (%.3fs) exceeds MaxTime (%.1fs). Stopping.", eventItem.Timestamp, s.config.MaxTime)
			break
		}
		// TODO: Condition 2: One team has no alive champion units (requires health/team check)

		// 2. Set simulation time = evt.Timestamp
		// Only advance time forward. If events are somehow scheduled in the past (shouldn't happen with jitter), log it.
		if eventItem.Timestamp < s.currentTime {
			log.Printf("WARN: Event timestamp %.3fs is before current time %.3fs. Processing anyway.", eventItem.Timestamp, s.currentTime)
		}
		s.currentTime = eventItem.Timestamp

		if s.config.DebugMode {
			log.Printf("[T=%.3fs] Dequeued: %T", s.currentTime, eventItem.Event)
		}

		// 3. Handle the event (Dispatch to registered handlers)
		// Event handlers might enqueue subsequent events.
		s.eventBus.Dispatch(eventItem.Event)

		// 4. Save event to RecordQueue (if implemented)
		s.recordQueue = append(s.recordQueue, eventItem)

	} // End of event loop

	elapsed := time.Since(startTime)
	log.Printf("\nSimulation Ended (Time: %.3fs, Events processed. Real time: %v)\n", s.currentTime, elapsed)
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
