package usecases

import (
	"database/sql"
	"fmt"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/consts"
	customErrors "lectures-2022-1/06_databases/99_hw/redditclone/internal/helpers/errors"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/models"
	mock_user "lectures-2022-1/06_databases/99_hw/redditclone/internal/user/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_RegiserUser(t *testing.T) {
	type mockBehaviour func(userRep *mock_user.MockUserPgRepository, lastID uint64)
	t.Parallel()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inUser        *models.User
		outUserID     uint64
		expError      *customErrors.Error
	}{
		{
			name: "OK",
			mockBehaviour: func(userRep *mock_user.MockUserPgRepository, lastID uint64) {
				userRep.
					EXPECT().
					Insert(gomock.Any()).
					Return(lastID, nil)
			},
			inUser: &models.User{
				Username: "test_user",
				Password: "test_password",
			},
			outUserID: 1,
			expError:  nil,
		},
		{
			name: "Error: CodeUsernameAlreadyTaken",
			mockBehaviour: func(userRep *mock_user.MockUserPgRepository, lastID uint64) {
				userRep.
					EXPECT().
					Insert(gomock.Any()).
					Return(lastID, fmt.Errorf("sql error"))
			},
			inUser: &models.User{
				Username: "test_user",
				Password: "test_password",
			},
			expError: customErrors.Get(consts.CodeUsernameAlreadyTaken),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userRep := mock_user.NewMockUserPgRepository(ctrl)
			testCase.mockBehaviour(userRep, testCase.outUserID)
			userUse := NewUserUsecase(userRep)

			lastID, err := userUse.RegiserUser(testCase.inUser)

			assert.Equal(t, lastID, testCase.outUserID)
			assert.Equal(t, err, testCase.expError)
		})
	}
}

func Test_LoginUser(t *testing.T) {
	type mockBehaviour func(userRep *mock_user.MockUserPgRepository, user *models.User)
	t.Parallel()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inUser        *models.User
		outUserID     uint64
		expError      *customErrors.Error
	}{
		{
			name: "OK",
			mockBehaviour: func(userRep *mock_user.MockUserPgRepository, user *models.User) {
				userRep.
					EXPECT().
					SelectByUsername(gomock.Any()).
					Return(user, nil)
			},
			inUser: &models.User{
				ID:       1,
				Username: "test_user",
				Password: "test_password",
			},
			outUserID: 1,
			expError: nil,
		},
		{
			name: "Error: CodeUserDoesntExist",
			mockBehaviour: func(userRep *mock_user.MockUserPgRepository, user *models.User) {
				userRep.
					EXPECT().
					SelectByUsername(gomock.Any()).
					Return(nil, sql.ErrNoRows)
			},
			inUser: &models.User{
				ID:       1,
				Username: "test_user",
				Password: "test_password",
			},
			outUserID: 0,
			expError: customErrors.Get(consts.CodeUserDoesntExist),
		},
		{
			name: "Error: CodeUsernameAlreadyTaken",
			mockBehaviour: func(userRep *mock_user.MockUserPgRepository, user *models.User) {
				userRep.
					EXPECT().
					SelectByUsername(gomock.Any()).
					Return(nil, fmt.Errorf("sql error"))
			},
			inUser: &models.User{
				ID:       1,
				Username: "test_user",
				Password: "test_password",
			},
			outUserID: 0,
			expError: customErrors.Get(consts.CodeUsernameAlreadyTaken),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			userRep := mock_user.NewMockUserPgRepository(ctrl)
			testCase.mockBehaviour(userRep, testCase.inUser)
			userUse := NewUserUsecase(userRep)

			userID, err := userUse.LoginUser(testCase.inUser)

			assert.Equal(t, userID, testCase.outUserID)
			assert.Equal(t, err, testCase.expError)
		})
	}
}
