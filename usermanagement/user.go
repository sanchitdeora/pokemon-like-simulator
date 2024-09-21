package usermanagement

import (
	"log/slog"

	"github.com/sanchitdeora/PokeSim/data"
	"github.com/sanchitdeora/PokeSim/errors"
	"github.com/sanchitdeora/PokeSim/utils"
)

type User interface {
	GetUser() (*data.User, error)
	PostBattleUpdate(user *data.User, report *data.BattleReport) error
	PostWildUpdate(user *data.User, win bool, caught *data.Pokemon) error
	SaveUser(user *data.User) error
}

type UserOpts struct {
	SavedUserPath string
}

type UserImpl struct {
	opts UserOpts
}

func NewUserService(opts UserOpts) User {
	return &UserImpl{opts: opts}
}

func (u *UserImpl) GetUser() (*data.User, error) {
	var savedUser data.UserSave
	var err error

	if utils.CheckPathExists(u.opts.SavedUserPath) {
		savedUser, err = utils.ReadJsonFromFile[data.UserSave](u.opts.SavedUserPath)
		if err != nil {
			slog.Error("could not read from saved file", "error", err)
			return nil, errors.ErrCouldNotReadFromFile
		}
	} else {
		slog.Error("file does not exist...")
		return nil, errors.ErrFileDoesNotExist
	}
	return savedUser.ToUser(), nil
}

func (u *UserImpl) PostBattleUpdate(user *data.User, report *data.BattleReport) error {
	user.Stats.Battles ++

	if report.UserWin {
		user.Stats.Wins ++
		if report.BadgeEarned != nil {
			user.Stats.Badges = append(user.Stats.Badges, *report.BadgeEarned)
		}
		
		if len(report.BonusItems) > 0 {
			for _, item := range report.BonusItems {
				bagItem := data.GetItemFromBag(user.Bag, item)
				if bagItem != nil {
					bagItem.Count ++
				} else {
					item.Count ++
					user.Bag = append(user.Bag, item)
				}
			}
		}

	} else {
		user.Stats.Losses ++
	}

	if err := u.SaveUser(user); err != nil {
		slog.Error("error updating user after trainer battle", "error", err)
		return err
	}

	return nil
}

func (u *UserImpl) PostWildUpdate(user *data.User, win bool, caught *data.Pokemon) error {
	user.Stats.Battles ++
	if win {
		user.Stats.Wins ++
		if caught != nil {
			user.Stats.Catches ++
			//TODO: update pokedex with caught.
		}
	} else {
		user.Stats.Losses ++
	}

	if err := u.SaveUser(user); err != nil {
		slog.Error("error updating user after wild pokemon battle", "error", err)
		return err
	}

	return nil
}

func (u *UserImpl) SaveUser(user *data.User) error {
	err := utils.WriteJsonToFile[data.UserSave](u.opts.SavedUserPath, *user.ToUserSave())
	if err != nil {
		slog.Error("could not write to saved file", "error", err)
		return errors.ErrCouldNotReadFromFile
	}

	return nil
}
