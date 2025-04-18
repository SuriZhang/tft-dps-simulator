package data

import (
	"fmt"
	// "strings" // Import if case-insensitive lookup is needed
)

// Global variable to store items for quick lookup by API Name
var SetActiveAugments map[string]*Item

// GetAugmentByApiName returns an item by its API name or nil if not found
func GetAugmentByApiName(apiName string) *Item {
	// Case insensitive lookup could be added with strings.ToLower if needed
	item, exists := SetActiveAugments[apiName]
	if exists {
		return item
	}
	return nil
}

// InitializeItems loads item data into the global map for quick access.
// It assumes the relevant set data (containing items) is at index 0 after loading.
func InitializeSetActiveAugments(setData *TFTSetData, filePath string) error {
	allItemsData, err := LoadItemDataFromFile(filePath)
	if err != nil {
		return fmt.Errorf("error loading item data: %v", err)
	}

	setActiveItemNames := setData.SetData[0].SetItems

	SetActiveAugments = make(map[string]*Item)
	// Create a map for faster lookup of all items by API name
	allItemsAndAugments = make(map[string]*Item, len(allItemsData))
	for i := range allItemsData {
		// Store pointer to the item in the map
		allItemsAndAugments[allItemsData[i].ApiName] = &allItemsData[i]
	}

	// Iterate through the API names of the items active in the set
	for _, apiName := range setActiveItemNames {
		// Look up the item in the pre-built map
		if item, found := allItemsAndAugments[apiName]; found {
			SetActiveAugments[apiName] = item // Add the found item pointer to the result map
		} else {
			// Handle case where an active item name is not found in the loaded item data
			fmt.Printf("Warning: Set active item '%s' not found in loaded item data from %s\n", apiName, filePath)
		}
	}
	return nil
}

// GetAugmentByName searches by display name (less reliable for uniqueness)
func GetAugmentByName(name string) *Item {
	for _, item := range SetActiveAugments {
		if item.Name == name {
			return item // Return the first match
		}
	}
	return nil
}
