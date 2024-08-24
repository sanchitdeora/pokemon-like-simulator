package data

import "log/slog"

type PokemonType int

const (
	Normal PokemonType = iota
	Fire
	Water
	Electric
	Grass
	Ice
	Fighting
	Poison
	Ground
	Flying
	Psychic
	Bug
	Rock
	Ghost
	Dragon
	Dark
	Steel
	Fairy
)

type PokemonTypeName string

const (
	NormalType   PokemonTypeName = "normal"
	FireType     PokemonTypeName = "fire"
	WaterType    PokemonTypeName = "water"
	ElectricType PokemonTypeName = "electric"
	GrassType    PokemonTypeName = "grass"
	IceType      PokemonTypeName = "ice"
	FightingType PokemonTypeName = "fighting"
	PoisonType   PokemonTypeName = "poison"
	GroundType   PokemonTypeName = "ground"
	FlyingType   PokemonTypeName = "flying"
	PsychicType  PokemonTypeName = "psychic"
	BugType      PokemonTypeName = "bug"
	RockType     PokemonTypeName = "rock"
	GhostType    PokemonTypeName = "ghost"
	DragonType   PokemonTypeName = "dragon"
	DarkType     PokemonTypeName = "dark"
	SteelType    PokemonTypeName = "steel"
	FairyType    PokemonTypeName = "fairy"
)

func (i PokemonType) ToString() PokemonTypeName {
	switch i {
	case Normal:
		return NormalType
	case Fire:
		return FireType
	case Water:
		return WaterType
	case Electric:
		return ElectricType
	case Grass:
		return GrassType
	case Ice:
		return IceType
	case Fighting:
		return FightingType
	case Poison:
		return PoisonType
	case Ground:
		return GroundType
	case Flying:
		return FlyingType
	case Psychic:
		return PsychicType
	case Bug:
		return BugType
	case Rock:
		return RockType
	case Ghost:
		return GhostType
	case Dragon:
		return DragonType
	case Dark:
		return DarkType
	case Steel:
		return SteelType
	case Fairy:
		return FairyType
	}
	return ""
}

type TypeEffective float32

const (
	NOE TypeEffective = 0    // NO EFFECT (0%)
	MNE TypeEffective = 0.25 // MINIMAL EFFECT (25%)
	NVR TypeEffective = 0.5  // NOT VERY EFFECTIVE (50%)
	NOR TypeEffective = 1    // NORMAL (100%)
	SUP TypeEffective = 2    // SUPER EFFECTIVE (200%)
	HYP TypeEffective = 4    // HYPER EFFECTIVE (400%)
)

