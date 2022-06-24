package usecase

import (
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"gitlab.com/mailru-go/lectures-2022-1/03/99_hw/coverage/internal/models"
	mock_user "gitlab.com/mailru-go/lectures-2022-1/03/99_hw/coverage/internal/user/mocks"
	"gitlab.com/mailru-go/lectures-2022-1/03/99_hw/coverage/tools"
)

func TestUserUsecase_GetUsers(t *testing.T) {
	type mockBehaviour func(usecaseRep *mock_user.MockUserRepository, users *models.SearchUsers)
	t.Parallel()

	testTable := []struct {
		name           string
		mockBehaviour  mockBehaviour
		outSearchUsers *models.SearchUsers
		expUsers       []models.User
		expError       error
	}{
		{
			name: "OK",
			mockBehaviour: func(usecaseRep *mock_user.MockUserRepository, users *models.SearchUsers) {
				usecaseRep.
					EXPECT().
					SelectUsers().
					Return(users, nil)
			},
			outSearchUsers: &models.SearchUsers{
				Users: []models.SearchUser{
					{
						ID:            1,
						GUID:          "FS42",
						IsActive:      false,
						Balance:       "100",
						Picture:       "img",
						Age:           22,
						EyeColor:      "green",
						FirstName:     "Bazil",
						LastName:      "Homov",
						Gender:        "male",
						Company:       "TSMC",
						Email:         "testmail@ru",
						Phone:         "+14241244",
						Address:       "Moscow",
						About:         "about",
						Registered:    "today",
						FavoriteFruit: "apple",
					},
				},
			},
			expUsers: []models.User{
				{
					ID:     1,
					Name:   "Bazil Homov",
					Age:    22,
					About:  "about",
					Gender: "male",
				},
			},
			expError: nil,
		},
		{
			name: "Error",
			mockBehaviour: func(usecaseRep *mock_user.MockUserRepository, users *models.SearchUsers) {
				usecaseRep.
					EXPECT().
					SelectUsers().
					Return(nil, fmt.Errorf("handling error"))
			},
			expUsers: nil,
			expError: fmt.Errorf("handling error"),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userRep := mock_user.NewMockUserRepository(ctrl)
			testCase.mockBehaviour(userRep, testCase.outSearchUsers)
			userUse := NewUserUsecase(userRep)

			users, err := userUse.GetUsers()

			assert.Equal(t, users, testCase.expUsers)
			assert.Equal(t, err, testCase.expError)
		})
	}
}

