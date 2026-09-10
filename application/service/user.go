package service

import (
	"user-service/application/cache"
	"user-service/application/repository"
)

type UserService interface {
	UpdateNickname(userID uint, nickname string) error
	UpdatePassword(userID uint, password string) error
	UpdateProfilePicture(userID uint, pfp string) error
}

type userService struct {
	userRep   repository.UserRepository
	userCache cache.UserCache
}

func NewUserService(
	userRep repository.UserRepository,
	userCache cache.UserCache,
) UserService {
	return &userService{
		userRep: userRep,
		userCache: userCache,
	}
}

func (u *userService) refreshCache(userID uint) error {
	user, err := u.userRep.GetUser(userID)
	if err != nil { return err }
	return u.userCache.SetUser(userID, user)
}

func (u *userService) UpdateNickname(userID uint, nickname string) error {
	err := u.userRep.UpdateNickname(userID, nickname)
	if err != nil { return err }
	return u.refreshCache(userID)
}

func (u *userService) UpdatePassword(userID uint, password string) error {
	err := u.userRep.UpdatePassword(userID, password)
	if err != nil { return err }
	return u.refreshCache(userID)
}

func (u *userService) UpdateProfilePicture(userID uint, pfp string) error {
	err := u.userRep.UpdateProfilePicture(userID, pfp)
	if err != nil { return err }
	return u.refreshCache(userID)
}
