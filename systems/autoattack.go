package systems

import (
    "fmt"
    "math/rand"
    "reflect"

    "github.com/suriz/tft-dps-simulator/components"
    "github.com/suriz/tft-dps-simulator/ecs"
)

// AutoAttackSystem handles champion auto attacks
type AutoAttackSystem struct {
    world   *ecs.World
    // Store current simulation time in seconds
    currentTime float64
}

// NewAutoAttackSystem creates a new auto attack system
func NewAutoAttackSystem(world *ecs.World) *AutoAttackSystem {
    return &AutoAttackSystem{
        world:   world,
        currentTime: 0.0,
    }
}

// Update processes all auto attacks for the current timestep
// deltaTime is the time passed since last update in seconds
func (s *AutoAttackSystem) Update(deltaTime float64) {
    // Update current time
    s.currentTime += deltaTime
    
    // Get all entities with Attack and Health components
    attackType := reflect.TypeOf(components.Attack{})
    healthType := reflect.TypeOf(components.Health{})
    
    // Find all attackers (units that can attack)
    attackers := s.world.GetEntitiesWithComponents(attackType, healthType)
    
    // Find all possible targets (units that can be attacked)
    targets := s.world.GetEntitiesWithComponents(healthType)
    
    // Stop if no targets available
    if len(targets) == 0 {
        return
    }
    
    // Process each attacking entity
    for _, attacker := range attackers {
        // Get attack component
        attackComp, _ := s.world.GetComponent(attacker, attackType)
        attack := attackComp.(components.Attack)
        
        // Check if attack is off cooldown
        attackCooldown := 1.0 / attack.Speed // Time between attacks in seconds
        if s.currentTime - attack.LastAttackTime < attackCooldown {
            continue // Not ready to attack yet
        }
        
        // Find a target (in a real implementation, this would use targeting logic)
        // For now, just pick the first target
        if len(targets) == 0 {
            continue
        }
        target := targets[0]
        
        // Perform the attack
        s.performAttack(attacker, target, attack)
        
        // Update last attack time
        attack.LastAttackTime = s.currentTime
        s.world.AddComponent(attacker, attack)
    }
}

// performAttack handles the logic of an entity attacking another
func (s *AutoAttackSystem) performAttack(attacker, target ecs.Entity, attack components.Attack) {
    // Get champion info
    infoComp, hasInfo := s.world.GetComponent(attacker, reflect.TypeOf(components.ChampionInfo{}))
    attackerName := "Unknown"
    if hasInfo {
        info := infoComp.(components.ChampionInfo)
        attackerName = info.Name
    }
    
    // Get target info
    targetInfoComp, hasTargetInfo := s.world.GetComponent(target, reflect.TypeOf(components.ChampionInfo{}))
    targetName := "Unknown"
    if hasTargetInfo {
        targetInfo := targetInfoComp.(components.ChampionInfo)
        targetName = targetInfo.Name
    }
    
    // Get target health
    targetHealthComp, hasHealth := s.world.GetComponent(target, reflect.TypeOf(components.Health{}))
    if !hasHealth {
        return
    }
    targetHealth := targetHealthComp.(components.Health)
    
    // Calculate damage
    damage := attack.Damage
    
    // Check for critical hit
    isCrit := rand.Float64() < attack.CritChance
    if isCrit {
        damage = damage * attack.CritMultiplier
    }
    
    // Apply armor formula: damage * (100 / (100 + armor))
    // This gives diminishing returns on armor
    damageReduction := 100.0 / (100.0 + targetHealth.Armor)
    finalDamage := damage * damageReduction
    
    // Apply damage to target
    targetHealth.Current -= finalDamage
    
    // Update target health
    s.world.AddComponent(target, targetHealth)
    
    // Debug output
    critText := ""
    if isCrit {
        critText = "CRIT! "
    }
    
    fmt.Printf("%s attacks %s for %.1f damage %s(%.1f HP remaining)\n",
        attackerName, targetName, finalDamage, critText, targetHealth.Current)
    
    // Check if target died
    if targetHealth.Current <= 0 {
        fmt.Printf("%s has been defeated!\n", targetName)
    }
}
