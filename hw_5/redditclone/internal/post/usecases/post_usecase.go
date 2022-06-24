package usecases

import (
	"fmt"
	"lectures-2022-1/05_web_app/99_hw/redditclone/internal/consts"
	customErrors "lectures-2022-1/05_web_app/99_hw/redditclone/internal/helpers/errors"
	"lectures-2022-1/05_web_app/99_hw/redditclone/internal/models"
	"lectures-2022-1/05_web_app/99_hw/redditclone/internal/post"
	"strings"
)

type PostUsecase struct {
	postRep post.PostRepository
}

func NewPostUsecase(rep post.PostRepository) post.PostUsecase {
	return &PostUsecase{
		postRep: rep,
	}
}

func (u *PostUsecase) CreatePost(post models.Post) (uint64, *customErrors.Error) {
	if strings.ReplaceAll(post.Title, " ", "") == "" || (strings.ReplaceAll(post.Text, " ", "") == "" && strings.ReplaceAll(post.URL, " ", "") == "") {
		return 0, customErrors.Get(consts.CodeBadRequest)
	}

	postID, err := u.postRep.InsertPost(post)
	if err != nil {
		return 0, customErrors.Get(consts.CodeInternalError)
	}
	return postID, nil
}

func (u *PostUsecase) GetPosts() ([]models.Post, error) {
	return u.postRep.SelectPosts()
}

func (u *PostUsecase) GetPostsByCategory(categoryName string) ([]models.Post, error) {
	return u.postRep.SelectPostsByCategory(categoryName)
}

func (u *PostUsecase) GetPostByID(postID string) (models.Post, error) {
	return u.postRep.SelectPostByID(postID)
}

func (u *PostUsecase) AddCommentToPost(postID string, comment models.Comment) (models.Post, error) {
	return u.postRep.InsertComment(postID, comment)
}

func (u *PostUsecase) DeleteComment(user models.User, postID string, commentID string) (models.Post, error) {
	authorComment, err := u.postRep.SelectAuthorComment(postID, commentID)
	if err != nil {
		return models.Post{}, err
	}
	if authorComment.ID != user.ID {
		return models.Post{}, fmt.Errorf("user can delete only own comments")
	}

	post, err := u.postRep.DeleteComment(postID, commentID)
	if err != nil {
		return models.Post{}, err
	}

	return post, nil
}

func (u *PostUsecase) UpvotePost(user models.User, postID string) (models.Post, error) {
	return u.postRep.UpdateVotePost(user, postID, 1)
}

func (u *PostUsecase) UnvotePost(user models.User, postID string) (models.Post, error) {
	return u.postRep.UnvotePost(user, postID)
}

func (u *PostUsecase) DownvotePost(user models.User, postID string) (models.Post, error) {
	return u.postRep.UpdateVotePost(user, postID, -1)
}

func (u *PostUsecase) GetPostsByUser(userLogin string) ([]models.Post, error) {
	return u.postRep.SelectPostsByUser(userLogin)
}

func (u *PostUsecase) DeletePost(user models.User, postID string) error {
	authorPost, err := u.postRep.SelectAuthorPost(postID)
	if err != nil {
		return err
	}
	if authorPost.ID != user.ID {
		return fmt.Errorf("user can delete only own posts")
	}

	err = u.postRep.DeletePost(postID)
	if err != nil {
		return err
	}

	return nil
}
