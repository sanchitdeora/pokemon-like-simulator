package battle

// import (
// 	"errors"
// 	"fmt"
// 	"log/slog"

// 	"github.com/sanchitdeora/PokeSim/data"
// )

// type WildBattleOpts struct {
// 	TrainerLogPrefix string
// }

// type WildBattle struct {
// 	UserTrainerInfo *BattleTrainerInfo
// 	WildPokemon     *InBattlePokemon
// }

// type WildBattleImpl struct {
// 	*WildBattleOpts
// 	*WildBattle
// }

// func NewWildBattle(wildBasePokemon *data.BasePokemon, u *data.Trainer) BattleIFace {
// 	return &WildBattleImpl{
// 		WildBattleOpts: &WildBattleOpts{TrainerLogPrefix: "Enemy"},
// 		WildBattle: &WildBattle{
// 			UserTrainerInfo: createTrainerInfo(u, true),
// 			WildPokemon: &InBattlePokemon{
// 				Pokemon: wildPokemon,
// 			},
// 		},
// 	}
// }

// func (wb *WildBattleImpl) InitiateBattleSequence() {
// 	slog.Info(fmt.Sprintf("%s, I choose you!", wb.UserTrainerInfo.ActivePokemon.Pokemon.BasePokemon.Name))
// 	battleCompleteFlag := false
// 	for {
// 		slog.Info(fmt.Sprintf("%s Health: %v", wb.UserTrainerInfo.ActivePokemon.Pokemon.BasePokemon.Name, wb.UserTrainerInfo.ActivePokemon.BattleHP))
// 		slog.Info(fmt.Sprintf("Enemy %s Health: %v", wb.EnemyTrainerInfo.ActivePokemon.Pokemon.BasePokemon.Name, wb.EnemyTrainerInfo.ActivePokemon.BattleHP))

// 		battleInputs := wb.GetPokemonAttackOrder()

// 		for _, input := range battleInputs {
// 			wb.Turn(input)
// 			if battleCompleteFlag = wb.IsBattleOver(); battleCompleteFlag {
// 				slog.Info("Battle is completed!")
// 				break
// 			}
// 		}
// 		if battleCompleteFlag {
// 			break
// 		}
// 	}
// 	wb.BattleReport()
// }

// func (wb *WildBattleImpl) GetPokemonAttackOrder() (inputs []*BattleInput) {
// 	userInput := waitForInput(wb.UserTrainerInfo.ActivePokemon, wb.EnemyTrainerInfo.ActivePokemon)
// 	enemyInput := waitForInput(wb.EnemyTrainerInfo.ActivePokemon, wb.UserTrainerInfo.ActivePokemon)

// 	// if any pokemon's move is higher priority, it will go first. If equal, check speed
// 	if userInput.Move.Priority != enemyInput.Move.Priority {
// 		if userInput.Move.Priority > enemyInput.Move.Priority {
// 			return append(inputs, userInput, enemyInput)
// 		} else {
// 			return append(inputs, enemyInput, userInput)
// 		}
// 	}

// 	// if user's active pokemon speed >= enemy's, then user goes first.
// 	if wb.UserTrainerInfo.ActivePokemon.Pokemon.Stats.Speed >= wb.EnemyTrainerInfo.ActivePokemon.Pokemon.Stats.Speed {
// 		return append(inputs, userInput, enemyInput)
// 	} else {
// 		return append(inputs, enemyInput, userInput)
// 	}
// }

// func (wb *WildBattleImpl) Turn(userInput *BattleInput) {
// 	switch userInput.Type {
// 	case Switch:
// 		slog.Info(fmt.Sprintf("%s is switching %s for %s", wb.getTrgeainerName(userInput.IsUser), userInput.CurrentPokemon.Pokemon.BasePokemon.Name, userInput.Target.Pokemon.BasePokemon.Name))
// 		wb.SwitchPokemon(userInput.Target.Pokemon.PartyOrder, userInput.IsUser, true)

// 	case Item:
// 		if userInput.Item != nil && userInput.Item.ItemType == data.MedicalItems {
// 			wb.UseItem(userInput.Target, userInput.Item)
// 		}

// 	case Attack:
// 		wb.Attack(userInput.CurrentPokemon, userInput.Target, userInput.Move)

// 		if wb.UserTrainerInfo.ActivePokemon.IsFainted {
// 			slog.Info(fmt.Sprintf("%s has fainted!", wb.UserTrainerInfo.ActivePokemon.Pokemon.BasePokemon.Name))
// 			wb.UserTrainerInfo.UnfaintedPartyCount -= 1

// 			if wb.UserTrainerInfo.UnfaintedPartyCount > 0 {
// 				// switch active pokemon with first in party and push fainted pokemon at end
// 				nextPokemonIndex := 0
// 				for i, pokemon := range wb.EnemyTrainerInfo.InBattleParty {
// 					if !pokemon.IsFainted {
// 						nextPokemonIndex = i
// 						break
// 					}
// 				}
// 				wb.SwitchPokemon(nextPokemonIndex, true, true)
// 			}
// 		}

// 		if wb.EnemyTrainerInfo.ActivePokemon.IsFainted {
// 			slog.Info(fmt.Sprintf("Enemy %s has fainted!", wb.EnemyTrainerInfo.ActivePokemon.Pokemon.BasePokemon.Name))
// 			wb.EnemyTrainerInfo.UnfaintedPartyCount -= 1

