package data

import "github.com/sanchitdeora/PokeSim/utils"

type PokemonSave struct {
	BasePokemonURL string       `json:"base_pokemon_url"`
	Stats          PokemonStats `json:"stats"`
	Level          int          `json:"level"`
	ExperienceLeft int          `json:"experience_left"`
	Moveset        Moveset      `json:"moveset"`
}

type Pokemon struct {
	BasePokemon
	BasePokemonURL string                `json:"base_pokemon_url"`
	EvolutionChain map[int][]BasePokemon `json:"evolution_chain"`
	Stats          PokemonStats          `json:"stats"`
	Level          int                   `json:"level"`
	ExperienceLeft int                   `json:"experience_left"`
	Moveset        Moveset               `json:"moveset"`
}

type BasePokemon struct {
	ID                  int              `json:"id"`
	Name                string           `json:"name"`
	BaseExperience      int              `json:"base_experience"`
	GrowthRate          GrowthRateTypes  `json:"growth_rate"`
	MovesLearnedByLevel map[int]Moves    `json:"moves_learned_by_level"`
	SpritesURL          string           `json:"sprites"`
	BaseStats           BasePokemonStats `json:"base_stats"`
	EVYield             BasePokemonStats `json:"ev_yield"`
	Type1               PokemonTypeName `json:"type1"`
	Type2               PokemonTypeName `json:"type2,omitempty"`
}

type BasePokemonStats struct {
	Speed          int `json:"speed"`
	Attack         int `json:"attack"`
	Defense        int `json:"defense"`
	SpecialAttack  int `json:"special_attack"`
	SpecialDefense int `json:"special_defense"`
	HP             int `json:"hp"`
}

type PokemonStats struct {
	Speed          PokemonStat `json:"speed"`
	Attack         PokemonStat `json:"attack"`
	Defense        PokemonStat `json:"defense"`
	SpecialAttack  PokemonStat `json:"special_attack"`
	SpecialDefense PokemonStat `json:"special_defense"`
	HP             PokemonStat `json:"hp"`
}

type PokemonStat struct {
	Value int `json:"value"`
	IV    int `json:"iv"`
	EV    int `json:"ev"`
}

type Moves struct {
	ID          int             `json:"id"`
	Name        string          `json:"name"`
	Accuracy    int             `json:"accuracy"`
	Priority    int             `json:"priority"`
	Power       int             `json:"power"`
	DamageClass MoveDamageClass `json:"damage_class"`
	Type        PokemonTypeName `json:"type"`
}

type Moveset struct {
	Move1 *Moves `json:"move1"`
	Move2 *Moves `json:"move2"`
	Move3 *Moves `json:"move3"`
	Move4 *Moves `json:"move4"`
}

type GrowthRateTypes string

const (
	Erratic     GrowthRateTypes = "erratic"
	Fast        GrowthRateTypes = "fast"
	MediumFast  GrowthRateTypes = "medium-fast"
	MediumSlow  GrowthRateTypes = "medium-slow"
	Slow        GrowthRateTypes = "slow"
	Fluctuating GrowthRateTypes = "fluctuating"
)

type MoveDamageClass string

const (
	Physical MoveDamageClass = "physical"
	Status   MoveDamageClass = "status"
	Special  MoveDamageClass = "special"
)

func (s *PokemonSave) ToPokemon() *Pokemon {
	basePokemon, _ := utils.ReadJsonFromFile[BasePokemon](s.BasePokemonURL)
	return &Pokemon{
		BasePokemon:    basePokemon,
		BasePokemonURL: s.BasePokemonURL,
		Stats:          s.Stats,
		Level:          s.Level,
		ExperienceLeft: s.ExperienceLeft,
		Moveset:        s.Moveset,
	}
}

func (p *Pokemon) ToPokemonSave() *PokemonSave {
	return &PokemonSave{
		BasePokemonURL: p.BasePokemonURL,
		Stats:          p.Stats,
		Level:          p.Level,
		ExperienceLeft: p.ExperienceLeft,
		Moveset:        p.Moveset,
	}
}
