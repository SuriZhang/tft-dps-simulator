package itemhandlers

import (
	"log"
	"math"
	"reflect"

	"tft-dps-simulator/internal/core/components"
	"tft-dps-simulator/internal/core/components/debuffs"
	"tft-dps-simulator/internal/core/data"
	"tft-dps-simulator/internal/core/ecs"
	"tft-dps-simulator/internal/core/entity"
	eventsys "tft-dps-simulator/internal/core/systems/events"
	itemsys "tft-dps-simulator/internal/core/systems/items"
)

type EvenshroudHandler struct{}

// Register the handler
func init() {
    itemsys.RegisterItemHandler(data.TFT_Item_Evenshroud, &EvenshroudHandler{})
}

// OnEquip implements itemsys.ItemHandler.
func (h *EvenshroudHandler) OnEquip(entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
    log.Printf("EvenshroudHandler: OnEquip for entity %d", entity)

	// Apply sunder aura immediately when equipped
    h.applySunderAura(entity, world, eventBus, 0.0)
    
    // Schedule the resistance bonus activation at combat start
    activateEvent := eventsys.EvenshroudResistActivateEvent{
        Entity:    entity,
        Timestamp: 0.0, // Combat start
    }
    eventBus.Enqueue(activateEvent, 0.0)
}

// OnUnequip implements itemsys.ItemHandler.
func (h *EvenshroudHandler) OnUnequip(entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
    log.Printf("EvenshroudHandler: OnUnequip for entity %d", entity)
    
    // Deactivate resistance bonus if active
    if effect, ok := world.GetEvenshroudEffect(entity); ok && effect.IsResistBonusActive() {
        deactivateEvent := eventsys.EvenshroudResistDeactivateEvent{
            Entity:    entity,
            Timestamp: 0.0, // Immediate
        }
        eventBus.Enqueue(deactivateEvent, 0.0)
    }
}

// ProcessEvent implements itemsys.ItemHandler.
func (h *EvenshroudHandler) ProcessEvent(event interface{}, entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
    switch evt := event.(type) {
    case eventsys.EvenshroudResistActivateEvent:
        if evt.Entity == entity {
            h.handleResistActivate(evt, entity, world, eventBus)
        }
    case eventsys.EvenshroudResistDeactivateEvent:
        if evt.Entity == entity {
            h.handleResistDeactivate(evt, entity, world, eventBus)
        }
    }
}

// handleResistActivate activates the temporary resistance bonus
func (h *EvenshroudHandler) handleResistActivate(evt eventsys.EvenshroudResistActivateEvent, entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
    // Verify entity still has Evenshroud
    equipment, okEq := world.GetEquipment(entity)
    if !okEq || equipment.GetItemCount(data.TFT_Item_Evenshroud) == 0 {
        log.Printf("EvenshroudHandler (handleResistActivate): Entity %d no longer has Evenshroud at %.3fs.", entity, evt.Timestamp)
        return
    }

    effect, okE := world.GetEvenshroudEffect(entity)
    if !okE {
        log.Printf("EvenshroudHandler (handleResistActivate): Entity %d missing EvenshroudEffect component at %.3fs.", entity, evt.Timestamp)
        return
    }

    health, okHealth := world.GetHealth(entity)
    if !okHealth {
        log.Printf("EvenshroudHandler (handleResistActivate): Entity %d missing Health component at %.3fs.", entity, evt.Timestamp)
        return
    }

    // Activate the resistance bonus
    effect.ActivateResistBonus()

    // Add temporary armor and magic resist bonuses
    bonusResists := effect.GetBonusResists()
    health.AddBonusArmor(bonusResists)
    health.AddBonusMR(bonusResists)

    log.Printf("EvenshroudHandler: Entity %d gained +%.0f Armor and MR for %.1fs at %.3fs",
        entity, bonusResists, effect.GetBonusResistDuration(), evt.Timestamp)

    // Enqueue event to recalculate stats
    recalcEvent := eventsys.RecalculateStatsEvent{Entity: entity, Timestamp: evt.Timestamp}
    eventBus.Enqueue(recalcEvent, evt.Timestamp)

    // Schedule deactivation
    deactivateEvent := eventsys.EvenshroudResistDeactivateEvent{
        Entity:    entity,
        Timestamp: evt.Timestamp + effect.GetBonusResistDuration(),
    }
    eventBus.Enqueue(deactivateEvent, evt.Timestamp + effect.GetBonusResistDuration())
}

