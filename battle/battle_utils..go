package battle

import (
	"math"
	"time"
	"math/rand"

	"github.com/sanchitdeora/PokeSim/data"
)

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

// calculate experience gained
func calculateExperienceGained(faintedPokemonLevel int, faintedPokemonBaseExp int, userPokemonLevel int) int {
	var totalExp float32
	totalExp = float32(faintedPokemonBaseExp*faintedPokemonLevel) / 5.0
	totalExp *= float32(math.Pow((float64((2*faintedPokemonLevel)+10)/float64(faintedPokemonLevel+userPokemonLevel+10)), 2.5)) + 1.0

	// TODO: determine if userPokemonLevel is on/past the next evolution level. If, yes -> x1.2

	return int(math.Round(float64(totalExp)))
}

func isCritHit() bool {
	return randomGenerator(0, 1) < (1.0 / 24.0)
}

func attackCoefficient() float32 {
	return randomGenerator(0.85, 1)
}

func isStab(attackMove *data.Moves, attackPokemon *data.InBattlePokemon) bool {
	return attackMove.Type == attackPokemon.Pokemon.BasePokemon.Type1 || attackMove.Type == attackPokemon.Pokemon.BasePokemon.Type2
}

func randomGenerator(min float32, max float32) float32 {
	randIndex := rand.Float32()
	return (min + randIndex*(max-min))
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

	input := randomGenerator(0, float32(len(moves)))
	return moves[int(input)]
}