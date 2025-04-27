package components

type ChampionActionState int

const (
	Casting           ChampionActionState = iota // Currently executing a spell cast
	AttackStartingUp                             // In the wind-up phase of an auto-attack
	AttackRecovering                             // In the recovery phase after an attack landed/fired
	AttackCoolingDown                            // Waiting for the attack speed timer after recovery
	Idle                                         // Not performing any action (can overlap with CC/Stun)
)

var stateNames = map[ChampionActionState]string{
	Casting:           "Casting",
	AttackStartingUp:  "AttackStartingUp",
	AttackCoolingDown: "AttackCoolingDown",
	AttackRecovering:  "AttackRecovering",
	Idle:              "Idle",
}

// State holds the current action state and status effects of a champion.
type State struct {
	// Status Effects
	IsStunned   bool    // Specific flag if stunned (prevents actions)
	CCStartTime float64 // Simulation time when the current CC started
	CCDuration  float64 // Duration of the current CC

	// Action States (Mutually Exclusive except for CC)
	PreviousState ChampionActionState // Previous action state (for logging/debugging)
	CurrentState  ChampionActionState // Current action state (e.g., Idle, Attack, Cast)

	// Timing for Current Action State
	ActionStartTime float64 // Simulation time when the current action state began
	ActionDuration  float64 // Expected duration of the current action state (e.g., cast time, startup time)
}

// NewState creates a default State component.
func NewState() *State {
	return &State{
		PreviousState: Idle,
		CurrentState:  Idle,
		IsStunned:     false,
	}
}

// --- Getters (Add more as needed) ---

func (s *State) GetIsStunned() bool {
	return s.IsStunned
}

// --- Setters / State Transition Helpers (Add more complex logic as needed) ---

// Example: StartAttack sets the state for attack startup.
// Needs current simulation time and calculated startup duration.
func (s *State) StartAttack(currentTime, startupDuration float64) {
	s.PreviousState = s.CurrentState
	s.CurrentState = AttackStartingUp
	s.ActionStartTime = currentTime
	s.ActionDuration = startupDuration
}

func (s *State) StartAttackRecovery(currentTime, recoveryDuration float64) {
	s.PreviousState = s.CurrentState
	s.CurrentState = AttackRecovering
	s.ActionStartTime = currentTime
	s.ActionDuration = recoveryDuration
}

func (s *State) StartAttackCooldown(currentTime, cooldownDuration float64) {
	s.PreviousState = s.CurrentState
	s.CurrentState = AttackCoolingDown
	s.ActionStartTime = currentTime
	s.ActionDuration = cooldownDuration
}

// Example: StartCast sets the state for spell casting.
// Needs current simulation time and calculated cast duration.
func (s *State) StartCast(currentTime, castDuration float64) {
	s.PreviousState = s.CurrentState
	s.CurrentState = Casting
	s.ActionStartTime = currentTime
	s.ActionDuration = castDuration
}

func (s *State) StartActionCheck(currentTime float64) {
    s.PreviousState = s.CurrentState
    s.CurrentState = Idle
    s.ActionStartTime = currentTime
    s.ActionDuration = 0.0 // No action duration in idle state
}