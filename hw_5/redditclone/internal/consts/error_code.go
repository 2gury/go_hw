package consts

type ErrorCode uint16

const (
	CodeInternalError ErrorCode = 101 + iota
	CodeBadRequest
	CodeUsernameAlreadyTaken
	CodeIncorrectUserPassword
	CodeUserDoesntExist
)