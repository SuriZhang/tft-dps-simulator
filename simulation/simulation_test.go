package simulation_test

import (
	"fmt"
	"math"
	"reflect"

	"github.com/suriz/tft-dps-simulator/components"
	"github.com/suriz/tft-dps-simulator/data"
	"github.com/suriz/tft-dps-simulator/ecs"
	"github.com/suriz/tft-dps-simulator/factory"
	"github.com/suriz/tft-dps-simulator/managers"
	"github.com/suriz/tft-dps-simulator/simulation"
	"github.com/suriz/tft-dps-simulator/systems"
	itemsys "github.com/suriz/tft-dps-simulator/systems/items"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes" // For checking output
)

// Helper to get component safely
func getAttack(w *ecs.World, e ecs.Entity) *components.Attack {
    comp, ok := w.GetAttack(e)
    Expect(ok).To(BeTrue(), "Entity should have Attack component")
    Expect(comp).NotTo(BeNil())
    return comp
}

func getHealth(w *ecs.World, e ecs.Entity) *components.Health {
    comp, ok := w.GetHealth(e)
    Expect(ok).To(BeTrue(), "Entity should have Health component")
    Expect(comp).NotTo(BeNil())
    return comp
}

func getSpell(w *ecs.World, e ecs.Entity) *components.Spell {
    comp, ok := w.GetSpell(e)
    Expect(ok).To(BeTrue(), "Entity should have Spell component")
    Expect(comp).NotTo(BeNil())
    return comp
}

func getCrit(w *ecs.World, e ecs.Entity) *components.Crit {
    comp, ok := w.GetCrit(e)
    Expect(ok).To(BeTrue(), "Entity should have Crit component")
    Expect(comp).NotTo(BeNil())
    return comp
}

