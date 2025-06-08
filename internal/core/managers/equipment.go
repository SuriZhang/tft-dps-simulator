package managers

import (
	"fmt"
	"log"
	"reflect"

	"tft-dps-simulator/internal/core/components/items"
	"tft-dps-simulator/internal/core/data"
	"tft-dps-simulator/internal/core/ecs"
	"tft-dps-simulator/internal/core/entity"
)

// EquipmentManager handles adding/removing items and calculating their effects.
type EquipmentManager struct {
	world *ecs.World
}

// NewEquipmentManager creates a new EquipmentManager.
func NewEquipmentManager(world *ecs.World) *EquipmentManager {
    return &EquipmentManager{
        world:    world,
    }
}

// AddItemToChampion adds an item to a champion's equipment if there's space.
// It also adds specific effect components for dynamic items.
func (em *EquipmentManager) AddItemToChampion(champion entity.Entity, itemApiName string) error {
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
		if _, exists := em.world.GetArchangelsStaffEffect(champion); !exists {
			// Fetch values from item data
			interval := item.Effects["IntervalSeconds"] // Default to 0 if not found
			apPerStack := item.Effects["APPerInterval"] // Default to 0 if not found

			// Call the updated constructor with fetched values
			archangelsEffect := items.NewArchangelsEffect(interval, apPerStack)
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

			quicksilverEffect := items.NewQuicksilverEffect(duration, procAS, procInterval)
			err := em.world.AddComponent(champion, quicksilverEffect)
			if err != nil {
				log.Printf("Warning: Failed to add QuicksilverEffect component for champion %s: %v", championName, err)
			} else {
				log.Printf("Added QuicksilverEffect component to champion %s (Duration: %.1f, ProcAS: %.2f, Interval: %.1f)",
					championName, duration, procAS, procInterval)
				// TODO: Add IsImmuneToCC marker component if implemented
				// em.world.AddComponent(champion, itemsIsImmuneToCC{})
			}
		}
	case data.TFT_Item_TitansResolve:
		if _, exists := em.world.GetTitansResolveEffect(champion); !exists {
			// Fetch values from item data, converting StackCap to int
			maxStacks := item.Effects["StackCap"]
			adPerStack := item.Effects["StackingAD"]
			apPerStack := item.Effects["StackingSP"]
			bonusResists := item.Effects["BonusResistsAtStackCap"]

			titansEffect := items.NewTitansResolveEffect(maxStacks, adPerStack, apPerStack, bonusResists)
			err := em.world.AddComponent(champion, titansEffect)
			if err != nil {
				log.Printf("Warning: Failed to add TitansResolveEffect component for champion %s: %v", championName, err)
			} else {
				log.Printf("Added TitansResolveEffect component to champion %s (Stacks: %d, AD/s: %.2f%%, AP/s: %.1f, Res@Max: %.0f)",
					championName, int(maxStacks), adPerStack*100, apPerStack, bonusResists)
			}
		}
	case data.TFT_Item_SpiritVisage: // NOTE: New case for Spirit Visage, this may change.
	if _, exists := em.world.GetSpiritVisageEffect(champion); !exists {
			missingHealthHeal := item.Effects["MissingHealthHeal"]
			healTickRate := item.Effects["HealTickRate"]
			maxHeal := item.Effects["MaxHeal"]

			spiritVisageEffect := items.NewSpiritVisageEffect(missingHealthHeal, healTickRate, maxHeal)
			err := em.world.AddComponent(champion, spiritVisageEffect)
			if err != nil {
				log.Printf("Warning: Failed to add SpiritVisageEffect component for champion %s: %v", championName, err)
			} else {
				log.Printf("Added SpiritVisageEffect component to champion %s (Heal: %.3f%% missing HP / %.1fs)",
					championName, missingHealthHeal*100, healTickRate)
			}
		}
	case data.TFT_Item_BlueBuff:
		if _, exists := em.world.GetBlueBuffEffect(champion); !exists {
			blueBuff := items.NewBlueBuff()
			err := em.world.AddComponent(champion, blueBuff)
			if err != nil {
				log.Printf("Warning: Failed to add BlueBuff component for champion %s: %v", championName, err)
			} else {
				log.Printf("Added BlueBuff component to champion %s (ManaRefund: %.1f, DamageAmp: %.1f%%, Timer: %.1fs)",
					championName, blueBuff.ManaRefund, blueBuff.DamageAmp*100, blueBuff.TakedownTimer)
			}
		}

	// Add cases for other dynamic items that need specific components
	case data.TFT_Item_GuinsoosRageblade:
		if _, exists := em.world.GetGuinsoosRagebladeEffect(champion); !exists {
			// Fetch the correct value from item data
			asPerStack := item.Effects["AttackSpeedPerStack"]
			// TODO: fix the effect field name when the data is updated
			intervalSeconds, ok := item.Effects["IntervalSeconds"]
			if !ok {
				intervalSeconds = 1.0 // Default to 1.0 if not found
			}
			
			// Create the effect component
			ragebladeEffect := items.NewGuinsoosRagebladeEffect(intervalSeconds, asPerStack / 100)
			err := em.world.AddComponent(champion, ragebladeEffect)
			if err != nil {
				log.Printf("Warning: Failed to add GuinsoosRagebladeEffect component for champion %s: %v", championName, err)
			} else {
				log.Printf("Added GuinsoosRagebladeEffect component to champion %s (Interval: %.1f, AS/Stack: %.2f%%)",
					championName, intervalSeconds, asPerStack)
			}
		}
	
	case data.TFT_Item_KrakensFury:
        if _, exists := em.world.GetKrakensFuryEffect(champion); !exists {
            adPerStack := item.Effects["ADOnAttack"]

            krakensFuryEffect := items.NewKrakensFuryEffect(adPerStack)
            err := em.world.AddComponent(champion, krakensFuryEffect)
            if err != nil {
                log.Printf("Warning: Failed to add KrakensFuryEffect component for champion %s: %v", championName, err)
            } else {
                log.Printf("Added KrakensFuryEffect component to champion %s (AD/Stack: %.2f%%)",
                    championName, adPerStack*100)
            }
        }
	case data.TFT_Item_SpearOfShojin:
        if _, exists := em.world.GetSpearOfShojinEffect(champion); !exists {
            // Fetch the correct value from item data
            flatManaRestore := item.Effects["FlatManaRestore"]
            
            // Create the effect component
            spearEffect := items.NewSpearOfShojinEffect(flatManaRestore)
            em.world.AddComponent(champion, spearEffect)
            log.Printf("EquipmentManager: Added SpearOfShojinEffect component to Entity %d with %.1f mana restore per attack", champion, flatManaRestore)
        }
	case data.TFT_Item_Artifact_NavoriFlickerblades:
        if _, exists := em.world.GetFlickerbladeEffect(champion); !exists {
            asPerStackVal := item.Effects["ASPerStack"]
            adPerBonusVal := item.Effects["ADPerBonus"]
            apPerBonusVal := item.Effects["APPerBonus"]
            stacksPerBonusVal := item.Effects["StacksPerBonus"]

            flickerbladeEffect := items.NewFlickerbladeEffect(asPerStackVal, adPerBonusVal, apPerBonusVal, stacksPerBonusVal)
            err := em.world.AddComponent(champion, flickerbladeEffect)
            if err != nil {
                return fmt.Errorf("champion %s: failed to add FlickerbladeEffect: %w", championName, err)
            }
            log.Printf("Added FlickerbladeEffect to champion %s (AS/stack: %.2f%%, AD/bonus: %.2f%%, AP/bonus: %.1f, Stacks/bonus: %.0f)",
                championName, asPerStackVal*100, adPerBonusVal*100, apPerBonusVal, stacksPerBonusVal)
        }
	case data.TFT_Item_NashorsTooth:
        if _, exists := em.world.GetNashorsToothEffect(champion); !exists {
            // Fetch values from item data
            attackSpeedToGive := item.Effects["AttackSpeedToGive"] / 100.0 // Convert percentage to decimal
            duration := item.Effects["ASDuration"]
            
            nashorsEffect := items.NewNashorsToothEffect(attackSpeedToGive, duration)
            err := em.world.AddComponent(champion, nashorsEffect)
            if err != nil {
                log.Printf("Warning: Failed to add NashorsToothEffect component for champion %s: %v", championName, err)
            } else {
                log.Printf("Added NashorsToothEffect component to champion %s (AS Gain: %.1f%%, Duration: %.1fs)",
                    championName, attackSpeedToGive*100, duration)
            }
        }
	case data.TFT_Item_VoidStaff:
		if _, exists := em.world.GetVoidStaffEffect(champion); !exists {
            // Fetch values from item data
            mrShred := item.Effects["MRShred"]         // 30.0
            duration := item.Effects["MRShredDuration"] // 3.0

            voidStaffEffect := items.NewVoidStaffEffect(mrShred, duration)
            err := em.world.AddComponent(champion, voidStaffEffect)
            if err != nil {
                log.Printf("Warning: Failed to add VoidStaffEffect component for champion %s: %v", championName, err)
            } else {
                log.Printf("Added VoidStaffEffect component to champion %s (MR Shred: %.1f%%, Duration: %.1fs)",
                    championName, mrShred, duration)
            }
        }
	case data.TFT_Item_RedBuff:
		if _, exists := em.world.GetRedBuffEffect(champion); !exists {
            // Fetch values from item data
            burnPercent := item.Effects["BurnPercent"]             // 1.0
            healingReductionPct := item.Effects["HealingReductionPct"] // 33.0
            duration := item.Effects["Duration"]                  // 5.0

            redBuffEffect := items.NewRedBuffEffect(burnPercent, healingReductionPct, duration)
            err := em.world.AddComponent(champion, redBuffEffect)
            if err != nil {
                log.Printf("Warning: Failed to add RedBuffEffect component for champion %s: %v", championName, err)
            } else {
                log.Printf("Added RedBuffEffect component to champion %s (Burn: %.1f%%, Wound: %.1f%%, Duration: %.1fs)",
                    championName, burnPercent, healingReductionPct, duration)
            }
        }
	}

    log.Printf("Updating static item effects for champion %s after adding %s.", championName, itemApiName)
    err = em.calculateAndUpdateStaticItemEffects(champion)
    if err != nil {
        // Attempt to remove the item if static effect calculation fails to revert state
        equipment.RemoveItem(itemApiName) // Best effort cleanup
        // Also potentially remove the specific effect component added above
        return fmt.Errorf("failed to calculate item effects for champion %s after adding %s: %w. Item addition reverted", championName, itemApiName, err)
    }

    return nil
}

