package data

// Global variable to store champions for quick lookup
var Champions map[string]*Champion

// GetChampionByApiName returns a champion by name or nil if not found
func GetChampionByApiName(apiName string) *Champion {
	// Case insensitive lookup could be added with strings.ToLower if needed
	champion, exists := Champions[apiName]
	if exists {
		return champion
	}
	return nil
}

// InitializeChampions loads champion data into the global map for quick access
func InitializeChampions(setData *TFTSetData) {
	Champions = make(map[string]*Champion, len(setData.SetData[0].Champions))

	// Populate the map from your set data
	for i, champion := range setData.SetData[0].Champions {
		// Store a pointer to the champion in the map
		Champions[champion.ApiName] = &setData.SetData[0].Champions[i]
	}
}
