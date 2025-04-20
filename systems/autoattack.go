package systems

import (
	"log"
	"reflect"

	"github.com/suriz/tft-dps-simulator/components"
	"github.com/suriz/tft-dps-simulator/ecs"
	eventsys "github.com/suriz/tft-dps-simulator/systems/events"
	"github.com/suriz/tft-dps-simulator/utils"
)

// AutoAttackSystem handles champion auto attacks based on Team ID.
type AutoAttackSystem struct {
	world       *ecs.World
	eventBus    eventsys.EventBus // Add event bus
	currentTime float64
}

// NewAutoAttackSystem creates a new auto attack system.
func NewAutoAttackSystem(world *ecs.World, bus eventsys.EventBus) *AutoAttackSystem { // Accept event bus
	return &AutoAttackSystem{
		world:       world,
		eventBus:    bus, // Store event bus
		currentTime: 0.0,
	}
}

// TriggerAutoAttack processes auto attacks for the current timestep.
func (s *AutoAttackSystem) TriggerAutoAttack(deltaTime float64) {
	s.currentTime += deltaTime

	attackType := reflect.TypeOf(components.Attack{})
	teamType := reflect.TypeOf(components.Team{})
	posType := reflect.TypeOf(components.Position{})
	healthType := reflect.TypeOf(components.Health{})

	potentialAttackers := s.world.GetEntitiesWithComponents(attackType, teamType, posType, healthType)

	for _, attacker := range potentialAttackers {
		attack, okAtk := s.world.GetAttack(attacker)
		team, okTeam := s.world.GetTeam(attacker)
		pos, okPos := s.world.GetPosition(attacker)
		health, okHealth := s.world.GetHealth(attacker)

		if !okAtk || !okTeam || !okPos || !okHealth || health.CurrentHP <= 0 {
			continue // Skip if missing components or dead
		}

		finalAS := attack.GetFinalAttackSpeed()
		if finalAS <= 0 {
			// log.Printf("Attacker %d has 0 AS, skipping.\n", attacker) // Reduce log spam
			continue
		}

		// --- Check Post Spell Cast Recovery ---
		if spell, okSpell := s.world.GetSpell(attacker); okSpell {
			if spell.GetCurrentRecovery() > 0 {
				attack.SetAttackStartupEndTime(-1.0) // Reset attack schedule
				log.Printf("Attacker %d is locked out by spell cast recovery (%.2f).\n", attacker, spell.GetCurrentRecovery())
				continue
			}
		}


		// --- Simplified Attack Scheduling Logic ---
		// Check for scheduling the VERY FIRST attack
        // Use LastAttackTime == 0.0 and AttackStartupEndTime == -1.0 as the trigger
        if attack.GetLastAttackTime() == 0.0 && attack.GetAttackStartupEndTime() == -1.0 {
            // This block runs only if *not* locked out this tick
            finalAS := attack.GetFinalAttackSpeed()
            if finalAS <= 0 {
                log.Printf("Attacker %d has 0 AS when trying to schedule first attack.", attacker)
                continue
            }
            totalCycleTime := 1.0 / finalAS

            // --- Determine Effective Start Time ---
            var effectiveStartTime float64
            // If current time is just the first delta, assume true start at t=0
            // Otherwise, assume a state (like lockout) just ended, so the cycle *should* have started at the beginning of this tick.
            if s.currentTime <= deltaTime { // Check if it's the first time TriggerAutoAttack runs *without lockout*
                effectiveStartTime = 0.0
            } else {
                // Assume lockout/other state ended just before this tick started
                effectiveStartTime = s.currentTime - deltaTime
            }

            firstLandingTime := effectiveStartTime + totalCycleTime

            attack.SetAttackStartupEndTime(firstLandingTime) // Schedule the first landing
            attack.SetAttackCycleEndTime(firstLandingTime) // Next cycle can start checking after this landing

            log.Printf("AutoAttackSystem: Attacker %d scheduling FIRST attack to land at %.3f (effective start: %.3f, current time: %.3f).", attacker, firstLandingTime, effectiveStartTime, s.currentTime)
        }

        // Check for LANDING an attack and scheduling the NEXT one
        scheduledLandingTime := attack.GetAttackStartupEndTime()

		if scheduledLandingTime != -1.0 && s.currentTime >= scheduledLandingTime {
			log.Printf("DEBUG: currentTime=%.3f, scheduledLandingTime=%.3f\n", s.currentTime, scheduledLandingTime)
			target, foundTarget := utils.FindNearestEnemy(s.world, attacker, team.ID)

			// *** RESERVED INTERRUPTION CHECK FOR CC DURING STARTUP***
            // if s.world.IsAttackPrevented(attacker) { // Assumes this function exists in your world/ECS
            //     log.Printf("AutoAttackSystem: Attacker %d attack scheduled for %.3f CANCELED due to interruption (e.g., CC) at time %.3f.", attacker, scheduledLandingTime, s.currentTime)

            //     // Reset the attack state to cancel the current attack and allow rescheduling after CC
            //     attack.SetAttackStartupEndTime(-1.0) // Mark scheduled attack as processed/canceled
            //     attack.SetAttackCycleEndTime(s.currentTime) // Allow starting a new cycle check immediately after CC wears off
            //     // LastAttackTime remains unchanged, as the attack didn't land.

            //     continue // Skip the rest of the landing logic for this attacker this tick
            // }
            // // *** END INTERRUPTION CHECK ***

			if foundTarget {
				targetPos, okTargetPos := s.world.GetPosition(target)
				if okTargetPos {
					dx := targetPos.GetX() - pos.GetX()
					dy := targetPos.GetY() - pos.GetY()
					distSq := dx*dx + dy*dy
					attackRange := attack.GetFinalRange()
					rangeSq := attackRange * attackRange

					if distSq <= rangeSq {
						// Target in range, enqueue event
						log.Printf("AutoAttackSystem: Attacker %d LANDED attack on %d at %.3f.", attacker, target, scheduledLandingTime)
						attackEvent := eventsys.AttackLandedEvent{
							Source:     attacker,
							Target:     target,
							BaseDamage: attack.GetFinalAD(), // Use current AD at time of landing
							Timestamp:  scheduledLandingTime,
						}
						s.eventBus.Enqueue(attackEvent)

						// Schedule the NEXT attack
						totalCycleTime := 1.0 / finalAS // Recalculate in case AS changed
						nextLandingTime := scheduledLandingTime + totalCycleTime

						attack.SetAttackStartupEndTime(nextLandingTime) // Schedule the *next* landing
						attack.SetLastAttackTime(scheduledLandingTime) // Record when the *last* attack landed
						attack.SetAttackCycleEndTime(nextLandingTime) // Update cycle end time

						log.Printf("AutoAttackSystem: Attacker %d scheduled next attack to land at %.3f.", attacker, nextLandingTime)

					} else {
						// Target found but out of range at landing time
						log.Printf("AutoAttackSystem: Attacker %d attack scheduled for %.3f missed target %d (out of range).", attacker, scheduledLandingTime, target)
						// Don't schedule next attack yet, wait for target to come in range?
						// For now, let's reset to allow rescheduling on next opportunity (might need refinement)
						attack.SetAttackStartupEndTime(-1.0) // Reset schedule
						attack.SetLastAttackTime(s.currentTime) // Update last attempt time
						attack.SetAttackCycleEndTime(s.currentTime) // Allow immediate re-check
					}
				} else {
					// Target has no position? Should not happen if FindNearestEnemy worked.
					log.Printf("AutoAttackSystem: Attacker %d attack scheduled for %.3f missed target %d (no position).", attacker, scheduledLandingTime, target)
					attack.SetAttackStartupEndTime(-1.0)
					attack.SetLastAttackTime(s.currentTime)
					attack.SetAttackCycleEndTime(s.currentTime)
				}
			} else {
				// No target found at landing time
				log.Printf("AutoAttackSystem: Attacker %d attack scheduled for %.3f missed (no target found).", attacker, scheduledLandingTime)
				attack.SetAttackStartupEndTime(-1.0)
				attack.SetLastAttackTime(s.currentTime)
				attack.SetAttackCycleEndTime(s.currentTime)
			}
		} // End check for landing attack
	}
}

// SetCurrentTime sets the current time for the system.
func (s *AutoAttackSystem) SetCurrentTime(time float64) {
	s.currentTime = time
}
