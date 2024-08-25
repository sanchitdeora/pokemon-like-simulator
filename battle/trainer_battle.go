package battle

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/sanchitdeora/PokeSim/data"
)

type TrainerBattleOpts struct{}

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
	*TrainerBattle
}

func NewTrainerBattle(opts *TrainerBattleOpts, trainer *data.Trainer, user *data.User) BattleIFace {
	// prepare user
	var userActivePokemon *data.InBattlePokemon
	userInBattlePokemonParty := make([]*data.InBattlePokemon, 0, 6)

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

		battleInputs := tb.GetPokemonAttackOrder()

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
	return tb.BattleReport()
}

func (tb *TrainerBattleImpl) GetPokemonAttackOrder() (inputs []*data.BattleInput) {
	userInput := waitForInput(tb.UserActivePokemon, tb.TrainerActivePokemon, true)
	trainerInput := waitForInput(tb.TrainerActivePokemon, tb.UserActivePokemon, false)

	// if any pokemon's move is higher priority, it will go first. If equal, check speed
	if userInput.Move.Priority != trainerInput.Move.Priority {
		if userInput.Move.Priority > trainerInput.Move.Priority {
			return append(inputs, userInput, trainerInput)
		} else {
			return append(inputs, trainerInput, userInput)
		}
	}

	// if user's active pokemon speed >= opposing pokemon's, then user goes first.
	if tb.UserActivePokemon.Pokemon.Stats.Speed >= tb.TrainerActivePokemon.Pokemon.Stats.Speed {
		return append(inputs, userInput, trainerInput)
	} else {
		return append(inputs, trainerInput, userInput)
	}
}

func (tb *TrainerBattleImpl) Turn(userInput *data.BattleInput) {
	switch userInput.Type {
	case data.Switch:
		slog.Info(fmt.Sprintf("%s is switching %s for %s", tb.getTrainerName(userInput.IsUser), userInput.CurrentPokemon.Pokemon.Name, userInput.Target.Pokemon.Name))
		tb.SwitchPokemon(userInput.Target.Pokemon, userInput.IsUser)

	case data.Bag:
		if userInput.Item != nil && userInput.Item.ItemType == data.MedicalItems {
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
				tb.SwitchPokemonWithIndex(nextUnfaintedPokemonIndex, true)
			}
		}

		if tb.TrainerActivePokemon.IsFainted {
			slog.Info(fmt.Sprintf("%s has fainted!", tb.getActivePokemonName(false)))
			tb.TrainerUnfaintedPartyCount -= 1

			battleExperienceGain(tb.TrainerActivePokemon.Pokemon, tb.TrainerPokemonFacedExp[tb.TrainerActivePokemon.Pokemon])

			if tb.TrainerUnfaintedPartyCount > 0 {
				// switch active pokemon with first in party and push fainted pokemon at end
				nextUnfaintedPokemonIndex := 0
				for i, pokemon := range tb.TrainerInBattleParty {
					if !pokemon.IsFainted {
						nextUnfaintedPokemonIndex = i
						break
					}
				}
				tb.SwitchPokemonWithIndex(nextUnfaintedPokemonIndex, false)
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
	fmt.Println()
}

func (tb *TrainerBattleImpl) SwitchPokemon(switchingPokemon *data.Pokemon, isUser bool) {
	var inBattleParty []*data.InBattlePokemon
	if isUser {
		inBattleParty = tb.UserInBattleParty
	} else {
		inBattleParty = tb.TrainerInBattleParty
	}

	for i, pokemon := range inBattleParty {
		if pokemon.Pokemon == switchingPokemon {
			tb.SwitchPokemonWithIndex(i, isUser)
		}
	}
}

func (tb *TrainerBattleImpl) SwitchPokemonWithIndex(switchingPokemonIndex int, isUser bool) {
	slog.Info(fmt.Sprintf("%s is switching their pokemon!", tb.getTrainerName(isUser)))
	if isUser {
		currentPokemon := tb.UserActivePokemon
		nextPokemon := tb.UserInBattleParty[switchingPokemonIndex]
		tb.UserActivePokemon = nil

		// only get active pokemon if unfainted pokemon available; else add to the list
		if !nextPokemon.IsFainted {
			tb.UserActivePokemon = tb.UserInBattleParty[switchingPokemonIndex]
			tb.UserInBattleParty = append(tb.UserInBattleParty[:switchingPokemonIndex], tb.UserInBattleParty[switchingPokemonIndex+1:]...)
		}
		tb.UserInBattleParty = append(tb.UserInBattleParty, currentPokemon)

		slog.Info(fmt.Sprintf("%s, I choose you!", tb.getActivePokemonName(true)))

	} else {
		currentPokemon := tb.TrainerActivePokemon
		nextPokemon := tb.TrainerInBattleParty[switchingPokemonIndex]
		tb.TrainerActivePokemon = nil

		// only get active pokemon if unfainted pokemon available; else add to the list
		if !nextPokemon.IsFainted {
			tb.TrainerActivePokemon = tb.TrainerInBattleParty[switchingPokemonIndex]
			tb.TrainerInBattleParty = append(tb.TrainerInBattleParty[:switchingPokemonIndex], tb.TrainerInBattleParty[switchingPokemonIndex+1:]...)
		}
		tb.TrainerInBattleParty = append(tb.TrainerInBattleParty, currentPokemon)

		slog.Info(fmt.Sprintf("%s, chooses %s!", tb.getTrainerName(false), tb.TrainerActivePokemon.Pokemon.Name))

	}

	// update map of pokemon facing
	tb.AddToTrainerPokemonFacedExp(tb.UserActivePokemon)
}

func (tb *TrainerBattleImpl) UseItem(targetPokemon *data.InBattlePokemon, item *data.Item) {
	switch item.ItemType {
	case data.MedicalItems:
		tb.HealPokemon(targetPokemon, item)
	case data.PokeBalls:
		tb.CatchPokemon(targetPokemon, item)
	}
}

func (tb *TrainerBattleImpl) HealPokemon(targetPokemon *data.InBattlePokemon, item *data.Item) {
	targetPokemon.BattleHP += targetPokemon.BattleHP + item.Attributes

	if targetPokemon.BattleHP > targetPokemon.Pokemon.Stats.HP {
		targetPokemon.BattleHP = targetPokemon.Pokemon.Stats.HP
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

	var result data.BattleReport
	if tb.UserUnfaintedPartyCount == 0 {
		result.PrizeMoney = data.GetMoneyLost(tb.User)

		slog.Info(fmt.Sprintf("You lost $%v!", result.PrizeMoney))
		slog.Info(fmt.Sprintf("You lost the battle to %s!", tb.Trainer.Name))
		result.UserWin = false
	} else {
		result.UserWin = true
		result.PrizeMoney = data.GetPrizeMoney(tb.Trainer)

		slog.Info(fmt.Sprintf("%s has won the battle!", tb.User.Name))
		slog.Info(fmt.Sprintf("You got $%v!", result.PrizeMoney))

		// if gym battle; earn badge
		if tb.Trainer.Type == data.GymLeaderPrefix {
			result.BadgesEarned = &tb.Trainer.Rewards.Badge
			slog.Info(fmt.Sprintf("You earned a $%v!", tb.Trainer.Rewards.Badge.Name))
		}
	}

	return &result, nil
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
