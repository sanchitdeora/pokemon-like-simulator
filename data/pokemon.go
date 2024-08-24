package data

type Pokemon struct {
	BasePokemon    BasePokemon  `json:"base_pokemon"`
	Stats          PokemonStats `json:"stats"`
	Level          int          `json:"level"`
	ExperienceLeft int          `json:"experience_left"`
	Moveset        Moveset      `json:"moveset"`
	PartyOrder     int          `json:"party_order"`
}

// TODO add level at which pokemon learns
type BasePokemon struct {
	ID                  int             `json:"id"`
	Name                string          `json:"name"`
	BaseExperience      int             `json:"base_experience"`
	LearnByLevelUpMoves []Moves         `json:"moves"`
	SpritesURL          string          `json:"sprites"`
	BaseStats           PokemonStats    `json:"stats"`
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
