package session

import (
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/models"
	customErrors "lectures-2022-1/06_databases/99_hw/redditclone/internal/helpers/errors"
)

type SessionUsecase interface {
	Create(user *models.User) (*models.Session, *customErrors.Error)
	Check(sessValue string) (*models.User, *customErrors.Error)
	Delete(sessValue string) *customErrors.Error
}