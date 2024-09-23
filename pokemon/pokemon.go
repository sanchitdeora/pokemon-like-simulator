package pokemon

import (
	"fmt"
	"log/slog"

	"github.com/sanchitdeora/PokeSim/data"
)

type Service interface {
	LevelUp(pokemon *data.Pokemon)
	Evolve(pokemon *data.Pokemon)
	LearnNewMoves(move *data.Moves)
	ExperienceGain(expGain int, pokemon *data.Pokemon) bool
}

type PokemonOpts struct{}

type PokemonImpl struct {
	opts PokemonOpts
}

func NewPokemon(opts PokemonOpts) Service {
	return &PokemonImpl{
		opts: opts,
	}
}

func (p *PokemonImpl) LevelUp(pokemon *data.Pokemon) {
	pokemon.Level++

	// calculate stats upgrade
	statUpgrades(pokemon)

	// calculate pokemon experience left
	pokemon.ExperienceLeft = getExperienceRequiredForNextLevel(pokemon)

	// TODO: get user input if pokemon learns moves; if yes which moveset?
	// if move, exists := pokemon.MovesLearnedByLevel[pokemon.Level]; exists {
	// 	p.LearnNewMoves(&move)
	// }
}

// // TODO: add should evolve method. Evolve after battle.
func (p *PokemonImpl) Evolve(pokemon *data.Pokemon) {
	if canPokemonEvolve(pokemon) {
		evolvedBasePokemon := pokemon.EvolutionChain[pokemon.Level]
		pokemon.BasePokemon = evolvedBasePokemon[0]
		statUpgrades(pokemon)
	}
}

func (p *PokemonImpl) LearnNewMoves(move *data.Moves) {
	// input which move should be replaced.
}

// TODO: add EVs to pokemon with target faints EV Yield
func (p *PokemonImpl) PostTargetPokemonFaint(expGain int, pokemon *data.Pokemon, evYield *data.BasePokemonStats) {
	// update EVs

	// Experience
}

func (p *PokemonImpl) ExperienceGain(expGain int, pokemon *data.Pokemon) bool {
	var canEvolve bool

	for {
		if pokemon.ExperienceLeft > expGain {
			pokemon.ExperienceLeft -= expGain
			break
		}
		expGain -= pokemon.ExperienceLeft
		p.LevelUp(pokemon)

		if canPokemonEvolve(pokemon) {
			canEvolve = true
		}
	}
	return canEvolve
}

func statUpgrades(pokemon *data.Pokemon) {
	hp := calculateHPStatUpgrade(pokemon.BaseStats.HP, &pokemon.Stats.HP, pokemon.Level)

	attack := calculateOtherStatUpgrade(pokemon.BaseStats.Attack, &pokemon.Stats.Attack, pokemon.Level)
	defense := calculateOtherStatUpgrade(pokemon.BaseStats.Defense, &pokemon.Stats.Defense, pokemon.Level)

	speed := calculateOtherStatUpgrade(pokemon.BaseStats.Speed, &pokemon.Stats.Speed, pokemon.Level)

	spAttack := calculateOtherStatUpgrade(pokemon.BaseStats.SpecialAttack, &pokemon.Stats.SpecialAttack, pokemon.Level)
	spDefense := calculateOtherStatUpgrade(pokemon.BaseStats.SpecialDefense, &pokemon.Stats.SpecialDefense, pokemon.Level)

	slog.Info("Stat upgrade for pokemon:")
	slog.Info(fmt.Sprintf("HP: +%v", hp-pokemon.Stats.HP.Value))
	slog.Info(fmt.Sprintf("Attack: +%v", attack-pokemon.Stats.Attack.Value))
	slog.Info(fmt.Sprintf("Defense: +%v", defense-pokemon.Stats.Defense.Value))
	slog.Info(fmt.Sprintf("Special Attack: +%v", spAttack-pokemon.Stats.SpecialAttack.Value))
	slog.Info(fmt.Sprintf("Special Defence: +%v", spDefense-pokemon.Stats.SpecialDefense.Value))
	slog.Info(fmt.Sprintf("Speed: +%v", speed-pokemon.Stats.Speed.Value))

	pokemon.Stats.HP.Value = hp
	pokemon.Stats.Attack.Value = attack
	pokemon.Stats.Defense.Value = defense
	pokemon.Stats.SpecialAttack.Value = spAttack
	pokemon.Stats.SpecialDefense.Value = spDefense
	pokemon.Stats.Speed.Value = speed
}

// calculate next level exp required
func getExperienceRequiredForNextLevel(pokemon *data.Pokemon) int {
	switch pokemon.GrowthRate {
	case data.Erratic:
		return nextLevelErraticExp(pokemon.Level)
	case data.Fast:
		return nextLevelFastExp(pokemon.Level)
	case data.MediumFast:
		return nextLevelMediumFastExp(pokemon.Level)
	case data.MediumSlow:
		return nextLevelMediumSlowExp(pokemon.Level)
	case data.Slow:
		return nextLevelSlowExp(pokemon.Level)
	case data.Fluctuating:
		return nextLevelFluctuatingExp(pokemon.Level)
	default:
		slog.Error("invalid growth rate type found", "growth rate type", pokemon.GrowthRate)
	}

	return -1
}

func canPokemonEvolve(pokemon *data.Pokemon) bool {
	_, exists := pokemon.EvolutionChain[pokemon.Level]
	return exists
}
