package battle

import (
	"math"
	"math/rand"
	"time"

	"github.com/sanchitdeora/PokeSim/data"
)

// calculate experience gained
func calculateExperienceGained(faintedPokemonLevel int, faintedPokemonBaseExp int, userPokemonLevel int) int {
	var totalExp float64
	totalExp = float64(faintedPokemonBaseExp*faintedPokemonLevel) / 5.0
	totalExp *= math.Pow((float64((2*faintedPokemonLevel)+10)/float64(faintedPokemonLevel+userPokemonLevel+10)), 2.5) + 1.0

	// TODO: determine if userPokemonLevel is on/past the next evolution level. If, yes -> x1.2

	return int(math.Round(totalExp))
}

func isCritHit() bool {
	return randomGenerator(0, 1) < (1.0 / 24.0)
}

func attackCoefficient() float64 {
	return randomGenerator(0.85, 1)
}

func isStab(attackMove *data.Moves, attackPokemon *data.InBattlePokemon) bool {
	return attackMove.Type == attackPokemon.Pokemon.BasePokemon.Type1 || attackMove.Type == attackPokemon.Pokemon.BasePokemon.Type2
}

func randomGenerator(min float64, max float64) float64 {
	randIndex := rand.Float64()
	return (min + randIndex*(max-min))
}

// TODO: move to a different package later
func waitForInput(pokemon *data.InBattlePokemon, target *data.InBattlePokemon, isUser bool) *data.BattleInput {
	time.Sleep(time.Second * 0)
	return &data.BattleInput{
		Type:           data.Attack,
		Move:           randomMove(&pokemon.Pokemon.Moveset),
		CurrentPokemon: pokemon,
		Target:         target,
		IsUser:         isUser,
	}
}

// generate random move
func randomMove(moveset *data.Moveset) *data.Moves {
	moves := []*data.Moves{}
	if moveset.Move1 != nil {
		moves = append(moves, moveset.Move1)
	}
	if moveset.Move2 != nil {
		moves = append(moves, moveset.Move2)
	}
	if moveset.Move3 != nil {
		moves = append(moves, moveset.Move3)
	}
	if moveset.Move4 != nil {
		moves = append(moves, moveset.Move4)
	}

	if len(moves) == 0 {
		return nil
	}

	input := randomGenerator(0, float64(len(moves)))
	return moves[int(input)]
}
