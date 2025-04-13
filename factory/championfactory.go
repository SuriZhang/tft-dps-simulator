package factory

import (
	"fmt"
	"log"

	"github.com/suriz/tft-dps-simulator/components"
	"github.com/suriz/tft-dps-simulator/data"
	"github.com/suriz/tft-dps-simulator/ecs"
)

// ChampionFactory creates champion entities from champion data.
type ChampionFactory struct {
	world *ecs.World
}

// NewChampionFactory creates a new ChampionFactory.
func NewChampionFactory(world *ecs.World) *ChampionFactory {
	return &ChampionFactory{
		world: world,
	}
}

// StarMultiplier returns the stat multiplier for a given star level.
func StarMultiplier(starLevel int) float64 {
	switch starLevel {
	case 1:
		return 1.0
	case 2:
		return 1.8
	case 3:
		return 3.24
	default:
		fmt.Printf("Warning: Invalid star level %d requested, using 1-star multiplier.\n", starLevel)
		return 1.0
	}
}

// CreateChampion creates a champion entity with components from data.
// It now returns an error if adding any component fails.
func (cf *ChampionFactory) CreateChampion(championData data.Champion, starLevel int) (ecs.Entity, error) {
	// Create entity
	entity := cf.world.CreateEntity()
	var err error // Variable to hold potential errors

	// Apply star level multiplier
	multiplier := StarMultiplier(starLevel)

	// --- Add components, checking for errors after each AddComponent call ---

	// Health
	healthComp := components.NewHealth(
		championData.Stats.HP*multiplier,
		championData.Stats.Armor,
		championData.Stats.MagicResist,
	)
	err = cf.world.AddComponent(entity, healthComp)
	if err != nil {
		return 0, fmt.Errorf("failed to add Health component to %s: %w", championData.Name, err)
	}

	// Attack
	attackComp := components.NewAttack(
		championData.Stats.Damage*multiplier,
		championData.Stats.AttackSpeed,
		championData.Stats.Range,
		championData.Stats.CritChance,
		championData.Stats.CritMultiplier,
	)
	err = cf.world.AddComponent(entity, attackComp)
	if err != nil {
		return 0, fmt.Errorf("failed to add Attack component to %s: %w", championData.Name, err)
	}

	// Mana
	manaComp := components.NewMana(
		championData.Stats.Mana,
		championData.Stats.InitialMana,
	)
	err = cf.world.AddComponent(entity, manaComp)
	if err != nil {
		return 0, fmt.Errorf("failed to add Mana component to %s: %w", championData.Name, err)
	}

	// ChampionInfo
	infoComp := components.NewChampionInfo(
		championData.ApiName,
		championData.Name,
		championData.Cost,
		starLevel,
	)
	err = cf.world.AddComponent(entity, infoComp)
	if err != nil {
		return 0, fmt.Errorf("failed to add ChampionInfo component to %s: %w", championData.Name, err)
	}

	// Traits (only if they exist)
	if len(championData.Traits) > 0 {
		traitsComp := components.NewTraits(championData.Traits)
		err = cf.world.AddComponent(entity, traitsComp)
		if err != nil {
			return 0, fmt.Errorf("failed to add Traits component to %s: %w", championData.Name, err)
		}
	}

	// create empty inventory
	err = cf.world.AddComponent(entity, components.NewEquipment())
	if err != nil {
		return 0, fmt.Errorf("failed to add Inventory component to %s: %w", championData.Name, err)
	}

	err = cf.world.AddComponent(entity, components.NewPosition(0, 0))
	if err != nil {
		return 0, fmt.Errorf("failed to add Position component to %s: %w", championData.Name, err)
	}
	// Add other essential components here (Position, Team are usually added later/externally)

	// If we reached here, all essential components were added successfully
	return entity, nil
}

// CreateChampionByName creates a champion entity by searching for it by name.
// It now propagates errors from CreateChampion.
func (cf *ChampionFactory) CreateChampionByName(name string, starLevel int, team int) (ecs.Entity, error) {
	// Find champion data by name using the function from the data package
	championData := data.GetChampionByName(name)
	if championData == nil {
		return 0, fmt.Errorf("champion data for '%s' not found", name)
	}

	// Create entity using the found data, passing any error up
	entity, err := cf.CreateChampion(*championData, starLevel)
	if err != nil {
		return 0, err
	}

	// assign a team
	err = cf.world.AddComponent(entity, components.NewTeam(team))
	if err != nil {
		return 0, fmt.Errorf("failed to add Team component to %s: %w", championData.Name, err)
	}

	return entity, nil
}

// CreatePlayerChampion creates a champion entity for the player team (team ID 0).
func (cf *ChampionFactory) CreatePlayerChampion(name string, starLevel int) (ecs.Entity, error) {
	return cf.CreateChampionByName(name, starLevel, 0)
}

