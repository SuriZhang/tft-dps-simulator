package data

import (
	"fmt"
	"log"
)

// Global variable to store items for quick lookup by API Name
var SetActiveItems map[string]*Item

var allItemsAndAugments map[string]*Item // Package private, key is API name

type ItemClassification string

const (
	StaticItem     ItemClassification = "static"
	DynamicTime    ItemClassification = "dynamic_time"
	DynamicEvent   ItemClassification = "dynamic_event"
	DynamicComplex ItemClassification = "dynamic_complex"
	UnknownItem    ItemClassification = "unknown" // Default or for items not yet classified

	TFT_Item_BFSword            = "TFT_Item_BFSword"
	TFT_Item_ChainVest          = "TFT_Item_ChainVest"
	TFT_Item_GiantsBelt         = "TFT_Item_GiantsBelt"
	TFT_Item_NeedlesslyLargeRod = "TFT_Item_NeedlesslyLargeRod"
	TFT_Item_NegatronCloak      = "TFT_Item_NegatronCloak"
	TFT_Item_RecurveBow         = "TFT_Item_RecurveBow"
	TFT_Item_SparringGloves     = "TFT_Item_SparringGloves"
	TFT_Item_Spatula            = "TFT_Item_Spatula"
	TFT_Item_TearOfTheGoddess   = "TFT_Item_TearOfTheGoddess"
	TFT_Item_RabadonsDeathcap   = "TFT_Item_RabadonsDeathcap"
	TFT_Item_Deathblade         = "TFT_Item_Deathblade"
	TFT_Item_WarmogsArmor       = "TFT_Item_WarmogsArmor"
	TFT_Item_ArchangelsStaff    = "TFT_Item_ArchangelsStaff"
	TFT_Item_Quicksilver        = "TFT_Item_Quicksilver"
	TFT_Item_TitansResolve      = "TFT_Item_TitansResolve"
	TFT_Item_GuinsoosRageblade  = "TFT_Item_GuinsoosRageblade"
	TFT_Item_Redemption         = "TFT_Item_Redemption"
	TFT_Item_InfinityEdge       = "TFT_Item_InfinityEdge"
	TFT_Item_JeweledGauntlet    = "TFT_Item_JeweledGauntlet"
)

var itemClassificationMap = map[string]ItemClassification{
	// --- Static Items ---
	// Components (Add all component ApiNames here)
	"TFT_Item_BFSword":            StaticItem,
	"TFT_Item_ChainVest":          StaticItem,
	"TFT_Item_GiantsBelt":         StaticItem,
	"TFT_Item_NeedlesslyLargeRod": StaticItem,
	"TFT_Item_NegatronCloak":      StaticItem,
	"TFT_Item_RecurveBow":         StaticItem,
	"TFT_Item_SparringGloves":     StaticItem,
	"TFT_Item_Spatula":            StaticItem,
	"TFT_Item_TearOfTheGoddess":   StaticItem,
	// Completed Static Items
	"TFT_Item_RabadonsDeathcap": StaticItem,
	"TFT_Item_Deathblade":       StaticItem,
	"TFT_Item_WarmogsArmor":     StaticItem,

	// --- Dynamic Items (Examples) ---
	"TFT_Item_ArchangelsStaff":   DynamicTime,    // Gains AP over time
	"TFT_Item_TitansResolve":     DynamicEvent,   // Gains stats on hit/being hit
	"TFT_Item_GuinsoosRageblade": DynamicComplex, // Stacks AS on attack
	"TFT_Item_Quicksilver":       DynamicTime,    // CC immunity duration
	"TFT_Item_Redemption":        DynamicTime,    // AoE heal after delay (time/event mix) - classify carefully
	"TFT_Item_AdaptiveHelm":      DynamicEvent,   // Different effect based on row (event)

	// Add all other items and their classifications...
}

// Function to get classification (can be part of factory or a helper)
func GetItemClassification(apiName string) ItemClassification {
	classification, found := itemClassificationMap[apiName]
	if !found {
		// Handle items not explicitly listed - maybe default to static or log a warning
		// log.Printf("Warning: Classification not found for item %s, defaulting to UnknownItem", apiName)
		return UnknownItem
	}
	return classification
}

// GetItemByApiName returns an item by its API name or nil if not found
func GetItemByApiName(apiName string) *Item {
	// Case insensitive lookup could be added with strings.ToLower if needed
	item, exists := SetActiveItems[apiName]
	if exists {
		return item
	}
	return nil
}

// InitializeItems loads item data into the global map for quick access.
// It assumes the relevant set data (containing items) is at index 0 after loading.
func InitializeSetActiveItems(setData *TFTSetData, filePath string) error {
	allItemsData, err := LoadItemDataFromFile(filePath)
	if err != nil {
		return fmt.Errorf("error loading item data: %v", err)
	}

	setActiveItemNames := setData.SetData[0].SetItems

	SetActiveItems = make(map[string]*Item)
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
			SetActiveItems[apiName] = item // Add the found item pointer to the result map
		} else {
			// Handle case where an active item name is not found in the loaded item data
			log.Printf("Warning: Set active item '%s' not found in loaded item data from %s\n", apiName, filePath)
		}
	}
	return nil
}

// GetItemByName searches by display name (less reliable for uniqueness)
func GetItemByName(name string) *Item {
	for _, item := range SetActiveItems {
		if item.Name == name {
			return item // Return the first match
		}
	}
	return nil
}
