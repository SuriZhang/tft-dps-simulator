// filepath: components/state.go
package components

// State holds the current action state and status effects of a champion.
type State struct {
    // Status Effects
    IsStunned         bool    // Specific flag if stunned (prevents actions)
    CCStartTime       float64 // Simulation time when the current CC started
    CCDuration        float64 // Duration of the current CC

    // Action States (Mutually Exclusive except for CC)
    IsCasting         bool    // Currently executing a spell cast
    IsAttackStartingUp bool    // In the wind-up phase of an auto-attack
    IsAttackRecovering bool    // In the recovery phase after an attack landed/fired
    IsAttackCoolingDown bool    // Waiting for the attack speed timer after recovery
    IsIdle            bool    // Not performing any action (can overlap with CC/Stun)

    // Timing for Current Action State
    ActionStartTime   float64 // Simulation time when the current action state began
    ActionDuration    float64 // Expected duration of the current action state (e.g., cast time, startup time)
}

// NewState creates a default State component.
func NewState() *State {
    return &State{
        IsIdle: true, // Start as idle
        // Initialize other fields to default zero/false values
    }
}

// --- Getters (Add more as needed) ---

func (s *State) GetIsStunned() bool {
    return s.IsStunned
}

func (s *State) GetIsAttackRecovering() bool {
    return s.IsAttackRecovering
}

func (s *State) GetIsAttackCoolingDown() bool {
    return s.IsAttackCoolingDown
}

// --- Setters / State Transition Helpers (Add more complex logic as needed) ---

// Example: StartAttack sets the state for attack startup.
// Needs current simulation time and calculated startup duration.
func (s *State) StartAttack(currentTime, startupDuration float64) {
    s.resetActionStates()
    s.IsAttackStartingUp = true
    s.ActionStartTime = currentTime
    s.ActionDuration = startupDuration
}

// Example: StartCast sets the state for spell casting.
// Needs current simulation time and calculated cast duration.
func (s *State) StartCast(currentTime, castDuration float64) {
    s.resetActionStates()
    s.IsCasting = true
    s.ActionStartTime = currentTime
    s.ActionDuration = castDuration
}

// resetActionStates clears flags before setting a new one.
func (s *State) resetActionStates() {
    s.IsCasting = false
    s.IsAttackStartingUp = false
    s.IsAttackRecovering = false
    s.IsAttackCoolingDown = false
    s.IsIdle = false
}

// TODO: Add methods for handling CC application and expiry.
// TODO: Add methods for transitioning between attack states (Startup -> Recovery -> Cooldown -> Idle/Startup).