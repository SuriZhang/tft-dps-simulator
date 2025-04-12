# tft-dps-simulator

``` files structure
/simulator
│
├── main.go                      # Entry point to run sim
│
├── /ecs                         # Core ECS framework
│   ├── world.go                 # Holds all entities + components
│   ├── entity.go                # ID generation & basic helpers
│
├── /components                 # Pure data structs (no logic)
│   ├── health.go               # Health component
│   ├── mana.go                 # Mana component
│   ├── attack.go               # Attack speed/damage
│   ├── traits.go               # Trait list component
│   ├── spell.go                # Spell cast function or ID
│   ├── defense.go              # Armor/MR component
│   └── buffs.go                # Trait/item bonuses
│
├── /systems                    # Game logic processors
│   ├── autoattack.go           # Handles base attacks
│   ├── ability.go              # Handles mana & spell casts
│   ├── trait.go                # Counts & applies trait bonuses
│   ├── item.go                 # Handles on-hit or passive items
│   └── damage.go               # Applies actual damage logic
│
├── /factory                    # Champion/unit construction
│   ├── kaisa.go                # Creates a specific champion
│   └── champion.go             # General helpers to spawn units
│
├── /data                       # Static game data
│   ├── traits.go               # Trait thresholds & effects
│   ├── champions.go            # Base stats per champ (could load from JSON)
│   ├── items.go                # Item definitions
|   ├── loader.go               # Main loader entry point
|   ├── champions.go            # Handles parsing of champion stats/spells/traits
|   ├── traits.go               # Parses traits & activation thresholds
|   ├── augments.go             # (future) parses augment info
|   ├── models.go                # Common raw struct definitions
│
├── /utils                      # Math helpers, targeting logic, RNG
│   └── targeting.go            # FindNearestEnemy, etc.
```