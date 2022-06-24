package usecases

import (
	"lectures-2022-1/05_web_app/99_hw/redditclone/internal/consts"
	customErrors "lectures-2022-1/05_web_app/99_hw/redditclone/internal/helpers/errors"
	"lectures-2022-1/05_web_app/99_hw/redditclone/internal/models"
	"lectures-2022-1/05_web_app/99_hw/redditclone/internal/user"
)

type UserUsecase struct {
	userRep user.UserRepository
}

func NewUserUsecase(rep user.UserRepository) user.UserUsecase {
	return &UserUsecase{
		userRep: rep,
	}
}

func (u *UserUsecase) RegiserUser(user models.User) (string, *customErrors.Error) {
	if exist := u.userRep.IsUsernameExist(user.Username); exist {
		return "", customErrors.Get(consts.CodeUsernameAlreadyTaken)
	}

	lastID := u.userRep.InsertUser(user)
	return lastID, nil
}

func (u *UserUsecase) LoginUser(user models.User) (string, *customErrors.Error) {
	if exist := u.userRep.IsUsernameExist(user.Username); !exist {
		return "", customErrors.Get(consts.CodeUserDoesntExist)
	}

	userID, err := u.userRep.CheckPassword(user)
	if err != nil {
		return "", customErrors.Get(consts.CodeIncorrectUserPassword)
	}
	return userID, nil
}