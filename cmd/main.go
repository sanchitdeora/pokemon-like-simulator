package main

import (
	"fmt"

	"github.com/sanchitdeora/PokeSim/battle"
	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/usermanagement"
	"github.com/sanchitdeora/PokeSim/utils"
)

const (
	SAVED_USER_PATH = "C:\\Projects\\Go-projects\\src\\PokéSim\\saved\\user.json"
)

func main() {
	// move, err := utils.ReadJsonFromFile[*data.Moves]("C:\\Projects\\Go-projects\\src\\PokéSim\\transformed_firePunchMove.json")
	// fmt.Println(move, err)

	userService := usermanagement.NewUserService(usermanagement.UserOpts{SavedUserPath: SAVED_USER_PATH})

	squirtle, err := utils.ReadJsonFromFile[*data.Pokemon]("C:\\Projects\\Go-projects\\src\\PokéSim\\testfiles\\transformed_squirtlePokemon.json")
	fmt.Println(squirtle, err)

	bulbasaur, err := utils.ReadJsonFromFile[*data.Pokemon]("C:\\Projects\\Go-projects\\src\\PokéSim\\testfiles\\transformed_bulbasaurPokemon.json")
	fmt.Println(bulbasaur, err)

	trainerBattle := battle.NewTrainerBattle(&battle.TrainerBattleOpts{
		UserService: userService,
	},
		&data.Trainer{
			BaseTrainer: data.BaseTrainer{
				Name:  "Bash kechtup",
				Party: [6]*data.Pokemon{bulbasaur},
			},
			Type:    data.TrainerPrefix,
			Rewards: &data.Rewards{},
		},
	)

	trainerBattle.InitiateBattleSequence()
}
