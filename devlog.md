# Random thoughts and learnings along the way

## 20250412

What is ECS (Entity-Component-System)?  
ECS is an architectural pattern commonly used in game engines and simulations. Instead of modeling objects with inheritance (OOP), it decouples data (components) from behavior (systems).

Here’s how it breaks down:

Entity: A unique ID or container (e.g., a champion)

Component: A pure data struct (e.g., Health, Mana, Position, AttackSpeed)

System: A processor that acts on entities with specific components (e.g., a DamageSystem acts on entities with Health and Damage)

Why Consider ECS in Go?
Go’s lack of inheritance makes ECS more natural here than in traditional OOP languages. ECS provides:

Pros:
High Flexibility / Extensibility
Want to add a new status effect, item proc, or trait behavior? Just add a new component or system—no need to tangle with inheritance trees.

Separation of Concerns
You keep logic out of your entities—systems handle all the behavior. This makes it easier to test and reason about.

Performance Friendly
ECS can batch-process entities in tight loops (though that’s more critical in real-time games, still nice for simulation).

Easier Scaling with Complex Interactions
As TFT sets evolve and mechanics get wilder (think augments, modifiers, stacking effects), ECS handles that variety better.

Read about ECS: https://www.richardlord.net/blog/ecs/what-is-an-entity-framework
space station game ECS in practice: https://docs.spacestation14.com/en/robust-toolbox/ecs.html

DONE Today:

- [x] set up project scaffold
- [x] implemented data parsing logic
- [x] implemented basic auto attack system

TODO:

- [ ] clean up main.go, wrap the auto attack testing code to another function or something
- [ ] traits, augments, items, spell cast

## 20250413

Goal:

- implement simple items

handling items:
Equipment Component:

Its main job is to hold the list of specific items equipped by a champion (e.g., [Tear of the Goddess, Needlessly Large Rod]).
The logic for enforcing the 3-item limit happens when you try to add an item to this component (likely within your AddItemToChampion function). So yes, it's central to the process of "putting an item on a champion".
ItemEffect Component:

This component stores the combined passive stat bonuses from all items currently in the Equipment. For example, if Tear gives +15 Mana and Rod gives +10 AP, this component would store { Mana: 15, AbilityPower: 10, ... }.
It does not contain the logic for individual item effects or how they are calculated. It's just a data container for the result of combining passive stats.
ItemSystem (specifically the ApplyStats method):

