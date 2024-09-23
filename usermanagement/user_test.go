package usermanagement_test

import (
	"testing"

	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/usermanagement"
	"github.com/stretchr/testify/assert"
)

const (
	TestUserPath = "C:\\Projects\\Go-projects\\src\\PokéSim\\testfiles\\user\\test_user.json"
	TestUserPath1 = "C:\\Projects\\Go-projects\\src\\PokéSim\\saved\\user.json"
)

func createUserService(savedPath string) usermanagement.User {
	return usermanagement.NewUserService(usermanagement.UserOpts{SavedUserPath: savedPath})
}

// LoadUser

func TestLoadUser(t *testing.T) {
	userService := createUserService(TestUserPath)

	user, err := userService.LoadUser()
	assert.NoError(t, err)
	
	assert.Equal(t, "John Cena", user.Name)
	assert.Equal(t, "bulbasaur", user.Party[0].Name)
	assert.Equal(t, "squirtle", user.Party[1].Name)
}

func TestPostBattleUpdate(t *testing.T) {
	userService := createUserService(TestUserPath)

	user, _ := userService.LoadUser()
	copyOfUser, _ := userService.LoadUser()

	report := &data.BattleReport{
		UserWin: true,
		Money: 100,
		BonusItems:[]*data.Item{{
				Name: "Potion",
				Category: data.MedicalItems,
			},
		},
		BadgeEarned: &data.BadgeType{
			Name: "test",
		},
	}

	err := userService.PostBattleUpdate(user, report)
	
	assert.NoError(t, err)
	assert.Equal(t, 1, user.Stats.Battles)
	assert.Equal(t, 1, user.Stats.Wins)
	assert.Equal(t, 1, len(user.Stats.Badges))
	assert.Equal(t, 1, len(user.Bag))

	// clean up after
	userService.SaveUser(copyOfUser)
}


func TestPostWildUpdate(t *testing.T) {
	userService := createUserService(TestUserPath)

	user, _ := userService.LoadUser()

	err := userService.PostWildUpdate(user, true, &data.Pokemon{})
	assert.NoError(t, err)
	assert.Equal(t, 1, user.Stats.Battles)
	assert.Equal(t, 1, user.Stats.Wins)
	assert.Equal(t, 1, user.Stats.Catches)

	// clean up after
	user.Stats = &data.TrainerStats{}
	userService.SaveUser(user)
}