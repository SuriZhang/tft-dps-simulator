package traitsys

import (
	"log"
	"reflect"

	// Correct import path for trait components
	traitcomps "github.com/suriz/tft-dps-simulator/components/traits"
	"github.com/suriz/tft-dps-simulator/data"
	"github.com/suriz/tft-dps-simulator/ecs"
	eventsys "github.com/suriz/tft-dps-simulator/systems/events"
)

// DynamicEventTraitSystem manages the activation, event handling, and deactivation of dynamic traits via components.
// It assumes trait tiers are determined once at combat start and do not change mid-combat.
// It handles trait effects triggered by game events.
type DynamicEventTraitSystem struct {
    world      *ecs.World
    traitState *TeamTraitState
    eventBus   eventsys.EventBus
}

// NewDynamicEventTraitSystem creates a new DynamicEventTraitSystem.
func NewDynamicEventTraitSystem(world *ecs.World, state *TeamTraitState, bus eventsys.EventBus) *DynamicEventTraitSystem {
    return &DynamicEventTraitSystem{
        world:      world,
        traitState: state,
        eventBus:   bus, // Store the event bus
    }
}

// ActivateTraits performs the initial activation of dynamic traits based on the calculated tiers.
// Should run once in setupCombat after TraitCounterSystem updates tiers.
func (s *DynamicEventTraitSystem) ActivateTraits() {
    log.Println("DynamicEventTraitSystem: Activating dynamic traits...")

    currentActiveTiers := s.traitState.activeTier // Get current state

    // Check for activations based on the initial calculation
    for teamID, currentTiers := range currentActiveTiers {
        for traitName, currentTierIndex := range currentTiers {
            if currentTierIndex == -1 {
                continue // Not active currently
            }

            handler, exists := GetTraitHandler(traitName)
            if !exists {
				log.Printf("Warning: No handler found for trait '%s'", traitName)
                continue
            }

            // Get trait data and call OnActivate
            traitData, exists := data.Traits[traitName]
            if exists && currentTierIndex < len(traitData.Effects) {
                activeEffect := traitData.Effects[currentTierIndex]
                log.Printf("  Activating dynamic event trait '%s' for Team %d (TierIndex %d)", traitName, teamID, currentTierIndex)
                // OnActivate is responsible for adding components or applying initial effects
                // Pass eventBus if OnActivate needs it (optional, depends on trait needs)
                handler.OnActivate(teamID, activeEffect, s.world)
            } else {
                log.Printf("Error: Invalid data for activating dynamic trait '%s' (TierIndex %d)", traitName, currentTierIndex)
            }
        }
    }

    log.Println("DynamicEventTraitSystem: Finished activating dynamic traits.")
}

// ResetAllTraits calls the Reset method on all registered dynamic trait handlers.
func (s *DynamicEventTraitSystem) ResetAllTraits() {
    log.Println("DynamicEventTraitSystem: Resetting all dynamic trait states...")
    for traitApiName, handler := range TraitRegistry {
        log.Printf("  Resetting state for trait '%s'", traitApiName)
        handler.Reset(s.world)
    }
    log.Println("DynamicEventTraitSystem: Finished resetting all dynamic trait states.")
}

// HandleEvent dispatches incoming game events to relevant active dynamic trait handlers.
func (s *DynamicEventTraitSystem) HandleEvent(event interface{}) {
    involvedEntities := s.determineInvolvedEntities(event)
    if len(involvedEntities) == 0 {
        return
    }

    processedHandlers := make(map[string]struct{})

    for _, entity := range involvedEntities {
        for traitApiName, handler := range TraitRegistry {
            // Check if the entity has the component associated with this trait.
            // Example for Rapidfire:
            if traitApiName == "TFT14_Swift" { // TODO: Replace with a better mapping mechanism
                // Use the correct component type name 'RapidfireEffect'
                if s.world.HasComponent(entity, reflect.TypeOf(traitcomps.RapidfireEffect{})) {
                    if _, done := processedHandlers[traitApiName]; !done {
                        // Pass the event bus to the handler
                        handler.Handle(event, entity, s.world, s.eventBus)
                        processedHandlers[traitApiName] = struct{}{}
                    }
                }
            }
            // Add checks for other dynamic trait components...
        }
    }
}

// determineInvolvedEntities extracts entities from an event.
func (s *DynamicEventTraitSystem) determineInvolvedEntities(event interface{}) []ecs.Entity {
    entities := make(map[ecs.Entity]struct{}) // Use map for uniqueness

    addEntity := func(entity ecs.Entity) {
        if entity != 0 {
            entities[entity] = struct{}{}
        }
    }

    switch evt := event.(type) {
    case eventsys.AttackLandedEvent:
        addEntity(evt.Source)
    default:
        return nil
    }

    result := make([]ecs.Entity, 0, len(entities))
    for entity := range entities {
        result = append(result, entity)
    }
    return result
}

// CanHandle checks if the system should process this event type.
func (s *DynamicEventTraitSystem) CanHandle(evt interface{}) bool {
    switch evt.(type) {
    case eventsys.AttackLandedEvent:
        return true
    default:
        return false
    }
}