package delivery

import (
	"context"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/consts"
	customContext "lectures-2022-1/06_databases/99_hw/redditclone/internal/helpers/context"
	customErrors "lectures-2022-1/06_databases/99_hw/redditclone/internal/helpers/errors"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/models"
	mock_post "lectures-2022-1/06_databases/99_hw/redditclone/internal/post/mocks"
	"lectures-2022-1/06_databases/99_hw/redditclone/tools/converter"
	"lectures-2022-1/06_databases/99_hw/redditclone/tools/response"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
)

func Test_DeletePost(t *testing.T) {
	type mockBehaviour func(postUse *mock_post.MockPostUsecase)

	t.Parallel()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inPath        string
		inParams      map[string]string
		inUser        *models.User
		expStatusCode int
		expRespBody   response.Body
	}{
		{
			name: "OK",
			mockBehaviour: func(postUse *mock_post.MockPostUsecase) {
				postUse.
					EXPECT().
					DeletePost(gomock.Any(), uint64(1)).
					Return(nil)
			},
			inPath: "/api/post/1",
			inParams: map[string]string{
				"postID": "1",
			},
			inUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			expStatusCode: http.StatusOK,
			expRespBody: response.Body{
				"message": "success",
			},
		},
		{
			name: "Error: CodeInternalError",
			mockBehaviour: func(postUse *mock_post.MockPostUsecase) {
				postUse.
					EXPECT().
					DeletePost(gomock.Any(), uint64(1)).
					Return(customErrors.Get(consts.CodeInternalError))
			},
			inPath: "/api/post/1",
			inParams: map[string]string{
				"postID": "1",
			},
			inUser: &models.User{
				ID:       1,
				Username: "testuser",
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
			r := httptest.NewRequest("DELETE", testCase.inPath, nil)
			r = mux.SetURLVars(r, testCase.inParams)
			w := httptest.NewRecorder()
			ctx := r.Context()
			ctx = context.WithValue(ctx,
				customContext.UserID, testCase.inUser.ID,
			)
			ctx = context.WithValue(ctx,
				customContext.Username, testCase.inUser.Username,
			)
			postUse := mock_post.NewMockPostUsecase(ctrl)

			testCase.mockBehaviour(postUse)
			postHnd := NewPostHandler(postUse)
			postHnd.Configure(mx, nil)

			postHnd.DeletePost()(w, r.WithContext(ctx))
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

func Test_GetPostsByUser(t *testing.T) {
	type mockBehaviour func(postUse *mock_post.MockPostUsecase, posts []*models.Post)

	t.Parallel()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inPath        string
		inParams      map[string]string
		outPosts      []*models.Post
		expStatusCode int
		expRespBody   interface{}
	}{
		{
			name: "OK",
			mockBehaviour: func(postUse *mock_post.MockPostUsecase, posts []*models.Post) {
				postUse.
					EXPECT().
					GetPostsByUser(gomock.Any()).
					Return(posts, nil)
			},
			inPath: "/api/user/testuser",
			inParams: map[string]string{
				"userLogin": "testuser",
			},
			outPosts: []*models.Post{
				{
					UserID:           1,
					Category:         "music",
					CreatedAt:        time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC),
					Score:            0,
					Text:             "text",
					Title:            "title",
					Type:             "type",
					UpvotePercentage: 100,
					URL:              "",
					Views:            0,
				},
				{
					UserID:           2,
					Category:         "fashion",
					CreatedAt:        time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC),
					Score:            0,
					Text:             "text",
					Title:            "title",
					Type:             "type",
					UpvotePercentage: 100,
					URL:              "",
					Views:            0,
				},
			},
			expStatusCode: http.StatusOK,
			expRespBody: []*models.Post{
				{
					UserID:           1,
					Category:         "music",
					CreatedAt:        time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC),
					Score:            0,
					Text:             "text",
					Title:            "title",
					Type:             "type",
					UpvotePercentage: 100,
					URL:              "",
					Views:            0,
				},
				{
					UserID:           2,
					Category:         "fashion",
					CreatedAt:        time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC),
					Score:            0,
					Text:             "text",
					Title:            "title",
					Type:             "type",
					UpvotePercentage: 100,
					URL:              "",
					Views:            0,
				},
			},
		},
		{
			name: "Error: CodeInternalError",
			mockBehaviour: func(postUse *mock_post.MockPostUsecase, posts []*models.Post) {
				postUse.
					EXPECT().
					GetPostsByUser(gomock.Any()).
					Return(nil, customErrors.Get(consts.CodeInternalError))
			},
			inPath: "/api/user/testuser",
			inParams: map[string]string{
				"userLogin": "testuser",
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
			r := httptest.NewRequest("GET", testCase.inPath, nil)
			r = mux.SetURLVars(r, testCase.inParams)
			w := httptest.NewRecorder()
			ctx := r.Context()
			postUse := mock_post.NewMockPostUsecase(ctrl)

			testCase.mockBehaviour(postUse, testCase.outPosts)
			postHnd := NewPostHandler(postUse)
			postHnd.Configure(mx, nil)

			postHnd.GetPostsByUser()(w, r.WithContext(ctx))
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

func Test_DownvotePost(t *testing.T) {
	type mockBehaviour func(postUse *mock_post.MockPostUsecase, post *models.Post)

	t.Parallel()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inPath        string
		inParams      map[string]string
		inUser        *models.User
		outPost       *models.Post
		expStatusCode int
		expRespBody   interface{}
	}{
		{
			name: "OK",
			mockBehaviour: func(postUse *mock_post.MockPostUsecase, post *models.Post) {
				postUse.
					EXPECT().
					DownvotePost(gomock.Any(), uint64(1)).
					Return(post, nil)
			},
			inPath: "/api/post/1/downvote",
			inParams: map[string]string{
				"postID": "1",
			},
			inUser: &models.User{
				ID:       0,
				Username: "testuser",
			},
			outPost: &models.Post{
				UserID:           1,
				Category:         "music",
				CreatedAt:        time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC),
				Score:            0,
				Text:             "text",
				Title:            "title",
				Type:             "type",
				UpvotePercentage: 100,
				URL:              "",
				Views:            0,
			},
			expStatusCode: http.StatusOK,
			expRespBody: &models.Post{
				UserID:           1,
				Category:         "music",
				CreatedAt:        time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC),
				Score:            0,
				Text:             "text",
				Title:            "title",
				Type:             "type",
				UpvotePercentage: 100,
				URL:              "",
				Views:            0,
			},
		},
		{
			name: "Error: CodeInternalError",
			mockBehaviour: func(postUse *mock_post.MockPostUsecase, post *models.Post) {
				postUse.
					EXPECT().
					DownvotePost(gomock.Any(), uint64(1)).
					Return(nil, customErrors.Get(consts.CodeInternalError))
			},
			inPath: "/api/post/1/downvote",
			inParams: map[string]string{
				"postID": "1",
			},
			inUser: &models.User{
				ID:       0,
				Username: "testuser",
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
			r := httptest.NewRequest("GET", testCase.inPath, nil)
			r = mux.SetURLVars(r, testCase.inParams)
			w := httptest.NewRecorder()
			ctx := r.Context()
			ctx = context.WithValue(ctx,
				customContext.UserID, testCase.inUser.ID,
			)
			ctx = context.WithValue(ctx,
				customContext.Username, testCase.inUser.Username,
			)
			postUse := mock_post.NewMockPostUsecase(ctrl)

			testCase.mockBehaviour(postUse, testCase.outPost)
			postHnd := NewPostHandler(postUse)
			postHnd.Configure(mx, nil)

			postHnd.DownvotePost()(w, r.WithContext(ctx))
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

func Test_UnvotePost(t *testing.T) {
	type mockBehaviour func(postUse *mock_post.MockPostUsecase, post *models.Post)

	t.Parallel()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inPath        string
		inParams      map[string]string
		inUser        *models.User
		outPost       *models.Post
		expStatusCode int
		expRespBody   interface{}
	}{
		{
			name: "OK",
			mockBehaviour: func(postUse *mock_post.MockPostUsecase, post *models.Post) {
				postUse.
					EXPECT().
					UnvotePost(gomock.Any(), uint64(1)).
					Return(post, nil)
			},
			inPath: "/api/post/1/downvote",
			inParams: map[string]string{
				"postID": "1",
			},
			inUser: &models.User{
				ID:       0,
				Username: "testuser",
			},
			outPost: &models.Post{
				UserID:           1,
				Category:         "music",
				CreatedAt:        time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC),
				Score:            0,
				Text:             "text",
				Title:            "title",
				Type:             "type",
				UpvotePercentage: 100,
				URL:              "",
				Views:            0,
			},
			expStatusCode: http.StatusOK,
			expRespBody: &models.Post{
				UserID:           1,
				Category:         "music",
				CreatedAt:        time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC),
				Score:            0,
				Text:             "text",
				Title:            "title",
				Type:             "type",
				UpvotePercentage: 100,
				URL:              "",
				Views:            0,
			},
		},
		{
			name: "Error: CodeInternalError",
			mockBehaviour: func(postUse *mock_post.MockPostUsecase, post *models.Post) {
				postUse.
					EXPECT().
					UnvotePost(gomock.Any(), uint64(1)).
					Return(nil, customErrors.Get(consts.CodeInternalError))
			},
			inPath: "/api/post/1/downvote",
			inParams: map[string]string{
				"postID": "1",
			},
			inUser: &models.User{
				ID:       0,
				Username: "testuser",
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
			r := httptest.NewRequest("GET", testCase.inPath, nil)
			r = mux.SetURLVars(r, testCase.inParams)
			w := httptest.NewRecorder()
			ctx := r.Context()
			ctx = context.WithValue(ctx,
				customContext.UserID, testCase.inUser.ID,
			)
			ctx = context.WithValue(ctx,
				customContext.Username, testCase.inUser.Username,
			)
			postUse := mock_post.NewMockPostUsecase(ctrl)

			testCase.mockBehaviour(postUse, testCase.outPost)
			postHnd := NewPostHandler(postUse)
			postHnd.Configure(mx, nil)

			postHnd.UnvotePost()(w, r.WithContext(ctx))
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

func Test_UpvotePost(t *testing.T) {
	type mockBehaviour func(postUse *mock_post.MockPostUsecase, post *models.Post)

	t.Parallel()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inPath        string
		inParams      map[string]string
		inUser        *models.User
		outPost       *models.Post
		expStatusCode int
		expRespBody   interface{}
	}{
		{
			name: "OK",
			mockBehaviour: func(postUse *mock_post.MockPostUsecase, post *models.Post) {
				postUse.
					EXPECT().
					UpvotePost(gomock.Any(), uint64(1)).
					Return(post, nil)
			},
			inPath: "/api/post/1/downvote",
			inParams: map[string]string{
				"postID": "1",
			},
			inUser: &models.User{
				ID:       0,
				Username: "testuser",
			},
			outPost: &models.Post{
				UserID:           1,
				Category:         "music",
				CreatedAt:        time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC),
				Score:            0,
				Text:             "text",
				Title:            "title",
				Type:             "type",
				UpvotePercentage: 100,
				URL:              "",
				Views:            0,
			},
			expStatusCode: http.StatusOK,
			expRespBody: &models.Post{
				UserID:           1,
				Category:         "music",
				CreatedAt:        time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC),
				Score:            0,
				Text:             "text",
				Title:            "title",
				Type:             "type",
				UpvotePercentage: 100,
				URL:              "",
				Views:            0,
			},
		},
		{
			name: "Error: CodeInternalError",
			mockBehaviour: func(postUse *mock_post.MockPostUsecase, post *models.Post) {
				postUse.
					EXPECT().
					UpvotePost(gomock.Any(), uint64(1)).
					Return(nil, customErrors.Get(consts.CodeInternalError))
			},
			inPath: "/api/post/1/downvote",
			inParams: map[string]string{
				"postID": "1",
			},
			inUser: &models.User{
				ID:       0,
				Username: "testuser",
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
			r := httptest.NewRequest("GET", testCase.inPath, nil)
			r = mux.SetURLVars(r, testCase.inParams)
			w := httptest.NewRecorder()
			ctx := r.Context()
			ctx = context.WithValue(ctx,
				customContext.UserID, testCase.inUser.ID,
			)
			ctx = context.WithValue(ctx,
				customContext.Username, testCase.inUser.Username,
			)
			postUse := mock_post.NewMockPostUsecase(ctrl)

			testCase.mockBehaviour(postUse, testCase.outPost)
			postHnd := NewPostHandler(postUse)
			postHnd.Configure(mx, nil)

			postHnd.UpvotePost()(w, r.WithContext(ctx))
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

func Test_DeleteComment(t *testing.T) {
	type mockBehaviour func(postUse *mock_post.MockPostUsecase, post *models.Post)

	t.Parallel()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inPath        string
		inParams      map[string]string
		inUser        *models.User
		outPost       *models.Post
		expStatusCode int
		expRespBody   interface{}
	}{
		{
			name: "OK",
			mockBehaviour: func(postUse *mock_post.MockPostUsecase, post *models.Post) {
				postUse.
					EXPECT().
					DeleteComment(gomock.Any(), uint64(5), uint64(6)).
					Return(post, nil)
			},
			inPath: "/api/post/5/6",
			inParams: map[string]string{
				"postID":    "5",
				"commentID": "6",
			},
			inUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			outPost: &models.Post{
				UserID:           5,
				Category:         "music",
				CreatedAt:        time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC),
				Score:            0,
				Text:             "text",
				Title:            "title",
				Type:             "type",
				UpvotePercentage: 100,
				URL:              "",
				Views:            0,
			},
			expStatusCode: http.StatusOK,
			expRespBody: &models.Post{
				UserID:           5,
				Category:         "music",
				CreatedAt:        time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC),
				Score:            0,
				Text:             "text",
				Title:            "title",
				Type:             "type",
				UpvotePercentage: 100,
				URL:              "",
				Views:            0,
			},
		},
		{
			name: "Error",
			mockBehaviour: func(postUse *mock_post.MockPostUsecase, post *models.Post) {
				postUse.
					EXPECT().
					DeleteComment(gomock.Any(), uint64(5), uint64(6)).
					Return(nil, customErrors.Get(consts.CodeInternalError))
			},
			inPath: "/api/post/5/6",
			inParams: map[string]string{
				"postID":    "5",
				"commentID": "6",
			},
			inUser: &models.User{
				ID:       1,
				Username: "testuser",
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
			r := httptest.NewRequest("DELETE", testCase.inPath, nil)
			r = mux.SetURLVars(r, testCase.inParams)
			w := httptest.NewRecorder()
			ctx := r.Context()
			ctx = context.WithValue(ctx,
				customContext.UserID, testCase.inUser.ID,
			)
			ctx = context.WithValue(ctx,
				customContext.Username, testCase.inUser.Username,
			)
			postUse := mock_post.NewMockPostUsecase(ctrl)

			testCase.mockBehaviour(postUse, testCase.outPost)
			postHnd := NewPostHandler(postUse)
			postHnd.Configure(mx, nil)

			postHnd.DeleteComment()(w, r.WithContext(ctx))
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

func Test_AddComment(t *testing.T) {
	type mockBehaviour func(postUse *mock_post.MockPostUsecase, post *models.Post)

	type Request struct {
		Comment string `json:"comment"`
	}

	t.Parallel()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inPath        string
		inParams      map[string]string
		inRequest     *Request
		inUser        *models.User
		outPost       *models.Post
		expStatusCode int
		expRespBody   interface{}
	}{
		{
			name: "OK",
			mockBehaviour: func(postUse *mock_post.MockPostUsecase, post *models.Post) {
				postUse.
					EXPECT().
					AddComment(gomock.Any()).
					Return(post, nil)
			},
			inPath: "/api/post/1",
			inParams: map[string]string{
				"postID": "1",
			},
			inRequest: &Request{
				Comment: "комментарий",
			},
			inUser: &models.User{
				ID:       0,
				Username: "testuser",
			},
			outPost: &models.Post{
				UserID:           1,
				Category:         "music",
				CreatedAt:        time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC),
				Score:            0,
				Text:             "text",
				Title:            "title",
				Type:             "type",
				UpvotePercentage: 100,
				URL:              "",
				Views:            0,
			},
			expStatusCode: http.StatusOK,
			expRespBody: &models.Post{
				UserID:           1,
				Category:         "music",
				CreatedAt:        time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC),
				Score:            0,
				Text:             "text",
				Title:            "title",
				Type:             "type",
				UpvotePercentage: 100,
				URL:              "",
				Views:            0,
			},
		},
		{
			name: "Error: CodeInternalError",
			mockBehaviour: func(postUse *mock_post.MockPostUsecase, post *models.Post) {
				postUse.
					EXPECT().
					AddComment(gomock.Any()).
					Return(nil, customErrors.Get(consts.CodeInternalError))
			},
			inPath: "/api/post/1",
			inParams: map[string]string{
				"postID": "1",
			},
			inRequest: &Request{
				Comment: "комментарий",
			},
			inUser: &models.User{
				ID:       0,
				Username: "testuser",
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
			r := httptest.NewRequest("POST", testCase.inPath, converter.AnyBytesToString(testCase.inRequest))
			r = mux.SetURLVars(r, testCase.inParams)
			w := httptest.NewRecorder()
			ctx := r.Context()
			ctx = context.WithValue(ctx,
				customContext.UserID, testCase.inUser.ID,
			)
			ctx = context.WithValue(ctx,
				customContext.Username, testCase.inUser.Username,
			)
			postUse := mock_post.NewMockPostUsecase(ctrl)

			testCase.mockBehaviour(postUse, testCase.outPost)
			postHnd := NewPostHandler(postUse)
			postHnd.Configure(mx, nil)

			postHnd.AddComment()(w, r.WithContext(ctx))
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

func Test_GetDetailedPost(t *testing.T) {
	type mockBehaviour func(postUse *mock_post.MockPostUsecase, post *models.Post)

	t.Parallel()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inPath        string
		inParams      map[string]string
		outPost       *models.Post
		expStatusCode int
		expRespBody   interface{}
	}{
		{
			name: "OK",
			mockBehaviour: func(postUse *mock_post.MockPostUsecase, post *models.Post) {
				postUse.
					EXPECT().
					GetPostByID(uint64(1)).
					Return(post, nil)
			},
			inPath: "/api/post/1",
			inParams: map[string]string{
				"postID": "1",
			},
			outPost: &models.Post{
				UserID:           1,
				Category:         "music",
				CreatedAt:        time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC),
				Score:            0,
				Text:             "text",
				Title:            "title",
				Type:             "type",
				UpvotePercentage: 100,
				URL:              "",
				Views:            0,
			},
			expStatusCode: http.StatusOK,
			expRespBody: &models.Post{
				UserID:           1,
				Category:         "music",
				CreatedAt:        time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC),
				Score:            0,
				Text:             "text",
				Title:            "title",
				Type:             "type",
				UpvotePercentage: 100,
				URL:              "",
				Views:            0,
			},
		},
		{
			name: "Error: ",
			mockBehaviour: func(postUse *mock_post.MockPostUsecase, post *models.Post) {
				postUse.
					EXPECT().
					GetPostByID(uint64(1)).
					Return(nil, customErrors.Get(consts.CodeInternalError))
			},
			inPath: "/api/post/1",
			inParams: map[string]string{
				"postID": "1",
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
			r := httptest.NewRequest("GET", testCase.inPath, nil)
			r = mux.SetURLVars(r, testCase.inParams)
			w := httptest.NewRecorder()
			postUse := mock_post.NewMockPostUsecase(ctrl)

			testCase.mockBehaviour(postUse, testCase.outPost)
			postHnd := NewPostHandler(postUse)
			postHnd.Configure(mx, nil)

			postHnd.GetDetailedPost()(w, r)
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

func Test_GetPosts(t *testing.T) {
	type mockBehaviour func(postUse *mock_post.MockPostUsecase, posts []*models.Post)

	t.Parallel()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inPath        string
		outPosts      []*models.Post
		expStatusCode int
		expRespBody   interface{}
	}{
		{
			name: "OK",
			mockBehaviour: func(postUse *mock_post.MockPostUsecase, posts []*models.Post) {
				postUse.
					EXPECT().
					GetPosts().
					Return(posts, nil)
			},
			inPath: "/api/posts/",
			outPosts: []*models.Post{
				{
					UserID:           1,
					Category:         "music",
					CreatedAt:        time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC),
					Score:            0,
					Text:             "text",
					Title:            "title",
					Type:             "type",
					UpvotePercentage: 100,
					URL:              "",
					Views:            0,
				},
				{
					UserID:           2,
					Category:         "fashion",
					CreatedAt:        time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC),
					Score:            0,
					Text:             "text",
					Title:            "title",
					Type:             "type",
					UpvotePercentage: 100,
					URL:              "",
					Views:            0,
				},
			},
			expStatusCode: http.StatusOK,
			expRespBody: []*models.Post{
				{
					UserID:           1,
					Category:         "music",
					CreatedAt:        time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC),
					Score:            0,
					Text:             "text",
					Title:            "title",
					Type:             "type",
					UpvotePercentage: 100,
					URL:              "",
					Views:            0,
				},
				{
					UserID:           2,
					Category:         "fashion",
					CreatedAt:        time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC),
					Score:            0,
					Text:             "text",
					Title:            "title",
					Type:             "type",
					UpvotePercentage: 100,
					URL:              "",
					Views:            0,
				},
			},
		},
		{
			name: "Error: CodeInternalError",
			mockBehaviour: func(postUse *mock_post.MockPostUsecase, posts []*models.Post) {
				postUse.
					EXPECT().
					GetPosts().
					Return(nil, customErrors.Get(consts.CodeInternalError))
			},
			inPath:        "/api/posts/",
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
			r := httptest.NewRequest("GET", testCase.inPath, nil)
			w := httptest.NewRecorder()
			postUse := mock_post.NewMockPostUsecase(ctrl)

			testCase.mockBehaviour(postUse, testCase.outPosts)
			postHnd := NewPostHandler(postUse)
			postHnd.Configure(mx, nil)

			postHnd.GetPosts()(w, r)
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

func Test_GetPostsByCategory(t *testing.T) {
	type mockBehaviour func(postUse *mock_post.MockPostUsecase, categoryPost string, posts []*models.Post)

	t.Parallel()

	testTable := []struct {
		name           string
		mockBehaviour  mockBehaviour
		inPath         string
		inParams       map[string]string
		inPostCategory string
		outPosts       []*models.Post
		expStatusCode  int
		expRespBody    interface{}
	}{
		{
			name: "OK",
			mockBehaviour: func(postUse *mock_post.MockPostUsecase, categoryPost string, posts []*models.Post) {
				postUse.
					EXPECT().
					GetPostsByCategory(categoryPost).
					Return(posts, nil)
			},
			inPath: "/api/posts/music",
			inParams: map[string]string{
				"categoryName": "music",
			},
			inPostCategory: "music",
			outPosts: []*models.Post{
				{
					UserID:           1,
					Category:         "music",
					CreatedAt:        time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC),
					Score:            0,
					Text:             "text",
					Title:            "title",
					Type:             "type",
					UpvotePercentage: 100,
					URL:              "",
					Views:            0,
				},
				{
					UserID:           2,
					Category:         "fashion",
					CreatedAt:        time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC),
					Score:            0,
					Text:             "text",
					Title:            "title",
					Type:             "type",
					UpvotePercentage: 100,
					URL:              "",
					Views:            0,
				},
			},
			expStatusCode: http.StatusOK,
			expRespBody: []*models.Post{
				{
					UserID:           1,
					Category:         "music",
					CreatedAt:        time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC),
					Score:            0,
					Text:             "text",
					Title:            "title",
					Type:             "type",
					UpvotePercentage: 100,
					URL:              "",
					Views:            0,
				},
				{
					UserID:           2,
					Category:         "fashion",
					CreatedAt:        time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC),
					Score:            0,
					Text:             "text",
					Title:            "title",
					Type:             "type",
					UpvotePercentage: 100,
					URL:              "",
					Views:            0,
				},
			},
		},
		{
			name: "Error: CodeInternalError",
			mockBehaviour: func(postUse *mock_post.MockPostUsecase, categoryPost string, posts []*models.Post) {
				postUse.
					EXPECT().
					GetPostsByCategory(categoryPost).
					Return(nil, customErrors.Get(consts.CodeInternalError))
			},
			inPath: "/api/posts/music",
			inParams: map[string]string{
				"categoryName": "music",
			},
			inPostCategory: "music",
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
			r := httptest.NewRequest("GET", testCase.inPath, nil)
			r = mux.SetURLVars(r, testCase.inParams)
			w := httptest.NewRecorder()
			postUse := mock_post.NewMockPostUsecase(ctrl)

			testCase.mockBehaviour(postUse, testCase.inPostCategory, testCase.outPosts)
			postHnd := NewPostHandler(postUse)
			postHnd.Configure(mx, nil)

			postHnd.GetPostsByCategory()(w, r)
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

func Test_AddPost(t *testing.T) {
	type mockBehaviour func(postUse *mock_post.MockPostUsecase, post *models.Post)

	type Request struct {
		Category string `json:"category"`
		Text     string `json:"text,omitempty"`
		Title    string `json:"title"`
		Type     string `json:"type"`
		URL      string `json:"url,omitempty"`
	}

	t.Parallel()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inPath        string
		inRequest     *Request
		inUser        *models.User
		outPost       *models.Post
		expStatusCode int
		expRespBody   interface{}
	}{
		{
			name: "OK",
			mockBehaviour: func(postUse *mock_post.MockPostUsecase, post *models.Post) {
				postUse.
					EXPECT().
					CreatePost(gomock.Any(), gomock.Any()).
					Return(post, nil)
			},
			inPath: "/api/posts",
			inRequest: &Request{
				Category: "music",
				Text: "text",
				Title: "title",
				Type: "type",
				URL: "",
			},
			inUser: &models.User{
				ID:       0,
				Username: "testuser",
			},
			outPost: &models.Post{
				UserID:           1,
				Category:         "music",
				CreatedAt:        time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC),
				Score:            0,
				Text:             "text",
				Title:            "title",
				Type:             "type",
				UpvotePercentage: 100,
				URL:              "",
				Views:            0,
			},
			expStatusCode: http.StatusCreated,
			expRespBody: &models.Post{
				UserID:           1,
				Category:         "music",
				CreatedAt:        time.Date(2021, time.May, 10, 23, 0, 0, 0, time.UTC),
				Score:            0,
				Text:             "text",
				Title:            "title",
				Type:             "type",
				UpvotePercentage: 100,
				URL:              "",
				Views:            0,
			},
		},
		{
			name: "Error",
			mockBehaviour: func(postUse *mock_post.MockPostUsecase, post *models.Post) {
				postUse.
					EXPECT().
					CreatePost(gomock.Any(), gomock.Any()).
					Return(nil, customErrors.Get(consts.CodeInternalError))
			},
			inPath: "/api/posts",
			inRequest: &Request{
				Category: "music",
				Text: "text",
				Title: "title",
				Type: "type",
				URL: "",
			},
			inUser: &models.User{
				ID:       0,
				Username: "testuser",
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
			r := httptest.NewRequest("POST", testCase.inPath, converter.AnyBytesToString(testCase.inRequest))
			w := httptest.NewRecorder()
			ctx := r.Context()
			ctx = context.WithValue(ctx,
				customContext.UserID, testCase.inUser.ID,
			)
			ctx = context.WithValue(ctx,
				customContext.Username, testCase.inUser.Username,
			)
			postUse := mock_post.NewMockPostUsecase(ctrl)

			testCase.mockBehaviour(postUse, testCase.outPost)
			postHnd := NewPostHandler(postUse)
			postHnd.Configure(mx, nil)

			postHnd.AddPost()(w, r.WithContext(ctx))
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