This system reads the aggregated stats from the ItemEffect component.
It then modifies the champion's effective stats (like Health.Current, Attack.Damage, Attack.Speed) based on these aggregated bonuses and the champion's base stats.
It does not handle active item effects (like Guinsoo's stacking on attack or Statikk Shiv proc). Those would require different systems that react to specific game events (like attacks).
So, to summarize:

You try to add an item -> Check Equipment for space.
If space, add item to Equipment -> Recalculate combined stats -> Store result in ItemEffect.
ItemSystem.ApplyStats runs -> Reads ItemEffect -> Updates champion's final combat stats based on base stats + item bonuses.

The logic for individual item effects, especially those that are not simple passive stat bonuses, should be implemented in dedicated Systems.

Here's a breakdown:

1. systems/items/static should holds the logic for all items that only modifies static champion stats. this is easy.
2. passive dynamic effects items: there are two sub-categories,
   one for the time-dependent effect items (in description it should be "gain X every X seconds"). examples are quick silver, evenshroud, archangel's staff, and adaptive helmet for back two rows  
   another category is the simple effect is triggered on some events (gain X when attacking or taking damge).e xamples are titan's resolve., adaptive helmet for the front two rows
3. complex tiggered effects that has impact on the simulation timeline. e.g. guiunsoo's regebald (it stacks attackspeed and the attackspeed at t+1 should be based on attackspeed at time t)

DONE Today:

- [x] implemented base static item system

TODO:

- [ ] implement StatCalculationSystem to calcluate bonus component stats and update final stats
- [x] maybe refactor item related code from championfactory to itemfactory.
- [ ] clean up main.go, wrap the auto attack testing code to another function or something

## 20250414

Goal: implement IE & JG

IE & JG's effect are the same, unique to all other items.

DONE Today:

- [x] implement StatCalculationSystem to calcluate bonus component stats and update final stats
- [x] implement logic for handling IE & JG

TODO:

- [x] think about the edge case when a champion wears more than 1 IE and/or JG at the same time
- [ ] clean up main.go, wrap the auto attack testing code to another function or something
- [x] write tests

## 20250415-0416

Lessons learned from debugging IE/JG edge cases:  
To prevent minor errors being ignored and causing me more problem in the future, I'll need to write tests.

DONE Today:

- [x] fixed IE & JG logic when a champion wears more than 1 item that can cause their ability to crit.
- [x] added tests!!
- [x] clean up main.go, wrap the logic and simulation config to `/simulation` directory

TODO:

- [ ] implement QickSilver and Archagel's Staff

## 20250417

DONE Today:

- [x] implemented QuickSilver and Archangels Staff, along with all DynamicTime item system.
- [x] Added Spell components, currently not completed, but stores Archagel's bonus AP. will add more details when it comes to implementing the champion Ability/Spell

TODO:

- [ ] implement logic when two component items are added, they form a composition item according to the formula --> not in MVP
- [ ] implement game event systems, sets up scaffold for DynamicEvent-typed items
- [ ] implement DynamicEvent items

## 20250418

Debug...

DONE Today:

- [x] Set up event system scaffold
- [x] Added more tests, and make sure they all passed...
- [x] Correctly updated dynamic time items bonus in simulation loops

TODO:

- [ ] implement very basic spell cast system, to prepare us for more item effects
- [ ] implement DynamicEvent items
- [ ] implement logic when two component items are added, they form a composition item according to the formula --> deprioritized, not in MVP

## 20250419

Simplifed Spell Cast cool down handling:
the term "cooldown" in TFT usually refers to the time after the spell animation finishes before the next spell can be cast. The period where the champion is locked out of auto-attacking is the cast animation time or cast lockout.

Implementing precise cast animation times for every champion adds significant complexity, as these times vary.

Simplified Approach (Using Cooldown as Lockout):

For now, we can implement a simpler version where we treat the Spell.Cooldown value itself as the lockout period during which the champion cannot auto-attack.

SpellCastSystem: When a spell is cast, set Spell.CurrentCooldown to Spell.Cooldown.
AttackSystem: Modify the AttackSystem.Update function. Before an entity performs an auto-attack, check if it has a Spell component. If it does, check if Spell.GetCurrentCooldown() > 0. If the cooldown is active, prevent the entity from auto-attacking in that frame.
Pros:

Relatively simple to implement using existing components and systems.
Achieves the basic goal: casting prevents auto-attacking for a duration.
Cons:

Inaccurate simulation: Uses the spell cooldown duration instead of the actual cast animation time.
Might feel clunky for spells with long cooldowns but short cast times, or vice-versa.
Recommendation:

Given the complexity of true cast times, let's proceed with the simplified approach for now. We can refine it later if needed by adding a dedicated CastTime or AnimationLockout field to the Spell component.

Might be buggy later (revisit later): when handling the edge cases where a champion died jsut before they about to cast a spell or auto attack.

DONE Today:

- [x] implement very basic spell cast system, to prepare us for more item effects
- [x] implement DynamicEvent items system
- [x] implement Titan's Resolve

TODO:

- [ ] implement Adaptive Helmet
- [ ] implement logic when two component items are added, they form a composition item according to the formula --> deprioritized, not in MVP

## 20250420

Goal: implement Guinsoo's Rageblade

Added attack and cast startup and recovery fields, though we dont have data yet, all default to 0. This may cause some logic refinement but let's move on for now.

Now auto attack is split up into startup + recovery periods

Autoattack system logic:
Initialization: An entity starts with LastAttackTime = 0.0 and AttackStartupEndTime = -1.0. AttackCycleEndTime is also often initialized to 0.0 or -1.0.
First Attack Scheduling (Lines 72-81):
This block if attack.GetLastAttackTime() == 0.0 && attack.GetAttackStartupEndTime() == -1.0 is designed to run only once per entity, right at the beginning before any attacks have happened.
It calculates totalCycleTime = 1.0 / attack.GetFinalAttackSpeed().
It sets firstLandingTime = totalCycleTime. This variable represents the desired time for the first attack to land, aligning with your requirement (e.g., 2.0s for AS=0.5).
Crucially, line 75: attack.SetAttackStartupEndTime(firstLandingTime + attack.GetCurrentAttackStartup()) schedules the actual landing time. This is where the confusion lies. It takes the desired landing time (firstLandingTime) and adds the attack's startup time again.
The log message then prints firstLandingTime, which might be misleading because the actual scheduled time stored in AttackStartupEndTime is later.
Subsequent Attack Cycle Start (Lines 85-113):
This block if s.currentTime >= attack.GetAttackCycleEndTime() runs when the entity is ready to start a new attack cycle (meaning the previous one has fully completed). This won't run for the very first attack because AttackCycleEndTime is initially 0 or -1.
It calculates the totalCycleTime and the startupTime for this specific cycle.
It schedules the landing time for this attack: attack.SetAttackStartupEndTime(s.currentTime + startupTime). This means the attack lands startupTime seconds after the current cycle begins.
It schedules when the next cycle can begin: attack.SetAttackCycleEndTime(s.currentTime + totalCycleTime).
It records the time this cycle started: attack.SetLastAttackTime(s.currentTime).
Attack Landing Check (Lines 116 onwards):
This block if attack.GetAttackStartupEndTime() != -1.0 && s.currentTime >= attack.GetAttackStartupEndTime() checks if the current simulation time has reached or passed the scheduled landing time (AttackStartupEndTime).
If it has, it finds a target, checks range, and enqueues the AttackLandedEvent.
It then resets attack.SetAttackStartupEndTime(-1.0) to mark the landing as processed.

Note to self: in auto attack system, the logic for checking CC/interuption during startup is already added (but commented out)

DONE Today:

- [x] added Rageblade
- [x] revamped auto attack system to include startup and recovery times, make sure all tests passed
- [x] when calculating damange, use the finalAD at damage calculation time, instead of at attack landed time
- [x] add mana gain for entity taking damage

TODO:

- [ ] add spell cast startup and recovery into the actual spellcast system
- [ ] revisit logic for mana gain when champion taking damage, this should depend on if the champion's mana is locked during cast (also mana gain from other sources)
- [ ] implement logic when two component items are added, they form a composition item according to the formula --> deprioritized, not in MVP

## 20250426

Revamped Architecture design to address some edge cases and current flaws:
Buff/Debuff system (for a champion unit):

for each buff, we should record the following information:

- source
- duration
- startTime
- buff/debuff specific stat

there’s a function according to each buff/debuff that would apply the buff to the champion

Global Buff (for the game simulation, applies for maybe more than one champion unit):

e.g., effects of ionic spark, Hex Augment (the team get 10% Attack Speed per 3 seconds in the combat)

Damage Statistics:

- per champion, we should record each damage dealt/taken and it’s type (AD/AP/True Damage and from Auto Attack/Spell Cast/Burn/Hurricane/Traits)
- how many spell casted during the simulation

Champion State: should record the state of the champion at time t of the simulation

- possible states are:
  - isUnderStun (independent from the rest, because a champion can be under Stun but still spelling)
  - isCasting
  - isAttackStartingUp
  - isAttackRecovering
  - isAttackCoolingDown
  - isIdle (must be under Stun)
- we should also store the start time and (expected) duration of the state

Champion Action Handler

champion state at any simulation time t:

→ check is stun?

    →Y: pass, do nothing until stun is over

        → N: check is full mana?

            → Y: check isAttackRecovering?

                →Y: pass, do nothing until attack recovery is over

                → N: should cast a spell

            → N: check isAttackCoolingDown?

        → Y: should auto attack

    → N: pass, do nothing until attack CD is over

Champion Action Cycle breakdown:

(1) targeting → (2) attack startup → (3) check if mana full1? → (4) no: trigger AttackFiredEvent → (5) attack recovery → (6) check if mana full2? → (7) no: attack CD → go to (1)

at any mana full check, if yes → (a) spell cast start up (has a higher priority than CC) → (b) trigger SpellCastEvent → (c) spell cast recovery → (d) check attack CD is over? → yes: go to (1)/ no: go to (7) to wait until remaining attack CD is over

Event structure:

- basic attributes: source, target, stats (damage, etc)
- Timestamp (accurate timestamp for the event)
- EnqueueTimestamp (timestamp with a delta t in range U(-10^-5, +10^-5) to help resolve simultaneity)

Simulation Steps:

- before combat:
  1. resolve starting-of-combat effects, including
     1. handle item gain, e.g., Thief’s Gloves, Sponging (Combat start: Up to 6 champions with 1 or fewer items gain a copy of a random completed item from the nearest itemized ally.)
     2. any static item/augments/traits effects
     3. enqueue time effects (e.g., archangel’s staff AP gain every 5s in combat, should enqueue at t=5, t-10, t=15, etc)
     4. other special handlings (e.g., S14 Overlord: The Overlord takes a bite out of the unit in the hex behind him, dealing 40% of their max Health as true damage. He gains 40% of their Health and 25% of their Attack Damage.)
- at t=0, enqueue all champion’s first action (auto attack or cast)
- simulation start:

  ```
  while (! combatEnds) {

  1. evt = EventQueue.dequeue()
  2. set simulation time = evt.Timestamp
  3. handleEvent(evt)
      1. resolve event
      2. enqueue subsequent event (when enqueuing new events to the EventQueue, find it’s position using evt.EnqueueTimestamp)
  4. save evt to RecordQueue (if evt should be saved to help replay/analyze the combat when simulation is over, types of event should be saved TBD)
  }

  combatEnds: 1. simulation time passed 30s; 2. one team has no alive champion units
  ```

DONE Today:

- [x] adapt new event-driven simulation design
- [x] refactored event system

TODO:

- [ ] refactor DynamicTimeSystem to enqueue and handle events
- [ ] fix ChampionActionSystem to correctly enqueue AttackCooldownStartEvent based on champion state
- [ ] fix SimulationTests (currently failing 8 tests)

## 20250427

DONE Today:

- [x] major event-driven refactor is done!!
- [x] refactor DynamicTimeSystem to enqueue and handle events
- [x] fix ChampionActionSystem to correctly enqueue AttackCooldownStartEvent based on champion state
- [x] all tests passed

## 20250501

DONE Today:

- [x] Traits system scaffold, traits are referenced by Name (instead oPf ApiName) in the simulation
- [x] Rapidfire trait effects

Next Step:

1. implement one or two simple abilities to complete the first milestone of MVP
2. refactor the code base to make it a proper web backend
   here is a checklist to guide the refactoring process, ordered roughly from lowest to highest workload:

Setup & Dependencies:

- [x] Add Fiber dependency (go get github.com/gofiber/fiber/v2) [new] [low]
- [x] Add Redis client dependency (e.g., go get github.com/go-redis/redis/v8) [new] [low]
- [ ]Create cmd/server/ directory [new] [low]
- [x] Create internal/ directory [new] [low]
- [ ]Create internal/api/, internal/service/, internal/store/ subdirectories [new] [low]
- [x] Create basic Fiber app setup in main.go [new] [low]
- [x] Implement Redis client initialization logic (e.g., in internal/store/ or main.go) [new] [low]
- [ ] Move one-time data loading (data.InitializeChampions, etc.) from main.go to main.go [refactor] [low]
- [ ] Ensure data loading happens only once on server start [refactor] [low]
- [ ] Define SimulationSetup and ChampionConfig structs (likely in internal/store/ or a shared types package) [new] [low]
- [ ] Define SimulationService struct signature in internal/service/ [new] [low]
- [ ] Define basic API handler function signatures in internal/api/ [new] [low]
- [ ] Instantiate store, service, handlers in main.go [new] [low]
- [x] Add basic server start logic (app.Listen()) in main.go [new] [low]

Core Logic Implementation & Refactoring:

- [x] Move existing simulation packages (ecs, components, systems, simulation, managers, factory, data, utils) into internal/simcore/ [refactor] [medium]
- [x] Update all import paths affected by the move to internal/simcore/ [refactor] [medium]
- [ ] Implement SaveSimulationSetup function in internal/store/ (handle serialization, Redis SET) [new] [medium]
- [ ] Implement GetSimulationSetup function in internal/store/ (handle Redis GET, deserialization) [new] [medium]
- [ ] Implement request body/parameter parsing in API handlers [new] [medium]
- [ ] Implement JSON response formatting in API handlers [new] [medium]
- [ ] Implement basic error handling and HTTP status code responses in API handlers [new] [medium]
- [x] Define Fiber routes in main.go connecting paths to API handlers [new] [medium]
- Implement AddChampion method in SimulationService (call store, update setup) [new] [medium]
- Implement AddItemToChampion method in SimulationService (call store, update setup) [new] [medium]
- Implement ChangeChampionStarLevel method in SimulationService (call store, update setup) [new] [medium]
- [ ] Refactor Simulation (in internal/simcore/simulation) to return structured results instead of printing [refactor] [medium]
- [ ] Implement RunSimulation method in SimulationService (get setup, create world, use factory/managers, run sim, return results) [new] [high]

## 20250502-0504

Mainly focused on getting the frontend up and look nice.

Future TODOs (nice-to-have frontend animations using Animejs):

1. Champion Drag-and-Drop animations

- Animate champion icons smoothly snapping to the board grid with a slight "bounce" effect.

- Add scaling or glow transitions on hover, drag start, and drop.

2. Trait Activation Animation
   When a trait becomes active, show:

- A glowing pulse on the synergy tracker bar.

- A quick expanding-ring or lightwave animation on the champion avatars that contribute to the trait.

3. Stat Panel & Tooltip Reveal

- Fade in tooltips with anime.js when hovering champions or items.

- Slide stat panels or damage breakdowns into view with ease-in curves.

DONE:

- [x] Frontend Scaffold (main layout, core components)
- [x] Data loading from json
- [x] wrote script to download champion square icons from community dragon
- [x] frontend trait activation
- [x] mock simulation endpoints, triggered by clicking combat start button, damage stats panel implemented

WIP:
- [ ] basic run simulation endpoints
- 

TODO:

- [ ] frontend: champion, item, trait tooltip polishing, properly parse html tags in the description/effects
- [ ] endpoint: add champion to board (show champion stats, trait activation)
- [ ] endpoint: add item to champion (show updated champion stats)
- [ ] endpoint: advanced run simulation endpoint, show time sequence data
