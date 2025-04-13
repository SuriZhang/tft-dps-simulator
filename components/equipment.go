package components

import "github.com/suriz/tft-dps-simulator/data" 

const MaxItems = 3

// Equipment holds the items equipped by an entity.
type Equipment struct {
    Items    []*data.Item // Slice to store pointers to the actual item data
    MaxSlots int
}

func NewEquipment() *Equipment {
    return &Equipment{
        Items:    make([]*data.Item, 0, MaxItems), // Pre-allocate capacity
        MaxSlots: MaxItems,
    }
}

// CanAddItem checks if there's space for another item.
func (eq *Equipment) CanAddItem() bool {
    return len(eq.Items) < eq.MaxSlots
}

// HasItemSlots adds an item if there is space. Returns true if successful, false otherwise.
func (eq *Equipment) HasItemSlots(item *data.Item) bool {
    if eq.CanAddItem() {
        eq.Items = append(eq.Items, item)
        return true
    } 
    return false
}

// IsUniqueItem checks if user is attempting to add another unique item of the same type.
func (eq *Equipment) IsDuplicateUniqueItem(itemApiName string) bool {
    // Check if the item is unique
    item := data.GetItemByApiName(itemApiName)
    if item == nil {
        return false // Item not found
    }
    if item.Unique {
        for _, equippedItem := range eq.Items {
            if equippedItem.ApiName == itemApiName {
                return true // Duplicate unique item found
            }
        }
    }
    return false // No duplicates found
}

// RemoveItem removes an item by its API name. Returns true if successful, false if not found.
func (eq *Equipment) RemoveItem(itemApiName string) bool { 
	for i, item := range eq.Items {
		if item.ApiName == itemApiName {
			eq.Items = append(eq.Items[:i], eq.Items[i+1:]...)
			return true
		}
	}
	return false
}