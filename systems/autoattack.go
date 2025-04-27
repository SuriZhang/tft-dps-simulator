package systems

import (
	"log"
	"math" // Needed for range check

	"github.com/suriz/tft-dps-simulator/ecs"
	eventsys "github.com/suriz/tft-dps-simulator/systems/events"
	"github.com/suriz/tft-dps-simulator/utils"
)

// AutoAttackSystem handles the auto-attack cycle based on events.
type AutoAttackSystem struct {
	world    *ecs.World
	eventBus eventsys.EventBus
}

// NewAutoAttackSystem creates a new auto attack system.
func NewAutoAttackSystem(world *ecs.World, bus eventsys.EventBus) *AutoAttackSystem {
	return &AutoAttackSystem{
		world:    world,
		eventBus: bus,
	}
}

// CanHandle checks if the system can process the given event type.
func (s *AutoAttackSystem) CanHandle(evt interface{}) bool {
	switch evt.(type) {
	case eventsys.AttackStartupEvent,
		eventsys.AttackFiredEvent,
		eventsys.AttackRecoveryEndEvent,
		eventsys.AttackCooldownStartEvent,
		eventsys.AttackLandedEvent,
		eventsys.AttackCooldownEndEvent:
		return true
	default:
		return false
	}
}

// HandleEvent processes events related to the auto-attack cycle.
func (s *AutoAttackSystem) HandleEvent(evt interface{}) {
	switch event := evt.(type) {
	case eventsys.AttackStartupEvent:
		s.handleAttackStart(event)
	case eventsys.AttackFiredEvent:
		s.handleAttackFired(event)
	case eventsys.AttackLandedEvent:
		s.handleAttackLanded(event)
	case eventsys.AttackRecoveryEndEvent:
		s.handleAttackRecoveryEnd(event)
	case eventsys.AttackCooldownStartEvent:
		s.handleAttackCooldownStart(event)
	case eventsys.AttackCooldownEndEvent:
		s.handleAttackCooldownEnd(event)
	}
}

// handleAttackStart initiates the attack wind-up.
func (s *AutoAttackSystem) handleAttackStart(evt eventsys.AttackStartupEvent) {
	attacker := evt.Entity
	currentTime := evt.Timestamp

	// --- Get Components ---
	state, okState := s.world.GetState(attacker)
	attack, okAttack := s.world.GetAttack(attacker)
	_, okPos := s.world.GetPosition(attacker)
	team, okTeam := s.world.GetTeam(attacker)
	health, okHealth := s.world.GetHealth(attacker)

	if !okState || !okAttack || !okPos || !okTeam || !okHealth || health.GetCurrentHP() <= 0 {
		log.Printf("AutoAttackSystem (Start): Entity %d missing components or dead at %.3fs. Canceling attack start.", attacker, currentTime)
		return // Cannot attack without necessary components or if dead
	}

	// --- Targeting ---
	target, foundTarget := utils.FindNearestEnemy(s.world, attacker, team.ID)
	if !foundTarget {
		log.Printf("AutoAttackSystem (Start): Entity %d found no target at %.3fs. Attack canceled.", attacker, currentTime)
		// TODO: What should happen if no target? Enqueue next ChampionActionEvent after a delay?
		// For now, just cancel and wait for the next trigger.
		// state.EndAction() // Go back to idle
		return
	}

	// --- Update State ---
	baseStartup := attack.GetBaseAttackStartup()
	baseRecovery := attack.GetBaseAttackRecovery()
	baseAS := attack.GetBaseAttackSpeed()
	finalAS := attack.GetFinalAttackSpeed()

	if baseAS == 0 || finalAS == 0 {
		log.Printf("WARN: AutoAttackSystem (Start): Entity %d has non-positive attack speed (Base: %.2f, Final: %.2f) at %.3fs. Cannot start attack.", attacker, baseAS, finalAS, currentTime)
		return // Cannot attack with invalid speeds
	}

	var startupDuration, recoveryDuration, cooldownDuration float64
	// Adjust startup time based on attack speed ratio
	// Faster attack speed (finalAS > baseAS) reduces startup time.
	if finalAS > 0 && baseAS > 0 {
		startupDuration = baseStartup * (baseAS / finalAS)
		recoveryDuration = baseRecovery * (baseAS / finalAS)
		cooldownDuration = 1/finalAS - startupDuration - recoveryDuration
	} else {
		// Fallback or error handling if attack speeds are invalid
		log.Printf("WARN: AutoAttackSystem (Start): Entity %d has invalid attack speed (Base: %.2f, Final: %.2f). Using base startup time.", attacker, baseAS, finalAS)
		startupDuration = baseStartup
	}
	attack.SetCurrentAttackStartup(startupDuration)
	attack.SetCurrentAttackRecovery(recoveryDuration)
	attack.SetCurrentAttackCooldown(cooldownDuration)
	state.StartAttack(currentTime, startupDuration)

	// --- Enqueue AttackFiredEvent ---
	fireTime := currentTime + startupDuration
	attackFiredEvent := eventsys.AttackFiredEvent{
		Source:    attacker,
		Target:    target, // Target locked at start of attack
		Timestamp: fireTime,
	}
	s.eventBus.Enqueue(attackFiredEvent, fireTime)

	log.Printf("AutoAttackSystem (Start): Entity %d started attack (-> %d) at %.3fs. Firing at %.3fs.", attacker, target, currentTime, fireTime)
}

