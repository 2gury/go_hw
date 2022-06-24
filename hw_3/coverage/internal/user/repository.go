package user

import "gitlab.com/mailru-go/lectures-2022-1/03/99_hw/coverage/internal/models"

type UserRepository interface {
	SelectUsers() (*models.SearchUsers, error)
}
