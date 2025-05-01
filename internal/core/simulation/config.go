package simulation

import "fmt"

// SimulationConfig holds all configuration parameters for the simulation
type SimulationConfig struct {
	// Time settings
	MaxTime  float64 // Maximum simulation time in seconds
	TimeStep float64 // Simulation time step in seconds

	// Simulation behavior flags
	DebugMode         bool    // Enables detailed logging during simulation
	ReportingInterval float64 // How often to output status (in simulation seconds)

	// System behavior controls
	EnableAutoAttacks  bool // Whether auto-attack system is active
	EnableSpellCasts   bool // Whether spell-cast system is active
	EnableItemEffects  bool // Whether item procs/effects are active
	EnableTraitEffects bool // Whether trait effects are active

	// Performance settings
	MaxEntities        int  // Upper limit on entities for memory pre-allocation
	ParallelProcessing bool // Whether to use goroutines for system updates
}

// DefaultConfig returns a configuration with sensible defaults
func DefaultConfig() SimulationConfig {
	return SimulationConfig{
		MaxTime:            30.0,
		TimeStep:           0.1,
		DebugMode:          false,
		ReportingInterval:  5.0,
		EnableAutoAttacks:  true,
		EnableSpellCasts:   true,
		EnableItemEffects:  true,
		EnableTraitEffects: true,
		MaxEntities:        100,
		ParallelProcessing: false,
	}
}

// Validate checks if the configuration has valid values
func (c SimulationConfig) Validate() error {
	if c.MaxTime <= 0 {
		return fmt.Errorf("MaxTime must be positive")
	}
	if c.TimeStep <= 0 {
		return fmt.Errorf("TimeStep must be positive")
	}
	if c.TimeStep > c.MaxTime {
		return fmt.Errorf("TimeStep cannot be larger than MaxTime")
	}
	return nil
}

// WithMaxTime returns a copy of the config with updated max time
func (c SimulationConfig) WithMaxTime(maxTime float64) SimulationConfig {
	c.MaxTime = maxTime
	return c
}

// WithTimeStep returns a copy of the config with updated time step
func (c SimulationConfig) WithTimeStep(timeStep float64) SimulationConfig {
	c.TimeStep = timeStep
	return c
}

// WithDebugMode returns a copy of the config with debug mode set
func (c SimulationConfig) WithDebugMode(enabled bool) SimulationConfig {
	c.DebugMode = enabled
	return c
}

// WithReportingInterval returns a copy of the config with updated reporting interval
func (c SimulationConfig) WithReportingInterval(interval float64) SimulationConfig {
	c.ReportingInterval = interval
	return c
}