// RemoveItemFromChampion removes an item from a champion's equipment by its API name.
// It also removes associated effect components for specific dynamic items.
// TODO: Need to handle all dynamic items that have effects
func (em *EquipmentManager) RemoveItemFromChampion(champion entity.Entity, itemApiName string) error {
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
		if _, exists := em.world.GetArchangelsStaffEffect(champion); exists {
			em.world.RemoveComponent(champion, reflect.TypeOf(items.ArchangelsStaffEffect{}))
			log.Printf("Removed ArchangelsEffect component from champion %s", championName)
		}
	case data.TFT_Item_Quicksilver:
		if _, exists := em.world.GetQuicksilverEffect(champion); exists {
			em.world.RemoveComponent(champion, reflect.TypeOf(items.QuicksilverEffect{}))
			log.Printf("Removed QuicksilverEffect component from champion %s", championName)
			// TODO: Remove IsImmuneToCC marker component if implemented
			// em.world.RemoveComponent(champion, reflect.TypeOf(effects.IsImmuneToCC{}))
		}
	case data.TFT_Item_TitansResolve: 
		if effect, exists := em.world.GetTitansResolveEffect(champion); exists {
			log.Printf("Removing Titan's Resolve effect from champion %s.", championName)
			// Get total bonuses provided by the effect *before* removing
			totalBonusAD := effect.GetCurrentBonusAD()
			totalBonusAP := effect.GetCurrentBonusAP()
			totalBonusArmor := effect.GetBonusArmorAtMax()
			totalBonusMR := effect.GetBonusMRAtMax()

			// Subtract these bonuses from the core components
			if attackComp, ok := em.world.GetAttack(champion); ok && totalBonusAD > 0 {
				attackComp.AddBonusPercentAD(-totalBonusAD)
				log.Printf("  Reversed Titan's AD: -%.2f%%. Total Bonus AD now: %.2f%%", totalBonusAD*100, attackComp.GetBonusPercentAD()*100)
			}
			if spellComp, ok := em.world.GetSpell(champion); ok && totalBonusAP > 0 {
				spellComp.AddBonusAP(-totalBonusAP)
				log.Printf("  Reversed Titan's AP: -%.1f. Total Bonus AP now: %.1f", totalBonusAP, spellComp.GetBonusAP())
			}
			if healthComp, ok := em.world.GetHealth(champion); ok {
				if totalBonusArmor > 0 {
					healthComp.AddBonusArmor(-totalBonusArmor)
				}
				if totalBonusMR > 0 {
					healthComp.AddBonusMR(-totalBonusMR)
				}
				log.Printf("  Reversed Titan's Resists: -%.0f Armor, -%.0f MR.", totalBonusArmor, totalBonusMR)
			}
		
			// Remove the effect component itself
			em.world.RemoveComponent(champion, reflect.TypeOf(items.TitansResolveEffect{}))
			log.Printf("Removed TitansResolveEffect component from champion %s", championName)
		}
	case data.TFT_Item_GuinsoosRageblade:
		if effect, exists := em.world.GetGuinsoosRagebladeEffect(champion); exists {
			log.Printf("Removing Guinsoo's Rageblade effect from champion %s.", championName)
			// Get total bonus AS provided by the effect *before* removing
			totalBonusAS := effect.GetCurrentBonusAS()

			// Subtract this bonus from the core Attack component
			if attackComp, ok := em.world.GetAttack(champion); ok && totalBonusAS > 0 {
				// Use the existing AddBonusPercentAttackSpeed with a negative value
				attackComp.AddBonusPercentAttackSpeed(-totalBonusAS)
				log.Printf("  Reversed Guinsoo's AS: -%.2f%%. Total Bonus AS now: %.2f%%", totalBonusAS*100, attackComp.GetBonusPercentAttackSpeed()*100)
			}
		
			// Remove the effect component itself
			em.world.RemoveComponent(champion, reflect.TypeOf(items.GuinsoosRagebladeEffect{}))
			log.Printf("Removed GuinsoosRagebladeEffect component from champion %s", championName)
		}
		// Add cases for other dynamic items
	}

	// --- Update Static Item Effects ---
	log.Printf("Updating static item effects for champion %s after removing %s.", championName, itemApiName)
	err := em.calculateAndUpdateStaticItemEffects(champion) // Recalculate remaining static passive stats
	if err != nil {
		log.Printf("Error updating static item effects for champion %s after removing %s: %v", championName, itemApiName, err)
		// return fmt.Errorf("failed to calculate item effects for champion %s: %w", championName, err) // Decide if this should be fatal
	}

	return nil
}

