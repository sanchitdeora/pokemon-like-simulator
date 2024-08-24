package battle

import (
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"time"

	"github.com/sanchitdeora/PokeSim/data"
)

type BattleIFace interface {
	InitiateBattleSequence()
	Attack(attackPokemon *InBattlePokemon, targetPokemon *InBattlePokemon, attackMove *data.Moves)
	CatchPokemon(targetPokemon *InBattlePokemon, item *data.Item)
	Run()
	IsBattleOver() bool
	BattleReport() (*BattleResult, error)
}

type InBattlePokemon struct {
	Pokemon   *data.Pokemon
	BattleHP  int
	IsFainted bool
}

type BattleTrainerInfo struct {
	Trainer             *data.Trainer
	ActivePokemon       *InBattlePokemon
	InBattleParty       []*InBattlePokemon
	EnemyPokemonFaced   map[*data.Pokemon][]*data.Pokemon
	UnfaintedPartyCount int
	IsUser              bool
}

type BattleResult struct {
	UserWin    bool
	PrizeMoney int
	BonusItems []*data.Item
}

type BattleOpts struct {
	Type string `json:"type"`
}

type BattleInput struct {
	Type           BattleInputType
	CurrentPokemon *InBattlePokemon
	Target         *InBattlePokemon
	Move           *data.Moves
	Item           *data.Item
	IsUser         bool
}

type BattleInputType string

const (
	Attack BattleInputType = "attack"
	Switch BattleInputType = "switch"
	Item   BattleInputType = "item"
	Run    BattleInputType = "run"
)

func getUserTrainerInfo() *data.Trainer {
	return nil
}

func createTrainerInfo(trainer *data.Trainer, isUser bool) *BattleTrainerInfo {
	var activePokemon InBattlePokemon
	inBattlePokemonParty := make([]*InBattlePokemon, 0)
	enemyPokemonFaced := make(map[*data.Pokemon][]*data.Pokemon, 0)

	if !isUser {
		for _, party := range trainer.Party {
			enemyPokemonFaced[party] = make([]*data.Pokemon, 0)
		}
	} else {
		enemyPokemonFaced = nil
	}

	for _, pokemon := range trainer.Party {
		if pokemon.PartyOrder == 1 {
			activePokemon.Pokemon = pokemon
			activePokemon.BattleHP = pokemon.Stats.HP
			activePokemon.IsFainted = false
		} else if len(inBattlePokemonParty) < 5 {
			inBattlePokemonParty = append(inBattlePokemonParty, &InBattlePokemon{
				Pokemon:   pokemon,
				BattleHP:  pokemon.Stats.HP,
				IsFainted: false,
			})
		}
	}
	return &BattleTrainerInfo{
		Trainer:             trainer,
		ActivePokemon:       &activePokemon,
		InBattleParty:       inBattlePokemonParty,
		UnfaintedPartyCount: len(trainer.Party),
		EnemyPokemonFaced:   enemyPokemonFaced,
		IsUser:              isUser,
	}
}

func calculateAttackDamage(attackPokemon *InBattlePokemon, targetPokemon *InBattlePokemon, attackMove *data.Moves, battleTypeAttackCoeff float32) int {
	var totalDamage float32

	var attackStat float32 = 0
	var defenseStat float32 = 0

	if attackMove.DamageClass == data.Physical {
		attackStat = float32(attackPokemon.Pokemon.Stats.Attack)
		defenseStat = float32(targetPokemon.Pokemon.Stats.Defense)
	} else if attackMove.DamageClass == data.Special {
		attackStat = float32(attackPokemon.Pokemon.Stats.SpecialAttack)
		defenseStat = float32(targetPokemon.Pokemon.Stats.SpecialDefense)
	} else {
		slog.Warn("Move Damage class not supported", "move damage class", attackMove.DamageClass)
	}

	pokemonTypes := []data.PokemonTypeName{targetPokemon.Pokemon.BasePokemon.Type1}

	if targetPokemon.Pokemon.BasePokemon.Type2 != "" {
		pokemonTypes = append(pokemonTypes, targetPokemon.Pokemon.BasePokemon.Type2)
	}
	moveEffect := data.GetMoveEffect(attackMove.Type, pokemonTypes...)

	slog.Debug(fmt.Sprintf("Calculating damage: ( ( (((2 * {%v})/5) + 2) * {%v} * ({%v} / {%v}) ) / 50 ) + 2", attackPokemon.Pokemon.Level, attackMove.Power, attackStat, defenseStat))

	totalDamage = (((((2.0 * float32(attackPokemon.Pokemon.Level)) / 5.0) + 2.0) * float32(attackMove.Power) * (attackStat / defenseStat)) / 50) + 2

	// for battles with more than 1 enemy, coefficient = 0.75
	slog.Debug(fmt.Sprintf("Battle Attack Coeff: {%f} * {%f}", totalDamage, battleTypeAttackCoeff))
	totalDamage *= battleTypeAttackCoeff

	if isCritHit() {
		slog.Info("Critical Hit!")
		slog.Debug(fmt.Sprintf("Critical hit: {%f} * 1.5", totalDamage))
		totalDamage *= 1.5
	}

	// random attack power coeffecient ranging from (0.85 - 1.00)
	slog.Debug(fmt.Sprintf("Attack Coeff 0.85 -- 1.00: {%f} * {%f}", totalDamage, attackCoefficient()))
	totalDamage *= attackCoefficient()

	if isStab(attackMove, attackPokemon) {
		slog.Debug(fmt.Sprintf("STAB: {%f} * 1.5", totalDamage))
		totalDamage *= 1.5
	}

	if moveEffect != data.NOR {
		if moveEffect == data.MNE || moveEffect == data.NVR {
			slog.Info("Not very effective!", "effect", moveEffect)
		} else if moveEffect == data.SUP || moveEffect == data.HYP {
			slog.Info("Super effective!", "effect", moveEffect)
		} else {
			slog.Info("This move has No effect to the target!", "effect", moveEffect)
		}
	}

	slog.Debug(fmt.Sprintf("Move Effect Coeff: {%f} * {%f} == {%f}", totalDamage, moveEffect, (totalDamage * float32(moveEffect))))
	return int(math.Round(float64(totalDamage * float32(moveEffect))))
}

func isCritHit() bool {
	return randomGenerator(0, 1) < (1.0 / 24.0)
}

func attackCoefficient() float32 {
	return randomGenerator(0.85, 1)
}

func isStab(attackMove *data.Moves, attackPokemon *InBattlePokemon) bool {
	return attackMove.Type == attackPokemon.Pokemon.BasePokemon.Type1 || attackMove.Type == attackPokemon.Pokemon.BasePokemon.Type2
}

func randomGenerator(min float32, max float32) float32 {
	randIndex := rand.Float32()
	return (min + randIndex*(max-min))
}

func waitForInput(pokemon *InBattlePokemon, target *InBattlePokemon) *BattleInput {
	time.Sleep(time.Second * 0)
	return &BattleInput{
		Type:           Attack,
		Move:           randomMove(&pokemon.Pokemon.Moveset),
		CurrentPokemon: pokemon,
		Target:         target,
	}
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

func BattleExperienceGain(faintedPokemon *data.Pokemon, pokemonFaced []*data.Pokemon) {
	for _, pokemon := range pokemonFaced {
		expGain := calculateExperienceGained(faintedPokemon.Level, faintedPokemon.BasePokemon.BaseExperience, pokemon.Level)
		slog.Info(fmt.Sprintf("%s gained %v experience points", pokemon.BasePokemon.Name, expGain))

		// TODO: create logic
		// updatePokemonExperience(expGain, pokemon)
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
