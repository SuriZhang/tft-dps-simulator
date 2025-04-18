package itemsys

import (
	"log"
	"reflect"

	"github.com/suriz/tft-dps-simulator/components"
	effects "github.com/suriz/tft-dps-simulator/components/effects"
	"github.com/suriz/tft-dps-simulator/data"
	"github.com/suriz/tft-dps-simulator/ecs"
)

// DynamicTimeItemSystem handles items whose effects change over time.
type DynamicTimeItemSystem struct {
	world *ecs.World
}

// NewDynamicTimeItemSystem creates a new system instance.
func NewDynamicTimeItemSystem(world *ecs.World) *DynamicTimeItemSystem {
	return &DynamicTimeItemSystem{
		world: world,
	}
}

// Update processes time-based item effects for relevant entities.
func (s *DynamicTimeItemSystem) Update(deltaTime float64) {
	// Query entities that have equipment (potential item carriers)
	// We don't query for the specific effect components because they might be added/removed by this system.
	entities := s.world.GetEntitiesWithComponents(reflect.TypeOf(components.Equipment{}))

	for _, entity := range entities {
		equipment, ok := s.world.GetEquipment(entity)
		if !ok {
			continue // Should not happen based on query, but good practice
		}

		// --- Archangel's Staff ---
		if equipment.HasItem(data.TFT_Item_ArchangelsStaff) {
			s.updateArchangels(entity, deltaTime)
		} else {
			// If the item was removed, ensure the effect component is also removed
			if _, exists := s.world.GetArchangelsEffect(entity); exists {
				s.world.RemoveComponent(entity, reflect.TypeOf(effects.ArchangelsEffect{}))
			}
		}

		// --- Quicksilver ---
		if equipment.HasItem(data.TFT_Item_Quicksilver) {
			s.updateQuicksilver(entity, deltaTime)
		} else {
			// If the item was removed, ensure the effect component is also removed
			if _, exists := s.world.GetQuicksilverEffect(entity); exists {
				s.world.RemoveComponent(entity, reflect.TypeOf(effects.QuicksilverEffect{}))
				// TODO: Also remove IsImmuneToCC marker if implemented
			}
		}

		// --- Add logic for other DynamicTime items here ---
		// Example: Redemption (needs careful handling of activation trigger/delay)
		// if equipment.HasItem("TFT_Item_Redemption") { ... }
	}
}

// updateArchangels handles the logic for Archangel's Staff.
func (s *DynamicTimeItemSystem) updateArchangels(entity ecs.Entity, deltaTime float64) {
	// Check how many Archangel's Staffs are equipped first
	equipComp, equipOk := s.world.GetEquipment(entity)
	if !equipOk {
		return // No equipment component, shouldn't happen if item was added
	}
	archangelsCount := equipComp.GetItemCount(data.TFT_Item_ArchangelsStaff)
	if archangelsCount == 0 {
		// No Archangel's Staff equipped, nothing to do
		return
	}

	effect, exists := s.world.GetArchangelsEffect(entity)
	if !exists {
		// Component should have been added by equipment manager when item equipped.
		return
	}

	// --- Stacking Logic (uses the single effect component) ---
	effect.AddTimer(deltaTime)
	procOccurredThisFrame := false
	if effect.GetTimer() >= effect.GetInterval() {
		// Calculate how many intervals passed (handles large deltaTime)
		intervalsPassed := int(effect.GetTimer() / effect.GetInterval())
		effect.AddStacks(intervalsPassed) // Assuming AddStacks(n int) exists
		// Subtract only the time for the intervals that passed
		effect.SetTimer(effect.GetTimer() - float64(intervalsPassed)*effect.GetInterval()) // Assuming SetTimer exists
		procOccurredThisFrame = true
		log.Printf("Entity %d Archangel's Proc! Stacks: %d", entity, effect.GetStacks())
	}

	// --- Apply Bonus AP (multiplied by item count) ---
	spellComp, spellOk := s.world.GetSpell(entity)
	if spellOk {
		// Calculate bonus AP from *one* staff based on current stacks
		singleBonus := float64(effect.GetStacks()) * effect.GetAPPerInterval()
		// Multiply by the number of Archangel's Staffs equipped
		totalBonus := singleBonus * float64(archangelsCount)

		// Apply the total bonus (AddBonusAP should sum bonuses from different sources)
		if totalBonus > 0 {
			spellComp.SetBonusAP(totalBonus) // Assuming AddBonusAP exists
			if procOccurredThisFrame {
				log.Printf("Entity %d Archangel's Stacks: %d, Count: %d, Applied Total Bonus AP: %.1f, Total Bonus AP now: %.1f\n", entity, effect.GetStacks(), archangelsCount, totalBonus, spellComp.GetBonusAP())
			}
		}
	}
}

