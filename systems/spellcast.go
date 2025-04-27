package systems

import (
	"log"

	"github.com/suriz/tft-dps-simulator/ecs"
	eventsys "github.com/suriz/tft-dps-simulator/systems/events"
	"github.com/suriz/tft-dps-simulator/utils"
)

// SpellCastSystem handles spell casting based on events.
type SpellCastSystem struct {
	world    *ecs.World
	eventBus eventsys.EventBus
}

// NewSpellCastSystem creates a new spell cast system.
func NewSpellCastSystem(world *ecs.World, bus eventsys.EventBus) *SpellCastSystem {
	return &SpellCastSystem{
		world:    world,
		eventBus: bus,
	}
}

// CanHandle checks if the system can process the given event type.
func (s *SpellCastSystem) CanHandle(evt interface{}) bool {
	switch evt.(type) {
	case eventsys.SpellCastCycleStartEvent, eventsys.SpellLandedEvent, eventsys.SpellRecoveryEndEvent: // Add more as needed
		return true
	default:
		return false
	}
}

// HandleEvent processes events related to the spell cast cycle.
func (s *SpellCastSystem) HandleEvent(evt interface{}) {
	switch event := evt.(type) {
	case eventsys.SpellCastCycleStartEvent:
		s.handleSpellCastStart(event)
	case eventsys.SpellLandedEvent:
		s.handleSpellLanded(event)
    case eventsys.SpellRecoveryEndEvent:
        s.handleSpellRecoveryEnd(event)
	// TODO: Add handler for SpellRecoveryEndEvent
	}
}

// handleSpellCastStart initiates the spell cast.
func (s *SpellCastSystem) handleSpellCastStart(evt eventsys.SpellCastCycleStartEvent) {
	caster := evt.Entity
	currentTime := evt.Timestamp

	// --- Get Components ---
	_, okState := s.world.GetState(caster)
	spell, okSpell := s.world.GetSpell(caster)
	mana, okMana := s.world.GetMana(caster)
	health, okHealth := s.world.GetHealth(caster)
	// Targeting components if needed
	_, okPos := s.world.GetPosition(caster)
	team, okTeam := s.world.GetTeam(caster)

	if !okState || !okSpell || !okMana || !okHealth || health.GetCurrentHP() <= 0 || !okPos || !okTeam {
		log.Printf("SpellCastSystem (Start): Entity %d missing components or dead at %.3fs. Canceling spell start.", caster, currentTime)
		return
	}

	// --- Targeting (Example: Assume spell targets nearest enemy) ---
	// TODO: Implement spell-specific targeting logic. Some spells might not need a target,
	// some might target self, some might have complex AoE rules.
	target, foundTarget := utils.FindNearestEnemy(s.world, caster, team.ID)
	if !foundTarget {
		log.Printf("SpellCastSystem (Start): Entity %d found no target for spell at %.3fs. Spell canceled.", caster, currentTime)
		return
	}

    currentMana := mana.GetCurrentMana() - mana.GetMaxMana()
	mana.SetCurrentMana(currentMana)

	// --- Enqueue SpellLandedEvent ---
	landTime := currentTime + spell.GetCastStartUp()
	spellLandedEvent := eventsys.SpellLandedEvent{
		Source:    caster,
		Target:    target, // Use the determined target
		SpellName: spell.GetName(),
		Timestamp: landTime,
	}
	s.eventBus.Enqueue(spellLandedEvent, landTime)
	
	spell.IncrementSpellCount()

	log.Printf("SpellCastSystem (Start): Entity %d started casting '%s' (-> %d) at %.3fs. Landing at %.3fs. Mana set to %.3f.", caster, spell.GetName(), target, currentTime, landTime, currentMana)
}

// handleSpellLanded applies spell effects (via DamageSystem usually) and starts recovery.
func (s *SpellCastSystem) handleSpellLanded(evt eventsys.SpellLandedEvent) {
	caster := evt.Source
	landTime := evt.Timestamp

	// --- Get Components ---
	_, okState := s.world.GetState(caster)
	spell, okSpell := s.world.GetSpell(caster)
	health, okHealth := s.world.GetHealth(caster) // Check if caster is still alive

	if !okState || !okSpell || !okHealth || health.GetCurrentHP() <= 0 {
		log.Printf("SpellCastSystem (Landed): Entity %d missing components or dead at %.3fs. Spell effect fizzles.", caster, landTime)
		// If the caster died during cast time, the spell doesn't happen.
		return
	}

	// Note: The actual spell *effect* (damage, healing, buffs) should typically be handled
	// by other systems listening for SpellLandedEvent (like DamageSystem).
	// This handler focuses on the state transition *after* the spell lands.

	log.Printf("SpellCastSystem (Landed): Entity %d landed spell '%s' at %.3fs.", caster, evt.SpellName, landTime)

	// --- Update State & Schedule Recovery End ---
	recoveryDuration := spell.GetCastRecovery()
	// no state update needed for spell, as casting cannot be interrupted by other actions.

	recoveryEndTime := landTime + recoveryDuration
	recoveryEndEvent := eventsys.SpellRecoveryEndEvent{
		Entity:    caster,
		Timestamp: recoveryEndTime,
	}
	s.eventBus.Enqueue(recoveryEndEvent, recoveryEndTime)
	log.Printf("SpellCastSystem (Landed): Entity %d starting spell recovery at %.3fs. Recovery ends at %.3fs.", caster, landTime, recoveryEndTime)
}

func (s *SpellCastSystem) handleSpellRecoveryEnd(evt eventsys.SpellRecoveryEndEvent) {
    caster := evt.Entity
    recoveryEndTime := evt.Timestamp

    // --- Get Components ---
    _, okState := s.world.GetState(caster)
    _, okSpell := s.world.GetSpell(caster)
    health, okHealth := s.world.GetHealth(caster) // Check if caster is still alive

    if !okState || !okSpell || !okHealth || health.GetCurrentHP() <= 0 {
        log.Printf("SpellCastSystem (RecoveryEnd): Entity %d missing components or dead at %.3fs. Recovery end ignored.", caster, recoveryEndTime)
        return
    }

    log.Printf("SpellCastSystem (RecoveryEnd): Entity %d finished spell recovery at %.3fs. Triggering action check.", caster, recoveryEndTime)

    // Enqueue ChampionActionEvent for the Action System to decide next step (AttackCooldown or AttackStart)
	actionCheckEvent := eventsys.ChampionActionEvent{
		Entity:    caster,
		Timestamp: recoveryEndTime, // Check action immediately
	}
	s.eventBus.Enqueue(actionCheckEvent, recoveryEndTime)
}
