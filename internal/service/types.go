package service

import (
	"tft-dps-simulator/internal/core/components"
	"tft-dps-simulator/internal/core/ecs"
	eventsys "tft-dps-simulator/internal/core/systems/events"
)

// BoardPosition matches frontend/src/utils/types.ts
type BoardPosition struct {
	Row int `json:"row"`
	Col int `json:"col"`
}

// Item matches frontend/src/utils/types.ts (simplified for simulation input)
type Item struct {
	ApiName string `json:"apiName"` // Use ApiName as identifier
	// Name, Description, Image, Type might not be strictly needed for the simulation input
	// if the core simulation logic can look them up by ID from its loaded data.
	// Let's keep it simple for now and assume ID is sufficient.
}

// BoardChampion matches frontend/src/utils/types.ts (simplified for simulation input)
type BoardChampion struct {
	ApiName  string        `json:"apiName"` // Use ApiName as identifier
	Stars    int           `json:"stars"` // Default to 1 if missing? Frontend uses 1 | 2 | 3
	Items    []Item        `json:"items"`
	Position BoardPosition `json:"position"`
	// Name, Cost, Traits, Image might not be needed if ID is sufficient for lookup.
}

// RunSimulationRequest is the expected request body structure
type RunSimulationRequest struct {
	BoardChampions []BoardChampion `json:"boardChampions"`
	// We could add other context later if needed, like selected Augments
	// SelectedAugments []Augment `json:"selectedAugments"`
}

// ChampionSimulationResult holds the results for a single champion
type ChampionSimulationResult struct {
	ChampionApiName  string      `json:"championApiName"` // Match the ApiName sent in the request
	ChampionEntityID ecs.Entity `json:"championEntityId"` // Entity ID in ECS world
	DamageStats components.DamageStats `json:"damageStats"`
}

type ArchivedEvent struct {
	EventItem eventsys.EventItem `json:"eventItem"`
	EventType string `json:"eventType"` // Type of event (e.g., "damage", "heal", etc.)
}

// RunSimulationResponse is the structure of the response body
type RunSimulationResponse struct {
	Results []ChampionSimulationResult `json:"results"`
	ArchieveEvents []ArchivedEvent `json:"archieveEvents"`
}