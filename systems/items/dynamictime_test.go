package itemsys_test // Use _test package convention

import (
	"github.com/suriz/tft-dps-simulator/components"
	"github.com/suriz/tft-dps-simulator/data"
	"github.com/suriz/tft-dps-simulator/ecs"
	"github.com/suriz/tft-dps-simulator/factory"
	"github.com/suriz/tft-dps-simulator/managers"
	itemsys "github.com/suriz/tft-dps-simulator/systems/items" // Alias the package

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// Helper function to simulate time and reset bonus stats
func simulateTime(system *itemsys.DynamicTimeItemSystem, world *ecs.World, entity ecs.Entity, duration float64, dt float64) {
	spellComp, spellOk := world.GetSpell(entity)
	attackComp, attackOk := world.GetAttack(entity)

	for t := 0.0; t < duration; t += dt {
		// Reset bonus stats before each system update (mimics main loop)
		if spellOk {
			spellComp.SetBonusAP(0) // Assuming SetBonusAP exists
		}
		if attackOk {
			attackComp.SetBonusPercentAttackSpeed(0) // Assuming SetBonusAttackSpeed exists
		}

		// Update the system
		system.Update(dt)

		// Optional: Stat Calculation System would run here in a real sim
		// For testing this system in isolation, we check the bonus stats directly
		// or the effect component's state.
	}
}

var _ = Describe("DynamicTimeItemSystem", func() {
	var (
		world             *ecs.World
		championFactory   *factory.ChampionFactory
		equipmentManager  *managers.EquipmentManager
		dynamicTimeSystem *itemsys.DynamicTimeItemSystem
		entity            ecs.Entity
		attackComp        *components.Attack
		spellComp         *components.Spell
		archangelsData    *data.Item // Store item data for assertions
		quicksilverData   *data.Item // Store item data for assertions
		ok                bool
		err               error
		deltaTime         float64 = 0.1 // Simulation time step
	)

	BeforeEach(func() {
		// Get item data for assertions
		archangelsData = data.GetItemByApiName(data.TFT_Item_ArchangelsStaff)
		Expect(archangelsData).NotTo(BeNil())
		quicksilverData = data.GetItemByApiName(data.TFT_Item_Quicksilver)
		Expect(quicksilverData).NotTo(BeNil())

		// Setup ECS and Systems
		world = ecs.NewWorld()
		championFactory = factory.NewChampionFactory(world)
		equipmentManager = managers.NewEquipmentManager(world)
		dynamicTimeSystem = itemsys.NewDynamicTimeItemSystem(world)

		// Create a test entity (e.g., a basic champion)
		// Use a champion known to exist in your data
		entity, err = championFactory.CreatePlayerChampion("TFT_TrainingDummy", 1) // Use a simple entity
		Expect(err).NotTo(HaveOccurred())

		// Ensure necessary components exist (factory should add them)
		attackComp, ok = world.GetAttack(entity)
		Expect(ok).To(BeTrue(), "Entity should have Attack component")
		spellComp, ok = world.GetSpell(entity)
		Expect(ok).To(BeTrue(), "Entity should have Spell component")
		_, ok = world.GetEquipment(entity)
		Expect(ok).To(BeTrue(), "Entity should have Equipment component")

		// Initialize base stats if needed (factory might do this)
		spellComp.SetBaseAP(50)            // Example base AP
		attackComp.SetBaseAttackSpeed(0.7) // Example base AS
	})

	Context("Archangel's Staff", func() {
		var (
			interval      float64
			apPerInterval float64
		)

		BeforeEach(func() {
			// Get expected values from loaded data
			interval = archangelsData.Effects["IntervalSeconds"]
			apPerInterval = archangelsData.Effects["APPerInterval"]
			Expect(interval).To(BeNumerically(">", 0))
			Expect(apPerInterval).To(BeNumerically(">", 0))
		})

		It("should not add effect or AP if item is not equipped", func() {
			simulateTime(dynamicTimeSystem, world, entity, 10.0, deltaTime)

			_, exists := world.GetArchangelsEffect(entity)
			Expect(exists).To(BeFalse())
			Expect(spellComp.GetBonusAP()).To(BeNumerically("~", 0.0)) // Check GetBonusAP
		})

		It("should initialize effect with correct values from data", func() {
			err := equipmentManager.AddItemToChampion(entity, "TFT_Item_ArchangelsStaff")
			Expect(err).NotTo(HaveOccurred())

			effect, exists := world.GetArchangelsEffect(entity)
			Expect(exists).To(BeTrue())
			Expect(effect.GetInterval()).To(BeNumerically("~", interval))
			Expect(effect.GetAPPerInterval()).To(BeNumerically("~", apPerInterval)) // Check GetapPerInterval
		})

		It("should stack AP correctly over time", func() {
			err := equipmentManager.AddItemToChampion(entity, "TFT_Item_ArchangelsStaff")
			Expect(err).NotTo(HaveOccurred())

			// Time just before first stack
			simulateTime(dynamicTimeSystem, world, entity, interval-(deltaTime/2), deltaTime)
			effect, _ := world.GetArchangelsEffect(entity)
			Expect(effect.GetStacks()).To(Equal(0))
			// Bonus AP should be 0 *before* the system runs for the interval tick
			Expect(spellComp.GetBonusAP()).To(BeNumerically("~", 0.0))

			// Time just after first stack
			simulateTime(dynamicTimeSystem, world, entity, deltaTime, deltaTime) // Simulate one more step
			effect, _ = world.GetArchangelsEffect(entity)
			Expect(effect.GetStacks()).To(Equal(1))
			// After the system runs for the interval tick, BonusAP should reflect the new stack
			Expect(spellComp.GetBonusAP()).To(BeNumerically("~", 1*apPerInterval, 0.01))
			Expect(effect.GetTimer()).To(BeNumerically("~", deltaTime, 0.001)) // Timer resets with overflow

			// Time just after second stack
			simulateTime(dynamicTimeSystem, world, entity, interval-effect.GetTimer(), deltaTime) // Simulate remaining time until next interval
			simulateTime(dynamicTimeSystem, world, entity, deltaTime, deltaTime)                  // Simulate the step crossing the interval
			effect, _ = world.GetArchangelsEffect(entity)
			Expect(effect.GetStacks()).To(Equal(2))
			Expect(spellComp.GetBonusAP()).To(BeNumerically("~", 2*apPerInterval, 0.01))

			// Longer duration (e.g., 23 seconds -> 4 stacks)
			simulateTime(dynamicTimeSystem, world, entity, (23.0 - (2*interval + deltaTime)), deltaTime) // Simulate remaining time
			effect, _ = world.GetArchangelsEffect(entity)
			Expect(effect.GetStacks()).To(Equal(4))
			Expect(spellComp.GetBonusAP()).To(BeNumerically("~", 4*apPerInterval, 0.01))
		})

		It("should stack AP correctly with multiple instances (e.g., 3)", func() {
			numStaffs := 3
			for i := 0; i < numStaffs; i++ {
				err := equipmentManager.AddItemToChampion(entity, "TFT_Item_ArchangelsStaff")
				Expect(err).NotTo(HaveOccurred())
			}

			// Time just after second stack (should have 2 stacks internally)
			simulateTime(dynamicTimeSystem, world, entity, 2*interval+deltaTime, deltaTime)
			effect, _ := world.GetArchangelsEffect(entity) // Re-fetch in case of changes (though unlikely here)
			Expect(effect.GetStacks()).To(Equal(2), "Internal stack count should reflect time passed, not item count")

			// Bonus AP should be stacks * AP_per_stack * number_of_staffs
			expectedBonusAP := float64(2) * apPerInterval * float64(numStaffs)
			Expect(spellComp.GetBonusAP()).To(BeNumerically("~", expectedBonusAP, 0.01), "Bonus AP should be scaled by the number of staffs")

			// Simulate more time (e.g., 23 seconds -> 4 stacks internally)
			// Total time simulated so far: 2*interval + deltaTime
			// Additional time needed: 23.0 - (2*interval + deltaTime)
			simulateTime(dynamicTimeSystem, world, entity, (23.0 - (2*interval + deltaTime)), deltaTime)
			effect, _ = world.GetArchangelsEffect(entity)
			Expect(effect.GetStacks()).To(Equal(4), "Internal stack count should be 4 after 23 seconds")

			expectedBonusAP = float64(4) * apPerInterval * float64(numStaffs)
			Expect(spellComp.GetBonusAP()).To(BeNumerically("~", expectedBonusAP, 0.01), "Bonus AP should be scaled by 3 staffs after 4 stacks")
		})

		It("should stop stacking AP after item removal", func() {
			err := equipmentManager.AddItemToChampion(entity, "TFT_Item_ArchangelsStaff")
			Expect(err).NotTo(HaveOccurred())

			// Gain 2 stacks
			simulateTime(dynamicTimeSystem, world, entity, 2*interval+deltaTime, deltaTime)
			effect, _ := world.GetArchangelsEffect(entity)
			Expect(effect.GetStacks()).To(Equal(2))
			Expect(spellComp.GetBonusAP()).To(BeNumerically("~", 2*apPerInterval, 0.01))
			apBeforeRemoval := spellComp.GetBonusAP()
			Expect(apBeforeRemoval).To(BeNumerically("~", 2*apPerInterval, 0.01))

			// Remove item
			err = equipmentManager.RemoveItemFromChampion(entity, "TFT_Item_ArchangelsStaff")
			Expect(err).NotTo(HaveOccurred())

			// Check component removed by manager
			_, exists := world.GetArchangelsEffect(entity)
			Expect(exists).To(BeFalse(), "ArchangelsEffect component should be removed by EquipmentManager")

			// Simulate more time
			simulateTime(dynamicTimeSystem, world, entity, 10.0, deltaTime)

			// Bonus AP should not increase further (it gets reset by simulateTime helper)
			Expect(spellComp.GetBonusAP()).To(BeNumerically("~", 0.0), "Bonus AP should be 0 after reset and no item effect")

			// Double check component is still gone
			_, exists = world.GetArchangelsEffect(entity)
			Expect(exists).To(BeFalse())
		})
	})

	Context("Quicksilver", func() {
		var (
			duration     float64
			procInterval float64
			procAS       float64
		)

		BeforeEach(func() {
			duration = quicksilverData.Effects["SpellShieldDuration"]
			procInterval = quicksilverData.Effects["ProcInterval"]
			procAS = quicksilverData.Effects["ProcAttackSpeed"] // This is the decimal value
			Expect(duration).To(BeNumerically(">", 0))
			Expect(procInterval).To(BeNumerically(">", 0))
			Expect(procAS).To(BeNumerically(">", 0))
		})

		It("should not add effect or AS if item is not equipped", func() {
			simulateTime(dynamicTimeSystem, world, entity, duration+5.0, deltaTime) // Simulate past duration

			_, exists := world.GetQuicksilverEffect(entity)
			Expect(exists).To(BeFalse())
			Expect(attackComp.GetBonusPercentAttackSpeed()).To(BeNumerically("~", 0.0)) // Check GetBonusPercentAttackSpeed
		})

		It("should initialize effect with correct values from data", func() {
			err := equipmentManager.AddItemToChampion(entity, "TFT_Item_Quicksilver")
			Expect(err).NotTo(HaveOccurred())

			effect, exists := world.GetQuicksilverEffect(entity)
			Expect(exists).To(BeTrue())
			Expect(effect.GetRemainingDuration()).To(BeNumerically("~", duration))
			Expect(effect.GetProcInterval()).To(BeNumerically("~", procInterval))
			Expect(effect.GetProcAttackSpeed()).To(BeNumerically("~", procAS))
			Expect(effect.IsActive()).To(BeTrue())
		})

		It("should track immunity duration correctly and deactivate", func() {
			err := equipmentManager.AddItemToChampion(entity, "TFT_Item_Quicksilver")
			Expect(err).NotTo(HaveOccurred())

			// Time just before expiry
			simulateTime(dynamicTimeSystem, world, entity, duration-(deltaTime/2), deltaTime)
			effect, exists := world.GetQuicksilverEffect(entity)
			Expect(exists).To(BeTrue())
			Expect(effect.IsActive()).To(BeTrue())
			Expect(effect.GetRemainingDuration()).To(BeNumerically(">", 0))

			// Time just after expiry
			simulateTime(dynamicTimeSystem, world, entity, deltaTime, deltaTime) // Simulate one more step
			quicksilverEffect, exists := world.GetQuicksilverEffect(entity)
			// The system removes the component when duration hits 0
			Expect(exists).To(BeTrue(), "QuicksilverEffect component still exists after duration expires")
			Expect(quicksilverEffect.IsActive()).To(BeFalse(), "QuicksilverEffect should be inactive after duration expires")
		})

		It("should stack Attack Speed correctly during first 18s in combat", func() {
			err := equipmentManager.AddItemToChampion(entity, "TFT_Item_Quicksilver")
			Expect(err).NotTo(HaveOccurred())

			// Time just before first proc
			simulateTime(dynamicTimeSystem, world, entity, procInterval-deltaTime, deltaTime)
			effect, _ := world.GetQuicksilverEffect(entity)
			Expect(effect.GetStacks()).To(Equal(0))
			Expect(effect.GetCurrentBonusAS()).To(BeNumerically("~", 0.0))
			Expect(attackComp.GetBonusPercentAttackSpeed()).To(BeNumerically("~", 0.0))

            // Time just after first proc (simulate one more step to cross the interval)
            simulateTime(dynamicTimeSystem, world, entity, deltaTime, deltaTime)
            effect, _ = world.GetQuicksilverEffect(entity)
            // Check internal effect state first
            Expect(effect.GetStacks()).To(Equal(1)) // Check internal stack count
            Expect(effect.GetCurrentBonusAS()).To(BeNumerically("~", 1*procAS, 0.001))
            // Then check the applied stat on the champion component
            Expect(attackComp.GetBonusPercentAttackSpeed()).To(BeNumerically("~", 1*procAS, 0.001))
            Expect(effect.GetProcTimer()).To(BeNumerically("~", 0.0, 0.001)) // Timer resets exactly to 0

			// Time just after second proc
            // Simulate remaining time until the next interval (procInterval - current timer)
            // Current timer is ~0.0, so simulate procInterval
            simulateTime(dynamicTimeSystem, world, entity, procInterval-effect.GetProcTimer()-deltaTime, deltaTime) // Simulate up to just before next proc
            simulateTime(dynamicTimeSystem, world, entity, deltaTime, deltaTime)                                   // Simulate step crossing interval
            effect, _ = world.GetQuicksilverEffect(entity)
            Expect(effect.GetStacks()).To(Equal(2)) // Check internal stack count
            Expect(effect.GetCurrentBonusAS()).To(BeNumerically("~", 2*procAS, 0.001))
            Expect(attackComp.GetBonusPercentAttackSpeed()).To(BeNumerically("~", 2*procAS, 0.001))

			// Check near end of duration (e.g., 17.9s for 18s duration, 2s interval -> 8 procs)
            // Total time simulated so far: 4.0s (see detailed calculation above)
            // Need to simulate remaining time: 17.9 - 4.0 = 13.9s
            simulateTime(dynamicTimeSystem, world, entity, 13.9, deltaTime) // Simulate remaining time precisely
            effect, _ = world.GetQuicksilverEffect(entity)
            // Total time = 4.0 + 13.9 = 17.9s. Procs at 2, 4, 6, 8, 10, 12, 14, 16. -> 8 procs.
            expectedStacks := 8
            Expect(effect.GetStacks()).To(Equal(expectedStacks)) // Check line 302
            Expect(effect.GetCurrentBonusAS()).To(BeNumerically("~", float64(expectedStacks)*procAS, 0.001))
            Expect(attackComp.GetBonusPercentAttackSpeed()).To(BeNumerically("~", float64(expectedStacks)*procAS, 0.001))
            // Timer should be 1.9s (17.9s total time - 16.0s last proc time)
            Expect(effect.GetProcTimer()).To(BeNumerically("~", 1.9, 0.001))
		})

		It("should stop stacking Attack Speed after immunity expires", func() {
			err := equipmentManager.AddItemToChampion(entity, "TFT_Item_Quicksilver")
			Expect(err).NotTo(HaveOccurred())

			// Simulate just past duration
			simulateTime(dynamicTimeSystem, world, entity, duration+deltaTime, deltaTime)
			_, exists := world.GetQuicksilverEffect(entity)
			Expect(exists).To(BeTrue()) // Component should be gone

			// Record bonus AS (which should be 0 now due to reset in simulateTime)
			bonusASAfterExpiry := attackComp.GetBonusPercentAttackSpeed()
			Expect(bonusASAfterExpiry).To(BeNumerically("~", 0.0))

			// Simulate more time
			simulateTime(dynamicTimeSystem, world, entity, 5.0, deltaTime)

			// Bonus AS should remain 0 as the component is gone
			Expect(attackComp.GetBonusPercentAttackSpeed()).To(BeNumerically("~", 0.0))
		})

		It("should stop stacking Attack Speed after item removal", func() {
			err := equipmentManager.AddItemToChampion(entity, "TFT_Item_Quicksilver")
			Expect(err).NotTo(HaveOccurred())

			// Gain some AS stacks (e.g., 5 seconds -> 2 stacks)
			simulateTime(dynamicTimeSystem, world, entity, 5.0, deltaTime)
			effect, _ := world.GetQuicksilverEffect(entity)
			Expect(effect.GetStacks()).To(Equal(2)) // Check internal stack count
			Expect(effect.GetCurrentBonusAS()).To(BeNumerically("~", 2*procAS, 0.001))
			Expect(attackComp.GetBonusPercentAttackSpeed()).To(BeNumerically("~", 2*procAS, 0.001))

			// Remove item
			err = equipmentManager.RemoveItemFromChampion(entity, "TFT_Item_Quicksilver")
			Expect(err).NotTo(HaveOccurred())

			// Check component removed by manager
			_, exists := world.GetQuicksilverEffect(entity)
			Expect(exists).To(BeFalse(), "QuicksilverEffect component should be removed by EquipmentManager")

			// Simulate more time
			simulateTime(dynamicTimeSystem, world, entity, 10.0, deltaTime)

			// Bonus AS should be 0 after reset and no item effect
			Expect(attackComp.GetBonusPercentAttackSpeed()).To(BeNumerically("~", 0.0))

			// Double check component is still gone
			_, exists = world.GetQuicksilverEffect(entity)
			Expect(exists).To(BeFalse())
		})
	})
})
