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
	BasePokemonURL string       `json:"base_pokemon_url"`
	Stats          PokemonStats `json:"stats"`
	Level          int          `json:"level"`
	ExperienceLeft int          `json:"experience_left"`
	Moveset        Moveset      `json:"moveset"`
}

type BasePokemon struct {
	ID                  int             `json:"id"`
	Name                string          `json:"name"`
	BaseExperience      int             `json:"base_experience"`
	MovesLearnedByLevel map[int]Moves   `json:"moves_learned_by_level"`
	SpritesURL          string          `json:"sprites"`
	BaseStats           PokemonStats    `json:"base_stats"`
	Type1               PokemonTypeName `json:"type1"`
	Type2               PokemonTypeName `json:"type2,omitempty"`
}

type PokemonStats struct {
	Speed          int `json:"speed"`
	Attack         int `json:"attack"`
	Defense        int `json:"defense"`
	SpecialAttack  int `json:"special_attack"`
	SpecialDefense int `json:"special_defense"`
	HP             int `json:"hp"`
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
