package debuffs

import "tft-dps-simulator/internal/core/entity"

// ShredEffect reduces target's Magic Resist
type ShredEffect struct {
	mrReduction  float64    // Amount of MR reduced
	duration     float64    // How long the effect lasts
	endTime      float64    // When the effect expires
	sourceEntity entity.Entity // Who applied this debuff
	sourceType   string     // "Item", "Trait", "Augment", "Spell"
	sourceId     string     // Specific identifier (e.g., "TFT_Item_LastWhisper")
}

func NewShredEffect(mrReduction, duration, endTime float64, sourceEntity entity.Entity, sourceType, sourceId string) *ShredEffect {
	return &ShredEffect{
		mrReduction:  mrReduction,
		duration:     duration,
		endTime:      endTime,
		sourceEntity: sourceEntity,
		sourceType:   sourceType,
		sourceId:     sourceId,
	}
}

func (s *ShredEffect) GetMRReduction() float64 { return s.mrReduction }
func (s *ShredEffect) GetDuration() float64    { return s.duration }
func (s *ShredEffect) GetEndTime() float64     { return s.endTime }
func (s *ShredEffect) GetSourceEntity() entity.Entity { return s.sourceEntity }
func (s *ShredEffect) GetSourceType() string   { return s.sourceType }
func (s *ShredEffect) GetSourceId() string     { return s.sourceId }

func (s *ShredEffect) IsActive(currentTime float64) bool {
	return currentTime < s.endTime
}

func (s *ShredEffect) UpdateFromStrongerEffect(newReduction, newDuration, newEndTime float64, newSource entity.Entity, newSourceType, newSourceId string) {
	if newReduction > s.mrReduction {
		s.mrReduction = newReduction
		s.duration = newDuration
		s.endTime = newEndTime
		s.sourceEntity = newSource
		s.sourceType = newSourceType
		s.sourceId = newSourceId
	} else if newReduction == s.mrReduction && newEndTime > s.endTime {
		// Same strength, extend duration
		s.endTime = newEndTime
		s.duration = newDuration
	}
}


