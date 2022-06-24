package delivery

import (
	"encoding/json"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/consts"
	customErrors "lectures-2022-1/06_databases/99_hw/redditclone/internal/helpers/errors"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/models"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/session"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/user"
	"lectures-2022-1/06_databases/99_hw/redditclone/tools/response"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type UserHandler struct {
	userUse user.UserUsecase
	sessUse session.SessionUsecase
}

func NewUserHandler(userUsecase user.UserUsecase, sessionUsecase session.SessionUsecase) *UserHandler {
	return &UserHandler{
		userUse: userUsecase,
		sessUse: sessionUsecase,
	}
}

func (h *UserHandler) Configure(mx *mux.Router) {
	mx.HandleFunc("/api/register", h.RegiserUser()).Methods(http.MethodPost)
	mx.HandleFunc("/api/login", h.LoginUser()).Methods(http.MethodPost)
}

func (h *UserHandler) RegiserUser() http.HandlerFunc {
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
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
		user := &models.User{
			Username: req.Username,
			Password: req.Password,
		}
		userID, cusErr := h.userUse.RegiserUser(user)
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		user.ID = userID
		sess, cusErr := h.sessUse.Create(user)
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}

		response.WriteStatusCode(w, ctx, http.StatusCreated)
		err = json.NewEncoder(w).Encode(response.Body{
			"token": sess.Value,
		})

		if err != nil {
			log.Println(err)
		}
	}
}

func (h *UserHandler) LoginUser() http.HandlerFunc {
	type Request struct {
		Username string `json:"username"`
		Password string `json:"password"`
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
		user := &models.User{
			Username: req.Username,
			Password: req.Password,
		}
		userID, cusErr := h.userUse.LoginUser(user)
		log.Println(cusErr)
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}
		user.ID = userID
		sess, cusErr := h.sessUse.Create(user)
		if cusErr != nil {
			response.WriteErrorResponse(w, ctx, cusErr)
			return
		}

		response.WriteStatusCode(w, ctx, http.StatusCreated)
		err = json.NewEncoder(w).Encode(response.Body{
			"token": sess.Value,
		})
		if err != nil {
			log.Println(err)
		}
	}
}
