package traithandlers

import (
	"log"
	"math"
	"reflect" // Needed for component type

	"tft-dps-simulator/internal/core/components/traits"
	"tft-dps-simulator/internal/core/data"
	"tft-dps-simulator/internal/core/ecs"
	eventsys "tft-dps-simulator/internal/core/systems/events"
	traitsys "tft-dps-simulator/internal/core/systems/traits"
)

// RapidfireHandler implements dynamic logic for the Rapidfire trait using a dedicated component.
type RapidfireHandler struct{}

// Static check to ensure interface implementation.
var _ traitsys.TraitHandler = (*RapidfireHandler)(nil)

func init() {
	traitsys.RegisterTraitHandler(data.TFT14_Rapidfire, &RapidfireHandler{})
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

	teamChampions := traitsys.GetChampionsByTeam(world, teamID)
	for _, entity := range teamChampions {
		// Apply static team-wide bonus directly
		if teamASBonus != 0 {
			if attack, compOk := world.GetAttack(entity); compOk {
				attack.AddBonusPercentAttackSpeed(teamASBonus)
				log.Printf("RapidfireHandler (Team %d): Applied static +%.1f%% AS for Entity %d", teamID, teamASBonus*100, entity)
			}
		}

		// Add RapidfireEffect component ONLY to champions with the Rapidfire trait
		if traitComp, ok := world.GetTraits(entity); ok && traitComp.HasTrait(data.TFT14_Rapidfire) {
			if !world.HasComponent(entity, reflect.TypeOf(traits.RapidfireEffect{})) {
				rapidfireEffect := traits.NewRapidfireEffect(maxStacks, asPerStack)
				world.AddComponent(entity, rapidfireEffect)
				log.Printf("RapidfireHandler (Team %d): Added RapidfireEffect component to Entity %d", teamID, entity)
			} else {
				log.Printf("RapidfireHandler (Team %d): RapidfireEffect component already exists for Entity %d", teamID, entity)
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
	rapidfireEffect, ok := world.GetRapidfireEffect(entity)
	if !ok {
		// log.Printf("  Rapidfire: Entity %d does not have RapidfireEffect component.", entity)
		return
	}

	// Increment stacks within the component
	_, reachedMax := rapidfireEffect.IncrementStacks()
	if !reachedMax {
		log.Printf("RapidfireHandler: Entity %d attacked. Stacks: %d/%d.",
			entity, rapidfireEffect.GetCurrentStacks(), rapidfireEffect.GetMaxStacks())
		
		attack, ok := world.GetAttack(entity)
		if !ok {
			log.Printf("RapidfireHandler: Entity %d does not have Attack component.", entity)
			return
		}
		attack.AddBonusPercentAttackSpeed(rapidfireEffect.GetAttackSpeedPerStack())

		// --- Enqueue RecalculateStatsEvent ---
		recalcEvent := eventsys.RecalculateStatsEvent{
			Entity:    entity,
			Timestamp: attackEvt.Timestamp,
		}
		// Enqueue immediately at the current time
		eventBus.Enqueue(recalcEvent, attackEvt.Timestamp)
		log.Printf("RapidfireHandler: Enqueued RecalculateStatsEvent for Entity %d at t=%.4f", entity, attackEvt.Timestamp)
	} else {
		log.Printf("RapidfireHandler: Entity %d reached max stacks (%d). No further stack increment.",
			entity, rapidfireEffect.GetMaxStacks())
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
	rapidfireEffectType := reflect.TypeOf(traits.RapidfireEffect{})
	entities := world.GetEntitiesWithComponents(rapidfireEffectType)
	for _, entity := range entities {
		world.RemoveComponent(entity, rapidfireEffectType)
		log.Printf("  Rapidfire: Removed RapidfireEffect component from Entity %d during reset.", entity)
	}
	// Static bonuses are assumed to be reset by the main stat reset mechanism before applying new ones.
}