// handleResistDeactivate deactivates the temporary resistance bonus
func (h *EvenshroudHandler) handleResistDeactivate(evt eventsys.EvenshroudResistDeactivateEvent, entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus) {
    effect, okE := world.GetEvenshroudEffect(entity)
    if !okE {
        return
    }

    health, okHealth := world.GetHealth(entity)
    if !okHealth {
        log.Printf("EvenshroudHandler (handleResistDeactivate): Entity %d missing Health component at %.3fs.", entity, evt.Timestamp)
        return
    }

    // Deactivate the resistance bonus
    effect.DeactivateResistBonus()

    // Remove temporary armor and magic resist bonuses
    bonusResists := effect.GetBonusResists()
    health.AddBonusArmor(-bonusResists) // Subtract the bonus
    health.AddBonusMR(-bonusResists)    // Subtract the bonus

    log.Printf("EvenshroudHandler: Entity %d lost +%.0f Armor and MR temporary resistance bonus at %.3fs", 
        entity, bonusResists, evt.Timestamp)

    // Enqueue event to recalculate stats
    recalcEvent := eventsys.RecalculateStatsEvent{Entity: entity, Timestamp: evt.Timestamp}
    eventBus.Enqueue(recalcEvent, evt.Timestamp)
}

// applySunderAura applies sunder debuff to enemies within range
func (h *EvenshroudHandler) applySunderAura(entity entity.Entity, world *ecs.World, eventBus eventsys.EventBus, timestamp float64) {
    // Verify entity still has Evenshroud
    equipment, okEq := world.GetEquipment(entity)
    if !okEq || equipment.GetItemCount(data.TFT_Item_Evenshroud) == 0 {
        return
    }

    effect, okE := world.GetEvenshroudEffect(entity)
    if !okE {
        return
    }

    // Get holder's position
    holderPos, okPos := world.GetPosition(entity)
    if !okPos {
        return
    }

    // Get holder's team
    holderTeam, okTeam := world.GetTeam(entity)
    if !okTeam {
        return
    }

    hexRange := effect.GetHexRange()
    arReduction := effect.GetARReductionAmount()

    // Find all enemies within range
    enemies := h.getEnemiesInRange(entity, holderPos, holderTeam, hexRange, world)

    // Apply sunder to each enemy
    for _, enemyEntity := range enemies {
        sunderEvent := eventsys.ApplyDebuffEvent{
            Target:     enemyEntity,
            Source:     entity,
            DebuffType: debuffs.Sunder,
            Value:      arReduction,
            Duration:   999.0, // Permanent for the duration of combat
            Timestamp:  timestamp,
            SourceType: "Item",
            SourceId:   data.TFT_Item_Evenshroud,
        }
        eventBus.Enqueue(sunderEvent, timestamp)
    }

    if len(enemies) > 0 {
        log.Printf("EvenshroudHandler: Entity %d applied %.1f%% sunder to %d enemies within %.1f hexes at %.3fs",
            entity, arReduction*100, len(enemies), hexRange, timestamp)
    }
}

// getEnemiesInRange finds all enemy entities within the specified range
func (h *EvenshroudHandler) getEnemiesInRange(sourceEntity entity.Entity, sourcePos *components.Position, sourceTeam *components.Team, hexRange float64, world *ecs.World) []entity.Entity {
    var enemies []entity.Entity

    // Get all entities with position and team components
    positionEntities := world.GetEntitiesWithComponents(reflect.TypeOf(components.Position{}), reflect.TypeOf(components.Team{}))

    for _, targetEntity := range positionEntities {
        if targetEntity == sourceEntity {
            continue
        }

        targetPos, okPos := world.GetPosition(targetEntity)
        if !okPos {
            continue
        }

        targetTeam, okTeam := world.GetTeam(targetEntity)
        if !okTeam {
            continue
        }

        // Check if it's an enemy
        if targetTeam.ID == sourceTeam.ID {
            continue
        }

        // Calculate distance
        dx := targetPos.GetX() - sourcePos.GetX()
        dy := targetPos.GetY() - sourcePos.GetY()
        distance := math.Sqrt(float64(dx*dx + dy*dy))

        // Check if within range
        if distance <= hexRange {
            enemies = append(enemies, targetEntity)
        }
    }

    return enemies
}
