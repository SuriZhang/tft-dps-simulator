package systems

import (
	"log"
	"reflect"

	"tft-dps-simulator/internal/core/components/debuffs"
	"tft-dps-simulator/internal/core/ecs"
	"tft-dps-simulator/internal/core/entity"
	eventsys "tft-dps-simulator/internal/core/systems/events"
)

type DebuffSystem struct {
    world    *ecs.World
    eventBus eventsys.EventBus
}

func NewDebuffSystem(world *ecs.World, eventBus eventsys.EventBus) *DebuffSystem {
    return &DebuffSystem{
        world:    world,
        eventBus: eventBus,
    }
}

func (s *DebuffSystem) CanHandle(evt interface{}) bool {
    switch evt.(type) {
    case eventsys.ApplyDebuffEvent, eventsys.RemoveDebuffEvent, eventsys.BurnTickEvent, eventsys.DebuffExpiredEvent:
        return true
    default:
        return false
    }
}

func (s *DebuffSystem) HandleEvent(evt interface{}) {
    switch event := evt.(type) {
    case eventsys.ApplyDebuffEvent:
        s.applyDebuff(event)
    case eventsys.RemoveDebuffEvent:
        s.removeDebuff(event)
    case eventsys.BurnTickEvent:
        s.handleBurnTick(event)
    case eventsys.DebuffExpiredEvent:
        s.handleDebuffExpiration(event)
    }
}

func (s *DebuffSystem) applyDebuff(evt eventsys.ApplyDebuffEvent) {
    target := evt.Target
    endTime := evt.Timestamp + evt.Duration

    switch evt.DebuffType {
    case debuffs.Shred:
        s.applyShred(target, evt.Value, evt.Duration, endTime, evt.Source, evt.SourceType, evt.SourceId, evt.Timestamp)
    case debuffs.Sunder:
        s.applySunder(target, evt.Value, evt.Duration, endTime, evt.Source, evt.SourceType, evt.SourceId, evt.Timestamp)
    case debuffs.Wound:
        s.applyWound(target, evt.Value, evt.Duration, endTime, evt.Source, evt.SourceType, evt.SourceId, evt.Timestamp)
    case debuffs.Burn:
        s.applyBurn(target, evt.Value, evt.Duration, endTime, evt.Source, evt.SourceType, evt.SourceId, evt.Timestamp)
    default:
        log.Printf("DebuffSystem: Unknown debuff type %s", evt.DebuffType)
    }
}

func (s *DebuffSystem) applyShred(target entity.Entity, mrReduction, duration, endTime float64, source entity.Entity, sourceType, sourceId string, timestamp float64) {
	mrToReduce := mrReduction
    if existingShred, exists := s.world.GetShredEffect(target); exists {
        existingShred.UpdateFromStrongerEffect(mrReduction, duration, endTime, source, sourceType, sourceId)
        log.Printf("DebuffSystem: Updated Shred on entity %d (%.1f MR reduction, %.1fs duration) from %s", 
            target, existingShred.GetMRReduction(), existingShred.GetDuration(), sourceId)
		mrToReduce = existingShred.GetMRReduction()
    } else {
        shredEffect := debuffs.NewShredEffect(mrReduction, duration, endTime, source, sourceType, sourceId)
        s.world.AddComponent(target, shredEffect)
        log.Printf("DebuffSystem: Applied Shred to entity %d (%.1f MR reduction, %.1fs duration) from %s", 
            target, mrReduction, duration, sourceId)
    }

	healthComp, ok := s.world.GetHealth(target)
	if !ok {
		log.Printf("DebuffSystem: Target entity %d has no Health component, cannot apply Shred", target)
		return
	}
	
	healthComp.AddBonusMR(-mrToReduce)
	s.enqueueStatsRecalculation(target, timestamp)

    expireEvent := eventsys.DebuffExpiredEvent{
        Target:     target,
        DebuffType: debuffs.Shred,
        Timestamp:  endTime,
        SourceId:   sourceId,
    }
    s.eventBus.Enqueue(expireEvent, endTime)
	log.Printf("DebuffSystem: Enqueued Shred expiration event for entity %d at t=%.3fs", target, endTime)
}

