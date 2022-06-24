package usecases

import (
	"fmt"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/consts"
	customErrors "lectures-2022-1/06_databases/99_hw/redditclone/internal/helpers/errors"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/models"
	mock_post "lectures-2022-1/06_databases/99_hw/redditclone/internal/post/mocks"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func Test_CreatePost(t *testing.T) {
	type mockBehaviourInsertPost func(postRep *mock_post.MockPostRepository)
	type mockBehaviourInsertVote func(postRep *mock_post.MockPostRepository)
	type mockBehaviourSelectPostByID func(postRep *mock_post.MockPostRepository, post *models.Post)
	type mockBehaviourSelectCommentsByPostID func(postRep *mock_post.MockPostRepository, comments []*models.Comment)
	type mockBehaviourSelectVotesByPostID func(postRep *mock_post.MockPostRepository, votes []*models.Vote)
	t.Parallel()

	testTable := []struct {
		name                                string
		mockBehaviourInsertPost             mockBehaviourInsertPost
		mockBehaviourInsertVote             mockBehaviourInsertVote
		mockBehaviourSelectPostByID         mockBehaviourSelectPostByID
		mockBehaviourSelectCommentsByPostID mockBehaviourSelectCommentsByPostID
		mockBehaviourSelectVotesByPostID    mockBehaviourSelectVotesByPostID
		inUser                              *models.User
		inPost                              *models.Post
		outPost                             *models.Post
		outComments                         []*models.Comment
		outVotes                            []*models.Vote
		expPost                             *models.Post
		expError                            *customErrors.Error
	}{
		{
			name: "OK",
			mockBehaviourInsertPost: func(postRep *mock_post.MockPostRepository) {
				postRep.
					EXPECT().
					InsertPost(gomock.Any()).
					Return(uint64(0), nil)
			},
			mockBehaviourInsertVote: func(postRep *mock_post.MockPostRepository) {
				postRep.
					EXPECT().
					InsertVote(gomock.Any()).
					Return(uint64(0), nil)
			},
			mockBehaviourSelectPostByID: func(postRep *mock_post.MockPostRepository, post *models.Post) {
				postRep.
					EXPECT().
					SelectPostByID(gomock.Any()).
					Return(post, nil)
			},
			mockBehaviourSelectCommentsByPostID: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {
				postRep.
					EXPECT().
					SelectCommentsByPostID(gomock.Any()).
					Return(comments, nil)
			},
			mockBehaviourSelectVotesByPostID: func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {
				postRep.
					EXPECT().
					SelectVotesByPostID(gomock.Any()).
					Return(votes, nil)
			},
			inUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			inPost: &models.Post{
				ID:    1,
				Title: "title",
				Text:  "text",
			},
			outPost: &models.Post{
				ID:    1,
				Title: "title",
				Text:  "text",
			},
			outComments: []*models.Comment{
				{
					ID: 3,
				},
			},
			outVotes: []*models.Vote{
				{
					ID: 5,
				},
			},
			expPost: &models.Post{
				ID:    1,
				Title: "title",
				Text:  "text",
				Comments: []*models.Comment{
					{
						ID: 3,
					},
				},
				Votes: []*models.Vote{
					{
						ID: 5,
					},
				},
			},
			expError: nil,
		},
		{
			name: "Error: InsertPost",
			mockBehaviourInsertPost: func(postRep *mock_post.MockPostRepository) {
				postRep.
					EXPECT().
					InsertPost(gomock.Any()).
					Return(uint64(0), fmt.Errorf("sql error"))
			},
			mockBehaviourInsertVote:             func(postRep *mock_post.MockPostRepository) {},
			mockBehaviourSelectPostByID:         func(postRep *mock_post.MockPostRepository, post *models.Post) {},
			mockBehaviourSelectCommentsByPostID: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {},
			mockBehaviourSelectVotesByPostID:    func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {},
			inUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			inPost: &models.Post{
				ID:    1,
				Title: "title",
				Text:  "text",
			},
			expError: customErrors.Get(consts.CodeInternalError),
		},
		{
			name: "Error: InsertVote",
			mockBehaviourInsertPost: func(postRep *mock_post.MockPostRepository) {
				postRep.
					EXPECT().
					InsertPost(gomock.Any()).
					Return(uint64(0), nil)
			},
			mockBehaviourInsertVote: func(postRep *mock_post.MockPostRepository) {
				postRep.
					EXPECT().
					InsertVote(gomock.Any()).
					Return(uint64(0), fmt.Errorf("sql error"))
			},
			mockBehaviourSelectPostByID:         func(postRep *mock_post.MockPostRepository, post *models.Post) {},
			mockBehaviourSelectCommentsByPostID: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {},
			mockBehaviourSelectVotesByPostID:    func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {},
			inUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			inPost: &models.Post{
				ID:    1,
				Title: "title",
				Text:  "text",
			},
			expError: customErrors.Get(consts.CodeInternalError),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			postRep := mock_post.NewMockPostRepository(ctrl)
			testCase.mockBehaviourInsertPost(postRep)
			testCase.mockBehaviourInsertVote(postRep)
			testCase.mockBehaviourSelectPostByID(postRep, testCase.outPost)
			testCase.mockBehaviourSelectCommentsByPostID(postRep, testCase.outComments)
			testCase.mockBehaviourSelectVotesByPostID(postRep, testCase.outVotes)
			postUse := NewPostUsecase(postRep)

			post, err := postUse.CreatePost(testCase.inUser, testCase.inPost)
			assert.Equal(t, post, testCase.expPost)
			assert.Equal(t, err, testCase.expError)
		})
	}
}

