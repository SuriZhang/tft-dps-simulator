package data

// Global variable to store traits for quick lookup
var Traits map[string]*Trait

// GetTraitByName returns a trait by name or nil if not found
func GetTraitByName(name string) *Trait {
    // Case insensitive lookup could be added with strings.ToLower if needed
    trait, exists := Traits[name]
    if exists {
        return trait
    }
    return nil
}

// InitializeTraits loads trait data into the global map for quick access
func InitializeTraits(setData *TFTSetData) {
    Traits = make(map[string]*Trait)
    
    // Populate the map from your set data
    for i, trait := range setData.SetData[0].Traits {
        // Store a pointer to the trait in the map
        Traits[trait.Name] = &setData.SetData[0].Traits[i]
    }
}