package main

import (
	"fmt"

	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/battle"
	"github.com/sanchitdeora/PokeSim/utils"
)

func main() {
	// move, err := utils.ReadJsonFromFile[*data.Moves]("C:\\Projects\\Go-projects\\src\\PokéSim\\transformed_firePunchMove.json")
	// fmt.Println(move, err)

	squirtle, err := utils.ReadJsonFromFile[*data.Pokemon]("C:\\Projects\\Go-projects\\src\\PokéSim\\testfiles\\transformed_squirtlePokemon.json")
	fmt.Println(squirtle, err)

	bulbasaur, err := utils.ReadJsonFromFile[*data.Pokemon]("C:\\Projects\\Go-projects\\src\\PokéSim\\testfiles\\transformed_bulbasaurPokemon.json")
	fmt.Println(bulbasaur, err)

	trainerBattle := battle.NewTrainerBattle(&data.Trainer{
		Name: "Bash kechtup",
		Type: data.NormalTrainer,
		Party: []*data.Pokemon{bulbasaur},
		AdditionalInfo: &data.AdditionalInfo{
			PrizeMoney: 100,
		},
	},
	&data.Trainer{
		Name: "John Cena",
		Type: data.NormalTrainer,
		Party: []*data.Pokemon{squirtle},
	})

	trainerBattle.InitiateBattleSequence()
}