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
		world.AddComponent(player, components.NewSpell("TestLockoutSpell", "",40,0, 1.0)) // Name, ManaCost, MaxMana, startup=0s, recovery=1.0s

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
		targetHealth.SetCurrentHP(550.00) // Ensure target is alive

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

		// --- Reset Timers/Recoverys ---
		playerAttack.SetLastAttackTime(0) // Ready to attack immediately (if not locked out)
		playerSpell.SetCurrentRecovery(0)      // Not locked out by default
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

			// Assertions on the event
			Expect(eventBus.EnqueuedEvents).To(HaveLen(1), "One AttackLandedEvent should be enqueued")
			event, ok := eventBus.GetLastEvent().(eventsys.AttackLandedEvent)
			Expect(ok).To(BeTrue(), "Enqueued event should be of type AttackLandedEvent")
			Expect(event.Source).To(Equal(player), "Event source should be the player")
			Expect(event.Target).To(Equal(target), "Event target should be the dummy")
			Expect(event.BaseDamage).To(BeNumerically("~", expectedBaseDamage, 0.01), "Event base damage should match player's final AD")
			// We could also check event.Timestamp if needed
			Expect(event.Timestamp).To(BeNumerically("~", attackDelay, 0.01), "Event timestamp should match the simulation time when the attack landed")
		})

		It("should enqueue subsequent attacks based on attack speed", func() {
            Expect(playerAttack).NotTo(BeNil())
            Expect(playerAttack.GetFinalAttackSpeed()).To(BeNumerically(">", 0))
            attackDelay := 1.0 / playerAttack.GetFinalAttackSpeed() // Approx 1.818s
            dt := 0.1 // Use a small dt for finer control

            // Simulate until the first attack lands
            var firstAttackTime float64
            var currentTimeOffset float64 = 0.0 // Track total simulated time
            for {
                autoAttackSystem.TriggerAutoAttack(dt)
                currentTimeOffset += dt
                if len(eventBus.EnqueuedEvents) > 0 {
                    Expect(eventBus.EnqueuedEvents).To(HaveLen(1)) // Ensure only one event so far
                    event, ok := eventBus.GetLastEvent().(eventsys.AttackLandedEvent)
                    Expect(ok).To(BeTrue())
                    firstAttackTime = event.Timestamp
                    Expect(firstAttackTime).To(BeNumerically("~", attackDelay, dt*1.5), "First attack should land around the expected delay") // Check first attack time roughly
                    break // Exit loop once the first event is enqueued
                }
                // Safety break
                if currentTimeOffset > attackDelay*2 {
                    Fail("First attack event was not generated within expected time")
                }
            }

            Expect(eventBus.EnqueuedEvents).To(HaveLen(1), "Should have 1 event after first attack time")
            eventBus.ClearEvents() // Clear events after checking the first one

            // Calculate expected second attack time based on the *actual* first attack time
            expectedSecondAttackTime := firstAttackTime + attackDelay

            // Simulate time slightly less than the *next* attack delay, starting from the first attack time
            timeToSimulateBeforeSecondAttack := (expectedSecondAttackTime - currentTimeOffset) - dt/2

            // Only simulate if there's significant time before the next expected attack
            if timeToSimulateBeforeSecondAttack > dt/10 { // Use a small threshold
                for t := 0.0; t < timeToSimulateBeforeSecondAttack; t += dt {
                    autoAttackSystem.TriggerAutoAttack(dt)
                    // Check that the event hasn't landed *prematurely*
                    Expect(eventBus.EnqueuedEvents).To(BeEmpty(), "No event should be enqueued before the second attack time (simulating step %.3f / %.3f)", t+dt, timeToSimulateBeforeSecondAttack)
                    currentTimeOffset += dt
                }
            } else {
                // // If the next attack is expected very soon, log it but don't loop
                // log.Printf("Skipping intermediate simulation loop as expected second attack time (%.3f) is very close to current time (%.3f)", expectedSecondAttackTime, currentTimeOffset)
            }


            // Simulate the time step where the second attack should land (or has just landed)
            // It might have landed in the last step of the previous loop if dt aligned perfectly
            if len(eventBus.EnqueuedEvents) == 0 { // Only run this step if the event hasn't already landed
                autoAttackSystem.TriggerAutoAttack(dt)
                currentTimeOffset += dt
            }


            Expect(eventBus.EnqueuedEvents).To(HaveLen(1), "One AttackLandedEvent should be enqueued for the second attack")
            event, ok := eventBus.GetLastEvent().(eventsys.AttackLandedEvent)
            Expect(ok).To(BeTrue())
            Expect(event.Source).To(Equal(player)) // Check second event details
            Expect(event.Target).To(Equal(target))
            Expect(event.BaseDamage).To(BeNumerically("~", expectedBaseDamage, 0.01))
            // Check the timestamp against the calculated expected time
            Expect(event.Timestamp).To(BeNumerically("~", expectedSecondAttackTime, dt*1.5), "Second event timestamp should be around firstAttackTime + attackDelay")
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
			playerHealth.SetCurrentHP(0) // Attacker starts dead
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
			targetHealth.SetCurrentHP(0) // Target starts dead
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
           
            Expect(playerSpell).NotTo(BeNil())
            playerSpell.SetCurrentRecovery(lockoutTime) // Apply lockout
			Expect(playerSpell.GetCurrentRecovery()).To(BeNumerically("~", lockoutTime, 0.01), "Initial recovery should be set to lockout time")
        })

        It("should not enqueue an AttackLandedEvent during the lockout period", func() {
            dt := 0.1
            currentTime := 0.0
            for t := 0.0; t < lockoutTime-dt/2; t += dt {
                currentCD := playerSpell.GetCurrentRecovery()
                playerSpell.SetCurrentRecovery(math.Max(0, currentCD-dt)) // Manually decrease recovery
                autoAttackSystem.TriggerAutoAttack(dt) // Use Update
                currentTime += dt
                Expect(eventBus.EnqueuedEvents).To(BeEmpty(), "No attack event while locked out at time %.2f", currentTime)
            }
        })

        It("should enqueue an AttackLandedEvent one cycle after the lockout period ends", func() {
            dt := 0.1
            currentTime := 0.0
            attackCycle := 1.0 / playerAttack.GetFinalAttackSpeed() // Should be 1.0s
            // Calculate the expected landing time: lockout ends + one full attack cycle
            expectedAttackTime := lockoutTime + attackCycle // Should be 1.0 + 1.0 = 2.0

            // Simulate slightly past the expected attack time
            simulationEndTime := expectedAttackTime + dt

			eventBus.ClearEvents() // Clear any previous events

			t := 0.0
            for t = 0.0; t < simulationEndTime; t += dt {
				autoAttackSystem.TriggerAutoAttack(dt) // Attack system runs
                // Decrease recovery only while it's positive
                currentCD := playerSpell.GetCurrentRecovery()
				currentTime += dt
                
				if currentCD > 0 {
                     playerSpell.SetCurrentRecovery(math.Max(0, currentCD-dt))
                }
                // Check that no event happens *before* the expected time (allowing for dt tolerance)
                if currentTime < expectedAttackTime - dt/2 {
                     Expect(eventBus.EnqueuedEvents).To(BeEmpty(), "No attack event should occur before expected time %.2f (current: %.2f)", expectedAttackTime, currentTime)
                }
            }
			
			Expect(t).To(BeNumerically("~", 2.1, dt), "Simulation should end around 2.1s")

            // --- Assertions after the loop ---
            Expect(playerSpell.GetCurrentRecovery()).To(BeNumerically("~", 0.0, dt/10), "Recovery should be zero long after lockout ended")

            // Check that exactly one event occurred
            Expect(eventBus.EnqueuedEvents).To(HaveLen(1), "One attack event should be enqueued after lockout + one cycle")

            event, ok := eventBus.GetLastEvent().(eventsys.AttackLandedEvent)
            Expect(ok).To(BeTrue())

            // Check the timestamp matches the expected time (lockout end + 1 cycle)
            Expect(event.Timestamp).To(BeNumerically("~", expectedAttackTime, dt*1.5), "Timestamp should be lockout time + attack cycle time") // Allow slightly larger tolerance due to dt steps
        })
	})

})
