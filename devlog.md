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