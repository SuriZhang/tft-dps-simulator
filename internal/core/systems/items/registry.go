package itemsys

import (
	"log"
    
)

// ItemRegistry maps item API names to their corresponding ItemHandler implementations.
var ItemRegistry = make(map[string]ItemHandler)

// RegisterItemHandler registers a handler for a specific item API name.
func RegisterItemHandler(itemApiName string, handler ItemHandler) {
    if _, exists := ItemRegistry[itemApiName]; exists {
        log.Printf("Warning: Overwriting existing item handler for %s", itemApiName)
    }
    ItemRegistry[itemApiName] = handler
    log.Printf("Registered item handler for %s", itemApiName)
}

// GetItemHandler retrieves the handler for a specific item API name.
func GetItemHandler(itemApiName string) (ItemHandler, bool) {
    handler, exists := ItemRegistry[itemApiName]
    return handler, exists
}