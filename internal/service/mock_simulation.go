package service

import (
    "tft-dps-simulator/internal/core/components"
)

// MockRunSimulation returns a fixed dummy response regardless of input
func (s *SimulationService) MockRunSimulation(requestChampions []BoardChampion) (*RunSimulationResponse, error) {
    // Create dummy response
    resp := &RunSimulationResponse{
        Results: []ChampionSimulationResult{
            {
                ChampionApiName: "TFT14_Jax",
                DamageStats: components.DamageStats{
                    TotalDamage:           3500.0,
                    DamagePerSecond:       116.67,
                    TotalADDamage:         2200.0,
                    TotalAPDamage:         1000.0,
                    TotalTrueDamage:       300.0,
                    TotalAutoAttackCounts: 15,
                    TotalSpellCastCounts:  2,
                },
            },
            {
                ChampionApiName: "TFT14_KogMaw",
                DamageStats: components.DamageStats{
                    TotalDamage:           4200.0,
                    DamagePerSecond:       140.0,
                    TotalADDamage:         2800.0,
                    TotalAPDamage:         1200.0,
                    TotalTrueDamage:       200.0,
                    TotalAutoAttackCounts: 18,
                    TotalSpellCastCounts:  3,
                },
            },
        },
    }

    return resp, nil
}