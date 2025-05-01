package traitsys

import (
	"log"
	"math"
	"reflect" // Needed for component type

	"github.com/suriz/tft-dps-simulator/components/traits"
	"github.com/suriz/tft-dps-simulator/data"
	"github.com/suriz/tft-dps-simulator/ecs"
	eventsys "github.com/suriz/tft-dps-simulator/systems/events" // Keep for AttackLandedEvent
)

// RapidfireHandler implements dynamic logic for the Rapidfire trait using a dedicated component.
type RapidfireHandler struct{}

// Static check to ensure interface implementation.
var _ TraitHandler = (*RapidfireHandler)(nil)

func init() {
    RegisterTraitHandler(data.Rapidfire, &RapidfireHandler{})
}

// OnActivate adds the RapidfireEffect component to champions with the trait and applies static team bonus.
func (h *RapidfireHandler) OnActivate(teamID int, effect data.Effect, world *ecs.World) {
    log.Printf("RapidfireHandler: Activating for Team %d (Style %d)", teamID, effect.Style)

    teamASBonus, okAS := effect.Variables["{b6739a03}"] // Mapped from "{b6739a03}"
    asPerStack, okStackAS := effect.Variables["AttackSpeed"]
    maxStacksFloat, okMax := effect.Variables["MaxStacks"]

    if !okAS || !okStackAS || !okMax {
        log.Printf("Warning: Rapidfire (Team %d) missing required variables in effect data.", teamID)
        return
    }
    maxStacks := int(math.Round(maxStacksFloat))

    teamChampions := getChampionsByTeam(world, teamID)
    for _, entity := range teamChampions {
        // Apply static team-wide bonus directly
        if teamASBonus != 0 {
            if attack, compOk := world.GetAttack(entity); compOk {
                attack.AddBonusPercentAttackSpeed(teamASBonus)
                log.Printf("  Rapidfire (Team %d): Applied static +%.1f%% AS for Entity %d", teamID, teamASBonus*100, entity)
            }
        }

        // Add RapidfireEffect component ONLY to champions with the Rapidfire trait
        if traitComp, ok := world.GetTraits(entity); ok && traitComp.HasTrait(data.Rapidfire) {
            if !world.HasComponent(entity, reflect.TypeOf(traits.RapidfireEffect{})) {
                newState := traits.NewRapidfireEffect(maxStacks, asPerStack)
                world.AddComponent(entity, newState)
                log.Printf("  Rapidfire (Team %d): Added RapidfireEffect component to Entity %d", teamID, entity)
            } else {
                log.Printf("  Rapidfire (Team %d): RapidfireEffect component already exists for Entity %d", teamID, entity)
            }
        }
    }
}

// Handle processes AttackLandedEvent for entities with RapidfireEffect component.
// Now receives eventBus to enqueue RecalculateStatsEvent.
func (h *RapidfireHandler) Handle(event interface{}, entity ecs.Entity, world *ecs.World, eventBus eventsys.EventBus) {
    attackEvt, ok := event.(eventsys.AttackLandedEvent)
    if !ok || attackEvt.Source != entity {
        return // Not an attack landed event from this entity
    }

    // Get the RapidfireEffect component
    stateCompType := reflect.TypeOf(traits.RapidfireEffect{})
    stateComp, ok := world.GetComponent(entity, stateCompType)
    if !ok {
        return // Entity doesn't have the component
    }
    rapidfireState := stateComp.(*traits.RapidfireEffect)

    // Increment stacks within the component
    if rapidfireState.IncrementStacks() {
        log.Printf("  Rapidfire: Entity %d attacked. Stacks: %d/%d.",
            entity, rapidfireState.CurrentStacks, rapidfireState.MaxStacks)

        // --- Enqueue RecalculateStatsEvent ---
        // The bonus application happens when StatCalculationSystem reads the component state.
        // We need to trigger that recalculation now.
        recalcEvent := eventsys.RecalculateStatsEvent{
            Entity:    entity,
            Timestamp: attackEvt.Timestamp, // Use the timestamp of the triggering event
        }
        // Enqueue immediately at the current time
        eventBus.Enqueue(recalcEvent, attackEvt.Timestamp)
        log.Printf("  Rapidfire: Enqueued RecalculateStatsEvent for Entity %d at t=%.4f", entity, attackEvt.Timestamp)
    }
}

// OnDeactivate removes the RapidfireEffect component and reverses static bonus.
func (h *RapidfireHandler) OnDeactivate(teamID int, effect data.Effect, world *ecs.World) {
    log.Printf("RapidfireHandler: Deactivating for Team %d", teamID)
	// TODO: Deactivate static bonus for team
}

// Reset removes RapidfireEffect components from all entities.
func (h *RapidfireHandler) Reset(world *ecs.World) {
    log.Printf("RapidfireHandler: Resetting all states.")
    rapidfireStateType := reflect.TypeOf(traits.RapidfireEffect{})
    entities := world.GetEntitiesWithComponents(rapidfireStateType)
    for _, entity := range entities {
        world.RemoveComponent(entity, rapidfireStateType)
        log.Printf("  Rapidfire: Removed RapidfireEffect component from Entity %d during reset.", entity)
    }
    // Static bonuses are assumed to be reset by the main stat reset mechanism before applying new ones.
}