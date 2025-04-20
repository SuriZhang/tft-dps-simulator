package systems

import (
	"log"
	"math"
	"reflect"

	"github.com/suriz/tft-dps-simulator/components"
	"github.com/suriz/tft-dps-simulator/ecs"
	eventsys "github.com/suriz/tft-dps-simulator/systems/events"
	"github.com/suriz/tft-dps-simulator/utils"
)

type SpellCastSystem struct {
    world       *ecs.World
    eventBus    eventsys.EventBus
    currentTime float64 // Track simulation time - needs to be updated externally or passed in
}

func NewSpellCastSystem(world *ecs.World, bus eventsys.EventBus) *SpellCastSystem {
    return &SpellCastSystem{world: world, eventBus: bus, currentTime: 0.0}
}

// SetCurrentTime allows the simulation loop to update the system's internal time
func (s *SpellCastSystem) SetCurrentTime(time float64) {
    s.currentTime = time
}

// Update checks for spell casts and manages cooldowns.
func (s *SpellCastSystem) TriggerSpellCast(deltaTime float64) {
    // --- 1. Cooldown Reduction ---
    spellType := reflect.TypeOf(components.Spell{})
    entitiesWithSpells := s.world.GetEntitiesWithComponents(spellType)
    for _, entity := range entitiesWithSpells {
        if spell, ok := s.world.GetSpell(entity); ok {
            if spell.GetCurrentRecovery() > 0 {
                newCooldown := math.Max(0, spell.GetCurrentRecovery()-deltaTime)
                spell.SetCurrentRecovery(newCooldown)
            }
        }
    }

    // --- 2. Casting Logic ---
    manaType := reflect.TypeOf(components.Mana{})
    teamType := reflect.TypeOf(components.Team{})
    spellType = reflect.TypeOf(components.Spell{}) // Ensure spellType is defined here too

    potentialCasters := s.world.GetEntitiesWithComponents(spellType, manaType, teamType)

    for _, caster := range potentialCasters {
        spell, okSpell := s.world.GetSpell(caster)
        mana, okMana := s.world.GetMana(caster)
        team, okTeam := s.world.GetTeam(caster)

        // Ensure components exist and entity is on the player team
        if !okSpell || !okMana || !okTeam || team.ID != 0 {
            continue
        }

        // Check conditions: Off cooldown and enough mana
        if spell.GetCurrentRecovery() <= 0 && mana.GetCurrentMana() >= spell.GetManaCost() {

            // --- Find Target using Utility Function ---
            target, foundTarget := utils.FindNearestEnemy(s.world, caster, team.ID)
            // --- End Find Target ---

            if foundTarget { // Check the boolean return value
                // Execute Cast
                log.Printf("SpellCastSystem: Entity %d casting spell '%s' on nearest target %d at time %.2f", caster, spell.GetName(), target, s.currentTime)

            } else { 
                log.Printf("SpellCastSystem: Entity %d has mana for spell '%s' but no valid target found. Still enqueue an event", caster, spell.GetName())
				target = 0 // Set target to 0 if no valid target is found
            }
            mana.SetCurrentMana(mana.GetCurrentMana() - spell.GetManaCost())
            spell.SetCurrentRecovery(spell.GetCastRecovery()) // Reset cooldown

            // Enqueue Event
            castEvent := eventsys.SpellCastEvent{
                Source:    caster,
                Target:    target,
                Timestamp: s.currentTime,
            }
            s.eventBus.Enqueue(castEvent)
        }
    }
}

