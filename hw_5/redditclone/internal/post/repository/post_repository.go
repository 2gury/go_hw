package repository

import (
	"fmt"
	"lectures-2022-1/05_web_app/99_hw/redditclone/internal/models"
	"lectures-2022-1/05_web_app/99_hw/redditclone/internal/post"
	"log"
	"sort"
	"sync"
)

type PostRepository struct {
	posts         []models.Post
	mx            *sync.Mutex
	lastPostID    uint64
	lastCommentID uint64
}

func NewPostRepository() post.PostRepository {
	return &PostRepository{
		posts:         []models.Post{},
		mx:            &sync.Mutex{},
		lastPostID:    0,
		lastCommentID: 0,
	}
}

func (r *PostRepository) UpdatePost(postIndex int) {
	countUpVotes := 0
	countDownVotes := 0
	
	for _, vote := range r.posts[postIndex].Votes {
		if vote.Vote == 1 {
			countUpVotes++
		} else {
			countDownVotes++
		}
	}

	switch {
	case countUpVotes == 0 && countDownVotes == 0:
		r.posts[postIndex].UpvotePercentage = 0
	case countUpVotes == 0:
		r.posts[postIndex].UpvotePercentage = 0
	case countDownVotes == 0:
		r.posts[postIndex].UpvotePercentage = 100
	default:
		r.posts[postIndex].UpvotePercentage = float64(countUpVotes / (countUpVotes + countDownVotes))
	}

	r.posts[postIndex].Score = int64(countUpVotes + (-1) * countDownVotes)
	sort.SliceStable(r.posts, func(i, j int) bool {
        return r.posts[i].Score > r.posts[j].Score 
})
}

func (r *PostRepository) DeletePost(postID string) error {
	r.mx.Lock()
	defer r.mx.Unlock()

	postIndex := r.GetPostIndexByPostID(postID)
	log.Println(postIndex)
	if postIndex == -1 {
		return fmt.Errorf("not found post")
	}

	if len(r.posts) == 1 {
		r.posts= r.posts[:0]
	} else {
		r.posts= append(r.posts[:postIndex], r.posts[postIndex+1:]...)
	}

	return nil
}

func (r *PostRepository) SelectAuthorPost(postID string) (models.User, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	postIndex := r.GetPostIndexByPostID(postID)
	if postIndex == -1 {
		return models.User{}, fmt.Errorf("not found post")
	}

	return r.posts[postIndex].Author, nil
}

func (r *PostRepository) SelectPostsByUser(userLogin string) ([]models.Post, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	postsByUser := []models.Post{}
	for _, post := range r.posts {
		if userLogin == post.Author.Username {
			postsByUser = append(postsByUser, post)
		}
	}
	return postsByUser, nil
}

func (r *PostRepository) GetPostIndexByPostID(postID string) int {
	postIndex := -1
	for id, post := range r.posts {
		if postID == post.ID {
			postIndex = id
			break
		}
	}
	return postIndex
}

func (r *PostRepository) GetVoteIndexByUserID(userID string, postIndex int) int {
	voteIndex := -1
	for id, vote := range r.posts[postIndex].Votes {
		if vote.UserID == userID {
			voteIndex = id
			break
		}
	}
	return voteIndex
}

func (r *PostRepository) GetCommentIndexByCommentID(commentID string, postIndex int) int {
	commentIndex := -1
	for id, comment := range r.posts[postIndex].Comments {
		if commentID == comment.ID {
			commentIndex = id
			break
		}
	}
	return commentIndex
}

func (r *PostRepository) UpdateVotePost(user models.User, postID string, value int64) (models.Post, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	postIndex := r.GetPostIndexByPostID(postID)
	if postIndex == -1 {
		return models.Post{}, fmt.Errorf("not found post")
	}

	voteIndex := r.GetVoteIndexByUserID(user.ID, postIndex)
	if voteIndex == -1 {
		r.posts[postIndex].Votes = append(r.posts[postIndex].Votes, models.Vote{
			UserID: user.ID,
			Vote:   value,
		})
	} else {
		r.posts[postIndex].Votes[voteIndex].Vote = value
	}
	r.UpdatePost(postIndex)

	return r.posts[postIndex], nil
}

