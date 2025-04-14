package managers

import (
	"fmt"
	"log"

	"github.com/suriz/tft-dps-simulator/components/effects"
	"github.com/suriz/tft-dps-simulator/data"
	"github.com/suriz/tft-dps-simulator/ecs"
)

// EquipmentManager handles adding/removing items and calculating their effects.
type EquipmentManager struct {
    world *ecs.World
}

// NewEquipmentManager creates a new EquipmentManager.
func NewEquipmentManager(world *ecs.World) *EquipmentManager {
    return &EquipmentManager{world: world}
}


// AddItemToChampion adds an item to a champion's equipment if there's space.
// It now returns an error if the item cannot be added.
func (em *EquipmentManager) AddItemToChampion(champion ecs.Entity, itemApiName string) error {
	// Get the item data by API name
	item := data.GetItemByApiName(itemApiName)
	if item == nil {
		return fmt.Errorf("item with API name '%s' not found", itemApiName)
	}

	championInfo, ok := em.world.GetChampionInfo(champion)
	if !ok {
		return fmt.Errorf("champion %d has no ChampionInfo component", champion)
	}
	// Get the Equipment component
	equipment, ok := em.world.GetEquipment(champion)
	if !ok {
		return fmt.Errorf("champion %s has no Equipment component", championInfo.Name)
	}

	// Attempt to add the item to the equipment
	if !equipment.HasItemSlots() {
		return fmt.Errorf("no space to add item %s to champion %s", item.ApiName, championInfo.Name)
	}

	if equipment.IsDuplicateUniqueItem(item.ApiName) {
		return fmt.Errorf("item %s is unique and already equipped on champion %s", item.ApiName, championInfo.Name)
	}

	equipment.AddItem(item) 
	log.Printf("Adding item '%s' to champion %s and updating item effects.", itemApiName, championInfo.Name)
	// Calculate the item stats and apply them to the champion, update ItemEffect component
	err := em.calculateAndUpdateItemEffects(champion)
	if err != nil {
		return fmt.Errorf("failed to calculate item effects for champion %s: %w", championInfo.Name, err)
	}

	return nil
}

// RemoveItemFromChampion removes an item from a champion's equipment by its API name.
// It now returns an error if the item cannot be removed.
func (em *EquipmentManager) RemoveItemFromChampion(champion ecs.Entity, itemApiName string) error {
	championInfo, ok := em.world.GetChampionInfo(champion)
	if !ok {
		return fmt.Errorf("champion %d has no ChampionInfo component", champion)
	}

	// Get the Equipment component
	equipment, ok := em.world.GetEquipment(champion)
	if !ok {
		return fmt.Errorf("champion %s has no Equipment component", championInfo.Name)
	}

	// Attempt to remove the item from the equipment
	if !equipment.RemoveItem(itemApiName) {
		return fmt.Errorf("item %s not found in champion %s's equipment", itemApiName, championInfo.Name)
	}

	log.Printf("Removing item '%s' from champion %s and updating item effects...", itemApiName, championInfo.Name)
	// Remove item effects from champion stats, update ItemEffect component
	err := em.calculateAndUpdateItemEffects(champion)
	if err != nil {
		return fmt.Errorf("failed to calculate item effects for champion %s: %w", championInfo.Name, err)
	}

	return nil
}

