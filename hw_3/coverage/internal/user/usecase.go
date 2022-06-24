package user

import (
	"gitlab.com/mailru-go/lectures-2022-1/03/99_hw/coverage/internal/models"
	"gitlab.com/mailru-go/lectures-2022-1/03/99_hw/coverage/tools"
)

type UserUsecase interface {
	GetUsers() ([]models.User, error)
	SortUsers([]models.User, *tools.QueryParams) []models.User
}