// CreateEnemyChampion creates a champion entity for the enemy team (team ID 1).
func (cf *ChampionFactory) CreateEnemyChampion(name string, starLevel int) (ecs.Entity, error) {
	return cf.CreateChampionByName(name, starLevel, 1)
}

// AddItemToChampion adds an item to a champion's equipment if there's space.
// It now returns an error if the item cannot be added.
func (cf *ChampionFactory) AddItemToChampion(champion ecs.Entity, itemApiName string) error {
	// Get the item data by API name
	item := data.GetItemByApiName(itemApiName)
	if item == nil {
		return fmt.Errorf("item with API name '%s' not found", itemApiName)
	}

	championInfo, ok := cf.world.GetChampionInfo(champion)
	if !ok {
		return fmt.Errorf("champion %d has no ChampionInfo component", champion)
	}
	// Get the Equipment component
	equipment, ok := cf.world.GetEquipment(champion)
	if !ok {
		return fmt.Errorf("champion %s has no Equipment component", championInfo.Name)
	}

	// Attempt to add the item to the equipment
	if !equipment.HasItemSlots(item) {
		return fmt.Errorf("no space to add item %s to champion %s", item.ApiName, championInfo.Name)
	}

	if equipment.IsDuplicateUniqueItem(item.ApiName) {
		return fmt.Errorf("item %s is unique and already equipped on champion %s", item.ApiName, championInfo.Name)
	}

	log.Printf("Adding item '%s' to champion %s and updating item effects.", itemApiName, championInfo.Name)
	// Calculate the item stats and apply them to the champion, update ItemEffect component
	err := cf.calculateAndUpdateItemEffects(champion)
	if err != nil {
		return fmt.Errorf("failed to calculate item effects for champion %s: %w", championInfo.Name, err)
	}

	return nil
}

// RemoveItemFromChampion removes an item from a champion's equipment by its API name.
// It now returns an error if the item cannot be removed.
func (cf *ChampionFactory) RemoveItemFromChampion(champion ecs.Entity, itemApiName string) error {
	championInfo, ok := cf.world.GetChampionInfo(champion)
	if !ok {
		return fmt.Errorf("champion %d has no ChampionInfo component", champion)
	}

	// Get the Equipment component
	equipment, ok := cf.world.GetEquipment(champion)
	if !ok {
		return fmt.Errorf("champion %s has no Equipment component", championInfo.Name)
	}

	// Attempt to remove the item from the equipment
	if !equipment.RemoveItem(itemApiName) {
		return fmt.Errorf("item %s not found in champion %s's equipment", itemApiName, championInfo.Name)
	}

	log.Printf("Removing item '%s' from champion %s and updating item effects.", itemApiName, championInfo.Name)
	// Remove item effects from champion stats, update ItemEffect component
	err := cf.calculateAndUpdateItemEffects(champion)
	if err != nil {
		return fmt.Errorf("failed to calculate item effects for champion %s: %w", championInfo.Name, err)
	}

	return nil
}

// calculateAndUpdateItemEffects calculates the total passive stats from equipped items
// and updates the champion's ItemEffect component.
func (cf *ChampionFactory) calculateAndUpdateItemEffects(champion ecs.Entity) error {
	championInfo, ok := cf.world.GetChampionInfo(champion)
	if !ok {
		return fmt.Errorf("champion %d has no ChampionInfo component", champion)
	}

	equipment, ok := cf.world.GetEquipment(champion)
	if !ok {
		// This shouldn't happen if called after ensuring equipment exists, but good practice to check.
		return fmt.Errorf("cannot caculate item effects: champion %s has no Equipment component", championInfo.Name)
	}

	if len(equipment.Items) == 0 {
		log.Printf("Champion %s has no items equipped, skipping item effect calculation.", championInfo.Name)
		return nil
	}

	log.Printf("Champion %s has %d items equipped.", championInfo.Name, len(equipment.Items))

	// Get or create the ItemEffect component
	itemEffect, ok := cf.world.GetItemEffect(champion)
	if !ok {
		// If no ItemEffect component exists, create a new one
		newItemEffect := components.NewItemEffect()
		err := cf.world.AddComponent(champion, newItemEffect)
		if err != nil {
			return fmt.Errorf("failed to add ItemEffect component to champion %s: %w", championInfo.Name, err)
		}
		itemEffect = newItemEffect // Use the newly added component
	}

	// Reset the aggregated stats before caculating
	itemEffect.ResetStats() // Add this method to ItemEffect component

	// Iterate through all items in the equipment and aggregate their stats
	for _, item := range equipment.Items {
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
	// cf.updateItemMarkerComponents(champion, equipment.Items)

	return nil
}
