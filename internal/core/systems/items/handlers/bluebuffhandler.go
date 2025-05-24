package itemhandlers

import (
	"log"

	"tft-dps-simulator/internal/core/data"
	"tft-dps-simulator/internal/core/ecs"
	eventsys "tft-dps-simulator/internal/core/systems/events"
	itemsys "tft-dps-simulator/internal/core/systems/items"
)

// BlueBuffHandler handler implements the Blue Buff item effects
type BlueBuffHandler struct{}

func init() {
	itemsys.RegisterItemHandler(data.TFT_Item_BlueBuff, &BlueBuffHandler{})
}

// OnEquip implements itemsys.ItemHandler
func (h *BlueBuffHandler) OnEquip(entity ecs.Entity, world *ecs.World, eventBus eventsys.EventBus) {
	log.Printf("BlueBuff: OnEquip for entity %d", entity)

	// Reset the effect state when equipped
	if blueBuff, exists := world.GetBlueBuffEffect(entity); exists {
		blueBuff.DeactivateAmplification()
		log.Printf("BlueBuff: Reset Blue Buff state for entity %d on equip.", entity)
	} else {
		log.Printf("BlueBuff: WARNING - BlueBuff component not found for entity %d on equip.", entity)
	}
}

// GetHandledEvents returns the list of events handled by this handler
func (h *BlueBuffHandler) GetHandledEvents() []string {
	return []string{
		"SpellLandedEvent",
		"KillEvent",
		"AssistEvent",
		"BlueBuffDamageAmpActivateEvent",
		"BlueBuffDamageAmpDeactivateEvent", // Handle deactivation events
	}
}

// ProcessEvent implements itemsys.ItemHandler
func (h *BlueBuffHandler) ProcessEvent(event interface{}, entity ecs.Entity, world *ecs.World, eventBus eventsys.EventBus) {
	switch evt := event.(type) {
	case eventsys.SpellLandedEvent:
		if evt.Source == entity {
			h.handleSpellLanded(evt, entity, world)
		}
	case eventsys.KillEvent:
		if evt.Killer == entity {
			h.enqueueTakedownActivation(entity, evt.Timestamp, world, eventBus)
		}
	case eventsys.AssistEvent:
		if evt.Assistor == entity {
			h.enqueueTakedownActivation(entity, evt.Timestamp, world, eventBus)
		}
	case eventsys.BlueBuffDamageAmpActivateEvent:
		if evt.Entity == entity {
			h.handleDamageAmpActivation(evt, entity, world, eventBus)
		}
	case eventsys.BlueBuffDamageAmpDeactivateEvent:
		if evt.Entity == entity {
			h.handleDamageAmpDeactivation(evt, entity, world, eventBus)
		}
	}
}

// enqueueTakedownActivation enqueues the activation event when takedown occurs
func (h *BlueBuffHandler) enqueueTakedownActivation(entity ecs.Entity, timestamp float64, world *ecs.World, eventBus eventsys.EventBus) {
	// Verify entity still has Blue Buff
	equipment, okEq := world.GetEquipment(entity)
	if !okEq || equipment.GetItemCount(data.TFT_Item_BlueBuff) == 0 {
		log.Printf("BlueBuff (enqueueTakedownActivation): Entity %d no longer has Blue Buff at %.3fs.", entity, timestamp)
		return
	}

	// Enqueue activation event immediately
	activationEvent := eventsys.BlueBuffDamageAmpActivateEvent{
		Entity:    entity,
		Timestamp: timestamp,
	}
	eventBus.Enqueue(activationEvent, timestamp)
	log.Printf("BlueBuff (enqueueTakedownActivation): Enqueued BlueBuffDamageAmpActivateEvent for entity %d at %.3fs", entity, timestamp)
}

// handleSpellLanded grants mana after casting
func (h *BlueBuffHandler) handleSpellLanded(evt eventsys.SpellLandedEvent, entity ecs.Entity, world *ecs.World) {
	// Verify entity still has Blue Buff
	equipment, okEq := world.GetEquipment(entity)
	if !okEq || equipment.GetItemCount(data.TFT_Item_BlueBuff) == 0 {
		log.Printf("BlueBuff (handleSpellLanded): Entity %d no longer has Blue Buff at %.3fs.", entity, evt.Timestamp)
		return
	}

	blueBuff, okBB := world.GetBlueBuffEffect(entity)
	if !okBB {
		log.Printf("BlueBuff (handleSpellLanded): Entity %d missing BlueBuff component at %.3fs.", entity, evt.Timestamp)
		return
	}

	mana, okMana := world.GetMana(entity)
	if !okMana {
		log.Printf("BlueBuff (handleSpellLanded): Entity %d missing Mana component at %.3fs.", entity, evt.Timestamp)
		return
	}

	// Grant mana refund
	initialMana := mana.GetCurrentMana()
	mana.AddCurrentMana(blueBuff.GetManaRefund())

	log.Printf("BlueBuff (handleSpellLanded): Entity %d gained %.1f mana after casting at %.3fs (%.1f -> %.1f / %.1f)",
		entity, blueBuff.GetManaRefund(), evt.Timestamp, initialMana, mana.GetCurrentMana(), mana.GetMaxMana())
}

