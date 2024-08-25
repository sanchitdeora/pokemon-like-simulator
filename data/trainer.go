package data

type User struct {
	BaseTrainer
	Stats *TrainerStats
}

type BaseTrainer struct {
	Name  string      `json:"name"`
	Party [6]*Pokemon `json:"party"`
	Bag   []*Item     `json:"bag"`
}

type Trainer struct {
	BaseTrainer
	Type    TrainerClass
	Rewards *Rewards
}

type Rewards struct {
	Items []*Item   `json:"items"`
	Badge BadgeType `json:"badge_type,omitempty"`
}

type BadgeType struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

type TrainerClass string

const (
	TrainerPrefix    TrainerClass = "Trainer"
	GymLeaderPrefix  TrainerClass = "Gym Leader"
	TournamentPrefix TrainerClass = "Tournament Trainer"
)

type TrainerStats struct {
	Badges  []BadgeType `json:"badges,omitempty"`
	Fights  int         `json:"fights"`
	Wins    int         `json:"wins"`
	Catches int         `json:"catches"`
	Losses  int         `json:"losses"`
	PokeDEX int         `json:"pokedex"`
}

type Item struct {
	Count       int          `json:"count"`
	ItemType    ItemCategory `json:"item_category"`
	Cost        int          `json:"cost"`
	Attributes  int          `json:"attributes"`
	Description string       `json:"description"`
}

type ItemCategory string

const (
	MedicalItems ItemCategory = "MedicalItems"
	PokeBalls    ItemCategory = "PokeBalls"
)

var BasePayoutTable map[TrainerClass]int = map[TrainerClass]int{
	TrainerPrefix:    80,
	GymLeaderPrefix:  160,
	TournamentPrefix: 160,
	// TODO: Add more as needed
}

var BlackOutPayoutTable map[int]int = map[int]int{0: 8, 1: 16, 2: 24, 3: 36, 4: 48, 5: 64, 6: 80, 7: 100, 8: 120}

func GetPrizeMoney(trainer *Trainer) int {
	highestLevel := 0
	for _, pokemon := range trainer.Party {
		if pokemon == nil {
			continue
		}
		if pokemon.Level > highestLevel {
			highestLevel = pokemon.Level
		}
	}
	return BasePayoutTable[trainer.Type] * highestLevel
}

func GetMoneyLost(user *User) int {
	numBadges := 0
	if user.Stats != nil && user.Stats.Badges != nil {
		numBadges = len(user.Stats.Badges)
	}
	
	highestLevel := 0
	for _, pokemon := range user.Party {
		if pokemon == nil {
			continue
		}
		if pokemon.Level > highestLevel {
			highestLevel = pokemon.Level
		}
	}
	return BlackOutPayoutTable[numBadges] * highestLevel
}
