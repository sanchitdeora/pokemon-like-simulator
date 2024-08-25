package main

import (
	"fmt"

	"github.com/sanchitdeora/PokeSim/battle"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/utils"
)

func main() {
	// move, err := utils.ReadJsonFromFile[*data.Moves]("C:\\Projects\\Go-projects\\src\\PokéSim\\transformed_firePunchMove.json")
	// fmt.Println(move, err)

	squirtle, err := utils.ReadJsonFromFile[*data.Pokemon]("C:\\Projects\\Go-projects\\src\\PokéSim\\testfiles\\transformed_squirtlePokemon.json")
	fmt.Println(squirtle, err)

	bulbasaur, err := utils.ReadJsonFromFile[*data.Pokemon]("C:\\Projects\\Go-projects\\src\\PokéSim\\testfiles\\transformed_bulbasaurPokemon.json")
	fmt.Println(bulbasaur, err)

	trainerBattle := battle.NewTrainerBattle(nil,
		&data.Trainer{
			BaseTrainer: data.BaseTrainer{
				Name:  "Bash kechtup",
				Party: [6]*data.Pokemon{bulbasaur},
			},
			Type:    data.TrainerPrefix,
			Rewards: &data.Rewards{},
		},
		&data.User{
			BaseTrainer: data.BaseTrainer{
				Name:  "John Cena",
				Party: [6]*data.Pokemon{squirtle},
			},
		},
	)

	trainerBattle.InitiateBattleSequence()
}
