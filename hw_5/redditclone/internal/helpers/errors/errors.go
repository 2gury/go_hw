package errors

import (
	"lectures-2022-1/05_web_app/99_hw/redditclone/internal/consts"
	"net/http"
)

type Error struct {
	HTTPCode    int    `json:"-"`
	Message     string `json:"message"`
}

var WrongErrorCode = &Error{
	HTTPCode:    http.StatusTeapot,
	Message:     "Технические неполадки. Уже чиним",
}

func Get(code consts.ErrorCode) *Error {
	customErr, has := Errors[code]
	if !has {
		return WrongErrorCode
	}
	return customErr
}

var Errors = map[consts.ErrorCode]*Error{
	consts.CodeInternalError: {
		HTTPCode:    http.StatusInternalServerError,
		Message:     "Что-то пошло не так",
	},
	consts.CodeBadRequest: {
		HTTPCode:    http.StatusBadRequest,
		Message:     "Неверный формат запроса",
	},
	consts.CodeUsernameAlreadyTaken: {
		HTTPCode:    http.StatusBadRequest,
		Message:     "Пользователь с таким именем уже существует",
	},
	consts.CodeIncorrectUserPassword: {
		HTTPCode:    http.StatusBadRequest,
		Message:     "Некорректный пароль",
	},
	consts.CodeUserDoesntExist: {
		HTTPCode:    http.StatusBadRequest,
		Message:     "Пользователь с таким именем не существует",
	},
}