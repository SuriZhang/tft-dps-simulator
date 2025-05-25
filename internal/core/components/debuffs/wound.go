package debuffs

import "tft-dps-simulator/internal/core/entity"

// WoundEffect reduces healing received
type WoundEffect struct {
    healingReduction float64    // Percentage reduction (0.0-1.0)
    duration         float64    // How long the effect lasts
    endTime          float64    // When the effect expires
    sourceEntity     entity.Entity // Who applied this debuff
    sourceType       string     // "Item", "Trait", "Augment", "Spell"
    sourceId         string     // Specific identifier
}

func NewWoundEffect(healingReduction, duration, endTime float64, sourceEntity entity.Entity, sourceType, sourceId string) *WoundEffect {
    return &WoundEffect{
        healingReduction: healingReduction,
        duration:         duration,
        endTime:          endTime,
        sourceEntity:     sourceEntity,
        sourceType:       sourceType,
        sourceId:         sourceId,
    }
}

func (w *WoundEffect) GetHealingReduction() float64 { return w.healingReduction }
func (w *WoundEffect) GetDuration() float64 { return w.duration }
func (w *WoundEffect) GetEndTime() float64 { return w.endTime }
func (w *WoundEffect) GetSourceEntity() entity.Entity { return w.sourceEntity }
func (w *WoundEffect) GetSourceType() string { return w.sourceType }
func (w *WoundEffect) GetSourceId() string { return w.sourceId }

func (w *WoundEffect) IsActive(currentTime float64) bool {
    return currentTime < w.endTime
}

func (w *WoundEffect) UpdateFromStrongerEffect(newReduction, newDuration, newEndTime float64, newSource entity.Entity, newSourceType, newSourceId string) {
    if newReduction > w.healingReduction {
        w.healingReduction = newReduction
        w.duration = newDuration
        w.endTime = newEndTime
        w.sourceEntity = newSource
        w.sourceType = newSourceType
        w.sourceId = newSourceId
    } else if newReduction == w.healingReduction && newEndTime > w.endTime {
        w.endTime = newEndTime
        w.duration = newDuration
    }
}