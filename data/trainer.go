package data

type User struct {
	Trainer Trainer
	Stats   TrainerStats
}

type Trainer struct {
	Name           string          `json:"name"`
	Party          []*Pokemon      `json:"party"`
	Bag            []*Item         `json:"bag"`
	Type           TrainerType     `json:"type"`
	AdditionalInfo *AdditionalInfo `json:"additional_info"`
}

type AdditionalInfo struct {
	PrizeMoney       int       `json:"prize_money"`
	Badge            BadgeType `json:"badge"`
	ExperienceGained int       `json:"experience_given"`
}

type BadgeType struct {
	Name   string `json:"name"`
	Region string `json:"region"`
}

type TrainerType string

const (
	UserTrainer       TrainerType = "user"
	NormalTrainer     TrainerType = "trainer"
	GymLeaderTrainer  TrainerType = "leader"
	TournamentTrainer TrainerType = "tournament"
)

type TrainerStats struct {
	Fights  int `json:"fights"`
	Wins    int `json:"wins"`
	Catches int `json:"catches"`
	Losses  int `json:"losses"`
	PokeDEX int `json:"pokedex"`
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
