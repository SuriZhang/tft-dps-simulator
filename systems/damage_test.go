package systems_test

import (
	"github.com/suriz/tft-dps-simulator/components"
	"github.com/suriz/tft-dps-simulator/ecs"
	"github.com/suriz/tft-dps-simulator/factory"
	"github.com/suriz/tft-dps-simulator/systems"
	eventsys "github.com/suriz/tft-dps-simulator/systems/events"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// --- Test Suite ---
var _ = Describe("DamageSystem", func() {
	var (
		world           *ecs.World
		eventBus        *MockEventBus
		championFactory *factory.ChampionFactory
		damageSystem    *systems.DamageSystem
		attacker        ecs.Entity
		target          ecs.Entity
		attackerAttack  *components.Attack
		attackerCrit	*components.Crit
		attackerMana    *components.Mana
		targetHealth    *components.Health
		ok              bool
	)

	BeforeEach(func() {
		world = ecs.NewWorld()
		eventBus = NewMockEventBus()
		championFactory = factory.NewChampionFactory(world)
		damageSystem = systems.NewDamageSystem(world, eventBus) // Pass bus to system

		// --- Create Attacker ---
		var err error
		attacker, err = championFactory.CreatePlayerChampion("TFT_BlueGolem", 1) // Using a known champ
		Expect(err).NotTo(HaveOccurred())
		// Add mana component if not present by default
		if _, ok := world.GetMana(attacker); !ok {
			world.AddComponent(attacker, components.NewMana(0, 100)) // Example mana
		}

		// --- Create Target ---
		target, err = championFactory.CreateEnemyChampion("TFT_TrainingDummy", 1)
		Expect(err).NotTo(HaveOccurred())

		// --- Get Components ---
		attackerAttack, ok = world.GetAttack(attacker)
		Expect(ok).To(BeTrue())
		attackerMana, ok = world.GetMana(attacker)
		Expect(ok).To(BeTrue())
		attackerCrit, ok = world.GetCrit(attacker)
		Expect(ok).To(BeTrue())
		targetHealth, ok = world.GetHealth(target)
		Expect(ok).To(BeTrue())

		// --- Set Base Stats for Predictable Calculations ---
		attackerAttack.SetFinalAD(100.0)
		attackerAttack.SetFinalAttackSpeed(1.0) // Needed for LastAttackTime update
		attackerCrit.SetFinalCritChance(0.25)
		attackerCrit.SetFinalCritMultiplier(1.5) // 50% bonus crit damage
		attackerAttack.SetFinalDamageAmp(0.0)      // No initial amp

		attackerMana.SetCurrentMana(0)
		attackerMana.SetMaxMana(100)

		targetHealth.SetCurrentHealth(1000.0)
		targetHealth.SetFinalArmor(50.0)
		targetHealth.SetFinalMR(50.0) // Not used in AD test, but good practice
		targetHealth.SetFinalDurability(0.0)
	})

	Describe("handling AttackLandedEvent", func() {
		var (
			attackEvent eventsys.AttackLandedEvent
			eventTime   float64 = 1.23 // Arbitrary time for the event
		)

		BeforeEach(func() {
			attackEvent = eventsys.AttackLandedEvent{
				Source:     attacker,
				Target:     target,
				BaseDamage: attackerAttack.GetFinalAD(), // Use the AD set in outer BeforeEach
				Timestamp:  eventTime,
			}
			eventBus.ClearEvents() // Ensure no stale events
		})

		It("should calculate final damage considering crit EV and armor", func() {
			damageSystem.HandleEvent(attackEvent)

			Expect(eventBus.EnqueuedEvents).To(HaveLen(1), "Should enqueue one DamageAppliedEvent")
			damageAppliedEvent, ok := eventBus.GetLastEvent().(eventsys.DamageAppliedEvent)
			Expect(ok).To(BeTrue(), "Enqueued event should be DamageAppliedEvent")

			// Manual Calculation:
			// Base = 100
			// Crit EV Multiplier = (1 - 0.25) + (0.25 * 1.5) = 0.75 + 0.375 = 1.125
			// Damage after Crit EV = 100 * 1.125 = 112.5
			// Armor Reduction = 100 / (100 + 50) = 100 / 150 = 0.666...
			// Final Damage = 112.5 * 0.666... = 75.0
			expectedFinalDamage := 75.0

			Expect(damageAppliedEvent.Source).To(Equal(attacker))
			Expect(damageAppliedEvent.Target).To(Equal(target))
			Expect(damageAppliedEvent.FinalDamage).To(BeNumerically("~", expectedFinalDamage, 0.01))
			Expect(damageAppliedEvent.Timestamp).To(Equal(eventTime)) // Timestamp should be passed through
		})

		It("should update the attacker's LastAttackTime", func() {
			initialLastAttackTime := attackerAttack.GetLastAttackTime()
			Expect(initialLastAttackTime).NotTo(Equal(eventTime)) // Ensure it's not already set

			damageSystem.HandleEvent(attackEvent)

			Expect(attackerAttack.GetLastAttackTime()).To(Equal(eventTime))
		})

		It("should consider damage amplification", func() {
			attackerAttack.SetFinalDamageAmp(0.10)               // 10% damage amp
			attackEvent.BaseDamage = attackerAttack.GetFinalAD() // Update event base damage if needed

			damageSystem.HandleEvent(attackEvent)

			Expect(eventBus.EnqueuedEvents).To(HaveLen(1))
			damageAppliedEvent, _ := eventBus.GetLastEvent().(eventsys.DamageAppliedEvent)

			// Manual Calculation:
			// Base = 100
			// Amp = 10% -> 1.1 multiplier
			// Damage after Amp = 100 * 1.1 = 110
			// Crit EV Multiplier = 1.125
			// Damage after Crit EV = 110 * 1.125 = 123.75
			// Armor Reduction = 0.666...
			// Final Damage = 123.75 * 0.666... = 82.5
			expectedFinalDamage := 82.5

			Expect(damageAppliedEvent.FinalDamage).To(BeNumerically("~", expectedFinalDamage, 0.01))
		})

		It("should consider durability", func() {
			targetHealth.SetFinalDurability(0.20) // 20% durability

			damageSystem.HandleEvent(attackEvent)

			Expect(eventBus.EnqueuedEvents).To(HaveLen(1))
			damageAppliedEvent, _ := eventBus.GetLastEvent().(eventsys.DamageAppliedEvent)

			// Manual Calculation:
			// Base = 100
			// Crit EV Multiplier = 1.125
			// Damage after Crit EV = 112.5
			// Armor Reduction = 100 / (100 + 50) = 0.666...
			// Durability Multiplier = (1 - 0.20) = 0.80
			// Final Damage = 112.5 * 0.666... * 0.80 = 75.0 * 0.8 = 60.0
			expectedFinalDamage := 60.0

			Expect(damageAppliedEvent.FinalDamage).To(BeNumerically("~", expectedFinalDamage, 0.01))
		})
	})

	Describe("handling DamageAppliedEvent", func() {
		var (
			damageEvent eventsys.DamageAppliedEvent
			eventTime   float64 = 4.56 // Arbitrary time
		)

		BeforeEach(func() {
			damageEvent = eventsys.DamageAppliedEvent{
				Source:      attacker,
				Target:      target,
				FinalDamage: 50.0, // Simple damage value for these tests
				Timestamp:   eventTime,
			}
			eventBus.ClearEvents()
			attackerMana.SetCurrentMana(10)    // Reset mana
			targetHealth.SetCurrentHealth(200) // Reset health
		})

		It("should decrease target health by FinalDamage", func() {
			initialHP := targetHealth.CurrentHP
			damageSystem.HandleEvent(damageEvent)
			Expect(targetHealth.CurrentHP).To(Equal(initialHP - damageEvent.FinalDamage))
		})

		It("should increase attacker mana by the standard amount", func() {
			initialMana := attackerMana.GetCurrentMana()
			expectedManaGain := 10.0 // Standard gain defined in DamageSystem
			damageSystem.HandleEvent(damageEvent)
			Expect(attackerMana.GetCurrentMana()).To(Equal(initialMana + expectedManaGain))
		})

		It("should clamp attacker mana gain to MaxMana", func() {
			attackerMana.SetCurrentMana(95) // Start close to max
			attackerMana.SetMaxMana(100)
			damageSystem.HandleEvent(damageEvent)                  // Tries to add 10 mana
			Expect(attackerMana.GetCurrentMana()).To(Equal(100.0)) // Should be clamped
		})

		It("should not enqueue a DeathEvent if target survives", func() {
			targetHealth.SetCurrentHealth(damageEvent.FinalDamage + 1) // Ensure survival
			damageSystem.HandleEvent(damageEvent)
			Expect(eventBus.EnqueuedEvents).To(BeEmpty())
		})

		Context("when damage is lethal", func() {
			BeforeEach(func() {
				targetHealth.SetCurrentHealth(damageEvent.FinalDamage - 1) // Ensure lethal damage
			})

			It("should set target health to zero or below", func() {
				damageSystem.HandleEvent(damageEvent)
				Expect(targetHealth.CurrentHP).To(BeNumerically("<=", 0))
			})

			It("should enqueue a DeathEvent and KillEvent", func() {
				damageSystem.HandleEvent(damageEvent)
				events := eventBus.GetAllEvents()
				Expect(events).To(HaveLen(2))
				// Check events contain one DeathEvent and one KillEvent
				Expect(events[0]).To(BeAssignableToTypeOf(eventsys.DeathEvent{}), "First event should be DeathEvent")
				Expect(events[1]).To(BeAssignableToTypeOf(eventsys.KillEvent{}), "Second event should be KillEvent")
				deathEvent, _ := events[0].(eventsys.DeathEvent)
				Expect(deathEvent.Target).To(Equal(target))
				Expect(deathEvent.Timestamp).To(Equal(eventTime)) 

				killEvent, _ := events[1].(eventsys.KillEvent)
				Expect(killEvent.Killer).To(Equal(attacker))
				Expect(killEvent.Victim).To(Equal(target))
				Expect(killEvent.Timestamp).To(Equal(eventTime)) // Should match the damage event timestamp
			})

			It("should still grant mana to the attacker", func() {
				initialMana := attackerMana.GetCurrentMana()
				expectedManaGain := 10.0
				damageSystem.HandleEvent(damageEvent)
				Expect(attackerMana.GetCurrentMana()).To(Equal(initialMana + expectedManaGain))
			})

			It("should only enqueue one DeathEvent and KillEvent even if called again", func() {
				damageSystem.HandleEvent(damageEvent) // First lethal hit
				events := eventBus.GetAllEvents()
				Expect(events).To(HaveLen(2))
				// Check events contain one DeathEvent and one KillEvent
				Expect(events[0]).To(BeAssignableToTypeOf(eventsys.DeathEvent{}))
				Expect(events[1]).To(BeAssignableToTypeOf(eventsys.KillEvent{}))

				// Simulate another damage event hitting the already dead target
				eventBus.ClearEvents()
				damageSystem.HandleEvent(damageEvent)
				Expect(eventBus.EnqueuedEvents).To(BeEmpty(), "No additional DeathEvent should be sent for an already dead target")
			})
		})
	})
})
