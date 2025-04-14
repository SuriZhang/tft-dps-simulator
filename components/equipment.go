package components

import (
	"fmt"

	"github.com/suriz/tft-dps-simulator/data"
)

const MaxItems = 3

// Equipment holds the items equipped by an entity.
type Equipment struct {
    // Use a map where the key is the item's ApiName and the value is the item data.
    Items    map[string]*data.Item
    MaxSlots int
}

// NewEquipment creates a new Equipment component.
func NewEquipment() *Equipment {
    return &Equipment{
        Items:    make(map[string]*data.Item, MaxItems),
        MaxSlots: MaxItems,
    }
}

// HasItemSlots checks if there's space for another item.
func (eq *Equipment) HasItemSlots() bool {
    return len(eq.Items) < eq.MaxSlots
}

// AddItem adds an item to the equipment map.
// Assumes checks for space and uniqueness have already been done.
// Returns true if added, false otherwise (e.g., if item is nil).
func (eq *Equipment) AddItem(item *data.Item) error {
    if item == nil {
        return fmt.Errorf("cannot add a nil item") 
    }

    eq.Items[item.ApiName] = item
    return nil
}

// HasItem checks if an item with the given ApiName is already equipped.
func (eq *Equipment) HasItem(itemApiName string) bool {
    _, exists := eq.Items[itemApiName]
    return exists
}

// IsDuplicateUniqueItem checks if the user is attempting to add a unique item
// that is already present.
// Note: This logic is slightly simplified because HasItem now exists.
func (eq *Equipment) IsDuplicateUniqueItem(itemApiName string) bool {
    // Check if the item to be added is unique
    itemToAdd := data.GetItemByApiName(itemApiName)
    if itemToAdd == nil || !itemToAdd.Unique {
        return false // Item not found or not unique, so no duplicate *unique* issue
    }

    // Check if an item with the same ApiName already exists in the map
    return eq.HasItem(itemApiName)
}

// RemoveItem removes an item by its API name. Returns true if successful, false if not found.
func (eq *Equipment) RemoveItem(itemApiName string) bool {
    if !eq.HasItem(itemApiName) {
        return false // Item not found
    }
    delete(eq.Items, itemApiName) // Remove the item from the map
    return true
}

// GetItem retrieves an item by its API name. Returns the item and true if found, nil and false otherwise.
func (eq *Equipment) GetItem(itemApiName string) (*data.Item, bool) {
    item, exists := eq.Items[itemApiName]
    return item, exists
}

// GetAllItems returns a slice containing all equipped items.
// Useful for systems that need to iterate over all items.
func (eq *Equipment) GetAllItems() []*data.Item {
    itemsSlice := make([]*data.Item, 0, len(eq.Items))
    for _, item := range eq.Items {
        itemsSlice = append(itemsSlice, item)
    }
    return itemsSlice
}

// IsFull checks if the equipment slots are full.
func (eq *Equipment) IsFull() bool {
    return len(eq.Items) >= eq.MaxSlots
}