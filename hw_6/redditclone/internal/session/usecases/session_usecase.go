package usecases

import (
	"fmt"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/consts"
	customErrors "lectures-2022-1/06_databases/99_hw/redditclone/internal/helpers/errors"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/models"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/session"
	"lectures-2022-1/06_databases/99_hw/redditclone/tools/response"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var secret = []byte("ASF9FA00KS2MLSC3L3NKFSA9")

type SessionUsecase struct {
	sessRep session.SessionRepository
}

func NewSessionUsecase(rep session.SessionRepository) session.SessionUsecase {
	return &SessionUsecase{
		sessRep: rep,
	}
}

func (u *SessionUsecase) Create(user *models.User) (*models.Session, *customErrors.Error) {
	strToken, err := u.NewJwtSession(user)
	if err != nil {
		return nil, customErrors.Get(consts.CodeInternalError)
	}

	sess := models.NewSession(strToken)

	err = u.sessRep.Create(sess)
	if err != nil {
		return nil, customErrors.Get(consts.CodeInternalError)
	}

	return sess, nil
}

func (u *SessionUsecase) Check(sessValue string) (*models.User, *customErrors.Error) {
	sess, err := u.sessRep.Get(sessValue)
	if err != nil {
		return nil, customErrors.Get(consts.CodeInternalError)
	}

	user, err := u.CheckJwtSession(sess.Value)
	if err != nil {
		return nil, customErrors.Get(consts.CodeInternalError)
	}

	return user, nil
}

func (u *SessionUsecase) Delete(sessValue string) *customErrors.Error {
	err := u.sessRep.Delete(sessValue)
	if err != nil {
		return customErrors.Get(consts.CodeInternalError)
	}

	return nil
}

func (u *SessionUsecase) NewJwtSession(user *models.User) (string, error) {
	curTime := time.Now()
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"iat": curTime.Unix(),
		"exp": curTime.Add(consts.ExpiresDuration).Unix(),
		"user": response.Body{
			"username": user.Username,
			"id":       fmt.Sprintf("%d", user.ID),
		},
	})
	strToken, err := jwtToken.SignedString(secret)
	if err != nil {
		return "", err
	}
	return strToken, nil
}

func (u *SessionUsecase) CheckJwtSession(tokenValue string) (*models.User, error) {
	token, err := jwt.Parse(tokenValue, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unsupported signing method: %v", token.Header["alg"])
		}
		return secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userRawData := claims["user"].(map[string]interface{})
		intUserID, err := strconv.Atoi(userRawData["id"].(string))
		if err != nil {
			return nil, fmt.Errorf("jwt token is not valid")
		}
		return &models.User{
			ID:       uint64(intUserID),
			Username: userRawData["username"].(string),
		}, nil
	}
	return nil, fmt.Errorf("jwt token is not valid")

}
