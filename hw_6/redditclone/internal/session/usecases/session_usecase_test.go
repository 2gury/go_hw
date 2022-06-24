package usecases

import (
	"fmt"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/consts"
	customErrors "lectures-2022-1/06_databases/99_hw/redditclone/internal/helpers/errors"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/models"
	mock_sess "lectures-2022-1/06_databases/99_hw/redditclone/internal/session/mocks"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/tkuchiki/faketime"
)

func Test_Create(t *testing.T) {
	type mockBehaviour func(sessRep *mock_sess.MockSessionRepository)
	// t.Parallel()
	jwt.TimeFunc = func() time.Time {
				return time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC)
			}
	mockTime := faketime.NewFaketime(2021, time.May, 10, 23, 0, 0, 0, time.UTC)
	defer mockTime.Undo()
	mockTime.Do()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inUser        *models.User
		outSess       *models.Session
		expError      *customErrors.Error
	}{
		{
			name: "OK",
			mockBehaviour: func(sessRep *mock_sess.MockSessionRepository) {
				sessRep.
					EXPECT().
					Create(gomock.Any()).
					Return(nil)
			},
			inUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			outSess: &models.Session{
				Value:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjA3MjM2MDAsImlhdCI6MTYyMDY4NzYwMCwidXNlciI6eyJpZCI6IjEiLCJ1c2VybmFtZSI6InRlc3R1c2VyIn19.LfH8WMJPK3P52fYutI0bn3t2n98vnyVb0fx7Cjq2qkA",
				TimeDuration: consts.ExpiresDuration,
			},
			expError: nil,
		},
		{
			name: "Error: CodeInternalError",
			mockBehaviour: func(sessRep *mock_sess.MockSessionRepository) {
				sessRep.
					EXPECT().
					Create(gomock.Any()).
					Return(fmt.Errorf("redis error"))
			},
			inUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			expError: customErrors.Get(consts.CodeInternalError),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			sessRep := mock_sess.NewMockSessionRepository(ctrl)
			testCase.mockBehaviour(sessRep)
			sessUse := NewSessionUsecase(sessRep)

			sess, err := sessUse.Create(testCase.inUser)

			assert.Equal(t, sess, testCase.outSess)
			assert.Equal(t, err, testCase.expError)
		})
	}
}

func Test_Check(t *testing.T) {
	type mockBehaviour func(sessRep *mock_sess.MockSessionRepository, sess *models.Session)
	t.Parallel()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inSessValue   string
		outSess       *models.Session
		outUser       *models.User
		expError      *customErrors.Error
	}{
		{
			name: "OK",
			mockBehaviour: func(sessRep *mock_sess.MockSessionRepository, sess *models.Session) {
				sessRep.
					EXPECT().
					Get(gomock.Any()).
					Return(sess, nil)
			},
			inSessValue: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjA3MjM2MDAsImlhdCI6MTYyMDY4NzYwMCwidXNlciI6eyJpZCI6IjEiLCJ1c2VybmFtZSI6InRlc3R1c2VyIn19.LfH8WMJPK3P52fYutI0bn3t2n98vnyVb0fx7Cjq2qkA",
			outSess: &models.Session{
				Value:        "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjA3MjM2MDAsImlhdCI6MTYyMDY4NzYwMCwidXNlciI6eyJpZCI6IjEiLCJ1c2VybmFtZSI6InRlc3R1c2VyIn19.LfH8WMJPK3P52fYutI0bn3t2n98vnyVb0fx7Cjq2qkA",
				TimeDuration: consts.ExpiresDuration,
			},
			outUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			expError: nil,
		},
		{
			name: "Error: CodeInternalError",
			mockBehaviour: func(sessRep *mock_sess.MockSessionRepository, sess *models.Session) {
				sessRep.
					EXPECT().
					Get(gomock.Any()).
					Return(nil, fmt.Errorf("redis error"))
			},
			inSessValue: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjA3MjM2MDAsImlhdCI6MTYyMDY4NzYwMCwidXNlciI6eyJpZCI6IjEiLCJ1c2VybmFtZSI6InRlc3R1c2VyIn19.LfH8WMJPK3P52fYutI0bn3t2n98vnyVb0fx7Cjq2qkA",
			expError:    customErrors.Get(consts.CodeInternalError),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			jwt.TimeFunc = func() time.Time {
				return time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC)
			}
			mockTime := faketime.NewFaketime(2021, time.May, 10, 23, 0, 0, 0, time.UTC)
			defer mockTime.Undo()
			mockTime.Do()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			sessRep := mock_sess.NewMockSessionRepository(ctrl)
			testCase.mockBehaviour(sessRep, testCase.outSess)
			sessUse := NewSessionUsecase(sessRep)

			user, err := sessUse.Check(testCase.inSessValue)
			assert.Equal(t, user, testCase.outUser)
			assert.Equal(t, err, testCase.expError)
		})
	}
}

