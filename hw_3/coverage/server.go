package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"gitlab.com/mailru-go/lectures-2022-1/03/99_hw/coverage/internal/user"
	"gitlab.com/mailru-go/lectures-2022-1/03/99_hw/coverage/internal/user/repository"
	"gitlab.com/mailru-go/lectures-2022-1/03/99_hw/coverage/internal/user/usecase"
	"gitlab.com/mailru-go/lectures-2022-1/03/99_hw/coverage/tools"
)

type UserHandler struct {
	userUse     user.UserUsecase
	accessToken string
}

func NewUserHandler(use user.UserUsecase, token string) *UserHandler {
	return &UserHandler{
		userUse:     use,
		accessToken: token,
	}
}

func (h *UserHandler) SearchServer(w http.ResponseWriter, r *http.Request) {
	queryValues := r.URL.Query()
	limit := queryValues.Get("limit")
	offset := queryValues.Get("offset")
	query := queryValues.Get("query")
	orderField := queryValues.Get("order_field")
	orderBy := queryValues.Get("order_by")

	if r.Header.Get("AccessToken") != h.accessToken {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	queryParams, err := tools.NewQueryParams(limit, offset, query, orderField, orderBy)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	users, err := h.userUse.GetUsers()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	clientUsers := h.userUse.SortUsers(users, queryParams)

	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(clientUsers)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func main() {
	rep := repository.NewUserRepository("dataset.xml")
	use := usecase.NewUserUsecase(rep)
	hnd := NewUserHandler(use, "qwerty")

	mux := mux.NewRouter()
	mux.HandleFunc("/users/search/", hnd.SearchServer)
	log.Println("starting server at :8080 port")
	err := http.ListenAndServe(":8080", mux)
	log.Println(err)
}
