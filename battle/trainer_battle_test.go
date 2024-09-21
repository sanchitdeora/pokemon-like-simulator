package battle_test

import (
	"testing"

	"github.com/sanchitdeora/PokeSim/battle"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/utils"
	"github.com/stretchr/testify/assert"
)

func createTrainerBattle() battle.BattleIFace {
	squirtle, _ := utils.ReadJsonFromFile[*data.Pokemon]("C:\\Projects\\Go-projects\\src\\PokéSim\\testfiles\\transformed_squirtlePokemon.json")
	bulbasaur, _ := utils.ReadJsonFromFile[*data.Pokemon]("C:\\Projects\\Go-projects\\src\\PokéSim\\testfiles\\transformed_bulbasaurPokemon.json")

	return battle.NewTrainerBattle(nil,
		&data.Trainer{
			BaseTrainer: data.BaseTrainer{
				Name:  "Bash kechtup",
				Party: [6]*data.Pokemon{squirtle, bulbasaur},
			},
			Type:    data.TrainerPrefix,
			Rewards: &data.Rewards{},
		},
	)
}

func TestBattle(t *testing.T) {
	battle := createTrainerBattle()

	report, err := battle.InitiateBattleSequence()
	assert.NoError(t, err)
	if report.UserWin {
		assert.Equal(t, 6000, report.Money)
	} else {
		assert.Equal(t, 600, report.Money)
	}
}
