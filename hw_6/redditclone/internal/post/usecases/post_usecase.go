package usecases

import (
	"database/sql"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/consts"
	customErrors "lectures-2022-1/06_databases/99_hw/redditclone/internal/helpers/errors"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/models"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/post"
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

func (u *PostUsecase) CreatePost(user *models.User, post *models.Post) (*models.Post, *customErrors.Error) {
	if strings.ReplaceAll(post.Title, " ", "") == "" || (strings.ReplaceAll(post.Text, " ", "") == "" && strings.ReplaceAll(post.URL, " ", "") == "") {
		return nil, customErrors.Get(consts.CodeBadRequest)
	}

	postID, err := u.postRep.InsertPost(post)
	if err != nil {
		return nil, customErrors.Get(consts.CodeInternalError)
	}

	vote := &models.Vote{
		Vote: 1,
		User: user.ID,
		UserID: user.ID,
		PostID: postID,
	}
	_, err = u.postRep.InsertVote(vote)
	if err != nil {
		return nil, customErrors.Get(consts.CodeInternalError)
	}

	post, cusErr := u.GetPostByID(postID)
	if cusErr != nil {
		return nil, cusErr
	}

	return post, nil
}

func (u *PostUsecase) GetPosts() ([]*models.Post, *customErrors.Error) {
	posts, err := u.postRep.SelectAllPosts()
	if err != nil {
		return nil, customErrors.Get(consts.CodeInternalError)
	}

	mapPosts := map[uint64]*models.Post{}
	for _, post := range posts {
		mapPosts[post.ID] = post
	}

	comments, err := u.postRep.SelectAllComments()
	if err != nil {
		return nil, customErrors.Get(consts.CodeInternalError)
	}
	for _, comment := range comments {
		mapPosts[comment.PostID].Comments = append(mapPosts[comment.PostID].Comments, comment)
	}

	votes, err := u.postRep.SelectAllVotes()
	if err != nil {
		return nil, customErrors.Get(consts.CodeInternalError)
	}
	for _, vote := range votes {
		mapPosts[vote.PostID].Votes = append(mapPosts[vote.PostID].Votes, vote)
	}

	return posts, nil
}

func (u *PostUsecase) GetPostsByCategory(categoryName string) ([]*models.Post, *customErrors.Error) {
	posts, err := u.postRep.SelectPostsByCategory(categoryName)
	if err != nil {
		return nil, customErrors.Get(consts.CodeInternalError)
	}

	mapPosts := map[uint64]*models.Post{}
	for _, post := range posts {
		mapPosts[post.ID] = post
	}

	comments, err := u.postRep.SelectCommentsByCategory(categoryName)
	if err != nil {
		return nil, customErrors.Get(consts.CodeInternalError)
	}
	for _, comment := range comments {
		mapPosts[comment.PostID].Comments = append(mapPosts[comment.PostID].Comments, comment)
	}

	votes, err := u.postRep.SelectVotesByCategory(categoryName)
	if err != nil {
		return nil, customErrors.Get(consts.CodeInternalError)
	}
	for _, vote := range votes {
		mapPosts[vote.PostID].Votes = append(mapPosts[vote.PostID].Votes, vote)
	}

	return posts, nil
}

func (u *PostUsecase) GetPostByID(postID uint64) (*models.Post, *customErrors.Error) {
	post, err := u.postRep.SelectPostByID(postID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, customErrors.Get(consts.CodePostDoesntExist)
		}
		return nil, customErrors.Get(consts.CodeInternalError)
	}

	comments, err := u.postRep.SelectCommentsByPostID(postID)
	switch {
	case err == sql.ErrNoRows:
		post.Comments = []*models.Comment{}
	case err != nil:
		return nil, customErrors.Get(consts.CodeInternalError)
	default:
		post.Comments = comments
	}

	votes, err := u.postRep.SelectVotesByPostID(postID)
	switch {
	case err == sql.ErrNoRows:
		post.Votes = []*models.Vote{}
	case err != nil:
		return nil, customErrors.Get(consts.CodeInternalError)
	default:
		post.Votes = votes
	}

	return post, nil
}