// handleAttackFired marks the point the attack connects/projectile launches.
func (s *AutoAttackSystem) handleAttackFired(evt eventsys.AttackFiredEvent) {
	attacker := evt.Source
	target := evt.Target // Use the target from the event
	fireTime := evt.Timestamp

	// --- Get Components ---
	state, okState := s.world.GetState(attacker)
	attack, okAttack := s.world.GetAttack(attacker)
	health, okHealth := s.world.GetHealth(attacker) // Check if attacker is still alive

	if !okState || !okAttack || !okHealth || health.GetCurrentHP() <= 0 {
		log.Printf("AutoAttackSystem (Fired): Entity %d missing components or dead at %.3fs. Attack fizzles.", attacker, fireTime)
		// If the attacker died during startup, the attack doesn't happen.
		// State might need cleanup depending on how death is handled.
		return
	}

	// --- Check Target Validity (Still alive? Still in range?) ---
	targetHealth, okTargetHealth := s.world.GetHealth(target)
	targetPos, okTargetPos := s.world.GetPosition(target)
	attackerPos, okAttackerPos := s.world.GetPosition(attacker)

	if !okTargetHealth || targetHealth.GetCurrentHP() <= 0 {
		log.Printf("AutoAttackSystem (Fired): Target %d for attack by %d is dead at %.3fs. Attack fizzles.", target, attacker, fireTime)
		// TODO: Handle attack fizzling - need to schedule recovery/cooldown anyway?
		// For now, let's assume the cycle continues to recovery/cooldown.
	} else if okTargetPos && okAttackerPos {
		// Check range at the moment of firing (TFT rule?) or landing? Assuming firing.
		distSq := utils.DistSq(attackerPos.GetX(), attackerPos.GetY(), targetPos.GetX(), targetPos.GetY())
		rangeSq := math.Pow(float64(attack.GetFinalRange()), 2)

		if distSq <= rangeSq {
			landedEvent := eventsys.AttackLandedEvent{
				Source:     attacker,
				Target:     target,
				BaseDamage: attack.GetFinalAD(), // AD at time of firing/landing? Using current AD.
				Timestamp:  fireTime,            // Assuming instant hit for now
			}
			s.eventBus.Enqueue(landedEvent, fireTime)
			// Target valid and in range, enqueue AttackLandedEvent
			log.Printf("AutoAttackSystem (Fired): Entity %d fired attack at %d at %.3fs. Enqueued AttackLandedEvent.", attacker, target, fireTime)
		} else {
			log.Printf("AutoAttackSystem (Fired): Target %d moved out of range of %d at %.3fs. Attack fizzles.", target, attacker, fireTime)
			// Attack fizzles, but recovery/cooldown still happens.

		}
	} else {
		log.Printf("AutoAttackSystem (Fired): Target %d or Attacker %d missing position at %.3fs. Attack fizzles.", target, attacker, fireTime)
		// Attack fizzles, recovery/cooldown still happens.
	}
	attack.IncrementAttackCount()
	state.StartAttackRecovery(fireTime, attack.GetCurrentAttackRecovery())
}

