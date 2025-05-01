package systems_test

import (
	"reflect"

	"tft-dps-simulator/internal/core/components"
	"tft-dps-simulator/internal/core/ecs"
	"tft-dps-simulator/internal/core/factory"
	"tft-dps-simulator/internal/core/systems"
	eventsys "tft-dps-simulator/internal/core/systems/events"
	"tft-dps-simulator/internal/core/utils"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DamageSystem", func() {
    var (
        world           *ecs.World
        mockEventBus        *utils.MockEventBus
        championFactory *factory.ChampionFactory
        damageSystem    *systems.DamageSystem
        attacker        ecs.Entity
        target          ecs.Entity
        attackerAttack  *components.Attack
        attackerCrit    *components.Crit
        attackerMana    *components.Mana
        attackerSpell   *components.Spell // Added for spell tests
        targetHealth    *components.Health
        ok              bool
    )

    BeforeEach(func() {
        world = ecs.NewWorld()
        mockEventBus = utils.NewMockEventBus()
        championFactory = factory.NewChampionFactory(world) // Factory now adds Crit component
        damageSystem = systems.NewDamageSystem(world, mockEventBus)
		mockEventBus.RegisterHandler(damageSystem)

        // --- Create Attacker ---
        var err error
        // Use a champion known to have Attack, Crit, Mana, Spell components
        attacker, err = championFactory.CreatePlayerChampion("TFT14_Kindred", 1)
        Expect(err).NotTo(HaveOccurred())

        // --- Create Target ---
        target, err = championFactory.CreateEnemyChampion("TFT_TrainingDummy", 1)
        Expect(err).NotTo(HaveOccurred())

        // --- Get Components ---
        attackerAttack, ok = world.GetAttack(attacker)
        Expect(ok).To(BeTrue())
        attackerCrit, ok = world.GetCrit(attacker) // Expect Crit component now
        Expect(ok).To(BeTrue())
        attackerMana, ok = world.GetMana(attacker)
        Expect(ok).To(BeTrue())
        attackerSpell, ok = world.GetSpell(attacker) // Expect Spell component
        Expect(ok).To(BeTrue())
        targetHealth, ok = world.GetHealth(target)
        Expect(ok).To(BeTrue())

        // --- Set Base Stats for Predictable Calculations ---
        // Attack Stats
        attackerAttack.SetFinalAD(100.0)
        attackerAttack.SetFinalAttackSpeed(1.0)
        attackerAttack.SetFinalDamageAmp(0.0) // No initial amp

        // Crit Stats (Shared)
        attackerCrit.SetFinalCritChance(0.25)
        attackerCrit.SetFinalCritMultiplier(1.5) // Base 1.4 + 0.1 bonus = 1.5

        // Spell Stats
        attackerSpell.SetFinalAP(100.0) // Example AP

        // Mana Stats
        attackerMana.SetCurrentMana(0)
        attackerMana.SetMaxMana(100)

        // Target Stats
        targetHealth.SetCurrentHP(1000.0)
        targetHealth.SetFinalArmor(50.0)
        targetHealth.SetFinalMR(50.0)
        targetHealth.SetFinalDurability(0.0) // No initial durability
    })

    Describe("handling AttackLandedEvent", func() {
        var (
            attackEvent eventsys.AttackLandedEvent
            eventTime   float64 = 1.23
        )

        // Only create the event structure in BeforeEach
        BeforeEach(func() {
            attackEvent = eventsys.AttackLandedEvent{
                Source:     attacker,
                Target:     target,
                BaseDamage: attackerAttack.GetFinalAD(), // Base AD is 100 from outer BeforeEach
                Timestamp:  eventTime,
            }
            // Reset potentially modified stats before each test in this context
            attackerAttack.SetFinalDamageAmp(0.0)
            targetHealth.SetFinalDurability(0.0)
            mockEventBus.ClearEvents() // Clear events from any previous test runs
        })

        It("should calculate final physical damage considering crit EV and armor", func() {
            // --- Arrange ---
            // Initial stats (AD=100, CritC=0.25, CritM=1.5, Armor=50) are set in the outer BeforeEach

            // --- Act ---
            // Enqueue and process the event *within this specific test*
            mockEventBus.Enqueue(attackEvent, eventTime)
            mockEventBus.ProcessNext() // Triggers damageSystem.HandleEvent

            // --- Assert ---
            // Check the events enqueued *by the handler* during ProcessNext
            enqueuedEvents := mockEventBus.GetAllEvents() // Get events added since ClearEvents
            Expect(enqueuedEvents).To(HaveLen(1), "Should enqueue one DamageAppliedEvent")

            damageAppliedEvent, ok := enqueuedEvents[0].(eventsys.DamageAppliedEvent)
            Expect(ok).To(BeTrue(), "Enqueued event should be DamageAppliedEvent")

            // Manual Calculation:
            // Base Physical = 100
            // Crit EV Multiplier = (1 - 0.25) + (0.25 * 1.5) = 0.75 + 0.375 = 1.125
            // Amp Multiplier = 1.0
            // PreMitigation Physical = 100 * 1.125 * 1.0 = 112.5
            // Armor Reduction = 100 / (100 + 50) = 0.666...
            // Durability Multiplier = 1.0
            // Final Physical = 112.5 * 0.666... * 1.0 = 75.0
            expectedPreMitPhys := 112.5
            expectedFinalPhys := 75.0
            expectedTotalDamage := expectedFinalPhys // Only physical damage

            Expect(damageAppliedEvent.Source).To(Equal(attacker))
            Expect(damageAppliedEvent.Target).To(Equal(target))
            Expect(damageAppliedEvent.Timestamp).To(Equal(eventTime))
            Expect(damageAppliedEvent.DamageSource).To(Equal("Attack"))
            Expect(damageAppliedEvent.PreMitigationDamage).To(BeNumerically("~", expectedPreMitPhys, 0.01))
            Expect(damageAppliedEvent.FinalTotalDamage).To(BeNumerically("~", expectedFinalPhys, 0.01))
            Expect(damageAppliedEvent.FinalTotalDamage).To(BeNumerically("~", expectedTotalDamage, 0.01)) // This assertion seems redundant with the one above it
        })

        It("should consider damage amplification", func() {
            // --- Arrange ---
            attackerAttack.SetFinalDamageAmp(0.10) // 10% damage amp

            // --- Act ---
            mockEventBus.Enqueue(attackEvent, eventTime)
            mockEventBus.ProcessNext() // Process the event

            // --- Assert ---
            enqueuedEvents := mockEventBus.GetAllEvents()
            Expect(enqueuedEvents).To(HaveLen(1))
            damageAppliedEvent, ok := enqueuedEvents[0].(eventsys.DamageAppliedEvent)
            Expect(ok).To(BeTrue())


            // Manual Calculation:
            // Base Physical = 100
            // Crit EV Multiplier = 1.125
            // Amp Multiplier = 1.1
            // PreMitigation Physical = 100 * 1.125 * 1.1 = 123.75
            // Armor Reduction = 0.666...
            // Durability Multiplier = 1.0
            // Final Physical = 123.75 * 0.666... * 1.0 = 82.5
            expectedPreMitPhys := 123.75
            // expectedFinalPhys := 82.5 // This was missing in original, needed for TotalDamage check
            expectedTotalDamage := 82.5

            Expect(damageAppliedEvent.PreMitigationDamage).To(BeNumerically("~", expectedPreMitPhys, 0.01))
            Expect(damageAppliedEvent.FinalTotalDamage).To(BeNumerically("~", expectedTotalDamage, 0.01))
        })

        It("should consider durability", func() {
            // --- Arrange ---
            targetHealth.SetFinalDurability(0.20) // 20% durability

            // --- Act ---
            mockEventBus.Enqueue(attackEvent, eventTime)
            mockEventBus.ProcessNext() // Process the event

            // --- Assert ---
            enqueuedEvents := mockEventBus.GetAllEvents()
            Expect(enqueuedEvents).To(HaveLen(1))
            damageAppliedEvent, ok := enqueuedEvents[0].(eventsys.DamageAppliedEvent)
            Expect(ok).To(BeTrue())

            // Manual Calculation:
            // Base Physical = 100
            // Crit EV Multiplier = 1.125
            // Amp Multiplier = 1.0
            // PreMitigation Physical = 100 * 1.125 * 1.0 = 112.5
            // Armor Reduction = 0.666...
            // Durability Multiplier = (1 - 0.20) = 0.80
            // Final Physical = 112.5 * 0.666... * 0.80 = 75.0 * 0.8 = 60.0
            expectedPreMitPhys := 112.5
            // expectedFinalPhys := 60.0 // This was missing in original, needed for TotalDamage check
            expectedTotalDamage := 60.0

            Expect(damageAppliedEvent.PreMitigationDamage).To(BeNumerically("~", expectedPreMitPhys, 0.01))
            Expect(damageAppliedEvent.FinalTotalDamage).To(BeNumerically("~", expectedTotalDamage, 0.01))
        })
    })

	Describe("handling SpellCastEvent", func() {
        var (
            spellCastEvent eventsys.SpellLandedEvent
            eventTime      float64 = 2.50
        )

        BeforeEach(func() {
            spellCastEvent = eventsys.SpellLandedEvent{
                Source:    attacker,
                Target:    target,
                Timestamp: eventTime,
            }
            mockEventBus.ClearEvents()
            // Reset crit ability flags
            world.RemoveComponent(attacker, reflect.TypeOf(components.CanAbilityCritFromItems{}))
            world.RemoveComponent(attacker, reflect.TypeOf(components.CanAbilityCritFromTraits{}))
        })

        It("should calculate final magic damage considering MR", func() {
            // --- Act ---
            // damageSystem.HandleEvent(spellCastEvent) // Use direct handling if event bus isn't needed for this interaction
            // OR if using bus:
            mockEventBus.Enqueue(spellCastEvent, eventTime)
            mockEventBus.ProcessNext()

            // --- Assert ---
            enqueuedEvents := mockEventBus.GetAllEvents()
            Expect(enqueuedEvents).To(HaveLen(1))
            damageAppliedEvent, ok := enqueuedEvents[0].(eventsys.DamageAppliedEvent)
            Expect(ok).To(BeTrue())

            // ... rest of assertions ...
            expectedPreMitMag := 100.0
            expectedTotalDamage := 66.66 // Approx 100 * (100 / (100 + 50))

            Expect(damageAppliedEvent.Source).To(Equal(attacker))
            Expect(damageAppliedEvent.Target).To(Equal(target))
            Expect(damageAppliedEvent.Timestamp).To(Equal(eventTime))
            Expect(damageAppliedEvent.DamageSource).To(Equal("Spell"))
            Expect(damageAppliedEvent.DamageType).To(Equal("AP")) // Assuming simple AP spell for now
            Expect(damageAppliedEvent.PreMitigationDamage).To(BeNumerically("~", expectedPreMitMag, 0.01))
            Expect(damageAppliedEvent.FinalTotalDamage).To(BeNumerically("~", expectedTotalDamage, 0.01))
        })

        Context("when abilities can crit (e.g., Jeweled Gauntlet)", func() {
            BeforeEach(func() {
                // Add marker component to allow spell crit
                world.AddComponent(attacker, components.CanAbilityCritFromItems{})
            })

            It("should calculate final magic damage considering spell crit EV and MR", func() {
                 // --- Act ---
                // damageSystem.HandleEvent(spellCastEvent) // Or use bus
                mockEventBus.Enqueue(spellCastEvent, eventTime)
                mockEventBus.ProcessNext()

                // --- Assert ---
                enqueuedEvents := mockEventBus.GetAllEvents()
                Expect(enqueuedEvents).To(HaveLen(1))
                damageAppliedEvent, ok := enqueuedEvents[0].(eventsys.DamageAppliedEvent)
                Expect(ok).To(BeTrue())


                // Manual Calculation (Spell damage = 100 Magic):
                // Base Magic = 100
                // Crit Chance = 0.25, Crit Multiplier = 1.5 -> EV = 1.125
                // PreMitigation Magic = 100 * 1.125 = 112.5
                // Final Magic = 112.5 * (100 / (100 + 50)) = 112.5 * 0.666... = 75.0
                expectedPreMitMag := 112.5
                expectedTotalDamage := 75.0

                Expect(damageAppliedEvent.DamageSource).To(Equal("Spell"))
                Expect(damageAppliedEvent.PreMitigationDamage).To(BeNumerically("~", expectedPreMitMag, 0.01))
                Expect(damageAppliedEvent.DamageType).To(Equal("AP")) // Assuming simple AP spell
                Expect(damageAppliedEvent.FinalTotalDamage).To(BeNumerically("~", expectedTotalDamage, 0.01))
            })
        })
	})

    Describe("handling DamageAppliedEvent", func() {
        var (
            damageEvent eventsys.DamageAppliedEvent
            eventTime   float64 = 4.56
        )

        BeforeEach(func() {
            damageEvent = eventsys.DamageAppliedEvent{
                Source:              attacker,
                Target:              target,
                Timestamp:           eventTime,
                DamageType: 	  "AD",
				DamageSource:    "Attack",
				RawDamage: 	 50.0,
                PreMitigationDamage: 50.0,  
                FinalTotalDamage:    50.0, // Sum of physical and magic
            }
            mockEventBus.ClearEvents()
            attackerMana.SetCurrentMana(10)    // Reset mana
            targetHealth.SetCurrentHP(200) // Reset health
        })

        It("should decrease target health by FinalTotalDamage", func() {
            initialHP := targetHealth.CurrentHP
            damageSystem.HandleEvent(damageEvent)
            Expect(targetHealth.CurrentHP).To(Equal(initialHP - damageEvent.FinalTotalDamage))
        })

        Context("when damage is from an attack (IsSpell is false)", func() {
            BeforeEach(func() {
                damageEvent.DamageSource = "Attack"
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
        })

        Context("when damage is from a spell (IsSpell is true)", func() {
            BeforeEach(func() {
                damageEvent.DamageSource = "Spell"
            })

            It("should NOT increase attacker mana", func() {
                initialMana := attackerMana.GetCurrentMana()
                damageSystem.HandleEvent(damageEvent)
                Expect(attackerMana.GetCurrentMana()).To(Equal(initialMana), "Mana should not change from spell damage application")
            })
        })

        It("should not enqueue a DeathEvent if target survives", func() {
            targetHealth.SetCurrentHP(damageEvent.FinalTotalDamage + 1) // Ensure survival
            damageSystem.HandleEvent(damageEvent)
            // Expect no *new* events from this handler (Death/Kill are the only ones it sends)
            Expect(mockEventBus.GetAllEvents()).To(BeEmpty())
        })

        Context("when damage is lethal", func() {
            BeforeEach(func() {
                targetHealth.SetCurrentHP(damageEvent.FinalTotalDamage - 1) // Ensure lethal damage
                damageEvent.DamageSource = "Attack" // Ensure mana gain check runs
            })

            It("should set target health to zero or below", func() {
                damageSystem.HandleEvent(damageEvent)
                Expect(targetHealth.CurrentHP).To(BeNumerically("<=", 0))
            })

            It("should enqueue a DeathEvent and KillEvent", func() {
                damageSystem.HandleEvent(damageEvent)
                events := mockEventBus.GetAllEvents()
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
                Expect(killEvent.Timestamp).To(Equal(eventTime))
            })

            It("should still grant mana to the attacker if damage was from an attack", func() {
                initialMana := attackerMana.GetCurrentMana()
                expectedManaGain := 10.0
                damageSystem.HandleEvent(damageEvent)
                Expect(attackerMana.GetCurrentMana()).To(Equal(initialMana + expectedManaGain))
            })

            It("should only enqueue one DeathEvent and KillEvent even if called again", func() {
                damageSystem.HandleEvent(damageEvent) // First lethal hit
                Expect(mockEventBus.GetAllEvents()).To(HaveLen(2))

                // Simulate another damage event hitting the already dead target
                mockEventBus.ClearEvents()
                damageSystem.HandleEvent(damageEvent)
                Expect(mockEventBus.GetAllEvents()).To(BeEmpty(), "No additional Death/Kill events should be sent for an already dead target")
            })
        })

        // TODO: Add test for target gaining mana when hit
    })
})