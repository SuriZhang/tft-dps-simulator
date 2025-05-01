package traitsys

import (
    "log"
)

// TraitRegistry maps trait API names to their corresponding TraitHandler implementations.
var TraitRegistry = make(map[string]TraitHandler)

// RegisterTraitHandler registers a handler for a specific trait API name.
func RegisterTraitHandler(traitName string, handler TraitHandler) {
    if _, exists := TraitRegistry[traitName]; exists {
        log.Printf("Warning: Overwriting existing trait handler for %s", traitName)
    }
    TraitRegistry[traitName] = handler
    log.Printf("Registered trait handler for %s", traitName)
}

// GetTraitHandler retrieves the handler for a specific trait API name.
func GetTraitHandler(traitApiName string) (TraitHandler, bool) {
    handler, exists := TraitRegistry[traitApiName]
    return handler, exists
}
