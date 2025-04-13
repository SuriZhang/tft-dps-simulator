package data

type TFTSetData struct {
	SetData []Set `json:"setData"`
}

type Set struct {
	Champions []Champion `json:"champions"`
	SetItems  []string   `json:"items"`
	Mutator   string     `json:"mutator"`
	Name      string     `json:"name"`
	Number    int        `json:"number"`
	Traits    []Trait    `json:"traits"`
	SetAugments  []string   `json:"augments"`
}

type Champion struct {
	Ability       Ability  `json:"ability"`
	ApiName       string   `json:"apiName"`
	CharacterName string   `json:"characterName"`
	Cost          int      `json:"cost"`
	Icon          string   `json:"icon"`
	Name          string   `json:"name"`
	Role          string   `json:"role"`
	SquareIcon    string   `json:"squareIcon"`
	Stats         Stats    `json:"stats"`
	TileIcon      string   `json:"tileIcon"`
	Traits        []string `json:"traits"`
}

type Ability struct {
	Desc      string            `json:"desc"`
	Icon      string            `json:"icon"`
	Name      string            `json:"name"`
	Variables []AbilityVariable `json:"variables"`
}

type AbilityVariable struct {
	Name  string    `json:"name"`
	Value []float64 `json:"value"`
}

type Stats struct {
	Armor          float64 `json:"armor"`
	AttackSpeed    float64 `json:"attackSpeed"`
	CritChance     float64 `json:"critChance"`
	CritMultiplier float64 `json:"critMultiplier"`
	Damage         float64 `json:"damage"`
	HP             float64 `json:"hp"`
	InitialMana    float64 `json:"initialMana"`
	MagicResist    float64 `json:"magicResist"`
	Mana           float64 `json:"mana"`
	Range          float64 `json:"range"`
}

type Trait struct {
	ApiName string   `json:"apiName"`
	Desc    string   `json:"desc"`
	Effects []Effect `json:"effects"`
	Icon    string   `json:"icon"`
	Name    string   `json:"name"`
}

type Effect struct {
	MaxUnits  int                `json:"maxUnits"`
	MinUnits  int                `json:"minUnits"`
	Style     int                `json:"style"`
	Variables map[string]float64 `json:"variables"` // Using map to handle dynamic keys
}

// Item represents a single item (or augment, based on JSON structure) in TFT data.
type Item struct {
	ApiName            string             `json:"apiName"`
	AssociatedTraits   []string           `json:"associatedTraits"`
	Composition        []string           `json:"composition"` // Often empty for augments
	Desc               string             `json:"desc"`
	Effects            map[string]float64 `json:"effects"`
	From               []string           `json:"from"` // Can be null, list of components
	Icon               string             `json:"icon"`
	ID                 *string            `json:"id"` // Use pointer to handle potential null
	IncompatibleTraits []string           `json:"incompatibleTraits"`
	Name               string             `json:"name"`
	Tags               []string           `json:"tags"`
	Unique             bool               `json:"unique"`
}