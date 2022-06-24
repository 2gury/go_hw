package post

import (
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/models"
	customErrors "lectures-2022-1/06_databases/99_hw/redditclone/internal/helpers/errors"
)

type PostUsecase interface {
	CreatePost(user *models.User, post *models.Post) (*models.Post, *customErrors.Error)
	GetPosts() ([]*models.Post, *customErrors.Error)
	GetPostsByCategory(categoryName string) ([]*models.Post, *customErrors.Error) 
	GetPostByID(postID uint64) (*models.Post, *customErrors.Error)
	AddComment(comment *models.Comment) (*models.Post, *customErrors.Error)
	DeleteComment(user *models.User, postID uint64, commentID uint64) (*models.Post, *customErrors.Error)
	UpvotePost(user *models.User, postID uint64) (*models.Post, *customErrors.Error)
	UnvotePost(user *models.User, postID uint64) (*models.Post, *customErrors.Error)
	DownvotePost(user *models.User, postID uint64) (*models.Post, *customErrors.Error)
	GetPostsByUser(userLogin string) ([]*models.Post, *customErrors.Error) 
	DeletePost(user *models.User, postID uint64) *customErrors.Error
}