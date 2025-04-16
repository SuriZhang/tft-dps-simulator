package managers_test

import (
	"github.com/suriz/tft-dps-simulator/components"
	"github.com/suriz/tft-dps-simulator/data"
	"github.com/suriz/tft-dps-simulator/ecs"
	"github.com/suriz/tft-dps-simulator/factory"
	"github.com/suriz/tft-dps-simulator/managers"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("EquipmentManager", func() {
	var (
		world            *ecs.World
		equipmentManager *managers.EquipmentManager
		champion         ecs.Entity
		bfSword          *data.Item
		deathblade       *data.Item
		inifityEdge      *data.Item
		tear             *data.Item
		bluebuff 	 *data.Item
	)

	BeforeEach(func() {
		world = ecs.NewWorld()
		championFactory := factory.NewChampionFactory(world)
		equipmentManager = managers.NewEquipmentManager(world)
		// statCalculationSystem := systems.NewStatCalculationSystem(world)
		// abilityCritSystem := itemsys.NewAbilityCritSystem(world)
		// baseStaticItemSystem := itemsys.NewBaseStaticItemSystem(world)

		// Create a champion entity with necessary components
		champion, _ = championFactory.CreatePlayerChampion("TFT_BlueGolem", 1)

		// Get some item data for testing
		bfSword = data.GetItemByApiName("TFT_Item_BFSword")
		Expect(bfSword).NotTo(BeNil(), "BF Sword item data should be loaded")
		deathblade = data.GetItemByApiName("TFT_Item_Deathblade") // Assuming Deathblade is unique
		Expect(deathblade).NotTo(BeNil(), "Deathblade item data should be loaded")
		Expect(deathblade.Unique).To(BeFalse(), "Deathblade should be marked as unique for this test") // Verify assumption
		tear = data.GetItemByApiName("TFT_Item_TearOfTheGoddess")
		Expect(tear).NotTo(BeNil(), "Tear item data should be loaded")
		inifityEdge = data.GetItemByApiName("TFT_Item_InfinityEdge")
		Expect(inifityEdge).NotTo(BeNil(), "Infinity Edge item data should be loaded")
		bluebuff = data.GetItemByApiName("TFT_Item_BlueBuff")
		Expect(bluebuff).NotTo(BeNil(), "Blue Buff item data should be loaded")

	})

	Describe("AddItemToChampion", func() {
		Context("when adding a valid item", func() {
			It("should add the item to the champion's equipment", func() {
				err := equipmentManager.AddItemToChampion(champion, bfSword.ApiName)
				Expect(err).NotTo(HaveOccurred())

				eq, ok := world.GetEquipment(champion)
				Expect(ok).To(BeTrue())
				Expect(eq.Items).To(HaveLen(1))
				Expect(eq.Items[0].ApiName).To(Equal(bfSword.ApiName))
			})

			It("should update the ItemEffect component", func() {
				err := equipmentManager.AddItemToChampion(champion, bfSword.ApiName)
				Expect(err).NotTo(HaveOccurred())

				itemEffect, ok := world.GetItemEffect(champion)
				Expect(ok).To(BeTrue())
				// Assuming BF Sword gives BonusPercentAD
				Expect(itemEffect.GetBonusPercentAD()).To(BeNumerically(">", 0))
			})

			It("should allow adding multiple non-unique items", func() {
				err := equipmentManager.AddItemToChampion(champion, bfSword.ApiName)
				Expect(err).NotTo(HaveOccurred())
				err = equipmentManager.AddItemToChampion(champion, bfSword.ApiName) // Add second BF Sword
				Expect(err).NotTo(HaveOccurred())

				eq, ok := world.GetEquipment(champion)
				Expect(ok).To(BeTrue())
				Expect(eq.Items).To(HaveLen(2))
				Expect(eq.GetItemCount(bfSword.ApiName)).To(Equal(2))

				itemEffect, ok := world.GetItemEffect(champion)
				Expect(ok).To(BeTrue())
				// Check if stats are doubled (adjust based on actual BF sword stats)
				expectedStat := bfSword.Effects["AD"] * 2 // Assuming "AD" is the key for BonusPercentAD
				Expect(itemEffect.GetBonusPercentAD()).To(BeNumerically("~", expectedStat))
			})
		})

		Context("when adding an item to full equipment", func() {
			BeforeEach(func() {
				// Fill equipment slots
				Expect(equipmentManager.AddItemToChampion(champion, bfSword.ApiName)).To(Succeed())
				Expect(equipmentManager.AddItemToChampion(champion, bfSword.ApiName)).To(Succeed())
				Expect(equipmentManager.AddItemToChampion(champion, tear.ApiName)).To(Succeed())
			})

			It("should return an error", func() {
				err := equipmentManager.AddItemToChampion(champion, bfSword.ApiName) // Try adding a 4th item
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("no space to add item"))

				eq, ok := world.GetEquipment(champion)
				Expect(ok).To(BeTrue())
				Expect(eq.Items).To(HaveLen(components.MaxItems)) // Should still have only MaxItems
			})
		})

		Context("when adding a unique item that is already equipped", func() {
			BeforeEach(func() {
				Expect(equipmentManager.AddItemToChampion(champion, bluebuff.ApiName)).To(Succeed())
			})

			It("should return an error", func() {
				err := equipmentManager.AddItemToChampion(champion, bluebuff.ApiName) 
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("is unique and already equipped"))

				eq, ok := world.GetEquipment(champion)
				Expect(ok).To(BeTrue())
				Expect(eq.Items).To(HaveLen(1)) 
			})
		})

		Context("when adding a non-existent item", func() {
			It("should return an error", func() {
				err := equipmentManager.AddItemToChampion(champion, "TFT_Item_NonExistent")
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("not found"))
			})
		})

		Context("when the champion entity is invalid or lacks components", func() {
			It("should return an error if champion has no Equipment component", func() {
				invalidChamp := world.NewEntity()
				err := world.AddComponent(invalidChamp, &components.ChampionInfo{Name: "NoEquipChamp"})
				Expect(err).NotTo(HaveOccurred())

				err = equipmentManager.AddItemToChampion(invalidChamp, bfSword.ApiName)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("has no Equipment component"))
			})
			// Add similar test for missing ChampionInfo if needed
		})
	})

	Describe("RemoveItemFromChampion", func() {
		BeforeEach(func() {
			// Add some items to remove
			Expect(equipmentManager.AddItemToChampion(champion, bfSword.ApiName)).To(Succeed())
			Expect(equipmentManager.AddItemToChampion(champion, tear.ApiName)).To(Succeed())
			Expect(equipmentManager.AddItemToChampion(champion, bfSword.ApiName)).To(Succeed()) // Add a second BF Sword
		})

		Context("when removing an existing item", func() {
			It("should remove the first instance of the item", func() {
				err := equipmentManager.RemoveItemFromChampion(champion, bfSword.ApiName)
				Expect(err).NotTo(HaveOccurred())

				eq, ok := world.GetEquipment(champion)
				Expect(ok).To(BeTrue())
				Expect(eq.Items).To(HaveLen(2))
				Expect(eq.GetItemCount(bfSword.ApiName)).To(Equal(1)) // One BF Sword should remain
				Expect(eq.HasItem(tear.ApiName)).To(BeTrue())         // Tear should remain
			})

			It("should update the ItemEffect component", func() {
				// Get initial effect
				itemEffectBefore, _ := world.GetItemEffect(champion)
				initialADBonus := itemEffectBefore.GetBonusPercentAD()
				Expect(initialADBonus).To(BeNumerically(">", 0))

				// Remove one BF Sword
				err := equipmentManager.RemoveItemFromChampion(champion, bfSword.ApiName)
				Expect(err).NotTo(HaveOccurred())

				// Check effect after removal
				itemEffectAfter, ok := world.GetItemEffect(champion)
				Expect(ok).To(BeTrue())
				// AD Bonus should be halved (approx)
				expectedStat := bfSword.Effects["AD"] // Assuming "AD" is the key
				Expect(itemEffectAfter.GetBonusPercentAD()).To(BeNumerically("~", expectedStat))
			})

			It("should remove the last item correctly", func() {
				Expect(equipmentManager.RemoveItemFromChampion(champion, bfSword.ApiName)).To(Succeed())
				Expect(equipmentManager.RemoveItemFromChampion(champion, tear.ApiName)).To(Succeed())
				Expect(equipmentManager.RemoveItemFromChampion(champion, bfSword.ApiName)).To(Succeed()) // Remove the last item

				eq, ok := world.GetEquipment(champion)
				Expect(ok).To(BeTrue())
				Expect(eq.Items).To(HaveLen(0))

				// Item effects should be reset (or close to zero)
				itemEffect, ok := world.GetItemEffect(champion)
				Expect(ok).To(BeTrue())
				Expect(itemEffect.GetBonusPercentAD()).To(BeNumerically("~", 0.0))
				Expect(itemEffect.GetBonusInitialMana()).To(BeNumerically("~", 0.0)) // Assuming Tear gives Mana
			})
		})

		Context("when removing an item that is not equipped", func() {
			It("should return an error", func() {
				err := equipmentManager.RemoveItemFromChampion(champion, deathblade.ApiName) // Deathblade was not added
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("not found in champion"))

				eq, ok := world.GetEquipment(champion)
				Expect(ok).To(BeTrue())
				Expect(eq.Items).To(HaveLen(3)) // Items should remain unchanged
			})
		})

		Context("when the champion entity is invalid or lacks components", func() {
			It("should return an error if champion has no Equipment component", func() {
				invalidChamp := world.NewEntity()
				err := world.AddComponent(invalidChamp, &components.ChampionInfo{Name: "NoEquipChamp"})
				Expect(err).NotTo(HaveOccurred())

				err = equipmentManager.RemoveItemFromChampion(invalidChamp, bfSword.ApiName)
				Expect(err).To(HaveOccurred())
				Expect(err.Error()).To(ContainSubstring("has no Equipment component"))
			})
			// Add similar test for missing ChampionInfo if needed
		})
	})

	// calculateAndUpdateItemEffects is implicitly tested by Add/Remove,
	// but you could add direct tests if needed, especially for complex aggregation logic.
})
