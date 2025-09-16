package cache

import "user-service/domain/entities"

type UserCache interface {
	SetUser(userID uint, user entities.User) error
	GetUser(userID uint) (entities.User, error)
}
