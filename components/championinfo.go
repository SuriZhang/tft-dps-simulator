package components

// Identity contains champion identification information
type ChampionInfo struct {
	Name      string
	Cost      int
	StarLevel int
}

// NewIdentity creates an Identity component
func NewChampionInfo(name string, cost, starLevel int) ChampionInfo {
	return ChampionInfo{
		Name:      name,
		Cost:      cost,
		StarLevel: starLevel,
	}
}