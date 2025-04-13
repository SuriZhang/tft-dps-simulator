package components

// Identity contains champion identification information
type ChampionInfo struct {
	ApiName  string
	Name      string
	Cost      int
	StarLevel int
}

// NewIdentity creates an Identity component
func NewChampionInfo(apiName string, name string, cost, starLevel int) ChampionInfo {
	return ChampionInfo{
		ApiName:  apiName,
		Name:      name,
		Cost:      cost,
		StarLevel: starLevel,
	}
}