package battle

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/pokemon"
	"github.com/sanchitdeora/PokeSim/usermanagement"
)

type TrainerBattleOpts struct {
	UserService usermanagement.User
	PokemonService pokemon.Service
}

type TrainerBattle struct {
	// User
	User                    *data.User
	UserActivePokemon       *data.InBattlePokemon
	UserInBattleParty       []*data.InBattlePokemon
	UserUnfaintedPartyCount int

	// Trainer
	Trainer                    *data.Trainer
	TrainerActivePokemon       *data.InBattlePokemon
	TrainerInBattleParty       []*data.InBattlePokemon
	TrainerPokemonFacedExp     map[*data.Pokemon][]*data.InBattlePokemon
	TrainerUnfaintedPartyCount int
}

type TrainerBattleImpl struct {
	*TrainerBattleOpts
	UserService usermanagement.User
	PokemonService pokemon.Service
	*TrainerBattle
}

func NewTrainerBattle(opts *TrainerBattleOpts, trainer *data.Trainer) BattleIFace {
	// prepare user
	var userActivePokemon *data.InBattlePokemon
	userInBattlePokemonParty := make([]*data.InBattlePokemon, 0, 6)

	// create trainer
	var trainerActivePokemon *data.InBattlePokemon
	trainerInBattlePokemonParty := make([]*data.InBattlePokemon, 0, 6)

	trainerPokemonFacedExp := make(map[*data.Pokemon][]*data.InBattlePokemon, 0)
	for partyIndex, pokemon := range trainer.Party {
		if pokemon == nil {
			continue
		}
		inBattlePokemon := data.CreateNewInBattlePokemon(pokemon)
		if partyIndex == 0 {
			trainerActivePokemon = inBattlePokemon
		} else if len(trainerInBattlePokemonParty) < 5 {
			trainerInBattlePokemonParty = append(trainerInBattlePokemonParty, inBattlePokemon)
		}
		trainerPokemonFacedExp[pokemon] = make([]*data.InBattlePokemon, 0, 6)
	}

	// Get User
	user, err := opts.UserService.LoadUser()
	if err != nil {
		slog.Warn("did not get user", "user", user, "error", err)
		user = &data.User{}
	}

	for partyIndex, pokemon := range user.Party {
		if pokemon == nil {
			continue
		}
		inBattlePokemon := data.CreateNewInBattlePokemon(pokemon)
		if partyIndex == 0 {
			userActivePokemon = inBattlePokemon
		} else if len(userInBattlePokemonParty) <= 6 {
			userInBattlePokemonParty = append(userInBattlePokemonParty, inBattlePokemon)
		}
	}

	return &TrainerBattleImpl{
		TrainerBattleOpts: opts,
		TrainerBattle: &TrainerBattle{
			User:                    user,
			UserActivePokemon:       userActivePokemon,
			UserInBattleParty:       userInBattlePokemonParty,
			UserUnfaintedPartyCount: len(userInBattlePokemonParty) + 1, // +1 for active pokemon

			Trainer:                    trainer,
			TrainerActivePokemon:       trainerActivePokemon,
			TrainerInBattleParty:       trainerInBattlePokemonParty,
			TrainerPokemonFacedExp:     trainerPokemonFacedExp,
			TrainerUnfaintedPartyCount: len(trainerInBattlePokemonParty) + 1, // +1 for active pokemon
		},
	}
}

func (tb *TrainerBattleImpl) InitiateBattleSequence() (*data.BattleReport, error) {
	slog.Info(fmt.Sprintf("%s chooses %s!", tb.getTrainerName(false), tb.TrainerActivePokemon.Pokemon.Name))
	slog.Info(fmt.Sprintf("%s, I choose you!\n", tb.getActivePokemonName(true)))

	battleCompleteFlag := false
	tb.AddToTrainerPokemonFacedExp(tb.UserActivePokemon)

	for {
		slog.Info(fmt.Sprintf("%s Health: %v", tb.getActivePokemonName(true), tb.UserActivePokemon.BattleHP))
		slog.Info(fmt.Sprintf("%s Health: %v\n", tb.getActivePokemonName(false), tb.TrainerActivePokemon.BattleHP))

		battleInputs := getPokemonAttackOrder(tb.UserActivePokemon, tb.TrainerActivePokemon)

		for _, input := range battleInputs {
			tb.Turn(input)
			if battleCompleteFlag = tb.IsBattleOver(); battleCompleteFlag {
				slog.Info("Battle is completed!")
				break
			}
		}
		if battleCompleteFlag {
			break
		}
	}
	report, err := tb.BattleReport()
	if err != nil {
		slog.Error("error getting battle report", "error", err)
	}

	// Update User and save after battle completed
	tb.UserService.PostBattleUpdate(tb.User, report)

	// Pokemon Evolutions
	for _, inBattlePokemon := range tb.UserInBattleParty {
		if inBattlePokemon.CanEvolve {
			tb.PokemonService.Evolve(inBattlePokemon.Pokemon)
		}
	}

	return report, err 
}

