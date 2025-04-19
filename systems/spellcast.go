package systems

import (
    "github.com/suriz/tft-dps-simulator/ecs"
    eventsys "github.com/suriz/tft-dps-simulator/systems/events"
)

type SpellCastSystem struct {
    world    *ecs.World
    eventBus eventsys.EventBus
}

func NewSpellCastSystem(world *ecs.World, bus eventsys.EventBus) *SpellCastSystem {
    return &SpellCastSystem{world: world, eventBus: bus}
}