func (r *PostRepository) UnvotePost(user models.User, postID string) (models.Post, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	postIndex := r.GetPostIndexByPostID(postID)
	if postIndex == -1 {
		return models.Post{}, fmt.Errorf("not found post")
	}

	voteIndex := r.GetVoteIndexByUserID(user.ID, postIndex)
	if voteIndex == -1 {
		return models.Post{}, fmt.Errorf("not found vote")
	}

	if len(r.posts[postIndex].Votes) == 1 {
		r.posts[postIndex].Votes = r.posts[postIndex].Votes[:0]
	} else {
		r.posts[postIndex].Votes = append(r.posts[postIndex].Votes[:voteIndex], r.posts[postIndex].Votes[voteIndex+1:]...)
	}
	r.UpdatePost(postIndex)

	return r.posts[postIndex], nil
}

func (r *PostRepository) InsertPost(post models.Post) (uint64, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	post.ID = fmt.Sprintf("%d", r.lastPostID)
	r.posts = append(r.posts, post)
	r.lastPostID++
	return r.lastPostID - 1, nil
}

func (r *PostRepository) SelectPosts() ([]models.Post, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	return r.posts, nil
}

func (r *PostRepository) SelectPostsByCategory(catetgoryName string) ([]models.Post, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	postsByCategory := []models.Post{}
	for _, post := range r.posts {
		if catetgoryName == post.Category {
			postsByCategory = append(postsByCategory, post)
		}
	}
	return postsByCategory, nil
}

func (r *PostRepository) SelectPostByID(postID string) (models.Post, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	postIndex := r.GetPostIndexByPostID(postID)
	if postIndex == -1 {
		return models.Post{}, fmt.Errorf("not found post")
	}
	r.posts[postIndex].Views++

	return r.posts[postIndex], nil
}

func (r *PostRepository) InsertComment(postID string, comment models.Comment) (models.Post, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	postIndex := r.GetPostIndexByPostID(postID)
	if postIndex == -1 {
		return models.Post{}, fmt.Errorf("not found post")
	}

	comment.ID = fmt.Sprintf("%d", r.lastCommentID)
	r.lastCommentID++

	r.posts[postIndex].Comments = append([]models.Comment{comment}, r.posts[postIndex].Comments...)
	return r.posts[postIndex], nil
}

func (r *PostRepository) SelectAuthorComment(postID string, commentID string) (models.User, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	postIndex := r.GetPostIndexByPostID(postID)
	if postIndex == -1 {
		return models.User{}, fmt.Errorf("not found post")
	}

	commentIndex := r.GetCommentIndexByCommentID(commentID, postIndex)
	if commentIndex == -1 {
		return models.User{}, fmt.Errorf("not found post")
	}

	return r.posts[postIndex].Comments[commentIndex].Author, nil
}

func (r *PostRepository) DeleteComment(postID string, commentID string) (models.Post, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	postIndex := r.GetPostIndexByPostID(postID)
	log.Println(postIndex)
	if postIndex == -1 {
		return models.Post{}, fmt.Errorf("not found post")
	}

	commentIndex := r.GetCommentIndexByCommentID(commentID, postIndex)
	if commentIndex == -1 {
		return models.Post{}, fmt.Errorf("not found comment")
	}

	if len(r.posts[postIndex].Comments) == 1 {
		r.posts[postIndex].Comments = r.posts[postIndex].Comments[:0]
	} else {
		r.posts[postIndex].Comments = append(r.posts[postIndex].Comments[:commentIndex], r.posts[postIndex].Comments[commentIndex+1:]...)
	}

	return r.posts[postIndex], nil
}