func (tb *TrainerBattleImpl) Turn(userInput *data.BattleInput) {
	switch userInput.Type {
	case data.Switch:
		slog.Info(fmt.Sprintf("%s is switching %s for %s", tb.getTrainerName(userInput.IsUser), userInput.CurrentPokemon.Pokemon.Name, userInput.Target.Pokemon.Name))
		tb.SwitchPokemon(userInput.Target.Pokemon, userInput.IsUser)

	case data.Bag:
		if userInput.Item != nil && userInput.Item.Category == data.MedicalItems {
			tb.UseItem(userInput.Target, userInput.Item)
		}

	case data.Attack:
		tb.Attack(userInput.CurrentPokemon, userInput.Target, userInput.Move, userInput.IsUser)

		if tb.UserActivePokemon.IsFainted {
			slog.Info(fmt.Sprintf("%s has fainted!", tb.getActivePokemonName(true)))
			tb.UserUnfaintedPartyCount -= 1

			if tb.UserUnfaintedPartyCount > 0 {
				// switch active pokemon with first in party and push fainted pokemon at end
				nextUnfaintedPokemonIndex := 0
				for i, pokemon := range tb.TrainerInBattleParty {
					if !pokemon.IsFainted {
						nextUnfaintedPokemonIndex = i
						break
					}
				}
				switchPokemonWithIndex(nextUnfaintedPokemonIndex, tb.UserActivePokemon, &tb.UserInBattleParty)
				// tb.SwitchPokemonWithIndex(nextUnfaintedPokemonIndex, true)
			}
		}

		if tb.TrainerActivePokemon.IsFainted {
			slog.Info(fmt.Sprintf("%s has fainted!", tb.getActivePokemonName(false)))
			tb.TrainerUnfaintedPartyCount -= 1

			tb.BattleExperienceGain(tb.TrainerActivePokemon.Pokemon, tb.TrainerPokemonFacedExp[tb.TrainerActivePokemon.Pokemon])

			if tb.TrainerUnfaintedPartyCount > 0 {
				// switch active pokemon with first in party and push fainted pokemon at end
				nextUnfaintedPokemonIndex := 0
				for i, pokemon := range tb.TrainerInBattleParty {
					if !pokemon.IsFainted {
						nextUnfaintedPokemonIndex = i
						break
					}
				}
				switchPokemonWithIndex(nextUnfaintedPokemonIndex, tb.TrainerActivePokemon, &tb.TrainerInBattleParty)
				// tb.SwitchPokemonWithIndex(nextUnfaintedPokemonIndex, false)
			}
		}

	// TODO: turn should not be counted
	case data.Run:
		tb.Run()
	}
}

func (tb *TrainerBattleImpl) Attack(attackPokemon *data.InBattlePokemon, targetPokemon *data.InBattlePokemon, attackMove *data.Moves, isUser bool) {
	slog.Info(fmt.Sprintf("%s used %s", tb.getActivePokemonName(isUser), attackMove.Name))

	damagePoints := calculateAttackDamage(attackPokemon, targetPokemon, attackMove, 1.0)

	if damagePoints >= targetPokemon.BattleHP {
		targetPokemon.BattleHP = 0
		targetPokemon.IsFainted = true
	} else {
		targetPokemon.BattleHP -= damagePoints
	}

	slog.Info(fmt.Sprintf("%s did %v points of damage to %s", tb.getActivePokemonName(isUser), damagePoints, tb.getActivePokemonName(!isUser)))
}