func Test_GetPosts(t *testing.T) {
	type mockBehaviourSelectAllPosts func(postRep *mock_post.MockPostRepository, posts []*models.Post)
	type mockBehaviourSelectAllComments func(postRep *mock_post.MockPostRepository, comments []*models.Comment)
	type mockBehaviourSelectAllVotes func(postRep *mock_post.MockPostRepository, votes []*models.Vote)
	t.Parallel()

	testTable := []struct {
		name                           string
		mockBehaviourSelectAllPosts    mockBehaviourSelectAllPosts
		mockBehaviourSelectAllComments mockBehaviourSelectAllComments
		mockBehaviourSelectAllVotes    mockBehaviourSelectAllVotes
		outPosts                       []*models.Post
		outComments                    []*models.Comment
		outVotes                       []*models.Vote
		expPosts                       []*models.Post
		expError                       *customErrors.Error
	}{
		{
			name: "OK",
			mockBehaviourSelectAllPosts: func(postRep *mock_post.MockPostRepository, posts []*models.Post) {
				postRep.
					EXPECT().
					SelectAllPosts().
					Return(posts, nil)
			},
			mockBehaviourSelectAllComments: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {
				postRep.
					EXPECT().
					SelectAllComments().
					Return(comments, nil)
			},
			mockBehaviourSelectAllVotes: func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {
				postRep.
					EXPECT().
					SelectAllVotes().
					Return(votes, nil)
			},
			outPosts: []*models.Post{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			},
			outComments: []*models.Comment{
				{
					ID:     7,
					PostID: 2,
				},
				{
					ID:     8,
					PostID: 1,
				},
			},
			outVotes: []*models.Vote{
				{
					ID:     5,
					PostID: 1,
				},
				{
					ID:     6,
					PostID: 2,
				},
			},
			expPosts: []*models.Post{
				{
					ID: 1,
					Comments: []*models.Comment{
						{
							ID:     8,
							PostID: 1,
						},
					},
					Votes: []*models.Vote{
						{
							ID:     5,
							PostID: 1,
						},
					},
				},
				{
					ID: 2,
					Comments: []*models.Comment{
						{
							ID:     7,
							PostID: 2,
						},
					},
					Votes: []*models.Vote{
						{
							ID:     6,
							PostID: 2,
						},
					},
				},
			},
			expError: nil,
		},
		{
			name: "Error: SelectAllVotes",
			mockBehaviourSelectAllPosts: func(postRep *mock_post.MockPostRepository, posts []*models.Post) {
				postRep.
					EXPECT().
					SelectAllPosts().
					Return(posts, nil)
			},
			mockBehaviourSelectAllComments: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {
				postRep.
					EXPECT().
					SelectAllComments().
					Return(comments, nil)
			},
			mockBehaviourSelectAllVotes: func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {
				postRep.
					EXPECT().
					SelectAllVotes().
					Return(nil, fmt.Errorf("sql error"))
			},
			outPosts: []*models.Post{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			},
			outComments: []*models.Comment{
				{
					ID:     7,
					PostID: 2,
				},
				{
					ID:     8,
					PostID: 1,
				},
			},
			expError: customErrors.Get(consts.CodeInternalError),
		},
		{
			name: "Error: SelectAllComments:",
			mockBehaviourSelectAllPosts: func(postRep *mock_post.MockPostRepository, posts []*models.Post) {
				postRep.
					EXPECT().
					SelectAllPosts().
					Return(posts, nil)
			},
			mockBehaviourSelectAllComments: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {
				postRep.
					EXPECT().
					SelectAllComments().
					Return(nil, fmt.Errorf("sql error"))
			},
			mockBehaviourSelectAllVotes: func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {},
			outPosts: []*models.Post{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			},
			expError: customErrors.Get(consts.CodeInternalError),
		},
		{
			name: "Error: SelectAllPosts:",
			mockBehaviourSelectAllPosts: func(postRep *mock_post.MockPostRepository, posts []*models.Post) {
				postRep.
					EXPECT().
					SelectAllPosts().
					Return(nil, fmt.Errorf("sql error"))
			},
			mockBehaviourSelectAllComments: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {},
			mockBehaviourSelectAllVotes:    func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {},
			expError:                       customErrors.Get(consts.CodeInternalError),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			postRep := mock_post.NewMockPostRepository(ctrl)
			testCase.mockBehaviourSelectAllPosts(postRep, testCase.outPosts)
			testCase.mockBehaviourSelectAllComments(postRep, testCase.outComments)
			testCase.mockBehaviourSelectAllVotes(postRep, testCase.outVotes)
			postUse := NewPostUsecase(postRep)

			posts, err := postUse.GetPosts()
			assert.Equal(t, posts, testCase.expPosts)
			assert.Equal(t, err, testCase.expError)
		})
	}
}

