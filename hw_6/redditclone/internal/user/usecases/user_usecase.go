package usecases

import (
	"database/sql"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/consts"
	customErrors "lectures-2022-1/06_databases/99_hw/redditclone/internal/helpers/errors"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/models"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/user"
)

type UserUsecase struct {
	userRep user.UserPgRepository
}

func NewUserUsecase(rep user.UserPgRepository) user.UserUsecase {
	return &UserUsecase{
		userRep: rep,
	}
}

func (u *UserUsecase) RegiserUser(user *models.User) (uint64, *customErrors.Error) {
	lastID, err := u.userRep.Insert(user)
	if err != nil {
		return 0, customErrors.Get(consts.CodeUsernameAlreadyTaken)
	}

	return lastID, nil
}

func (u *UserUsecase) LoginUser(user *models.User) (uint64, *customErrors.Error) {
	user, err := u.userRep.SelectByUsername(user.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, customErrors.Get(consts.CodeUserDoesntExist)
		}
		return 0, customErrors.Get(consts.CodeUsernameAlreadyTaken)
	}

	return user.ID, nil
}