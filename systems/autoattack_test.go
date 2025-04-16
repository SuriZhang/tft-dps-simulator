package systems_test

import (
	"github.com/suriz/tft-dps-simulator/components"
	"github.com/suriz/tft-dps-simulator/ecs"
	"github.com/suriz/tft-dps-simulator/factory" // Import factory
	"github.com/suriz/tft-dps-simulator/systems"
	// "github.com/suriz/tft-dps-simulator/utils"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("AutoAttackSystem", func() {
	var (
		world            *ecs.World
		championFactory  *factory.ChampionFactory
		autoAttackSystem *systems.AutoAttackSystem
		player           ecs.Entity // Blue Golem
		target           ecs.Entity // Training Dummy
		playerAttack     *components.Attack
		targetHealth     *components.Health
		ok               bool
		// Define the pre-calculated expected damage based on provided stats
		// *** IMPORTANT: Recalculate this value if stats in BeforeEach change! ***
		// Calculation: 55 * ((1-0.25) + 0.25*1.4) * (1+0/100) * (100/(100+30)) * (1-0) = 46.538...
		preCalculatedExpectedDamage float64 = 46.54
	)

	BeforeEach(func() {

		world = ecs.NewWorld()
		championFactory = factory.NewChampionFactory(world) // Create factory
		autoAttackSystem = systems.NewAutoAttackSystem(world)

		// --- Create Player (Blue Golem) using Factory ---
		var err error
		// Use the correct API Name for Blue Golem from your data files
		player, err = championFactory.CreatePlayerChampion("TFT_BlueGolem", 1)
		Expect(err).NotTo(HaveOccurred())
		world.AddComponent(player, components.NewPosition(0,1)) 

		// utils.PrintChampionStats(world, player) 

		// --- Create Target (Training Dummy) using Factory ---
		// Use the correct API Name for Training Dummy
		target, err = championFactory.CreateEnemyChampion("TFT_TrainingDummy", 1)
		Expect(err).NotTo(HaveOccurred())
		world.AddComponent(target, components.NewPosition(1,1))
		// utils.PrintChampionStats(world, target)

		// --- Get Components ---
		playerAttack, ok = world.GetAttack(player)
		Expect(ok).To(BeTrue(), "Player should have Attack component")
		Expect(playerAttack).NotTo(BeNil(), "Player Attack component pointer should not be nil after retrieval")
		targetHealth, ok = world.GetHealth(target)
		Expect(ok).To(BeTrue(), "Target should have Health component")
		Expect(targetHealth).NotTo(BeNil(), "Target Health component pointer should not be nil after retrieval")

		// --- Manually Set Final Stats (as StatCalculationSystem isn't run) ---
		// This ensures attack speed calculations in the test are correct.
		// The factory sets base stats, but we need final stats for the system.
		playerAttack.SetFinalAD(55.00)
		playerAttack.SetFinalAttackSpeed(0.550)
		playerAttack.SetFinalRange(1.00)
		playerAttack.SetFinalCritChance(0.25) // Add other relevant final stats
		playerAttack.SetFinalCritMultiplier(1.40)
		playerAttack.SetFinalDamageAmp(0.00)

		targetHealth.SetFinalArmor(30.00)
		targetHealth.SetFinalMR(30.00) // Set even if not used by AD calc
		targetHealth.SetFinalDurability(0.00)
		// Set CurrentHP high enough for multiple hits
		targetHealth.SetCurrentHealth(550.00) // Use target's base max HP

		// --- Set Positions ---
		playerPosition, ok := world.GetPosition(player)
		Expect(ok).To(BeTrue())
		Expect(playerPosition).NotTo(BeNil())
		playerPosition.X = 0
		playerPosition.Y = 0

		targetPosition, ok := world.GetPosition(target)
		Expect(ok).To(BeTrue())
		Expect(targetPosition).NotTo(BeNil())
		targetPosition.X = 1 // Place within default range (adjust if needed based on actual range)
		targetPosition.Y = 0

		// --- Ensure Dummy doesn't attack ---
		targetAttack, ok := world.GetAttack(target)
		if ok { // Training Dummy might not even have an Attack component depending on factory setup
			Expect(targetAttack).NotTo(BeNil())
			targetAttack.SetBaseAttackSpeed(0)
			targetAttack.SetFinalAttackSpeed(0)
		}
	})

	// --- Contexts and It blocks remain the same ---
	Context("when a valid target is in range", func() {
		It("should perform an attack after the initial delay", func() {
			Expect(playerAttack).NotTo(BeNil(), "playerAttack should not be nil at the start of the test")
			Expect(targetHealth).NotTo(BeNil(), "targetHealth should not be nil at the start of the test")

			initialTargetHP := targetHealth.GetCurrentHP()
			// Ensure FinalAttackSpeed is not zero before calculating delay

			GinkgoWriter.Println("*********** Player Attack Speed:", playerAttack.GetFinalAttackSpeed())

			Expect(playerAttack.GetFinalAttackSpeed()).To(BeNumerically(">", 0), "Player FinalAttackSpeed must be > 0 for attack delay calculation")
			attackDelay := 1.0 / playerAttack.GetFinalAttackSpeed() // Time for one attack cycle
			dt := 0.1                                      // Small time step

			// Simulate time slightly less than the attack delay
			for t := 0.0; t < attackDelay-dt/2; t += dt {
				autoAttackSystem.Update(dt)
				Expect(targetHealth.GetCurrentHP()).To(Equal(initialTargetHP), "Target HP should not change before the attack lands")
			}

			// Simulate the time step where the attack should land
			autoAttackSystem.Update(dt)
			expectedHP := initialTargetHP - preCalculatedExpectedDamage // Assumes no damage reduction for simplicity
			Expect(targetHealth.GetCurrentHP()).To(BeNumerically("~", expectedHP, 0.01), "Target HP should decrease after the attack lands")
		})

		It("should reset the attack timer after attacking", func() {
			Expect(playerAttack).NotTo(BeNil())
			Expect(targetHealth).NotTo(BeNil())
			Expect(playerAttack.GetFinalAttackSpeed()).To(BeNumerically(">", 0))
			attackDelay := 1.0 / playerAttack.GetFinalAttackSpeed()
			dt := attackDelay + 0.01 // Simulate enough time for one attack

			// First attack
			autoAttackSystem.Update(dt)
			hpAfterFirstAttack := targetHealth.GetCurrentHP()

			// Simulate time slightly less than the *next* attack delay
			smallDt := 0.1 // Use smaller steps for inner loop
			for t := 0.0; t < attackDelay-smallDt/2; t += smallDt {
				autoAttackSystem.Update(smallDt)
				Expect(targetHealth.GetCurrentHP()).To(BeNumerically("~", hpAfterFirstAttack, 0.01), "Target HP should not change again before the second attack lands")
			}

			// Simulate the time step where the second attack should land
			autoAttackSystem.Update(smallDt)
			expectedHP := hpAfterFirstAttack - preCalculatedExpectedDamage
			Expect(targetHealth.GetCurrentHP()).To(BeNumerically("~", expectedHP, 0.01), "Target HP should decrease after the second attack lands")
		})

		It("should generate mana for the attacker on attack", func() {
			Expect(playerAttack).NotTo(BeNil())
			playerMana, ok := world.GetMana(player)
			Expect(ok).To(BeTrue())
			initialMana := playerMana.GetCurrentMana()
			// Expect(initialMana).To(BeNumerically("~", 0.0)) // Initial mana might not be 0 depending on factory

			Expect(playerAttack.GetFinalAttackSpeed()).To(BeNumerically(">", 0))
			attackDelay := 1.0 / playerAttack.GetFinalAttackSpeed()
			dt := attackDelay + 0.01 // Simulate enough time for one attack

			autoAttackSystem.Update(dt) // Perform the attack

			// Mana gain per attack (TFT standard is 10, but check components.Attack if different)
			expectedManaGain := 10.0
			Expect(playerMana.GetCurrentMana()).To(BeNumerically("~", initialMana+expectedManaGain), "Player should gain mana after attacking")
		})

		It("should stop attacking if the target dies", func() {
			Expect(playerAttack).NotTo(BeNil())
			Expect(targetHealth).NotTo(BeNil())
			// Reduce target HP so it dies in one hit
			Expect(playerAttack.GetFinalAD()).To(BeNumerically(">", 0)) // Ensure AD is positive
			targetHealth.SetCurrentHealth(preCalculatedExpectedDamage - 1)
			// hpBeforeAttack := targetHealth.GetCurrentHP() // Not strictly needed

			Expect(playerAttack.GetFinalAttackSpeed()).To(BeNumerically(">", 0))
			attackDelay := 1.0 / playerAttack.GetFinalAttackSpeed()
			dt := attackDelay + 0.01 // Simulate enough time for one attack

			// First attack (kills target)
			autoAttackSystem.Update(dt)
			Expect(targetHealth.GetCurrentHP()).To(BeNumerically("<=", 0), "Target should be dead or below 0 HP")
			hpAfterDeath := targetHealth.GetCurrentHP() // Store HP after the killing blow

			// Simulate more time, enough for several more attacks if target were alive
			for i := 0; i < 5; i++ {
				autoAttackSystem.Update(attackDelay)
			}

			// Target HP should not change further after death
			Expect(targetHealth.GetCurrentHP()).To(BeNumerically("~", hpAfterDeath, 0.01), "Target HP should not change after it dies")
		})
	})

	Context("when no valid target is in range", func() {
		BeforeEach(func() {
			Expect(playerAttack).NotTo(BeNil())
			// Move target out of range
			targetPos, ok := world.GetPosition(target)
			Expect(ok).To(BeTrue())
			Expect(targetPos).NotTo(BeNil())
			targetPos.X = playerAttack.GetFinalRange() + 1 // Just outside range
		})

		It("should not perform an attack", func() {
			Expect(targetHealth).NotTo(BeNil())
			initialTargetHP := targetHealth.GetCurrentHP()
			dt := 0.1
			simulationTime := 3.0 // Simulate for 3 seconds

			// Simulate time
			for t := 0.0; t < simulationTime; t += dt {
				autoAttackSystem.Update(dt)
			}

			Expect(targetHealth.GetCurrentHP()).To(Equal(initialTargetHP), "Target HP should not change when out of range")
		})

		It("should not generate mana for the attacker", func() {
			playerMana, ok := world.GetMana(player)
			Expect(ok).To(BeTrue())
			initialMana := playerMana.GetCurrentMana()

			dt := 0.1
			simulationTime := 3.0 // Simulate for 3 seconds

			// Simulate time
			for t := 0.0; t < simulationTime; t += dt {
				autoAttackSystem.Update(dt)
			}

			Expect(playerMana.GetCurrentMana()).To(BeNumerically("~", initialMana), "Player mana should not change when no target is attacked")
		})
	})

	Context("when the attacker has 0 Attack Speed", func() {
		BeforeEach(func() {
			Expect(playerAttack).NotTo(BeNil())
			playerAttack.SetBaseAttackSpeed(0)
			playerAttack.SetFinalAttackSpeed(0) // Ensure final AS is also 0
		})

		It("should not perform an attack", func() {
			Expect(targetHealth).NotTo(BeNil())
			initialTargetHP := targetHealth.GetCurrentHP()
			dt := 0.1

			// Simulate a significant amount of time
			for t := 0.0; t < 5.0; t += dt {
				autoAttackSystem.Update(dt)
			}

			Expect(targetHealth.GetCurrentHP()).To(Equal(initialTargetHP), "Target HP should not change when attacker AS is 0")
		})
	})

})
