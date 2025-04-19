package itemsys_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/suriz/tft-dps-simulator/components"
	"github.com/suriz/tft-dps-simulator/components/effects"
	"github.com/suriz/tft-dps-simulator/data"
	"github.com/suriz/tft-dps-simulator/ecs"
	"github.com/suriz/tft-dps-simulator/factory"
	"github.com/suriz/tft-dps-simulator/managers"
	eventsys "github.com/suriz/tft-dps-simulator/systems/events"
	itemsys "github.com/suriz/tft-dps-simulator/systems/items"
	"github.com/suriz/tft-dps-simulator/utils"
)

// only dymanic bonuses are tested here, static bonuses are tested in BaseStaticItemSystem tests
var _ = Describe("DynamicEventItemSystem", func() {
	var (
		world              *ecs.World
		mockEventBus       *utils.MockEventBus
		championFactory    *factory.ChampionFactory
		equipmentManager   *managers.EquipmentManager
		system             *itemsys.DynamicEventItemSystem
		player             ecs.Entity
		target             ecs.Entity
		playerAttack       *components.Attack
		playerSpell        *components.Spell
		playerHealth       *components.Health
		playerTitansEffect *effects.TitansResolveEffect
		titansData         *data.Item
		err                error

		// Titan's config values
		titansMaxStacks         int
		titansADPerStack        float64
		titansAPPerStack        float64
		titansBonusResistsAtCap float64
	)

	BeforeEach(func() {
		world = ecs.NewWorld()
		mockEventBus = utils.NewMockEventBus()
		championFactory = factory.NewChampionFactory(world)
		equipmentManager = managers.NewEquipmentManager(world)

		// Create and register the system
		system = itemsys.NewDynamicEventItemSystem(world, mockEventBus)
		mockEventBus.RegisterHandler(system) // Register the system to handle events

		// Create Player (Blue Golem)
		player, err = championFactory.CreatePlayerChampion("TFT_BlueGolem", 1)
		Expect(err).NotTo(HaveOccurred())

		// Create Target (Training Dummy)
		target, err = championFactory.CreateEnemyChampion("TFT_TrainingDummy", 1)
		Expect(err).NotTo(HaveOccurred())

		// Get Titan's data
		titansData = data.GetItemByApiName(data.TFT_Item_TitansResolve)
		Expect(titansData).NotTo(BeNil())
		titansMaxStacks = int(titansData.Effects["StackCap"])
		titansADPerStack = titansData.Effects["StackingAD"]
		titansAPPerStack = titansData.Effects["StackingSP"]
		titansBonusResistsAtCap = titansData.Effects["BonusResistsAtStackCap"]

		// Add Titan's Resolve to player
		err = equipmentManager.AddItemToChampion(player, data.TFT_Item_TitansResolve)
		Expect(err).NotTo(HaveOccurred())

		// Get player components
		var ok bool
		playerAttack, ok = world.GetAttack(player)
		Expect(ok).To(BeTrue())
		playerSpell, ok = world.GetSpell(player)
		Expect(ok).To(BeTrue())
		playerHealth, ok = world.GetHealth(player)
		Expect(ok).To(BeTrue())
		playerTitansEffect, ok = world.GetTitansResolveEffect(player)
		Expect(ok).To(BeTrue())
		Expect(playerTitansEffect).NotTo(BeNil())

		// Verify initial state after adding item
		Expect(playerTitansEffect.GetCurrentStacks()).To(Equal(0))
		Expect(playerTitansEffect.IsMaxStacks).To(BeFalse())

		// Check dynamic bonuses are initially zero
		Expect(playerAttack.GetBonusPercentAD()).To(BeNumerically("~", 0, 0.001))
		Expect(playerSpell.GetBonusAP()).To(BeNumerically("~", 0, 0.001))
	})

	Context("Stacking Mechanics", func() {
		It("should gain 1 stack, AD, and AP on AttackLandedEvent", func() {
			initialBonusAD := playerAttack.GetBonusPercentAD()
			initialBonusAP := playerSpell.GetBonusAP()

			mockEventBus.SimulateAndProcessEvent(eventsys.AttackLandedEvent{Source: player, Target: target})

			Expect(playerTitansEffect.GetCurrentStacks()).To(Equal(1))
			Expect(playerAttack.GetBonusPercentAD()).To(BeNumerically("~", initialBonusAD+titansADPerStack, 0.001))
			Expect(playerSpell.GetBonusAP()).To(BeNumerically("~", initialBonusAP+titansAPPerStack, 0.001))
			Expect(playerHealth.GetBonusArmor()).To(BeNumerically("~", 0.0, 0.001)) // No max stack bonus yet
			Expect(playerHealth.GetBonusMR()).To(BeNumerically("~", 0, 0.001)) // No max stack bonus yet
		})

		It("should gain 1 stack, AD, and AP on DamageAppliedEvent", func() {
			initialBonusAD := playerAttack.GetBonusPercentAD()
			initialBonusAP := playerSpell.GetBonusAP()

			mockEventBus.SimulateAndProcessEvent(eventsys.DamageAppliedEvent{Source: target, Target: player, FinalTotalDamage: 10})

			Expect(playerTitansEffect.GetCurrentStacks()).To(Equal(1))
			Expect(playerAttack.GetBonusPercentAD()).To(BeNumerically("~", initialBonusAD+titansADPerStack, 0.001))
			Expect(playerSpell.GetBonusAP()).To(BeNumerically("~", initialBonusAP+titansAPPerStack, 0.001))
			Expect(playerHealth.GetBonusArmor()).To(BeNumerically("~", 0.0, 0.001))
			Expect(playerHealth.GetBonusMR()).To(BeNumerically("~", 0, 0.001))
		})

		It("should stack bonuses additively from multiple events", func() {
			initialBonusAD := playerAttack.GetBonusPercentAD()
			initialBonusAP := playerSpell.GetBonusAP()

			mockEventBus.SimulateAndProcessEvent(eventsys.AttackLandedEvent{Source: player, Target: target})
			mockEventBus.SimulateAndProcessEvent(eventsys.DamageAppliedEvent{Source: target, Target: player, FinalTotalDamage: 10})

			Expect(playerTitansEffect.GetCurrentStacks()).To(Equal(2))
			Expect(playerAttack.GetBonusPercentAD()).To(BeNumerically("~", initialBonusAD+titansADPerStack*2, 0.001))
			Expect(playerSpell.GetBonusAP()).To(BeNumerically("~", initialBonusAP+titansAPPerStack*2, 0.001))
		})
	})

	Context("Max Stacks", func() {
		It("should apply bonus resists only when max stacks are reached", func() {
			// Simulate events until just before max stacks
			for i := 0; i < titansMaxStacks-1; i++ {
				mockEventBus.SimulateAndProcessEvent(eventsys.AttackLandedEvent{Source: player, Target: target})
			}
			Expect(playerTitansEffect.GetCurrentStacks()).To(Equal(titansMaxStacks - 1))
			Expect(playerTitansEffect.IsMaxStacks).To(BeFalse())
			Expect(playerHealth.GetBonusArmor()).To(BeNumerically("~", 0.0, 0.001)) // Only static armor
			Expect(playerHealth.GetBonusMR()).To(BeNumerically("~", 0, 0.001)) // No bonus MR yet

			// Simulate the final event to reach max stacks
			mockEventBus.SimulateAndProcessEvent(eventsys.DamageAppliedEvent{Source: target, Target: player, FinalTotalDamage: 10})

			Expect(playerTitansEffect.GetCurrentStacks()).To(Equal(titansMaxStacks))
			Expect(playerTitansEffect.IsMaxStacks).To(BeTrue())
			// Check final stats include static + max stack bonus
			Expect(playerAttack.GetBonusPercentAD()).To(BeNumerically("~", titansADPerStack*float64(titansMaxStacks), 0.001))
			Expect(playerSpell.GetBonusAP()).To(BeNumerically("~", titansAPPerStack*float64(titansMaxStacks), 0.001))
			Expect(playerHealth.GetBonusArmor()).To(BeNumerically("~", titansBonusResistsAtCap, 0.001))
			Expect(playerHealth.GetBonusMR()).To(BeNumerically("~", titansBonusResistsAtCap, 0.001))
		})

		It("should not stack beyond the maximum limit", func() {
			// Reach max stacks
			for i := 0; i < titansMaxStacks; i++ {
				mockEventBus.SimulateAndProcessEvent(eventsys.AttackLandedEvent{Source: player, Target: target})
			}
			Expect(playerTitansEffect.GetCurrentStacks()).To(Equal(titansMaxStacks))

			// Record stats at max stacks
			maxStackBonusAD := playerAttack.GetBonusPercentAD()
			maxStackBonusAP := playerSpell.GetBonusAP()
			maxStackBonusArmor := playerHealth.GetBonusArmor()
			maxStackBonusMR := playerHealth.GetBonusMR()

			// Simulate another event
			mockEventBus.SimulateAndProcessEvent(eventsys.DamageAppliedEvent{Source: target, Target: player, FinalTotalDamage: 10})

			// Assert stats haven't changed
			Expect(playerTitansEffect.GetCurrentStacks()).To(Equal(titansMaxStacks))
			Expect(playerAttack.GetBonusPercentAD()).To(BeNumerically("~", maxStackBonusAD, 0.001))
			Expect(playerSpell.GetBonusAP()).To(BeNumerically("~", maxStackBonusAP, 0.001))
			Expect(playerHealth.GetBonusArmor()).To(BeNumerically("~", maxStackBonusArmor, 0.001))
			Expect(playerHealth.GetBonusMR()).To(BeNumerically("~", maxStackBonusMR, 0.001))
		})
	})

	Context("Item Absence / Removal", func() {
		var otherPlayer ecs.Entity
		var otherPlayerAttack *components.Attack
		var otherPlayerSpell *components.Spell
		var otherPlayerHealth *components.Health

		BeforeEach(func() {
			// Create another player without Titan's
			otherPlayer, err = championFactory.CreatePlayerChampion("TFT_BlueGolem", 1)
			Expect(err).NotTo(HaveOccurred())
			// Ensure it has stats components
			var ok bool
			otherPlayerAttack, ok = world.GetAttack(otherPlayer)
			Expect(ok).To(BeTrue())
			otherPlayerSpell, ok = world.GetSpell(otherPlayer)
			Expect(ok).To(BeTrue())
			otherPlayerHealth, ok = world.GetHealth(otherPlayer)
			Expect(ok).To(BeTrue())
		})

		It("should have no effect if the item is not equipped", func() {
			initialAD := otherPlayerAttack.GetBonusPercentAD()
			initialAP := otherPlayerSpell.GetBonusAP()
			initialArmor := otherPlayerHealth.GetBonusArmor()
			initialMR := otherPlayerHealth.GetBonusMR()

			// Simulate events involving the champion without the item
			mockEventBus.SimulateAndProcessEvent(eventsys.AttackLandedEvent{Source: otherPlayer, Target: target})
			mockEventBus.SimulateAndProcessEvent(eventsys.DamageAppliedEvent{Source: target, Target: otherPlayer, FinalTotalDamage: 10})

			// Check component doesn't exist and stats are unchanged
			_, ok := world.GetTitansResolveEffect(otherPlayer)
			Expect(ok).To(BeFalse())
			Expect(otherPlayerAttack.GetBonusPercentAD()).To(BeNumerically("~", initialAD, 0.001))
			Expect(otherPlayerSpell.GetBonusAP()).To(BeNumerically("~", initialAP, 0.001))
			Expect(otherPlayerHealth.GetBonusArmor()).To(BeNumerically("~", initialArmor, 0.001))
			Expect(otherPlayerHealth.GetBonusMR()).To(BeNumerically("~", initialMR, 0.001))
		})

		It("should remove all dynamic and static bonuses when the item is removed (before max stacks)", func() {
			// Gain some stacks (e.g., 10)
			stacksToGain := 10
			if titansMaxStacks < stacksToGain {
				stacksToGain = titansMaxStacks / 2 // Ensure we don't accidentally hit max
			}
			for i := 0; i < stacksToGain; i++ {
				mockEventBus.SimulateAndProcessEvent(eventsys.AttackLandedEvent{Source: player, Target: target})
			}
			Expect(playerTitansEffect.GetCurrentStacks()).To(Equal(stacksToGain))

			// Record stats BEFORE removal
			adBefore := playerAttack.GetBonusPercentAD()
			apBefore := playerSpell.GetBonusAP()
			armorBefore := playerHealth.GetBonusArmor()
			mrBefore := playerHealth.GetBonusMR()

			// Calculate expected dynamic bonuses that were applied
			expectedDynamicAD := titansADPerStack * float64(stacksToGain)
			expectedDynamicAP := titansAPPerStack * float64(stacksToGain)

			// Remove the item
			err = equipmentManager.RemoveItemFromChampion(player, data.TFT_Item_TitansResolve)
			Expect(err).NotTo(HaveOccurred())

			// Assert component is gone
			_, ok := world.GetTitansResolveEffect(player)
			Expect(ok).To(BeFalse())

			// Assert stats AFTER removal (should be original stats minus all bonuses from this item)
			Expect(playerAttack.GetBonusPercentAD()).To(BeNumerically("~", adBefore-expectedDynamicAD, 0.001), "Dynamic AD bonus should be removed")
			Expect(playerSpell.GetBonusAP()).To(BeNumerically("~", apBefore-expectedDynamicAP, 0.001), "Dynamic AP bonus should be removed")
			Expect(playerHealth.GetBonusArmor()).To(BeNumerically("~", armorBefore, 0.001), "Armor should be unchanged as no dynamic Armor was added")
			Expect(playerHealth.GetBonusMR()).To(BeNumerically("~", mrBefore, 0.001), "MR should be unchanged as no dynamic MR was added") // MR was 0 before max stacks
		})

		It("should remove all dynamic and static bonuses when the item is removed (at max stacks)", func() {
			// Reach max stacks
			for i := 0; i < titansMaxStacks; i++ {
				mockEventBus.SimulateAndProcessEvent(eventsys.AttackLandedEvent{Source: player, Target: target})
			}
			Expect(playerTitansEffect.GetCurrentStacks()).To(Equal(titansMaxStacks))
			Expect(playerTitansEffect.IsMaxStacks).To(BeTrue())

			// Record stats BEFORE removal
			adBefore := playerAttack.GetBonusPercentAD()
			apBefore := playerSpell.GetBonusAP()
			armorBefore := playerHealth.GetBonusArmor()
			mrBefore := playerHealth.GetBonusMR()

			// Calculate expected dynamic bonuses that were applied
			expectedDynamicAD := titansADPerStack * float64(titansMaxStacks)
			expectedDynamicAP := titansAPPerStack * float64(titansMaxStacks)
			expectedDynamicResists := titansBonusResistsAtCap // This was added at max stacks

			// Remove the item
			err = equipmentManager.RemoveItemFromChampion(player, data.TFT_Item_TitansResolve)
			Expect(err).NotTo(HaveOccurred())

			// Assert component is gone
			_, ok := world.GetTitansResolveEffect(player)
			Expect(ok).To(BeFalse())

			// Assert stats AFTER removal (should be original stats minus all bonuses from this item)
			Expect(playerAttack.GetBonusPercentAD()).To(BeNumerically("~", adBefore-expectedDynamicAD, 0.001), "Dynamic AD bonus should be removed")
			Expect(playerSpell.GetBonusAP()).To(BeNumerically("~", apBefore-expectedDynamicAP, 0.001), "Dynamic AP bonus should be removed")
			Expect(playerHealth.GetBonusArmor()).To(BeNumerically("~", armorBefore-expectedDynamicResists, 0.001), "Dynamic Armor bonuses should be removed")
			Expect(playerHealth.GetBonusMR()).To(BeNumerically("~", mrBefore-expectedDynamicResists, 0.001), "Dynamic MR bonus should be removed")
		})
	})
})
