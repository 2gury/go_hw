package delivery

import (
	"encoding/json"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/consts"
	customContext "lectures-2022-1/06_databases/99_hw/redditclone/internal/helpers/context"
	customErrors "lectures-2022-1/06_databases/99_hw/redditclone/internal/helpers/errors"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/models"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/mwares"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/post"
	"lectures-2022-1/06_databases/99_hw/redditclone/tools/response"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

type PostHandler struct {
	postUse post.PostUsecase
}

func NewPostHandler(use post.PostUsecase) *PostHandler {
	return &PostHandler{
		postUse: use,
	}
}

func (h *PostHandler) Configure(mx *mux.Router, mm *mwares.MiddlewareManager) {
	mx.HandleFunc("/api/posts/", h.GetPosts()).Methods(http.MethodGet)
	mx.HandleFunc("/api/posts/{categoryName}", h.GetPostsByCategory()).Methods(http.MethodGet)
	mx.HandleFunc("/api/post/{postID:[0-9]+}", h.GetDetailedPost()).Methods(http.MethodGet)
	mx.HandleFunc("/api/user/{userLogin}", h.GetPostsByUser()).Methods(http.MethodGet)

	customMux := mx.PathPrefix("/api").Subrouter()
	customMux.Use(mm.CheckAuth)
	customMux.Path("/post/{postID:[0-9]+}").HandlerFunc(h.DeletePost()).Methods(http.MethodDelete)
	customMux.Path("/post/{postID:[0-9]+}/downvote").HandlerFunc(h.DownvotePost()).Methods(http.MethodGet)
	customMux.Path("/post/{postID:[0-9]+}/unvote").HandlerFunc(h.UnvotePost()).Methods(http.MethodGet)
	customMux.Path("/post/{postID:[0-9]+}/upvote").HandlerFunc(h.UpvotePost()).Methods(http.MethodGet)
	customMux.Path("/post/{postID:[0-9]+}/{commentID:[0-9]+}").HandlerFunc(h.DeleteComment()).Methods(http.MethodDelete)
	customMux.Path("/post/{postID:[0-9]+}").HandlerFunc(h.AddComment()).Methods(http.MethodPost)
	customMux.Path("/posts").HandlerFunc(h.AddPost()).Methods(http.MethodPost)
}

func (h *PostHandler) DeletePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		params := mux.Vars(r)
		postID, ok := params["postID"]
		if !ok {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeBadRequest))
			return
		}
		intPostID, err := strconv.Atoi(postID)
		if err != nil {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeBadRequest))
			return
		}

		userID, cusErr := customContext.GetUserID(ctx)
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		username, cusErr := customContext.GetUsername(ctx)
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		user := &models.User{
			ID:       userID,
			Username: username,
		}

		cusErr = h.postUse.DeletePost(user, uint64(intPostID))
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		response.WriteStatusCode(w, ctx, http.StatusOK)
		err = json.NewEncoder(w).Encode(&response.Body{
			"message": "success",
		})
		if err != nil {
			log.Println(err)
		}
	}
}

func (h *PostHandler) GetPostsByUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		params := mux.Vars(r)
		userLogin, ok := params["userLogin"]
		if !ok {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeBadRequest))
			return
		}

		posts, cusErr := h.postUse.GetPostsByUser(userLogin)
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		response.WriteStatusCode(w, ctx, http.StatusOK)
		err := json.NewEncoder(w).Encode(posts)
		if err != nil {
			log.Println(err)
		}
	}
}

func (h *PostHandler) DownvotePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		params := mux.Vars(r)
		postID, ok := params["postID"]
		if !ok {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeBadRequest))
			return
		}
		intPostID, err := strconv.Atoi(postID)
		if err != nil {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeBadRequest))
			return
		}

		userID, cusErr := customContext.GetUserID(ctx)
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		username, cusErr := customContext.GetUsername(ctx)
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		user := &models.User{
			ID: userID,
			Username: username,
		}

		post, cusErr := h.postUse.DownvotePost(user, uint64(intPostID))
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		response.WriteStatusCode(w, ctx, http.StatusOK)
		err = json.NewEncoder(w).Encode(post)
		if err != nil {
			log.Println(err)
		}
	}
}

func (h *PostHandler) UnvotePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		params := mux.Vars(r)
		postID, ok := params["postID"]
		if !ok {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeBadRequest))
			return
		}
		intPostID, err := strconv.Atoi(postID)
		if err != nil {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeBadRequest))
			return
		}

		userID, cusErr := customContext.GetUserID(ctx)
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		username, cusErr := customContext.GetUsername(ctx)
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		user := &models.User{
			ID: userID,
			Username: username,
		}

		post, cusErr := h.postUse.UnvotePost(user, uint64(intPostID))
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		response.WriteStatusCode(w, ctx, http.StatusOK)
		err = json.NewEncoder(w).Encode(post)
		if err != nil {
			log.Println(err)
		}
	}
}