func TestUserUsecase_SelectUsers(t *testing.T) {
	t.Parallel()

	testTable := []struct {
		name     string
		inUsers  []models.User
		inParams *tools.QueryParams
		expUsers []models.User
	}{
		{
			name: "OK: Search by query",
			inUsers: []models.User{
				{
					ID:     1,
					Name:   "Bazil Hotov",
					Age:    22,
					About:  "about",
					Gender: "male",
				},
				{
					ID:     2,
					Name:   "Kate Jannifer",
					Age:    35,
					About:  "not filtered",
					Gender: "female",
				},
			},
			inParams: &tools.QueryParams{
				Limit:      2,
				Offset:     0,
				Query:      "ab",
				OrderField: "name",
				OrderBy:    -1,
			},
			expUsers: []models.User{
				{
					ID:     1,
					Name:   "Bazil Hotov",
					Age:    22,
					About:  "about",
					Gender: "male",
				},
			},
		},
		{
			name: "OK: Search by query. Sort name asc",
			inUsers: []models.User{
				{
					ID:     2,
					Name:   "Kate Jannifer",
					Age:    35,
					About:  "not filtered",
					Gender: "female",
				},
				{
					ID:     1,
					Name:   "Bazil Hotov",
					Age:    22,
					About:  "about",
					Gender: "male",
				},
				{
					ID:     3,
					Name:   "Anifer Veikub",
					Age:    29,
					About:  "cant talk))",
					Gender: "female",
				},
			},
			inParams: &tools.QueryParams{
				Limit:      3,
				Offset:     0,
				OrderField: "name",
				OrderBy:    1,
			},
			expUsers: []models.User{
				{
					ID:     3,
					Name:   "Anifer Veikub",
					Age:    29,
					About:  "cant talk))",
					Gender: "female",
				},
				{
					ID:     1,
					Name:   "Bazil Hotov",
					Age:    22,
					About:  "about",
					Gender: "male",
				},
				{
					ID:     2,
					Name:   "Kate Jannifer",
					Age:    35,
					About:  "not filtered",
					Gender: "female",
				},
			},
		},
		{
			name: "OK: Search by query. Sort name desc",
			inUsers: []models.User{
				{
					ID:     1,
					Name:   "Bazil Hotov",
					Age:    22,
					About:  "about",
					Gender: "male",
				},
				{
					ID:     2,
					Name:   "Kate Jannifer",
					Age:    35,
					About:  "not filtered",
					Gender: "female",
				},
			},
			inParams: &tools.QueryParams{
				Limit:      2,
				Offset:     0,
				OrderField: "name",
				OrderBy:    -1,
			},
			expUsers: []models.User{
				{
					ID:     2,
					Name:   "Kate Jannifer",
					Age:    35,
					About:  "not filtered",
					Gender: "female",
				},
				{
					ID:     1,
					Name:   "Bazil Hotov",
					Age:    22,
					About:  "about",
					Gender: "male",
				},
			},
		},
		{
			name: "OK: Search by query. Sort name default",
			inUsers: []models.User{
				{
					ID:     1,
					Name:   "Bazil Hotov",
					Age:    22,
					About:  "about",
					Gender: "male",
				},
				{
					ID:     2,
					Name:   "Kate Jannifer",
					Age:    35,
					About:  "not filtered",
					Gender: "female",
				},
			},
			inParams: &tools.QueryParams{
				Limit:      2,
				Offset:     0,
				OrderField: "name",
			},
			expUsers: []models.User{
				{
					ID:     1,
					Name:   "Bazil Hotov",
					Age:    22,
					About:  "about",
					Gender: "male",
				},
				{
					ID:     2,
					Name:   "Kate Jannifer",
					Age:    35,
					About:  "not filtered",
					Gender: "female",
				},
			},
		},
		{
			name: "OK: Sort id asc",
			inUsers: []models.User{
				{
					ID:     2,
					Name:   "Kate Jannifer",
					Age:    35,
					About:  "not filtered",
					Gender: "female",
				},
				{
					ID:     1,
					Name:   "Bazil Hotov",
					Age:    22,
					About:  "about",
					Gender: "male",
				},
			},
			inParams: &tools.QueryParams{
				Limit:      2,
				Offset:     0,
				OrderBy:    1,
				OrderField: "id",
			},
			expUsers: []models.User{
				{
					ID:     1,
					Name:   "Bazil Hotov",
					Age:    22,
					About:  "about",
					Gender: "male",
				},
				{
					ID:     2,
					Name:   "Kate Jannifer",
					Age:    35,
					About:  "not filtered",
					Gender: "female",
				},
			},
		},
		{
			name: "OK: Sort id desc",
			inUsers: []models.User{
				{
					ID:     1,
					Name:   "Bazil Hotov",
					Age:    22,
					About:  "about",
					Gender: "male",
				},
				{
					ID:     2,
					Name:   "Kate Jannifer",
					Age:    35,
					About:  "not filtered",
					Gender: "female",
				},
			},
			inParams: &tools.QueryParams{
				Limit:      2,
				Offset:     0,
				OrderBy:    -1,
				OrderField: "id",
			},
			expUsers: []models.User{
				{
					ID:     2,
					Name:   "Kate Jannifer",
					Age:    35,
					About:  "not filtered",
					Gender: "female",
				},
				{
					ID:     1,
					Name:   "Bazil Hotov",
					Age:    22,
					About:  "about",
					Gender: "male",
				},
			},
		},
		{
			name: "OK: Sort id default",
			inUsers: []models.User{
				{
					ID:     1,
					Name:   "Bazil Hotov",
					Age:    22,
					About:  "about",
					Gender: "male",
				},
				{
					ID:     2,
					Name:   "Kate Jannifer",
					Age:    35,
					About:  "not filtered",
					Gender: "female",
				},
			},
			inParams: &tools.QueryParams{
				Limit:      2,
				Offset:     0,
				OrderBy:    0,
				OrderField: "id",
			},
			expUsers: []models.User{
				{
					ID:     1,
					Name:   "Bazil Hotov",
					Age:    22,
					About:  "about",
					Gender: "male",
				},
				{
					ID:     2,
					Name:   "Kate Jannifer",
					Age:    35,
					About:  "not filtered",
					Gender: "female",
				},
			},
		},
		{
			name: "OK: Sort age desc",
			inUsers: []models.User{
				{
					ID:     1,
					Name:   "Bazil Hotov",
					Age:    22,
					About:  "about",
					Gender: "male",
				},
				{
					ID:     2,
					Name:   "Kate Jannifer",
					Age:    35,
					About:  "not filtered",
					Gender: "female",
				},
			},
			inParams: &tools.QueryParams{
				Limit:      2,
				Offset:     0,
				OrderBy:    -1,
				OrderField: "age",
			},
			expUsers: []models.User{
				{
					ID:     2,
					Name:   "Kate Jannifer",
					Age:    35,
					About:  "not filtered",
					Gender: "female",
				},
				{
					ID:     1,
					Name:   "Bazil Hotov",
					Age:    22,
					About:  "about",
					Gender: "male",
				},
			},
		},
		{
			name: "OK: Sort age asc",
			inUsers: []models.User{
				{
					ID:     2,
					Name:   "Kate Jannifer",
					Age:    35,
					About:  "not filtered",
					Gender: "female",
				},
				{
					ID:     1,
					Name:   "Bazil Hotov",
					Age:    22,
					About:  "about",
					Gender: "male",
				},
			},
			inParams: &tools.QueryParams{
				Limit:      2,
				Offset:     0,
				OrderBy:    1,
				OrderField: "age",
			},
			expUsers: []models.User{
				{
					ID:     1,
					Name:   "Bazil Hotov",
					Age:    22,
					About:  "about",
					Gender: "male",
				},
				{
					ID:     2,
					Name:   "Kate Jannifer",
					Age:    35,
					About:  "not filtered",
					Gender: "female",
				},
			},
		},
		{
			name: "OK: Sort age asc",
			inUsers: []models.User{
				{
					ID:     1,
					Name:   "Bazil Hotov",
					Age:    22,
					About:  "about",
					Gender: "male",
				},
				{
					ID:     2,
					Name:   "Kate Jannifer",
					Age:    35,
					About:  "not filtered",
					Gender: "female",
				},
			},
			inParams: &tools.QueryParams{
				Limit:      2,
				Offset:     0,
				OrderBy:    0,
				OrderField: "age",
			},
			expUsers: []models.User{
				{
					ID:     1,
					Name:   "Bazil Hotov",
					Age:    22,
					About:  "about",
					Gender: "male",
				},
				{
					ID:     2,
					Name:   "Kate Jannifer",
					Age:    35,
					About:  "not filtered",
					Gender: "female",
				},
			},
		},
		{
			name: "OK: Offset > lenUsers",
			inUsers: []models.User{
				{
					ID:     1,
					Name:   "Bazil Hotov",
					Age:    22,
					About:  "about",
					Gender: "male",
				},
			},
			inParams: &tools.QueryParams{
				Limit:      2,
				Offset:     100,
				OrderBy:    0,
				OrderField: "age",
			},
			expUsers: []models.User{},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userRep := mock_user.NewMockUserRepository(ctrl)
			userUse := NewUserUsecase(userRep)

			users := userUse.SortUsers(testCase.inUsers, testCase.inParams)

			assert.Equal(t, users, testCase.expUsers)
		})
	}
}
