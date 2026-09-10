package service

import (
	"time"
	"user-service/application/cache"
	"user-service/application/dtos"
	"user-service/application/repository"
	"user-service/domain/entities"
	"user-service/domain/errs"

	"github.com/phuslu/log"
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
	userRep   repository.UserRepository
	userCache cache.UserCache
	jwtSecret []byte
}

func NewAuthService(
	userRep   repository.UserRepository,
	userCache cache.UserCache,
	jwtSecret string,
) AuthService {
	return &authService{
		userRep:   userRep,
		userCache: userCache,
		jwtSecret: []byte(jwtSecret),
	}
}

func (a *authService) extractUserID(tokenString string) (uint, error) {
	log.Debug().Str("token", tokenString).Msg("received token")

	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		return a.jwtSecret, nil
	})

	if err != nil {
		log.Debug().Msg(err.Error())
		return 0, errs.ErrWrongToken
	}

	if exp, ok := claims["exp"].(float64); ok {
		if time.Now().Unix() > int64(exp) {
			log.Debug().Msg("token expired")
			return 0, errs.ErrWrongToken
		}
	} else {
		log.Debug().Msg("exp parsing error")
		return 0, errs.ErrWrongToken
	}

	if userID, ok := claims["userID"].(float64); ok {
		return uint(userID), nil
	} else {
		log.Debug().Msg("userID parsing error")
		return 0, errs.ErrWrongToken
	}
}

func (a *authService) generateToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"userID": userID,
		"exp": time.Now().Add(time.Hour * 24 * 6).Unix(),
	}

	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := t.SignedString(a.jwtSecret)
	log.Debug().Str("jwt", token).Msg("generated token")
	return token, err
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
		log.Debug().Msg(err.Error())
		return entities.User{}, err
	}

	user, err := a.userCache.GetUser(userID)
	if err != nil {
		if err != errs.ErrUserNotFound {
			log.Debug().Msg(err.Error())
			return entities.User{}, err
		}

		user, err = a.userRep.GetUser(userID)
		if err != nil {
			log.Debug().Msg(err.Error())
			return entities.User{}, err
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
