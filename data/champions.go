package data

// Global variable to store champions for quick lookup
var Champions map[string]*Champion

// GetChampionByName returns a champion by name or nil if not found
func GetChampionByName(name string) *Champion {
    // Case insensitive lookup could be added with strings.ToLower if needed
    champion, exists := Champions[name]
    if exists {
        return champion
    }
    return nil
}

// InitializeChampions loads champion data into the global map for quick access
func InitializeChampions(setData *TFTSetData) {
    Champions = make(map[string]*Champion)
    
    // Populate the map from your set data
    for i, champion := range setData.SetData[0].Champions {
        // Store a pointer to the champion in the map
        Champions[champion.Name] = &setData.SetData[0].Champions[i]
    }
}