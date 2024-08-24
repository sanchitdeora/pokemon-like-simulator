package battle_test

import (
	"testing"

	"github.com/sanchitdeora/PokeSim/battle"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/utils"
)

func createTrainerBattle() battle.BattleIFace {
	squirtle, _ := utils.ReadJsonFromFile[*data.Pokemon]("C:\\Projects\\Go-projects\\src\\PokéSim\\testfiles\\transformed_squirtlePokemon.json")
	bulbasaur, _ := utils.ReadJsonFromFile[*data.Pokemon]("C:\\Projects\\Go-projects\\src\\PokéSim\\testfiles\\transformed_bulbasaurPokemon.json")

	copyOfSquirtle := *squirtle
	copyOfSquirtle.PartyOrder = 2
	return battle.NewTrainerBattle(&data.Trainer{
		Name: "Bash kechtup",
		Type: data.NormalTrainer,
		Party: []*data.Pokemon{squirtle},
		AdditionalInfo: &data.AdditionalInfo{
			PrizeMoney: 100,
		},
	},
	&data.Trainer{
		Name: "John Cena",
		Type: data.NormalTrainer,
		Party: []*data.Pokemon{bulbasaur, &copyOfSquirtle},
	})
}

func TestBattle(t *testing.T) {
	battle := createTrainerBattle()

	battle.InitiateBattleSequence()
}