func (h *PostHandler) UpvotePost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		params := mux.Vars(r)
		postID, ok := params["postID"]
		if !ok {
			log.Println("lol")
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeBadRequest))
			return
		}
		intPostID, err := strconv.Atoi(postID)
		if err != nil {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeBadRequest))
			return
		}

		userID, cusErr := customContext.GetUserID(ctx)
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		username, cusErr := customContext.GetUsername(ctx)
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		user := &models.User{
			ID: userID,
			Username: username,
		}

		post, cusErr := h.postUse.UpvotePost(user, uint64(intPostID))
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		response.WriteStatusCode(w, ctx, http.StatusOK)
		err = json.NewEncoder(w).Encode(post)
		if err != nil {
			log.Println(err)
		}
	}
}

func (h *PostHandler) DeleteComment() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		params := mux.Vars(r)
		postID, ok := params["postID"]
		if !ok {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeBadRequest))
			return
		}
		intPostID, err := strconv.Atoi(postID)
		if err != nil {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeBadRequest))
			return
		}

		commentID, ok := params["commentID"]
		if !ok {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeBadRequest))
			return
		}
		intCommentID, err := strconv.Atoi(commentID)
		if err != nil {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeBadRequest))
			return
		}

		userID, cusErr := customContext.GetUserID(ctx)
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		username, cusErr := customContext.GetUsername(ctx)
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		user := &models.User{
			ID: userID,
			Username: username,
		}

		post, cusErr := h.postUse.DeleteComment(user, uint64(intPostID), uint64(intCommentID))
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		response.WriteStatusCode(w, ctx, http.StatusOK)
		err = json.NewEncoder(w).Encode(post)
		if err != nil {
			log.Println(err)
		}
	}
}

func (h *PostHandler) AddComment() http.HandlerFunc {
	type Request struct {
		Comment string `json:"comment"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		params := mux.Vars(r)
		postID, ok := params["postID"]
		if !ok {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeBadRequest))
			return
		}
		intPostID, err := strconv.Atoi(postID)
		if err != nil {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeBadRequest))
			return
		}

		defer r.Body.Close()
		req := &Request{}

		err = json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeBadRequest))
			return
		}

		userID, cusErr := customContext.GetUserID(ctx)
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		username, cusErr := customContext.GetUsername(ctx)
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		user := &models.User{
			ID:       userID,
			Username: username,
		}

		comment := &models.Comment{
			UserID:    user.ID,
			PostID:    uint64(intPostID),
			Body:      req.Comment,
			CreatedAt: time.Now(),
		}

		post, cusErr := h.postUse.AddComment(comment)
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		response.WriteStatusCode(w, ctx, http.StatusOK)
		err = json.NewEncoder(w).Encode(post)
		if err != nil {
			log.Println(err)
		}
	}
}

func (h *PostHandler) GetDetailedPost() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		params := mux.Vars(r)
		postID, ok := params["postID"]
		if !ok {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeBadRequest))
			return
		}
		intPostID, err := strconv.Atoi(postID)
		if err != nil {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeBadRequest))
			return
		}

		post, cusErr := h.postUse.GetPostByID(uint64(intPostID))
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		response.WriteStatusCode(w, ctx, http.StatusOK)
		err = json.NewEncoder(w).Encode(post)
		if err != nil {
			log.Println(err)
		}
	}
}

func (h *PostHandler) GetPosts() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		posts, cusErr := h.postUse.GetPosts()
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		response.WriteStatusCode(w, ctx, http.StatusOK)
		err := json.NewEncoder(w).Encode(posts)
		if err != nil {
			log.Println(err)
		}
	}
}

func (h *PostHandler) GetPostsByCategory() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		params := mux.Vars(r)
		categoryName, ok := params["categoryName"]
		log.Println(categoryName)
		if !ok {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeBadRequest))
			return
		}

		posts, cusErr := h.postUse.GetPostsByCategory(categoryName)
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		response.WriteStatusCode(w, ctx, http.StatusOK)
		err := json.NewEncoder(w).Encode(posts)
		if err != nil {
			log.Println(err)
		}
	}
}

func (h *PostHandler) AddPost() http.HandlerFunc {
	type Request struct {
		Category string `json:"category"`
		Text     string `json:"text,omitempty"`
		Title    string `json:"title"`
		Type     string `json:"type"`
		URL      string `json:"url,omitempty"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		defer r.Body.Close()
		req := &Request{}

		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeBadRequest))
			return
		}

		userID, cusErr := customContext.GetUserID(ctx)
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		username, cusErr := customContext.GetUsername(ctx)
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		user := &models.User{
			ID:       userID,
			Username: username,
		}

		postInfo := &models.Post{
			UserID:           user.ID,
			Category:         req.Category,
			CreatedAt:        time.Now(),
			Score:            0,
			Text:             req.Text,
			Title:            req.Title,
			Type:             req.Type,
			UpvotePercentage: 100,
			URL:              req.URL,
			Views:            0,
		}
		post, cusErr := h.postUse.CreatePost(user, postInfo)
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		response.WriteStatusCode(w, ctx, http.StatusCreated)
		err = json.NewEncoder(w).Encode(post)
		if err != nil {
			log.Println(err)
		}
	}
}