func Test_Delete(t *testing.T) {
	type mockBehaviour func(sessRep *mock_sess.MockSessionRepository)
	t.Parallel()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inSessValue   string
		expError      *customErrors.Error
	}{
		{
			name: "OK",
			mockBehaviour: func(sessRep *mock_sess.MockSessionRepository) {
				sessRep.
					EXPECT().
					Delete(gomock.Any()).
					Return(nil)
			},
			inSessValue: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjA3MjM2MDAsImlhdCI6MTYyMDY4NzYwMCwidXNlciI6eyJpZCI6IjEiLCJ1c2VybmFtZSI6InRlc3R1c2VyIn19.LfH8WMJPK3P52fYutI0bn3t2n98vnyVb0fx7Cjq2qkA",
			expError:    nil,
		},
		{
			name: "Error: CodeInternalError",
			mockBehaviour: func(sessRep *mock_sess.MockSessionRepository) {
				sessRep.
					EXPECT().
					Delete(gomock.Any()).
					Return(fmt.Errorf("redis error"))
			},
			inSessValue: "eyJhbGciOiIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjA3MjM2MDAsImlhdCI6MTYyMDY4NzYwMCwidXNlciI6eyJpZCI6IjEiLCJ1c2VybmFtZSI6InRlc3R1c2VyIn19.LfH8WMJPK3P52fYutI0bn3t2n98vnyVb0fx7Cjq2qkA",
			expError:    customErrors.Get(consts.CodeInternalError),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			jwt.TimeFunc = func() time.Time {
				return time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC)
			}
			mockTime := faketime.NewFaketime(2021, time.May, 10, 23, 0, 0, 0, time.UTC)
			defer mockTime.Undo()
			mockTime.Do()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			sessRep := mock_sess.NewMockSessionRepository(ctrl)
			testCase.mockBehaviour(sessRep)
			sessUse := NewSessionUsecase(sessRep)

			err := sessUse.Delete(testCase.inSessValue)
			assert.Equal(t, err, testCase.expError)
		})
	}
}

func Test_NewJwtSession(t *testing.T) {
	// t.Parallel()
	jwt.TimeFunc = func() time.Time {
		return time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC)
	}
	mockTime := faketime.NewFaketime(2021, time.May, 10, 23, 0, 0, 0, time.UTC)
	defer mockTime.Undo()
	mockTime.Do()

	testTable := []struct {
		name         string
		inUser       *models.User
		outSessValue string
		expError     error
	}{
		{
			name: "OK",
			inUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			outSessValue: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2MjA3MjM2MDAsImlhdCI6MTYyMDY4NzYwMCwidXNlciI6eyJpZCI6IjEiLCJ1c2VybmFtZSI6InRlc3R1c2VyIn19.LfH8WMJPK3P52fYutI0bn3t2n98vnyVb0fx7Cjq2qkA",
			expError:     nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			sessRep := mock_sess.NewMockSessionRepository(ctrl)
			sessUse := &SessionUsecase{
				sessRep: sessRep,
			}

			sessValue, err := sessUse.NewJwtSession(testCase.inUser)
			assert.Equal(t, sessValue, testCase.outSessValue)
			assert.Equal(t, err, testCase.expError)
		})
	}
}
