package repository

import (
	"database/sql/driver"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/models"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func Test_InsertPost(t *testing.T) {
	type mockBehaviour func(mock sqlmock.Sqlmock, post *models.Post, postID uint64)

	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inPost        *models.Post
		outPostID     uint64
		expError      error
	}{
		{
			name: "OK",
			mockBehaviour: func(mock sqlmock.Sqlmock, post *models.Post, postID uint64) {
				mock.ExpectBegin()
				lastID := sqlmock.NewRows([]string{"id"}).AddRow(postID)
				mock.ExpectQuery(`INSERT INTO posts`).
					WithArgs(post.UserID, post.Category, post.CreatedAt, post.Score, post.Text, post.Title, post.Type, post.UpvotePercentage, post.URL, post.Views).
					WillReturnRows(lastID)
				mock.ExpectCommit()
			},
			inPost: &models.Post{
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
			outPostID: 1,
			expError:  nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			postRep := NewPostRepository(db)
			testCase.mockBehaviour(mock, testCase.inPost, testCase.outPostID)
			postID, err := postRep.InsertPost(testCase.inPost)

			assert.Equal(t, postID, testCase.outPostID)
			assert.Equal(t, err, testCase.expError)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("were met expectation: %s", err)
			}
		})
	}
}

func Test_DeletePostByID(t *testing.T) {
	type mockBehaviour func(mock sqlmock.Sqlmock, inPostID uint64)

	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inPostID      uint64
		expError      error
	}{
		{
			name: "OK",
			mockBehaviour: func(mock sqlmock.Sqlmock, inPostID uint64) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM posts`).
					WithArgs(inPostID).
					WillReturnResult(driver.ResultNoRows)

				mock.ExpectExec(`DELETE FROM comments`).
					WithArgs(inPostID).
					WillReturnResult(driver.ResultNoRows)

				mock.ExpectExec(`DELETE FROM votes`).
					WithArgs(inPostID).
					WillReturnResult(driver.ResultNoRows)
				mock.ExpectCommit()
			},
			inPostID: 1,
			expError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			postRep := NewPostRepository(db)
			testCase.mockBehaviour(mock, testCase.inPostID)
			err := postRep.DeletePostByID(testCase.inPostID)

			assert.Equal(t, err, testCase.expError)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("were met expectation: %s", err)
			}
		})
	}
}

func Test_SelectAuthorPost(t *testing.T) {
	type mockBehaviour func(mock sqlmock.Sqlmock, user *models.User, postID uint64)

	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inPostID      uint64
		outUser       *models.User
		expError      error
	}{
		{
			name: "OK",
			mockBehaviour: func(mock sqlmock.Sqlmock, user *models.User, postID uint64) {
				rows := sqlmock.NewRows([]string{"id", "username"})
				rows.AddRow(user.ID, user.Username)
				mock.ExpectQuery(`SELECT`).WithArgs(postID).WillReturnRows(rows)
			},
			inPostID: 1,
			outUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			expError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			postRep := NewPostRepository(db)
			testCase.mockBehaviour(mock, testCase.outUser, testCase.inPostID)
			user, err := postRep.SelectAuthorPost(testCase.inPostID)

			assert.Equal(t, user, testCase.outUser)
			assert.Equal(t, err, testCase.expError)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("were met expectation: %s", err)
			}
		})
	}
}

func Test_SelectAuthorComment(t *testing.T) {
	type mockBehaviour func(mock sqlmock.Sqlmock, user *models.User, postID uint64)

	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inCommentID   uint64
		outUser       *models.User
		expError      error
	}{
		{
			name: "OK",
			mockBehaviour: func(mock sqlmock.Sqlmock, user *models.User, commentID uint64) {
				rows := sqlmock.NewRows([]string{"id", "username"})
				rows.AddRow(user.ID, user.Username)
				mock.ExpectQuery(`SELECT`).WithArgs(commentID).WillReturnRows(rows)
			},
			inCommentID: 1,
			outUser: &models.User{
				ID:       1,
				Username: "testuser",
			},
			expError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			postRep := NewPostRepository(db)
			testCase.mockBehaviour(mock, testCase.outUser, testCase.inCommentID)
			user, err := postRep.SelectAuthorComment(testCase.inCommentID)

			assert.Equal(t, user, testCase.outUser)
			assert.Equal(t, err, testCase.expError)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("were met expectation: %s", err)
			}
		})
	}
}

func Test_DeleteVoteFromPostByUserID(t *testing.T) {
	type mockBehaviour func(mock sqlmock.Sqlmock, inUserID, inPostID uint64)

	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inUserID      uint64
		inPostID      uint64
		expError      error
	}{
		{
			name: "OK",
			mockBehaviour: func(mock sqlmock.Sqlmock, inUserID, inPostID uint64) {
				mock.ExpectBegin()
				mock.ExpectExec(`DELETE FROM votes`).
					WithArgs(inUserID, inPostID).
					WillReturnResult(driver.ResultNoRows)
				mock.ExpectCommit()
			},
			inUserID: 1,
			inPostID: 1,
			expError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			postRep := NewPostRepository(db)
			testCase.mockBehaviour(mock, testCase.inUserID, testCase.inPostID)
			err := postRep.DeleteVoteFromPostByUserID(testCase.inUserID, testCase.inPostID)

			assert.Equal(t, err, testCase.expError)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("were met expectation: %s", err)
			}
		})
	}
}

func Test_InsertVote(t *testing.T) {
	type mockBehaviour func(mock sqlmock.Sqlmock, vote *models.Vote, voteID uint64)

	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inVote        *models.Vote
		outVoteID     uint64
		expError      error
	}{
		{
			name: "OK",
			mockBehaviour: func(mock sqlmock.Sqlmock, vote *models.Vote, voteID uint64) {
				mock.ExpectBegin()
				lastID := sqlmock.NewRows([]string{"id"}).AddRow(voteID)
				mock.ExpectQuery(`INSERT INTO votes`).
					WithArgs(vote.UserID, vote.PostID, vote.Vote).
					WillReturnRows(lastID)
				mock.ExpectCommit()
			},
			inVote: &models.Vote{
				UserID: 1,
				PostID: 2,
				Vote:   -1,
			},
			outVoteID: 1,
			expError:  nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			postRep := NewPostRepository(db)
			testCase.mockBehaviour(mock, testCase.inVote, testCase.outVoteID)
			postID, err := postRep.InsertVote(testCase.inVote)

			assert.Equal(t, postID, testCase.outVoteID)
			assert.Equal(t, err, testCase.expError)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("were met expectation: %s", err)
			}
		})
	}
}

func Test_SelectAllPosts(t *testing.T) {
	type mockBehaviour func(mock sqlmock.Sqlmock, users []*models.User, posts []*models.Post)
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inPosts       []*models.Post
		inUsers       []*models.User
		expPosts      []*models.Post
		expError      error
	}{
		{
			name: "OK",
			mockBehaviour: func(mock sqlmock.Sqlmock, users []*models.User, posts []*models.Post) {
				rows := sqlmock.NewRows([]string{"usr.id", "usr.username", "pst.id", "pst.user_id", "pst.category", "pst.created", "pst.score", "pst.text", "pst.title", "pst.type", "pst.upvote_percentage", "pst.url", "pst.views"})
				for i := range posts {
					rows.AddRow(users[i].ID, users[i].Username, posts[i].ID, posts[i].UserID, posts[i].Category, posts[i].CreatedAt, posts[i].Score, posts[i].Text, posts[i].Title, posts[i].Type, posts[i].UpvotePercentage, posts[i].URL, posts[i].Views)
				}
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows)
			},
			inPosts: []*models.Post{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			},
			inUsers: []*models.User{
				{
					ID: 5,
				},
				{
					ID: 6,
				},
			},
			expPosts: []*models.Post{
				{
					ID: 1,
					Author: &models.User{
						ID: 5,
					},
				},
				{
					ID: 2,
					Author: &models.User{
						ID: 6,
					},
				},
			},
			expError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			postRep := NewPostRepository(db)
			testCase.mockBehaviour(mock, testCase.inUsers, testCase.inPosts)
			posts, err := postRep.SelectAllPosts()

			assert.Equal(t, posts, testCase.expPosts)
			assert.Equal(t, err, testCase.expError)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("were met expectation: %s", err)
			}
		})
	}
}

func Test_SelectAllComments(t *testing.T) {
	type mockBehaviour func(mock sqlmock.Sqlmock, users []*models.User, comments []*models.Comment)
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inComments    []*models.Comment
		inUsers       []*models.User
		expComments   []*models.Comment
		expError      error
	}{
		{
			name: "OK",
			mockBehaviour: func(mock sqlmock.Sqlmock, users []*models.User, comments []*models.Comment) {
				rows := sqlmock.NewRows([]string{"usr.id", "usr.username", "cmt.post_id", "cmt.id", "cmt.body", "cmt.created"})
				for i := range comments {
					rows.AddRow(users[i].ID, users[i].Username, comments[i].PostID, comments[i].ID, comments[i].Body, comments[i].CreatedAt)
				}
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows)
			},
			inComments: []*models.Comment{
				{
					ID: 1,
				},
			},
			inUsers: []*models.User{
				{
					ID: 5,
				},
			},
			expComments: []*models.Comment{
				{
					ID: 1,
					Author: &models.User{
						ID: 5,
					},
				},
			},
			expError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			postRep := NewPostRepository(db)
			testCase.mockBehaviour(mock, testCase.inUsers, testCase.inComments)
			posts, err := postRep.SelectAllComments()

			assert.Equal(t, posts, testCase.expComments)
			assert.Equal(t, err, testCase.expError)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("were met expectation: %s", err)
			}
		})
	}
}

func Test_SelectAllVotes(t *testing.T) {
	type mockBehaviour func(mock sqlmock.Sqlmock, votes []*models.Vote)
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inVotes       []*models.Vote
		expVotes      []*models.Vote
		expError      error
	}{
		{
			name: "OK",
			mockBehaviour: func(mock sqlmock.Sqlmock, votes []*models.Vote) {
				rows := sqlmock.NewRows([]string{"vts.post_id", "vts.user_id", "vts.vote"})
				for _, vote := range votes {
					rows.AddRow(vote.PostID, vote.User, vote.Vote)
				}
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows)
			},
			inVotes: []*models.Vote{
				{
					Vote: 1,
				},
			},
			expVotes: []*models.Vote{
				{
					Vote: 1,
				},
			},
			expError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			postRep := NewPostRepository(db)
			testCase.mockBehaviour(mock, testCase.inVotes)
			posts, err := postRep.SelectAllVotes()

			assert.Equal(t, posts, testCase.expVotes)
			assert.Equal(t, err, testCase.expError)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("were met expectation: %s", err)
			}
		})
	}
}

func Test_SelectPostsByCategory(t *testing.T) {
	type mockBehaviour func(mock sqlmock.Sqlmock, users []*models.User, posts []*models.Post, postCategory string)
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testTable := []struct {
		name           string
		mockBehaviour  mockBehaviour
		inPostCategory string
		inPosts        []*models.Post
		inUsers        []*models.User
		expPosts       []*models.Post
		expError       error
	}{
		{
			name: "OK",
			mockBehaviour: func(mock sqlmock.Sqlmock, users []*models.User, posts []*models.Post, postCategory string) {
				rows := sqlmock.NewRows([]string{"usr.id", "usr.username", "pst.id", "pst.user_id", "pst.category", "pst.created", "pst.score", "pst.text", "pst.title", "pst.type", "pst.upvote_percentage", "pst.url", "pst.views"})
				for i := range posts {
					rows.AddRow(users[i].ID, users[i].Username, posts[i].ID, posts[i].UserID, posts[i].Category, posts[i].CreatedAt, posts[i].Score, posts[i].Text, posts[i].Title, posts[i].Type, posts[i].UpvotePercentage, posts[i].URL, posts[i].Views)
				}
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows).WithArgs(postCategory)
			},
			inPostCategory: "music",
			inPosts: []*models.Post{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			},
			inUsers: []*models.User{
				{
					ID: 5,
				},
				{
					ID: 6,
				},
			},
			expPosts: []*models.Post{
				{
					ID: 1,
					Author: &models.User{
						ID: 5,
					},
				},
				{
					ID: 2,
					Author: &models.User{
						ID: 6,
					},
				},
			},
			expError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			postRep := NewPostRepository(db)
			testCase.mockBehaviour(mock, testCase.inUsers, testCase.inPosts, testCase.inPostCategory)
			posts, err := postRep.SelectPostsByCategory(testCase.inPostCategory)

			assert.Equal(t, posts, testCase.expPosts)
			assert.Equal(t, err, testCase.expError)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("were met expectation: %s", err)
			}
		})
	}
}

func Test_SelectCommentsByCategory(t *testing.T) {
	type mockBehaviour func(mock sqlmock.Sqlmock, users []*models.User, comments []*models.Comment, postCategory string)
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testTable := []struct {
		name           string
		mockBehaviour  mockBehaviour
		inPostCategory string
		inComments     []*models.Comment
		inUsers        []*models.User
		expComments    []*models.Comment
		expError       error
	}{
		{
			name: "OK",
			mockBehaviour: func(mock sqlmock.Sqlmock, users []*models.User, comments []*models.Comment, postCategory string) {
				rows := sqlmock.NewRows([]string{"usr.id", "usr.username", "cmt.post_id", "cmt.id", "cmt.body", "cmt.created"})
				for i := range comments {
					rows.AddRow(users[i].ID, users[i].Username, comments[i].PostID, comments[i].ID, comments[i].Body, comments[i].CreatedAt)
				}
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows).WithArgs(postCategory)
			},
			inPostCategory: "music",
			inComments: []*models.Comment{
				{
					ID: 1,
				},
			},
			inUsers: []*models.User{
				{
					ID: 5,
				},
			},
			expComments: []*models.Comment{
				{
					ID: 1,
					Author: &models.User{
						ID: 5,
					},
				},
			},
			expError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			postRep := NewPostRepository(db)
			testCase.mockBehaviour(mock, testCase.inUsers, testCase.inComments, testCase.inPostCategory)
			posts, err := postRep.SelectCommentsByCategory(testCase.inPostCategory)

			assert.Equal(t, posts, testCase.expComments)
			assert.Equal(t, err, testCase.expError)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("were met expectation: %s", err)
			}
		})
	}
}

func Test_SelectVotesByCategory(t *testing.T) {
	type mockBehaviour func(mock sqlmock.Sqlmock, votes []*models.Vote, postCategory string)
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testTable := []struct {
		name           string
		mockBehaviour  mockBehaviour
		inPostCategory string
		inVotes        []*models.Vote
		expVotes       []*models.Vote
		expError       error
	}{
		{
			name: "OK",
			mockBehaviour: func(mock sqlmock.Sqlmock, votes []*models.Vote, postCategory string) {
				rows := sqlmock.NewRows([]string{"vts.post_id", "vts.user_id", "vts.vote"})
				for _, vote := range votes {
					rows.AddRow(vote.PostID, vote.User, vote.Vote)
				}
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows).WithArgs(postCategory)
			},
			inPostCategory: "music",
			inVotes: []*models.Vote{
				{
					Vote: 1,
				},
			},
			expVotes: []*models.Vote{
				{
					Vote: 1,
				},
			},
			expError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			postRep := NewPostRepository(db)
			testCase.mockBehaviour(mock, testCase.inVotes, testCase.inPostCategory)
			posts, err := postRep.SelectVotesByCategory(testCase.inPostCategory)

			assert.Equal(t, posts, testCase.expVotes)
			assert.Equal(t, err, testCase.expError)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("were met expectation: %s", err)
			}
		})
	}
}

func Test_SelectPostsByUsername(t *testing.T) {
	type mockBehaviour func(mock sqlmock.Sqlmock, users []*models.User, posts []*models.Post, postCategory string)
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inUsername    string
		inPosts       []*models.Post
		inUsers       []*models.User
		expPosts      []*models.Post
		expError      error
	}{
		{
			name: "OK",
			mockBehaviour: func(mock sqlmock.Sqlmock, users []*models.User, posts []*models.Post, username string) {
				rows := sqlmock.NewRows([]string{"usr.id", "usr.username", "pst.id", "pst.user_id", "pst.category", "pst.created", "pst.score", "pst.text", "pst.title", "pst.type", "pst.upvote_percentage", "pst.url", "pst.views"})
				for i := range posts {
					rows.AddRow(users[i].ID, users[i].Username, posts[i].ID, posts[i].UserID, posts[i].Category, posts[i].CreatedAt, posts[i].Score, posts[i].Text, posts[i].Title, posts[i].Type, posts[i].UpvotePercentage, posts[i].URL, posts[i].Views)
				}
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows).WithArgs(username)
			},
			inUsername: "testuser",
			inPosts: []*models.Post{
				{
					ID: 1,
				},
				{
					ID: 2,
				},
			},
			inUsers: []*models.User{
				{
					ID: 5,
				},
				{
					ID: 6,
				},
			},
			expPosts: []*models.Post{
				{
					ID: 1,
					Author: &models.User{
						ID: 5,
					},
				},
				{
					ID: 2,
					Author: &models.User{
						ID: 6,
					},
				},
			},
			expError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			postRep := NewPostRepository(db)
			testCase.mockBehaviour(mock, testCase.inUsers, testCase.inPosts, testCase.inUsername)
			posts, err := postRep.SelectPostsByUsername(testCase.inUsername)

			assert.Equal(t, posts, testCase.expPosts)
			assert.Equal(t, err, testCase.expError)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("were met expectation: %s", err)
			}
		})
	}
}

func Test_SelectCommentsByUsername(t *testing.T) {
	type mockBehaviour func(mock sqlmock.Sqlmock, users []*models.User, comments []*models.Comment, postCategory string)
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inUsername    string
		inComments    []*models.Comment
		inUsers       []*models.User
		expComments   []*models.Comment
		expError      error
	}{
		{
			name: "OK",
			mockBehaviour: func(mock sqlmock.Sqlmock, users []*models.User, comments []*models.Comment, username string) {
				rows := sqlmock.NewRows([]string{"usr.id", "usr.username", "cmt.post_id", "cmt.id", "cmt.body", "cmt.created"})
				for i := range comments {
					rows.AddRow(users[i].ID, users[i].Username, comments[i].PostID, comments[i].ID, comments[i].Body, comments[i].CreatedAt)
				}
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows).WithArgs(username)
			},
			inUsername: "testuser",
			inComments: []*models.Comment{
				{
					ID: 1,
				},
			},
			inUsers: []*models.User{
				{
					ID: 5,
				},
			},
			expComments: []*models.Comment{
				{
					ID: 1,
					Author: &models.User{
						ID: 5,
					},
				},
			},
			expError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			postRep := NewPostRepository(db)
			testCase.mockBehaviour(mock, testCase.inUsers, testCase.inComments, testCase.inUsername)
			posts, err := postRep.SelectCommentsByUsername(testCase.inUsername)

			assert.Equal(t, posts, testCase.expComments)
			assert.Equal(t, err, testCase.expError)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("were met expectation: %s", err)
			}
		})
	}
}

func Test_SelectVotesByUsername(t *testing.T) {
	type mockBehaviour func(mock sqlmock.Sqlmock, votes []*models.Vote, postCategory string)
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inUsername    string
		inVotes       []*models.Vote
		expVotes      []*models.Vote
		expError      error
	}{
		{
			name: "OK",
			mockBehaviour: func(mock sqlmock.Sqlmock, votes []*models.Vote, username string) {
				rows := sqlmock.NewRows([]string{"vts.post_id", "vts.user_id", "vts.vote"})
				for _, vote := range votes {
					rows.AddRow(vote.PostID, vote.User, vote.Vote)
				}
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows).WithArgs(username)
			},
			inUsername: "testuser",
			inVotes: []*models.Vote{
				{
					Vote: 1,
				},
			},
			expVotes: []*models.Vote{
				{
					Vote: 1,
				},
			},
			expError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			postRep := NewPostRepository(db)
			testCase.mockBehaviour(mock, testCase.inVotes, testCase.inUsername)
			posts, err := postRep.SelectVotesByUsername(testCase.inUsername)

			assert.Equal(t, posts, testCase.expVotes)
			assert.Equal(t, err, testCase.expError)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("were met expectation: %s", err)
			}
		})
	}
}

func Test_SelectPostByID(t *testing.T) {
	type mockBehaviour func(mock sqlmock.Sqlmock, user *models.User, post *models.Post, postID uint64)

	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inPostID      uint64
		inPost        *models.Post
		inUser        *models.User
		outPost       *models.Post
		expError      error
	}{
		{
			name: "OK",
			mockBehaviour: func(mock sqlmock.Sqlmock, user *models.User, post *models.Post, postID uint64) {
				rows := sqlmock.NewRows([]string{"usr.id", "usr.username", "pst.id", "pst.user_id", "pst.category", "pst.created", "pst.score", "pst.text", "pst.title", "pst.type", "pst.upvote_percentage", "pst.url", "pst.views"})
				rows.AddRow(user.ID, user.Username, post.ID, post.UserID, post.Category, post.CreatedAt, post.Score, post.Text, post.Title, post.Type, post.UpvotePercentage, post.URL, post.Views)
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows).WithArgs(postID)
			},
			inPostID: 1,
			inPost: &models.Post{
				ID: 1,
			},
			inUser: &models.User{
				ID: 6,
			},
			outPost: &models.Post{
				ID: 1,
				Author: &models.User{
					ID: 6,
				},
			},
			expError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			postRep := NewPostRepository(db)
			testCase.mockBehaviour(mock, testCase.inUser, testCase.inPost, testCase.inPostID)
			post, err := postRep.SelectPostByID(testCase.inPostID)

			assert.Equal(t, post, testCase.outPost)
			assert.Equal(t, err, testCase.expError)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("were met expectation: %s", err)
			}
		})
	}
}

func Test_SelectCommentsByPostID(t *testing.T) {
	type mockBehaviour func(mock sqlmock.Sqlmock, users []*models.User, comments []*models.Comment, postID uint64)
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inPostID      uint64
		inComments    []*models.Comment
		inUsers       []*models.User
		expComments   []*models.Comment
		expError      error
	}{
		{
			name: "OK",
			mockBehaviour: func(mock sqlmock.Sqlmock, users []*models.User, comments []*models.Comment, postID uint64) {
				rows := sqlmock.NewRows([]string{"usr.id", "usr.username", "cmt.id", "cmt.body", "cmt.created"})
				for i := range comments {
					rows.AddRow(users[i].ID, users[i].Username, comments[i].ID, comments[i].Body, comments[i].CreatedAt)
				}
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows).WithArgs(postID)
			},
			inPostID: 1,
			inComments: []*models.Comment{
				{
					ID: 1,
				},
			},
			inUsers: []*models.User{
				{
					ID: 5,
				},
			},
			expComments: []*models.Comment{
				{
					ID: 1,
					Author: &models.User{
						ID: 5,
					},
				},
			},
			expError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			postRep := NewPostRepository(db)
			testCase.mockBehaviour(mock, testCase.inUsers, testCase.inComments, testCase.inPostID)
			posts, err := postRep.SelectCommentsByPostID(testCase.inPostID)

			assert.Equal(t, posts, testCase.expComments)
			assert.Equal(t, err, testCase.expError)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("were met expectation: %s", err)
			}
		})
	}
}

func Test_SelectVotesByPostID(t *testing.T) {
	type mockBehaviour func(mock sqlmock.Sqlmock, votes []*models.Vote, postID uint64)
	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inPostID      uint64
		inVotes       []*models.Vote
		expVotes      []*models.Vote
		expError      error
	}{
		{
			name: "OK",
			mockBehaviour: func(mock sqlmock.Sqlmock, votes []*models.Vote, postID uint64) {
				rows := sqlmock.NewRows([]string{"vts.user_id", "vts.vote"})
				for _, vote := range votes {
					rows.AddRow(vote.User, vote.Vote)
				}
				mock.ExpectQuery(`SELECT`).WillReturnRows(rows).WithArgs(postID)
			},
			inPostID: 1,
			inVotes: []*models.Vote{
				{
					Vote: 1,
				},
			},
			expVotes: []*models.Vote{
				{
					Vote: 1,
				},
			},
			expError: nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			postRep := NewPostRepository(db)
			testCase.mockBehaviour(mock, testCase.inVotes, testCase.inPostID)
			posts, err := postRep.SelectVotesByPostID(testCase.inPostID)

			assert.Equal(t, posts, testCase.expVotes)
			assert.Equal(t, err, testCase.expError)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("were met expectation: %s", err)
			}
		})
	}
}

func Test_InsertComment(t *testing.T) {
	type mockBehaviour func(mock sqlmock.Sqlmock, comment *models.Comment, commentID uint64)

	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inComment     *models.Comment
		outCommentID  uint64
		expError      error
	}{
		{
			name: "OK",
			mockBehaviour: func(mock sqlmock.Sqlmock, comment *models.Comment, commentID uint64) {
				mock.ExpectBegin()
				lastID := sqlmock.NewRows([]string{"id"}).AddRow(commentID)
				mock.ExpectQuery(`INSERT INTO comments`).
					WithArgs(comment.UserID, comment.PostID, comment.Body, comment.CreatedAt).
					WillReturnRows(lastID)
				mock.ExpectCommit()
			},
			inComment: &models.Comment{
				ID: 1,
			},
			outCommentID: 1,
			expError:     nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			postRep := NewPostRepository(db)
			testCase.mockBehaviour(mock, testCase.inComment, testCase.outCommentID)
			postID, err := postRep.InsertComment(testCase.inComment)

			assert.Equal(t, postID, testCase.outCommentID)
			assert.Equal(t, err, testCase.expError)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("were met expectation: %s", err)
			}
		})
	}
}

func Test_DeleteCommentByID(t *testing.T) {
	type mockBehaviour func(mock sqlmock.Sqlmock, inCommentID uint64)

	t.Parallel()
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	testTable := []struct {
		name          string
		mockBehaviour mockBehaviour
		inCommentID   uint64
		expError      error
	}{
		{
			name: "OK",
			mockBehaviour: func(mock sqlmock.Sqlmock, inCommentID uint64) {
				mock.ExpectBegin()

				mock.ExpectExec(`DELETE FROM comments`).
					WithArgs(inCommentID).
					WillReturnResult(driver.ResultNoRows)

				mock.ExpectCommit()
			},
			inCommentID: 1,
			expError:    nil,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			postRep := NewPostRepository(db)
			testCase.mockBehaviour(mock, testCase.inCommentID)
			err := postRep.DeleteCommentByID(testCase.inCommentID)

			assert.Equal(t, err, testCase.expError)

			if err := mock.ExpectationsWereMet(); err != nil {
				t.Errorf("were met expectation: %s", err)
			}
		})
	}
}
