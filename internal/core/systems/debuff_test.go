package systems_test

import (
	"testing"

	"tft-dps-simulator/internal/core/components/debuffs"
	"tft-dps-simulator/internal/core/ecs"
	"tft-dps-simulator/internal/core/entity"
	"tft-dps-simulator/internal/core/systems"
	eventsys "tft-dps-simulator/internal/core/systems/events"
	"tft-dps-simulator/internal/core/utils/mock" // Assuming a mock package exists

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

func TestDebuffSystem(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Debuff System Suite")
}

var _ = Describe("DebuffSystem", func() {
	var (
		world         *ecs.World
		mockEventBus  *mock.EventBus // Assuming a mock.EventBus exists
		debuffSystem  *systems.DebuffSystem
		sourceEntity entity.Entity
		targetEntity entity.Entity
	)

	BeforeEach(func() {
		world = ecs.NewWorld()
		mockEventBus = new(mock.EventBus) // Initialize your mock event bus
		debuffSystem = systems.NewDebuffSystem(world, mockEventBus)

		// Register Health component for entities (required by debuff system)
		world.RegisterComponent(ecs.HealthComponentID)
		world.RegisterComponent(ecs.DebuffComponentID) // Assuming a generic DebuffComponentID or specific ones

		// Create source and target entities for tests
		sourceEntity = world.CreateEntity()
		targetEntity = world.CreateEntity()

		// Add Health component to target entity (as debuff system interacts with it)
		hc := &ecs.Health{BaseHP: 1000, CurrentHP: 1000, MaxHP: 1000, BaseArmor: 50, BaseMR: 50, FinalArmor: 50, FinalMR: 50}
		world.AddComponent(targetEntity, hc)
	})

	// Tests for different debuff types will go here
	// e.g., Context("when applying Shred debuff", func() { ... })

	Context("when applying Shred debuff", func() {
		It("should apply Shred and reduce MR", func() {
			applyShredEvent := eventsys.ApplyDebuffEvent{
				Target:     targetEntity,
				Source:     sourceEntity,
				DebuffType: debuffs.Shred,
				Value:      0.3, // 30% MR reduction
				Duration:   5.0,
				Timestamp:  0.0,
				SourceId:   "testSource",
			}
			debuffSystem.HandleEvent(applyShredEvent)

			healthComp, _ := world.GetHealth(targetEntity)
			Expect(healthComp.GetFinalMR()).To(BeNumerically("~", 50*(1-0.3), 0.01)) // Base MR 50

			shredEffect, exists := world.GetShredEffect(targetEntity)
			Expect(exists).To(BeTrue())
			Expect(shredEffect.GetMRReduction()).To(Equal(0.3))
			Expect(shredEffect.GetDuration()).To(Equal(5.0))
			Expect(shredEffect.GetEndTime()).To(Equal(5.0))

			// Check for expiration event
			Expect(mockEventBus.EnqueuedEvents).To(ContainElement(SatisfyAll(
				BeAssignableToTypeOf(eventsys.DebuffExpiredEvent{}),
				WithTransform(func(ev eventsys.Event) eventsys.DebuffExpiredEvent { return ev.(eventsys.DebuffExpiredEvent) },
					SatisfyAll(
						HaveField("Target", targetEntity),
						HaveField("DebuffType", debuffs.Shred),
						HaveField("Timestamp", 5.0),
						HaveField("SourceId", "testSource"),
					),
				),
			)))
		})

		It("should update Shred if a stronger one is applied", func() {
			// Apply initial Shred
			initialShredEvent := eventsys.ApplyDebuffEvent{
				Target: targetEntity, Source: sourceEntity, DebuffType: debuffs.Shred,
				Value: 0.2, Duration: 5.0, Timestamp: 0.0, SourceId: "source1",
			}
			debuffSystem.HandleEvent(initialShredEvent)
			
			healthComp, _ := world.GetHealth(targetEntity)
			Expect(healthComp.GetFinalMR()).To(BeNumerically("~", 50*(1-0.2), 0.01))


			// Apply stronger Shred
			strongerShredEvent := eventsys.ApplyDebuffEvent{
				Target: targetEntity, Source: sourceEntity, DebuffType: debuffs.Shred,
				Value: 0.4, Duration: 3.0, Timestamp: 1.0, SourceId: "source2",
			}
			debuffSystem.HandleEvent(strongerShredEvent)
			
			Expect(healthComp.GetFinalMR()).To(BeNumerically("~", 50*(1-0.4), 0.01))


			shredEffect, _ := world.GetShredEffect(targetEntity)
			Expect(shredEffect.GetMRReduction()).To(Equal(0.4))
			Expect(shredEffect.GetDuration()).To(Equal(3.0))
			Expect(shredEffect.GetEndTime()).To(Equal(1.0 + 3.0)) // Timestamp + new duration
			Expect(shredEffect.GetSourceId()).To(Equal("source2"))
		})

		It("should restore MR when Shred expires", func() {
			applyShredEvent := eventsys.ApplyDebuffEvent{
				Target: targetEntity, Source: sourceEntity, DebuffType: debuffs.Shred,
				Value: 0.3, Duration: 1.0, Timestamp: 0.0, SourceId: "testSource",
			}
			debuffSystem.HandleEvent(applyShredEvent)

			healthComp, _ := world.GetHealth(targetEntity)
			Expect(healthComp.GetFinalMR()).To(BeNumerically("~", 50*(1-0.3), 0.01))

			// Simulate time passing and handle expiration
			debuffExpiredEvent := eventsys.DebuffExpiredEvent{
				Target:     targetEntity,
				DebuffType: debuffs.Shred,
				Timestamp:  1.0, // Expiration time
				SourceId:   "testSource",
			}
			debuffSystem.HandleEvent(debuffExpiredEvent)
			
			// MR should be restored
			Expect(healthComp.GetFinalMR()).To(BeNumerically("~", 50.0, 0.01))
			_, exists := world.GetShredEffect(targetEntity)
			Expect(exists).To(BeFalse())
		})

		It("should remove Shred and restore MR on RemoveDebuffEvent", func() {
			applyShredEvent := eventsys.ApplyDebuffEvent{
				Target: targetEntity, Source: sourceEntity, DebuffType: debuffs.Shred,
				Value: 0.3, Duration: 5.0, Timestamp: 0.0, SourceId: "testSource",
			}
			debuffSystem.HandleEvent(applyShredEvent)

			healthComp, _ := world.GetHealth(targetEntity)
			Expect(healthComp.GetFinalMR()).To(BeNumerically("~", 50*(1-0.3), 0.01))

			removeDebuffEvent := eventsys.RemoveDebuffEvent{
				Target:     targetEntity,
				DebuffType: debuffs.Shred,
				Timestamp:  1.0,
				SourceId:   "testSource",
			}
			debuffSystem.HandleEvent(removeDebuffEvent)

			Expect(healthComp.GetFinalMR()).To(BeNumerically("~", 50.0, 0.01))
			_, exists := world.GetShredEffect(targetEntity)
			Expect(exists).To(BeFalse())
		})
	})

	Context("when applying Sunder debuff", func() {
		It("should apply Sunder and reduce Armor", func() {
			applySunderEvent := eventsys.ApplyDebuffEvent{
				Target:     targetEntity,
				Source:     sourceEntity,
				DebuffType: debuffs.Sunder,
				Value:      0.25, // 25% Armor reduction
				Duration:   4.0,
				Timestamp:  0.0,
				SourceId:   "testSunderSource",
			}
			debuffSystem.HandleEvent(applySunderEvent)

			healthComp, _ := world.GetHealth(targetEntity)
			Expect(healthComp.GetFinalArmor()).To(BeNumerically("~", 50*(1-0.25), 0.01)) // Base Armor 50

			sunderEffect, exists := world.GetSunderEffect(targetEntity)
			Expect(exists).To(BeTrue())
			Expect(sunderEffect.GetArmorReduction()).To(Equal(0.25))
			Expect(sunderEffect.GetDuration()).To(Equal(4.0))
			Expect(sunderEffect.GetEndTime()).To(Equal(4.0))

			// Check for expiration event
			Expect(mockEventBus.EnqueuedEvents).To(ContainElement(SatisfyAll(
				BeAssignableToTypeOf(eventsys.DebuffExpiredEvent{}),
				WithTransform(func(ev eventsys.Event) eventsys.DebuffExpiredEvent { return ev.(eventsys.DebuffExpiredEvent) },
					SatisfyAll(
						HaveField("Target", targetEntity),
						HaveField("DebuffType", debuffs.Sunder),
						HaveField("Timestamp", 4.0),
						HaveField("SourceId", "testSunderSource"),
					),
				),
			)))
		})

		It("should update Sunder if a stronger one is applied", func() {
			// Apply initial Sunder
			initialSunderEvent := eventsys.ApplyDebuffEvent{
				Target: targetEntity, Source: sourceEntity, DebuffType: debuffs.Sunder,
				Value: 0.2, Duration: 5.0, Timestamp: 0.0, SourceId: "sourceS1",
			}
			debuffSystem.HandleEvent(initialSunderEvent)
			healthComp, _ := world.GetHealth(targetEntity)
			Expect(healthComp.GetFinalArmor()).To(BeNumerically("~", 50*(1-0.2), 0.01))

			// Apply stronger Sunder
			strongerSunderEvent := eventsys.ApplyDebuffEvent{
				Target: targetEntity, Source: sourceEntity, DebuffType: debuffs.Sunder,
				Value: 0.5, Duration: 2.0, Timestamp: 1.0, SourceId: "sourceS2",
			}
			debuffSystem.HandleEvent(strongerSunderEvent)
			Expect(healthComp.GetFinalArmor()).To(BeNumerically("~", 50*(1-0.5), 0.01))

			sunderEffect, _ := world.GetSunderEffect(targetEntity)
			Expect(sunderEffect.GetArmorReduction()).To(Equal(0.5))
			Expect(sunderEffect.GetDuration()).To(Equal(2.0))
			Expect(sunderEffect.GetEndTime()).To(Equal(1.0 + 2.0))
			Expect(sunderEffect.GetSourceId()).To(Equal("sourceS2"))
		})

		It("should restore Armor when Sunder expires", func() {
			applySunderEvent := eventsys.ApplyDebuffEvent{
				Target: targetEntity, Source: sourceEntity, DebuffType: debuffs.Sunder,
				Value: 0.25, Duration: 1.0, Timestamp: 0.0, SourceId: "testSunderSource",
			}
			debuffSystem.HandleEvent(applySunderEvent)
			healthComp, _ := world.GetHealth(targetEntity)
			Expect(healthComp.GetFinalArmor()).To(BeNumerically("~", 50*(1-0.25), 0.01))

			debuffExpiredEvent := eventsys.DebuffExpiredEvent{
				Target:     targetEntity,
				DebuffType: debuffs.Sunder,
				Timestamp:  1.0, // Expiration time
				SourceId:   "testSunderSource",
			}
			debuffSystem.HandleEvent(debuffExpiredEvent)
			
			Expect(healthComp.GetFinalArmor()).To(BeNumerically("~", 50.0, 0.01))
			_, exists := world.GetSunderEffect(targetEntity)
			Expect(exists).To(BeFalse())
		})

		It("should remove Sunder and restore Armor on RemoveDebuffEvent", func() {
			applySunderEvent := eventsys.ApplyDebuffEvent{
				Target: targetEntity, Source: sourceEntity, DebuffType: debuffs.Sunder,
				Value: 0.25, Duration: 5.0, Timestamp: 0.0, SourceId: "testSunderSource",
			}
			debuffSystem.HandleEvent(applySunderEvent)
			healthComp, _ := world.GetHealth(targetEntity)
			Expect(healthComp.GetFinalArmor()).To(BeNumerically("~", 50*(1-0.25), 0.01))

			removeDebuffEvent := eventsys.RemoveDebuffEvent{
				Target:     targetEntity,
				DebuffType: debuffs.Sunder,
				Timestamp:  1.0,
				SourceId:   "testSunderSource",
			}
			debuffSystem.HandleEvent(removeDebuffEvent)

			Expect(healthComp.GetFinalArmor()).To(BeNumerically("~", 50.0, 0.01))
			_, exists := world.GetSunderEffect(targetEntity)
			Expect(exists).To(BeFalse())
		})
	})

	Context("when applying Wound debuff", func() {
		It("should apply Wound and set healing reduction", func() {
			applyWoundEvent := eventsys.ApplyDebuffEvent{
				Target:     targetEntity,
				Source:     sourceEntity,
				DebuffType: debuffs.Wound,
				Value:      0.5, // 50% healing reduction
				Duration:   3.0,
				Timestamp:  0.0,
				SourceId:   "testWoundSource",
			}
			debuffSystem.HandleEvent(applyWoundEvent)

			healthComp, _ := world.GetHealth(targetEntity)
			Expect(healthComp.GetHealReduction()).To(Equal(0.5))

			woundEffect, exists := world.GetWoundEffect(targetEntity)
			Expect(exists).To(BeTrue())
			Expect(woundEffect.GetHealingReduction()).To(Equal(0.5))
			Expect(woundEffect.GetDuration()).To(Equal(3.0))
			Expect(woundEffect.GetEndTime()).To(Equal(3.0))

			// Check for expiration event
			Expect(mockEventBus.EnqueuedEvents).To(ContainElement(SatisfyAll(
				BeAssignableToTypeOf(eventsys.DebuffExpiredEvent{}),
				WithTransform(func(ev eventsys.Event) eventsys.DebuffExpiredEvent { return ev.(eventsys.DebuffExpiredEvent) },
					SatisfyAll(
						HaveField("Target", targetEntity),
						HaveField("DebuffType", debuffs.Wound),
						HaveField("Timestamp", 3.0),
						HaveField("SourceId", "testWoundSource"),
					),
				),
			)))
		})

		It("should update Wound if a stronger one is applied", func() {
			// Apply initial Wound
			initialWoundEvent := eventsys.ApplyDebuffEvent{
				Target: targetEntity, Source: sourceEntity, DebuffType: debuffs.Wound,
				Value: 0.3, Duration: 5.0, Timestamp: 0.0, SourceId: "sourceW1",
			}
			debuffSystem.HandleEvent(initialWoundEvent)
			healthComp, _ := world.GetHealth(targetEntity)
			Expect(healthComp.GetHealReduction()).To(Equal(0.3))

			// Apply stronger Wound
			strongerWoundEvent := eventsys.ApplyDebuffEvent{
				Target: targetEntity, Source: sourceEntity, DebuffType: debuffs.Wound,
				Value: 0.6, Duration: 2.5, Timestamp: 1.0, SourceId: "sourceW2",
			}
			debuffSystem.HandleEvent(strongerWoundEvent)
			Expect(healthComp.GetHealReduction()).To(Equal(0.6))

			woundEffect, _ := world.GetWoundEffect(targetEntity)
			Expect(woundEffect.GetHealingReduction()).To(Equal(0.6))
			Expect(woundEffect.GetDuration()).To(Equal(2.5))
			Expect(woundEffect.GetEndTime()).To(Equal(1.0 + 2.5))
			Expect(woundEffect.GetSourceId()).To(Equal("sourceW2"))
		})

		It("should remove healing reduction when Wound expires", func() {
			applyWoundEvent := eventsys.ApplyDebuffEvent{
				Target: targetEntity, Source: sourceEntity, DebuffType: debuffs.Wound,
				Value: 0.5, Duration: 1.0, Timestamp: 0.0, SourceId: "testWoundSource",
			}
			debuffSystem.HandleEvent(applyWoundEvent)
			healthComp, _ := world.GetHealth(targetEntity)
			Expect(healthComp.GetHealReduction()).To(Equal(0.5))

			debuffExpiredEvent := eventsys.DebuffExpiredEvent{
				Target:     targetEntity,
				DebuffType: debuffs.Wound,
				Timestamp:  1.0, // Expiration time
				SourceId:   "testWoundSource",
			}
			debuffSystem.HandleEvent(debuffExpiredEvent)
			
			Expect(healthComp.GetHealReduction()).To(Equal(0.0)) // Should be reset
			_, exists := world.GetWoundEffect(targetEntity)
			Expect(exists).To(BeFalse())
		})

		It("should remove Wound and healing reduction on RemoveDebuffEvent", func() {
			applyWoundEvent := eventsys.ApplyDebuffEvent{
				Target: targetEntity, Source: sourceEntity, DebuffType: debuffs.Wound,
				Value: 0.5, Duration: 5.0, Timestamp: 0.0, SourceId: "testWoundSource",
			}
			debuffSystem.HandleEvent(applyWoundEvent)
			healthComp, _ := world.GetHealth(targetEntity)
			Expect(healthComp.GetHealReduction()).To(Equal(0.5))

			removeDebuffEvent := eventsys.RemoveDebuffEvent{
				Target:     targetEntity,
				DebuffType: debuffs.Wound,
				Timestamp:  1.0,
				SourceId:   "testWoundSource",
			}
			debuffSystem.HandleEvent(removeDebuffEvent)

			Expect(healthComp.GetHealReduction()).To(Equal(0.0))
			_, exists := world.GetWoundEffect(targetEntity)
			Expect(exists).To(BeFalse())
		})
	})

	Context("when applying Burn debuff", func() {
		It("should apply Burn and enqueue initial tick and expiration", func() {
			applyBurnEvent := eventsys.ApplyDebuffEvent{
				Target:     targetEntity,
				Source:     sourceEntity,
				DebuffType: debuffs.Burn,
				Value:      0.01, // 1% max HP per second
				Duration:   3.0,
				Timestamp:  0.0,
				SourceId:   "testBurnSource",
			}
			debuffSystem.HandleEvent(applyBurnEvent)

			burnEffect, exists := world.GetBurnEffect(targetEntity)
			Expect(exists).To(BeTrue())
			Expect(burnEffect.GetDamagePercent()).To(Equal(0.01))
			Expect(burnEffect.GetDuration()).To(Equal(3.0))
			Expect(burnEffect.GetEndTime()).To(Equal(3.0))

			// Check for initial BurnTickEvent
			Expect(mockEventBus.EnqueuedEvents).To(ContainElement(SatisfyAll(
				BeAssignableToTypeOf(eventsys.BurnTickEvent{}),
				WithTransform(func(ev eventsys.Event) eventsys.BurnTickEvent { return ev.(eventsys.BurnTickEvent) },
					SatisfyAll(
						HaveField("Target", targetEntity),
						HaveField("Source", sourceEntity),
						HaveField("Timestamp", 1.0), // First tick after 1s
						HaveField("SourceId", "testBurnSource"),
					),
				),
			)))

			// Check for expiration event
			Expect(mockEventBus.EnqueuedEvents).To(ContainElement(SatisfyAll(
				BeAssignableToTypeOf(eventsys.DebuffExpiredEvent{}),
				WithTransform(func(ev eventsys.Event) eventsys.DebuffExpiredEvent { return ev.(eventsys.DebuffExpiredEvent) },
					SatisfyAll(
						HaveField("Target", targetEntity),
						HaveField("DebuffType", debuffs.Burn),
						HaveField("Timestamp", 3.0),
						HaveField("SourceId", "testBurnSource"),
					),
				),
			)))
		})

		It("should update Burn if a stronger one is applied", func() {
			// Apply initial Burn
			initialBurnEvent := eventsys.ApplyDebuffEvent{
				Target: targetEntity, Source: sourceEntity, DebuffType: debuffs.Burn,
				Value: 0.01, Duration: 5.0, Timestamp: 0.0, SourceId: "sourceB1",
			}
			debuffSystem.HandleEvent(initialBurnEvent)

			// Apply stronger Burn
			strongerBurnEvent := eventsys.ApplyDebuffEvent{
				Target: targetEntity, Source: sourceEntity, DebuffType: debuffs.Burn,
				Value: 0.02, Duration: 2.0, Timestamp: 1.0, SourceId: "sourceB2",
			}
			debuffSystem.HandleEvent(strongerBurnEvent)

			burnEffect, _ := world.GetBurnEffect(targetEntity)
			Expect(burnEffect.GetDamagePercent()).To(Equal(0.02))
			Expect(burnEffect.GetDuration()).To(Equal(2.0))
			Expect(burnEffect.GetEndTime()).To(Equal(1.0 + 2.0))
			Expect(burnEffect.GetSourceId()).To(Equal("sourceB2"))

			// Ensure new expiration is enqueued for the stronger burn
			Expect(mockEventBus.EnqueuedEvents).To(ContainElement(SatisfyAll(
				BeAssignableToTypeOf(eventsys.DebuffExpiredEvent{}),
				WithTransform(func(ev eventsys.Event) eventsys.DebuffExpiredEvent { return ev.(eventsys.DebuffExpiredEvent) },
					SatisfyAll(
						HaveField("Target", targetEntity),
						HaveField("DebuffType", debuffs.Burn),
						HaveField("Timestamp", 1.0 + 2.0), // New end time
						HaveField("SourceId", "sourceB2"),
					),
				),
			)))
		})
		
		It("should deal damage on BurnTickEvent and enqueue next tick", func() {
			// Apply Burn first
			applyBurnEvent := eventsys.ApplyDebuffEvent{
				Target:     targetEntity,
				Source:     sourceEntity,
				DebuffType: debuffs.Burn,
				Value:      0.01, // 1% max HP
				Duration:   3.0,  // Will tick at 1.0, 2.0
				Timestamp:  0.0,
				SourceId:   "testBurnTickSource",
			}
			debuffSystem.HandleEvent(applyBurnEvent)
			mockEventBus.Clear() // Clear apply-related events

			// Simulate first tick
			burnTickEvent := eventsys.BurnTickEvent{
				Target:    targetEntity,
				Source:    sourceEntity,
				Timestamp: 1.0,
				SourceId:  "testBurnTickSource",
			}
			debuffSystem.HandleEvent(burnTickEvent)

			// Check for DamageAppliedEvent
			healthComp, _ := world.GetHealth(targetEntity) // MaxHP is 1000
			expectedDamage := 1000 * 0.01
			Expect(mockEventBus.EnqueuedEvents).To(ContainElement(SatisfyAll(
				BeAssignableToTypeOf(eventsys.DamageAppliedEvent{}),
				WithTransform(func(ev eventsys.Event) eventsys.DamageAppliedEvent { return ev.(eventsys.DamageAppliedEvent) },
					SatisfyAll(
						HaveField("Target", targetEntity),
						HaveField("Source", sourceEntity),
						HaveField("DamageType", "True"),
						HaveField("DamageSource", "Burn"),
						HaveField("FinalTotalDamage", expectedDamage),
					),
				),
			)))

			// Check for next BurnTickEvent
			Expect(mockEventBus.EnqueuedEvents).To(ContainElement(SatisfyAll(
				BeAssignableToTypeOf(eventsys.BurnTickEvent{}),
				WithTransform(func(ev eventsys.Event) eventsys.BurnTickEvent { return ev.(eventsys.BurnTickEvent) },
					SatisfyAll(
						HaveField("Target", targetEntity),
						HaveField("Source", sourceEntity),
						HaveField("Timestamp", 2.0), // Next tick
						HaveField("SourceId", "testBurnTickSource"),
					),
				),
			)))
		})

		It("should stop ticking after Burn expires", func() {
			applyBurnEvent := eventsys.ApplyDebuffEvent{
				Target:     targetEntity,
				Source:     sourceEntity,
				DebuffType: debuffs.Burn,
				Value:      0.01,
				Duration:   1.5, // Should tick at 1.0, expire at 1.5
				Timestamp:  0.0,
				SourceId:   "testBurnExpireSource",
			}
			debuffSystem.HandleEvent(applyBurnEvent)
			mockEventBus.Clear()

			// First tick
			burnTickEvent1 := eventsys.BurnTickEvent{
				Target: targetEntity, Source: sourceEntity, Timestamp: 1.0, SourceId: "testBurnExpireSource",
			}
			debuffSystem.HandleEvent(burnTickEvent1)
			Expect(mockEventBus.EnqueuedEvents).To(ContainElement(WithTransform(func(ev eventsys.Event) string { return ev.(eventsys.DamageAppliedEvent).DamageSource }, Equal("Burn"))))
			
			// Expiration
			debuffExpiredEvent := eventsys.DebuffExpiredEvent{
				Target: targetEntity, DebuffType: debuffs.Burn, Timestamp: 1.5, SourceId: "testBurnExpireSource",
			}
			debuffSystem.HandleEvent(debuffExpiredEvent) // Process expiration
			mockEventBus.Clear() // Clear previous events

			// Attempt to process a tick that would have occurred after expiration
			burnTickEvent2 := eventsys.BurnTickEvent{
				Target: targetEntity, Source: sourceEntity, Timestamp: 2.0, SourceId: "testBurnExpireSource",
			}
			debuffSystem.HandleEvent(burnTickEvent2)
			
			// No damage event should be queued because burn is gone
			Expect(mockEventBus.EnqueuedEvents).NotTo(ContainElement(BeAssignableToTypeOf(eventsys.DamageAppliedEvent{})))
			// No next tick event should be queued
			Expect(mockEventBus.EnqueuedEvents).NotTo(ContainElement(BeAssignableToTypeOf(eventsys.BurnTickEvent{})))
		})

		It("should remove Burn on RemoveDebuffEvent and stop ticks", func() {
			applyBurnEvent := eventsys.ApplyDebuffEvent{
				Target:     targetEntity,
				Source:     sourceEntity,
				DebuffType: debuffs.Burn,
				Value:      0.01,
				Duration:   3.0,
				Timestamp:  0.0,
				SourceId:   "testBurnRemoveSource",
			}
			debuffSystem.HandleEvent(applyBurnEvent)
			_, exists := world.GetBurnEffect(targetEntity)
			Expect(exists).To(BeTrue())
			mockEventBus.Clear()

			removeDebuffEvent := eventsys.RemoveDebuffEvent{
				Target:     targetEntity,
				DebuffType: debuffs.Burn,
				Timestamp:  0.5, // Remove before first tick
				SourceId:   "testBurnRemoveSource",
			}
			debuffSystem.HandleEvent(removeDebuffEvent)

			_, exists = world.GetBurnEffect(targetEntity)
			Expect(exists).To(BeFalse())

			// Attempt to process a tick that would have occurred if not removed
			burnTickEvent := eventsys.BurnTickEvent{
				Target: targetEntity, Source: sourceEntity, Timestamp: 1.0, SourceId: "testBurnRemoveSource",
			}
			debuffSystem.HandleEvent(burnTickEvent)
			Expect(mockEventBus.EnqueuedEvents).NotTo(ContainElement(BeAssignableToTypeOf(eventsys.DamageAppliedEvent{})))
			Expect(mockEventBus.EnqueuedEvents).NotTo(ContainElement(BeAssignableToTypeOf(eventsys.BurnTickEvent{})))
		})
	})
})