func (u *PostUsecase) AddComment(comment *models.Comment) (*models.Post, *customErrors.Error) {
	_, err := u.postRep.InsertComment(comment)
	if err != nil {
		return nil, customErrors.Get(consts.CodeInternalError)
	}

	post, cusErr := u.GetPostByID(comment.PostID)
	if cusErr != nil {
		return nil, cusErr
	}

	return post, nil
}

func (u *PostUsecase) DeleteComment(user *models.User, postID uint64, commentID uint64) (*models.Post, *customErrors.Error) {
	authorPost, err := u.postRep.SelectAuthorComment(commentID)
	if err != nil {
		return nil, customErrors.Get(consts.CodeInternalError)
	}
	if authorPost.ID != user.ID {
		return nil, customErrors.Get(consts.CodeBadAccess)
	}

	err = u.postRep.DeleteCommentByID(commentID)
	if err != nil {
		return nil, customErrors.Get(consts.CodeInternalError)
	}

	post, cusErr := u.GetPostByID(postID)
	if cusErr != nil {
		return nil, cusErr
	}

	return post, nil
}

func (u *PostUsecase) UpvotePost(user *models.User, postID uint64) (*models.Post, *customErrors.Error) {
	err := u.postRep.DeleteVoteFromPostByUserID(user.ID, postID)
	if err != nil {
		return nil, customErrors.Get(consts.CodeInternalError)
	}
	vote := &models.Vote{
		Vote: 1,
		User: user.ID,

		UserID: user.ID,
		PostID: postID,
	}
	_, err = u.postRep.InsertVote(vote)
	if err != nil {
		return nil, customErrors.Get(consts.CodeInternalError)
	}

	post, cusErr := u.GetPostByID(postID)
	if cusErr != nil {
		return nil, cusErr
	}
	return post, nil
}

func (u *PostUsecase) UnvotePost(user *models.User, postID uint64) (*models.Post, *customErrors.Error) {
	err := u.postRep.DeleteVoteFromPostByUserID(user.ID, postID)
	if err != nil {
		return nil, customErrors.Get(consts.CodeInternalError)
	}

	post, cusErr := u.GetPostByID(postID)
	if cusErr != nil {
		return nil, cusErr
	}
	return post, nil
}

func (u *PostUsecase) DownvotePost(user *models.User, postID uint64) (*models.Post, *customErrors.Error) {
	err := u.postRep.DeleteVoteFromPostByUserID(user.ID, postID)
	if err != nil {
		return nil, customErrors.Get(consts.CodeInternalError)
	}
	vote := &models.Vote{
		Vote: -1,
		User: user.ID,

		UserID: user.ID,
		PostID: postID,
	}
	_, err = u.postRep.InsertVote(vote)
	if err != nil {
		return nil, customErrors.Get(consts.CodeInternalError)
	}

	post, cusErr := u.GetPostByID(postID)
	if cusErr != nil {
		return nil, cusErr
	}
	return post, nil
}

func (u *PostUsecase) GetPostsByUser(userLogin string) ([]*models.Post, *customErrors.Error) {
	posts, err := u.postRep.SelectPostsByUsername(userLogin)
	if err != nil {
		return nil, customErrors.Get(consts.CodeInternalError)
	}

	mapPosts := map[uint64]*models.Post{}
	for _, post := range posts {
		mapPosts[post.ID] = post
	}

	comments, err := u.postRep.SelectCommentsByUsername(userLogin)
	if err != nil {
		return nil, customErrors.Get(consts.CodeInternalError)
	}
	for _, comment := range comments {
		mapPosts[comment.PostID].Comments = append(mapPosts[comment.PostID].Comments, comment)
	}

	votes, err := u.postRep.SelectVotesByUsername(userLogin)
	if err != nil {
		return nil, customErrors.Get(consts.CodeInternalError)
	}
	for _, vote := range votes {
		mapPosts[vote.PostID].Votes = append(mapPosts[vote.PostID].Votes, vote)
	}

	return posts, nil
}

func (u *PostUsecase) DeletePost(user *models.User, postID uint64) *customErrors.Error {
	authorPost, err := u.postRep.SelectAuthorPost(postID)
	if err != nil {
		return customErrors.Get(consts.CodeInternalError)
	}
	if authorPost.ID != user.ID {
		return customErrors.Get(consts.CodeBadAccess)
	}

	err = u.postRep.DeletePostByID(postID)
	if err != nil {
		return customErrors.Get(consts.CodeInternalError)
	}

	return nil
}
