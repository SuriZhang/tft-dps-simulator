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

	err := equipment.AddItem(item)
	if err != nil {
		return fmt.Errorf("failed to add item %s to champion %s: %w", item.ApiName, championInfo.Name, err)
	}
	log.Printf("Adding item '%s' to champion %s and updating item effects.", itemApiName, championInfo.Name)
	// Calculate the item stats and apply them to the champion, update ItemEffect component
	err = em.calculateAndUpdateItemEffects(champion)
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
		return fmt.Errorf("cannot calculate item effects: champion %s has no Equipment component", championInfo.Name)
	}

	// Get or create the ItemEffect component FIRST
	itemEffect, ok := em.world.GetItemEffect(champion)
	if !ok {
		// If no ItemEffect component exists, create a new one
		newItemEffect := effects.NewItemStaticEffect()
		err := em.world.AddComponent(champion, newItemEffect)
		if err != nil {
			return fmt.Errorf("failed to add ItemEffect component to champion %s: %w", championInfo.Name, err)
		}
		itemEffect = newItemEffect // Use the newly added component
		log.Printf("Created new ItemEffect component for champion %s.", championInfo.Name)
	}

	// Reset the aggregated stats regardless of whether items exist.
	// This ensures stats are cleared when the last item is removed.
	itemEffect.ResetStats()
	log.Printf("Reset ItemEffect stats for champion %s.", championInfo.Name)

	// --- Handle the case where there are no items ---
	if len(equipment.Items) == 0 {
		log.Printf("Champion %s has no items equipped. Item effects reset.", championInfo.Name)
		// No error, just return after resetting stats.
		return nil 
	}

	// --- Process items if they exist ---
	log.Printf("Champion %s has %d items equipped. Calculating effects...", championInfo.Name, len(equipment.Items))

	// Iterate through all items in the equipment and aggregate their stats
	for _, item := range equipment.GetAllItems() { // Use GetAllItems which returns a copy
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
				// log.Printf("Adding %f bonus health to champion %s from item %s", value, championInfo.Name, item.ApiName) // Reduce log verbosity
			case "BonusPercentHP":
				itemEffect.AddBonusPercentHp(value)
				// log.Printf("Adding %f bonus percent HP to champion %s from item %s", value, championInfo.Name, item.ApiName)
			case "Mana":
				itemEffect.AddBonusInitialMana(value)
				// log.Printf("Adding %f bonus initial mana to champion %s from item %s", value, championInfo.Name, item.ApiName)
			case "Armor":
				itemEffect.AddBonusArmor(value)
				// log.Printf("Adding %f bonus armor to champion %s from item %s", value, championInfo.Name, item.ApiName)
			case "MagicResist":
				itemEffect.AddBonusMR(value)
				// log.Printf("Adding %f bonus magic resist to champion %s from item %s", value, championInfo.Name, item.ApiName)
			case "AD":
				itemEffect.AddBonusPercentAD(value)
				// log.Printf("Adding %f bonus percent AD to champion %s from item %s", value, championInfo.Name, item.ApiName)
			case "AP":
				itemEffect.AddBonusAP(value)
				// log.Printf("Adding %f bonus AP to champion %s from item %s", value, championInfo.Name, item.ApiName)
			case "AS":
				itemEffect.AddBonusPercentAttackSpeed(value / 100)
				// log.Printf("Adding %f bonus percent attack speed to champion %s from item %s", value/100, championInfo.Name, item.ApiName)
			case "CritChance":
				itemEffect.AddBonusCritChance(value / 100)
				// log.Printf("Adding %f bonus crit chance to champion %s from item %s", value/100, championInfo.Name, item.ApiName)
			case "BonusDamage":
				itemEffect.AddBonusDamageAmp(value)
				// log.Printf("Adding %f bonus damage amp to champion %s from item %s", value, championInfo.Name, item.ApiName)
			case "CritDamageToGive": // specific to IE & JG
				itemEffect.AddCritDamageToGive(value)
				// log.Printf("Adding %f crit damage to give to champion %s from item %s", value, championInfo.Name, item.ApiName)
			default:
				log.Printf("Warning: Unrecognized item effect stat '%s' for item %s", statName, item.ApiName)
			}
		}
		// log.Printf("Successfully processed item '%s' to champion %s and updated item effects.", item.ApiName, championInfo.Name) // Reduce log verbosity
	}

	log.Printf("Finished calculating item effects for champion %s.", championInfo.Name)
	return nil
}
