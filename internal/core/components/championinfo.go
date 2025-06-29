package components

// Identity contains champion identification information
type ChampionInfo struct {
	ApiName   string
	Name      string
	Cost      int
	StarLevel int
	role 	string
}

// NewIdentity creates an Identity component
func NewChampionInfo(apiName, name, role string, cost, starLevel int) ChampionInfo {
	return ChampionInfo{
		ApiName:   apiName,
		Name:      name,
		Cost:      cost,
		role : role,
		StarLevel: starLevel,
	}
}
