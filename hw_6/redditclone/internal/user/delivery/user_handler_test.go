package delivery

import (
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/consts"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/helpers/errors"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/models"
	mock_sess "lectures-2022-1/06_databases/99_hw/redditclone/internal/session/mocks"
	mock_user "lectures-2022-1/06_databases/99_hw/redditclone/internal/user/mocks"
	"lectures-2022-1/06_databases/99_hw/redditclone/tools/converter"
	"lectures-2022-1/06_databases/99_hw/redditclone/tools/response"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Test_RegiserUser(t *testing.T) {
	type mockBehaviouRegistUser func(userUse *mock_user.MockUserUsecase)
	type mockBehaviourCreateSession func(sessUse *mock_sess.MockSessionUsecase, sess *models.Session)

	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	t.Parallel()

	testTable := []struct {
		name                       string
		mockBehaviouRegistUser     mockBehaviouRegistUser
		mockBehaviourCreateSession mockBehaviourCreateSession
		inRequest                  *Request
		outSession                 *models.Session
		expStatusCode              int
		expRespBody                response.Body
	}{
		{
			name: "OK",
			mockBehaviouRegistUser: func(userUse *mock_user.MockUserUsecase) {
				userUse.
					EXPECT().
					RegiserUser(gomock.Any()).
					Return(uint64(0), nil)
			},
			mockBehaviourCreateSession: func(sessUse *mock_sess.MockSessionUsecase, sess *models.Session) {
				sessUse.
					EXPECT().
					Create(gomock.Any()).
					Return(sess, nil)
			},
			inRequest: &Request{
				Username: "testuser",
				Password: "testpassword",
			},
			outSession: &models.Session{
				Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpX",
			},
			expStatusCode: http.StatusCreated,
			expRespBody: response.Body{
				"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpX",
			},
		},
		{
			name: "Error: RegisterUser",
			mockBehaviouRegistUser: func(userUse *mock_user.MockUserUsecase) {
				userUse.
					EXPECT().
					RegiserUser(gomock.Any()).
					Return(uint64(0), errors.Get(consts.CodeInternalError))
			},
			mockBehaviourCreateSession: func(sessUse *mock_sess.MockSessionUsecase, sess *models.Session) {},
			inRequest: &Request{
				Username: "testuser",
				Password: "testpassword",
			},
			expStatusCode: http.StatusInternalServerError,
			expRespBody: response.Body{
				"message": "Что-то пошло не так",
			},
		},
		{
			name: "Error: CreateSession",
			mockBehaviouRegistUser: func(userUse *mock_user.MockUserUsecase) {
				userUse.
					EXPECT().
					RegiserUser(gomock.Any()).
					Return(uint64(0), nil)
			},
			mockBehaviourCreateSession: func(sessUse *mock_sess.MockSessionUsecase, sess *models.Session) {
				sessUse.
					EXPECT().
					Create(gomock.Any()).
					Return(sess, errors.Get(consts.CodeInternalError))
			},
			expStatusCode: http.StatusInternalServerError,
			expRespBody: response.Body{
				"message": "Что-то пошло не так",
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mx := mux.NewRouter()
			r := httptest.NewRequest("POST", "/api/register", converter.AnyBytesToString(testCase.inRequest))
			w := httptest.NewRecorder()
			userUse := mock_user.NewMockUserUsecase(ctrl)
			sessUse := mock_sess.NewMockSessionUsecase(ctrl)
			testCase.mockBehaviouRegistUser(userUse)
			testCase.mockBehaviourCreateSession(sessUse, testCase.outSession)
			userHandler := NewUserHandler(userUse, sessUse)
			userHandler.Configure(mx)

			userHandler.RegiserUser()(w, r)

			expResBody, err := converter.AnyToBytesBuffer(testCase.expRespBody)
			if err != nil {
				t.Error(err.Error())
				return
			}
			bytes := converter.ReadBytes(w.Body)
			assert.Equal(t, testCase.expStatusCode, w.Code)
			assert.JSONEq(t, expResBody.String(), string(bytes))
		})
	}
}

func Test_LoginUser(t *testing.T) {
	type mockBehaviourLoginUser func(userUse *mock_user.MockUserUsecase)
	type mockBehaviourCreateSession func(sessUse *mock_sess.MockSessionUsecase, sess *models.Session)

	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	t.Parallel()

	testTable := []struct {
		name                       string
		mockBehaviouRegistUser     mockBehaviourLoginUser
		mockBehaviourCreateSession mockBehaviourCreateSession
		inRequest                  *Request
		outSession                 *models.Session
		expStatusCode              int
		expRespBody                response.Body
	}{
		{
			name: "OK",
			mockBehaviouRegistUser: func(userUse *mock_user.MockUserUsecase) {
				userUse.
					EXPECT().
					LoginUser(gomock.Any()).
					Return(uint64(0), nil)
			},
			mockBehaviourCreateSession: func(sessUse *mock_sess.MockSessionUsecase, sess *models.Session) {
				sessUse.
					EXPECT().
					Create(gomock.Any()).
					Return(sess, nil)
			},
			inRequest: &Request{
				Username: "testuser",
				Password: "testpassword",
			},
			outSession: &models.Session{
				Value: "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpX",
			},
			expStatusCode: http.StatusCreated,
			expRespBody: response.Body{
				"token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpX",
			},
		},
		{
			name: "Error: RegisterUser",
			mockBehaviouRegistUser: func(userUse *mock_user.MockUserUsecase) {
				userUse.
					EXPECT().
					LoginUser(gomock.Any()).
					Return(uint64(0), errors.Get(consts.CodeInternalError))
			},
			mockBehaviourCreateSession: func(sessUse *mock_sess.MockSessionUsecase, sess *models.Session) {},
			inRequest: &Request{
				Username: "testuser",
				Password: "testpassword",
			},
			expStatusCode: http.StatusInternalServerError,
			expRespBody: response.Body{
				"message": "Что-то пошло не так",
			},
		},
		{
			name: "Error: CreateSession",
			mockBehaviouRegistUser: func(userUse *mock_user.MockUserUsecase) {
				userUse.
					EXPECT().
					LoginUser(gomock.Any()).
					Return(uint64(0), nil)
			},
			mockBehaviourCreateSession: func(sessUse *mock_sess.MockSessionUsecase, sess *models.Session) {
				sessUse.
					EXPECT().
					Create(gomock.Any()).
					Return(sess, errors.Get(consts.CodeInternalError))
			},
			expStatusCode: http.StatusInternalServerError,
			expRespBody: response.Body{
				"message": "Что-то пошло не так",
			},
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			mx := mux.NewRouter()
			r := httptest.NewRequest("POST", "/api/register", converter.AnyBytesToString(testCase.inRequest))
			w := httptest.NewRecorder()
			userUse := mock_user.NewMockUserUsecase(ctrl)
			sessUse := mock_sess.NewMockSessionUsecase(ctrl)
			testCase.mockBehaviouRegistUser(userUse)
			testCase.mockBehaviourCreateSession(sessUse, testCase.outSession)
			userHandler := NewUserHandler(userUse, sessUse)
			userHandler.Configure(mx)

			userHandler.LoginUser()(w, r)

			expResBody, err := converter.AnyToBytesBuffer(testCase.expRespBody)
			if err != nil {
				t.Error(err.Error())
				return
			}
			bytes := converter.ReadBytes(w.Body)
			assert.Equal(t, testCase.expStatusCode, w.Code)
			assert.JSONEq(t, expResBody.String(), string(bytes))
		})
	}
}
