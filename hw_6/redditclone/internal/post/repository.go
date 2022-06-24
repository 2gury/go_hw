package post

import "lectures-2022-1/06_databases/99_hw/redditclone/internal/models"

type PostRepository interface {
	InsertPost(post *models.Post) (uint64, error)
	SelectAllPosts() ([]*models.Post, error)
	SelectAllComments() ([]*models.Comment, error)
	SelectAllVotes() ([]*models.Vote, error) 
	// SelectPostsByCategory(catetgoryName string) ([]models.Post, error)
	SelectPostByID(postID uint64) (*models.Post, error)
	InsertComment(comment *models.Comment) (uint64, error)
	InsertVote(vote *models.Vote) (uint64, error)
	// SelectAuthorComment(postID string, commentID string) (models.User, error)
	DeleteCommentByID(commentID uint64) error
	// UnvotePost(user models.User, postID string) (models.Post, error)
	SelectAuthorPost(postID uint64) (*models.User, error)
	SelectAuthorComment(commentID uint64) (*models.User, error)
	DeletePostByID(postID uint64) error
	SelectCommentsByPostID(postID uint64) ([]*models.Comment, error)
	SelectVotesByPostID(postID uint64) ([]*models.Vote, error)
	SelectPostsByCategory(categoryName string) ([]*models.Post, error)
	SelectCommentsByCategory(categoryName string) ([]*models.Comment, error)
	SelectVotesByCategory(categoryName string) ([]*models.Vote, error)
	DeleteVoteFromPostByUserID(userID uint64, postID uint64) error

	SelectPostsByUsername(userLogin string) ([]*models.Post, error)
	SelectCommentsByUsername(userLogin string) ([]*models.Comment, error)
	SelectVotesByUsername(userLogin string) ([]*models.Vote, error)
}