func (s *DebuffSystem) applySunder(target entity.Entity, armorReduction, duration, endTime float64, source entity.Entity, sourceType, sourceId string, timestamp float64) {
	armorToReduce := armorReduction
    if existingSunder, exists := s.world.GetSunderEffect(target); exists {
        existingSunder.UpdateFromStrongerEffect(armorReduction, duration, endTime, source, sourceType, sourceId)
        log.Printf("DebuffSystem: Updated Sunder on entity %d (%.1f armor reduction, %.1fs duration) from %s", 
            target, existingSunder.GetArmorReduction(), existingSunder.GetDuration(), sourceId)
		armorToReduce = existingSunder.GetArmorReduction()
    } else {
        sunderEffect := debuffs.NewSunderEffect(armorReduction, duration, endTime, source, sourceType, sourceId)
        s.world.AddComponent(target, sunderEffect)
        log.Printf("DebuffSystem: Applied Sunder to entity %d (%.1f armor reduction, %.1fs duration) from %s", 
            target, armorReduction, duration, sourceId)
    }

	healthComp, ok := s.world.GetHealth(target)
	if !ok {
		log.Printf("DebuffSystem: Target entity %d has no Health component, cannot apply Sunder", target)
		return
	}
	healthComp.AddBonusArmor(-armorToReduce)
	s.enqueueStatsRecalculation(target, timestamp)

    expireEvent := eventsys.DebuffExpiredEvent{
        Target:     target,
        DebuffType: debuffs.Sunder,
        Timestamp:  endTime,
        SourceId:   sourceId,
    }
    s.eventBus.Enqueue(expireEvent, endTime)
    s.enqueueStatsRecalculation(target, timestamp)
}

func (s *DebuffSystem) applyWound(target entity.Entity, healingReduction, duration, endTime float64, source entity.Entity, sourceType, sourceId string, timestamp float64) {
	healingReductionToApply := healingReduction
    if existingWound, exists := s.world.GetWoundEffect(target); exists {
        existingWound.UpdateFromStrongerEffect(healingReduction, duration, endTime, source, sourceType, sourceId)
        log.Printf("DebuffSystem: Updated Wound on entity %d (%.1f%% healing reduction, %.1fs duration) from %s", 
            target, existingWound.GetHealingReduction()*100, existingWound.GetDuration(), sourceId)
		healingReductionToApply = existingWound.GetHealingReduction()
    } else {
        woundEffect := debuffs.NewWoundEffect(healingReduction, duration, endTime, source, sourceType, sourceId)
        s.world.AddComponent(target, woundEffect)
        log.Printf("DebuffSystem: Applied Wound to entity %d (%.1f%% healing reduction, %.1fs duration) from %s", 
            target, healingReduction*100, duration, sourceId)
    }

	healthComp, ok := s.world.GetHealth(target)
	if !ok {
		log.Printf("DebuffSystem: Target entity %d has no Health component, cannot apply Wound", target)
		return
	}
	healthComp.SetHealReduction(healingReductionToApply)
	s.enqueueStatsRecalculation(target, timestamp)

    expireEvent := eventsys.DebuffExpiredEvent{
        Target:     target,
        DebuffType: debuffs.Wound,
        Timestamp:  endTime,
        SourceId:   sourceId,
    }
    s.eventBus.Enqueue(expireEvent, endTime)
}

func (s *DebuffSystem) applyBurn(target entity.Entity, damagePercent, duration, endTime float64, source entity.Entity, sourceType, sourceId string, timestamp float64) {
    isNewBurn := false
    
    if existingBurn, exists := s.world.GetBurnEffect(target); exists {
        existingBurn.UpdateFromStrongerEffect(damagePercent, duration, endTime, source, sourceType, sourceId)
        log.Printf("DebuffSystem: Updated Burn on entity %d (%.1f%% max HP per second, %.1fs duration) from %s", 
            target, existingBurn.GetDamagePercent()*100, existingBurn.GetDuration(), sourceId)
    } else {
        burnEffect := debuffs.NewBurnEffect(damagePercent, duration, endTime, source, sourceType, sourceId)
        s.world.AddComponent(target, burnEffect)
        isNewBurn = true
        log.Printf("DebuffSystem: Applied Burn to entity %d (%.1f%% max HP per second, %.1fs duration) from %s", 
            target, damagePercent*100, duration, sourceId)
    }

    expireEvent := eventsys.DebuffExpiredEvent{
        Target:     target,
        DebuffType: debuffs.Burn,
        Timestamp:  endTime,
        SourceId:   sourceId,
    }
    s.eventBus.Enqueue(expireEvent, endTime)

    // Enqueue the first burn tick event (like Archangel's Staff pattern)
    if isNewBurn {
        firstBurnTickTime := timestamp + 1.0 // First tick after 1 second
        burnTickEvent := eventsys.BurnTickEvent{
            Target:    target,
            Source:    source,
            Timestamp: firstBurnTickTime,
            SourceId:  sourceId,
        }
        s.eventBus.Enqueue(burnTickEvent, firstBurnTickTime)
        log.Printf("DebuffSystem: Enqueued initial BurnTickEvent for entity %d at t=%.3fs", target, firstBurnTickTime)
    }
}

