package simulation_test

import (
	"fmt"
	"math"

	"github.com/suriz/tft-dps-simulator/components"
	"github.com/suriz/tft-dps-simulator/data"
	"github.com/suriz/tft-dps-simulator/ecs"
	"github.com/suriz/tft-dps-simulator/factory"
	"github.com/suriz/tft-dps-simulator/managers"
	"github.com/suriz/tft-dps-simulator/simulation"

	// Need systems for manual stat calc check
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
			WithTimeStep(0.1).
			WithMaxTime(1.0).    // Short max time for faster tests
			WithDebugMode(false) // Disable debug output by default

		championFactory = factory.NewChampionFactory(world)
		// Load real item data for the manager
		// Ensure item data is loaded correctly in your actual setup, maybe in a BeforeSuite
		// For this test, we assume item data is available.
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
			targetAttack.SetFinalAttackSpeed(0)
		}

		// Initialize sim here for tests that don't need specific pre-item setup
		sim = simulation.NewSimulationWithConfig(world, config)
		Expect(sim).NotTo(BeNil()) // Ensure sim is created
	})

	Describe("Initialization", func() {
		It("should create a simulation with default config", func() {
			defaultSim := simulation.NewSimulation(world)
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
			invalidConfig := simulation.DefaultConfig().WithTimeStep(0)
			Expect(func() { simulation.NewSimulationWithConfig(world, invalidConfig) }).To(Panic())
		})
	})

	Describe("Configuration", func() {
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
			invalidConfig := config.WithMaxTime(0)
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
			Expect(sim.GetConfig().TimeStep).To(Equal(0.05))
		})
	})

	// Removed the Describe("Step Method") block as we test its effects via RunSimulation

	Describe("RunSimulation Method", func() {
		// Uses the 'sim' instance created in the outer BeforeEach,
		// which has the real event bus and systems wired up internally.

		It("should run until MaxTime is reached and apply effects", func() {
			// Set simulation time long enough for at least one attack cycle
			sim.SetMaxTime(2.5) // AS=0.5, attack should land around t=2.0

			// Get initial state
			targetHealth := getHealth(world, target)
			initialHP := targetHealth.GetCurrentHP()

			Expect(initialHP).To(Equal(targetMaxHP)) // Should start at max HP

			// Run the simulation
			sim.RunSimulation()

			// Check effects after running
			// Target should have taken damage
			// Damage calculation: BaseAD=50. Assume 0 crit/amp/durability, Armor=0 (dummy default) -> Final=50
			expectedHPAfterOneAttack := targetMaxHP - 50.0
			Expect(targetHealth.GetCurrentHP()).To(BeNumerically("~", expectedHPAfterOneAttack, 0.1), "Target HP should decrease after attack")
		})

		It("should execute initial system updates before the loop", func() {
			// Check if stats are calculated correctly *before* RunSimulation starts.
			// NewSimulationWithConfig should have run initial updates.
			attackerAttack, ok := world.GetAttack(attacker)
			Expect(ok).To(BeTrue())
			// Check if a final stat (which depends on calculation) has a non-zero value
			// Base AD was 50, Final AD should be 50 after initial calculation.
			Expect(attackerAttack.GetFinalAD()).To(BeNumerically("~", 50.0, 0.1))
			Expect(attackerAttack.GetFinalAttackSpeed()).To(BeNumerically("~", 0.5, 0.01))

			// Now run the simulation - this implicitly tests that RunSimulation doesn't break pre-calculated stats
			sim.RunSimulation()
			Expect(attackerAttack.GetFinalAD()).To(BeNumerically("~", 50.0, 0.1)) // Should still be calculated
		})

		It("should handle multiple attack cycles correctly", func() {
			sim.SetMaxTime(4.5) // Enough time for attacks at t=1.0 and t=2.0

			targetHealth := getHealth(world, target)

			sim.RunSimulation()

			// Check effects after two attacks
			expectedHPAfterTwoAttacks := targetMaxHP - 50.0 - 50.0
			Expect(targetHealth.GetCurrentHP()).To(BeNumerically("~", expectedHPAfterTwoAttacks, 0.1))
		})

		// Optional: Test DebugMode output
		PIt("should print debug messages when DebugMode is enabled", func() {
			// This test requires capturing stdout, which can be complex.
			// Using gbytes.Say is one way with Ginkgo V2.
			buffer := gbytes.NewBuffer() // Capture output
			// Need a way to redirect stdout or pass the buffer to the simulation logger

			debugConfig := config.WithDebugMode(true).WithReportingInterval(0.5).WithMaxTime(0.6)
			// Create a new sim instance for this specific test
			debugSim := simulation.NewSimulationWithConfig(world, debugConfig)
			// TODO: Inject buffer into logger if possible

			debugSim.RunSimulation()

			Eventually(buffer).Should(gbytes.Say("Simulation time: 0.5s"))
		})

		Context("with Dynamic Time Items", func() {

			It("should apply Quicksilver initial AS bonus, allow attacks, and let bonus expire", func() {
				// Add Quicksilver to attacker
				err := equipmentManager.AddItemToChampion(attacker, data.TFT_Item_Quicksilver)
				Expect(err).NotTo(HaveOccurred())

				// Create simulation AFTER adding the item
				sim = simulation.NewSimulationWithConfig(world, config)
				Expect(sim).NotTo(BeNil())

				attackerAttack := getAttack(world, attacker)
				targetHealth := getHealth(world, target)
				baseAS := attackerAttack.GetBaseAttackSpeed()
				Expect(baseAS).To(BeNumerically("~", 0.5), "Base AS should be 0.5")
				finalStaticAS := attackerAttack.GetFinalAttackSpeed()
				Expect(finalStaticAS).To(BeNumerically("~", 0.65, 0.01), "Final AS should be 0.65 after static QS bonus")

				// Run simulation for a short time (e.g., 5 seconds)
				shortDuration := 5.0
				sim.SetMaxTime(shortDuration)
				sim.RunSimulation()

				// Verify attacks happened and AS bonus was active
				Expect(targetHealth.CurrentHP).To(BeNumerically("<", targetMaxHP), "Target should have taken damage within %f seconds", shortDuration)

				expectedASWithQS := baseAS * (1 + 0.3 + 2*0.03)
				Expect(attackerAttack.GetFinalAttackSpeed()).To(BeNumerically("~", expectedASWithQS, 0.01), "AS should still have QS bonus active")
				fmt.Fprintf(GinkgoWriter, "Debug QS Test (t=0): Base=%.2f, Expected=%.2f, Actual=%.2f\n", baseAS, expectedASWithQS, attackerAttack.GetFinalAttackSpeed())

				healthAfterShortRun := targetHealth.CurrentHP

				attackerAttack.ResetBonuses()
				quicksilverEffect, ok := world.GetQuicksilverEffect(attacker)
				Expect(ok).To(BeTrue(), "Quicksilver effect component should be present")
				quicksilverEffect.ResetEffects()
				
				// attackerAttack.SetBonusPercentAttackSpeed(0.3)
				// Reset simulation time and run past expiry (QS duration is 18)
				longDuration := 19.0
				// Re-create sim with the *same world state* but reset time/systems for a clean run
				// NOTE: Recreating the sim resets internal system states but KEEPS the world components (like the added item and effect)
				sim = simulation.NewSimulationWithConfig(world, config)
				sim.SetMaxTime(longDuration)
				sim.RunSimulation()

				// Verify AS bonus persists after expiry (Static + Max Dynamic Stacks)
				// Static = 0.3, Max Dynamic = 9 stacks * 0.03 = 0.27
				expectedASAfterExpiry := baseAS * (1 + 0.3 + 0.27) // Additive stacking
				Expect(attackerAttack.GetFinalAttackSpeed()).To(BeNumerically("~", expectedASAfterExpiry, 0.01), "Final AS should reflect base + static + fully stacked dynamic bonus after QS expiry")
				fmt.Fprintf(GinkgoWriter, "Debug QS Test (t=19s): Base=%.2f, Expected (Additive, Post-Expiry)=%.3f, Actual=%.3f\n", baseAS, expectedASAfterExpiry, attackerAttack.GetFinalAttackSpeed())

				// Verify QS effect component state
				quicksilverEffects, ok := world.GetQuicksilverEffect(attacker)
				Expect(ok).To(BeTrue(), "Quicksilver effect component should still be present")
				Expect(quicksilverEffects).NotTo(BeNil(), "Quicksilver effect should not be nil")
				Expect(quicksilverEffects.IsActive()).To(BeFalse(), "Quicksilver effect should be inactive after expiry")
				Expect(quicksilverEffects.GetRemainingDuration()).To(BeNumerically("<=", 0.0), "Quicksilver effect should have <= 0 remaining duration after expiry")
				// Stacks accumulate until expiry (2, 4, ..., 18s -> 9 stacks) and persist
				Expect(quicksilverEffects.GetStacks()).To(Equal(9), "Quicksilver effect should retain its 9 accumulated stacks after expiry")

				// Verify attacks continued after expiry
				Expect(targetHealth.CurrentHP).To(BeNumerically("<", healthAfterShortRun), "Target should have taken more damage after QS expired")
			})

			It("should stack Archangel's Staff AP over time", func() {
				// Add Archangel's Staff to attacker
				err := equipmentManager.AddItemToChampion(attacker, "TFT_Item_ArchangelsStaff")
				Expect(err).NotTo(HaveOccurred())

				// Create simulation AFTER adding the item
				sim = simulation.NewSimulationWithConfig(world, config)
				Expect(sim).NotTo(BeNil())

				attackerSpell := getSpell(world, attacker)
				initialAP := attackerSpell.GetBaseAP() // Get base AP

				// Check initial AP immediately after simulation creation
				// Archangel's might grant some initial AP besides stacking
				apAfterCreation := attackerSpell.GetFinalAP()
				Expect(apAfterCreation).To(BeNumerically(">=", initialAP), "Initial AP should include base AP and any initial AP bonus")

				// --- Run for ~6 seconds (expect 1 stack) ---
				sim = simulation.NewSimulationWithConfig(world, config) // Reset sim
				sim.SetMaxTime(6.0)
				sim.RunSimulation()
				apAfter6Sec := attackerSpell.GetFinalAP()
				// We expect AP to be higher than initial AP after 1 stack interval (5s)
				Expect(apAfter6Sec).To(BeNumerically(">", apAfterCreation), "AP should increase after first stack interval (~5s)")

				// --- Run for ~11 seconds (expect 2 stacks) ---
				sim = simulation.NewSimulationWithConfig(world, config) // Reset sim
				sim.SetMaxTime(11.0)
				sim.RunSimulation()
				apAfter11Sec := attackerSpell.GetFinalAP()
				// We expect AP after 11s to be higher than AP after 6s
				Expect(apAfter11Sec).To(BeNumerically(">", apAfter6Sec), "AP should increase further after second stack interval (~10s)")

			})
		})
	})

	Describe("PrintResults", func() {
		// This primarily tests output formatting.
		It("should print stats for specified entities", func() {
			// Requires capturing stdout.
			buffer := gbytes.NewBuffer()
			// TODO: Redirect stdout or inject buffer

			// Run for a short time to have some state changes
			sim.SetMaxTime(1.5)
			sim.RunSimulation()

			// Call the method - Need a way to pass the buffer for capture
			// sim.PrintResults(buffer, attacker, target) // Hypothetical modification

			Skip("Skipping stdout capture test for PrintResults") // Skip until capture is implemented

			// Check if output contains expected names or stats (adjust based on actual output format)
			Eventually(buffer).Should(gbytes.Say("TFT_BlueGolem")) // Check if champion name is printed
			Eventually(buffer).Should(gbytes.Say("TFT_TrainingDummy"))
			Eventually(buffer).Should(gbytes.Say("HP:")) // Check for stat labels
			Eventually(buffer).Should(gbytes.Say("AD:"))
		})
	})
})
