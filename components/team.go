package components

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