func (s *DebuffSystem) handleBurnTick(evt eventsys.BurnTickEvent) {
    target := evt.Target
    currentTime := evt.Timestamp

    log.Printf("DebuffSystem: Processing BurnTickEvent for entity %d at t=%.3fs", target, currentTime)

    // Check if the burn effect still exists and is active
    burnEffect, exists := s.world.GetBurnEffect(target)
    if !exists {
        log.Printf("DebuffSystem: BurnTickEvent for entity %d - burn effect no longer exists at %.3fs. Stopping ticks.", target, currentTime)
        return
    }

    if !burnEffect.IsActive(currentTime) {
        log.Printf("DebuffSystem: BurnTickEvent for entity %d - burn effect expired at %.3fs. Stopping ticks.", target, currentTime)
        return
    }

    // // Check if the burn tick is from the same source
    // if evt.SourceId != burnEffect.GetSourceId() {
    //     log.Printf("DebuffSystem: BurnTickEvent for entity %d - source mismatch (%s vs %s). Stopping ticks.", target, evt.SourceId, burnEffect.GetSourceId())
    //     return
    // }

    // Check if entity is still alive
    health, healthExists := s.world.GetHealth(target)
    if !healthExists || health.GetCurrentHP() <= 0 {
        log.Printf("DebuffSystem: BurnTickEvent for entity %d - entity dead at %.3fs. Stopping ticks.", target, currentTime)
        return
    }

    // Apply burn damage
    maxHP := health.GetFinalMaxHP()
    burnDamage := maxHP * burnEffect.GetDamagePercent()

    burnDamageEvent := eventsys.DamageAppliedEvent{
        Source:              burnEffect.GetSourceEntity(),
        Target:              target,
        Timestamp:           currentTime,
        DamageType:          "True",
        DamageSource:        "Burn",
        RawDamage:           burnDamage,
        PreMitigationDamage: burnDamage,
        FinalTotalDamage:    burnDamage,
        MitigatedDamage:     0.0,
    }

    s.eventBus.Enqueue(burnDamageEvent, currentTime)
    burnEffect.SetLastTickTime(currentTime)

    log.Printf("DebuffSystem: Burn tick on entity %d for %.1f true damage (%.1f%% of %.1f max HP)", 
        target, burnDamage, burnEffect.GetDamagePercent()*100, maxHP)

    // Enqueue the next burn tick event
    nextTickTime := currentTime + burnEffect.GetTickInterval()
    if nextTickTime < burnEffect.GetEndTime() {
        nextBurnTickEvent := eventsys.BurnTickEvent{
            Target:    target,
            Source:    burnEffect.GetSourceEntity(),
            Timestamp: nextTickTime,
            SourceId:  burnEffect.GetSourceId(),
        }
        s.eventBus.Enqueue(nextBurnTickEvent, nextTickTime)
        log.Printf("DebuffSystem: Enqueued next BurnTickEvent for entity %d at t=%.3fs", target, nextTickTime)
    } else {
        log.Printf("DebuffSystem: No more burn ticks for entity %d - effect will expire at %.3fs", target, burnEffect.GetEndTime())
    }
}

