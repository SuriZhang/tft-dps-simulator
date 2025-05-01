package traitsys

import (
	"log"
	"reflect"
	"sort"

	"tft-dps-simulator/internal/core/components"
	"tft-dps-simulator/internal/core/data"
	"tft-dps-simulator/internal/core/ecs"
)

// TraitCounterSystem calculates active trait tiers based on unit counts.
type TraitCounterSystem struct {
    world      *ecs.World
    traitState *TeamTraitState
}

// NewTraitCounterSystem creates a new TraitCounterSystem.
func NewTraitCounterSystem(world *ecs.World, state *TeamTraitState) *TraitCounterSystem {
    return &TraitCounterSystem{
        world:      world,
        traitState: state,
    }
}

// UpdateCountsAndTiers recalculates trait counts and determines active tiers for all teams.
// This should be called on combat start or whenever team compositions change significantly.
func (s *TraitCounterSystem) UpdateCountsAndTiers() {
    log.Println("TraitCounterSystem: Updating counts and tiers...")
    s.traitState.ResetAll() // Clear previous counts and tiers

    teamType := reflect.TypeOf(components.Team{})
    traitsType := reflect.TypeOf(components.Traits{})
    ChampionInfoType := reflect.TypeOf(components.ChampionInfo{})
    entities := s.world.GetEntitiesWithComponents(teamType, traitsType, ChampionInfoType) 

    // 1. Identify unique champions per team and their traits
    uniqueChampsPerTeam := make(map[int]map[string][]string) // teamID -> championApiName -> traitsList

    for _, entity := range entities {
        team, okTeam := s.world.GetTeam(entity)
        traits, okTraits := s.world.GetTraits(entity)
        ChampionInfo, okInfo := s.world.GetChampionInfo(entity)

        if !okTeam || !okTraits || !okInfo {
            log.Printf("Warning: Entity %d missing required components for trait counting", entity)
            continue
        }

        if _, teamExists := uniqueChampsPerTeam[team.ID]; !teamExists {
            uniqueChampsPerTeam[team.ID] = make(map[string][]string)
        }

        // Store the traits list using the champion's name as the key.
        // This automatically handles duplicates - only one entry per champion name per team.
        uniqueChampsPerTeam[team.ID][ChampionInfo.ApiName] = traits.GetTraits()
    }

    // 2. Count traits based on unique champions
    for teamID, championsMap := range uniqueChampsPerTeam {
        s.traitState.EnsureTeam(teamID) // Make sure team maps are initialized in traitState

        log.Printf("TraitCounterSystem: Counting unique champion traits for Team %d", teamID)
        for champApiName, traitsList := range championsMap {
            for _, traitName := range traitsList {
                // Ensure the trait exists in the loaded data
                if _, exists := data.Traits[traitName]; exists {
                    s.traitState.unitCounts[teamID][traitName]++
                } else {
                    // Log warning only once per unique champion type if needed, but logging here is fine too.
                    log.Printf("Warning: Champion type '%s' has unknown trait '%s'", champApiName, traitName)
                }
            }
        }
		// Print active trait counts for debugging
		log.Printf("TraitCounterSystem: Team %d trait counts: %v", teamID, s.traitState.unitCounts[teamID])
    }


    // 3. Determine active tier for each trait per team (This part remains the same)
    for teamID, traitCounts := range s.traitState.unitCounts {
        log.Printf("TraitCounterSystem: Calculating active tiers for Team %d", teamID)
        for traitName, count := range traitCounts {
            traitData := data.Traits[traitName]

            activeTierIndex := -1 // Default to inactive
            // Sort effects by MinUnits to ensure correct tier detection
            // Optimization: This sort could happen once globally if traitData.Effects is never modified.
            sort.SliceStable(traitData.Effects, func(i, j int) bool {
                return traitData.Effects[i].MinUnits < traitData.Effects[j].MinUnits
            })

            for i, effect := range traitData.Effects {
                if count >= effect.MinUnits {
                    activeTierIndex = i // Found the highest active tier
                } else {
                    break // Since effects are sorted, no higher tier can be active
                }
            }

            // Store the active tier index
            s.traitState.activeTier[teamID][traitName] = activeTierIndex

            if activeTierIndex != -1 {
                log.Printf("  Team %d: Trait '%s' active at tier %d (Count: %d, MinUnits: %d, MaxUnits: %d, Style: %d)",
                    teamID, traitName, activeTierIndex, count, traitData.Effects[activeTierIndex].MinUnits, traitData.Effects[activeTierIndex].MaxUnits, traitData.Effects[activeTierIndex].Style)
            } 
        }
    }
	log.Printf("TraitCounterSystem: Final active tiers: %v", s.traitState.activeTier)
	log.Printf("TraitCounterSystem: Final unit counts: %v", s.traitState.unitCounts)
    log.Println("TraitCounterSystem: Finished updating counts and tiers.")
}