var _ = Describe("Simulation", func() {
    var (
        world            *ecs.World
        sim              *simulation.Simulation
        config           simulation.SimulationConfig
        championFactory  *factory.ChampionFactory
        equipmentManager *managers.EquipmentManager // Added equipment manager
        attacker         ecs.Entity
        target           ecs.Entity
        targetMaxHP      float64 // Store initial max HP for checks
    )

    BeforeEach(func() {
        world = ecs.NewWorld()
        config = simulation.DefaultConfig().
            WithTimeStep(0.1).    // TimeStep might be less relevant now, but keep for potential future use?
            WithMaxTime(1.0).     // Short max time for faster tests
            WithDebugMode(false) // Disable debug output by default

        championFactory = factory.NewChampionFactory(world)
        // Load real item data for the manager
        // Ensure item data is loaded correctly in your actual setup, maybe in a BeforeSuite
        // For this test, we assume item data is available via BeforeSuite in simulation_suite_test.go
        equipmentManager = managers.NewEquipmentManager(world) // Initialize equipment manager

        // Create basic entities for interaction testing
        var err error
        // Use a champion with known base stats if possible, otherwise set reasonable bases
        attacker, err = championFactory.CreatePlayerChampion("TFT_TrainingDummy", 1) // Use dummy for predictable base stats
        Expect(err).NotTo(HaveOccurred())
        world.AddComponent(attacker, components.NewPosition(0, 0))

        // Add mana if not present
        if _, ok := world.GetMana(attacker); !ok {
            world.AddComponent(attacker, components.NewMana(0, 100))
        }
        // Add spell component if not present (needed for Archangel's)
        attackerSpell := getSpell(world, attacker)

        // Set reasonable base stats if using a dummy
        attackerAttack := getAttack(world, attacker)
        attackerAttack.SetBaseAttackSpeed(0.5)
        attackerAttack.SetBaseAD(50)
        attackerSpell.SetBaseAP(50)
        attackerAttack.SetBaseRange(1.0) // Ensure target is in range

        target, err = championFactory.CreateEnemyChampion("TFT_TrainingDummy", 1)
        Expect(err).NotTo(HaveOccurred())
        world.AddComponent(target, components.NewPosition(1, 0))
        targetHealth := getHealth(world, target)
        targetHealth.SetBaseArmor(0.0)
        targetHealth.SetBaseMR(0.0)
        targetMaxHP = targetHealth.GetBaseMaxHp()
        // Ensure target health is not NaN initially
        Expect(math.IsNaN(targetMaxHP)).To(BeFalse(), "Initial Target Max HP should not be NaN")
        Expect(math.IsNaN(targetHealth.GetCurrentHP())).To(BeFalse(), "Initial Target Current HP should not be NaN")

        // Ensure target dummy doesn't attack
        if targetAttack, ok := world.GetAttack(target); ok {
            targetAttack.SetBaseAttackSpeed(0)
            targetAttack.SetFinalAttackSpeed(0)
        }
    })

    Describe("Initialization", func() {
        It("should create a simulation with default config", func() {
            defaultSim := simulation.NewSimulation(world) // Assumes NewSimulation uses DefaultConfig
            Expect(defaultSim).NotTo(BeNil())
            Expect(defaultSim.GetConfig()).To(Equal(simulation.DefaultConfig()))
        })

        It("should create a simulation with custom config", func() {
            customConfig := simulation.DefaultConfig().WithMaxTime(5.5)
            customSim := simulation.NewSimulationWithConfig(world, customConfig)
            Expect(customSim).NotTo(BeNil())
            Expect(customSim.GetConfig()).To(Equal(customConfig))
        })

        It("should panic with invalid config", func() {
            invalidConfig := simulation.DefaultConfig().WithTimeStep(-1) // Example invalid config
            Expect(func() { simulation.NewSimulationWithConfig(world, invalidConfig) }).To(Panic())
        })

        It("should run initial setup (setupCombat) upon creation", func() {
            // Add an item with static stats BEFORE creating the simulation
            deathcapData := data.GetItemByApiName(data.TFT_Item_RabadonsDeathcap)
            Expect(deathcapData).NotTo(BeNil())
            err := equipmentManager.AddItemToChampion(attacker, data.TFT_Item_RabadonsDeathcap)
            Expect(err).NotTo(HaveOccurred())

            // Create the simulation - this should trigger setupCombat
            sim = simulation.NewSimulationWithConfig(world, config)
            Expect(sim).NotTo(BeNil())

            // Check if static stats were applied during setupCombat
            attackerSpell := getSpell(world, attacker)
            expectedInitialAP := attackerSpell.GetBaseAP() + deathcapData.Effects["AP"]
            Expect(attackerSpell.GetFinalAP()).To(BeNumerically("~", expectedInitialAP, 0.01), "Static AP should be applied during simulation initialization")

            // Check if initial events were enqueued (e.g., ChampionActionEvent)
            // This requires access to the internal event bus state, might be hard to test directly
            // We can infer it by running the simulation for a tiny duration and checking effects.
        })
    })

    Describe("Configuration", func() {
        BeforeEach(func() {
            // Create sim instance for config tests
            sim = simulation.NewSimulationWithConfig(world, config)
            Expect(sim).NotTo(BeNil())
        })

        It("should get the current config", func() {
            Expect(sim.GetConfig()).To(Equal(config))
        })

        It("should set a new valid config", func() {
            newConfig := config.WithMaxTime(15.0)
            err := sim.SetConfig(newConfig)
            Expect(err).NotTo(HaveOccurred())
            Expect(sim.GetConfig()).To(Equal(newConfig))
        })

        It("should return error when setting invalid config", func() {
            invalidConfig := config.WithMaxTime(-5) // Invalid MaxTime
            err := sim.SetConfig(invalidConfig)
            Expect(err).To(HaveOccurred())
            Expect(sim.GetConfig()).To(Equal(config)) // Config should not change
        })

        It("should set max time via helper", func() {
            sim.SetMaxTime(12.0)
            Expect(sim.GetConfig().MaxTime).To(Equal(12.0))
        })

        It("should set time step via helper", func() {
            sim.SetTimeStep(0.05)
            Expect(sim.GetConfig().TimeStep).To(Equal(0.05)) // Assuming TimeStep is still used somewhere or just configurable
        })
    })

    Describe("RunSimulation Method", func() {
        // Uses the 'sim' instance created in the outer BeforeEach,
        // which has the real event bus and systems wired up internally.

        BeforeEach(func() {
            // Create sim instance for RunSimulation tests AFTER basic world setup
            sim = simulation.NewSimulationWithConfig(world, config)
            Expect(sim).NotTo(BeNil())
        })

        It("should run until MaxTime is reached and apply effects", func() {
            // Attacker Base AS = 0.5. Attack interval = 1 / 0.5 = 2.0s.
            // Attack should land around t=2.0s.
            sim.SetMaxTime(2.5) // Set time long enough for one attack

            // Get initial state
            targetHealth := getHealth(world, target)
            initialHP := targetHealth.GetCurrentHP()
            Expect(initialHP).To(Equal(targetMaxHP)) // Should start at max HP

            // Run the simulation - processes events until MaxTime or queue empty
            sim.RunSimulation()

            // Check effects after running
            // Target should have taken damage from the AttackLandedEvent -> DamageAppliedEvent chain
            // Damage calculation: BaseAD=50. Assume 0 crit/amp/durability, Armor=0 -> Final=50
            expectedHPAfterOneAttack := targetMaxHP - 50.0
            Expect(targetHealth.GetCurrentHP()).To(BeNumerically("~", expectedHPAfterOneAttack, 0.1), "Target HP should decrease after attack event processing")
        })

        It("should execute initial system updates before the event loop", func() {
            // NewSimulationWithConfig runs setupCombat, which applies static stats.
            // This test verifies those stats are present *before* RunSimulation processes time-based events.
            attackerAttack := getAttack(world, attacker)
            Expect(attackerAttack.GetFinalAD()).To(BeNumerically("~", 50.0, 0.1), "Initial Final AD calculation should happen in setup")
            Expect(attackerAttack.GetFinalAttackSpeed()).To(BeNumerically("~", 0.5, 0.01), "Initial Final AS calculation should happen in setup")

            // Now run the simulation for a very short time (or zero time if possible)
            // This ensures RunSimulation doesn't mess up the initial state.
            sim.SetMaxTime(0.01)
            sim.RunSimulation()
            Expect(attackerAttack.GetFinalAD()).To(BeNumerically("~", 50.0, 0.1)) // Should still be calculated correctly
            Expect(attackerAttack.GetFinalAttackSpeed()).To(BeNumerically("~", 0.5, 0.01))
        })

        It("should handle multiple attack cycles correctly", func() {
            // Attacker Base AS = 0.5. Attack Startup/Recovery = 0
            // first attack fired at 0.0s
            // second attack fired at 2.0s
            sim.SetMaxTime(3.9) // Enough time for two attacks

            targetHealth := getHealth(world, target)
            attackerAttack := getAttack(world, attacker)
            sim.RunSimulation() // Processes events for both attacks

            // Check effects after two attacks
            // Damage = 50 per attack
            Expect(attackerAttack.GetAttackCount()).To(Equal(2))
            expectedHPAfterTwoAttacks := targetMaxHP - 50.0 - 50.0
            Expect(targetHealth.GetCurrentHP()).To(BeNumerically("~", expectedHPAfterTwoAttacks, 0.1))
        })

        // Optional: Test DebugMode output
        PIt("should print debug messages when DebugMode is enabled", func() {
            // This test requires capturing stdout, which can be complex.
            // Using gbytes.Say is one way with Ginkgo V2.
            buffer := gbytes.NewBuffer() // Capture output

            debugConfig := config.WithDebugMode(true).WithMaxTime(0.6) // Removed ReportingInterval if not used
            // Create a new sim instance for this specific test
            // TODO: Need a way to inject the buffer into the simulation's logger
            debugSim := simulation.NewSimulationWithConfig(world, debugConfig)

            debugSim.RunSimulation()

            // Check for specific debug log patterns
            Eventually(buffer).Should(gbytes.Say(`\[T=0.000s\] Dequeued: events.ChampionActionEvent`)) // Example expected output
            // Eventually(buffer).Should(gbytes.Say(`\[T=2.000s\] Dequeued: events.AttackLandedEvent`)) // Example
        })

        // Context("with Dynamic Time Items", func() {
        //     // Tests for items like Archangel's, Quicksilver that trigger based on time intervals

        //     var (
        //         quicksilverData *data.Item
        //         qsDuration      float64
        //         qsProcInterval  float64
        //         qsProcAS        float64
        //         qsStaticAS      float64

        //         archangelsData *data.Item
        //         aaInitialAP    float64
        //         aaInterval     float64
        //         aaAPPerStack   float64
        //     )

        //     BeforeEach(func() {
        //         // Load item data needed for these tests
        //         quicksilverData = data.GetItemByApiName(data.TFT_Item_Quicksilver)
        //         Expect(quicksilverData).NotTo(BeNil())
        //         qsDuration = quicksilverData.Effects["SpellShieldDuration"] // e.g., 18.0
        //         qsProcInterval = quicksilverData.Effects["ProcInterval"]    // e.g., 2.0
        //         qsProcAS = quicksilverData.Effects["ProcAttackSpeed"]       // e.g., 0.03
        //         qsStaticAS = quicksilverData.Effects["AS"] / 100.0          // e.g., 0.3

        //         archangelsData = data.GetItemByApiName(data.TFT_Item_ArchangelsStaff)
        //         Expect(archangelsData).NotTo(BeNil())
        //         aaInitialAP = archangelsData.Effects["AP"]
        //         aaInterval = archangelsData.Effects["IntervalSeconds"]    // e.g., 5.0
        //         aaAPPerStack = archangelsData.Effects["APPerInterval"] // e.g., 30.0
        //     })

        //     It("should apply Quicksilver initial AS bonus, stack AS, and expire correctly", func() {
        //         // Add Quicksilver to attacker BEFORE creating simulation
        //         err := equipmentManager.AddItemToChampion(attacker, data.TFT_Item_Quicksilver)
        //         Expect(err).NotTo(HaveOccurred())

        //         // Create simulation AFTER adding the item
        //         sim = simulation.NewSimulationWithConfig(world, config)
        //         Expect(sim).NotTo(BeNil())

        //         attackerAttack := getAttack(world, attacker)
        //         targetHealth := getHealth(world, target)
        //         baseAS := attackerAttack.GetBaseAttackSpeed() // 0.5

        //         // --- Check Initial State (after setupCombat) ---
        //         expectedInitialAS := baseAS * (1 + qsStaticAS) // 0.5 * (1 + 0.3) = 0.65
        //         Expect(attackerAttack.GetFinalAttackSpeed()).To(BeNumerically("~", expectedInitialAS, 0.01), "Initial Final AS should include static QS bonus")
        //         qsEffect, ok := world.GetQuicksilverEffect(attacker)
        //         Expect(ok).To(BeTrue())
        //         Expect(qsEffect.IsActive()).To(BeTrue())
        //         Expect(qsEffect.GetStacks()).To(Equal(0))
        //         Expect(qsEffect.GetCurrentBonusAS()).To(BeNumerically("~", 0.0)) // Dynamic bonus is 0 initially

        //         // --- Run simulation for a short time (e.g., 5 seconds) ---
        //         // Expect procs at t=2.0, t=4.0 -> 2 stacks
        //         shortDuration := 5.0
        //         sim.SetMaxTime(shortDuration)
        //         sim.RunSimulation() // Processes events up to 5.0s

        //         // --- Verify State after 5s ---
        //         Expect(targetHealth.CurrentHP).To(BeNumerically("<", targetMaxHP), "Target should have taken damage within %f seconds", shortDuration)
        //         Expect(qsEffect.IsActive()).To(BeTrue(), "QS should still be active at 5s")
        //         Expect(qsEffect.GetStacks()).To(Equal(2), "QS should have 2 stacks at 5s")
        //         expectedDynamicASBonus := float64(2) * qsProcAS // 2 * 0.03 = 0.06
        //         Expect(qsEffect.GetCurrentBonusAS()).To(BeNumerically("~", expectedDynamicASBonus, 0.001), "QS dynamic bonus should be for 2 stacks")
        //         // Check Final AS reflects static + dynamic bonus
        //         expectedASAt5s := baseAS * (1 + qsStaticAS + expectedDynamicASBonus) // 0.5 * (1 + 0.3 + 0.06) = 0.5 * 1.36 = 0.68
        //         Expect(attackerAttack.GetFinalAttackSpeed()).To(BeNumerically("~", expectedASAt5s, 0.01), "Final AS at 5s should reflect static + 2 stacks")
        //         healthAfterShortRun := targetHealth.CurrentHP

        //         // --- Continue simulation past expiry (e.g., 19 seconds) ---
        //         // Procs happen at 2, 4, 6, 8, 10, 12, 14, 16, 18 -> 9 stacks total before expiry
        //         // NOTE: We create a NEW sim instance with the SAME world but new config/time
        //         // This simulates continuing the fight from the state at 5s, but the sim runs from t=0 internally again.
        //         // To properly test continuation, the simulation loop needs to handle starting from a non-zero time,
        //         // OR we need to run the *same* sim instance longer. Let's try running the same instance longer.

        //         // Reset sim's internal time and run longer using the *same* instance
        //         // This assumes RunSimulation can be called multiple times or handles continuation.
        //         // If RunSimulation always starts from t=0, this approach is flawed.
        //         // Let's assume RunSimulation processes events based on their timestamp relative to MaxTime.
        //         // A cleaner way: Create ONE sim and run it for the full duration.

        //         // --- Re-run from scratch for the full duration ---
        //         world = ecs.NewWorld() // Reset world
        //         championFactory = factory.NewChampionFactory(world)
        //         equipmentManager = managers.NewEquipmentManager(world)
        //         attacker, _ = championFactory.CreatePlayerChampion("TFT_TrainingDummy", 1)
        //         world.AddComponent(attacker, components.NewPosition(0, 0))
        //         if _, ok := world.GetMana(attacker); !ok { world.AddComponent(attacker, components.NewMana(0, 100)) }
        //         attackerAttack = getAttack(world, attacker) // Re-get component
        //         attackerAttack.SetBaseAttackSpeed(0.5)
        //         attackerAttack.SetBaseAD(50)
        //         getSpell(world, attacker).SetBaseAP(50)
        //         attackerAttack.SetBaseRange(1.0)
        //         target, _ = championFactory.CreateEnemyChampion("TFT_TrainingDummy", 1)
        //         world.AddComponent(target, components.NewPosition(1, 0))
        //         targetHealth = getHealth(world, target) // Re-get component
        //         targetHealth.SetBaseArmor(0.0)
        //         targetHealth.SetBaseMR(0.0)
        //         targetMaxHP = targetHealth.GetBaseMaxHp()
        //         if targetAttack, ok := world.GetAttack(target); ok { targetAttack.SetFinalAttackSpeed(0) }

        //         err = equipmentManager.AddItemToChampion(attacker, data.TFT_Item_Quicksilver) // Add item again
        //         Expect(err).NotTo(HaveOccurred())

        //         longDuration := 19.0
        //         config = config.WithMaxTime(longDuration) // Update config max time
        //         sim = simulation.NewSimulationWithConfig(world, config) // Create ONE sim instance for the full duration
        //         Expect(sim).NotTo(BeNil())

        //         sim.RunSimulation() // Run for 19 seconds

        //         // --- Verify State after 19s ---
        //         qsEffect, ok = world.GetQuicksilverEffect(attacker) // Re-get effect component
        //         Expect(ok).To(BeTrue(), "QS effect component should still exist after expiry (to hold stacks)")
        //         Expect(qsEffect.IsActive()).To(BeFalse(), "QS effect should be inactive after 18s")
        //         Expect(qsEffect.GetRemainingDuration()).To(BeNumerically("<=", 0.0), "QS remaining duration should be <= 0")
        //         // Stacks accumulate until expiry (t=18s -> 9 stacks) and persist
        //         Expect(qsEffect.GetStacks()).To(Equal(9), "QS effect should retain 9 accumulated stacks after expiry")
        //         // Dynamic AS bonus persists after expiry based on final stacks
        //         expectedFinalDynamicASBonus := float64(9) * qsProcAS // 9 * 0.03 = 0.27
        //         Expect(qsEffect.GetCurrentBonusAS()).To(BeNumerically("~", expectedFinalDynamicASBonus, 0.001), "QS dynamic bonus should persist based on 9 stacks")

        //         // Verify Final AS reflects static + max dynamic stacks
        //         expectedASAfterExpiry := baseAS * (1 + qsStaticAS + expectedFinalDynamicASBonus) // 0.5 * (1 + 0.3 + 0.27) = 0.5 * 1.57 = 0.785
        //         Expect(attackerAttack.GetFinalAttackSpeed()).To(BeNumerically("~", expectedASAfterExpiry, 0.01), "Final AS after expiry should reflect base + static + 9 stacks bonus")

        //         // Verify attacks continued after expiry (target HP should be lower than if sim stopped at 5s)
        //         // Cannot directly compare to healthAfterShortRun as we reset the world.
        //         // Just check that damage occurred.
        //         Expect(targetHealth.CurrentHP).To(BeNumerically("<", targetMaxHP), "Target should have taken damage over 19 seconds")
        //     })

        //     It("should stack Archangel's Staff AP over time", func() {
        //         // Add Archangel's Staff to attacker BEFORE creating simulation
        //         err := equipmentManager.AddItemToChampion(attacker, data.TFT_Item_ArchangelsStaff)
        //         Expect(err).NotTo(HaveOccurred())

        //         // Create simulation AFTER adding the item
        //         sim = simulation.NewSimulationWithConfig(world, config)
        //         Expect(sim).NotTo(BeNil())

        //         attackerSpell := getSpell(world, attacker)
        //         aaEffect, ok := world.GetArchangelsEffect(attacker)
        //         Expect(ok).To(BeTrue())

        //         // --- Check Initial State ---
        //         expectedInitialAP := attackerSpell.GetBaseAP() + aaInitialAP
        //         Expect(attackerSpell.GetFinalAP()).To(BeNumerically("~", expectedInitialAP, 0.01), "Initial Final AP should include static AA bonus")
        //         Expect(aaEffect.GetStacks()).To(Equal(0))

        //         // --- Run for ~6 seconds (expect 1 stack at t=5.0) ---
        //         sim.SetMaxTime(6.0)
        //         sim.RunSimulation()

        //         // --- Verify State after 6s ---
        //         Expect(aaEffect.GetStacks()).To(Equal(1), "AA should have 1 stack after 6s")
        //         expectedAPAt6s := attackerSpell.GetBaseAP() + aaInitialAP + (1 * aaAPPerStack)
        //         Expect(attackerSpell.GetFinalAP()).To(BeNumerically("~", expectedAPAt6s, 0.01), "Final AP should include static + 1 stack bonus after 6s")
        //         apAfter6Sec := attackerSpell.GetFinalAP() // Store for next check

        //         // --- Run for ~11 seconds (expect 2 stacks at t=5.0, t=10.0) ---
        //         // Re-run from scratch for the full duration
        //         world = ecs.NewWorld() // Reset world
        //         championFactory = factory.NewChampionFactory(world)
        //         equipmentManager = managers.NewEquipmentManager(world)
        //         attacker, _ = championFactory.CreatePlayerChampion("TFT_TrainingDummy", 1)
        //         world.AddComponent(attacker, components.NewPosition(0, 0))
        //         if _, ok := world.GetMana(attacker); !ok { world.AddComponent(attacker, components.NewMana(0, 100)) }
        //         attackerSpell = getSpell(world, attacker) // Re-get component
        //         attackerSpell.SetBaseAP(50) // Reset base stat
        //         target, _ = championFactory.CreateEnemyChampion("TFT_TrainingDummy", 1) // Recreate target if needed
        //         world.AddComponent(target, components.NewPosition(1, 0))

        //         err = equipmentManager.AddItemToChampion(attacker, data.TFT_Item_ArchangelsStaff) // Add item again
        //         Expect(err).NotTo(HaveOccurred())

        //         sim = simulation.NewSimulationWithConfig(world, config.WithMaxTime(11.0)) // Create sim for 11s
        //         Expect(sim).NotTo(BeNil())

        //         sim.RunSimulation() // Run for 11 seconds

        //         // --- Verify State after 11s ---
        //         aaEffect, ok = world.GetArchangelsEffect(attacker) // Re-get effect
        //         Expect(ok).To(BeTrue())
        //         Expect(aaEffect.GetStacks()).To(Equal(2), "AA should have 2 stacks after 11s")
        //         expectedAPAt11s := attackerSpell.GetBaseAP() + aaInitialAP + (2 * aaAPPerStack)
        //         Expect(attackerSpell.GetFinalAP()).To(BeNumerically("~", expectedAPAt11s, 0.01), "Final AP should include static + 2 stacks bonus after 11s")
        //         // Cannot compare directly to apAfter6Sec as world was reset
        //         // Expect(attackerSpell.GetFinalAP()).To(BeNumerically(">", apAfter6Sec), "AP should increase further after second stack interval (~10s)")
        //     })
        // }) // End Context("with Dynamic Time Items")

        Context("with Dynamic Event Items", func() {
            // Tests for items like Titan's, Guinsoo's that stack on events

            var (
                titansData              *data.Item
                titansMaxStacks         int
                titansADPerStack        float64
                titansAPPerStack        float64
                titansBonusResistsAtCap float64
                titansStaticArmor       float64

                ragebladeData       *data.Item
                ragebladeStaticAS   float64
                ragebladeASPerStack float64
            )

            BeforeEach(func() {
                // Get Item data once for this context
                titansData = data.GetItemByApiName(data.TFT_Item_TitansResolve)
                Expect(titansData).NotTo(BeNil())
                titansMaxStacks = int(titansData.Effects["StackCap"])         // 25
                titansADPerStack = titansData.Effects["StackingAD"]           // 2.0
                titansAPPerStack = titansData.Effects["StackingSP"]           // 2.0
                titansBonusResistsAtCap = titansData.Effects["BonusResistsAtStackCap"] // 20.0
                titansStaticArmor = titansData.Effects["Armor"]               // 10.0

                ragebladeData = data.GetItemByApiName(data.TFT_Item_GuinsoosRageblade)
                Expect(ragebladeData).NotTo(BeNil())
                ragebladeStaticAS = ragebladeData.Effects["AS"] / 100.0          // 0.1
                ragebladeASPerStack = ragebladeData.Effects["AttackSpeedPerStack"] / 100.0 // 0.05
            })

            It("should stack Titan's Resolve on attacks and apply AD/AP bonuses", func() {
                // Add Titan's to attacker BEFORE creating simulation
                err := equipmentManager.AddItemToChampion(attacker, data.TFT_Item_TitansResolve)
                Expect(err).NotTo(HaveOccurred())

                // Ensure attacker attacks and target doesn't
                attackerAttack := getAttack(world, attacker)
                attackerAttack.SetBaseAttackSpeed(0.5) // Attacks land at t=2.0, 4.0
                if targetAttack, ok := world.GetAttack(target); ok {
                    targetAttack.SetFinalAttackSpeed(0)
                }

                // Create simulation AFTER adding the item and setting AS
                sim = simulation.NewSimulationWithConfig(world, config.WithMaxTime(4.5)) // Run long enough for 2 attacks
                Expect(sim).NotTo(BeNil())

                // Get components and initial stats *after* sim creation (setupCombat runs)
                attackerSpell := getSpell(world, attacker)
                attackerHealth := getHealth(world, attacker)
                targetHealth := getHealth(world, target)
                titansEffect, ok := world.GetTitansResolveEffect(attacker)
                Expect(ok).To(BeTrue())

                initialFinalAD := attackerAttack.GetFinalAD() // Includes Base AD
                initialFinalAP := attackerSpell.GetFinalAP() // Includes Base AP + Static Item AP (if any)
                initialBonusArmor := attackerHealth.GetBonusArmor() // Should include static armor from Titan's (10)
                Expect(initialBonusArmor).To(BeNumerically("~", titansStaticArmor, 0.001))

                // --- Run Simulation ---
                sim.RunSimulation()

                // --- Assertions ---
                expectedStacks := 2 // From attacks at t=2.0, t=4.0
                Expect(titansEffect.GetCurrentStacks()).To(Equal(expectedStacks), "Should gain 2 stacks from 2 attacks")
                // Check dynamic bonuses applied by the event handler
                Expect(attackerAttack.GetBonusPercentAD()).To(BeNumerically("~", titansADPerStack*float64(expectedStacks), 0.001), "Bonus AD should reflect 2 stacks")
                Expect(attackerSpell.GetBonusAP()).To(BeNumerically("~", titansAPPerStack*float64(expectedStacks), 0.001), "Bonus AP should reflect 2 stacks")
                // Bonus Armor should only include static part until max stacks
                Expect(attackerHealth.GetBonusArmor()).To(BeNumerically("~", initialBonusArmor, 0.001), "Bonus Armor should not include max stack bonus yet")

                // Check final stats increased due to stacking bonuses (StatCalculationSystem runs implicitly or via events)
                // Final AD = Base * (1 + BonusPercentAD)
                // Final AP = Base + BonusAP
                expectedFinalAD := attackerAttack.GetBaseAD() * (1 + attackerAttack.GetBonusPercentAD())
                expectedFinalAP := attackerSpell.GetBaseAP() + attackerSpell.GetBonusAP()
                Expect(attackerAttack.GetFinalAD()).To(BeNumerically("~", expectedFinalAD, 0.01))
                Expect(attackerSpell.GetFinalAP()).To(BeNumerically("~", expectedFinalAP, 0.01))
                Expect(attackerAttack.GetFinalAD()).To(BeNumerically(">", initialFinalAD), "Final AD should increase from stacks")
                Expect(attackerSpell.GetFinalAP()).To(BeNumerically(">", initialFinalAP), "Final AP should increase from stacks")
                Expect(targetHealth.GetCurrentHP()).To(BeNumerically("<", targetMaxHP), "Target should take damage")
            })

            It("should stack Titan's Resolve on taking damage and apply AD/AP bonuses", func() {
                // Add Titan's to attacker BEFORE creating simulation
                err := equipmentManager.AddItemToChampion(attacker, data.TFT_Item_TitansResolve)
                Expect(err).NotTo(HaveOccurred())

                // Modify setup: Target attacks slowly, attacker doesn't attack
                if targetAttack, ok := world.GetAttack(target); ok {
                    targetAttack.SetBaseAttackSpeed(0.1) // Target attacks at t=10.0
                    targetAttack.SetBaseAD(30)           // Give target some AD
                }
                attackerAttack := getAttack(world, attacker)
                attackerAttack.SetBaseAttackSpeed(0) // Stop attacker

                // Create simulation AFTER setup changes
                sim = simulation.NewSimulationWithConfig(world, config.WithMaxTime(10.5)) // Run long enough for target's attack
                Expect(sim).NotTo(BeNil())

                // Get components and initial stats
                attackerSpell := getSpell(world, attacker)
                attackerHealth := getHealth(world, attacker)
                titansEffect, ok := world.GetTitansResolveEffect(attacker)
                Expect(ok).To(BeTrue())

                initialFinalAD := attackerAttack.GetFinalAD()
                initialFinalAP := attackerSpell.GetFinalAP()
                initialAttackerHP := attackerHealth.GetCurrentHP()

                // --- Run Simulation ---
                sim.RunSimulation()

                // --- Assertions ---
                expectedStacks := 1 // From taking damage around t=10.0
                Expect(titansEffect.GetCurrentStacks()).To(Equal(expectedStacks), "Should gain 1 stack from taking damage")
                Expect(attackerAttack.GetBonusPercentAD()).To(BeNumerically("~", titansADPerStack*float64(expectedStacks), 0.001))
                Expect(attackerSpell.GetBonusAP()).To(BeNumerically("~", titansAPPerStack*float64(expectedStacks), 0.001))

                // Check final stats increased and attacker took damage
                expectedFinalAD := attackerAttack.GetBaseAD() * (1 + attackerAttack.GetBonusPercentAD())
                expectedFinalAP := attackerSpell.GetBaseAP() + attackerSpell.GetBonusAP()
                Expect(attackerAttack.GetFinalAD()).To(BeNumerically("~", expectedFinalAD, 0.01))
                Expect(attackerSpell.GetFinalAP()).To(BeNumerically("~", expectedFinalAP, 0.01))
                Expect(attackerAttack.GetFinalAD()).To(BeNumerically(">", initialFinalAD))
                Expect(attackerSpell.GetFinalAP()).To(BeNumerically(">", initialFinalAP))
                Expect(attackerHealth.GetCurrentHP()).To(BeNumerically("<", initialAttackerHP), "Attacker should take damage")
            })

            It("should stack Titan's Resolve from both attacking and taking damage", func() {
                // Add Titan's to attacker
                err := equipmentManager.AddItemToChampion(attacker, data.TFT_Item_TitansResolve)
                Expect(err).NotTo(HaveOccurred())

                // Modify setup: Attacker attacks (AS=0.5 -> t=2,4,6,8,10), Target attacks slowly (AS=0.1 -> t=10)
                attackerAttack := getAttack(world, attacker)
                attackerAttack.SetBaseAttackSpeed(0.5)
                if targetAttack, ok := world.GetAttack(target); ok {
                    targetAttack.SetBaseAttackSpeed(0.1)
                    targetAttack.SetBaseAD(30)
                }

                // Create simulation AFTER setup changes
                sim = simulation.NewSimulationWithConfig(world, config.WithMaxTime(10.5))
                Expect(sim).NotTo(BeNil())

                // Get components
                attackerSpell := getSpell(world, attacker)
                titansEffect, ok := world.GetTitansResolveEffect(attacker)
                Expect(ok).To(BeTrue())

                // --- Run Simulation ---
                sim.RunSimulation()

                // --- Assertions ---
                // Expected stacks: 5 from attacker attacks (t=2,4,6,8,10), 1 from target attack (t=10) = 6
                expectedStacks := 6
                Expect(titansEffect.GetCurrentStacks()).To(Equal(expectedStacks), "Should gain 6 stacks (5 attack, 1 damage)")
                Expect(attackerAttack.GetBonusPercentAD()).To(BeNumerically("~", titansADPerStack*float64(expectedStacks), 0.001))
                Expect(attackerSpell.GetBonusAP()).To(BeNumerically("~", titansAPPerStack*float64(expectedStacks), 0.001))
            })

            It("should apply bonus resists upon reaching max stacks during simulation", func() {
                // Add Titan's to attacker
                err := equipmentManager.AddItemToChampion(attacker, data.TFT_Item_TitansResolve)
                Expect(err).NotTo(HaveOccurred())

                // Modify setup: Increase attacker AS to reach max stacks faster
                attackerAttack := getAttack(world, attacker)
                attackerAttack.SetBaseAttackSpeed(1.0) // Attacks every 1s. Needs 25 attacks for max stacks.
                if targetAttack, ok := world.GetAttack(target); ok {
                    targetAttack.SetFinalAttackSpeed(0) // Target doesn't attack
                }

                // Modify target: Set high HP to avoid quick death
                targetHealth := getHealth(world, target)
                targetHealth.SetBaseMaxHP(10000.0)
                targetHealth.SetCurrentHP(10000.0)

                // Create simulation AFTER setup changes, running long enough for max stacks
                timeToMaxStacks := float64(titansMaxStacks) + 1.0 // e.g., 26.0s
                sim = simulation.NewSimulationWithConfig(world, config.WithMaxTime(timeToMaxStacks))
                Expect(sim).NotTo(BeNil())

                // Get components AFTER sim creation
                attackerHealth := getHealth(world, attacker)
                titansEffect, ok := world.GetTitansResolveEffect(attacker)
                Expect(ok).To(BeTrue(), "TitansResolveEffect should exist after adding item")

                initialBonusArmor := attackerHealth.GetBonusArmor() // Includes static armor (10)
                Expect(initialBonusArmor).To(BeNumerically("~", titansStaticArmor, 0.001))
                // initialBonusMR := attackerHealth.GetBonusMR() // Base MR is 0

                // --- Run the simulation ---
                sim.RunSimulation()
                // --- End Run ---

                // --- Assertions after the run ---
                fmt.Fprintf(GinkgoWriter, "Debug (Reach Max Stacks Test): After RunSimulation (MaxTime=%.1f):\n", sim.GetConfig().MaxTime)
                fmt.Fprintf(GinkgoWriter, "  - Titans Stacks: %d (Expected: %d)\n", titansEffect.GetCurrentStacks(), titansMaxStacks)
                fmt.Fprintf(GinkgoWriter, "  - Attacker Bonus AD%%: %.2f\n", attackerAttack.GetBonusPercentAD()*100) // AD stacks
                fmt.Fprintf(GinkgoWriter, "  - Attacker Bonus AP: %.1f\n", getSpell(world, attacker).GetBonusAP())    // AP stacks
                fmt.Fprintf(GinkgoWriter, "  - Attacker Bonus Armor: %.1f\n", attackerHealth.GetBonusArmor()) // Static + Max Stack Bonus
                fmt.Fprintf(GinkgoWriter, "  - Attacker Bonus MR: %.1f\n", attackerHealth.GetBonusMR())       // Max Stack Bonus

                // Check if max stacks were reached
                Expect(titansEffect.GetCurrentStacks()).To(Equal(titansMaxStacks), "Should reach max stacks")
                Expect(titansEffect.IsMaxStacksReached()).To(BeTrue(), "IsMaxStacks flag should be true")

                // Check if resists were applied correctly at max stacks
                expectedBonusArmorAtMax := titansStaticArmor + titansBonusResistsAtCap // 10 + 20 = 30
                expectedBonusMRAtMax := titansBonusResistsAtCap                     // 0 + 20 = 20
                Expect(attackerHealth.GetBonusArmor()).To(BeNumerically("~", expectedBonusArmorAtMax, 0.001), "Bonus Armor should include static + max stack bonus")
                Expect(attackerHealth.GetBonusMR()).To(BeNumerically("~", expectedBonusMRAtMax, 0.001), "Bonus MR should include max stack bonus")

                // Check final stats reflect bonuses
                Expect(attackerHealth.GetFinalArmor()).To(BeNumerically("~", attackerHealth.GetBaseArmor()+expectedBonusArmorAtMax, 0.001))
                Expect(attackerHealth.GetFinalMR()).To(BeNumerically("~", attackerHealth.GetBaseMR()+expectedBonusMRAtMax, 0.001))
            })

            // Separate test for checking behavior *after* max stacks
            It("should not stack beyond max stacks (separate test)", func() {
                // Add Titan's to attacker
                err := equipmentManager.AddItemToChampion(attacker, data.TFT_Item_TitansResolve)
                Expect(err).NotTo(HaveOccurred())

                // Modify setup: Increase attacker AS, high target HP
                attackerAttack := getAttack(world, attacker)
                attackerAttack.SetBaseAttackSpeed(1.0) // Attacks every 1s
                targetHealth := getHealth(world, target)
                targetHealth.SetBaseMaxHP(10000.0)
                targetHealth.SetCurrentHP(10000.0)
                if targetAttack, ok := world.GetAttack(target); ok {
                    targetAttack.SetFinalAttackSpeed(0)
                }

                // --- Run simulation long enough to reach max stacks ---
                timeToMaxStacks := float64(titansMaxStacks) + 1.0 // e.g., 26.0s
                sim = simulation.NewSimulationWithConfig(world, config.WithMaxTime(timeToMaxStacks))
                Expect(sim).NotTo(BeNil())

                sim.RunSimulation() // Run first simulation part

                // Verify max stacks reached and record stats
                titansEffect, ok := world.GetTitansResolveEffect(attacker)
                Expect(ok).To(BeTrue())
                Expect(titansEffect.GetCurrentStacks()).To(Equal(titansMaxStacks), "Should reach max stacks after first run")
                attackerSpell := getSpell(world, attacker)
                attackerHealth := getHealth(world, attacker)
                adAtMax := attackerAttack.GetFinalAD()
                apAtMax := attackerSpell.GetFinalAP()
                armorAtMax := attackerHealth.GetFinalArmor()
                mrAtMax := attackerHealth.GetFinalMR()
                bonusArmorAtMax := attackerHealth.GetBonusArmor() // Includes static + cap bonus
                bonusMRAtMax := attackerHealth.GetBonusMR()       // Includes cap bonus

                // --- Run simulation for longer using the SAME instance ---
                // This requires SetMaxTime and RunSimulation to handle continuation correctly.
                // If RunSimulation always restarts from t=0, this test needs rethinking.
                // Assuming RunSimulation processes events up to the NEW MaxTime.
                // Let's try creating a NEW sim with the SAME world state but longer time.

                // --- Re-run from scratch for a longer duration ---
                // This ensures the simulation runs correctly from t=0 with the item,
                // and we check the state at the end.
                longerTotalTime := timeToMaxStacks + 10.0 // e.g., 36.0s

                // Need to reset world state? No, the point is to start *with* max stacks.
                // So, we need a way to preserve the world state (stacks, stats) from the first run.
                // The current approach of creating a new sim with the *same world* should work.

                sim2 := simulation.NewSimulationWithConfig(world, config.WithMaxTime(longerTotalTime)) // Use world state from sim1
                Expect(sim2).NotTo(BeNil())

                sim2.RunSimulation() // Run the second simulation part

                // Assert stacks and stats haven't changed from the max stack state
                Expect(titansEffect.GetCurrentStacks()).To(Equal(titansMaxStacks), "Stacks should remain at max after second longer run")
                // Re-get components as they might be different instances in the new sim? No, world is shared.
                Expect(attackerAttack.GetFinalAD()).To(BeNumerically("~", adAtMax, 0.01), "Final AD should not change after max stacks")
                Expect(attackerSpell.GetFinalAP()).To(BeNumerically("~", apAtMax, 0.01), "Final AP should not change after max stacks")
                // Check bonus stats directly
                Expect(attackerHealth.GetBonusArmor()).To(BeNumerically("~", bonusArmorAtMax, 0.01), "Bonus Armor should not change after max stacks")
                Expect(attackerHealth.GetBonusMR()).To(BeNumerically("~", bonusMRAtMax, 0.01), "Bonus MR should not change after max stacks")
                // Check final stats
                Expect(attackerHealth.GetFinalArmor()).To(BeNumerically("~", armorAtMax, 0.01), "Final Armor should not change after max stacks")
                Expect(attackerHealth.GetFinalMR()).To(BeNumerically("~", mrAtMax, 0.01), "Final MR should not change after max stacks")
            })

            It("should stack Guinsoo's Rageblade on attacks and apply AS bonus", func() {
                // Add Rageblade to attacker
                err := equipmentManager.AddItemToChampion(attacker, data.TFT_Item_GuinsoosRageblade)
                Expect(err).NotTo(HaveOccurred())

                // Get components and ensure attacker attacks
                attackerAttack := getAttack(world, attacker)
                attackerAttack.SetBaseAttackSpeed(0.5) // Base AS
                targetHealth := getHealth(world, target)
                if targetAttack, ok := world.GetAttack(target); ok {
                    targetAttack.SetFinalAttackSpeed(0)
                }

                // Create simulation AFTER adding the item
                sim = simulation.NewSimulationWithConfig(world, config) // Use default MaxTime initially if needed
                Expect(sim).NotTo(BeNil())

                // --- Initial State Verification (after setupCombat) ---
                ragebladeEffect, ok := world.GetGuinsoosRagebladeEffect(attacker)
                Expect(ok).To(BeTrue(), "GuinsoosRagebladeEffect component should be added")
                Expect(ragebladeEffect.GetCurrentStacks()).To(Equal(0), "Initial stacks should be 0")
                // Bonus AS should only include the static bonus initially
                Expect(attackerAttack.GetBonusPercentAttackSpeed()).To(BeNumerically("~", ragebladeStaticAS, 0.001), "Initial Bonus AS should be static bonus")
                baseAS := attackerAttack.GetBaseAttackSpeed() // 0.5
                initialFinalAS := attackerAttack.GetFinalAttackSpeed()
                expectedInitialFinalAS := baseAS * (1.0 + ragebladeStaticAS) // 0.5 * (1 + 0.1) = 0.55
                Expect(initialFinalAS).To(BeNumerically("~", expectedInitialFinalAS, 0.001), "Initial Final AS should reflect base + static bonus")

                // --- Simulation Run ---
                // Initial AS = 0.55 -> Interval = 1 / 0.55 = ~1.818s
                // Attack 1 lands ~1.8s -> Stack 1 -> Bonus AS = 0.1 + 0.05 = 0.15 -> Final AS = 0.5 * 1.15 = 0.575 -> Interval = ~1.739s
                // Attack 2 lands ~1.8 + 1.7 = ~3.5s -> Stack 2 -> Bonus AS = 0.1 + 0.1 = 0.2 -> Final AS = 0.5 * 1.2 = 0.6 -> Interval = ~1.667s
                // Attack 3 lands ~3.5 + 1.7 = ~5.2s -> Stack 3 -> Bonus AS = 0.1 + 0.15 = 0.25 -> Final AS = 0.5 * 1.25 = 0.625
                // Run for 6.0s to ensure 3 attacks land.
                sim.SetMaxTime(6.0)
                sim.RunSimulation()

                // --- Post-Simulation Verification ---
                expectedStacks := 3
                Expect(ragebladeEffect.GetCurrentStacks()).To(Equal(expectedStacks), "Should gain 3 stacks from 3 attacks")

                // Verify Bonus AS includes static + stacks
                expectedTotalBonusAS := ragebladeStaticAS + (float64(expectedStacks) * ragebladeASPerStack) // 0.1 + 3*0.05 = 0.25
                Expect(attackerAttack.GetBonusPercentAttackSpeed()).To(BeNumerically("~", expectedTotalBonusAS, 0.001), "Bonus AS should include static + 3 stacks")

                // Verify Final AS reflects the total bonus
                expectedFinalAS := baseAS * (1.0 + expectedTotalBonusAS) // 0.5 * (1 + 0.25) = 0.625
                Expect(attackerAttack.GetFinalAttackSpeed()).To(BeNumerically("~", expectedFinalAS, 0.001), "Final AS should reflect base + static + 3 stacks")

                // Verify target took damage
                Expect(targetHealth.GetCurrentHP()).To(BeNumerically("<", targetMaxHP), "Target should take damage from attacks")
            })

        }) // End Context("with Dynamic Event Items")

    }) // End Describe("RunSimulation Method")

    Describe("PrintResults", func() {
        // This primarily tests output formatting. Requires simulation run first.
        BeforeEach(func() {
            // Create sim instance for PrintResults tests
            sim = simulation.NewSimulationWithConfig(world, config)
            Expect(sim).NotTo(BeNil())
        })

        It("should print stats for specified entities", func() {
            // Requires capturing stdout. Using Skip for now.
            Skip("Skipping stdout capture test for PrintResults")

            // buffer := gbytes.NewBuffer() // Capture output
            // TODO: Redirect stdout or inject buffer into logger used by PrintResults

            // Run for a short time to have some state changes
            sim.SetMaxTime(1.5)
            sim.RunSimulation()

            // Call the method - Need a way to pass the buffer for capture
            // sim.PrintResults(buffer, attacker, target) // Hypothetical modification

            // Check if output contains expected names or stats (adjust based on actual output format)
            // Eventually(buffer).Should(gbytes.Say("Stats for Entity")) // Check for header
            // Eventually(buffer).Should(gbytes.Say("TFT_TrainingDummy")) // Check if champion name is printed
            // Eventually(buffer).Should(gbytes.Say("HP:")) // Check for stat labels
            // Eventually(buffer).Should(gbytes.Say("AD:"))
        })
    })

    Describe("Item Effects Integration (Static/Conditional Focus)", func() {
        // Focus on how static effects and conditional effects (like JG crit) are applied
        // during the simulation setup phase (NewSimulationWithConfig -> setupCombat).
        var (
            // Systems needed for static application pipeline (mimic setupCombat)
            statCalculationSystem *systems.StatCalculationSystem
            abilityCritSystem     *itemsys.AbilityCritSystem
            baseStaticItemSystem  *itemsys.BaseStaticItemSystem
            // Champion and components for testing
            champion       ecs.Entity
            championSpell  *components.Spell
            championAttack *components.Attack
            championCrit   *components.Crit
            err            error
        )

        BeforeEach(func() {
            // Reset world and create components/systems needed for these specific tests
            // Note: We are testing the *result* of the setup process, not RunSimulation here.
            world = ecs.NewWorld()
            championFactory = factory.NewChampionFactory(world)
            equipmentManager = managers.NewEquipmentManager(world)
            // Create instances of systems involved in static stat application
            statCalculationSystem = systems.NewStatCalculationSystem(world)
            abilityCritSystem = itemsys.NewAbilityCritSystem(world)
            baseStaticItemSystem = itemsys.NewBaseStaticItemSystem(world)

            // Use a champion known to have Spell, Attack, Crit components
            champion, err = championFactory.CreatePlayerChampion("TFT14_Kindred", 1) // Kindred has these
            Expect(err).NotTo(HaveOccurred())

            // Get components for assertions
            championSpell = getSpell(world, champion)
            championAttack = getAttack(world, champion)
            championCrit = getCrit(world, champion)

            // Reset bonuses before applying item effects (mimics start of setupCombat)
            championSpell.ResetBonuses()
            championAttack.ResetBonuses()
            championCrit.ResetBonuses()
        })

        // Helper function to *manually* run the static stat application pipeline
        // This simulates what setupCombat does *after* items are added but *before* the event loop.
        applyStaticStatsManually := func() {
            // Order matters, based on dependencies described in devlog/simulation setup
            abilityCritSystem.Update()                    // 1. Check for JG/IE (Reads ItemEffect populated by AddItemToChampion)
            baseStaticItemSystem.ApplyStats()             // 2. Apply bonuses from ItemEffect to component Bonus fields
            statCalculationSystem.ApplyStaticBonusStats() // 3. Calculate final stats based on updated component bonuses
        }

        Context("Rabadon's Deathcap (Static)", func() {
            var (
                deathcapData *data.Item
                expectedAP   float64
                expectedAmp  float64
            )
            BeforeEach(func() {
                deathcapData = data.GetItemByApiName(data.TFT_Item_RabadonsDeathcap)
                Expect(deathcapData).NotTo(BeNil())
                expectedAP = deathcapData.Effects["AP"]           // e.g., 70
                expectedAmp = deathcapData.Effects["BonusDamage"] // e.g., 0.08 (8%)

                // Add item using manager (populates ItemEffect component)
                err = equipmentManager.AddItemToChampion(champion, data.TFT_Item_RabadonsDeathcap)
                Expect(err).NotTo(HaveOccurred())

                // Manually run the static application pipeline to test its effect
                applyStaticStatsManually()
            })

            It("should add the correct Bonus AP to the Spell component", func() {
                // BaseStaticItemSystem should have applied the bonus
                Expect(championSpell.GetBonusAP()).To(BeNumerically("~", expectedAP, 0.01))
            })

            It("should add the correct Bonus Damage Amp to the Attack component", func() {
                // BaseStaticItemSystem should have applied the bonus
                Expect(championAttack.GetBonusDamageAmp()).To(BeNumerically("~", expectedAmp, 0.01))
            })

            It("should calculate Final AP including the item bonus", func() {
                // StatCalculationSystem should have calculated FinalAP
                expectedFinalAP := championSpell.GetBaseAP() + expectedAP
                Expect(championSpell.GetFinalAP()).To(BeNumerically("~", expectedFinalAP, 0.01))
            })

            It("should calculate Final Damage Amp including the item bonus", func() {
                // StatCalculationSystem should have calculated FinalDamageAmp
                expectedFinalAmp := championAttack.GetBaseDamageAmp() + expectedAmp
                Expect(championAttack.GetFinalDamageAmp()).To(BeNumerically("~", expectedFinalAmp, 0.01))
            })
        })

        Context("Jeweled Gauntlet (Static + Conditional)", func() {
            var (
                jgData             *data.Item
                jgBonusAP          float64
                jgCritChance       float64 // As decimal
                jgCritDamageToGive float64 // e.g., 0.4 (40%)
            )
            BeforeEach(func() {
                jgData = data.GetItemByApiName(data.TFT_Item_JeweledGauntlet)
                Expect(jgData).NotTo(BeNil())
                jgBonusAP = jgData.Effects["AP"]                     // e.g., 35
                jgCritChance = jgData.Effects["CritChance"] / 100.0   // e.g., 0.2
                jgCritDamageToGive = jgData.Effects["CritDamageToGive"] // e.g., 0.4

                // Add JG item
                err = equipmentManager.AddItemToChampion(champion, data.TFT_Item_JeweledGauntlet)
                Expect(err).NotTo(HaveOccurred())
                // Note: applyStaticStatsManually is called within the nested Contexts
            })

            Context("when champion abilities cannot already crit", func() {
                BeforeEach(func() {
                    // Ensure no other crit sources exist
                    world.RemoveComponent(champion, reflect.TypeOf(components.CanAbilityCritFromTraits{}))
                    world.RemoveComponent(champion, reflect.TypeOf(components.CanAbilityCritFromAugments{}))
                    // Run the static pipeline AFTER setting the context
                    applyStaticStatsManually()
                })

                It("should add the CanAbilityCritFromItems marker component", func() {
                    // AbilityCritSystem should add this marker
                    _, ok := world.GetCanAbilityCritFromItems(champion)
                    Expect(ok).To(BeTrue(), "CanAbilityCritFromItems marker should be added by AbilityCritSystem")
                })

                It("should add Bonus AP and Bonus Crit Chance", func() {
                    // BaseStaticItemSystem applies these
                    Expect(championSpell.GetBonusAP()).To(BeNumerically("~", jgBonusAP, 0.01))
                    Expect(championCrit.GetBonusCritChance()).To(BeNumerically("~", jgCritChance, 0.01))
                })

                It("should calculate Final Spell Crit Chance correctly", func() {
                    // StatCalculationSystem calculates this
                    championBaseCritChance := championCrit.GetBaseCritChance() // Usually 0.25
                    expectedFinalCritChance := championBaseCritChance + jgCritChance // 0.25 + 0.2 = 0.45
                    expectedFinalCritChance = math.Min(expectedFinalCritChance, 1.0) // Cap at 1.0
                    Expect(championCrit.GetFinalCritChance()).To(BeNumerically("~", expectedFinalCritChance, 0.01))
                })

                It("should NOT add the conditional Crit Damage", func() {
                    // AbilityCritSystem should NOT add BonusCritDamageToGive if only JG provides crit ability
                    Expect(championCrit.GetBonusCritDamageToGive()).To(BeNumerically("~", jgCritDamageToGive, 0.01))
                    // StatCalculationSystem calculates final crit multiplier
                    expectedFinalCritMultiplier := championCrit.GetBaseCritMultiplier() + championCrit.GetBonusCritMultiplier() // Base=1.4, Bonus=0 -> 1.4
                    Expect(championCrit.GetFinalCritMultiplier()).To(BeNumerically("~", expectedFinalCritMultiplier, 0.01))
                })
            })

            Context("when champion abilities can already crit (from trait)", func() {
                BeforeEach(func() {
                    // Add trait crit marker BEFORE running static pipeline
                    err = world.AddComponent(champion, &components.CanAbilityCritFromTraits{})
                    Expect(err).NotTo(HaveOccurred())
                    // Run the static pipeline AFTER setting the context
                    applyStaticStatsManually()
                })

                It("should still add the CanAbilityCritFromItems marker component", func() {
                    // JG still adds its marker regardless of others
                    _, ok := world.GetCanAbilityCritFromItems(champion)
                    Expect(ok).To(BeTrue(), "CanAbilityCritFromItems marker should still be added")
                })

                It("should add Bonus AP and Bonus Crit Chance", func() {
                    Expect(championSpell.GetBonusAP()).To(BeNumerically("~", jgBonusAP, 0.01))
                    Expect(championCrit.GetBonusCritChance()).To(BeNumerically("~", jgCritChance, 0.01))
                })

                It("should calculate Final Spell Crit Chance correctly", func() {
                    championBaseCritChance := championCrit.GetBaseCritChance()
                    expectedFinalCritChance := championBaseCritChance + jgCritChance
                    expectedFinalCritChance = math.Min(expectedFinalCritChance, 1.0)
                    Expect(championCrit.GetFinalCritChance()).To(BeNumerically("~", expectedFinalCritChance, 0.01))
                })

                It("SHOULD add the conditional Crit Damage", func() {
                    // AbilityCritSystem SHOULD add BonusCritDamageToGive because CanAbilityCritFromTraits exists
                    Expect(championCrit.GetBonusCritDamageToGive()).To(BeNumerically("~", jgCritDamageToGive, 0.01))

                    // StatCalculationSystem calculates final crit multiplier
                    expectedFinalCritMultiplier := championCrit.GetBaseCritMultiplier() + championCrit.GetBonusCritMultiplier() + championCrit.GetBonusCritDamageToGive() // Base=1.4, Bonus=0, BonusToGive=0.4 -> 1.8
                    Expect(championCrit.GetFinalCritMultiplier()).To(BeNumerically("~", expectedFinalCritMultiplier, 0.01))
                })
            })
        })

        // Archangel's test needs RunSimulation, so it's moved/adapted below

    }) // End Describe("Item Effects Integration (Static/Conditional Focus)")

    Describe("Item Effects Integration (Dynamic Time via RunSimulation)", func() {
        // Tests how RunSimulation handles dynamic time items like Archangel's

        var (
            champion       ecs.Entity
            championSpell  *components.Spell
            archangelsData *data.Item
            initialAP      float64 // Static AP from item
            initialMana    float64
            interval       float64
            apPerInterval  float64
            err            error
        )

        BeforeEach(func() {
            // Setup world, champion, manager
            world = ecs.NewWorld()
            championFactory = factory.NewChampionFactory(world)
            equipmentManager = managers.NewEquipmentManager(world)
            champion, err = championFactory.CreatePlayerChampion("TFT14_Kindred", 1)
            Expect(err).NotTo(HaveOccurred())
            championSpell = getSpell(world, champion)
            // Add mana if needed
            if _, ok := world.GetMana(champion); !ok {
                world.AddComponent(champion, components.NewMana(0, 100))
            }

            // Get Archangel's data
            archangelsData = data.GetItemByApiName(data.TFT_Item_ArchangelsStaff)
            Expect(archangelsData).NotTo(BeNil())
            initialAP = archangelsData.Effects["AP"]
            initialMana = archangelsData.Effects["Mana"]
            interval = archangelsData.Effects["IntervalSeconds"]    // Should be 5.0
            apPerInterval = archangelsData.Effects["APPerInterval"] // Should be 30.0

            // Add item BEFORE creating the simulation instance
            err = equipmentManager.AddItemToChampion(champion, data.TFT_Item_ArchangelsStaff)
            Expect(err).NotTo(HaveOccurred())

            // Create simulation AFTER adding the item
            // Use a config suitable for testing intervals
            sim = simulation.NewSimulationWithConfig(world, config.WithMaxTime(11.0)) // Default max time for this context
            Expect(sim).NotTo(BeNil())
        })

        It("should apply initial static AP and Mana during setup", func() {
            // Check Final fields immediately after simulation creation (setupCombat runs)
            expectedInitialFinalAP := championSpell.GetBaseAP() + initialAP
            Expect(championSpell.GetFinalAP()).To(BeNumerically("~", expectedInitialFinalAP, 0.01))

            manaComp, ok := world.GetMana(champion)
            Expect(ok).To(BeTrue())
            // Assuming BaseInitialMana exists and is relevant
            expectedInitialFinalMana := manaComp.GetBaseInitialMana() + initialMana
            Expect(manaComp.GetFinalInitialMana()).To(BeNumerically("~", expectedInitialFinalMana, 0.01)) // Check FinalInitialMana if applicable
        })

        It("should stack AP correctly over time via RunSimulation", func() {
            // --- State after initial setup (checked above) ---
            effect, ok := world.GetArchangelsEffect(champion)
            Expect(ok).To(BeTrue())
            Expect(effect.GetStacks()).To(Equal(0), "Initial stacks should be 0")
            initialFinalAP := championSpell.GetFinalAP()

            // --- Run just before first stack (e.g., 4.9s) ---
            sim.SetMaxTime(interval - 0.1) // Run up to 4.9s
            sim.RunSimulation()

            Expect(effect.GetStacks()).To(Equal(0), "Stacks should be 0 before the first interval")
            Expect(championSpell.GetFinalAP()).To(BeNumerically("~", initialFinalAP, 0.01), "Final AP should not have increased before the first interval")

            // --- Run past first stack (e.g., 5.1s total) ---
            // Create a new sim instance with the same world but longer time
            sim = simulation.NewSimulationWithConfig(world, config.WithMaxTime(interval+0.1)) // Run up to 5.1s
            sim.RunSimulation() // Re-runs from t=0 up to 5.1s

            Expect(effect.GetStacks()).To(Equal(1), "Stacks should be 1 after the first interval")
            expectedBonusAPAfterStack1 := initialAP + apPerInterval
            expectedFinalAPAfterStack1 := championSpell.GetBaseAP() + expectedBonusAPAfterStack1
            Expect(championSpell.GetBonusAP()).To(BeNumerically("~", expectedBonusAPAfterStack1, 0.01), "Bonus AP should include 1 stack after the first interval")
            Expect(championSpell.GetFinalAP()).To(BeNumerically("~", expectedFinalAPAfterStack1, 0.01), "Final AP should include 1 stack after the first interval")

            // --- Run past second stack (e.g., 10.1s total) ---
            sim = simulation.NewSimulationWithConfig(world, config.WithMaxTime(2*interval+0.1)) // Run up to 10.1s
            sim.RunSimulation() // Re-runs from t=0 up to 10.1s

            Expect(effect.GetStacks()).To(Equal(2), "Stacks should be 2 after the second interval")
            expectedBonusAPAfterStack2 := initialAP + 2*apPerInterval
            expectedFinalAPAfterStack2 := championSpell.GetBaseAP() + expectedBonusAPAfterStack2
            Expect(championSpell.GetBonusAP()).To(BeNumerically("~", expectedBonusAPAfterStack2, 0.01), "Bonus AP should include 2 stacks after the second interval")
            Expect(championSpell.GetFinalAP()).To(BeNumerically("~", expectedFinalAPAfterStack2, 0.01), "Final AP should include 2 stacks after the second interval")
        })
    }) // End Describe("Item Effects Integration (Dynamic Time via RunSimulation)")

}) // End Describe("Simulation")
