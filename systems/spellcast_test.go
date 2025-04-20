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

var _ = Describe("SpellCastSystem", func() {
    var (
        world           *ecs.World
        eventBus        *MockEventBus
        championFactory *factory.ChampionFactory
        spellCastSystem *systems.SpellCastSystem
        player          ecs.Entity // Blue Golem
        target          ecs.Entity // Training Dummy
        playerMana      *components.Mana
        playerSpell     *components.Spell
        targetHealth    *components.Health
        ok              bool

        // Define spell properties for Blue Golem in this test
        spellName     string  = "BlueGolemSpell"
        spellManaCost float64 = 40.0
        spellMaxMana  float64 = 80.0
		castStartUp float64 = 0.0 // cast start up time
        castRecovery float64 = 1.5 // cast recovery used as lockout duration
    )

    BeforeEach(func() {
        world = ecs.NewWorld()
        eventBus = NewMockEventBus()
        championFactory = factory.NewChampionFactory(world)
        spellCastSystem = systems.NewSpellCastSystem(world, eventBus)

        // --- Create Player (Blue Golem) ---
        var err error
        player, err = championFactory.CreatePlayerChampion("TFT_BlueGolem", 1)
        Expect(err).NotTo(HaveOccurred())
        // Ensure Position exists
        if _, ok := world.GetPosition(player); !ok {
            world.AddComponent(player, components.NewPosition(0, 0))
        } else {
            pos, _ := world.GetPosition(player)
            pos.SetPosition(0, 0)
        }
        // Add/Replace Mana and Spell components with test values
        world.AddComponent(player, components.NewMana(0, spellMaxMana)) // Start with 0 mana, set later
        world.AddComponent(player, components.NewSpell(spellName, "", spellMaxMana, castStartUp, castRecovery))

        // --- Create Target (Training Dummy) ---
        target, err = championFactory.CreateEnemyChampion("TFT_TrainingDummy", 1)
        Expect(err).NotTo(HaveOccurred())
        // Ensure Position exists
        if _, ok := world.GetPosition(target); !ok {
            world.AddComponent(target, components.NewPosition(1, 0)) // Ensure target exists and has position
        } else {
            pos, _ := world.GetPosition(target)
            pos.SetPosition(1, 0) // Ensure in range
        }

        // --- Get Components ---
        playerMana, ok = world.GetMana(player)
        Expect(ok).To(BeTrue())
        Expect(playerMana).NotTo(BeNil())
        playerSpell, ok = world.GetSpell(player)
        Expect(ok).To(BeTrue())
        Expect(playerSpell).NotTo(BeNil())
        targetHealth, ok = world.GetHealth(target)
        Expect(ok).To(BeTrue())
        Expect(targetHealth).NotTo(BeNil())
        targetHealth.SetCurrentHP(1000) // Ensure target is alive

        // --- Initialize System Time ---
        spellCastSystem.SetCurrentTime(0.0)

        // --- Reset State ---
        playerMana.SetCurrentMana(spellManaCost) // Default to having enough mana
        playerSpell.SetCurrentRecovery(0)        // Default to being off recovery
    })

    Context("when player has insufficient mana", func() {
        BeforeEach(func() {
            playerMana.SetCurrentMana(spellManaCost - 1.0) // Set mana just below cost
        })

        It("should not enqueue a SpellCastEvent", func() {
            dt := 0.1
            spellCastSystem.SetCurrentTime(dt)
            spellCastSystem.TriggerSpellCast(dt) // Use TriggerSpellCast
            Expect(eventBus.EnqueuedEvents).To(BeEmpty())
        })
    })

    Context("when spell is on recovery", func() {
        var initialRecovery float64 = 1.0
        BeforeEach(func() {
            playerMana.SetCurrentMana(spellMaxMana) // Ensure enough mana
            playerSpell.SetCurrentRecovery(initialRecovery)
        })

        It("should reduce recovery over time", func() {
            dt := 0.1
            spellCastSystem.SetCurrentTime(dt)
            spellCastSystem.TriggerSpellCast(dt) // Use TriggerSpellCast
            Expect(playerSpell.GetCurrentRecovery()).To(BeNumerically("~", initialRecovery-dt, 0.001))
            Expect(eventBus.EnqueuedEvents).To(BeEmpty()) // Should not cast yet

            spellCastSystem.SetCurrentTime(dt * 2)
            spellCastSystem.TriggerSpellCast(dt) // dt is delta, not absolute time
            Expect(playerSpell.GetCurrentRecovery()).To(BeNumerically("~", initialRecovery-dt*2, 0.001))
            Expect(eventBus.EnqueuedEvents).To(BeEmpty())
        })

        It("should not enqueue a SpellCastEvent while recovery > 0", func() {
            dt := 0.1
            currentTime := 0.0
            // Simulate just before recovery ends
            for t := 0.0; t < initialRecovery-dt/2; t += dt {
                spellCastSystem.SetCurrentTime(currentTime + dt)
                spellCastSystem.TriggerSpellCast(dt) // Use TriggerSpellCast
                currentTime += dt
                Expect(eventBus.EnqueuedEvents).To(BeEmpty(), "Should not cast while recovery is active at time %.2f", currentTime)
                Expect(playerSpell.GetCurrentRecovery()).To(BeNumerically(">", 0))
            }
        })
    })

    Context("when player has mana and spell is off recovery", func() {
        BeforeEach(func() {
            playerMana.SetCurrentMana(spellMaxMana) // Ensure plenty of mana
            playerSpell.SetCurrentRecovery(0)
        })

        It("should enqueue a SpellCastEvent and update state", func() {
            dt := 0.1
            currentTime := 0.0

            spellCastSystem.SetCurrentTime(currentTime + dt)
            spellCastSystem.TriggerSpellCast(dt) // Use TriggerSpellCast
            currentTime += dt

            // Check Event
            Expect(eventBus.EnqueuedEvents).To(HaveLen(1))
            event, ok := eventBus.GetLastEvent().(eventsys.SpellCastEvent)
            Expect(ok).To(BeTrue())
            Expect(event.Source).To(Equal(player))
            Expect(event.Target).To(Equal(target)) // Assumes FindNearestEnemy works and finds the dummy
            Expect(event.Timestamp).To(BeNumerically("~", currentTime, 0.001))

            // Check State Update
            Expect(playerMana.GetCurrentMana()).To(BeNumerically("~", 0, 0.001), "Mana should be reduced by cost")
            Expect(playerSpell.GetCurrentRecovery()).To(BeNumerically("~", castRecovery, 0.001), "Recovery should be reset")
        })
    })

    Context("when no valid target exists", func() {
        BeforeEach(func() {
            playerMana.SetCurrentMana(spellMaxMana)
            playerSpell.SetCurrentRecovery(0)
            targetHealth.SetCurrentHP(0) // Make target dead
        })

        It("should still enqueue a SpellCastEvent", func() {
            dt := 0.1
            currentTime := 0.0
            spellCastSystem.SetCurrentTime(currentTime + dt)
            spellCastSystem.TriggerSpellCast(dt) // Use TriggerSpellCast
            currentTime += dt

            Expect(eventBus.EnqueuedEvents).To(HaveLen(1))
            // Check state was NOT updated
            Expect(playerMana.GetCurrentMana()).To(BeNumerically("~", 0, 0.001))
            Expect(playerSpell.GetCurrentRecovery()).To(BeNumerically("~", castRecovery, 0.001))
        })
    })

})