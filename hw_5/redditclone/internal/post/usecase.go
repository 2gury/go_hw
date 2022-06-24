package post

import (
	"lectures-2022-1/05_web_app/99_hw/redditclone/internal/models"
	customErrors "lectures-2022-1/05_web_app/99_hw/redditclone/internal/helpers/errors"
)

type PostUsecase interface {
	CreatePost(post models.Post) (uint64, *customErrors.Error)
	GetPosts() ([]models.Post, error)
	GetPostsByCategory(catetgoryName string) ([]models.Post, error)
	GetPostByID(postID string) (models.Post, error)
	AddCommentToPost(postID string, comment models.Comment) (models.Post, error)
	DeleteComment(user models.User, postID string, commentID string) (models.Post, error)
	UpvotePost(user models.User, postID string) (models.Post, error)
	UnvotePost(user models.User, postID string) (models.Post, error)
	DownvotePost(user models.User, postID string) (models.Post, error)
	GetPostsByUser(userLogin string) ([]models.Post, error)
	DeletePost(user models.User, postID string) error
}