func Test_GetPostsByCategory(t *testing.T) {
	type mockBehaviourSelectPostsByCategory func(postRep *mock_post.MockPostRepository, posts []*models.Post)
	type mockBehaviourSelectCommentsByCategory func(postRep *mock_post.MockPostRepository, comments []*models.Comment)
	type mockBehaviourSelectVotesByCategory func(postRep *mock_post.MockPostRepository, votes []*models.Vote)
	t.Parallel()

	testTable := []struct {
		name                                  string
		mockBehaviourSelectPostsByCategory    mockBehaviourSelectPostsByCategory
		mockBehaviourSelectCommentsByCategory mockBehaviourSelectCommentsByCategory
		mockBehaviourSelectVotesByCategory    mockBehaviourSelectVotesByCategory
		inPostCategory                        string
		outPosts                              []*models.Post
		outComments                           []*models.Comment
		outVotes                              []*models.Vote
		expPosts                              []*models.Post
		expError                              *customErrors.Error
	}{
		{
			name: "OK",
			mockBehaviourSelectPostsByCategory: func(postRep *mock_post.MockPostRepository, posts []*models.Post) {
				postRep.
					EXPECT().
					SelectPostsByCategory(gomock.Any()).
					Return(posts, nil)
			},
			mockBehaviourSelectCommentsByCategory: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {
				postRep.
					EXPECT().
					SelectCommentsByCategory(gomock.Any()).
					Return(comments, nil)
			},
			mockBehaviourSelectVotesByCategory: func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {
				postRep.
					EXPECT().
					SelectVotesByCategory(gomock.Any()).
					Return(votes, nil)
			},
			inPostCategory: "music",
			outPosts: []*models.Post{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			},
			outComments: []*models.Comment{
				{
					ID:     7,
					PostID: 2,
				},
				{
					ID:     8,
					PostID: 1,
				},
			},
			outVotes: []*models.Vote{
				{
					ID:     5,
					PostID: 1,
				},
				{
					ID:     6,
					PostID: 2,
				},
			},
			expPosts: []*models.Post{
				{
					ID: 1,
					Comments: []*models.Comment{
						{
							ID:     8,
							PostID: 1,
						},
					},
					Votes: []*models.Vote{
						{
							ID:     5,
							PostID: 1,
						},
					},
				},
				{
					ID: 2,
					Comments: []*models.Comment{
						{
							ID:     7,
							PostID: 2,
						},
					},
					Votes: []*models.Vote{
						{
							ID:     6,
							PostID: 2,
						},
					},
				},
			},
			expError: nil,
		},
		{
			name: "Error: SelectAllVotes",
			mockBehaviourSelectPostsByCategory: func(postRep *mock_post.MockPostRepository, posts []*models.Post) {
				postRep.
					EXPECT().
					SelectPostsByCategory(gomock.Any()).
					Return(posts, nil)
			},
			mockBehaviourSelectCommentsByCategory: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {
				postRep.
					EXPECT().
					SelectCommentsByCategory(gomock.Any()).
					Return(comments, nil)
			},
			mockBehaviourSelectVotesByCategory: func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {
				postRep.
					EXPECT().
					SelectVotesByCategory(gomock.Any()).
					Return(nil, fmt.Errorf("sql error"))
			},
			inPostCategory: "music",
			outPosts: []*models.Post{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			},
			outComments: []*models.Comment{
				{
					ID:     7,
					PostID: 2,
				},
				{
					ID:     8,
					PostID: 1,
				},
			},
			expError: customErrors.Get(consts.CodeInternalError),
		},
		{
			name: "Error: SelectCommentsByCategory:",
			mockBehaviourSelectPostsByCategory: func(postRep *mock_post.MockPostRepository, posts []*models.Post) {
				postRep.
					EXPECT().
					SelectPostsByCategory(gomock.Any()).
					Return(posts, nil)
			},
			mockBehaviourSelectCommentsByCategory: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {
				postRep.
					EXPECT().
					SelectCommentsByCategory(gomock.Any()).
					Return(nil, fmt.Errorf("sql error"))
			},
			mockBehaviourSelectVotesByCategory: func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {},
			inPostCategory:                     "music",
			outPosts: []*models.Post{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			},
			expError: customErrors.Get(consts.CodeInternalError),
		},
		{
			name: "Error: SelectAllPosts:",
			mockBehaviourSelectPostsByCategory: func(postRep *mock_post.MockPostRepository, posts []*models.Post) {
				postRep.
					EXPECT().
					SelectPostsByCategory(gomock.Any()).
					Return(nil, fmt.Errorf("sql error"))
			},
			mockBehaviourSelectCommentsByCategory: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {},
			mockBehaviourSelectVotesByCategory:    func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {},
			inPostCategory:                        "music",
			expError:                              customErrors.Get(consts.CodeInternalError),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			postRep := mock_post.NewMockPostRepository(ctrl)
			testCase.mockBehaviourSelectPostsByCategory(postRep, testCase.outPosts)
			testCase.mockBehaviourSelectCommentsByCategory(postRep, testCase.outComments)
			testCase.mockBehaviourSelectVotesByCategory(postRep, testCase.outVotes)
			postUse := NewPostUsecase(postRep)

			posts, err := postUse.GetPostsByCategory(testCase.inPostCategory)
			assert.Equal(t, posts, testCase.expPosts)
			assert.Equal(t, err, testCase.expError)
		})
	}
}

func Test_GetPostByID(t *testing.T) {
	type mockBehaviourSelectPostByID func(postRep *mock_post.MockPostRepository, post *models.Post)
	type mockBehaviourSelectCommentsByPostID func(postRep *mock_post.MockPostRepository, comments []*models.Comment)
	type mockBehaviourSelectVotesByPostID func(postRep *mock_post.MockPostRepository, votes []*models.Vote)
	t.Parallel()

	testTable := []struct {
		name                                string
		mockBehaviourSelectPostByID         mockBehaviourSelectPostByID
		mockBehaviourSelectCommentsByPostID mockBehaviourSelectCommentsByPostID
		mockBehaviourSelectVotesByPostID    mockBehaviourSelectVotesByPostID
		inPostID                            uint64
		outPost                             *models.Post
		outComments                         []*models.Comment
		outVotes                            []*models.Vote
		expPost                             *models.Post
		expError                            *customErrors.Error
	}{
		{
			name: "OK",
			mockBehaviourSelectPostByID: func(postRep *mock_post.MockPostRepository, post *models.Post) {
				postRep.
					EXPECT().
					SelectPostByID(gomock.Any()).
					Return(post, nil)
			},
			mockBehaviourSelectCommentsByPostID: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {
				postRep.
					EXPECT().
					SelectCommentsByPostID(gomock.Any()).
					Return(comments, nil)
			},
			mockBehaviourSelectVotesByPostID: func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {
				postRep.
					EXPECT().
					SelectVotesByPostID(gomock.Any()).
					Return(votes, nil)
			},
			inPostID: 1,
			outPost: &models.Post{
				ID: 1,
			},
			outComments: []*models.Comment{
				{
					ID: 3,
				},
			},
			outVotes: []*models.Vote{
				{
					ID: 5,
				},
			},
			expPost: &models.Post{
				ID: 1,
				Comments: []*models.Comment{
					{
						ID: 3,
					},
				},
				Votes: []*models.Vote{
					{
						ID: 5,
					},
				},
			},
			expError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			postRep := mock_post.NewMockPostRepository(ctrl)
			testCase.mockBehaviourSelectPostByID(postRep, testCase.outPost)
			testCase.mockBehaviourSelectCommentsByPostID(postRep, testCase.outComments)
			testCase.mockBehaviourSelectVotesByPostID(postRep, testCase.outVotes)
			postUse := NewPostUsecase(postRep)

			post, err := postUse.GetPostByID(testCase.inPostID)
			assert.Equal(t, post, testCase.expPost)
			assert.Equal(t, err, testCase.expError)
		})
	}
}

