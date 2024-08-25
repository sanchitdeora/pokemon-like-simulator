package data

type InBattlePokemon struct {
	Pokemon   *Pokemon
	BattleHP  int
	IsFainted bool
}

type BattleReport struct {
	UserWin      bool
	PrizeMoney   int
	BonusItems   []*Item
	BadgesEarned *BadgeType
}

type BattleOpts struct {
	Type string `json:"type"`
}

type BattleInput struct {
	Type           BattleInputType
	CurrentPokemon *InBattlePokemon
	Target         *InBattlePokemon
	Move           *Moves
	Item           *Item
	IsUser         bool
}

type BattleInputType string

const (
	Attack BattleInputType = "attack"
	Switch BattleInputType = "switch"
	Bag    BattleInputType = "bag"
	Run    BattleInputType = "run"
)

func CreateNewInBattlePokemon(pokemon *Pokemon) *InBattlePokemon {
	return &InBattlePokemon{
		Pokemon:   pokemon,
		BattleHP:  pokemon.Stats.HP,
		IsFainted: false,
	}
}