var SingleEffectTypeChart [18][18]TypeEffective = [18][18]TypeEffective{
	//AT/DEF | NOR, FIR, WAT, ELE, GRA, ICE, FIG, POI, GRO, FLY, PSY, BUG, ROC, GHO, DRA, DAR, STE, FAI
	/* NOR */ {NOR, NOR, NOR, NOR, NOR, NOR, NOR, NOR, NOR, NOR, NOR, NOR, NVR, NOE, NOR, NOR, NVR, NOR},
	/* FIR */ {NOR, NVR, NVR, NOR, SUP, SUP, NOR, NOR, NOR, NOR, NOR, SUP, NVR, NOR, NVR, NOR, SUP, NOR},
	/* WAT */ {NOR, SUP, NVR, NOR, NVR, NOR, NOR, NOR, SUP, NOR, NOR, NOR, SUP, NOR, NVR, NOR, NOR, NOR},
	/* ELE */ {NOR, NOR, SUP, NVR, NVR, NOR, NOR, NOR, NOE, SUP, NOR, NOR, NOR, NOR, NVR, NOR, NOR, NOR},
	/* GRA */ {NOR, NVR, SUP, NOR, NVR, NOR, NOR, NVR, SUP, NVR, NOR, NVR, SUP, NOR, NVR, NOR, NVR, NOR},
	/* ICE */ {NOR, NVR, NVR, NOR, SUP, NVR, NOR, NOR, SUP, SUP, NOR, NOR, NOR, NOR, SUP, NOR, NVR, NOR},
	/* FIG */ {SUP, NOR, NOR, NOR, NOR, SUP, NOR, NVR, NOR, NVR, NVR, NVR, SUP, NOE, NOR, SUP, SUP, NVR},
	/* POI */ {NOR, NOR, NOR, NOR, SUP, NOR, NOR, NVR, NVR, NOR, NOR, NOR, NVR, NVR, NOR, NOR, NOE, SUP},
	/* GRO */ {NOR, SUP, NOR, SUP, NVR, NOR, NOR, SUP, NOR, NOE, NOR, NVR, SUP, NOR, NOR, NOR, SUP, NOR},
	/* FLY */ {NOR, NOR, NOR, NVR, SUP, NOR, SUP, NOR, NOR, NOR, NOR, SUP, NVR, NOR, NOR, NOR, NVR, NOR},
	/* PSY */ {NOR, NOR, NOR, NOR, NOR, NOR, SUP, SUP, NOR, NOR, NVR, NOR, NOR, NOR, NOR, NOE, NVR, NOR},
	/* BUG */ {NOR, NVR, NOR, NOR, SUP, NOR, NVR, NVR, NOR, NVR, SUP, NOR, NOR, NVR, NOR, SUP, NVR, NVR},
	/* ROC */ {NOR, SUP, NOR, NOR, NOR, SUP, NVR, NOR, NVR, SUP, NOR, SUP, NOR, NOR, NOR, NOR, NVR, NOR},
	/* GHO */ {NOE, NOR, NOR, NOR, NOR, NOR, NOR, NOR, NOR, NOR, SUP, NOR, NOR, SUP, NOR, NVR, NOR, NOR},
	/* DRA */ {NOR, NOR, NOR, NOR, NOR, NOR, NOR, NOR, NOR, NOR, NOR, NOR, NOR, NOR, SUP, NOR, NVR, NOE},
	/* DAR */ {NOR, NOR, NOR, NOR, NOR, NOR, NVR, NOR, NOR, NOR, SUP, NOR, NOR, SUP, NOR, NVR, NOR, NVR},
	/* STE */ {NOR, NVR, NVR, NVR, NOR, SUP, NOR, NOR, NOR, NOR, NOR, NOR, SUP, NOR, NOR, NOR, NVR, SUP},
	/* FAI */ {NOR, NVR, NOR, NOR, NOR, NOR, SUP, NVR, NOR, NOR, NOR, NOR, NOR, NOR, SUP, SUP, NVR, NOR},
}

func (i PokemonTypeName) ToIndex() PokemonType {
	switch i {
	case NormalType:
		return Normal
	case FireType:
		return Fire
	case WaterType:
		return Water
	case ElectricType:
		return Electric
	case GrassType:
		return Grass
	case IceType:
		return Ice
	case FightingType:
		return Fighting
	case PoisonType:
		return Poison
	case GroundType:
		return Ground
	case FlyingType:
		return Flying
	case PsychicType:
		return Psychic
	case BugType:
		return Bug
	case RockType:
		return Rock
	case GhostType:
		return Ghost
	case DragonType:
		return Dragon
	case DarkType:
		return Dark
	case SteelType:
		return Steel
	case FairyType:
		return Fairy
	}
	return -1
}

func GetMoveEffect(moveType PokemonTypeName, pokemonTypes ...PokemonTypeName) TypeEffective {
	slog.Debug("GetMoveEffect...", "moveType", moveType, "moveTypeIdx", moveType.ToIndex(), "pokemonTypes", pokemonTypes, "pokemonTypesIdx", moveType.ToIndex())
	
	moveEffect := NOR
	for _, pTypeIndex := range pokemonTypes {
		moveEffect *= SingleEffectTypeChart[moveType.ToIndex()][pTypeIndex.ToIndex()]
	}
	return moveEffect
}