func Test_AddComment(t *testing.T) {
	type mockBehaviourInsertComment func(postRep *mock_post.MockPostRepository)
	type mockBehaviourSelectPostByID func(postRep *mock_post.MockPostRepository, post *models.Post)
	type mockBehaviourSelectCommentsByPostID func(postRep *mock_post.MockPostRepository, comments []*models.Comment)
	type mockBehaviourSelectVotesByPostID func(postRep *mock_post.MockPostRepository, votes []*models.Vote)
	t.Parallel()

	testTable := []struct {
		name                                string
		mockBehaviourInsertComment          mockBehaviourInsertComment
		mockBehaviourSelectPostByID         mockBehaviourSelectPostByID
		mockBehaviourSelectCommentsByPostID mockBehaviourSelectCommentsByPostID
		mockBehaviourSelectVotesByPostID    mockBehaviourSelectVotesByPostID
		inComment                           *models.Comment
		outPost                             *models.Post
		outComments                         []*models.Comment
		outVotes                            []*models.Vote
		expPost                             *models.Post
		expError                            *customErrors.Error
	}{
		{
			name: "OK",
			mockBehaviourInsertComment: func(postRep *mock_post.MockPostRepository) {
				postRep.
					EXPECT().
					InsertComment(gomock.Any()).
					Return(uint64(0), nil)
			},
			mockBehaviourSelectPostByID: func(postRep *mock_post.MockPostRepository, post *models.Post) {
				postRep.
					EXPECT().
					SelectPostByID(gomock.Any()).
					Return(post, nil)
			},
			mockBehaviourSelectCommentsByPostID: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {
				postRep.
					EXPECT().
					SelectCommentsByPostID(gomock.Any()).
					Return(comments, nil)
			},
			mockBehaviourSelectVotesByPostID: func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {
				postRep.
					EXPECT().
					SelectVotesByPostID(gomock.Any()).
					Return(votes, nil)
			},
			inComment: &models.Comment{
				ID: 8,
			},
			outPost: &models.Post{
				ID:    1,
				Title: "title",
				Text:  "text",
			},
			outComments: []*models.Comment{
				{
					ID: 3,
				},
			},
			outVotes: []*models.Vote{
				{
					ID: 5,
				},
			},
			expPost: &models.Post{
				ID:    1,
				Title: "title",
				Text:  "text",
				Comments: []*models.Comment{
					{
						ID: 3,
					},
				},
				Votes: []*models.Vote{
					{
						ID: 5,
					},
				},
			},
			expError: nil,
		},
		{
			name: "Error: InsertComment",
			mockBehaviourInsertComment: func(postRep *mock_post.MockPostRepository) {
				postRep.
					EXPECT().
					InsertComment(gomock.Any()).
					Return(uint64(0), fmt.Errorf("sql error"))
			},
			mockBehaviourSelectPostByID:         func(postRep *mock_post.MockPostRepository, post *models.Post) {},
			mockBehaviourSelectCommentsByPostID: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {},
			mockBehaviourSelectVotesByPostID:    func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {},
			inComment: &models.Comment{
				ID: 1,
			},
			expError: customErrors.Get(consts.CodeInternalError),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			postRep := mock_post.NewMockPostRepository(ctrl)
			testCase.mockBehaviourInsertComment(postRep)
			testCase.mockBehaviourSelectPostByID(postRep, testCase.outPost)
			testCase.mockBehaviourSelectCommentsByPostID(postRep, testCase.outComments)
			testCase.mockBehaviourSelectVotesByPostID(postRep, testCase.outVotes)
			postUse := NewPostUsecase(postRep)

			post, err := postUse.AddComment(testCase.inComment)
			assert.Equal(t, post, testCase.expPost)
			assert.Equal(t, err, testCase.expError)
		})
	}
}

func Test_DeleteComment(t *testing.T) {
	type mockBehaviourSelectAuthorComment func(postRep *mock_post.MockPostRepository, user *models.User)
	type mockBehaviourDeleteCommentByID func(postRep *mock_post.MockPostRepository)
	type mockBehaviourSelectPostByID func(postRep *mock_post.MockPostRepository, post *models.Post)
	type mockBehaviourSelectCommentsByPostID func(postRep *mock_post.MockPostRepository, comments []*models.Comment)
	type mockBehaviourSelectVotesByPostID func(postRep *mock_post.MockPostRepository, votes []*models.Vote)
	t.Parallel()

	testTable := []struct {
		name                                string
		mockBehaviourSelectAuthorComment    mockBehaviourSelectAuthorComment
		mockBehaviourDeleteCommentByID      mockBehaviourDeleteCommentByID
		mockBehaviourSelectPostByID         mockBehaviourSelectPostByID
		mockBehaviourSelectCommentsByPostID mockBehaviourSelectCommentsByPostID
		mockBehaviourSelectVotesByPostID    mockBehaviourSelectVotesByPostID
		inUser                              *models.User
		inPostID                            uint64
		inCommentID                         uint64
		outUser                             *models.User
		outPost                             *models.Post
		outComments                         []*models.Comment
		outVotes                            []*models.Vote
		expPost                             *models.Post
		expError                            *customErrors.Error
	}{
		{
			name: "OK",
			mockBehaviourSelectAuthorComment: func(postRep *mock_post.MockPostRepository, user *models.User) {
				postRep.
					EXPECT().
					SelectAuthorComment(gomock.Any()).
					Return(user, nil)
			},
			mockBehaviourDeleteCommentByID: func(postRep *mock_post.MockPostRepository) {
				postRep.
					EXPECT().
					DeleteCommentByID(gomock.Any()).
					Return(nil)
			},
			mockBehaviourSelectPostByID: func(postRep *mock_post.MockPostRepository, post *models.Post) {
				postRep.
					EXPECT().
					SelectPostByID(gomock.Any()).
					Return(post, nil)
			},
			mockBehaviourSelectCommentsByPostID: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {
				postRep.
					EXPECT().
					SelectCommentsByPostID(gomock.Any()).
					Return(comments, nil)
			},
			mockBehaviourSelectVotesByPostID: func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {
				postRep.
					EXPECT().
					SelectVotesByPostID(gomock.Any()).
					Return(votes, nil)
			},
			inUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			inPostID:    1,
			inCommentID: 2,
			outUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			outPost: &models.Post{
				ID:    1,
				Title: "title",
				Text:  "text",
			},
			outComments: []*models.Comment{
				{
					ID: 3,
				},
			},
			outVotes: []*models.Vote{
				{
					ID: 5,
				},
			},
			expPost: &models.Post{
				ID:    1,
				Title: "title",
				Text:  "text",
				Comments: []*models.Comment{
					{
						ID: 3,
					},
				},
				Votes: []*models.Vote{
					{
						ID: 5,
					},
				},
			},
			expError: nil,
		},
		{
			name: "Error: DeleteCommentByID",
			mockBehaviourSelectAuthorComment: func(postRep *mock_post.MockPostRepository, user *models.User) {
				postRep.
					EXPECT().
					SelectAuthorComment(gomock.Any()).
					Return(user, nil)
			},
			mockBehaviourDeleteCommentByID: func(postRep *mock_post.MockPostRepository) {
				postRep.
					EXPECT().
					DeleteCommentByID(gomock.Any()).
					Return(fmt.Errorf("sql error"))
			},
			mockBehaviourSelectPostByID:         func(postRep *mock_post.MockPostRepository, post *models.Post) {},
			mockBehaviourSelectCommentsByPostID: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {},
			mockBehaviourSelectVotesByPostID:    func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {},
			inUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			inPostID:    1,
			inCommentID: 2,
			outUser: &models.User{
				ID:       1,
				Username: "testuser",
			},

			expError: customErrors.Get(consts.CodeInternalError),
		},
		{
			name: "Error: SelectAuthorComment",
			mockBehaviourSelectAuthorComment: func(postRep *mock_post.MockPostRepository, user *models.User) {
				postRep.
					EXPECT().
					SelectAuthorComment(gomock.Any()).
					Return(nil, fmt.Errorf("sql error"))
			},
			mockBehaviourDeleteCommentByID:      func(postRep *mock_post.MockPostRepository) {},
			mockBehaviourSelectPostByID:         func(postRep *mock_post.MockPostRepository, post *models.Post) {},
			mockBehaviourSelectCommentsByPostID: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {},
			mockBehaviourSelectVotesByPostID:    func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {},
			inUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			inPostID:    1,
			inCommentID: 2,
			expError:    customErrors.Get(consts.CodeInternalError),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			postRep := mock_post.NewMockPostRepository(ctrl)
			testCase.mockBehaviourSelectAuthorComment(postRep, testCase.outUser)
			testCase.mockBehaviourDeleteCommentByID(postRep)
			testCase.mockBehaviourSelectPostByID(postRep, testCase.outPost)
			testCase.mockBehaviourSelectCommentsByPostID(postRep, testCase.outComments)
			testCase.mockBehaviourSelectVotesByPostID(postRep, testCase.outVotes)
			postUse := NewPostUsecase(postRep)

			post, err := postUse.DeleteComment(testCase.inUser, testCase.inPostID, testCase.inCommentID)
			assert.Equal(t, post, testCase.expPost)
			assert.Equal(t, err, testCase.expError)
		})
	}
}

