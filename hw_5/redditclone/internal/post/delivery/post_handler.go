package delivery

import (
	"encoding/json"
	"fmt"
	"lectures-2022-1/05_web_app/99_hw/redditclone/internal/consts"
	customErrors "lectures-2022-1/05_web_app/99_hw/redditclone/internal/helpers/errors"
	customContext "lectures-2022-1/05_web_app/99_hw/redditclone/internal/helpers/context"
	"lectures-2022-1/05_web_app/99_hw/redditclone/internal/models"
	"lectures-2022-1/05_web_app/99_hw/redditclone/internal/mwares"
	"lectures-2022-1/05_web_app/99_hw/redditclone/internal/post"
	"lectures-2022-1/05_web_app/99_hw/redditclone/tools/response"
	"log"
	"net/http"
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

func (h *PostHandler) Configure(mx *mux.Router) {
	mx.HandleFunc("/api/posts/", h.GetPosts()).Methods(http.MethodGet)
	mx.HandleFunc("/api/posts/{categoryName}", h.GetPostsByCategory()).Methods(http.MethodGet)
	mx.HandleFunc("/api/post/{postID:[0-9]+}", h.GetDetailedPost()).Methods(http.MethodGet)
	mx.HandleFunc("/api/user/{userLogin}", h.GetPostsByUser()).Methods(http.MethodGet)

	customMux := mx.PathPrefix("/api").Subrouter()
	customMux.Use(mwares.CheckAuth)
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

		err := h.postUse.DeletePost(*user, postID)
		if err != nil {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeInternalError))
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

		posts, err := h.postUse.GetPostsByUser(userLogin)
		if err != nil {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeInternalError))
			return
		}
		response.WriteStatusCode(w, ctx, http.StatusOK)
		err = json.NewEncoder(w).Encode(posts)
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

		post, err := h.postUse.DownvotePost(*user, postID)
		if err != nil {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeInternalError))
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

		post, err := h.postUse.UnvotePost(*user, postID)
		if err != nil {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeInternalError))
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

		post, err := h.postUse.UpvotePost(*user, postID)
		if err != nil {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeInternalError))
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
		commentID, ok := params["commentID"]
		if !ok {
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

		post, err := h.postUse.DeleteComment(*user, postID, commentID)
		if err != nil {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeInternalError))
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
			ID: userID,
			Username: username,
		}

		comment := models.Comment{
			Author:    *user,
			Body:      req.Comment,
			CreatedAt: time.Now().UTC().String(),
		}

		post, err := h.postUse.AddCommentToPost(postID, comment)
		if err != nil {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeInternalError))
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

		post, err := h.postUse.GetPostByID(postID)
		if err != nil {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeInternalError))
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
		posts, err := h.postUse.GetPosts()
		if err != nil {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeInternalError))
			return
		}
		response.WriteStatusCode(w, ctx, http.StatusOK)
		err = json.NewEncoder(w).Encode(posts)
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

		posts, err := h.postUse.GetPostsByCategory(categoryName)
		if err != nil {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeInternalError))
			return
		}
		response.WriteStatusCode(w, ctx, http.StatusOK)
		err = json.NewEncoder(w).Encode(posts)
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
			ID: userID,
			Username: username,
		}

		post := models.Post{
			Author:           *user,
			Category:         req.Category,
			Comments:         []models.Comment{},
			CreatedAt:        time.Now().UTC().String(),
			Score:            1,
			Text:             req.Text,
			Title:            req.Title,
			Type:             req.Type,
			UpvotePercentage: 100,
			URL:              req.URL,
			Views:            0,
			Votes: []models.Vote{
				{
					UserID: user.ID,
					Vote:   1,
				},
			},
		}
		postID, cusErr := h.postUse.CreatePost(post)
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		post.ID = fmt.Sprintf("%d", postID)
		response.WriteStatusCode(w, ctx, http.StatusCreated)
		err = json.NewEncoder(w).Encode(post)
		if err != nil {
			log.Println(err)
		}
	}
}
