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
func (inv *Equipment) CanAddItem() bool {
    return len(inv.Items) < inv.MaxSlots
}

// AddItem adds an item if there is space. Returns true if successful, false otherwise.
func (inv *Equipment) AddItem(item *data.Item) bool {
    if inv.CanAddItem() {
        inv.Items = append(inv.Items, item)
        return true
    }
    return false
}

// RemoveItem removes an item by its API name. Returns true if successful, false if not found.
func (inv *Equipment) RemoveItem(itemApiName string) bool { 
	for i, item := range inv.Items {
		if item.ApiName == itemApiName {
			inv.Items = append(inv.Items[:i], inv.Items[i+1:]...)
			return true
		}
	}
	return false
}