func (tb *TrainerBattleImpl) SwitchPokemon(switchingPokemon *data.Pokemon, isUser bool) {
	var activePokemon *data.InBattlePokemon
	var inBattleParty []*data.InBattlePokemon

	if isUser {
		activePokemon = tb.UserActivePokemon
		inBattleParty = tb.UserInBattleParty
	} else {
		activePokemon = tb.TrainerActivePokemon
		inBattleParty = tb.TrainerInBattleParty
	}

	for i, pokemon := range inBattleParty {
		if pokemon.Pokemon == switchingPokemon && !pokemon.IsFainted {
			switchPokemonWithIndex(i, activePokemon, &inBattleParty)
			break
		}
	}

	slog.Info(fmt.Sprintf("%s is switching their pokemon!", tb.getTrainerName(isUser)))

	if isUser {
		slog.Info(fmt.Sprintf("%s, I choose you!", tb.getActivePokemonName(true)))
	} else {
		slog.Info(fmt.Sprintf("%s, chooses %s!", tb.getTrainerName(false), tb.TrainerActivePokemon.Pokemon.Name))
	}

	// update map of pokemon facing
	tb.AddToTrainerPokemonFacedExp(tb.UserActivePokemon)
}

func (tb *TrainerBattleImpl) UseItem(targetPokemon *data.InBattlePokemon, item *data.Item) {
	switch item.Category {
	case data.MedicalItems:
		healPokemon(targetPokemon, item)
	case data.PokeBalls:
		tb.CatchPokemon(targetPokemon, item)
	}
}

// TODO: turn should not be counted
func (tb *TrainerBattleImpl) CatchPokemon(targetPokemon *data.InBattlePokemon, item *data.Item) {
	slog.Info("Cannot catch a pokemon in trainer battle")
}

func (tb *TrainerBattleImpl) Run() {
	slog.Info("No! There's no running from a trainer battle!")
}

func (tb *TrainerBattleImpl) IsBattleOver() bool {
	if tb.UserUnfaintedPartyCount == 0 || tb.TrainerUnfaintedPartyCount == 0 {
		return true
	}
	return false
}

func (tb *TrainerBattleImpl) BattleReport() (*data.BattleReport, error) {
	if !tb.IsBattleOver() {
		return nil, errors.New("battle is ongoing")
	}

	var report data.BattleReport
	if tb.UserUnfaintedPartyCount == 0 {
		report.Money = data.GetMoneyLost(tb.User)

		slog.Info(fmt.Sprintf("You lost $%v!", report.Money))
		slog.Info(fmt.Sprintf("You lost the battle to %s!", tb.Trainer.Name))
		report.UserWin = false

	} else {
		report.UserWin = true
		report.Money = data.GetPrizeMoney(tb.Trainer)

		slog.Info(fmt.Sprintf("%s has won the battle!", tb.User.Name))
		slog.Info(fmt.Sprintf("You got $%v!", report.Money))

		// if gym battle; earn badge
		if tb.Trainer.Type == data.GymLeaderPrefix {
			report.BadgeEarned = &tb.Trainer.Rewards.Badge
			slog.Info(fmt.Sprintf("You earned a $%v!", tb.Trainer.Rewards.Badge.Name))
		}
	}

	return &report, nil
}

func (tb *TrainerBattleImpl) BattleExperienceGain(faintedPokemon *data.Pokemon, pokemonFaced []*data.InBattlePokemon) {
	for _, inBattlePokemon := range pokemonFaced {
		if !inBattlePokemon.IsFainted {
			expGain := calculateExperienceGained(faintedPokemon.Level, faintedPokemon.BasePokemon.BaseExperience, inBattlePokemon.Pokemon.Level)
			slog.Info(fmt.Sprintf("%s gained %v experience points", inBattlePokemon.Pokemon.Name, expGain))
			
			inBattlePokemon.CanEvolve = tb.PokemonService.ExperienceGain(expGain, inBattlePokemon.Pokemon)
		}
	}
}

func (tb *TrainerBattleImpl) AddToTrainerPokemonFacedExp(pokemon *data.InBattlePokemon) {
	tb.TrainerPokemonFacedExp[tb.TrainerActivePokemon.Pokemon] =
		append(tb.TrainerPokemonFacedExp[tb.TrainerActivePokemon.Pokemon], pokemon)
}

func (tb *TrainerBattleImpl) getTrainerName(isUser bool) string {
	if isUser {
		return tb.User.Name
	} else {
		return fmt.Sprintf("%s %s", tb.Trainer.Type, tb.Trainer.Name)
	}
}

func (tb *TrainerBattleImpl) getActivePokemonName(isUser bool) string {
	if isUser {
		return tb.UserActivePokemon.Pokemon.Name
	} else {
		return fmt.Sprintf("%s %s", "the opposing", tb.TrainerActivePokemon.Pokemon.Name)
	}
}
