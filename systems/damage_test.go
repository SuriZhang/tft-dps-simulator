package systems_test

import (
	"reflect"

	"github.com/suriz/tft-dps-simulator/components"
	"github.com/suriz/tft-dps-simulator/ecs"
	"github.com/suriz/tft-dps-simulator/factory"
	"github.com/suriz/tft-dps-simulator/systems"
	eventsys "github.com/suriz/tft-dps-simulator/systems/events"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("DamageSystem", func() {
    var (
        world           *ecs.World
        eventBus        *MockEventBus
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
        eventBus = NewMockEventBus()
        championFactory = factory.NewChampionFactory(world) // Factory now adds Crit component
        damageSystem = systems.NewDamageSystem(world, eventBus)

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
        // TODO: Set spell variables if needed for specific spell tests
        // attackerSpell.SetVarBaseDamage(50)
        // attackerSpell.SetVarAPScaling(0.8)

        // Mana Stats
        attackerMana.SetCurrentMana(0)
        attackerMana.SetMaxMana(100)

        // Target Stats
        targetHealth.SetCurrentHealth(1000.0)
        targetHealth.SetFinalArmor(50.0)
        targetHealth.SetFinalMR(50.0)
        targetHealth.SetFinalDurability(0.0) // No initial durability
    })

    Describe("handling AttackLandedEvent", func() {
        var (
            attackEvent eventsys.AttackLandedEvent
            eventTime   float64 = 1.23
        )

        BeforeEach(func() {
            attackEvent = eventsys.AttackLandedEvent{
                Source:     attacker,
                Target:     target,
                BaseDamage: attackerAttack.GetFinalAD(),
                Timestamp:  eventTime,
            }
            eventBus.ClearEvents()
        })

        It("should calculate final physical damage considering crit EV and armor", func() {
            damageSystem.HandleEvent(attackEvent)

            Expect(eventBus.EnqueuedEvents).To(HaveLen(1), "Should enqueue one DamageAppliedEvent")
            damageAppliedEvent, ok := eventBus.GetLastEvent().(eventsys.DamageAppliedEvent)
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
            Expect(damageAppliedEvent.IsSpell).To(BeFalse())
            Expect(damageAppliedEvent.PreMitigationPhysical).To(BeNumerically("~", expectedPreMitPhys, 0.01))
            Expect(damageAppliedEvent.PreMitigationMagic).To(BeNumerically("~", 0.0, 0.01))
            Expect(damageAppliedEvent.FinalPhysicalDamage).To(BeNumerically("~", expectedFinalPhys, 0.01))
            Expect(damageAppliedEvent.FinalMagicDamage).To(BeNumerically("~", 0.0, 0.01))
            Expect(damageAppliedEvent.FinalTotalDamage).To(BeNumerically("~", expectedTotalDamage, 0.01))
        })

        It("should update the attacker's LastAttackTime", func() {
            initialLastAttackTime := attackerAttack.GetLastAttackTime()
            Expect(initialLastAttackTime).NotTo(Equal(eventTime))

            damageSystem.HandleEvent(attackEvent)

            Expect(attackerAttack.GetLastAttackTime()).To(Equal(eventTime))
        })

        It("should consider damage amplification", func() {
            attackerAttack.SetFinalDamageAmp(0.10) // 10% damage amp

            damageSystem.HandleEvent(attackEvent)

            Expect(eventBus.EnqueuedEvents).To(HaveLen(1))
            damageAppliedEvent, _ := eventBus.GetLastEvent().(eventsys.DamageAppliedEvent)

            // Manual Calculation:
            // Base Physical = 100
            // Crit EV Multiplier = 1.125
            // Amp Multiplier = 1.1
            // PreMitigation Physical = 100 * 1.125 * 1.1 = 123.75
            // Armor Reduction = 0.666...
            // Durability Multiplier = 1.0
            // Final Physical = 123.75 * 0.666... * 1.0 = 82.5
            expectedPreMitPhys := 123.75
            expectedFinalPhys := 82.5
            expectedTotalDamage := expectedFinalPhys

            Expect(damageAppliedEvent.PreMitigationPhysical).To(BeNumerically("~", expectedPreMitPhys, 0.01))
            Expect(damageAppliedEvent.FinalPhysicalDamage).To(BeNumerically("~", expectedFinalPhys, 0.01))
            Expect(damageAppliedEvent.FinalTotalDamage).To(BeNumerically("~", expectedTotalDamage, 0.01))
        })

        It("should consider durability", func() {
            targetHealth.SetFinalDurability(0.20) // 20% durability

            damageSystem.HandleEvent(attackEvent)

            Expect(eventBus.EnqueuedEvents).To(HaveLen(1))
            damageAppliedEvent, _ := eventBus.GetLastEvent().(eventsys.DamageAppliedEvent)

            // Manual Calculation:
            // Base Physical = 100
            // Crit EV Multiplier = 1.125
            // Amp Multiplier = 1.0
            // PreMitigation Physical = 100 * 1.125 * 1.0 = 112.5
            // Armor Reduction = 0.666...
            // Durability Multiplier = (1 - 0.20) = 0.80
            // Final Physical = 112.5 * 0.666... * 0.80 = 75.0 * 0.8 = 60.0
            expectedPreMitPhys := 112.5
            expectedFinalPhys := 60.0
            expectedTotalDamage := expectedFinalPhys

            Expect(damageAppliedEvent.PreMitigationPhysical).To(BeNumerically("~", expectedPreMitPhys, 0.01))
            Expect(damageAppliedEvent.FinalPhysicalDamage).To(BeNumerically("~", expectedFinalPhys, 0.01))
            Expect(damageAppliedEvent.FinalTotalDamage).To(BeNumerically("~", expectedTotalDamage, 0.01))
        })
    })

    Describe("handling SpellCastEvent", func() {
        var (
            spellCastEvent eventsys.SpellCastEvent
            eventTime      float64 = 2.50
        )

        BeforeEach(func() {
            // Setup spell variables for a hypothetical spell
            // Example: 50 base magic damage + 80% AP scaling + 20% AD scaling
            // attackerSpell.SetVarBaseDamage(50.0)
            // attackerSpell.SetVarAPScaling(0.8)
            // attackerSpell.SetVarPercentADDamage(0.2)
            // For simplicity, let's assume spell does only magic damage based on FinalAP for now
            // rawMagicDamage := attackerSpell.GetFinalAP() // 100
            // rawPhysicalDamage := 0.0

            spellCastEvent = eventsys.SpellCastEvent{
                Source:    attacker,
                Target:    target,
                Timestamp: eventTime,
            }
            eventBus.ClearEvents()
            // Ensure abilities cannot crit by default for this test
            world.RemoveComponent(attacker, reflect.TypeOf(components.CanAbilityCritFromItems{}))
            world.RemoveComponent(attacker, reflect.TypeOf(components.CanAbilityCritFromTraits{}))
        })

        It("should calculate final magic damage considering MR", func() {
            damageSystem.HandleEvent(spellCastEvent)

            Expect(eventBus.EnqueuedEvents).To(HaveLen(1))
            damageAppliedEvent, ok := eventBus.GetLastEvent().(eventsys.DamageAppliedEvent)
            Expect(ok).To(BeTrue())

            // Manual Calculation (Assuming spell damage = FinalAP = 100 Magic Damage):
            // Base Magic = 100 (from FinalAP, assuming simple spell)
            // Base Physical = 0
            // Can Abilities Crit = false -> Crit Multiplier EV = 1.0
            // Amp Multiplier = 1.0 (using attack amp for now, which is 0)
            // PreMitigation Magic = 100 * 1.0 * 1.0 = 100
            // PreMitigation Physical = 0 * 1.0 * 1.0 = 0
            // MR Reduction = 100 / (100 + 50) = 0.666...
            // Armor Reduction = 100 / (100 + 50) = 0.666... (not applied to magic)
            // Durability Multiplier = 1.0
            // Final Magic = 100 * 0.666... * 1.0 = 66.66...
            // Final Physical = 0
            // Final Total = 66.66...
            expectedPreMitMag := 100.0
            expectedFinalMag := 66.66
            expectedTotalDamage := expectedFinalMag

            Expect(damageAppliedEvent.Source).To(Equal(attacker))
            Expect(damageAppliedEvent.Target).To(Equal(target))
            Expect(damageAppliedEvent.Timestamp).To(Equal(eventTime))
            Expect(damageAppliedEvent.IsSpell).To(BeTrue())
            Expect(damageAppliedEvent.PreMitigationPhysical).To(BeNumerically("~", 0.0, 0.01))
            Expect(damageAppliedEvent.PreMitigationMagic).To(BeNumerically("~", expectedPreMitMag, 0.01))
            Expect(damageAppliedEvent.FinalPhysicalDamage).To(BeNumerically("~", 0.0, 0.01))
            Expect(damageAppliedEvent.FinalMagicDamage).To(BeNumerically("~", expectedFinalMag, 0.01))
            Expect(damageAppliedEvent.FinalTotalDamage).To(BeNumerically("~", expectedTotalDamage, 0.01))
        })

        Context("when abilities can crit (e.g., Jeweled Gauntlet)", func() {
            BeforeEach(func() {
                // Add marker component to allow spell crit
                world.AddComponent(attacker, components.CanAbilityCritFromItems{})
                // Use the same crit stats as attack: 25% chance, 1.5 multiplier
            })

            It("should calculate final magic damage considering spell crit EV and MR", func() {
                damageSystem.HandleEvent(spellCastEvent)

                Expect(eventBus.EnqueuedEvents).To(HaveLen(1))
                damageAppliedEvent, _ := eventBus.GetLastEvent().(eventsys.DamageAppliedEvent)

                // Manual Calculation (Spell damage = 100 Magic):
                // Base Magic = 100
                // Crit Chance = 0.25
                // Crit Multiplier = 1.5
                // Crit EV Multiplier = (1 - 0.25) + (0.25 * 1.5) = 1.125
                // Amp Multiplier = 1.0
                // PreMitigation Magic = 100 * 1.125 * 1.0 = 112.5
                // MR Reduction = 0.666...
                // Durability Multiplier = 1.0
                // Final Magic = 112.5 * 0.666... * 1.0 = 75.0
                expectedPreMitMag := 112.5
                expectedFinalMag := 75.0
                expectedTotalDamage := expectedFinalMag

                Expect(damageAppliedEvent.IsSpell).To(BeTrue())
                Expect(damageAppliedEvent.PreMitigationMagic).To(BeNumerically("~", expectedPreMitMag, 0.01))
                Expect(damageAppliedEvent.FinalMagicDamage).To(BeNumerically("~", expectedFinalMag, 0.01))
                Expect(damageAppliedEvent.FinalTotalDamage).To(BeNumerically("~", expectedTotalDamage, 0.01))
            })
        })

        // TODO: Add tests for spells that deal physical damage (AD scaling)
        // TODO: Add tests for spells dealing mixed damage
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
                IsSpell:             false, // Default to attack damage for these tests
                FinalPhysicalDamage: 50.0,  // Example damage breakdown
                FinalMagicDamage:    0.0,
                FinalTotalDamage:    50.0, // Sum of physical and magic
            }
            eventBus.ClearEvents()
            attackerMana.SetCurrentMana(10)    // Reset mana
            targetHealth.SetCurrentHealth(200) // Reset health
        })

        It("should decrease target health by FinalTotalDamage", func() {
            initialHP := targetHealth.CurrentHP
            damageSystem.HandleEvent(damageEvent)
            Expect(targetHealth.CurrentHP).To(Equal(initialHP - damageEvent.FinalTotalDamage))
        })

        Context("when damage is from an attack (IsSpell is false)", func() {
            BeforeEach(func() {
                damageEvent.IsSpell = false
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
                damageEvent.IsSpell = true
            })

            It("should NOT increase attacker mana", func() {
                initialMana := attackerMana.GetCurrentMana()
                damageSystem.HandleEvent(damageEvent)
                Expect(attackerMana.GetCurrentMana()).To(Equal(initialMana), "Mana should not change from spell damage application")
            })
        })

        It("should not enqueue a DeathEvent if target survives", func() {
            targetHealth.SetCurrentHealth(damageEvent.FinalTotalDamage + 1) // Ensure survival
            damageSystem.HandleEvent(damageEvent)
            // Expect no *new* events from this handler (Death/Kill are the only ones it sends)
            Expect(eventBus.EnqueuedEvents).To(BeEmpty())
        })

        Context("when damage is lethal", func() {
            BeforeEach(func() {
                targetHealth.SetCurrentHealth(damageEvent.FinalTotalDamage - 1) // Ensure lethal damage
                damageEvent.IsSpell = false // Ensure mana gain check runs
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
                Expect(eventBus.GetAllEvents()).To(HaveLen(2))

                // Simulate another damage event hitting the already dead target
                eventBus.ClearEvents()
                damageSystem.HandleEvent(damageEvent)
                Expect(eventBus.EnqueuedEvents).To(BeEmpty(), "No additional Death/Kill events should be sent for an already dead target")
            })
        })

        // TODO: Add test for target gaining mana when hit
    })
})