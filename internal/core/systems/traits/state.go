package traitsys

// "tft-dps-simulator/internal/coreecs"

// TeamTraitState caches trait counts and active tiers per team.
type TeamTraitState struct {
    // unitCounts maps teamID -> trait ApiName -> count of units with that trait
    unitCounts map[int]map[string]int
    // activeTier maps teamID -> trait ApiName -> index of the highest active Effect in data.Trait.Effects
    // -1 indicates the trait is inactive for the team.
    activeTier map[int]map[string]int
}

// newTeamTraitState initializes the trait state cache.
func NewTeamTraitState() *TeamTraitState {
    return &TeamTraitState{
        unitCounts: make(map[int]map[string]int),
        activeTier: make(map[int]map[string]int),
    }
}

// EnsureTeam initializes maps for a given teamID if they don't exist.
func (tts *TeamTraitState) EnsureTeam(teamID int) {
    if _, ok := tts.unitCounts[teamID]; !ok {
        tts.unitCounts[teamID] = make(map[string]int)
        tts.activeTier[teamID] = make(map[string]int)
    }
}

// ResetTeam clears the state for a specific team.
func (tts *TeamTraitState) ResetTeam(teamID int) {
    delete(tts.unitCounts, teamID)
    delete(tts.activeTier, teamID)
    // Re-initialize maps to avoid nil pointer errors on subsequent access
    tts.EnsureTeam(teamID)
}

// ResetAll clears the state for all teams.
func (tts *TeamTraitState) ResetAll() {
    tts.unitCounts = make(map[int]map[string]int)
    tts.activeTier = make(map[int]map[string]int)
}

// GetUnitCount returns the count of a specific trait for a team.
func (tts *TeamTraitState) GetUnitCount(teamID int, traitApiName string) int {
	if counts, ok := tts.unitCounts[teamID]; ok {
		return counts[traitApiName]
	}
	return 0
}

// GetActiveTierForTrait returns the active tier index for a specific trait for a team.
func (tts *TeamTraitState) GetActiveTier(teamID int, traitApiName string) int {
	if tiers, ok := tts.activeTier[teamID]; ok {
		return tiers[traitApiName]
	}
	return -1 // Trait is inactive for the team
}

// GetActiveTiers returns the active tiers for all traits for a team.
func (tts *TeamTraitState) GetActiveTiers() map[int]map[string]int {
	return tts.activeTier
}