// handleDamageAmpActivation processes the activation event and manages the effect
func (h *BlueBuffHandler) handleDamageAmpActivation(evt eventsys.BlueBuffDamageAmpActivateEvent, entity ecs.Entity, world *ecs.World, eventBus eventsys.EventBus) {
	// Verify entity still has Blue Buff
	equipment, okEq := world.GetEquipment(entity)
	if !okEq || equipment.GetItemCount(data.TFT_Item_BlueBuff) == 0 {
		log.Printf("BlueBuff (handleDamageAmpActivation): Entity %d no longer has Blue Buff at %.3fs.", entity, evt.Timestamp)
		return
	}

	blueBuff, okBB := world.GetBlueBuffEffect(entity)
	if !okBB {
		log.Printf("BlueBuff (handleDamageAmpActivation): Entity %d missing BlueBuff component at %.3fs.", entity, evt.Timestamp)
		return
	}

	attack, okAttack := world.GetAttack(entity)
	if !okAttack {
		log.Printf("BlueBuff (handleDamageAmpActivation): Entity %d missing Attack component at %.3fs.", entity, evt.Timestamp)
		return
	}

	// Check if already active - this handles timer resets
	wasActive := blueBuff.IsActive(evt.Timestamp)

	// If already active, we need to cancel the pending deactivation event
	// We'll use sequence numbers to handle this
	blueBuff.IncremenetCurrentActivationSequence()

	// Always activate/reset the amplification
	blueBuff.ActivateAmplification(evt.Timestamp)

	// If not already active, add the damage amplification to stats
	if !wasActive {
		attack.AddBonusDamageAmp(blueBuff.DamageAmp)
		log.Printf("BlueBuff (handleDamageAmpActivation): Entity %d activated damage amplification at %.3fs (%.1f%% for %.1fs)",
			entity, evt.Timestamp, blueBuff.DamageAmp*100, blueBuff.TakedownTimer)
	} else {
		log.Printf("BlueBuff (handleDamageAmpActivation): Entity %d refreshed damage amplification at %.3fs (%.1f%% for %.1fs)",
			entity, evt.Timestamp, blueBuff.DamageAmp*100, blueBuff.TakedownTimer)
	}

	// Enqueue RecalculateStatsEvent to update champion stats
	recalcEvent := eventsys.RecalculateStatsEvent{Entity: entity, Timestamp: evt.Timestamp}
	eventBus.Enqueue(recalcEvent, evt.Timestamp)

	// Enqueue deactivation event with sequence number
	deactivationEvent := eventsys.BlueBuffDamageAmpDeactivateEvent{
		Entity:              entity,
		Timestamp:           evt.Timestamp + blueBuff.TakedownTimer,
		ActivationSequence: blueBuff.CurrentActivationSequence,
	}
	eventBus.Enqueue(deactivationEvent, evt.Timestamp+blueBuff.TakedownTimer)
	log.Printf("BlueBuff (handleDamageAmpActivation): Enqueued BlueBuffDamageAmpDeactivateEvent for entity %d at %.3fs (sequence %d)",
		entity, evt.Timestamp+blueBuff.TakedownTimer, blueBuff.CurrentActivationSequence)
}

// Add handleDamageAmpDeactivation method
func (h *BlueBuffHandler) handleDamageAmpDeactivation(evt eventsys.BlueBuffDamageAmpDeactivateEvent, entity ecs.Entity, world *ecs.World, eventBus eventsys.EventBus) {
	// Verify entity still has Blue Buff
	equipment, okEq := world.GetEquipment(entity)
	if !okEq || equipment.GetItemCount(data.TFT_Item_BlueBuff) == 0 {
		log.Printf("BlueBuff (handleDamageAmpDeactivation): Entity %d no longer has Blue Buff at %.3fs.", entity, evt.Timestamp)
		return
	}

	blueBuff, okBB := world.GetBlueBuffEffect(entity)
	if !okBB {
		log.Printf("BlueBuff (handleDamageAmpDeactivation): Entity %d missing BlueBuff component at %.3fs.", entity, evt.Timestamp)
		return
	}

	attack, okAttack := world.GetAttack(entity)
	if !okAttack {
		log.Printf("BlueBuff (handleDamageAmpDeactivation): Entity %d missing Attack component at %.3fs.", entity, evt.Timestamp)
		return
	}

	// Only deactivate if this event matches the current activation sequence
	if evt.ActivationSequence == blueBuff.CurrentActivationSequence {
		// Remove the damage amplification from stats
		attack.AddBonusDamageAmp(-blueBuff.DamageAmp)
		blueBuff.DeactivateAmplification()

		log.Printf("BlueBuff (handleDamageAmpDeactivation): Entity %d deactivated damage amplification at %.3fs",
			entity, evt.Timestamp)

		// Enqueue RecalculateStatsEvent to update champion stats
		recalcEvent := eventsys.RecalculateStatsEvent{Entity: entity, Timestamp: evt.Timestamp}
		eventBus.Enqueue(recalcEvent, evt.Timestamp)
	} else {
		log.Printf("BlueBuff (handleDamageAmpDeactivation): Ignoring stale deactivation event for entity %d (sequence %d vs current %d)",
			entity, evt.ActivationSequence, blueBuff.CurrentActivationSequence)
	}
}
