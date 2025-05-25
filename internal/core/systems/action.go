// filepath: systems/action.go
package systems

import (
	"log"

	"tft-dps-simulator/internal/core/components"
	"tft-dps-simulator/internal/core/ecs"
	"tft-dps-simulator/internal/core/entity"
	eventsys "tft-dps-simulator/internal/core/systems/events"
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
func (s *ChampionActionSystem) decideNextAction(entity entity.Entity, currentTime float64) {
	// --- Get necessary components ---
	state, hasState := s.world.GetState(entity)
	mana, hasMana := s.world.GetMana(entity)
	attack, hasAttack := s.world.GetAttack(entity)
	spell, hasSpell := s.world.GetSpell(entity)
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
	state.StartActionCheck(currentTime)

	// --- Implement Devlog Logic (devlog.md#L237) ---

	// 0. Check if already performing an uninterruptible action (simplification: only check stun for now)
	// More detailed checks needed: IsCasting, IsAttackStartingUp? Devlog implies spell startup has priority over CC (L269) - needs clarification.
	// Let's assume for now only stun stops a *new* action request. Ongoing actions might be handled differently.

	// 1. Check Stun
	if state.IsStunned { // Assumes IsStunned is managed correctly elsewhere
		log.Printf("ActionSystem: Entity %d is stunned at %.3fs, skipping action.", entity, currentTime)
		// TODO: Need a mechanism to re-trigger action check when stun ends (e.g., StunEndEvent)
		return
	}

	// 2. Check first action (attack or spell)
	if (state.CurrentState == components.Idle && state.PreviousState == components.Idle && currentTime == 0.0) || state.PreviousState == components.AttackCoolingDown {
		// Check if mana is full and spell is available
		if mana.CanCastSpell() {
			log.Printf("ActionSystem: Entity %d casting spell at %.3fs.", entity, currentTime)
			state.StartCast(currentTime, spell.GetCastStartUp() + spell.GetCastRecovery()) 
			s.eventBus.Enqueue(eventsys.SpellCastCycleStartEvent{Entity: entity, Timestamp: currentTime}, currentTime)
			return
		} else {
			if (attack.GetFinalAttackSpeed() > 0.0) {
			log.Printf("ActionSystem: Entity %d attacking at %.3fs.", entity, currentTime)

			state.StartAttack(currentTime, attack.GetCurrentAttackStartup())
			s.eventBus.Enqueue(eventsys.AttackStartupEvent{Entity: entity, Timestamp: currentTime}, currentTime)
			}
			return
		}
	}

	// 3. Check if previousState is AttackRecovering
	if state.PreviousState == components.AttackRecovering && state.CurrentState == components.Idle {
		// Check if mana is full and spell is available
		if mana.CanCastSpell() {
			log.Printf("ActionSystem: Entity %d casting spell at %.3fs.", entity, currentTime)
			state.StartCast(currentTime, spell.GetCastStartUp() + spell.GetCastRecovery()) 
			s.eventBus.Enqueue(eventsys.SpellCastCycleStartEvent{Entity: entity, Timestamp: currentTime}, currentTime)
			return
		} else {
			log.Printf("ActionSystem: Entity %d start attack coolingdown at %.3fs.", entity, currentTime)
			state.StartAttackCooldown(currentTime, attack.GetCurrentAttackCooldown())
			s.eventBus.Enqueue(eventsys.AttackCooldownStartEvent{Entity: entity, Timestamp: currentTime}, currentTime)
			return
		}
	}

	// 4. Check if previousState is Casting
	if state.PreviousState == components.Casting && state.CurrentState == components.Idle {
		// Check if attack cooldown is over
		if state.PreviousActionDuration <= attack.GetCurrentAttackCooldown() {
			log.Printf("DEBUG: state.ActionDuration: %.3fs, attack cooldown: %.3fs", state.ActionDuration, attack.GetCurrentAttackCooldown())
			log.Printf("ActionSystem: Entity %d start attack coolingdown at %.3fs.", entity, currentTime)

			
			remainingCooldown := attack.GetCurrentAttackCooldown() - state.PreviousActionDuration
			log.Printf("DEBUG: Attack cooldown remaining: %.3fs, normal attack cooldown: %.3fs", remainingCooldown, attack.GetCurrentAttackCooldown())

			state.StartAttackCooldown(currentTime, remainingCooldown)
			s.eventBus.Enqueue(eventsys.AttackCooldownStartEvent{Entity: entity, Timestamp: currentTime}, currentTime)
			return
		} else {
			log.Printf("ActionSystem: Entity %d start attack at %.3fs.", entity, currentTime)
			state.StartAttack(currentTime, attack.GetCurrentAttackStartup())
			s.eventBus.Enqueue(eventsys.AttackStartupEvent{Entity: entity, Timestamp: currentTime}, currentTime)
			return
		}
	}

}
