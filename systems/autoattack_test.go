package systems_test

import (
	"math" // Import math package

	"github.com/suriz/tft-dps-simulator/components"
	"github.com/suriz/tft-dps-simulator/ecs"
	"github.com/suriz/tft-dps-simulator/factory"
	"github.com/suriz/tft-dps-simulator/systems"
	eventsys "github.com/suriz/tft-dps-simulator/systems/events"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

// --- Test Suite ---
var _ = Describe("AutoAttackSystem", func() {
	var (
		world            *ecs.World
		eventBus         *MockEventBus // Use the mock event bus
		championFactory  *factory.ChampionFactory
		autoAttackSystem *systems.AutoAttackSystem
		player           ecs.Entity // Blue Golem
		target           ecs.Entity // Training Dummy
		playerAttack     *components.Attack
		playerCrit       *components.Crit
		playerSpell      *components.Spell // Added for lockout test
		targetHealth     *components.Health // Still needed to check if target is alive initially
		ok               bool
		// Base AD for event checking
		expectedBaseDamage float64 = 55.00
	)

	BeforeEach(func() {
		world = ecs.NewWorld()
		eventBus = NewMockEventBus() // Initialize mock bus
		championFactory = factory.NewChampionFactory(world)
		autoAttackSystem = systems.NewAutoAttackSystem(world, eventBus) // Pass bus to system

		// --- Create Player (Blue Golem) ---
		var err error
		player, err = championFactory.CreatePlayerChampion("TFT_BlueGolem", 1)
		Expect(err).NotTo(HaveOccurred())
		world.AddComponent(player, components.NewPosition(0, 1))
		// --- Add Spell Component ---
		world.AddComponent(player, components.NewSpell("TestLockoutSpell", "",40,1.0)) // Name, ManaCost, MaxMana, Cooldown=1.0s

		// --- Create Target (Training Dummy) ---
		target, err = championFactory.CreateEnemyChampion("TFT_TrainingDummy", 1)
		Expect(err).NotTo(HaveOccurred())
		world.AddComponent(target, components.NewPosition(1, 1))

		// --- Get Components ---
		playerAttack, ok = world.GetAttack(player)
		Expect(ok).To(BeTrue())
		Expect(playerAttack).NotTo(BeNil())
		playerCrit, ok = world.GetCrit(player)
		Expect(ok).To(BeTrue())
		Expect(playerCrit).NotTo(BeNil())
		playerSpell, ok = world.GetSpell(player) // Get the spell component
		Expect(ok).To(BeTrue())
		Expect(playerSpell).NotTo(BeNil())
		targetHealth, ok = world.GetHealth(target) // Still get health for setup/checks
		Expect(ok).To(BeTrue())
		Expect(targetHealth).NotTo(BeNil())

		// --- Manually Set Final Stats ---
		playerAttack.SetFinalAD(expectedBaseDamage) // Use the variable
		playerAttack.SetFinalAttackSpeed(0.550)
		playerAttack.SetFinalRange(1.00)
		playerAttack.SetFinalDamageAmp(0.00)
		
		playerCrit.SetFinalCritChance(0.25)
		playerCrit.SetFinalCritMultiplier(1.40)
		targetHealth.SetFinalArmor(30.00)
		targetHealth.SetFinalMR(30.00)
		targetHealth.SetFinalDurability(0.00)
		targetHealth.SetCurrentHealth(550.00) // Ensure target is alive

		// --- Set Positions ---
		playerPosition, ok := world.GetPosition(player)
		Expect(ok).To(BeTrue())
		Expect(playerPosition).NotTo(BeNil())
		playerPosition.SetPosition(0, 0)

		targetPosition, ok := world.GetPosition(target)
		Expect(ok).To(BeTrue())
		Expect(targetPosition).NotTo(BeNil())
		targetPosition.SetPosition(1, 0)

		// --- Ensure Dummy doesn't attack ---
		targetAttack, ok := world.GetAttack(target)
		if ok {
			Expect(targetAttack).NotTo(BeNil())
			targetAttack.SetBaseAttackSpeed(0)
			targetAttack.SetFinalAttackSpeed(0)
		}

		// --- Initialize System ---
		// Pass 0.0 as initial time, it will be updated in tests

		// --- Reset Timers/Cooldowns ---
		playerAttack.SetLastAttackTime(0) // Ready to attack immediately (if not locked out)
		playerSpell.SetCurrentCooldown(0)      // Not locked out by default
	})

	Context("when a valid target is in range", func() {
		It("should enqueue an AttackLandedEvent after the initial delay", func() {
			Expect(playerAttack).NotTo(BeNil())
			Expect(playerAttack.GetFinalAttackSpeed()).To(BeNumerically(">", 0))
			attackDelay := 1.0 / playerAttack.GetFinalAttackSpeed()
			dt := 0.1

			// Simulate time slightly less than the attack delay
			for t := 0.0; t < attackDelay-dt/2; t += dt {
				autoAttackSystem.TriggerAutoAttack(dt)
				Expect(eventBus.EnqueuedEvents).To(BeEmpty(), "No event should be enqueued before the attack time")
			}

			// Simulate the time step where the attack should land
			autoAttackSystem.TriggerAutoAttack(dt)
			// Calculate the simulation time when the event is generated
			// Time before this update = 18 iterations * 0.1 dt = 1.8
			// Current time = 1.8 + 0.1 = 1.9
			expectedEventTime := 1.9

			// Assertions on the event
			Expect(eventBus.EnqueuedEvents).To(HaveLen(1), "One AttackLandedEvent should be enqueued")
			event, ok := eventBus.GetLastEvent().(eventsys.AttackLandedEvent)
			Expect(ok).To(BeTrue(), "Enqueued event should be of type AttackLandedEvent")
			Expect(event.Source).To(Equal(player), "Event source should be the player")
			Expect(event.Target).To(Equal(target), "Event target should be the dummy")
			Expect(event.BaseDamage).To(BeNumerically("~", expectedBaseDamage, 0.01), "Event base damage should match player's final AD")
			// We could also check event.Timestamp if needed
			Expect(event.Timestamp).To(BeNumerically("~", expectedEventTime, 0.01), "Event timestamp should match the simulation time when the attack landed")
		})

		It("should enqueue subsequent attacks based on attack speed", func() {
			Expect(playerAttack).NotTo(BeNil())
			Expect(playerAttack.GetFinalAttackSpeed()).To(BeNumerically(">", 0))
			attackDelay := 1.0 / playerAttack.GetFinalAttackSpeed()
			dt := 0.1 // Use a small dt for finer control

			// Simulate until the first attack lands
			var firstAttackTime float64
			for t := 0.0; ; t += dt {
				autoAttackSystem.TriggerAutoAttack(dt)
				if len(eventBus.EnqueuedEvents) > 0 {
					firstAttackTime = eventBus.GetLastEvent().(eventsys.AttackLandedEvent).Timestamp
					break // Exit loop once the first event is enqueued
				}
				// Safety break to prevent infinite loops in case of error
				if t > attackDelay*2 {
					Fail("First attack event was not generated within expected time")
				}
			}

			Expect(eventBus.EnqueuedEvents).To(HaveLen(1), "Should have 1 event after first attack time")
			// --- Manually update LastAttackTime ---
			playerAttack.SetLastAttackTime(firstAttackTime)
			// --- End Manual Update ---
			eventBus.ClearEvents() // Clear events after checking the first one and updating time

			// Simulate time slightly less than the *next* attack delay, starting from the first attack time
			timeToSimulateBeforeSecondAttack := attackDelay - dt/2
			currentTimeOffset := firstAttackTime // Keep track of simulation time relative to start
			for t := 0.0; t < timeToSimulateBeforeSecondAttack; t += dt {
				autoAttackSystem.TriggerAutoAttack(dt)
				Expect(eventBus.EnqueuedEvents).To(BeEmpty(), "No event should be enqueued before the second attack time")
				currentTimeOffset += dt
			}

			// Simulate the time step where the second attack should land
			autoAttackSystem.TriggerAutoAttack(dt)
			currentTimeOffset += dt
			expectedSecondAttackTime := currentTimeOffset

			Expect(eventBus.EnqueuedEvents).To(HaveLen(1), "One AttackLandedEvent should be enqueued for the second attack")
			event, ok := eventBus.GetLastEvent().(eventsys.AttackLandedEvent)
			Expect(ok).To(BeTrue())
			Expect(event.Source).To(Equal(player)) // Check second event details
			Expect(event.Target).To(Equal(target))
			Expect(event.BaseDamage).To(BeNumerically("~", expectedBaseDamage, 0.01))
			Expect(event.Timestamp).To(BeNumerically("~", expectedSecondAttackTime, 0.01), "Second event timestamp should be correct")
		})
	})

	Context("when no valid target is in range", func() {
		BeforeEach(func() {
			Expect(playerAttack).NotTo(BeNil())
			targetPos, ok := world.GetPosition(target)
			Expect(ok).To(BeTrue())
			Expect(targetPos).NotTo(BeNil())
			targetPos.SetX(playerAttack.GetFinalRange() + 1) // Move out of range
		})

		It("should not enqueue an AttackLandedEvent", func() {
			dt := 0.1
			simulationTime := 3.0

			for t := 0.0; t < simulationTime; t += dt {
				autoAttackSystem.TriggerAutoAttack(dt)
			}

			Expect(eventBus.EnqueuedEvents).To(BeEmpty(), "No events should be enqueued when target is out of range")
		})
	})

	Context("when the attacker has 0 Attack Speed", func() {
		BeforeEach(func() {
			Expect(playerAttack).NotTo(BeNil())
			playerAttack.SetBaseAttackSpeed(0)
			playerAttack.SetFinalAttackSpeed(0)
		})

		It("should not enqueue an AttackLandedEvent", func() {
			dt := 0.1
			for t := 0.0; t < 5.0; t += dt {
				autoAttackSystem.TriggerAutoAttack(dt)
			}
			Expect(eventBus.EnqueuedEvents).To(BeEmpty(), "No events should be enqueued when attacker AS is 0")
		})
	})

	Context("when the attacker starts dead", func() {
		BeforeEach(func() {
			playerHealth, ok := world.GetHealth(player)
			Expect(ok).To(BeTrue())
			playerHealth.SetCurrentHealth(0) // Attacker starts dead
		})

		It("should not enqueue an AttackLandedEvent", func() {
			dt := 0.1
			for t := 0.0; t < 5.0; t += dt {
				autoAttackSystem.TriggerAutoAttack(dt)
			}
			Expect(eventBus.EnqueuedEvents).To(BeEmpty(), "No events should be enqueued when attacker starts dead")
		})
	})

	Context("when the target starts dead", func() {
		BeforeEach(func() {
			Expect(targetHealth).NotTo(BeNil())
			targetHealth.SetCurrentHealth(0) // Target starts dead
		})

		It("should not enqueue an AttackLandedEvent", func() {
			// Note: This relies on FindNearestEnemy filtering out dead targets,
			// or the system checking target health before enqueueing.
			// If FindNearestEnemy doesn't filter, the system's internal check should catch it.
			dt := 0.1
			for t := 0.0; t < 5.0; t += dt {
				autoAttackSystem.TriggerAutoAttack(dt)
			}
			Expect(eventBus.EnqueuedEvents).To(BeEmpty(), "No events should be enqueued when target starts dead")
		})
	})

	Context("when the attacker is under spell cast lockout", func() {
        var lockoutTime float64 = 1.0

        BeforeEach(func() {
            Expect(playerAttack).NotTo(BeNil())
            playerAttack.SetFinalAttackSpeed(1.0)
            playerAttack.SetLastAttackTime(-1000.0) // Ready to attack initially
            // ... (target setup) ...
            Expect(playerSpell).NotTo(BeNil())
            playerSpell.SetCurrentCooldown(lockoutTime) // Apply lockout
        })

        It("should not enqueue an AttackLandedEvent during the lockout period", func() {
            dt := 0.1
            currentTime := 0.0
            for t := 0.0; t < lockoutTime-dt/2; t += dt {
                currentCD := playerSpell.GetCurrentCooldown()
                playerSpell.SetCurrentCooldown(math.Max(0, currentCD-dt)) // Manually decrease cooldown
                autoAttackSystem.TriggerAutoAttack(dt) // Use Update
                currentTime += dt
                Expect(eventBus.EnqueuedEvents).To(BeEmpty(), "No attack event while locked out at time %.2f", currentTime)
            }
        })

        It("should enqueue an AttackLandedEvent immediately after the lockout period ends", func() {
            dt := 0.1
            currentTime := 0.0

            // Simulate through the lockout period
            for t := 0.0; t < lockoutTime; t += dt {
                currentCD := playerSpell.GetCurrentCooldown()
                playerSpell.SetCurrentCooldown(math.Max(0, currentCD-dt))
                autoAttackSystem.TriggerAutoAttack(dt) // Attack system runs
                currentTime += dt
                // Check that no event happens *before* the lockout ends
                if t < lockoutTime - dt { // Check all iterations except the last one
                     Expect(eventBus.EnqueuedEvents).To(BeEmpty(), "No attack event should occur strictly before lockout ends at t=%.2f", t+dt)
                }
            }

            // --- Assertions after the loop ---
            // The attack should have happened exactly at lockoutTime (the end of the last step)
            Expect(playerSpell.GetCurrentCooldown()).To(BeNumerically("~", 0.0, dt/10), "Cooldown should be zero after the loop")
            Expect(eventBus.EnqueuedEvents).To(HaveLen(1), "One attack event should be enqueued exactly when lockout ends")

            event, ok := eventBus.GetLastEvent().(eventsys.AttackLandedEvent)
            Expect(ok).To(BeTrue())
            // The timestamp should be the time at the end of the last loop iteration
            expectedAttackTime := lockoutTime
            Expect(event.Timestamp).To(BeNumerically("~", expectedAttackTime, dt), "Timestamp should be exactly the lockout time")
            
			// Should not verify the last attack time here, as it is updated by the event handler
        })
	})

})
