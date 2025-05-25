package itemhandlers

import (
	"log"
	"tft-dps-simulator/internal/core/data"
	"tft-dps-simulator/internal/core/ecs"
	"tft-dps-simulator/internal/core/entity"
	eventsys "tft-dps-simulator/internal/core/systems/events"
	itemsys "tft-dps-simulator/internal/core/systems/items"
)

type SpearOfShojinHandler struct{}

func init() {
    itemsys.RegisterItemHandler(data.TFT_Item_SpearOfShojin, &SpearOfShojinHandler{})
}

func (h *SpearOfShojinHandler) OnEquip(entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
    log.Printf("SpearOfShojinHandler: OnEquip for entity %d", entity)
    // No special equip logic needed - effect is handled in ProcessEvent
}

// func (h *SpearOfShojinHandler) OnUnequip(entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
//     log.Printf("SpearOfShojinHandler: OnUnequip for entity %d", entity)
//     // Remove the component when item is unequipped
//     world.RemoveComponent(entity, (*items.SpearOfShojinEffect)(nil))
// }

func (h *SpearOfShojinHandler) ProcessEvent(event interface{}, entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
    switch evt := event.(type) {
    case eventsys.AttackFiredEvent:
        h.handleAttackFired(evt, entity, world, eventBus)
    }
}

func (h *SpearOfShojinHandler) handleAttackFired(evt eventsys.AttackFiredEvent, entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
    // Only process if this entity is the attacker
    if evt.Source != entity {
        return
    }

    effect, exists := world.GetSpearOfShojinEffect(entity)
    if !exists {
        log.Printf("SpearOfShojinHandler: SpearOfShojinEffect component not found for entity %d", entity)
        return
    }

    equipment, equipExists := world.GetEquipment(entity)
    if !equipExists {
        log.Printf("SpearOfShojinHandler: Equipment component not found for entity %d", entity)
        return
    }

    mana, manaExists := world.GetMana(entity)
    if !manaExists {
        log.Printf("SpearOfShojinHandler: Mana component not found for entity %d", entity)
        return
    }

    // Apply mana restoration
    spearCount := equipment.GetItemCount(data.TFT_Item_SpearOfShojin)
    manaRestore := effect.GetFlatManaRestore() * float64(spearCount)
    
    initialMana := mana.GetCurrentMana()
    mana.AddCurrentMana(manaRestore)
    
    log.Printf("SpearOfShojinHandler (Attack): Entity %d restored %.1f mana on attack at %.3fs (%.1f -> %.1f / %.1f)",
        entity, manaRestore, evt.Timestamp, initialMana, mana.GetCurrentMana(), mana.GetMaxMana())
}