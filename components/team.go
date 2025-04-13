package components

import (
    "strings"
)

// Team represents which team a champion belongs to
type Team struct {
    ID int // 0 for player team, 1 for enemy team
}

// NewTeam creates a Team component
func NewTeam(id int) Team {
    return Team{
        ID: id,
    }
}

func (t *Team) String() string {
    var sb strings.Builder // Use strings.Builder for efficiency
    if t.ID == 0 {
        sb.WriteString("  Team: Player\n")
    } else {
        sb.WriteString("  Team: Enemy\n")
    }
    return sb.String()
}