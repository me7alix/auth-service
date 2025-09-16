package service

import (
	"user-service/application/repository"
)

type UserService interface {
	UpdateNickname(userID uint, nickname string) error
	UpdatePassword(userID uint, password string) error
	UpdateProfilePicture(userID uint, pfp string) error
}

type userService struct {
	userRep repository.UserRepository
}

func NewUserService(userRep repository.UserRepository) UserService {
	return &userService{
		userRep: userRep,
	}
}

func (u *userService) UpdateNickname(userID uint, nickname string) error {
	return u.userRep.UpdateNickname(userID, nickname)
}

func (u *userService) UpdatePassword(userID uint, password string) error {
	return u.userRep.UpdatePassword(userID, password)
}

func (u *userService) UpdateProfilePicture(userID uint, pfp string) error {
	return u.userRep.UpdateProfilePicture(userID, pfp)
}
