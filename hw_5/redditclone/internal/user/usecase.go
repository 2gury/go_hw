package user

import (
	customErrors "lectures-2022-1/05_web_app/99_hw/redditclone/internal/helpers/errors"
	"lectures-2022-1/05_web_app/99_hw/redditclone/internal/models"
)

type UserUsecase interface {
	RegiserUser(user models.User) (string, *customErrors.Error)
	LoginUser(user models.User) (string, *customErrors.Error)
}
