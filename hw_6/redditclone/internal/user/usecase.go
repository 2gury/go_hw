package user

import (
	customErrors "lectures-2022-1/06_databases/99_hw/redditclone/internal/helpers/errors"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/models"
)

type UserUsecase interface {
	RegiserUser(user *models.User) (uint64, *customErrors.Error)
	LoginUser(user *models.User) (uint64, *customErrors.Error)
}
