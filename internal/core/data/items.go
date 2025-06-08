package data

import (
	"fmt"
	"log"
)

// Global variable to store items for quick lookup by API Name
var SetActiveItems map[string]*Item

var allItemsAndAugments map[string]*Item // Package private, key is API name

const (
	TFT_Item_BFSword                      = "TFT_Item_BFSword"
	TFT_Item_ChainVest                    = "TFT_Item_ChainVest"
	TFT_Item_GiantsBelt                   = "TFT_Item_GiantsBelt"
	TFT_Item_NeedlesslyLargeRod           = "TFT_Item_NeedlesslyLargeRod"
	TFT_Item_NegatronCloak                = "TFT_Item_NegatronCloak"
	TFT_Item_RecurveBow                   = "TFT_Item_RecurveBow"
	TFT_Item_SparringGloves               = "TFT_Item_SparringGloves"
	TFT_Item_Spatula                      = "TFT_Item_Spatula"
	TFT_Item_TearOfTheGoddess             = "TFT_Item_TearOfTheGoddess"
	TFT_Item_RabadonsDeathcap             = "TFT_Item_RabadonsDeathcap"
	TFT_Item_Deathblade                   = "TFT_Item_Deathblade"
	TFT_Item_WarmogsArmor                 = "TFT_Item_WarmogsArmor"
	TFT_Item_ArchangelsStaff              = "TFT_Item_ArchangelsStaff"
	TFT_Item_Quicksilver                  = "TFT_Item_Quicksilver"
	TFT_Item_TitansResolve                = "TFT_Item_TitansResolve"
	TFT_Item_GuinsoosRageblade            = "TFT_Item_GuinsoosRageblade"
	TFT_Item_SpiritVisage                 = "TFT_Item_Redemption"
	TFT_Item_InfinityEdge                 = "TFT_Item_InfinityEdge"
	TFT_Item_JeweledGauntlet              = "TFT_Item_JeweledGauntlet"
	TFT_Item_KrakensFury                  = "TFT_Item_RunaansHurricane"
	TFT_Item_SpearOfShojin                = "TFT_Item_SpearOfShojin"
	TFT_Item_BlueBuff                     = "TFT_Item_BlueBuff"
	TFT_Item_Artifact_NavoriFlickerblades = "TFT_Item_Artifact_NavoriFlickerblades"
	TFT_Item_NashorsTooth                 = "TFT_Item_Leviathan"
	TFT_Item_VoidStaff                    = "TFT_Item_StatikkShiv"
	TFT_Item_RedBuff                      = "TFT_Item_RapidFireCannon"
)

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
