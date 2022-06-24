package usecase

import (
	"sort"
	"strings"

	"gitlab.com/mailru-go/lectures-2022-1/03/99_hw/coverage/internal/models"
	"gitlab.com/mailru-go/lectures-2022-1/03/99_hw/coverage/internal/user"
	"gitlab.com/mailru-go/lectures-2022-1/03/99_hw/coverage/tools"
)

type UserUsecase struct {
	userRep user.UserRepository
}

func NewUserUsecase(rep user.UserRepository) user.UserUsecase {
	return &UserUsecase{
		userRep: rep,
	}
}

func (u *UserUsecase) GetUsers() ([]models.User, error) {
	users, err := u.userRep.SelectUsers()
	if err != nil {
		return nil, err
	}

	clientUsers := ConvertUsers(users)

	return clientUsers, nil
}

func ConvertUsers(users *models.SearchUsers) []models.User {
	return models.ConverSearchUsersToClientUsers(users)
}

func (u *UserUsecase) SortUsers(users []models.User, params *tools.QueryParams) []models.User {
	if params.Query != "" {
		filteredUsers := []models.User{}
		for _, user := range users {
			if strings.Contains(user.Name+user.About, params.Query) {
				filteredUsers = append(filteredUsers, user)
			}
		}
		users = filteredUsers
	}

	switch params.OrderField {
	case "id":
		switch params.OrderBy {
		case 1:
			sort.SliceStable(users, func(i, j int) bool {
				return users[i].ID < users[j].ID
			})
		case -1:
			sort.SliceStable(users, func(i, j int) bool {
				return users[i].ID > users[j].ID
			})
		default:
			break
		}
	case "age":
		switch params.OrderBy {
		case 1:
			sort.SliceStable(users, func(i, j int) bool {
				return users[i].Age < users[j].Age
			})
		case -1:
			sort.SliceStable(users, func(i, j int) bool {
				return users[i].Age > users[j].Age
			})
		default:
			break
		}
	case "name":
		switch params.OrderBy {
		case 1:
			sort.SliceStable(users, func(i, j int) bool {
				return users[i].Name < users[j].Name
			})
		case -1:
			sort.SliceStable(users, func(i, j int) bool {
				return users[i].Name > users[j].Name
			})
		default:
			break
		}
	}

	if params.Offset != tools.ErrConvert {
		if params.Offset > len(users) {
			return []models.User{}
		}
		users = users[params.Offset:]
	}

	if params.Limit != tools.ErrConvert {
		if params.Limit > len(users) {
			params.Limit = len(users)
		}
		users = users[:params.Limit]
	}

	return users
}
