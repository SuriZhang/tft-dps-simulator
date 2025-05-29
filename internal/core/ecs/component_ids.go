package ecs

// ComponentID is a type alias for identifying component types
type ComponentID string

var nextComponentID ComponentID
var componentRegistry = make(map[string]ComponentID)

func registerComponent(name string) ComponentID {
	if id, exists := componentRegistry[name]; exists {
		return id
	}
	id := ComponentID(name)
	componentRegistry[name] = id
	return id
}

// Existing components (assuming some might exist, add more as needed)
var (
	HealthComponentID   ComponentID = registerComponent("Health")
	AttackComponentID   ComponentID = registerComponent("Attack")
	StatsComponentID    ComponentID = registerComponent("Stats")
	TargetComponentID   ComponentID = registerComponent("Target")
	DebuffComponentID   ComponentID = registerComponent("Debuff") // Generic debuff component

	// New Debuff Component IDs
	ShredEffectComponentID  ComponentID = registerComponent("ShredEffect")
	SunderEffectComponentID ComponentID = registerComponent("SunderEffect")
	WoundEffectComponentID  ComponentID = registerComponent("WoundEffect")
	BurnEffectComponentID   ComponentID = registerComponent("BurnEffect")
)
