package repository

import (
	"user-service/application/dtos"
	"user-service/domain/entities"
)

type UserRepository interface {
	Init() error
	Create(user dtos.Credentials) (uint, error)
	GetUser(userID uint) (entities.User, error)
	GetId(creds dtos.Credentials) (uint, error)
	UpdateNickname(userID uint, nickname string) error
	UpdatePassword(userID uint, password string) error
	UpdateProfilePicture(userID uint, pfp string) error
	Delete(userID uint) error
}