func Test_UpvotePost(t *testing.T) {
	type mockBehaviourDeleteVoteFromPostByUserID func(postRep *mock_post.MockPostRepository)
	type mockBehaviourInsertVote func(postRep *mock_post.MockPostRepository)
	type mockBehaviourSelectPostByID func(postRep *mock_post.MockPostRepository, post *models.Post)
	type mockBehaviourSelectCommentsByPostID func(postRep *mock_post.MockPostRepository, comments []*models.Comment)
	type mockBehaviourSelectVotesByPostID func(postRep *mock_post.MockPostRepository, votes []*models.Vote)
	t.Parallel()

	testTable := []struct {
		name                                    string
		mockBehaviourDeleteVoteFromPostByUserID mockBehaviourDeleteVoteFromPostByUserID
		mockBehaviourInsertVote                 mockBehaviourInsertVote
		mockBehaviourSelectPostByID             mockBehaviourSelectPostByID
		mockBehaviourSelectCommentsByPostID     mockBehaviourSelectCommentsByPostID
		mockBehaviourSelectVotesByPostID        mockBehaviourSelectVotesByPostID
		inUser                                  *models.User
		inPostID                                uint64
		outPost                                 *models.Post
		outComments                             []*models.Comment
		outVotes                                []*models.Vote
		expPost                                 *models.Post
		expError                                *customErrors.Error
	}{
		{
			name: "OK",
			mockBehaviourDeleteVoteFromPostByUserID: func(postRep *mock_post.MockPostRepository) {
				postRep.
					EXPECT().
					DeleteVoteFromPostByUserID(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			mockBehaviourInsertVote: func(postRep *mock_post.MockPostRepository) {
				postRep.
					EXPECT().
					InsertVote(gomock.Any()).
					Return(uint64(0), nil)
			},
			mockBehaviourSelectPostByID: func(postRep *mock_post.MockPostRepository, post *models.Post) {
				postRep.
					EXPECT().
					SelectPostByID(gomock.Any()).
					Return(post, nil)
			},
			mockBehaviourSelectCommentsByPostID: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {
				postRep.
					EXPECT().
					SelectCommentsByPostID(gomock.Any()).
					Return(comments, nil)
			},
			mockBehaviourSelectVotesByPostID: func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {
				postRep.
					EXPECT().
					SelectVotesByPostID(gomock.Any()).
					Return(votes, nil)
			},
			inUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			inPostID: 1,
			outPost: &models.Post{
				ID:    1,
				Title: "title",
				Text:  "text",
			},
			outComments: []*models.Comment{
				{
					ID: 3,
				},
			},
			outVotes: []*models.Vote{
				{
					ID: 5,
				},
			},
			expPost: &models.Post{
				ID:    1,
				Title: "title",
				Text:  "text",
				Comments: []*models.Comment{
					{
						ID: 3,
					},
				},
				Votes: []*models.Vote{
					{
						ID: 5,
					},
				},
			},
			expError: nil,
		},
		{
			name: "Error: InsertVote",
			mockBehaviourDeleteVoteFromPostByUserID: func(postRep *mock_post.MockPostRepository) {
				postRep.
					EXPECT().
					DeleteVoteFromPostByUserID(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			mockBehaviourInsertVote: func(postRep *mock_post.MockPostRepository) {
				postRep.
					EXPECT().
					InsertVote(gomock.Any()).
					Return(uint64(0), fmt.Errorf("sql error"))
			},
			mockBehaviourSelectPostByID:         func(postRep *mock_post.MockPostRepository, post *models.Post) {},
			mockBehaviourSelectCommentsByPostID: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {},
			mockBehaviourSelectVotesByPostID:    func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {},
			inUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			inPostID: 1,
			expError: customErrors.Get(consts.CodeInternalError),
		},
		{
			name: "Error: DeleteVoteFromPostByUserID",
			mockBehaviourDeleteVoteFromPostByUserID: func(postRep *mock_post.MockPostRepository) {
				postRep.
					EXPECT().
					DeleteVoteFromPostByUserID(gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("sql error"))
			},
			mockBehaviourInsertVote:             func(postRep *mock_post.MockPostRepository) {},
			mockBehaviourSelectPostByID:         func(postRep *mock_post.MockPostRepository, post *models.Post) {},
			mockBehaviourSelectCommentsByPostID: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {},
			mockBehaviourSelectVotesByPostID:    func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {},
			inUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			inPostID: 1,
			expError: customErrors.Get(consts.CodeInternalError),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			postRep := mock_post.NewMockPostRepository(ctrl)
			testCase.mockBehaviourDeleteVoteFromPostByUserID(postRep)
			testCase.mockBehaviourInsertVote(postRep)
			testCase.mockBehaviourSelectPostByID(postRep, testCase.outPost)
			testCase.mockBehaviourSelectCommentsByPostID(postRep, testCase.outComments)
			testCase.mockBehaviourSelectVotesByPostID(postRep, testCase.outVotes)
			postUse := NewPostUsecase(postRep)

			post, err := postUse.UpvotePost(testCase.inUser, testCase.inPostID)
			assert.Equal(t, post, testCase.expPost)
			assert.Equal(t, err, testCase.expError)
		})
	}
}

func Test_UnvotePost(t *testing.T) {
	type mockBehaviourDeleteVoteFromPostByUserID func(postRep *mock_post.MockPostRepository)
	type mockBehaviourSelectPostByID func(postRep *mock_post.MockPostRepository, post *models.Post)
	type mockBehaviourSelectCommentsByPostID func(postRep *mock_post.MockPostRepository, comments []*models.Comment)
	type mockBehaviourSelectVotesByPostID func(postRep *mock_post.MockPostRepository, votes []*models.Vote)
	t.Parallel()

	testTable := []struct {
		name                                    string
		mockBehaviourDeleteVoteFromPostByUserID mockBehaviourDeleteVoteFromPostByUserID
		mockBehaviourSelectPostByID             mockBehaviourSelectPostByID
		mockBehaviourSelectCommentsByPostID     mockBehaviourSelectCommentsByPostID
		mockBehaviourSelectVotesByPostID        mockBehaviourSelectVotesByPostID
		inUser                                  *models.User
		inPostID                                uint64
		outPost                                 *models.Post
		outComments                             []*models.Comment
		outVotes                                []*models.Vote
		expPost                                 *models.Post
		expError                                *customErrors.Error
	}{
		{
			name: "OK",
			mockBehaviourDeleteVoteFromPostByUserID: func(postRep *mock_post.MockPostRepository) {
				postRep.
					EXPECT().
					DeleteVoteFromPostByUserID(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			mockBehaviourSelectPostByID: func(postRep *mock_post.MockPostRepository, post *models.Post) {
				postRep.
					EXPECT().
					SelectPostByID(gomock.Any()).
					Return(post, nil)
			},
			mockBehaviourSelectCommentsByPostID: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {
				postRep.
					EXPECT().
					SelectCommentsByPostID(gomock.Any()).
					Return(comments, nil)
			},
			mockBehaviourSelectVotesByPostID: func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {
				postRep.
					EXPECT().
					SelectVotesByPostID(gomock.Any()).
					Return(votes, nil)
			},
			inUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			inPostID: 1,
			outPost: &models.Post{
				ID:    1,
				Title: "title",
				Text:  "text",
			},
			outComments: []*models.Comment{
				{
					ID: 3,
				},
			},
			outVotes: []*models.Vote{
				{
					ID: 5,
				},
			},
			expPost: &models.Post{
				ID:    1,
				Title: "title",
				Text:  "text",
				Comments: []*models.Comment{
					{
						ID: 3,
					},
				},
				Votes: []*models.Vote{
					{
						ID: 5,
					},
				},
			},
			expError: nil,
		},
		{
			name: "Error: DeleteVoteFromPostByUserID",
			mockBehaviourDeleteVoteFromPostByUserID: func(postRep *mock_post.MockPostRepository) {
				postRep.
					EXPECT().
					DeleteVoteFromPostByUserID(gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("sql error"))
			},
			mockBehaviourSelectPostByID:         func(postRep *mock_post.MockPostRepository, post *models.Post) {},
			mockBehaviourSelectCommentsByPostID: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {},
			mockBehaviourSelectVotesByPostID:    func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {},
			inUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			inPostID: 1,
			expError: customErrors.Get(consts.CodeInternalError),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			postRep := mock_post.NewMockPostRepository(ctrl)
			testCase.mockBehaviourDeleteVoteFromPostByUserID(postRep)
			testCase.mockBehaviourSelectPostByID(postRep, testCase.outPost)
			testCase.mockBehaviourSelectCommentsByPostID(postRep, testCase.outComments)
			testCase.mockBehaviourSelectVotesByPostID(postRep, testCase.outVotes)
			postUse := NewPostUsecase(postRep)

			post, err := postUse.UnvotePost(testCase.inUser, testCase.inPostID)
			assert.Equal(t, post, testCase.expPost)
			assert.Equal(t, err, testCase.expError)
		})
	}
}

func Test_DownvotePost(t *testing.T) {
	type mockBehaviourDeleteVoteFromPostByUserID func(postRep *mock_post.MockPostRepository)
	type mockBehaviourInsertVote func(postRep *mock_post.MockPostRepository)
	type mockBehaviourSelectPostByID func(postRep *mock_post.MockPostRepository, post *models.Post)
	type mockBehaviourSelectCommentsByPostID func(postRep *mock_post.MockPostRepository, comments []*models.Comment)
	type mockBehaviourSelectVotesByPostID func(postRep *mock_post.MockPostRepository, votes []*models.Vote)
	t.Parallel()

	testTable := []struct {
		name                                    string
		mockBehaviourDeleteVoteFromPostByUserID mockBehaviourDeleteVoteFromPostByUserID
		mockBehaviourInsertVote                 mockBehaviourInsertVote
		mockBehaviourSelectPostByID             mockBehaviourSelectPostByID
		mockBehaviourSelectCommentsByPostID     mockBehaviourSelectCommentsByPostID
		mockBehaviourSelectVotesByPostID        mockBehaviourSelectVotesByPostID
		inUser                                  *models.User
		inPostID                                uint64
		outPost                                 *models.Post
		outComments                             []*models.Comment
		outVotes                                []*models.Vote
		expPost                                 *models.Post
		expError                                *customErrors.Error
	}{
		{
			name: "OK",
			mockBehaviourDeleteVoteFromPostByUserID: func(postRep *mock_post.MockPostRepository) {
				postRep.
					EXPECT().
					DeleteVoteFromPostByUserID(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			mockBehaviourInsertVote: func(postRep *mock_post.MockPostRepository) {
				postRep.
					EXPECT().
					InsertVote(gomock.Any()).
					Return(uint64(0), nil)
			},
			mockBehaviourSelectPostByID: func(postRep *mock_post.MockPostRepository, post *models.Post) {
				postRep.
					EXPECT().
					SelectPostByID(gomock.Any()).
					Return(post, nil)
			},
			mockBehaviourSelectCommentsByPostID: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {
				postRep.
					EXPECT().
					SelectCommentsByPostID(gomock.Any()).
					Return(comments, nil)
			},
			mockBehaviourSelectVotesByPostID: func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {
				postRep.
					EXPECT().
					SelectVotesByPostID(gomock.Any()).
					Return(votes, nil)
			},
			inUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			inPostID: 1,
			outPost: &models.Post{
				ID:    1,
				Title: "title",
				Text:  "text",
			},
			outComments: []*models.Comment{
				{
					ID: 3,
				},
			},
			outVotes: []*models.Vote{
				{
					ID: 5,
				},
			},
			expPost: &models.Post{
				ID:    1,
				Title: "title",
				Text:  "text",
				Comments: []*models.Comment{
					{
						ID: 3,
					},
				},
				Votes: []*models.Vote{
					{
						ID: 5,
					},
				},
			},
			expError: nil,
		},
		{
			name: "Error: InsertVote",
			mockBehaviourDeleteVoteFromPostByUserID: func(postRep *mock_post.MockPostRepository) {
				postRep.
					EXPECT().
					DeleteVoteFromPostByUserID(gomock.Any(), gomock.Any()).
					Return(nil)
			},
			mockBehaviourInsertVote: func(postRep *mock_post.MockPostRepository) {
				postRep.
					EXPECT().
					InsertVote(gomock.Any()).
					Return(uint64(0), fmt.Errorf("sql error"))
			},
			mockBehaviourSelectPostByID:         func(postRep *mock_post.MockPostRepository, post *models.Post) {},
			mockBehaviourSelectCommentsByPostID: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {},
			mockBehaviourSelectVotesByPostID:    func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {},
			inUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			inPostID: 1,
			expError: customErrors.Get(consts.CodeInternalError),
		},
		{
			name: "Error: DeleteVoteFromPostByUserID",
			mockBehaviourDeleteVoteFromPostByUserID: func(postRep *mock_post.MockPostRepository) {
				postRep.
					EXPECT().
					DeleteVoteFromPostByUserID(gomock.Any(), gomock.Any()).
					Return(fmt.Errorf("sql error"))
			},
			mockBehaviourInsertVote:             func(postRep *mock_post.MockPostRepository) {},
			mockBehaviourSelectPostByID:         func(postRep *mock_post.MockPostRepository, post *models.Post) {},
			mockBehaviourSelectCommentsByPostID: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {},
			mockBehaviourSelectVotesByPostID:    func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {},
			inUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			inPostID: 1,
			expError: customErrors.Get(consts.CodeInternalError),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			postRep := mock_post.NewMockPostRepository(ctrl)
			testCase.mockBehaviourDeleteVoteFromPostByUserID(postRep)
			testCase.mockBehaviourInsertVote(postRep)
			testCase.mockBehaviourSelectPostByID(postRep, testCase.outPost)
			testCase.mockBehaviourSelectCommentsByPostID(postRep, testCase.outComments)
			testCase.mockBehaviourSelectVotesByPostID(postRep, testCase.outVotes)
			postUse := NewPostUsecase(postRep)

			post, err := postUse.DownvotePost(testCase.inUser, testCase.inPostID)
			assert.Equal(t, post, testCase.expPost)
			assert.Equal(t, err, testCase.expError)
		})
	}
}

func Test_GetPostsByUser(t *testing.T) {
	type mockBehaviourSelectPostsByUsername func(postRep *mock_post.MockPostRepository, posts []*models.Post)
	type mockBehaviourSelectCommentsByUsername func(postRep *mock_post.MockPostRepository, comments []*models.Comment)
	type mockBehaviourSelectVotesByUsername func(postRep *mock_post.MockPostRepository, votes []*models.Vote)
	t.Parallel()

	testTable := []struct {
		name                                  string
		mockBehaviourSelectPostsByUsername    mockBehaviourSelectPostsByUsername
		mockBehaviourSelectCommentsByUsername mockBehaviourSelectCommentsByUsername
		mockBehaviourSelectVotesByUsername    mockBehaviourSelectVotesByUsername
		inPostCategory                        string
		outPosts                              []*models.Post
		outComments                           []*models.Comment
		outVotes                              []*models.Vote
		expPosts                              []*models.Post
		expError                              *customErrors.Error
	}{
		{
			name: "OK",
			mockBehaviourSelectPostsByUsername: func(postRep *mock_post.MockPostRepository, posts []*models.Post) {
				postRep.
					EXPECT().
					SelectPostsByUsername(gomock.Any()).
					Return(posts, nil)
			},
			mockBehaviourSelectCommentsByUsername: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {
				postRep.
					EXPECT().
					SelectCommentsByUsername(gomock.Any()).
					Return(comments, nil)
			},
			mockBehaviourSelectVotesByUsername: func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {
				postRep.
					EXPECT().
					SelectVotesByUsername(gomock.Any()).
					Return(votes, nil)
			},
			inPostCategory: "music",
			outPosts: []*models.Post{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			},
			outComments: []*models.Comment{
				{
					ID:     7,
					PostID: 2,
				},
				{
					ID:     8,
					PostID: 1,
				},
			},
			outVotes: []*models.Vote{
				{
					ID:     5,
					PostID: 1,
				},
				{
					ID:     6,
					PostID: 2,
				},
			},
			expPosts: []*models.Post{
				{
					ID: 1,
					Comments: []*models.Comment{
						{
							ID:     8,
							PostID: 1,
						},
					},
					Votes: []*models.Vote{
						{
							ID:     5,
							PostID: 1,
						},
					},
				},
				{
					ID: 2,
					Comments: []*models.Comment{
						{
							ID:     7,
							PostID: 2,
						},
					},
					Votes: []*models.Vote{
						{
							ID:     6,
							PostID: 2,
						},
					},
				},
			},
			expError: nil,
		},
		{
			name: "Error: SelectAllVotes",
			mockBehaviourSelectPostsByUsername: func(postRep *mock_post.MockPostRepository, posts []*models.Post) {
				postRep.
					EXPECT().
					SelectPostsByUsername(gomock.Any()).
					Return(posts, nil)
			},
			mockBehaviourSelectCommentsByUsername: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {
				postRep.
					EXPECT().
					SelectCommentsByUsername(gomock.Any()).
					Return(comments, nil)
			},
			mockBehaviourSelectVotesByUsername: func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {
				postRep.
					EXPECT().
					SelectVotesByUsername(gomock.Any()).
					Return(nil, fmt.Errorf("sql error"))
			},
			inPostCategory: "music",
			outPosts: []*models.Post{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			},
			outComments: []*models.Comment{
				{
					ID:     7,
					PostID: 2,
				},
				{
					ID:     8,
					PostID: 1,
				},
			},
			expError: customErrors.Get(consts.CodeInternalError),
		},
		{
			name: "Error: SelectCommentsByCategory:",
			mockBehaviourSelectPostsByUsername: func(postRep *mock_post.MockPostRepository, posts []*models.Post) {
				postRep.
					EXPECT().
					SelectPostsByUsername(gomock.Any()).
					Return(posts, nil)
			},
			mockBehaviourSelectCommentsByUsername: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {
				postRep.
					EXPECT().
					SelectCommentsByUsername(gomock.Any()).
					Return(nil, fmt.Errorf("sql error"))
			},
			mockBehaviourSelectVotesByUsername: func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {},
			inPostCategory:                     "music",
			outPosts: []*models.Post{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			},
			expError: customErrors.Get(consts.CodeInternalError),
		},
		{
			name: "Error: SelectAllPosts:",
			mockBehaviourSelectPostsByUsername: func(postRep *mock_post.MockPostRepository, posts []*models.Post) {
				postRep.
					EXPECT().
					SelectPostsByUsername(gomock.Any()).
					Return(nil, fmt.Errorf("sql error"))
			},
			mockBehaviourSelectCommentsByUsername: func(postRep *mock_post.MockPostRepository, comments []*models.Comment) {},
			mockBehaviourSelectVotesByUsername:    func(postRep *mock_post.MockPostRepository, votes []*models.Vote) {},
			inPostCategory:                        "music",
			expError:                              customErrors.Get(consts.CodeInternalError),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			postRep := mock_post.NewMockPostRepository(ctrl)
			testCase.mockBehaviourSelectPostsByUsername(postRep, testCase.outPosts)
			testCase.mockBehaviourSelectCommentsByUsername(postRep, testCase.outComments)
			testCase.mockBehaviourSelectVotesByUsername(postRep, testCase.outVotes)
			postUse := NewPostUsecase(postRep)

			posts, err := postUse.GetPostsByUser(testCase.inPostCategory)
			assert.Equal(t, posts, testCase.expPosts)
			assert.Equal(t, err, testCase.expError)
		})
	}
}

func Test_DeletePost(t *testing.T) {
	type mockBehaviourSelectAuthorPost func(postRep *mock_post.MockPostRepository, user *models.User)
	type mockBehaviourDeletePostByID func(postRep *mock_post.MockPostRepository)

	t.Parallel()

	testTable := []struct {
		name                          string
		mockBehaviourSelectAuthorPost mockBehaviourSelectAuthorPost
		mockBehaviourDeletePostByID   mockBehaviourDeletePostByID
		inUser                        *models.User
		inPostID                      uint64
		outUser                       *models.User
		expError                      *customErrors.Error
	}{
		{
			name: "OK",
			mockBehaviourSelectAuthorPost: func(postRep *mock_post.MockPostRepository, user *models.User) {
				postRep.
					EXPECT().
					SelectAuthorPost(gomock.Any()).
					Return(user, nil)
			},
			mockBehaviourDeletePostByID: func(postRep *mock_post.MockPostRepository) {
				postRep.
					EXPECT().
					DeletePostByID(gomock.Any()).
					Return(nil)
			},
			inUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			inPostID: 1,
			outUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			expError: nil,
		},
		{
			name: "Error: DeletePostByID",
			mockBehaviourSelectAuthorPost: func(postRep *mock_post.MockPostRepository, user *models.User) {
				postRep.
					EXPECT().
					SelectAuthorPost(gomock.Any()).
					Return(user, nil)
			},
			mockBehaviourDeletePostByID: func(postRep *mock_post.MockPostRepository) {
				postRep.
					EXPECT().
					DeletePostByID(gomock.Any()).
					Return(fmt.Errorf("sql error"))
			},
			inUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			inPostID: 1,
			outUser: &models.User{
				ID:       1,
				Username: "testuser",
			},

			expError: customErrors.Get(consts.CodeInternalError),
		},
		{
			name: "Error: SelectAuthorComment",
			mockBehaviourSelectAuthorPost: func(postRep *mock_post.MockPostRepository, user *models.User) {
				postRep.
					EXPECT().
					SelectAuthorPost(gomock.Any()).
					Return(nil, fmt.Errorf("sql error"))
			},
			mockBehaviourDeletePostByID: func(postRep *mock_post.MockPostRepository) {},
			inUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			inPostID: 1,
			expError: customErrors.Get(consts.CodeInternalError),
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			postRep := mock_post.NewMockPostRepository(ctrl)
			testCase.mockBehaviourSelectAuthorPost(postRep, testCase.outUser)
			testCase.mockBehaviourDeletePostByID(postRep)
			postUse := NewPostUsecase(postRep)

			err := postUse.DeletePost(testCase.inUser, testCase.inPostID)
			assert.Equal(t, err, testCase.expError)
		})
	}
}