// updateQuicksilver handles the logic for Quicksilver.
func (s *DynamicTimeItemSystem) updateQuicksilver(entity ecs.Entity, deltaTime float64) {
	effect, exists := s.world.GetQuicksilverEffect(entity)
	if !exists {
		return // Component removed previously or never added
	}

	// Only process if the effect is marked as active.
	if effect.IsActive() {

		// --- Determine if a proc interval boundary is crossed THIS frame ---
		oldProcTimer := effect.GetProcTimer()
		newProcTimer := oldProcTimer + deltaTime
		procInterval := effect.GetProcInterval()
		procOccurredThisFrame := false
		newStacksToAdd := 0

		// Check how many full intervals were completed by the end of this frame
		// compared to the start. Use integer division for floor effect.
		intervalsAtStart := int(oldProcTimer / procInterval)
		intervalsAtEnd := int(newProcTimer / procInterval)

		if intervalsAtEnd > intervalsAtStart {
			newStacksToAdd = intervalsAtEnd - intervalsAtStart
			procOccurredThisFrame = true
			// Update the timer, carrying over the remainder past the last completed interval
			remainder := newProcTimer - float64(intervalsAtEnd)*procInterval
			effect.SetProcTimer(remainder) // Assuming SetProcTimer exists
		} else {
			// No full interval completed, just update the timer
			effect.SetProcTimer(newProcTimer)
		}

		// --- If proc occurred, update internal state ---
		if procOccurredThisFrame {
			// Assuming AddStacks updates the internal stack count
			effect.AddStacks(newStacksToAdd)
			effect.AddBonusAS(float64(newStacksToAdd) * effect.GetProcAttackSpeed())
			// Assuming GetCurrentBonusAS now calculates based on the new stack count
			log.Printf("Entity %d Quicksilver Proc! Stacks: %d, Internal Bonus AS now: %.2f%%", entity, effect.GetStacks(), effect.GetCurrentBonusAS()*100)
		}

		// --- Apply the *delta* bonus AS gained this frame ---
		// Only apply if a proc occurred this frame.
		if procOccurredThisFrame && newStacksToAdd > 0 {
			attackComp, attackOk := s.world.GetAttack(entity)
			if attackOk {
				// Calculate the bonus AS gained *this frame*
				deltaBonusAS := float64(newStacksToAdd) * effect.GetProcAttackSpeed()

				if deltaBonusAS > 0 {
					before := attackComp.GetBonusPercentAttackSpeed()
					log.Printf("Entity %d Quicksilver Before Delta AS: %.2f%%", entity, before*100)
					// Add the newly gained bonus AS to the attack component
					attackComp.AddBonusPercentAttackSpeed(deltaBonusAS)
					log.Printf("Entity %d Quicksilver Applied Delta AS: +%.2f%%, Total Bonus AS now: %.2f%%", entity, deltaBonusAS*100, attackComp.GetBonusPercentAttackSpeed()*100)
				}
			}
		}

		// --- Handle Duration Countdown ---
		effect.DecreaseRemainingDuration(deltaTime)
		// log.Printf("Entity %d Quicksilver remaining: %.2f\n", entity, effect.GetRemainingDuration()) // Log remaining duration

		// --- Check for Expiry AFTER applying effects for this frame ---
		// Use <= 0 to catch expiry exactly at 0.0
		if effect.GetRemainingDuration() <= 0 {
			// Duration expired THIS frame. Mark inactive.
			// Do NOT remove the component here.
			log.Printf("Entity %d Quicksilver expired (Duration <= 0), setting inactive.\n", entity)
			effect.SetIsActive(false)
			// TODO: Remove IsImmuneToCC marker if implemented
		}
	}
	// If effect.IsActive() is false, do nothing.
}
