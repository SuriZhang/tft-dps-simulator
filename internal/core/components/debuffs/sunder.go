package debuffs

import "tft-dps-simulator/internal/core/entity"

// SunderEffect reduces target's Armor
type SunderEffect struct {
    armorReduction float64    // Amount of armor reduced
    duration       float64    // How long the effect lasts
    endTime        float64    // When the effect expires
    sourceEntity   entity.Entity // Who applied this debuff
    sourceType     string     // "Item", "Trait", "Augment", "Spell"
    sourceId       string     // Specific identifier
}

func NewSunderEffect(armorReduction, duration, endTime float64, sourceEntity entity.Entity, sourceType, sourceId string) *SunderEffect {
    return &SunderEffect{
        armorReduction: armorReduction,
        duration:       duration,
        endTime:        endTime,
        sourceEntity:   sourceEntity,
        sourceType:     sourceType,
        sourceId:       sourceId,
    }
}

func (s *SunderEffect) GetArmorReduction() float64 { return s.armorReduction }
func (s *SunderEffect) GetDuration() float64 { return s.duration }
func (s *SunderEffect) GetEndTime() float64 { return s.endTime }
func (s *SunderEffect) GetSourceEntity() entity.Entity { return s.sourceEntity }
func (s *SunderEffect) GetSourceType() string { return s.sourceType }
func (s *SunderEffect) GetSourceId() string { return s.sourceId }

func (s *SunderEffect) IsActive(currentTime float64) bool {
    return currentTime < s.endTime
}

func (s *SunderEffect) UpdateFromStrongerEffect(newReduction, newDuration, newEndTime float64, newSource entity.Entity, newSourceType, newSourceId string) {
    if newReduction > s.armorReduction {
        s.armorReduction = newReduction
        s.duration = newDuration
        s.endTime = newEndTime
        s.sourceEntity = newSource
        s.sourceType = newSourceType
        s.sourceId = newSourceId
    } else if newReduction == s.armorReduction && newEndTime > s.endTime {
        s.endTime = newEndTime
        s.duration = newDuration
    }
}