// calculateAndUpdateItemEffects calculates the total passive stats from equipped items
// and updates the champion's ItemEffect component.
func (em *EquipmentManager) calculateAndUpdateItemEffects(champion ecs.Entity) error {
	championInfo, ok := em.world.GetChampionInfo(champion)
	if !ok {
		return fmt.Errorf("champion %d has no ChampionInfo component", champion)
	}

	equipment, ok := em.world.GetEquipment(champion)
	if !ok {
		// This shouldn't happen if called after ensuring equipment exists, but good practice to check.
		return fmt.Errorf("cannot caculate item effects: champion %s has no Equipment component", championInfo.Name)
	}

	if len(equipment.Items) == 0 {
		return fmt.Errorf("champion %s has no items equipped", championInfo.Name)
	}

	log.Printf("Champion %s has %d items equipped.", championInfo.Name, len(equipment.Items))

	// Get or create the ItemEffect component
	itemEffect, ok := em.world.GetItemEffect(champion)
	if !ok {
		// If no ItemEffect component exists, create a new one
		newItemEffect := effects.NewItemStaticEffect()
		err := em.world.AddComponent(champion, newItemEffect)
		if err != nil {
			return fmt.Errorf("failed to add ItemEffect component to champion %s: %w", championInfo.Name, err)
		}
		itemEffect = newItemEffect // Use the newly added component
	}

	// Reset the aggregated stats before caculating
	itemEffect.ResetStats() // Add this method to ItemEffect component

	// Iterate through all items in the equipment and aggregate their stats
	for _, item := range equipment.GetAllItems() {
		if item == nil || item.Effects == nil {
			log.Printf("Warning: Skipping item with nil data or nil effects in equipment for champion %s", championInfo.Name)
			continue
		}

		log.Printf("Processing item %s for champion %s", item.ApiName, championInfo.Name)

		// Add stats from this item to the aggregate
		for statName, value := range item.Effects {
			switch statName {
			case "Health":
				itemEffect.AddBonusHealth(value)
				log.Printf("Adding %f bonus health to champion %s from item %s", value, championInfo.Name, item.ApiName)
			case "BonusPercentHP":
				itemEffect.AddBonusPercentHp(value)
				log.Printf("Adding %f bonus percent HP to champion %s from item %s", value, championInfo.Name, item.ApiName)
			case "Mana":
				itemEffect.AddBonusInitialMana(value)
				log.Printf("Adding %f bonus initial mana to champion %s from item %s", value, championInfo.Name, item.ApiName)
			case "Armor":
				itemEffect.AddBonusArmor(value)
				log.Printf("Adding %f bonus armor to champion %s from item %s", value, championInfo.Name, item.ApiName)
			case "MagicResist":
				itemEffect.AddBonusMR(value)
				log.Printf("Adding %f bonus magic resist to champion %s from item %s", value, championInfo.Name, item.ApiName)
			case "AD":
				itemEffect.AddBonusPercentAD(value)
				log.Printf("Adding %f bonus percent AD to champion %s from item %s", value, championInfo.Name, item.ApiName)
			case "AP":
				itemEffect.AddBonusAP(value)
				log.Printf("Adding %f bonus AP to champion %s from item %s", value, championInfo.Name, item.ApiName)
			case "AS":
				itemEffect.AddBonusPercentAttackSpeed(value / 100)
				log.Printf("Adding %f bonus percent attack speed to champion %s from item %s", value/100, championInfo.Name, item.ApiName)
			case "CritChance":
				itemEffect.AddBonusCritChance(value / 100)
				log.Printf("Adding %f bonus crit chance to champion %s from item %s", value/100, championInfo.Name, item.ApiName)
			case "BonusDamage":
				itemEffect.AddBonusDamageAmp(value)
				log.Printf("Adding %f bonus damage amp to champion %s from item %s", value, championInfo.Name, item.ApiName)
			case "CritDamageToGive": // specific to IE & JG
				itemEffect.AddCritDamageToGive(value)
				log.Printf("Adding %f crit damage to give to champion %s from item %s", value, championInfo.Name, item.ApiName)
			// Add cases for other stats as needed...
			default:
				// Optional: Log or handle unrecognized stat names
				log.Printf("Warning: Unrecognized item effect stat '%s' for item %s", statName, item.ApiName)
			}
			// Add other stats as needed...
		}
		log.Printf("Successfully processed item '%s' to champion %s and updated item effects.", item.ApiName, championInfo.Name)
	}

	// Optional: Add/Remove marker components based on items (for future systems)
	// em.updateItemMarkerComponents(champion, equipment.Items)

	return nil
}
