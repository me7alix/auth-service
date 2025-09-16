package service

import (
	"log"
	"time"
	"user-service/application/cache"
	"user-service/application/dtos"
	"user-service/application/repository"
	"user-service/domain/entities"
	"user-service/domain/errs"

	"github.com/golang-jwt/jwt"
)

type AuthService interface {
	Login(creds dtos.Credentials) (string, error)
	Register(creds dtos.Credentials) (string, error)
	Refresh(reftok string) (string, error)
	IsAdmin(tok string) (bool, error)
	Check(tok string) error
	Get(tok string) (entities.User, error)
	GetID(tok string) (uint, error)
	Delete(tok string) error
}

type authService struct {
	userRep 	repository.UserRepository
	userCache 	cache.UserCache
	jwtSecret 	[]byte
}

func NewAuthService(
	userRep 	repository.UserRepository,
	userCache 	cache.UserCache,
	jwtSecret 	string,
) AuthService {
	return &authService{
		userRep: 	userRep,
		userCache: 	userCache,
		jwtSecret: 	[]byte(jwtSecret),
	}
}

func (a *authService) extractUserID(tokenString string) (uint, error) {
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return a.jwtSecret, nil
	})

	if err != nil {
		return 0, errs.ErrWrongToken
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			return 0, errs.ErrWrongToken
		}
	} else {
		return 0, errs.ErrWrongToken
	}

	if userID, ok := claims["userID"].(float64); ok {
		return uint(userID), nil
	} else {
		return 0, errs.ErrWrongToken
	}
}

func (a *authService) generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp":    time.Now().Add(time.Hour * 6).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.jwtSecret)
}

func (a *authService) Login(creds dtos.Credentials) (string, error) {
	userID, err := a.userRep.GetId(creds)
	if err != nil {
		return "", err
	}

	return a.generateToken(userID)
}

func (a *authService) Register(creds dtos.Credentials) (string, error) {
	user := dtos.Credentials{
		Nickname: creds.Nickname,
		Password: creds.Password,
	}

	userID, err := a.userRep.Create(user)
	if err != nil {
		return "", err
	}

	return a.generateToken(userID)
}

func (a *authService) Check(tok string) error {
	_, err := a.extractUserID(tok)
	return err
}

func (a *authService) Refresh(reftok string) (string, error) {
	userID, err := a.extractUserID(reftok)
	if err != nil {
		return "", err
	}

	return a.generateToken(userID)
}

func (a *authService) IsAdmin(tok string) (bool, error) {
	user, err := a.Get(tok)
	if err != nil {
		return false, err
	}

	return user.Role == entities.AdminRole, nil
}

func (a *authService) Get(tok string) (entities.User, error) {
	userID, err := a.extractUserID(tok)
	if err != nil {
		return entities.User{}, err
	}

	user, err := a.userCache.GetUser(userID)
	if err != nil {
		if err != errs.ErrUserNotFound {
			log.Println(err.Error())
		}

		user, err = a.userRep.GetUser(userID)
		if err != nil {
			return user, err
		}

		a.userCache.SetUser(userID, user)
		return user, nil
	}

	return user, nil
}

func (a *authService) GetID(tok string) (uint, error) {
	userID, err := a.extractUserID(tok)
	if err != nil {
		return 0, err
	}

	return userID, nil
}

func (a *authService) Delete(tok string) error {
	userID, err := a.extractUserID(tok)
	if err != nil {
		return err
	}

	return a.userRep.Delete(userID)
}
