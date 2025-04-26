// filepath: systems/action.go
package systems

import (
	"log"

	"github.com/suriz/tft-dps-simulator/ecs"
	eventsys "github.com/suriz/tft-dps-simulator/systems/events"
)

// ChampionActionSystem decides the next action for a champion based on its state.
type ChampionActionSystem struct {
	world    *ecs.World
	eventBus eventsys.EventBus
	// currentTime float64 // System needs to know the current simulation time
}

// NewChampionActionSystem creates a new ChampionActionSystem.
func NewChampionActionSystem(world *ecs.World, bus eventsys.EventBus) *ChampionActionSystem {
	return &ChampionActionSystem{
		world:    world,
		eventBus: bus,
	}
}

// CanHandle checks if the system can process the given event type.
func (s *ChampionActionSystem) CanHandle(evt interface{}) bool {
	switch evt.(type) {
	case eventsys.ChampionActionEvent:
		return true
	default:
		return false
	}
}

// HandleEvent processes events related to action decisions.
func (s *ChampionActionSystem) HandleEvent(evt interface{}) {
	switch event := evt.(type) {
	case eventsys.ChampionActionEvent:
		s.decideNextAction(event.Entity, event.Timestamp)
	}
}

// decideNextAction implements the logic from devlog.md#L237
func (s *ChampionActionSystem) decideNextAction(entity ecs.Entity, currentTime float64) {
	// --- Get necessary components ---
	state, hasState := s.world.GetState(entity)
	mana, hasMana := s.world.GetMana(entity)
	attack, hasAttack := s.world.GetAttack(entity)
	_, hasSpell := s.world.GetSpell(entity)
	health, hasHealth := s.world.GetHealth(entity)

	// Basic checks
	if !hasHealth || health.GetCurrentHP() <= 0 {
		// log.Printf("ActionSystem: Entity %d is dead, skipping action check at %.3fs.", entity, currentTime)
		return // Dead entity takes no actions
	}
	if !hasState || !hasMana || !hasAttack || !hasSpell {
		log.Printf("ActionSystem: Entity %d missing core components (State/Mana/Attack/Spell?) at %.3fs, skipping action.", entity, currentTime)
		return
	}

	// --- Implement Devlog Logic (devlog.md#L237) ---

	// 0. Check if already performing an uninterruptible action (simplification: only check stun for now)
	// More detailed checks needed: IsCasting, IsAttackStartingUp? Devlog implies spell startup has priority over CC (L269) - needs clarification.
	// Let's assume for now only stun stops a *new* action request. Ongoing actions might be handled differently.

	// 1. Check Stun
	if state.GetIsStunned() { // Assumes IsStunned is managed correctly elsewhere
		log.Printf("ActionSystem: Entity %d is stunned at %.3fs, skipping action.", entity, currentTime)
		// TODO: Need a mechanism to re-trigger action check when stun ends (e.g., StunEndEvent)
		return
	}

	// Check if currently busy with another action state that prevents starting a new one
	// (e.g., already casting, already starting up an attack)
	if state.IsCasting || state.IsAttackStartingUp {
		log.Printf("ActionSystem: Entity %d is already busy (%+v) at %.3fs, skipping new action decision.", entity, *state, currentTime)
		return
	}

	// 2. Check Mana Full
	if mana.IsFull() && mana.GetMaxMana() != 0.0 {
		// 2a. Check Attack Recovering (devlog.md#L245)
		if state.GetIsAttackRecovering() {
			log.Printf("ActionSystem: Entity %d has full mana but is recovering from attack at %.3fs, skipping spell.", entity, currentTime)
			// Wait for AttackRecoveryEndEvent to trigger next check
			return
		}

		// 2b. Should Cast Spell
		log.Printf("ActionSystem: Entity %d has full mana, enqueueing SpellCastStartEvent at %.3fs.", entity, currentTime)
		s.eventBus.Enqueue(eventsys.SpellCastStartEvent{Entity: entity, Timestamp: currentTime}, currentTime)
		// The SpellCastSystem will handle this event and update the state to IsCasting
		return // Action decided
	}



	// 3. Check Attack Cooling Down (devlog.md#L251)
	if state.GetIsAttackCoolingDown() {
		log.Printf("ActionSystem: Entity %d is waiting for attack cooldown at %.3fs, skipping attack.", entity, currentTime)
		// Wait for AttackCooldownEndEvent to trigger next check
		return
	}

	if attack.GetBaseAttackSpeed() == 0 || attack.GetFinalAttackSpeed() == 0 {
		log.Printf("ActionSystem: Entity %d has zero attack speed at %.3fs, skipping attack.", entity, currentTime)
		return // No attack speed, can't attack
	}

	// 4. Should Auto Attack
	log.Printf("ActionSystem: Entity %d enqueueing AttackStartEvent at %.3fs.", entity, currentTime)
	s.eventBus.Enqueue(eventsys.AttackStartEvent{Entity: entity, Timestamp: currentTime}, currentTime)
	// The AutoAttackSystem will handle this event and update the state to IsAttackStartingUp
}
