package managers

import (
	"fmt"
	"log"
	"reflect"

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
// It also adds specific effect components for dynamic items.
func (em *EquipmentManager) AddItemToChampion(champion ecs.Entity, itemApiName string) error {
	// Get the item data by API name
	item := data.GetItemByApiName(itemApiName)
	if item == nil {
		return fmt.Errorf("item with API name '%s' not found", itemApiName)
	}

	championInfo, ok := em.world.GetChampionInfo(champion)
	if !ok {
		// It's often better to ensure ChampionInfo exists before calling this
		log.Printf("Warning: Champion %d has no ChampionInfo component when adding item %s", champion, itemApiName)
		// Decide if this should be a hard error or just a log
		// return fmt.Errorf("champion %d has no ChampionInfo component", champion)
	}
	championName := fmt.Sprintf("Entity %d", champion)
	if championInfo != nil {
		championName = championInfo.Name
	}

	// Get the Equipment component
	equipment, ok := em.world.GetEquipment(champion)
	if !ok {
		return fmt.Errorf("champion %s has no Equipment component", championInfo.Name)
	}

	// Attempt to add the item to the equipment component's list
	if !equipment.HasItemSlots() {
		return fmt.Errorf("no space to add item %s to champion %s", item.ApiName, championName)
	}

	// Check for unique constraint *before* adding
	if equipment.IsDuplicateUniqueItem(item.ApiName) {
		return fmt.Errorf("item %s is unique and already equipped on champion %s", item.ApiName, championName)
	}

	// Add the item to the component
	err := equipment.AddItem(item) // This adds the *data.Item pointer
	if err != nil {
		return fmt.Errorf("failed to add item %s to champion %s: %w", item.ApiName, championName, err)
	}
	log.Printf("Adding item '%s' to champion %s and updating item effects.", itemApiName, championName)

	// --- Add Specific Effect Components for Dynamic Items ---
	switch itemApiName {
	case data.TFT_Item_ArchangelsStaff:
		if _, exists := em.world.GetArchangelsEffect(champion); !exists {
			// Fetch values from item data
			interval := item.Effects["IntervalSeconds"] // Default to 0 if not found
			apPerStack := item.Effects["APPerInterval"] // Default to 0 if not found

			// Call the updated constructor with fetched values
			archangelsEffect := effects.NewArchangelsEffect(interval, apPerStack)
			err := em.world.AddComponent(champion, archangelsEffect)
			if err != nil {
				log.Printf("Warning: Failed to add ArchangelsEffect component for champion %s: %v", championName, err)
			} else {
				log.Printf("Added ArchangelsEffect component to champion %s (Interval: %.1f, AP/Stack: %.1f)",
					championName, interval, apPerStack)
			}
		}
	case data.TFT_Item_Quicksilver:
		if _, exists := em.world.GetQuicksilverEffect(champion); !exists {
			// Fetch values from item data
			duration := item.Effects["SpellShieldDuration"] // Default to 0 if not found
			procAS := item.Effects["ProcAttackSpeed"]       // Default to 0 if not found
			procInterval := item.Effects["ProcInterval"]    // Default to 0 if not found

			quicksilverEffect := effects.NewQuicksilverEffect(duration, procAS, procInterval)
			err := em.world.AddComponent(champion, quicksilverEffect)
			if err != nil {
				log.Printf("Warning: Failed to add QuicksilverEffect component for champion %s: %v", championName, err)
			} else {
				log.Printf("Added QuicksilverEffect component to champion %s (Duration: %.1f, ProcAS: %.2f, Interval: %.1f)",
					championName, duration, procAS, procInterval)
				// TODO: Add IsImmuneToCC marker component if implemented
				// em.world.AddComponent(champion, effects.IsImmuneToCC{})
			}
		}
		// Add cases for other dynamic items that need specific components
	}

	// Calculate the item stats and apply them to the champion, update ItemEffect component
	err = em.calculateAndUpdateItemEffects(champion)
	if err != nil {
		return fmt.Errorf("failed to calculate item effects for champion %s: %w", championName, err)
	}

	return nil
}

// RemoveItemFromChampion removes an item from a champion's equipment by its API name.
// It also removes associated effect components for specific dynamic items.
func (em *EquipmentManager) RemoveItemFromChampion(champion ecs.Entity, itemApiName string) error {
	championInfo, ok := em.world.GetChampionInfo(champion)
	if !ok {
		log.Printf("Warning: Champion %d has no ChampionInfo component when removing item %s", champion, itemApiName)
	}
	championName := fmt.Sprintf("Entity %d", champion) // Default name
	if championInfo != nil {
		championName = championInfo.Name
	}

	// Get the Equipment component
	equipment, ok := em.world.GetEquipment(champion)
	if !ok {
		return fmt.Errorf("champion %s has no Equipment component, cannot remove item", championName)
	}

	// Attempt to remove the item from the equipment component's list
	if !equipment.RemoveItem(itemApiName) { // This removes by API name
		return fmt.Errorf("item %s not found in champion %s's equipment", itemApiName, championName)
	}
	log.Printf("Removed item '%s' from champion %s's equipment component.", itemApiName, championName)

	// --- Remove Specific Effect Components for Dynamic Items ---
	switch itemApiName {
	case data.TFT_Item_ArchangelsStaff:
		if _, exists := em.world.GetArchangelsEffect(champion); exists {
			em.world.RemoveComponent(champion, reflect.TypeOf(effects.ArchangelsEffect{}))
			log.Printf("Removed ArchangelsEffect component from champion %s", championName)
		}
	case data.TFT_Item_Quicksilver:
		if _, exists := em.world.GetQuicksilverEffect(champion); exists {
			em.world.RemoveComponent(champion, reflect.TypeOf(effects.QuicksilverEffect{}))
			log.Printf("Removed QuicksilverEffect component from champion %s", championName)
			// TODO: Remove IsImmuneToCC marker component if implemented
			// em.world.RemoveComponent(champion, reflect.TypeOf(effects.IsImmuneToCC{}))
		}
		// Add cases for other dynamic items
	}

	// --- Update Static Item Effects ---
	log.Printf("Updating static item effects for champion %s after removing %s.", championName, itemApiName)
	err := em.calculateAndUpdateItemEffects(champion) // Recalculate remaining static passive stats
	if err != nil {
		log.Printf("Error updating static item effects for champion %s after removing %s: %v", championName, itemApiName, err)
		// return fmt.Errorf("failed to calculate item effects for champion %s: %w", championName, err) // Decide if this should be fatal
	}

	return nil
}

// calculateAndUpdateItemEffects calculates the total passive stats from equipped items
// and updates the champion's ItemStaticEffect component.
func (em *EquipmentManager) calculateAndUpdateItemEffects(champion ecs.Entity) error {
	championInfo, ok := em.world.GetChampionInfo(champion)
	if !ok {
		return fmt.Errorf("champion %d has no ChampionInfo component", champion)
	}
	championName := championInfo.Name

	equipment, ok := em.world.GetEquipment(champion)
	if !ok {
		// This shouldn't happen if called after ensuring equipment exists, but good practice to check.
		return fmt.Errorf("cannot calculate item effects: champion %s has no Equipment component", championName)
	}

	// Get or create the ItemStaticEffect component FIRST
	itemEffect, ok := em.world.GetItemEffect(champion)
	if !ok {
		// If no ItemStaticEffect component exists, create a new one
		newItemEffect := effects.NewItemStaticEffect()
		err := em.world.AddComponent(champion, newItemEffect)
		if err != nil {
			return fmt.Errorf("failed to add ItemStaticEffect component to champion %s: %w", championName, err)
		}
		itemEffect = newItemEffect // Use the newly added component
		log.Printf("Created new ItemStaticEffect component for champion %s.", championName)
	}

	// Reset the aggregated stats regardless of whether items exist.
	// This ensures stats are cleared when the last item is removed.
	itemEffect.ResetStats()
	log.Printf("Reset ItemStaticEffect stats for champion %s.", championName)

	// --- Handle the case where there are no items ---
	if len(equipment.Items) == 0 {
		log.Printf("Champion %s has no items equipped. Item effects reset.", championName)
		// No error, just return after resetting stats.
		return nil
	}

	// --- Process items if they exist ---
	log.Printf("Champion %s has %d items equipped. Calculating effects...", championName, len(equipment.Items))

	// Iterate through all items in the equipment and aggregate their stats
	for _, item := range equipment.GetAllItems() { // Use GetAllItems which returns *data.Item pointers
		if item == nil || item.Effects == nil {
			log.Printf("Warning: Skipping item with nil data or nil effects in equipment for champion %s", championName)
			continue
		}

		log.Printf("Processing static effects for item %s for champion %s", item.ApiName, championName)

		// Add stats from this item to the aggregate
		for statName, value := range item.Effects {
			// Only process static stats here. Dynamic effects are handled by their systems.
			switch statName {
			case "Health":
				itemEffect.AddBonusHealth(value)
			case "BonusPercentHP":
				itemEffect.AddBonusPercentHp(value)
			case "Mana":
				itemEffect.AddBonusInitialMana(value)
			case "Armor":
				itemEffect.AddBonusArmor(value)
			case "MagicResist":
				itemEffect.AddBonusMR(value)
			case "AD":
				itemEffect.AddBonusPercentAD(value)
			case "AP":
				itemEffect.AddBonusAP(value)
			case "AS":
				itemEffect.AddBonusPercentAttackSpeed(value / 100)
			case "CritChance":
				itemEffect.AddBonusCritChance(value / 100)
			case "BonusDamage":
				itemEffect.AddBonusDamageAmp(value)
			case "CritDamageToGive": // specific to IE & JG
				itemEffect.AddCritDamageToGive(value)
			// Add other known static stats...
			default:
				log.Printf("Warning: Unrecognized item effect stat '%s' for item %s", statName, item.ApiName)
			}
		}
	}

	log.Printf("Finished calculating static item effects for champion %s.", championName)
	return nil
}