// 			if wb.EnemyTrainerInfo.UnfaintedPartyCount > 0 {
// 				// switch active pokemon with first in party and push fainted pokemon at end
// 				nextPokemonIndex := 0
// 				for i, pokemon := range wb.EnemyTrainerInfo.InBattleParty {
// 					if !pokemon.IsFainted {
// 						nextPokemonIndex = i
// 						break
// 					}
// 				}
// 				wb.SwitchPokemon(nextPokemonIndex, false, true)
// 			}
// 		}

// 	case Run:
// 		wb.Run()
// 	}
// }

// func (wb *WildBattleImpl) Attack(attackPokemon *InBattlePokemon, targetPokemon *InBattlePokemon, attackMove *data.Moves) {
// 	slog.Info(fmt.Sprintf("%s used %s", attackPokemon.Pokemon.BasePokemon.Name, attackMove.Name))

// 	damagePoints := calculateAttackDamage(attackPokemon, targetPokemon, attackMove, 1.0)

// 	if damagePoints >= targetPokemon.BattleHP {
// 		targetPokemon.BattleHP = 0
// 		targetPokemon.IsFainted = true
// 	} else {
// 		targetPokemon.BattleHP -= damagePoints
// 	}

// 	slog.Info(fmt.Sprintf("%s did %v points of damage to %s", attackPokemon.Pokemon.BasePokemon.Name, damagePoints, targetPokemon.Pokemon.BasePokemon.Name))
// 	fmt.Println()
// }

// func (wb *WildBattleImpl) SwitchPokemon(switchPokemonIndex int, isUser bool, enabled bool) {
// 	if isUser {
// 		switchPokemon := wb.UserTrainerInfo.ActivePokemon
// 		wb.UserTrainerInfo.ActivePokemon = wb.UserTrainerInfo.InBattleParty[switchPokemonIndex]
// 		wb.UserTrainerInfo.InBattleParty = append(wb.UserTrainerInfo.InBattleParty[:switchPokemonIndex], wb.UserTrainerInfo.InBattleParty[switchPokemonIndex+1:]...)
// 		wb.UserTrainerInfo.InBattleParty = append(wb.UserTrainerInfo.InBattleParty, switchPokemon)
// 	} else {
// 		switchPokemon := wb.EnemyTrainerInfo.ActivePokemon
// 		wb.EnemyTrainerInfo.ActivePokemon = wb.EnemyTrainerInfo.InBattleParty[switchPokemonIndex]
// 		wb.EnemyTrainerInfo.InBattleParty = append(wb.EnemyTrainerInfo.InBattleParty[:switchPokemonIndex], wb.EnemyTrainerInfo.InBattleParty[switchPokemonIndex+1:]...)
// 		wb.EnemyTrainerInfo.InBattleParty = append(wb.EnemyTrainerInfo.InBattleParty, switchPokemon)
// 	}
// }

// func (wb *WildBattleImpl) UseItem(targetPokemon *InBattlePokemon, item *data.Item) {
// 	switch item.ItemType {
// 	case data.MedicalItems:
// 		wb.HealPokemon(targetPokemon, item)
// 	case data.PokeBalls:
// 		wb.CatchPokemon(targetPokemon, item)
// 	}
// }

// func (wb *WildBattleImpl) HealPokemon(targetPokemon *InBattlePokemon, item *data.Item) {
// 	targetPokemon.BattleHP += targetPokemon.BattleHP + item.Attributes

// 	if targetPokemon.BattleHP > targetPokemon.Pokemon.Stats.HP {
// 		targetPokemon.BattleHP = targetPokemon.Pokemon.Stats.HP
// 	}
// }

// func (wb *WildBattleImpl) CatchPokemon(targetPokemon *InBattlePokemon, item *data.Item) {
// 	slog.Info("Cannot catch a pokemon in trainer battle")
// }

// func (wb *WildBattleImpl) Run() {
// 	slog.Info("No! There's no running from a trainer battle!")
// }

// func (wb *WildBattleImpl) IsBattleOver() bool {
// 	if wb.UserTrainerInfo.UnfaintedPartyCount == 0 || wb.EnemyTrainerInfo.UnfaintedPartyCount == 0 {
// 		return true
// 	}
// 	return false
// }

// func (wb *WildBattleImpl) BattleReport() (*BattleResult, error) {
// 	if !wb.IsBattleOver() {
// 		return nil, errors.New("battle is ongoing")
// 	}

// 	var result BattleResult
// 	if wb.UserTrainerInfo.UnfaintedPartyCount == 0 {
// 		slog.Info(fmt.Sprintf("Trainer %s has won the battle", wb.EnemyTrainerInfo.Trainer.Name))
// 		result.UserWin = false
// 	} else {
// 		slog.Info(fmt.Sprintf("%s has won the battle!", wb.UserTrainerInfo.Trainer.Name))
// 		slog.Info(fmt.Sprintf("You got $%v!", wb.EnemyTrainerInfo.Trainer.AdditionalInfo.PrizeMoney))

// 		result.UserWin = true
// 		result.PrizeMoney = wb.EnemyTrainerInfo.Trainer.AdditionalInfo.PrizeMoney
// 	}

// 	return &result, nil
// }

// func (wb *WildBattleImpl) getPokemonName(isUser bool) string {
// 	if isUser {
// 		return wb.UserTrainerInfo.Trainer.Name
// 	} else {
// 		return "Wild" + wb.UserTrainerInfo.Trainer.Name
// 	}
// }

// func wildPokemonAppeared(pokemon *data.BasePokemon) *data.Pokemon {
// 	// TODO: get a way to determine what range of level should the wild pokemon be.
// 	minLvl, maxLvl := getLevelRange()
// 	lvl := randomGenerator(float32(minLvl), float32(maxLvl))

// }

// func getLevelRange() (int, int) {
// 	return 5, 10
// }
