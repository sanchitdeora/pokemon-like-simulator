package pokemon

import (
	"math"

	"github.com/sanchitdeora/PokeSim/data"
)

const NatureCoeff = 1.0

// hp calculation is different from the other stats
func calculateHPStatUpgrade(baseHP int, pokemonStatHP *data.PokemonStat, level int) int {
	return int(math.Floor(((((2 * float64(baseHP)) + float64(pokemonStatHP.IV) + (float64(pokemonStatHP.EV) / 4)) * float64(level)) / 100) + float64(level) + 10))
}

// all other stats calculation
func calculateOtherStatUpgrade(baseStat int, pokemonStat *data.PokemonStat, level int) int {
	return int(math.Floor(((((2 * float64(baseStat)) + float64(pokemonStat.IV) + (float64(pokemonStat.EV) / 4)) * float64(level)) / 100) + 5) * NatureCoeff)
}

func nextLevelErraticExp(level int) int {
	if level < 50 {
		return int(math.Round(((math.Pow(float64(level), 3)) * (100 - float64(level))) / 50))
	} else if level < 68 {
		return int(math.Round(((math.Pow(float64(level), 3)) * (100 - float64(level))) / 50))
	} else if level < 98 {
		return int(math.Round(((math.Pow(float64(level), 3)) * math.Floor((1911 - (10 * float64(level))) / 3)) / 500))
	} else {
		return int(math.Round(((math.Pow(float64(level), 3)) * (100 - float64(level))) / 50))
	}
}

func nextLevelFastExp(level int) int {
	return ((int(math.Pow(float64(level), 3)) * 4) / 5)
}

func nextLevelMediumFastExp(level int) int {
	return int(math.Pow(float64(level), 3))
}

func nextLevelMediumSlowExp(level int) int {
	return int(math.Round((6 * math.Pow(float64(level), 3)) / 5) - (15 * math.Pow(float64(level), 2)) + (100 * float64(level)) - 140)
}

func nextLevelSlowExp(level int) int {
	return ((int(math.Pow(float64(level), 3)) * 5) / 4)
}

func nextLevelFluctuatingExp(level int) int {
	if level < 15 {
		return int(math.Round(math.Pow(float64(level), 3) * (((math.Floor((float64(level) + 1) / 3) + 24)) / 50)))
	} else if level < 36 {
		return int(math.Round((math.Pow(float64(level), 3) * (float64(level) + 14)) / 50))
	} else {
		return int(math.Round((math.Pow(float64(level), 3) * (math.Floor(float64(level) / 2) + 32)) / 50))
	}
}