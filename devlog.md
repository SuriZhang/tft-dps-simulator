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

Done Today:
- [x] set up project scaffold
- [x] implemented data parsing logic
- [x] implemented basic auto attack system

TODO:
- [ ] clean up main.go, wrap the auto attack testing code to another function or something
- [ ] traits, augments, items, spell cast