// calculateAndUpdateStaticItemEffects calculates the total passive stats from equipped items
// and updates the champion's ItemStaticEffect component.
func (em *EquipmentManager) calculateAndUpdateStaticItemEffects(champion entity.Entity) error {
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
		newItemEffect := items.NewItemStaticEffect()
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
				log.Printf("  [%s] Champion %d: Adding BonusHealth: %.1f", item.ApiName, champion, value)
			case "BonusPercentHP":
				itemEffect.AddBonusPercentHp(value)
				log.Printf("  [%s] Champion %d: Adding BonusPercentHP: %.1f%%", item.ApiName, champion, value*100)
			case "Mana":
				itemEffect.AddBonusInitialMana(value)
				log.Printf("  [%s] Champion %d: Adding Mana: %.1f", item.ApiName, champion, value)
			case "Armor":
				itemEffect.AddBonusArmor(value)
				log.Printf("  [%s] Champion %d: Adding Armor: %.1f", item.ApiName, champion, value)
			case "MagicResist":
				itemEffect.AddBonusMR(value)
				log.Printf("  [%s] Champion %d: Adding MagicResist: %.1f", item.ApiName, champion, value)
			case "AD":
				itemEffect.AddBonusPercentAD(value)
				log.Printf("  [%s] Champion %d: Adding BonusPercentAD: %.1f%%", item.ApiName, champion, value*100)
			case "AP":
				itemEffect.AddBonusAP(value)
				log.Printf("  [%s] Champion %d: Adding AP: %.1f", item.ApiName, champion, value)
			case "AS":
				itemEffect.AddBonusPercentAttackSpeed(value / 100)
				log.Printf("  [%s] Champion %d: Adding BonusPercentAttackSpeed: %.1f%%", item.ApiName, champion, value) // Log original percentage value
			case "CritChance":
				itemEffect.AddBonusCritChance(value / 100)
				log.Printf("  [%s] Champion %d: Adding CritChance: %.1f%%", item.ApiName, champion, value) // Log original percentage value
			case "BonusDamage":
				itemEffect.AddBonusDamageAmp(value)
				log.Printf("  [%s] Champion %d: Adding BonusDamageAmp: %.1f%%", item.ApiName, champion, value*100)
			case "CritDamageToGive": // specific to IE & JG
				itemEffect.AddCritDamageToGive(value)
				log.Printf("  [%s] Champion %d: Adding CritDamageToGive: %.1f%%", item.ApiName, champion, value*100)
            case "Durability": // NOTE: New case for Spirit Visage and other durability items
                itemEffect.AddDurability(value) // value is 0.10 for 10%
                log.Printf("  [%s] Champion %d: Adding Durability: %.1f%%", item.ApiName, champion, value*100)
			// Add other known static stats...
			default:
				log.Printf("Warning: Champion %d: Unrecognized or non-static item effect stat '%s' (value: %.2f) for item %s", champion, statName, value, item.ApiName)

			}
		}
	}

	log.Printf("Finished calculating static item effects for champion %s.", championName)
	return nil
}
