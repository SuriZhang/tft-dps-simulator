package itemsys

import (
	"log"
	"reflect"

	"github.com/suriz/tft-dps-simulator/components"
	"github.com/suriz/tft-dps-simulator/data"
	"github.com/suriz/tft-dps-simulator/ecs"
	eventsys "github.com/suriz/tft-dps-simulator/systems/events"
)

// DynamicTimeItemSystem handles items whose effects change over time via events.
type DynamicTimeItemSystem struct {
	world    *ecs.World
	eventBus eventsys.EventBus // Store the bus to enqueue subsequent events
}

// NewDynamicTimeItemSystem creates a new system instance.
func NewDynamicTimeItemSystem(world *ecs.World, bus eventsys.EventBus) *DynamicTimeItemSystem { // Added bus parameter
	return &DynamicTimeItemSystem{
		world:    world,
		eventBus: bus, // Store the bus
	}
}

// EnqueueInitialEvents checks equipped items and schedules the first timer events.
// This should be called during simulation setup (e.g., in Simulation.setupCombat).
func (s *DynamicTimeItemSystem) EnqueueInitialEvents() {
	log.Println("DynamicTimeItemSystem: Enqueueing initial timer events...")
	entities := s.world.GetEntitiesWithComponents(reflect.TypeOf(components.Equipment{}))

	for _, entity := range entities {
		equipment, ok := s.world.GetEquipment(entity)
		if !ok {
			continue
		}

		// --- Archangel's Staff ---
		if equipment.HasItem(data.TFT_Item_ArchangelsStaff) {
			effect, exists := s.world.GetArchangelsEffect(entity)

			if exists && effect.GetInterval() > 0 {
				effect.ResetEffects()
				firstTickTime := effect.GetInterval() // First tick happens after interval
				tickEvent := eventsys.ArchangelsTickEvent{Entity: entity, Timestamp: firstTickTime}
				s.eventBus.Enqueue(tickEvent, firstTickTime)
				log.Printf("  Enqueued initial ArchangelsTickEvent for entity %d at t=%.3fs", entity, firstTickTime)
			}
		}

		// --- Quicksilver ---
		if equipment.HasItem(data.TFT_Item_Quicksilver) {
			effect, exists := s.world.GetQuicksilverEffect(entity)
			if exists && effect.GetProcInterval() > 0 && effect.GetSpellShieldDuration() > 0 {
				effect.ResetEffects()
				// Enqueue first proc event
				firstProcTime := effect.GetProcInterval()
				if firstProcTime <= effect.GetSpellShieldDuration() { // Only if first proc happens before expiry
					procEvent := eventsys.QuicksilverProcEvent{Entity: entity, Timestamp: firstProcTime}
					s.eventBus.Enqueue(procEvent, firstProcTime)
					log.Printf("  Enqueued initial QuicksilverProcEvent for entity %d at t=%.3fs", entity, firstProcTime)
				}

				// Enqueue the end event
				endTime := effect.GetSpellShieldDuration() // Duration starts at t=0
				endEvent := eventsys.QuicksilverEndEvent{Entity: entity, Timestamp: endTime}
				s.eventBus.Enqueue(endEvent, endTime)
				log.Printf("  Enqueued QuicksilverEndEvent for entity %d at t=%.3fs", entity, endTime)
			}
		}
		// --- Add other items here ---
	}
}

// CanHandle checks if the system can process the given event type.
func (s *DynamicTimeItemSystem) CanHandle(evt interface{}) bool {
	switch evt.(type) {
	case eventsys.ArchangelsTickEvent,
		eventsys.QuicksilverProcEvent,
		eventsys.QuicksilverEndEvent:
		return true
	default:
		return false
	}
}

// HandleEvent processes incoming timer events.
func (s *DynamicTimeItemSystem) HandleEvent(evt interface{}) {
	switch event := evt.(type) {
	case eventsys.ArchangelsTickEvent:
		s.handleArchangelsTick(event)
	case eventsys.QuicksilverProcEvent:
		s.handleQuicksilverProc(event)
	case eventsys.QuicksilverEndEvent:
		s.handleQuicksilverEnd(event)
	}
}

