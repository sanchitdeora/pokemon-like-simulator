package battle

import (
	"errors"
	"fmt"
	"log/slog"

	"github.com/sanchitdeora/PokeSim/data"
)

type TrainerBattleOpts struct {
	TrainerLogPrefix string
}

type TrainerBattle struct {
	UserTrainerInfo  *BattleTrainerInfo
	EnemyTrainerInfo *BattleTrainerInfo
}

type TrainerBattleImpl struct {
	*TrainerBattleOpts
	*TrainerBattle
}

func NewTrainerBattle(enemyTrainer *data.Trainer, user *data.Trainer) BattleIFace {
	return &TrainerBattleImpl{
		TrainerBattleOpts: &TrainerBattleOpts{TrainerLogPrefix: "Trainer"},
		TrainerBattle: &TrainerBattle{
			UserTrainerInfo:  createTrainerInfo(user, true),
			EnemyTrainerInfo: createTrainerInfo(enemyTrainer, false),
		},
	}
}

func (tb *TrainerBattleImpl) InitiateBattleSequence() {
	slog.Info(fmt.Sprintf("%s, chooses %s!", tb.getTrainerName(false), tb.UserTrainerInfo.ActivePokemon.Pokemon.BasePokemon.Name))
	slog.Info(fmt.Sprintf("%s, I choose you!", tb.UserTrainerInfo.ActivePokemon.Pokemon.BasePokemon.Name))

	battleCompleteFlag := false
	tb.AddUserPokemonFacingEnemyActive(tb.UserTrainerInfo.ActivePokemon.Pokemon)

	for {
		slog.Info(fmt.Sprintf("%s Health: %v", tb.UserTrainerInfo.ActivePokemon.Pokemon.BasePokemon.Name, tb.UserTrainerInfo.ActivePokemon.BattleHP))
		slog.Info(fmt.Sprintf("Enemy %s Health: %v", tb.EnemyTrainerInfo.ActivePokemon.Pokemon.BasePokemon.Name, tb.EnemyTrainerInfo.ActivePokemon.BattleHP))

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
	tb.BattleReport()
}

func (tb *TrainerBattleImpl) GetPokemonAttackOrder() (inputs []*BattleInput) {
	userInput := waitForInput(tb.UserTrainerInfo.ActivePokemon, tb.EnemyTrainerInfo.ActivePokemon)
	enemyInput := waitForInput(tb.EnemyTrainerInfo.ActivePokemon, tb.UserTrainerInfo.ActivePokemon)

	// if any pokemon's move is higher priority, it will go first. If equal, check speed
	if userInput.Move.Priority != enemyInput.Move.Priority {
		if userInput.Move.Priority > enemyInput.Move.Priority {
			return append(inputs, userInput, enemyInput)
		} else {
			return append(inputs, enemyInput, userInput)
		}
	}

	// if user's active pokemon speed >= enemy's, then user goes first.
	if tb.UserTrainerInfo.ActivePokemon.Pokemon.Stats.Speed >= tb.EnemyTrainerInfo.ActivePokemon.Pokemon.Stats.Speed {
		return append(inputs, userInput, enemyInput)
	} else {
		return append(inputs, enemyInput, userInput)
	}
}

func (tb *TrainerBattleImpl) Turn(userInput *BattleInput) {
	switch userInput.Type {
	case Switch:
		slog.Info(fmt.Sprintf("%s is switching %s for %s", tb.getTrainerName(userInput.IsUser), userInput.CurrentPokemon.Pokemon.BasePokemon.Name, userInput.Target.Pokemon.BasePokemon.Name))
		tb.SwitchPokemon(userInput.Target.Pokemon.PartyOrder, userInput.IsUser, true)

	case Item:
		if userInput.Item != nil && userInput.Item.ItemType == data.MedicalItems {
			tb.UseItem(userInput.Target, userInput.Item)
		}

	case Attack:
		tb.Attack(userInput.CurrentPokemon, userInput.Target, userInput.Move)

		if tb.UserTrainerInfo.ActivePokemon.IsFainted {
			slog.Info(fmt.Sprintf("%s has fainted!", tb.UserTrainerInfo.ActivePokemon.Pokemon.BasePokemon.Name))
			tb.UserTrainerInfo.UnfaintedPartyCount -= 1

			if tb.UserTrainerInfo.UnfaintedPartyCount > 0 {
				// switch active pokemon with first in party and push fainted pokemon at end
				nextPokemonIndex := 0
				for i, pokemon := range tb.EnemyTrainerInfo.InBattleParty {
					if !pokemon.IsFainted {
						nextPokemonIndex = i
						break
					}
				}
				tb.SwitchPokemon(nextPokemonIndex, true, true)
			}
		}

		if tb.EnemyTrainerInfo.ActivePokemon.IsFainted {
			slog.Info(fmt.Sprintf("Enemy %s has fainted!", tb.EnemyTrainerInfo.ActivePokemon.Pokemon.BasePokemon.Name))
			tb.EnemyTrainerInfo.UnfaintedPartyCount -= 1

			BattleExperienceGain(tb.EnemyTrainerInfo.ActivePokemon.Pokemon, tb.EnemyTrainerInfo.EnemyPokemonFaced[tb.EnemyTrainerInfo.ActivePokemon.Pokemon])

			if tb.EnemyTrainerInfo.UnfaintedPartyCount > 0 {
				// switch active pokemon with first in party and push fainted pokemon at end
				nextPokemonIndex := 0
				for i, pokemon := range tb.EnemyTrainerInfo.InBattleParty {
					if !pokemon.IsFainted {
						nextPokemonIndex = i
						break
					}
				}
				tb.SwitchPokemon(nextPokemonIndex, false, true)
			}
		}

	case Run:
		tb.Run()
	}
}

func (tb *TrainerBattleImpl) Attack(attackPokemon *InBattlePokemon, targetPokemon *InBattlePokemon, attackMove *data.Moves) {
	slog.Info(fmt.Sprintf("%s used %s", attackPokemon.Pokemon.BasePokemon.Name, attackMove.Name))

	damagePoints := calculateAttackDamage(attackPokemon, targetPokemon, attackMove, 1.0)

	if damagePoints >= targetPokemon.BattleHP {
		targetPokemon.BattleHP = 0
		targetPokemon.IsFainted = true
	} else {
		targetPokemon.BattleHP -= damagePoints
	}

	slog.Info(fmt.Sprintf("%s did %v points of damage to %s", attackPokemon.Pokemon.BasePokemon.Name, damagePoints, targetPokemon.Pokemon.BasePokemon.Name))
	fmt.Println()
}

func (tb *TrainerBattleImpl) SwitchPokemon(switchPokemonIndex int, isUser bool, enabled bool) {
	if isUser {
		switchPokemon := tb.UserTrainerInfo.ActivePokemon
		tb.UserTrainerInfo.ActivePokemon = tb.UserTrainerInfo.InBattleParty[switchPokemonIndex]
		tb.UserTrainerInfo.InBattleParty = append(tb.UserTrainerInfo.InBattleParty[:switchPokemonIndex], tb.UserTrainerInfo.InBattleParty[switchPokemonIndex+1:]...)
		tb.UserTrainerInfo.InBattleParty = append(tb.UserTrainerInfo.InBattleParty, switchPokemon)

		slog.Info(fmt.Sprintf("%s, I choose you!", tb.UserTrainerInfo.ActivePokemon.Pokemon.BasePokemon.Name))
	} else {
		switchPokemon := tb.EnemyTrainerInfo.ActivePokemon
		tb.EnemyTrainerInfo.ActivePokemon = tb.EnemyTrainerInfo.InBattleParty[switchPokemonIndex]
		tb.EnemyTrainerInfo.InBattleParty = append(tb.EnemyTrainerInfo.InBattleParty[:switchPokemonIndex], tb.EnemyTrainerInfo.InBattleParty[switchPokemonIndex+1:]...)
		tb.EnemyTrainerInfo.InBattleParty = append(tb.EnemyTrainerInfo.InBattleParty, switchPokemon)

		slog.Info(fmt.Sprintf("%s, chooses %s!", tb.getTrainerName(false), tb.UserTrainerInfo.ActivePokemon.Pokemon.BasePokemon.Name))
	}

	// update map of pokemon facing
	tb.AddUserPokemonFacingEnemyActive(tb.UserTrainerInfo.ActivePokemon.Pokemon)
}

func (tb *TrainerBattleImpl) UseItem(targetPokemon *InBattlePokemon, item *data.Item) {
	switch item.ItemType {
	case data.MedicalItems:
		tb.HealPokemon(targetPokemon, item)
	case data.PokeBalls:
		tb.CatchPokemon(targetPokemon, item)
	}
}

func (tb *TrainerBattleImpl) HealPokemon(targetPokemon *InBattlePokemon, item *data.Item) {
	targetPokemon.BattleHP += targetPokemon.BattleHP + item.Attributes

	if targetPokemon.BattleHP > targetPokemon.Pokemon.Stats.HP {
		targetPokemon.BattleHP = targetPokemon.Pokemon.Stats.HP
	}
}

func (tb *TrainerBattleImpl) CatchPokemon(targetPokemon *InBattlePokemon, item *data.Item) {
	slog.Info("Cannot catch a pokemon in trainer battle")
}

func (tb *TrainerBattleImpl) Run() {
	slog.Info("No! There's no running from a trainer battle!")
}

func (tb *TrainerBattleImpl) IsBattleOver() bool {
	if tb.UserTrainerInfo.UnfaintedPartyCount == 0 || tb.EnemyTrainerInfo.UnfaintedPartyCount == 0 {
		return true
	}
	return false
}

func (tb *TrainerBattleImpl) BattleReport() (*BattleResult, error) {
	if !tb.IsBattleOver() {
		return nil, errors.New("battle is ongoing")
	}

	var result BattleResult
	if tb.UserTrainerInfo.UnfaintedPartyCount == 0 {
		slog.Info(fmt.Sprintf("Trainer %s has won the battle", tb.EnemyTrainerInfo.Trainer.Name))
		result.UserWin = false
	} else {
		slog.Info(fmt.Sprintf("%s has won the battle!", tb.UserTrainerInfo.Trainer.Name))
		slog.Info(fmt.Sprintf("You got $%v!", tb.EnemyTrainerInfo.Trainer.AdditionalInfo.PrizeMoney))

		result.UserWin = true
		result.PrizeMoney = tb.EnemyTrainerInfo.Trainer.AdditionalInfo.PrizeMoney
	}

	return &result, nil
}

func (tb *TrainerBattleImpl) AddUserPokemonFacingEnemyActive(pokemon *data.Pokemon) {
	tb.EnemyTrainerInfo.EnemyPokemonFaced[tb.EnemyTrainerInfo.ActivePokemon.Pokemon] =
		append(tb.EnemyTrainerInfo.EnemyPokemonFaced[tb.EnemyTrainerInfo.ActivePokemon.Pokemon], pokemon)
}

func (tb *TrainerBattleImpl) getTrainerName(isUser bool) string {
	if isUser {
		return tb.UserTrainerInfo.Trainer.Name
	} else {
		return tb.TrainerLogPrefix + tb.UserTrainerInfo.Trainer.Name
	}
}
