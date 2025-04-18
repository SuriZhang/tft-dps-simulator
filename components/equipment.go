package components

import (
	"fmt"

	"github.com/suriz/tft-dps-simulator/data"
)

const MaxItems = 3

// Equipment holds the items equipped by an entity using a slice.
type Equipment struct {
	Items    []*data.Item // Slice to hold items, allows duplicates
	MaxSlots int
}

// NewEquipment creates a new Equipment component.
func NewEquipment() *Equipment {
	return &Equipment{
		// Initialize slice with 0 length but capacity for MaxItems
		Items:    make([]*data.Item, 0, MaxItems),
		MaxSlots: MaxItems,
	}
}

// HasItemSlots checks if there's space for another item.
func (eq *Equipment) HasItemSlots() bool {
	return len(eq.Items) < eq.MaxSlots
}

// AddItem adds an item to the equipment slice.
// Assumes checks for space and uniqueness (if applicable) have already been done.
// Returns an error if the item is nil.
func (eq *Equipment) AddItem(item *data.Item) error {
	if item == nil {
		return fmt.Errorf("cannot add a nil item")
	}

	eq.Items = append(eq.Items, item)
	return nil
}

// HasItem checks if an item with the given ApiName is already equipped.
func (eq *Equipment) HasItem(itemApiName string) bool {
	for _, equippedItem := range eq.Items {
		if equippedItem != nil && equippedItem.ApiName == itemApiName {
			return true
		}
	}
	return false
}

// IsDuplicateUniqueItem checks if the user is attempting to add a unique item
// that is already present.
func (eq *Equipment) IsDuplicateUniqueItem(itemApiName string) bool {
	itemToAdd := data.GetItemByApiName(itemApiName) // Check the item definition
	if itemToAdd == nil || !itemToAdd.Unique {
		return false // Item not found or not unique, so no duplicate *unique* issue
	}

	// Check if an item with the same ApiName already exists in the slice
	return eq.HasItem(itemApiName)
}

// RemoveItem removes the *first* occurrence of an item by its API name.
// Returns true if successful, false if not found.
func (eq *Equipment) RemoveItem(itemApiName string) bool {
	removeIndex := -1
	for i, equippedItem := range eq.Items {
		if equippedItem != nil && equippedItem.ApiName == itemApiName {
			removeIndex = i
			break // Found the first occurrence
		}
	}

	if removeIndex == -1 {
		return false // Item not found
	}

	// Remove element by slicing (preserves order if needed, efficient for small slices)
	// See: https://github.com/golang/go/wiki/SliceTricks#delete
	eq.Items = append(eq.Items[:removeIndex], eq.Items[removeIndex+1:]...)
	return true
}

// GetItem retrieves the *first* occurrence of an item by its API name.
// Returns the item and true if found, nil and false otherwise.
func (eq *Equipment) GetItem(itemApiName string) (*data.Item, bool) {
	for _, equippedItem := range eq.Items {
		if equippedItem != nil && equippedItem.ApiName == itemApiName {
			return equippedItem, true // Return first match
		}
	}
	return nil, false
}

// GetAllItems returns a slice containing all equipped items.
// Returns a copy to prevent external modification of the internal slice.
func (eq *Equipment) GetAllItems() []*data.Item {
	// Return a copy to be safe
	itemsCopy := make([]*data.Item, len(eq.Items))
	copy(itemsCopy, eq.Items)
	return itemsCopy
}

// GetItemCount returns the number of instances of an item with the given ApiName.
func (eq *Equipment) GetItemCount(itemApiName string) int {
	count := 0
	for _, equippedItem := range eq.Items {
		if equippedItem != nil && equippedItem.ApiName == itemApiName {
			count++
		}
	}
	return count
}
