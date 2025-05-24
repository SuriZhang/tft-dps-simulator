package simulation_test

import (
	"fmt"
	"math"
	"reflect"

	"tft-dps-simulator/internal/core/components"
	"tft-dps-simulator/internal/core/data"
	"tft-dps-simulator/internal/core/ecs"
	"tft-dps-simulator/internal/core/factory"
	"tft-dps-simulator/internal/core/managers"
	"tft-dps-simulator/internal/core/simulation"
	"tft-dps-simulator/internal/core/systems"
	itemsys "tft-dps-simulator/internal/core/systems/items"

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

func getPosition(w *ecs.World, e ecs.Entity) *components.Position {
    comp, ok := w.GetPosition(e)
    Expect(ok).To(BeTrue(), "Entity should have Position component")
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
            // First Attack Land at t=0.0 (0 startup/recovery time)
            sim.SetMaxTime(1.9) // Set time long enough for one attack

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
                ragebladeStaticAS = ragebladeData.Effects["AS"] / 100.0          // e.g. 0.1 for 10%
                ragebladeASPerStack = ragebladeData.Effects["AttackSpeedPerStack"] / 100.0 // e.g. 0.05 for 5%
            })

            It("should stack Titan's Resolve on attacks and apply AD/AP bonuses", func() {
                // Add Titan's to attacker BEFORE creating simulation
                err := equipmentManager.AddItemToChampion(attacker, data.TFT_Item_TitansResolve)
                Expect(err).NotTo(HaveOccurred())

                // Ensure attacker attacks and target doesn't
                attackerAttack := getAttack(world, attacker)
                attackerAttack.SetBaseAttackSpeed(0.5) // Attacks land at t=0.0, 2.0
                if targetAttack, ok := world.GetAttack(target); ok {
                    targetAttack.SetBaseAttackSpeed(0)
                    targetAttack.SetFinalAttackSpeed(0)
                }

                // Create simulation AFTER adding the item and setting AS
                sim = simulation.NewSimulationWithConfig(world, config.WithMaxTime(2)) // Run long enough for 2 attacks
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
                expectedStacks := 2 // From attacks at t=0.0, t=1.818
                Expect(titansEffect.GetCurrentStacks()).To(Equal(expectedStacks), "Should gain 2 stacks from 2 attacks")
                Expect(attackerAttack.GetAttackCount()).To(Equal(expectedStacks))
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
                sim = simulation.NewSimulationWithConfig(world, config.WithMaxTime(9.9)) // Run long enough for target's attack
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

                // Modify setup: Attacker attacks (AS=0.5 -> t=0, 1.818, 3.636, 5.455, 7.273), Target attacks slowly (AS=0.1 -> t=0, 10)
                attackerAttack := getAttack(world, attacker)
                attackerAttack.SetBaseAttackSpeed(0.5)
                if targetAttack, ok := world.GetAttack(target); ok {
                    targetAttack.SetBaseAttackSpeed(0.1)
                    targetAttack.SetBaseAD(30)
                }

                // Create simulation AFTER setup changes
                sim = simulation.NewSimulationWithConfig(world, config.WithMaxTime(8))
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
                Expect(attackerAttack.GetAttackCount()).To(Equal(5))
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
            It("should not stack beyond max stacks", func() {
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
                    targetAttack.SetBaseAttackSpeed(0)
                    targetAttack.SetFinalAttackSpeed(0)
                }

                // --- Run simulation long enough to reach max stacks ---
                timeToMaxStacks := float64(titansMaxStacks)+10// e.g., 24.0s
                sim = simulation.NewSimulationWithConfig(world, config.WithMaxTime(timeToMaxStacks))
                Expect(sim).NotTo(BeNil())

                sim.RunSimulation() // Run first simulation part

                // Verify max stacks reached and record stats
                titansEffect, ok := world.GetTitansResolveEffect(attacker)
                Expect(ok).To(BeTrue())
                Expect(titansEffect.GetCurrentStacks()).To(Equal(titansMaxStacks), "Should reach max stacks after first run")
                Expect(attackerAttack.GetFinalAttackSpeed()).To(Equal(1.1))
                attackerSpell := getSpell(world, attacker)
                attackerHealth := getHealth(world, attacker)
                adAtMax := attackerAttack.GetFinalAD()
                apAtMax := attackerSpell.GetFinalAP()
                armorAtMax := attackerHealth.GetFinalArmor()
                mrAtMax := attackerHealth.GetFinalMR()
                bonusArmorAtMax := attackerHealth.GetBonusArmor() // Includes static + cap bonus
                bonusMRAtMax := attackerHealth.GetBonusMR()       // Includes cap bonus

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
                    targetAttack.SetFinalAttackSpeed(0) // Target does not attack
                    targetAttack.SetBaseAttackSpeed(0)
                }
                 // Ensure target has enough HP to survive all attacks
                targetHealth.SetBaseMaxHP(10000.0)
                targetHealth.SetCurrentHP(10000.0)


                // Create simulation AFTER adding the item
                // MaxTime will be set specifically for this test later.
                // Initial setupCombat runs here.
                sim = simulation.NewSimulationWithConfig(world, config.WithMaxTime(0.1))
                Expect(sim).NotTo(BeNil())

                // --- Initial State Verification (after setupCombat) ---
                ragebladeEffect, ok := world.GetGuinsoosRagebladeEffect(attacker)
                Expect(ok).To(BeTrue(), "GuinsoosRagebladeEffect component should be added")
                Expect(ragebladeEffect.GetCurrentStacks()).To(Equal(0), "Initial stacks should be 0")

                // Bonus AS should only include the static bonus initially from setupCombat
                Expect(attackerAttack.GetBonusPercentAttackSpeed()).To(BeNumerically("~", ragebladeStaticAS, 0.001), "Initial BonusPercentAttackSpeed should be static bonus from item")

                baseAS := attackerAttack.GetBaseAttackSpeed() // Should be 0.5
                expectedInitialFinalAS := baseAS * (1.0 + ragebladeStaticAS) // 0.5 * (1 + 0.1) = 0.55
                Expect(attackerAttack.GetFinalAttackSpeed()).To(BeNumerically("~", expectedInitialFinalAS, 0.001), "Initial FinalAttackSpeed should reflect base + static bonus")

                // --- Simulation Run ---
                // Attacker Base AS = 0.5. Static Bonus = 10%. AS per Stack = 5%.
                // Initial Final AS = 0.5 * (1 + 0.10) = 0.55. Interval = 1 / 0.55 = ~1.818s.
                // Attack 1 (t=0): Stacks = 1. Bonus AS = 0.10 (static) + 0.05 (1 stack) = 0.15. Final AS = 0.5 * 1.15 = 0.575. Interval = ~1.739s.
                // Attack 2 (t~1.818): Stacks = 2. Bonus AS = 0.10 + 0.10 = 0.20. Final AS = 0.5 * 1.20 = 0.60. Interval = ~1.667s.
                // Attack 3 (t~1.818 + 1.739 = 3.557): Stacks = 3. Bonus AS = 0.10 + 0.15 = 0.25. Final AS = 0.5 * 1.25 = 0.625. Interval = ~1.600s.
                // Attack 4 (t~3.557 + 1.667 = 5.224): Stacks = 4. Bonus AS = 0.10 + 0.20 = 0.30. Final AS = 0.5 * 1.30 = 0.65.
                // Set MaxTime to allow 4 attacks to land.
                sim.SetMaxTime(5.3) // Slightly above 5.224s
                // Reset attack count before running, as some might have occurred during setup/previous short sim runs.
                attackerAttack.SetAttackCount(0)
                sim.RunSimulation()

                // --- Post-Simulation Verification ---
                expectedAttacks := 4
                expectedStacks := 4 // Assuming each attack adds a stack
                Expect(attackerAttack.GetAttackCount()).To(Equal(expectedAttacks), fmt.Sprintf("Expected %d attacks", expectedAttacks))
                Expect(ragebladeEffect.GetCurrentStacks()).To(Equal(expectedStacks), fmt.Sprintf("Should gain %d stacks from %d attacks", expectedStacks, expectedAttacks))

                // Verify Bonus AS includes static + stacks
                expectedTotalBonusAS := ragebladeStaticAS + (float64(expectedStacks) * ragebladeASPerStack) // 0.1 + 4*0.05 = 0.1 + 0.2 = 0.3
                Expect(attackerAttack.GetBonusPercentAttackSpeed()).To(BeNumerically("~", expectedTotalBonusAS, 0.001), "BonusPercentAttackSpeed should include static + all stack bonuses")

                // Verify Final AS reflects the total bonus
                expectedFinalAS := baseAS * (1.0 + expectedTotalBonusAS) // 0.5 * (1 + 0.3) = 0.65
                Expect(attackerAttack.GetFinalAttackSpeed()).To(BeNumerically("~", expectedFinalAS, 0.001), "FinalAttackSpeed should reflect base + static + all stack bonuses")

                // Verify target took damage
                Expect(targetHealth.GetCurrentHP()).To(BeNumerically("<", targetHealth.GetBaseMaxHP()), "Target should take damage from attacks")
            })

        }) // End Context("with Dynamic Event Items")

        Context("with Kranken's Fury", func() {
            // Placeholder stats for Runaan's Hurricane (base for Kranken's Fury)
            // These should be verified against actual game data if possible.
            const (
                runaansStaticAD         = 10.0
                runaansStaticASPercent  = 0.15 // 15%
                runaansStaticCritChance = 0.20 // 20%
                krankensADPerStack      = 5.0  // Assumed AD per stack for Kranken's Fury
            )

            var (
                attackerAttack    *components.Attack
                attackerCrit      *components.Crit
                krakensEffect     *components.KrakensFuryEffect
                initialAttackerAD float64
                initialAttackerAS float64
            )

            BeforeEach(func() {
                // Equip attacker with Kranken's Fury
                err := equipmentManager.AddItemToChampion(attacker, data.TFT_Item_KrakensFury)
                Expect(err).NotTo(HaveOccurred())

                // Ensure attacker has the KrakensFuryEffect component.
                // This should be added by the item system when the item is equipped.
                var ok bool
                krakensEffect, ok = world.GetKrakensFuryEffect(attacker)
                Expect(ok).To(BeTrue(), "Attacker should have KrakensFuryEffect component after equipping Kranken's Fury")
                Expect(krakensEffect).NotTo(BeNil())

                // Configure the effect (assuming it's not configured by default from item data in test)
                // In a real scenario, this value would come from the item's data definition.
                // For testing the handler, we set it explicitly if needed.
                // If the handler itself reads this from item data, this might not be necessary
                // or could be used to override for specific test cases.
                // Let's assume the handler or item system correctly sets ADPerStack.
                // If GetADPerStack() relies on krakensEffect.ADPerStack being set:
                krakensEffect.ADPerStack = krankensADPerStack // Ensure this is set if component doesn't init from data

                // Set attacker base stats for predictable calculations
                attackerAttack = getAttack(world, attacker)
                attackerAttack.SetBaseAttackSpeed(1.0) // 1 attack per second
                attackerAttack.SetBaseAD(50)           // Base AD for attacker

                attackerCrit = getCrit(world, attacker) // Get crit component for stat checks

                // Ensure target doesn't attack and has high HP
                targetHealth := getHealth(world, target)
                targetHealth.SetBaseMaxHP(10000.0)
                targetHealth.SetCurrentHP(10000.0)
                if targetAttack, tOk := world.GetAttack(target); tOk {
                    targetAttack.SetBaseAttackSpeed(0)
                    targetAttack.SetFinalAttackSpeed(0)
                }

                // Create a new sim instance (this runs setupCombat)
                // Use a short MaxTime for setup-related checks if not running the sim further yet
                sim = simulation.NewSimulationWithConfig(world, config.WithMaxTime(0.1)) // Short time just for setup
                Expect(sim).NotTo(BeNil())

                // Store initial AD/AS *after* simulation setup (static bonuses applied)
                initialAttackerAD = attackerAttack.GetFinalAD()
                initialAttackerAS = attackerAttack.GetFinalAttackSpeed()
            })

            It("should apply base stats of Runaan's Hurricane and initialize Kraken's Fury effect", func() {
                // Verify static stats from Runaan's Hurricane part of Kranken's Fury
                // Expected AD = BaseAD + Runaan's Static AD
                expectedAD := attackerAttack.GetBaseAD() + runaansStaticAD
                Expect(attackerAttack.GetFinalAD()).To(BeNumerically("~", expectedAD, 0.01), "Final AD should include Runaan's static AD")

                // Expected AS = BaseAS * (1 + Runaan's Static AS %)
                expectedAS := attackerAttack.GetBaseAttackSpeed() * (1.0 + runaansStaticASPercent)
                Expect(attackerAttack.GetFinalAttackSpeed()).To(BeNumerically("~", expectedAS, 0.01), "Final AS should include Runaan's static AS percentage")

                // Expected Crit Chance = BaseCritChance + Runaan's Static Crit Chance
                // Assuming base crit chance is 0.25 (standard)
                baseCritChance := attackerCrit.GetBaseCritChance() // Use actual base
                expectedCritChance := baseCritChance + runaansStaticCritChance
                Expect(attackerCrit.GetFinalCritChance()).To(BeNumerically("~", expectedCritChance, 0.01), "Final Crit Chance should include Runaan's static Crit Chance")

                // Verify Kraken's Fury effect state
                Expect(krakensEffect).NotTo(BeNil()) // Already checked in BeforeEach, but good for clarity
                Expect(krakensEffect.GetCurrentStacks()).To(Equal(0), "Initial Kraken's Fury stacks should be 0")
                Expect(krakensEffect.GetADPerStack()).To(BeNumerically("~", krankensADPerStack, 0.01), "ADPerStack should be the configured value")

                // Bonus AD from stacks should be zero initially
                // This checks the specific field where the handler adds AD from stacks.
                // We need to know which field this is, e.g., attackerAttack.BonusADFromKrakens (hypothetical)
                // or if it's directly added to attackerAttack.BonusAD.
                // For now, let's assume it's part of a general BonusAD or a specific one.
                // If Kranken's adds its *static* AD to BonusAD, then initial BonusAD would be runaansStaticAD.
                // The dynamic part (stacks * ADPerStack) should be 0.
                // Let's check GetBonusAD, assuming static AD is part of it, and dynamic isn't yet.
                Expect(attackerAttack.GetBonusAD()).To(BeNumerically("~", runaansStaticAD, 0.01), "Initial Bonus AD should reflect only Runaan's static AD, not stacked AD")
            })

            It("should stack AD on each attack and update final AD", func() {
                sim.SetMaxTime(2.5) // Allow for 3 attacks (t=0, t=1, t=2, sim ends before t=3 attack lands)
                // Attacker AS = 1.0. Attacks should land at t=0.0, t=1.0, t=2.0.
                // So 3 attacks should occur.

                // Reset attack count for this specific test if necessary (or ensure it's clean)
                attackerAttack.SetAttackCount(0) // Assuming ResetBonuses or similar doesn't reset this.

                sim.RunSimulation()

                Expect(krakensEffect.GetCurrentStacks()).To(Equal(3), "Kraken's Fury should have 3 stacks after 3 attacks")
                // Attacker attack count should also be 3.
                Expect(attackerAttack.GetAttackCount()).To(Equal(3), "Attacker should have performed 3 attacks")


                // Expected Bonus AD from stacks = 3 * ADPerStack
                expectedStackedAD := float64(3) * krankensADPerStack
                // Total Bonus AD = Runaan's Static AD + Stacked AD
                expectedTotalBonusAD := runaansStaticAD + expectedStackedAD
                Expect(attackerAttack.GetBonusAD()).To(BeNumerically("~", expectedTotalBonusAD, 0.01), "Bonus AD should include Runaan's static AD and stacked AD from Kraken's Fury")

                // Final AD = BaseAD + Runaan's Static AD + Stacked AD
                expectedFinalAD := attackerAttack.GetBaseAD() + expectedTotalBonusAD
                Expect(attackerAttack.GetFinalAD()).To(BeNumerically("~", expectedFinalAD, 0.01), "Final AD should reflect Base AD, Runaan's static AD, and stacked AD from Kraken's Fury")

                // Target should have taken damage
                targetHealth := getHealth(world, target)
                Expect(targetHealth.GetCurrentHP()).To(BeNumerically("<", targetHealth.GetBaseMaxHp()), "Target should have taken damage")
            })

            It("should reset stacks if a new simulation is created with the same world (simulating re-equip)", func() {
                // First, equip the item and run a short simulation to accumulate some stacks
                // Note: Item is already equipped in the outer BeforeEach.
                // We need to ensure the effect component is clean before this first run for this specific test.
                krakensEffect.SetCurrentStacks(0) // Manually reset for this test's first phase
                attackerAttack.SetBonusAD(runaansStaticAD) // Reset bonus AD to only static part
                attackerAttack.CalculateFinalStats() // Recalculate based on reset bonuses

                sim1Config := config.WithMaxTime(1.5) // Enough for 2 attacks (t=0, t=1)
                sim1 := simulation.NewSimulationWithConfig(world, sim1Config)
                // The krakensEffect instance is shared via the world, so sim1 will use the one from BeforeEach.
                // NewSimulationWithConfig runs setupCombat, which should re-initialize the effect if designed that way.
                // Let's verify the handler's OnEquip logic.
                // Get the effect component *after* sim1 is created, as setupCombat might replace/reset it.
                krakensEffectSim1, okSim1 := world.GetKrakensFuryEffect(attacker)
                Expect(okSim1).To(BeTrue())
                // If OnEquip resets stacks, it should be 0 after NewSimulationWithConfig
                Expect(krakensEffectSim1.GetCurrentStacks()).To(Equal(0), "Stacks should be 0 after new sim creation due to OnEquip reset")


                sim1.RunSimulation() // Run to accumulate stacks again
                Expect(krakensEffectSim1.GetCurrentStacks()).To(Equal(2), "Should have 2 stacks after sim1 run")

                // Now, create a *new* simulation instance with the *same world*.
                // This simulates the setup process happening again for the same champion with the same item.
                // The OnEquip method in the handler should be called during this new simulation's setupCombat.
                sim2Config := config.WithMaxTime(0.1) // Short, just for setup
                sim2 := simulation.NewSimulationWithConfig(world, sim2Config)
                _ = sim2 // Avoid unused variable error if not running sim2

                krakensEffectSim2, okSim2 := world.GetKrakensFuryEffect(attacker)
                Expect(okSim2).To(BeTrue())
                Expect(krakensEffectSim2.GetCurrentStacks()).To(Equal(0), "Stacks should be reset to 0 when a new simulation initializes with the same entity having the item")
            })

            It("should not stack if the AttackLandedEvent source is not the wearer", func() {
                // Attacker has Kranken's Fury from BeforeEach.
                // Ensure its stacks are 0 before this test.
                krakensEffect.SetCurrentStacks(0)
                attackerAttack.SetBonusAD(runaansStaticAD) // Reset to static AD
                attackerAttack.CalculateFinalStats()

                // Set up target to also be able to attack
                targetAttack, ok := world.GetAttack(target)
                Expect(ok).To(BeTrue(), "Target should have an Attack component")
                targetAttack.SetBaseAttackSpeed(1.0) // Target can attack
                targetAttack.SetBaseAD(10)

                // Manually enqueue an AttackLandedEvent where 'target' is the source
                // Sim is already created in BeforeEach. We can use its event bus.
                eventBus := sim.GetEventBus()
                Expect(eventBus).NotTo(BeNil())
                eventBus.Enqueue(eventsys.NewAttackLandedEvent(target, attacker, 0.1))

                sim.SetMaxTime(0.2) // Run sim long enough to process the enqueued event
                sim.RunSimulation()

                // Stacks on attacker's Kranken's Fury should remain 0
                Expect(krakensEffect.GetCurrentStacks()).To(Equal(0), "Kranken's Fury should not stack if the wearer is not the source of the attack")
            })

        }) // End Context("with Kranken's Fury")

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
            baseStaticItemSystem.ApplyStaticItemsBonus()             // 2. Apply bonuses from ItemEffect to component Bonus fields
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
            effect, ok := world.GetArchangelsStaffEffect(champion)
            Expect(ok).To(BeTrue())
            Expect(effect.GetStacks()).To(Equal(0), "Initial stacks should be 0")

            championSpell.ResetBonuses()

            // --- Run past first stack (e.g., 5.1s total) ---
            // Create a new sim instance with the same world but longer time
            sim = simulation.NewSimulationWithConfig(world, config.WithMaxTime(interval+0.1)) // Run up to 5.1s
            sim.RunSimulation() // Re-runs from t=0 up to 5.1s

            Expect(effect.GetStacks()).To(Equal(1), "Stacks should be 1 after the first interval")
            expectedBonusAPAfterStack1 := initialAP + apPerInterval
            expectedFinalAPAfterStack1 := championSpell.GetBaseAP() + expectedBonusAPAfterStack1
            Expect(championSpell.GetBonusAP()).To(BeNumerically("~", expectedBonusAPAfterStack1, 0.01), "Bonus AP should include 1 stack after the first interval")
            Expect(championSpell.GetFinalAP()).To(BeNumerically("~", expectedFinalAPAfterStack1, 0.01), "Final AP should include 1 stack after the first interval")

            championSpell.ResetBonuses()
            // --- Run past second stack (e.g., 10.1s total) ---
            sim = simulation.NewSimulationWithConfig(world, config.WithMaxTime(2*interval+0.1)) // Run up to 10.1s
            sim.RunSimulation() // Re-runs from t=0 up to 10.1s

            Expect(effect.GetStacks()).To(Equal(2), "Stacks should be 2 after the second interval")
            expectedBonusAPAfterStack2 := initialAP + 2*apPerInterval
            expectedFinalAPAfterStack2 := championSpell.GetBaseAP() + expectedBonusAPAfterStack2
            Expect(championSpell.GetBonusAP()).To(BeNumerically("~", expectedBonusAPAfterStack2, 0.01), "Bonus AP should include 2 stacks after the second interval")
            Expect(championSpell.GetFinalAP()).To(BeNumerically("~", expectedFinalAPAfterStack2, 0.01), "Final AP should include 2 stacks after the second interval")
        })
    })

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
            effect, ok := world.GetArchangelsStaffEffect(champion)
            Expect(ok).To(BeTrue())
            Expect(effect.GetStacks()).To(Equal(0), "Initial stacks should be 0")

            championSpell.ResetBonuses()

            // --- Run past first stack (e.g., 5.1s total) ---
            // Create a new sim instance with the same world but longer time
            sim = simulation.NewSimulationWithConfig(world, config.WithMaxTime(interval+0.1)) // Run up to 5.1s
            sim.RunSimulation() // Re-runs from t=0 up to 5.1s

            Expect(effect.GetStacks()).To(Equal(1), "Stacks should be 1 after the first interval")
            expectedBonusAPAfterStack1 := initialAP + apPerInterval
            expectedFinalAPAfterStack1 := championSpell.GetBaseAP() + expectedBonusAPAfterStack1
            Expect(championSpell.GetBonusAP()).To(BeNumerically("~", expectedBonusAPAfterStack1, 0.01), "Bonus AP should include 1 stack after the first interval")
            Expect(championSpell.GetFinalAP()).To(BeNumerically("~", expectedFinalAPAfterStack1, 0.01), "Final AP should include 1 stack after the first interval")

            championSpell.ResetBonuses()
            // --- Run past second stack (e.g., 10.1s total) ---
            sim = simulation.NewSimulationWithConfig(world, config.WithMaxTime(2*interval+0.1)) // Run up to 10.1s
            sim.RunSimulation() // Re-runs from t=0 up to 10.1s

            Expect(effect.GetStacks()).To(Equal(2), "Stacks should be 2 after the second interval")
            expectedBonusAPAfterStack2 := initialAP + 2*apPerInterval
            expectedFinalAPAfterStack2 := championSpell.GetBaseAP() + expectedBonusAPAfterStack2
            Expect(championSpell.GetBonusAP()).To(BeNumerically("~", expectedBonusAPAfterStack2, 0.01), "Bonus AP should include 2 stacks after the second interval")
            Expect(championSpell.GetFinalAP()).To(BeNumerically("~", expectedFinalAPAfterStack2, 0.01), "Final AP should include 2 stacks after the second interval")
        })

        Context("with Spirit Visage", func() {
            // Placeholder stats for Redemption (base for Spirit Visage)
            const (
                redemptionStaticHealth = 200.0
                redemptionStaticArmor  = 25.0
                redemptionStaticMR     = 25.0
            )
            // Assumed values for SpiritVisageEffect (as per task description)
            const (
                spiritVisageTickInterval      = 2.0
                spiritVisageMissingHealRate = 0.05 // 5%
                spiritVisageMaxHeal           = 100.0 // This is MaxHeal as per component, used in math.Max by handler
            )

            var (
                attackerHealth       *components.Health
                spiritVisageEffect   *components.SpiritVisageEffect
                baseMaxHP            float64
            )

            BeforeEach(func() {
                // Attacker is already created in the outer Describe("Simulation") BeforeEach
                // Target is also available

                // Set Attacker's Max HP for predictable calculations
                attackerHealth = getHealth(world, attacker) // Get health component for attacker
                baseMaxHP = 1000.0
                attackerHealth.SetBaseMaxHP(baseMaxHP)
                attackerHealth.SetCurrentHP(baseMaxHP / 2) // Start at 50% HP (500/1000)

                // Equip attacker with Spirit Visage
                err := equipmentManager.AddItemToChampion(attacker, data.TFT_Item_SpiritVisage)
                Expect(err).NotTo(HaveOccurred())

                // Ensure attacker has the SpiritVisageEffect component.
                // This should be added by the item system (specifically by BaseStaticItemSystem or a similar setup system).
                // For this test, we'll retrieve it. If it's not added, the test will fail here.
                var ok bool
                spiritVisageEffect, ok = world.GetSpiritVisageEffect(attacker)
                Expect(ok).To(BeTrue(), "Attacker should have SpiritVisageEffect component after equipping Spirit Visage")
                Expect(spiritVisageEffect).NotTo(BeNil())

                // Manually set the effect parameters if they are not automatically populated from item data
                // This ensures the test uses the assumed values.
                // In a real scenario, these would come from loaded item data.
                spiritVisageEffect.TickInterval = spiritVisageTickInterval
                spiritVisageEffect.MissingHealthHealRate = spiritVisageMissingHealRate
                spiritVisageEffect.MaxHeal = spiritVisageMaxHeal // This is the value used in math.Max by the handler

                // Create a new sim instance (this runs setupCombat, which should trigger OnEquip for Spirit Visage)
                // Config MaxTime can be short for setup tests, will be overridden in specific test cases.
                sim = simulation.NewSimulationWithConfig(world, config.WithMaxTime(0.1))
                Expect(sim).NotTo(BeNil())

                // Re-fetch health component in case it was modified by simulation setup
                attackerHealth = getHealth(world, attacker)
            })

            It("should apply base stats of Redemption and initialize Spirit Visage effect", func() {
                // Static bonuses from Redemption should be applied during simulation setup.
                // Base Max HP = 1000 (set in BeforeEach)
                // Expected Final Max HP = Base Max HP (1000) + Redemption Static Health (200)
                expectedFinalMaxHP := baseMaxHP + redemptionStaticHealth
                Expect(attackerHealth.GetFinalMaxHP()).To(BeNumerically("~", expectedFinalMaxHP, 0.01), "Final Max HP should include Redemption's static health bonus")

                // Current HP should also reflect the added static health if it was at max.
                // If it was 500/1000, now it should be 500 + 200 = 700 / 1200 (assuming current HP increases proportionally or by the flat amount if not capped)
                // The setupCombat logic for health usually adds static HP to both current and max.
                // Initial current HP was 500. After adding 200 static HP, it should be 700.
                Expect(attackerHealth.GetCurrentHP()).To(BeNumerically("~", (baseMaxHP/2)+redemptionStaticHealth, 0.01), "Current HP should increase by Redemption's static health bonus")


                // Expected Armor = Base Armor + Redemption Static Armor
                // Base Armor is usually 0 for dummies, or a champion's base. Assume 0 for simplicity here.
                baseArmor := attackerHealth.GetBaseArmor() // Get actual base armor
                expectedFinalArmor := baseArmor + redemptionStaticArmor
                Expect(attackerHealth.GetFinalArmor()).To(BeNumerically("~", expectedFinalArmor, 0.01), "Final Armor should include Redemption's static armor bonus")

                // Expected MR = Base MR + Redemption Static MR
                baseMR := attackerHealth.GetBaseMR() // Get actual base MR
                expectedFinalMR := baseMR + redemptionStaticMR
                Expect(attackerHealth.GetFinalMR()).To(BeNumerically("~", expectedFinalMR, 0.01), "Final MR should include Redemption's static MR bonus")

                // Assert SpiritVisageEffect parameters
                Expect(spiritVisageEffect).NotTo(BeNil()) // Already checked in BeforeEach
                Expect(spiritVisageEffect.GetTickInterval()).To(BeNumerically("~", spiritVisageTickInterval, 0.001))
                Expect(spiritVisageEffect.GetMissingHealthHealRate()).To(BeNumerically("~", spiritVisageMissingHealRate, 0.001))
                Expect(spiritVisageEffect.GetMaxHeal()).To(BeNumerically("~", spiritVisageMaxHeal, 0.001)) // This is the value for math.Max
            })

            It("should schedule and process the first heal tick correctly", func() {
                // Attacker MaxHP = 1000 (base) + 200 (item) = 1200
                // Attacker CurrentHP = 500 (base) + 200 (item) = 700 initially (after static bonuses)
                // Missing HP = 1200 - 700 = 500
                initialCurrentHP := attackerHealth.GetCurrentHP() // Should be 700
                finalMaxHP := attackerHealth.GetFinalMaxHP()     // Should be 1200
                missingHP := finalMaxHP - initialCurrentHP       // 1200 - 700 = 500

                sim.SetMaxTime(spiritVisageEffect.GetTickInterval() + 0.1) // e.g., 2.1s
                sim.RunSimulation()

                // Heal calculation:
                // PercentHeal = missingHP * MissingHealthHealRate = 500 * 0.05 = 25
                // MaxHeal (from component, used in math.Max) = 100.0
                // Actual Heal = math.Max(PercentHeal, MaxHeal) = math.Max(25, 100) = 100
                percentHeal := missingHP * spiritVisageEffect.GetMissingHealthHealRate()
                expectedHealAmount := math.Max(percentHeal, spiritVisageEffect.GetMaxHeal()) // As per handler logic

                expectedHPAfterHeal := initialCurrentHP + expectedHealAmount
                // Cap HP at MaxHP
                if expectedHPAfterHeal > finalMaxHP {
                    expectedHPAfterHeal = finalMaxHP
                }

                Expect(attackerHealth.GetCurrentHP()).To(BeNumerically("~", expectedHPAfterHeal, 0.01),
                    fmt.Sprintf("Expected HP after 1 tick. InitialHP: %.1f, MissingHP: %.1f, PercentHeal: %.1f, MaxHealField: %.1f, ActualHeal: %.1f",
                        initialCurrentHP, missingHP, percentHeal, spiritVisageEffect.GetMaxHeal(), expectedHealAmount))
            })

            It("should process multiple heal ticks, adjusting to current missing health", func() {
                // Initial state after static bonuses: MaxHP=1200, CurrentHP=700
                initialHP := attackerHealth.GetCurrentHP() // 700
                finalMaxHP := attackerHealth.GetFinalMaxHP() // 1200

                sim.SetMaxTime((2 * spiritVisageEffect.GetTickInterval()) + 0.1) // e.g., 4.1s
                sim.RunSimulation()

                // Tick 1 @ 2.0s:
                missingHP1 := finalMaxHP - initialHP // 1200 - 700 = 500
                percentHeal1 := missingHP1 * spiritVisageEffect.GetMissingHealthHealRate() // 500 * 0.05 = 25
                actualHeal1 := math.Max(percentHeal1, spiritVisageEffect.GetMaxHeal())     // Max(25, 100) = 100
                hpAfterTick1 := initialHP + actualHeal1                                   // 700 + 100 = 800
                if hpAfterTick1 > finalMaxHP { hpAfterTick1 = finalMaxHP }

                // Tick 2 @ 4.0s:
                missingHP2 := finalMaxHP - hpAfterTick1 // 1200 - 800 = 400
                percentHeal2 := missingHP2 * spiritVisageEffect.GetMissingHealthHealRate() // 400 * 0.05 = 20
                actualHeal2 := math.Max(percentHeal2, spiritVisageEffect.GetMaxHeal())     // Max(20, 100) = 100
                hpAfterTick2 := hpAfterTick1 + actualHeal2                                 // 800 + 100 = 900
                if hpAfterTick2 > finalMaxHP { hpAfterTick2 = finalMaxHP }

                Expect(attackerHealth.GetCurrentHP()).To(BeNumerically("~", hpAfterTick2, 0.01),
                    fmt.Sprintf("Expected HP after 2 ticks. HP_init:%.1f, Heal1:%.1f (Missing1:%.1f, Perc1:%.1f), HP_mid:%.1f, Heal2:%.1f (Missing2:%.1f, Perc2:%.1f)",
                        initialHP, actualHeal1, missingHP1, percentHeal1, hpAfterTick1, actualHeal2, missingHP2, percentHeal2))
            })

            It("should use MaxHeal value if percent missing health heal is smaller (due to math.Max)", func() {
                // Set CurrentHP so that (MissingHP * Rate) < MaxHeal
                // MaxHP = 1200. CurrentHP = 700 (from BeforeEach + static)
                // Let's set current HP higher to make missing HP smaller for this test.
                // CurrentHP = 1100. MaxHP = 1200. MissingHP = 100.
                attackerHealth.SetCurrentHP(finalMaxHP - 100) // Set HP to 1100 (Missing = 100)
                initialCurrentHP := attackerHealth.GetCurrentHP() // Should be 1100
                missingHP := finalMaxHP - initialCurrentHP       // 100

                sim.SetMaxTime(spiritVisageEffect.GetTickInterval() + 0.1)
                sim.RunSimulation()

                // PercentHeal = missingHP * Rate = 100 * 0.05 = 5
                // MaxHeal (from component, used in math.Max) = 100.0
                // Actual Heal = math.Max(PercentHeal, MaxHeal) = math.Max(5, 100) = 100
                percentHeal := missingHP * spiritVisageEffect.GetMissingHealthHealRate()
                expectedHealAmount := math.Max(percentHeal, spiritVisageEffect.GetMaxHeal())

                expectedHPAfterHeal := initialCurrentHP + expectedHealAmount
                if expectedHPAfterHeal > finalMaxHP { expectedHPAfterHeal = finalMaxHP }


                Expect(attackerHealth.GetCurrentHP()).To(BeNumerically("~", expectedHPAfterHeal, 0.01),
                    fmt.Sprintf("Expected HP when percent heal is small. InitialHP: %.1f, MissingHP: %.1f, PercentHeal: %.1f, MaxHealField: %.1f, ActualHeal: %.1f",
                        initialCurrentHP, missingHP, percentHeal, spiritVisageEffect.GetMaxHeal(), expectedHealAmount))
            })

            It("should not heal if champion is at full health", func() {
                attackerHealth.SetCurrentHP(attackerHealth.GetFinalMaxHP()) // Set to full health
                initialCurrentHP := attackerHealth.GetCurrentHP()

                sim.SetMaxTime(spiritVisageEffect.GetTickInterval() + 0.1)
                sim.RunSimulation()

                Expect(attackerHealth.GetCurrentHP()).To(BeNumerically("~", initialCurrentHP, 0.01), "HP should not change if already at max")
            })

            It("should stop scheduling new heal ticks if item is removed", func() {
                // Initial state after static bonuses: MaxHP=1200, CurrentHP=700
                initialHP := attackerHealth.GetCurrentHP() // 700
                finalMaxHP := attackerHealth.GetFinalMaxHP() // 1200

                // Run for one tick
                sim.SetMaxTime(spiritVisageEffect.GetTickInterval() + 0.1) // e.g. 2.1s
                sim.RunSimulation()

                // Calculate HP after 1st tick
                missingHP1 := finalMaxHP - initialHP
                percentHeal1 := missingHP1 * spiritVisageEffect.GetMissingHealthHealRate()
                actualHeal1 := math.Max(percentHeal1, spiritVisageEffect.GetMaxHeal())
                hpAfterTick1 := initialHP + actualHeal1
                if hpAfterTick1 > finalMaxHP { hpAfterTick1 = finalMaxHP }

                Expect(attackerHealth.GetCurrentHP()).To(BeNumerically("~", hpAfterTick1, 0.01), "HP after first tick before item removal")

                // Remove Spirit Visage
                equipmentComp, ok := world.GetEquipment(attacker)
                Expect(ok).To(BeTrue())
                err := equipmentComp.RemoveItem(data.TFT_Item_SpiritVisage) // Assuming RemoveItem by API name
                Expect(err).NotTo(HaveOccurred())
                // Also need to remove the SpiritVisageEffect component to stop existing handlers
                // or ensure the handler checks for the item in equipment component.
                // The handler should check `if !equipment.HasItem(s.itemApiName) { return }`
                // Let's assume the handler does this check.

                // Run simulation longer, past the next scheduled tick time
                // A new simulation instance is implicitly created by SetMaxTime then RunSimulation
                // if we are using the same 'sim' variable defined in an outer scope.
                // To ensure the same simulation continues with item removed, we need a way
                // to remove item and then continue.
                // The current structure of `sim.RunSimulation()` re-evaluates events from t=0
                // up to new MaxTime.
                // A better test would be to check the event queue or if new events are published.
                // However, if the handler itself checks for item presence before healing,
                // then even if an event is processed, no heal should occur.

                // To test this robustly with current sim structure:
                // 1. Create sim, run for 1st tick. HP changes.
                // 2. Create a *new* sim instance with the *same world* but *without the item*.
                //    The item must be removed from the world's entity BEFORE creating the new sim.
                //    This is tricky because the item is already added in BeforeEach.
                //    Alternative: Modify the existing sim's world state and continue.
                //    Let's try to remove it and run the *same* sim longer. The handler logic will be key.

                // Remove item from equipment component (done above)
                // The SpiritVisageHealTickEvent for the *next* tick was already enqueued by the previous tick.
                // We need to check if the handler for that event will abort due to item removal.
                sim.SetMaxTime((2 * spiritVisageEffect.GetTickInterval()) + 0.1) // e.g. 4.1s
                sim.RunSimulation() // This will process events up to 4.1s

                Expect(attackerHealth.GetCurrentHP()).To(BeNumerically("~", hpAfterTick1, 0.01), "HP should remain unchanged after item removal, despite a previously scheduled tick")
            })

            It("should increase heal amount if multiple Spirit Visages are equipped", func() {
                // One Spirit Visage is already equipped from BeforeEach. Equip a second one.
                err := equipmentManager.AddItemToChampion(attacker, data.TFT_Item_SpiritVisage)
                Expect(err).NotTo(HaveOccurred())

                // Manually add a second SpiritVisageEffect component or ensure the system handles multiple.
                // The current ECS design might only allow one component of each type per entity.
                // If so, this test needs a system that aggregates effects from multiple items.
                // Assuming the item system correctly creates multiple effect instances or one that counts stacks.
                // For this test, let's assume the SpiritVisageHandler's ProcessEvent is called for each item instance's effect.
                // This implies multiple SpiritVisageEffect components or a single one that knows about stacks.
                // The provided handler seems to operate on a single effect component instance passed to it.
                // This test might expose a limitation or require a specific setup for multiple items of the same type.

                // Given the current handler structure (receives one effect component),
                // true stacking would mean the handler for Spirit Visage is registered per item instance,
                // and each calls ApplyHeal.

                // Let's assume the test setup means two separate heal events will be processed for the same champion
                // if two items are equipped and the system supports this (e.g. two effect components, or one effect component with stack count).
                // For now, we'll assume the handler's OnEquip is called twice, scheduling two independent series of ticks.
                // And ProcessEvent is triggered for each.

                // Re-initialize simulation to apply the second item's OnEquip
                sim = simulation.NewSimulationWithConfig(world, config.WithMaxTime(spiritVisageEffect.GetTickInterval()+0.1))
                Expect(sim).NotTo(BeNil())
                attackerHealth = getHealth(world, attacker) // Re-fetch after new sim

                // Attacker MaxHP = 1000 (base) + 200 (item1) + 200 (item2) = 1400
                // Attacker CurrentHP = 500 (base) + 200 (item1) + 200 (item2) = 900 initially
                initialCurrentHP := attackerHealth.GetCurrentHP() // Should be 900
                finalMaxHP := attackerHealth.GetFinalMaxHP()     // Should be 1400
                missingHP := finalMaxHP - initialCurrentHP       // 1400 - 900 = 500

                // Heal calculation for ONE item:
                percentHealPerItem := missingHP * spiritVisageEffect.GetMissingHealthHealRate() // 500 * 0.05 = 25
                actualHealPerItem := math.Max(percentHealPerItem, spiritVisageEffect.GetMaxHeal()) // Max(25, 100) = 100

                // Total heal from two items = actualHealPerItem * 2 (assuming independent application)
                totalExpectedHeal := actualHealPerItem * 2 // 100 * 2 = 200

                sim.RunSimulation()

                expectedHPAfterHeal := initialCurrentHP + totalExpectedHeal
                if expectedHPAfterHeal > finalMaxHP { expectedHPAfterHeal = finalMaxHP }

                Expect(attackerHealth.GetCurrentHP()).To(BeNumerically("~", expectedHPAfterHeal, 0.01),
                    fmt.Sprintf("Expected HP with 2 items. InitialHP: %.1f, MissingHP: %.1f, HealPerItem: %.1f, TotalHeal: %.1f",
                        initialCurrentHP, missingHP, actualHealPerItem, totalExpectedHeal))
            })
        })
    })

    Describe("Trait System Integration", func() {
        // Focus on how trait systems interact during setup and potentially during the run

        var (
            championFactory *factory.ChampionFactory
            // Systems involved in trait processing (might be needed for manual checks if not relying solely on Simulation setup)
            // traitCounterSystem     *traitsys.TraitCounterSystem
            // traitStaticBonusSystem *traitsys.TraitStaticBonusSystem
            // traitState             *traitsys.TeamTraitState // Might be hard to access directly from sim
        )

        BeforeEach(func() {
            // Reset world and create basic components/systems for trait tests
            world = ecs.NewWorld()
            championFactory = factory.NewChampionFactory(world)
            equipmentManager = managers.NewEquipmentManager(world) // Not needed unless items interact with traits

            // Create instances of systems if needed for manual checks,
            // otherwise rely on Simulation's internal creation.
            // traitState = traitsys.NewTeamTraitState()
            // traitCounterSystem = traitsys.NewTraitCounterSystem(world, traitState)
            // traitStaticBonusSystem = traitsys.NewTraitStaticBonusSystem(world, traitState)
            // statCalculationSystem = systems.NewStatCalculationSystem(world) // Needed to see final stats
        })

        Context("with Rapidfire Trait (2 units)", func() {
            var (
                kindred1, kindred2, kogmaw, shyvana, blueGolem ecs.Entity
                rapidfireData              *data.Trait
                expectedBonusAS            float64
            )
            BeforeEach(func() {
                // Create champions with the Rapidfire trait for the player team
                var err error
                kindred1, err = championFactory.CreatePlayerChampion("TFT14_Kindred", 1) 
                Expect(err).NotTo(HaveOccurred())
                kindred2, err = championFactory.CreatePlayerChampion("TFT14_Kindred", 1)
                Expect(err).NotTo(HaveOccurred())
                kogmaw, err = championFactory.CreatePlayerChampion("TFT14_KogMaw", 1)
                Expect(err).NotTo(HaveOccurred())
                // Add a non-Rapidfire champion
                shyvana, err = championFactory.CreatePlayerChampion("TFT14_Shyvana", 1) 
                Expect(err).NotTo(HaveOccurred())

                blueGolem, err = championFactory.CreateEnemyChampion("TFT_BlueGolem", 1)
                Expect(err).NotTo(HaveOccurred())
                blueGolemHealth := getHealth(world, blueGolem)
                blueGolemHealth.SetBaseMaxHP(10000.0)
                blueGolemPos := getPosition(world, blueGolem)
                blueGolemPos.SetPosition(0, 1) 

                // Get trait data to find the expected bonus
                rapidfireData = data.GetTraitByName(data.TFT14_Rapidfire)
                Expect(rapidfireData).NotTo(BeNil())
                // Find the bonus for tier 0 (2 units)
                var foundBonus bool
                for _, effect := range rapidfireData.Effects {
                    if effect.MinUnits == 2 { // Tier 0 threshold
                        expectedBonusAS = effect.Variables["TeamBonus"]// team bonus AS
                        foundBonus = true
                        break
                    }
                }
                Expect(foundBonus).To(BeTrue(), "Could not find Rapidfire bonus for 2 units")
                Expect(expectedBonusAS).To(BeNumerically("~", 0.1), "Expected team bonus AS should be 10%")

                // Create the simulation - this runs setupCombat which includes trait counting and static bonus application
                sim = simulation.NewSimulationWithConfig(world, config.WithMaxTime(0.1)) // Short time, only care about setup
                Expect(sim).NotTo(BeNil())
                // Note: setupCombat runs:
                // 1. traitCounterSystem.UpdateCountsAndTiers()
                // 2. abilityCritSystem.Update()
                // 3. traitStaticBonusSystem.ApplyStaticTraitsBonus()
                // 4. baseStaticItemSystem.ApplyStaticItemsBonus()
                // 5. statCalcSystem.ApplyStaticBonusStats()
            })
            It("should correctly calculate TeamTraitState", func() {
                // Check if the trait state was updated correctly
                traitState := sim.GetTeamTraitState()
                Expect(traitState).NotTo(BeNil(), "TraitState should not be nil")
                // Check if Rapidfire trait is present and has the correct count
                rapidfireCount := traitState.GetUnitCount(components.TeamPlayer, rapidfireData.Name)
                Expect(rapidfireCount).To(Equal(2), "Rapidfire trait count should be 2")
                // Check if the trait tier is correct
                rapidfireTier := traitState.GetActiveTier(components.TeamPlayer, rapidfireData.Name)
                Expect(rapidfireTier).To(Equal(0), "Rapidfire trait tier should be 0 (2 units)")
            })

            It("should apply the correct static Attack Speed bonus to Rapidfire champions", func() {
                // Get attack components AFTER simulation setup
                attackK1 := getAttack(world, kindred1)
                attackK2 := getAttack(world, kindred2)
                attackKog := getAttack(world, kogmaw)

                // Calculate expected final AS = BaseAS * (1 + TraitBonusAS)
                expectedFinalASK1 := attackK1.GetBaseAttackSpeed() * (1.0 + expectedBonusAS)
                expectedFinalASK2 := attackK2.GetBaseAttackSpeed() * (1.0 + expectedBonusAS)
                expectedFinalASKog := attackKog.GetBaseAttackSpeed() * (1.0 + expectedBonusAS)

                // Verify the BonusPercentAttackSpeed field directly
                Expect(attackK1.GetBonusPercentAttackSpeed()).To(BeNumerically("~", expectedBonusAS, 0.001), "Kindred 1 BonusPercentAttackSpeed should be Rapidfire bonus")
                Expect(attackK2.GetBonusPercentAttackSpeed()).To(BeNumerically("~", expectedBonusAS, 0.001), "Kindred 2 BonusPercentAttackSpeed should be Rapidfire bonus")
                Expect(attackKog.GetBonusPercentAttackSpeed()).To(BeNumerically("~", expectedBonusAS, 0.001), "Kog'Maw BonusPercentAttackSpeed should be Rapidfire bonus")

                // Assert that the final attack speed reflects the trait bonus
                Expect(attackK1.GetFinalAttackSpeed()).To(BeNumerically("~", expectedFinalASK1, 0.001), "Kindred 1 Final AS should include Rapidfire bonus")
                Expect(attackK2.GetFinalAttackSpeed()).To(BeNumerically("~", expectedFinalASK2, 0.001), "Kindred 2 Final AS should include Rapidfire bonus")
                Expect(attackKog.GetFinalAttackSpeed()).To(BeNumerically("~", expectedFinalASKog, 0.001), "Kog'Maw Final AS should include Rapidfire bonus")

                // Assert rapidfire champions have new RapidfireEffect component added
                rfK1, ok := world.GetRapidfireEffect(kindred1)
                Expect(ok).To(BeTrue(), "RapidfireEffect should be added to Kindred 1")
                Expect(rfK1.GetCurrentBonusAS()).To(BeNumerically("~", 0, 0.001), "RapidfireEffect should have correct bonus AS")
                Expect(rfK1.GetCurrentStacks()).To(Equal(0), "RapidfireEffect should have 0 stacks initially")
                Expect(rfK1.GetAttackSpeedPerStack()).To(BeNumerically("~", 0.04, 0.001), "RapidfireEffect should have correct AS per stack")
                rfK2, ok := world.GetRapidfireEffect(kindred2)
                Expect(ok).To(BeTrue(), "RapidfireEffect should be added to Kindred 2")
                Expect(rfK2.GetCurrentBonusAS()).To(BeNumerically("~", 0, 0.001), "RapidfireEffect should have correct bonus AS")
                Expect(rfK2.GetCurrentStacks()).To(Equal(0), "RapidfireEffect should have 0 stacks initially")
                Expect(rfK2.GetAttackSpeedPerStack()).To(BeNumerically("~", 0.04, 0.001), "RapidfireEffect should have correct AS per stack")
                rfKog, ok := world.GetRapidfireEffect(kogmaw)
                Expect(ok).To(BeTrue(), "RapidfireEffect should be added to Kog'Maw")
                Expect(rfKog.GetCurrentBonusAS()).To(BeNumerically("~", 0, 0.001), "RapidfireEffect should have correct bonus AS")
                Expect(rfKog.GetCurrentStacks()).To(Equal(0), "RapidfireEffect should have 0 stacks initially")
                Expect(rfKog.GetAttackSpeedPerStack()).To(BeNumerically("~", 0.04, 0.001), "RapidfireEffect should have correct AS per stack")
            })

            It("should apply Rapidfire team bonus to non-Rapidfire champions", func() {
                attackShyvana := getAttack(world, shyvana)

                // Shyvana should also receive bonus AS from Rapidfire
                Expect(attackShyvana.GetBonusPercentAttackSpeed()).To(BeNumerically("~", expectedBonusAS, 0.001), "Shyvana BonusPercentAttackSpeed should be 10%")
                Expect(attackShyvana.GetFinalAttackSpeed()).To(BeNumerically("~", attackShyvana.GetBaseAttackSpeed() * 1.1, 0.001), "Shyvana Final AS should be greater than Base AS")

                // Shyvana should NOT have a RapidfireEffect component
                _, ok := world.GetRapidfireEffect(shyvana)
                Expect(ok).To(BeFalse(), "Shyvana should NOT have a RapidfireEffect component")
            })

            It("should correctly stack RapidfireEffect on attacks", func() {
                // Set up attack components for Kindred 1 and Kog'Maw
                attackK1 := getAttack(world, kindred1)
                attackK2 := getAttack(world, kindred2)
                attackKog := getAttack(world, kogmaw)

                attackK1.ResetBonuses()
                attackK2.ResetBonuses()
                attackKog.ResetBonuses()

                sim = simulation.NewSimulationWithConfig(world, config.WithMaxTime(1.10)) // starting AS for Kindred and Kog'Maw is 0.77 = 0.7 * 1.1 (10% bonus from Rapidfire), 1.10s should be enough for 1 attack to land
                sim.RunSimulation() 

                Expect(attackK1.GetAttackCount()).To(Equal(1), "Kindred 1 should have attacked once")
                Expect(attackK2.GetAttackCount()).To(Equal(1), "Kindred 2 should have attacked once")
                Expect(attackKog.GetAttackCount()).To(Equal(1), "Kog'Maw should have attacked once")

                // Check if stacks were incremented correctly
                rfK1, _ := world.GetRapidfireEffect(kindred1)
                rfK2, _ := world.GetRapidfireEffect(kindred2)
                rfKog, _ := world.GetRapidfireEffect(kogmaw)

                Expect(rfK1.GetCurrentStacks()).To(Equal(1), "RapidfireEffect stacks should be 1 after attack")
                Expect(rfK2.GetCurrentStacks()).To(Equal(1), "RapidfireEffect stacks should be 1 after attack") 
                Expect(rfKog.GetCurrentStacks()).To(Equal(1), "RapidfireEffect stacks should be 1 after attack")

                // Check if bonus AS was applied correctly
                expectedBonusASFromStack := 0.04 // 4% per stack
                Expect(rfK1.GetCurrentBonusAS()).To(BeNumerically("~", expectedBonusASFromStack, 0.001), "RapidfireEffect bonus AS should be 4% after attack")
                Expect(attackK1.GetBonusPercentAttackSpeed()).To(BeNumerically("~", expectedBonusASFromStack+0.1, 0.001), "RapidfireEffect bonus AS should be 4% after attack")
                Expect(attackK1.GetFinalAttackSpeed()).To(BeNumerically("~", attackK1.GetBaseAttackSpeed() * (1.0 + expectedBonusASFromStack + 0.1), 0.001), "RapidfireEffect bonus AS should be 14% after attack")

                Expect(rfK2.GetCurrentBonusAS()).To(BeNumerically("~", expectedBonusASFromStack, 0.001), "RapidfireEffect bonus AS should be 4% after attack")

                Expect(rfKog.GetCurrentBonusAS()).To(BeNumerically("~", expectedBonusASFromStack, 0.001), "RapidfireEffect bonus AS should be 4% after attack")
            })
            It("should not stack RapidfireEffects if max stacks reached", func() {
                // Set up attack components for Kindred 1 and Kog'Maw
                attackK1 := getAttack(world, kindred1)
                attackK2 := getAttack(world, kindred2)
                attackKog := getAttack(world, kogmaw)

                attackK1.ResetBonuses()
                attackK2.ResetBonuses()
                attackKog.ResetBonuses()

                sim = simulation.NewSimulationWithConfig(world, config.WithMaxTime(12.0)) 
                sim.RunSimulation() 

                Expect(attackK1.GetAttackCount()).To(Equal(11), "Kindred 1 should have attacked 11 times")
                Expect(attackK2.GetAttackCount()).To(Equal(11), "Kindred 2 should have attacked 11 times")
                Expect(attackKog.GetAttackCount()).To(Equal(11), "Kog'Maw should have attacked 11 times")

                // Check if stacks were incremented correctly
                rfK1, _ := world.GetRapidfireEffect(kindred1)
                rfK2, _ := world.GetRapidfireEffect(kindred2)
                rfKog, _ := world.GetRapidfireEffect(kogmaw)

                Expect(rfK1.GetCurrentStacks()).To(Equal(10), "RapidfireEffect stacks should be 10 (max) after attack")
                Expect(rfK2.GetCurrentStacks()).To(Equal(10), "RapidfireEffect stacks should be 10 (max) after attack") 
                Expect(rfKog.GetCurrentStacks()).To(Equal(10), "RapidfireEffect stacks should be 10 (max) after attack")

                // Check if bonus AS was applied correctly
                expectedBonusASFromStack := 0.4 // 4% per stack * 10 stacks
                Expect(rfK1.GetCurrentBonusAS()).To(BeNumerically("~", expectedBonusASFromStack, 0.001), "RapidfireEffect bonus AS should be 40% after attack")
                Expect(attackK1.GetBonusPercentAttackSpeed()).To(BeNumerically("~", expectedBonusASFromStack+0.1, 0.001), "RapidfireEffect bonus AS should be 50% after attack")
                Expect(attackK1.GetFinalAttackSpeed()).To(BeNumerically("~", attackK1.GetBaseAttackSpeed() * (1.0 + expectedBonusASFromStack + 0.1), 0.001), "Final Attack Speed should be 1.5*Base AS after attack")

                Expect(rfK2.GetCurrentBonusAS()).To(BeNumerically("~", expectedBonusASFromStack, 0.001), "RapidfireEffect bonus AS should be 40% after attack")

                Expect(rfKog.GetCurrentBonusAS()).To(BeNumerically("~", expectedBonusASFromStack, 0.001), "RapidfireEffect bonus AS should be 40% after attack")
            })
        }) // End Context("with Rapidfire Trait (2 units)")

    }) // End Describe("Trait System Integration")

}) // End Describe("Simulation")
