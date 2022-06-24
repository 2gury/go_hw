package post

import "lectures-2022-1/05_web_app/99_hw/redditclone/internal/models"

type PostRepository interface {
	InsertPost(post models.Post) (uint64, error)
	SelectPosts() ([]models.Post, error)
	SelectPostsByCategory(catetgoryName string) ([]models.Post, error)
	SelectPostByID(postID string) (models.Post, error)
	InsertComment(postID string, comment models.Comment) (models.Post, error)
	SelectAuthorComment(postID string, commentID string) (models.User, error)
	DeleteComment(postID string, commentID string) (models.Post, error)
	UpdateVotePost(user models.User, postID string, value int64) (models.Post, error)
	UnvotePost(user models.User, postID string) (models.Post, error)
	SelectPostsByUser(userLogin string) ([]models.Post, error)
	SelectAuthorPost(postID string) (models.User, error)
	DeletePost(postID string) error
}