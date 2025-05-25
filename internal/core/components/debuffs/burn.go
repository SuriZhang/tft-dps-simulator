package debuffs

import "tft-dps-simulator/internal/core/entity"

// BurnEffect deals percentage of max HP as true damage over time
type BurnEffect struct {
    damagePercent float64    // Percentage of max HP per second
    duration      float64    // How long the effect lasts
    endTime       float64    // When the effect expires
    lastTickTime  float64    // When damage was last applied
    tickInterval  float64    // How often to apply damage (usually 1.0s)
    sourceEntity  entity.Entity // Who applied this debuff
    sourceType    string     // "Item", "Trait", "Augment", "Spell"
    sourceId      string     // Specific identifier
}

func NewBurnEffect(damagePercent, duration, endTime float64, sourceEntity entity.Entity, sourceType, sourceId string) *BurnEffect {
    return &BurnEffect{
        damagePercent: damagePercent,
        duration:      duration,
        endTime:       endTime,
        lastTickTime:  endTime - duration, // Start ticking immediately
        tickInterval:  1.0,                 // Default 1 second intervals
        sourceEntity:  sourceEntity,
        sourceType:    sourceType,
        sourceId:      sourceId,
    }
}

func (b *BurnEffect) GetDamagePercent() float64 { return b.damagePercent }
func (b *BurnEffect) GetDuration() float64 { return b.duration }
func (b *BurnEffect) GetEndTime() float64 { return b.endTime }
func (b *BurnEffect) GetLastTickTime() float64 { return b.lastTickTime }
func (b *BurnEffect) GetTickInterval() float64 { return b.tickInterval }
func (b *BurnEffect) GetSourceEntity() entity.Entity { return b.sourceEntity }
func (b *BurnEffect) GetSourceType() string { return b.sourceType }
func (b *BurnEffect) GetSourceId() string { return b.sourceId }

func (b *BurnEffect) SetLastTickTime(time float64) { b.lastTickTime = time }

func (b *BurnEffect) IsActive(currentTime float64) bool {
    return currentTime < b.endTime
}

func (b *BurnEffect) ShouldTick(currentTime float64) bool {
    return b.IsActive(currentTime) && (currentTime - b.lastTickTime) >= b.tickInterval
}

func (b *BurnEffect) UpdateFromStrongerEffect(newPercent, newDuration, newEndTime float64, newSource entity.Entity, newSourceType, newSourceId string) {
    if newPercent > b.damagePercent {
        b.damagePercent = newPercent
        b.duration = newDuration
        b.endTime = newEndTime
        b.sourceEntity = newSource
        b.sourceType = newSourceType
        b.sourceId = newSourceId
        // Don't reset lastTickTime to allow immediate ticking if stronger
    } else if newPercent == b.damagePercent && newEndTime > b.endTime {
        b.endTime = newEndTime
        b.duration = newDuration
    }
}