func (s *DebuffSystem) handleDebuffExpiration(evt eventsys.DebuffExpiredEvent) {
    switch evt.DebuffType {
    case debuffs.Shred:
        if shred, exists := s.world.GetShredEffect(evt.Target); exists && shred.GetSourceId() == evt.SourceId {
            s.world.RemoveComponent(evt.Target, reflect.TypeOf(debuffs.ShredEffect{}))
			// Remove the MR reduction from health component
			healthComp, ok := s.world.GetHealth(evt.Target)
			if ok {
				healthComp.AddBonusMR(shred.GetMRReduction())
			} else {
				log.Printf("DebuffSystem: Target entity %d has no Health component, cannot remove Shred MR reduction", evt.Target)
			}
            s.enqueueStatsRecalculation(evt.Target, evt.Timestamp)
            log.Printf("DebuffSystem: Removed expired Shred from entity %d", evt.Target)
        }
    case debuffs.Sunder:
        if sunder, exists := s.world.GetSunderEffect(evt.Target); exists && sunder.GetSourceId() == evt.SourceId {
            s.world.RemoveComponent(evt.Target, reflect.TypeOf(debuffs.SunderEffect{}))
			// Remove the armor reduction from health component
			healthComp, ok := s.world.GetHealth(evt.Target)
			if ok {
				healthComp.AddBonusArmor(sunder.GetArmorReduction())
			} else {
				log.Printf("DebuffSystem: Target entity %d has no Health component, cannot remove Sunder armor reduction", evt.Target)
			}
            s.enqueueStatsRecalculation(evt.Target, evt.Timestamp)
            log.Printf("DebuffSystem: Removed expired Sunder from entity %d", evt.Target)
        }
    case debuffs.Wound:
        if wound, exists := s.world.GetWoundEffect(evt.Target); exists && wound.GetSourceId() == evt.SourceId {
            s.world.RemoveComponent(evt.Target, reflect.TypeOf(debuffs.WoundEffect{}))
			// Remove the healing reduction from health component
			healthComp, ok := s.world.GetHealth(evt.Target)
			if ok {
				healthComp.SetHealReduction(0.0) // Reset healing reduction
			} else {
				log.Printf("DebuffSystem: Target entity %d has no Health component, cannot remove Wound healing reduction", evt.Target)
			}
            log.Printf("DebuffSystem: Removed expired Wound from entity %d", evt.Target)
        }
    case debuffs.Burn:
        if burn, exists := s.world.GetBurnEffect(evt.Target); exists && burn.GetSourceId() == evt.SourceId {
            s.world.RemoveComponent(evt.Target, reflect.TypeOf(debuffs.BurnEffect{}))
            log.Printf("DebuffSystem: Removed expired Burn from entity %d", evt.Target)
        }
    }
}

func (s *DebuffSystem) removeDebuff(evt eventsys.RemoveDebuffEvent) {
    target := evt.Target
    
    switch evt.DebuffType {
    case debuffs.Shred:
        if shred, exists := s.world.GetShredEffect(target); exists {
            if evt.SourceId == "" || shred.GetSourceId() == evt.SourceId {
                s.world.RemoveComponent(target, reflect.TypeOf(debuffs.ShredEffect{}))
				// Remove the MR reduction from health component
				healthComp, ok := s.world.GetHealth(target)
				if ok {
					healthComp.AddBonusMR(shred.GetMRReduction())
				} else {
					log.Printf("DebuffSystem: Target entity %d has no Health component, cannot remove Shred MR reduction", target)
				}
                s.enqueueStatsRecalculation(target, evt.Timestamp)
                log.Printf("DebuffSystem: Forcibly removed Shred from entity %d", target)
            }
        }
    case debuffs.Sunder:
        if sunder, exists := s.world.GetSunderEffect(target); exists {
            if evt.SourceId == "" || sunder.GetSourceId() == evt.SourceId {
                s.world.RemoveComponent(target, reflect.TypeOf(debuffs.SunderEffect{}))
				// Remove the armor reduction from health component
				healthComp, ok := s.world.GetHealth(target)
				if ok {
					healthComp.AddBonusArmor(sunder.GetArmorReduction())
				} else {
					log.Printf("DebuffSystem: Target entity %d has no Health component, cannot remove Sunder armor reduction", target)
				}
                s.enqueueStatsRecalculation(target, evt.Timestamp)
                log.Printf("DebuffSystem: Forcibly removed Sunder from entity %d", target)
            }
        }
    case debuffs.Wound:
        if wound, exists := s.world.GetWoundEffect(target); exists {
            if evt.SourceId == "" || wound.GetSourceId() == evt.SourceId {
                s.world.RemoveComponent(target, reflect.TypeOf(debuffs.WoundEffect{}))
				// Remove the healing reduction from health component
				healthComp, ok := s.world.GetHealth(target)
				if ok {
					healthComp.SetHealReduction(0.0) // Reset healing reduction
				} else {
					log.Printf("DebuffSystem: Target entity %d has no Health component, cannot remove Wound healing reduction", target)
				}
                log.Printf("DebuffSystem: Forcibly removed Wound from entity %d", target)
            }
        }
    case debuffs.Burn:
        if burn, exists := s.world.GetBurnEffect(target); exists {
            if evt.SourceId == "" || burn.GetSourceId() == evt.SourceId {
                s.world.RemoveComponent(target, reflect.TypeOf(debuffs.BurnEffect{}))
                log.Printf("DebuffSystem: Forcibly removed Burn from entity %d (ticks will stop automatically)", target)
            }
		}
    }
}

func (s *DebuffSystem) enqueueStatsRecalculation(entity entity.Entity, timestamp float64) {
    recalcEvent := eventsys.RecalculateStatsEvent{
        Entity:    entity,
        Timestamp: timestamp,
    }
    s.eventBus.Enqueue(recalcEvent, timestamp)
}