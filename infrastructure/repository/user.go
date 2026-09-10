package repository

import (
	"fmt"
	"database/sql"
	"user-service/application/dtos"
	"user-service/application/repository"
	"user-service/domain/entities"
	"user-service/domain/errs"
	db "user-service/infrastructure"

	"github.com/phuslu/log"
	"github.com/lib/pq"
)

type userRepository struct {
	dbClient *db.DBClient
}

func NewUserRepository(databaseURL string) repository.UserRepository {
	dbClient := db.NewDBClient(databaseURL)
	if dbClient == nil {
		panic(fmt.Errorf("wrong databse url"))
	}

	userRep := &userRepository{
		dbClient: dbClient,
	}

	if err := userRep.Init(); err != nil {
		panic(err)
	}

	return userRep
}

func handlePQError(err error) error {
	if err == nil { return nil }
	if pqErr, ok := err.(*pq.Error); ok {
		switch pqErr.Code.Name() {
		case "unique_violation":
			return errs.ErrNicknameAlreadyExists
		default:
			log.Debug().Msgf(err.Error())
			return errs.ErrInternalServer
		}
	}
	log.Debug().Msgf(err.Error())
	return errs.ErrInternalServer
}

func (r *userRepository) Init() error {
	request := `
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		role VARCHAR(32) NOT NULL,
		nickname VARCHAR(32) UNIQUE,
		password VARCHAR(60) NOT NULL,
		pfp VARCHAR(256) DEFAULT ''
	)`

	_, err := r.dbClient.Client.Exec(request)
	return err
}

func (r *userRepository) Create(user dtos.Credentials) (uint, error) {
	const sqlStmt = `
	INSERT INTO users (role, nickname, password)
	VALUES ($1, $2, $3)
	RETURNING id`
	var id uint
	err := r.dbClient.Client.QueryRow(sqlStmt, entities.UserRole, user.Nickname, user.Password).Scan(&id)
	if err != nil { return 0, handlePQError(err) }
	return id, nil
}

func (r *userRepository) Delete(userID uint) error {
	if _, err := r.dbClient.Client.Exec(
		"DELETE FROM users WHERE id = $1", userID,
	); err != nil {
		return errs.ErrInternalServer
	}
	return nil
}

func (r *userRepository) GetId(user dtos.Credentials) (uint, error) {
	var userID uint
	err := r.dbClient.Client.QueryRow(
		"SELECT id FROM users WHERE nickname = $1 AND password = $2",
		user.Nickname, user.Password,
	).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, errs.ErrInvalidCredentials
		}
		return 0, errs.ErrInternalServer
	}
	return userID, nil
}

func (r *userRepository) GetUser(userID uint) (entities.User, error) {
	var u entities.User
	err := r.dbClient.Client.QueryRow(
		`SELECT id, role, nickname, password, pfp
		 FROM users
		 WHERE id = $1`, userID,
	).Scan(&u.ID, &u.Role, &u.Nickname, &u.Password, &u.ProfilePicture)
	if err != nil {
		if err == sql.ErrNoRows {
			return entities.User{}, errs.ErrInvalidCredentials
		}
		return entities.User{}, errs.ErrInternalServer
	}
	return u, nil
}

func (r *userRepository) UpdateNickname(userID uint, nickname string) error {
	_, err := r.dbClient.Client.Exec(
		"UPDATE users SET nickname = $1 WHERE id = $2",
		nickname, userID)
	return handlePQError(err)
}

func (r *userRepository) UpdatePassword(userID uint, password string) error {
	_, err := r.dbClient.Client.Exec(
		"UPDATE users SET password = $1 WHERE id = $2",
		password, userID)
	return handlePQError(err)
}

func (r *userRepository) UpdateProfilePicture(userID uint, pfp string) error {
	_, err := r.dbClient.Client.Exec(
		"UPDATE users SET pfp = $1 WHERE id = $2",
		pfp, userID)
	return handlePQError(err)
}