// handleArchangelsTick applies AP bonus and enqueues the next tick.
func (s *DynamicTimeItemSystem) handleArchangelsTick(evt eventsys.ArchangelsTickEvent) {
	entity := evt.Entity
	currentTime := evt.Timestamp

	// Check if entity still exists and has the item/effect
	equipment, equipOk := s.world.GetEquipment(entity)
	effect, effectOk := s.world.GetArchangelsEffect(entity)
	spellComp, spellOk := s.world.GetSpell(entity)
	health, healthOk := s.world.GetHealth(entity) // Check if alive

	if !equipOk || !effectOk || !spellOk || !healthOk || health.GetCurrentHP() <= 0 || !equipment.HasItem(data.TFT_Item_ArchangelsStaff) {
		log.Printf("DynamicTimeItemSystem (ArchangelsTick): Entity %d no longer valid or item removed at %.3fs. Stopping ticks.", entity, currentTime)
		return
	}

	// Apply effect (add bonus AP)
	archangelsCount := equipment.GetItemCount(data.TFT_Item_ArchangelsStaff)
	apGain := effect.GetAPPerInterval() * float64(archangelsCount)
	spellComp.AddBonusAP(apGain)
	effect.AddStacks(1)

	log.Printf("DynamicTimeItemSystem (ArchangelsTick): Entity %d gained %.1f AP at %.3fs (Stacks: %d, Count: %d). Total Bonus AP: %.1f",
		entity, apGain, currentTime, effect.GetStacks(), archangelsCount, spellComp.GetBonusAP())

	// Enqueue event to recalculate stats
	recalcEvent := eventsys.RecalculateStatsEvent{Entity: entity, Timestamp: currentTime}
	s.eventBus.Enqueue(recalcEvent, currentTime)

	// Enqueue the next tick event
	nextTickTime := currentTime + effect.GetInterval()
	nextTickEvent := eventsys.ArchangelsTickEvent{Entity: entity, Timestamp: nextTickTime}
	s.eventBus.Enqueue(nextTickEvent, nextTickTime)
}

// handleQuicksilverProc applies AS bonus and enqueues the next proc if applicable.
func (s *DynamicTimeItemSystem) handleQuicksilverProc(evt eventsys.QuicksilverProcEvent) {
	entity := evt.Entity
	currentTime := evt.Timestamp

	// Check if entity still exists and has the item/effect and is active
	equipment, equipOk := s.world.GetEquipment(entity)
	effect, effectOk := s.world.GetQuicksilverEffect(entity)
	attackComp, attackOk := s.world.GetAttack(entity)
	health, healthOk := s.world.GetHealth(entity) // Check if alive

	// Crucially, also check if the effect is still active (duration hasn't ended)
	if !equipOk || !effectOk || !attackOk || !healthOk || health.GetCurrentHP() <= 0 || !equipment.HasItem(data.TFT_Item_Quicksilver) || !effect.IsActive() {
		log.Printf("DynamicTimeItemSystem (QuicksilverProc): Entity %d no longer valid, item removed, or effect inactive at %.3fs. Stopping procs.", entity, currentTime)
		return
	}

	// Apply effect (add bonus AS)
	quicksilverCount := equipment.GetItemCount(data.TFT_Item_Quicksilver) // Should always be 1 due to uniqueness, but check anyway
	asGain := effect.GetProcAttackSpeed() * float64(quicksilverCount)
	attackComp.AddBonusPercentAttackSpeed(asGain)
	effect.AddStacks(1) // Increment internal stack count

	log.Printf("DynamicTimeItemSystem (QuicksilverProc): Entity %d gained %.2f%% AS at %.3fs (Stacks: %d). Total Bonus AS: %.2f%%",
		entity, asGain*100, currentTime, effect.GetStacks(), attackComp.GetBonusPercentAttackSpeed()*100)

	// Enqueue event to recalculate stats
	recalcEvent := eventsys.RecalculateStatsEvent{Entity: entity, Timestamp: currentTime}
	s.eventBus.Enqueue(recalcEvent, currentTime)

	// Enqueue the next proc event ONLY if it happens before the effect expires
	nextProcTime := currentTime + effect.GetProcInterval()
	// Use the original duration stored in the effect, don't rely on RemainingDuration which might change
	expiryTime := effect.GetSpellShieldDuration()
	if nextProcTime <= expiryTime {
		nextProcEvent := eventsys.QuicksilverProcEvent{Entity: entity, Timestamp: nextProcTime}
		s.eventBus.Enqueue(nextProcEvent, nextProcTime)
	} else {
		log.Printf("DynamicTimeItemSystem (QuicksilverProc): Next proc for entity %d at %.3fs would be after expiry (%.3fs). Not enqueueing.", entity, nextProcTime, expiryTime)
	}
}

// handleQuicksilverEnd marks the effect as inactive.
func (s *DynamicTimeItemSystem) handleQuicksilverEnd(evt eventsys.QuicksilverEndEvent) {
	entity := evt.Entity
	currentTime := evt.Timestamp

	effect, effectOk := s.world.GetQuicksilverEffect(entity)
	if !effectOk {
		// Effect might have been removed if item was removed earlier
		return
	}

	if effect.IsActive() {
		log.Printf("DynamicTimeItemSystem (QuicksilverEnd): Entity %d Quicksilver duration ended at %.3fs. Marking inactive.", entity, currentTime)
		effect.SetIsActive(false)
		// Note: Bonus AS is NOT removed here. It persists but stops stacking.
		// If removal is desired, EquipmentManager should handle it on item removal.
	}
}