func (s *AutoAttackSystem) handleAttackLanded(evt eventsys.AttackLandedEvent) {
	attacker := evt.Source
	landedTime := evt.Timestamp

	// --- Get Components ---
	_, okState := s.world.GetState(attacker)
	attack, okAttack := s.world.GetAttack(attacker)
	health, okHealth := s.world.GetHealth(attacker) // Check if attacker is still alive

	if !okState || !okAttack || !okHealth || health.GetCurrentHP() <= 0 {
		log.Printf("AutoAttackSystem (Landed): Entity %d missing components or dead at %.3fs. Attack fizzles.", attacker, landedTime)
		// If the attacker died during startup, the attack doesn't happen.
		// State might need cleanup depending on how death is handled.
		return
	}

	recoveryEndTime := landedTime + attack.GetCurrentAttackRecovery()

	recoveryEndEvent := eventsys.AttackRecoveryEndEvent{
		Entity:    attacker,
		Timestamp: recoveryEndTime,
	}
	s.eventBus.Enqueue(recoveryEndEvent, recoveryEndTime)
	log.Printf("AutoAttackSystem (Landed): Entity %d starting recovery at %.3fs. Recovery ends at %.3fs.", attacker, landedTime, recoveryEndTime)
}

// handleAttackRecoveryEnd updates state and triggers ChampionActionSystem check.
func (s *AutoAttackSystem) handleAttackRecoveryEnd(evt eventsys.AttackRecoveryEndEvent) {
	entity := evt.Entity
	recoveryEndTime := evt.Timestamp

	_, okState := s.world.GetState(entity)
	health, okHealth := s.world.GetHealth(entity)

	if !okState || !okHealth || health.GetCurrentHP() <= 0 {
		log.Printf("AutoAttackSystem (RecoveryEnd): Entity %d missing components or dead at %.3fs.", entity, recoveryEndTime)
		return
	}

	log.Printf("AutoAttackSystem (RecoveryEnd): Entity %d finished recovery at %.3fs. Triggering action check.", entity, recoveryEndTime)
	// Enqueue ChampionActionEvent for the Action System to decide next step (Cast or CooldownStart)
	actionCheckEvent := eventsys.ChampionActionEvent{
		Entity:    entity,
		Timestamp: recoveryEndTime, // Check action immediately
	}
	s.eventBus.Enqueue(actionCheckEvent, recoveryEndTime)
}

// handleAttackCooldownStart calculates cooldown duration, updates state, and schedules cooldown end.
func (s *AutoAttackSystem) handleAttackCooldownStart(evt eventsys.AttackCooldownStartEvent) {
	entity := evt.Entity
	cooldownStartTime := evt.Timestamp

	state, okState := s.world.GetState(entity)
	attack, okAttack := s.world.GetAttack(entity)
	health, okHealth := s.world.GetHealth(entity)

	if !okState || !okAttack || !okHealth || health.GetCurrentHP() <= 0 {
		log.Printf("AutoAttackSystem (CooldownStart): Entity %d missing components or dead at %.3fs. Skipping cooldown.", entity, cooldownStartTime)
		return
	}

	// use ActonDuration because AttackCooldown might not be a full CD due to spell casted, calculation for total cooldown time is done in ActionSystem.
	cooldownEndTime := cooldownStartTime + state.ActionDuration
	cooldownEndEvent := eventsys.AttackCooldownEndEvent{
		Entity:    entity,
		Timestamp: cooldownEndTime,
	}
	s.eventBus.Enqueue(cooldownEndEvent, cooldownEndTime)
	log.Printf("AutoAttackSystem (CooldownStart): Entity %d starting cooldown at %.3fs (duration %.3fs). Cooldown ends at %.3fs.", entity, cooldownStartTime, attack.GetCurrentAttackCooldown(), cooldownEndTime)
}

// handleAttackCooldownEnd transitions state to idle and triggers ChampionActionSystem check.
func (s *AutoAttackSystem) handleAttackCooldownEnd(evt eventsys.AttackCooldownEndEvent) {
	entity := evt.Entity
	cooldownEndTime := evt.Timestamp

	_, okState := s.world.GetState(entity)
	health, okHealth := s.world.GetHealth(entity)

	if !okState || !okHealth || health.GetCurrentHP() <= 0 {
		log.Printf("AutoAttackSystem (CooldownEnd): Entity %d missing state or dead at %.3fs.", entity, cooldownEndTime)
		return
	}

	log.Printf("AutoAttackSystem (CooldownEnd): Entity %d finished cooldown at %.3fs. Triggering action check.", entity, cooldownEndTime)

	// Enqueue ChampionActionEvent for the Action System to decide next step (Cast or AttackStart)
	actionCheckEvent := eventsys.ChampionActionEvent{
		Entity:    entity,
		Timestamp: cooldownEndTime, // Check action immediately
	}
	s.eventBus.Enqueue(actionCheckEvent, cooldownEndTime)
}
