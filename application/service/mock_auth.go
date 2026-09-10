package service

import (
	"fmt"
	"user-service/domain/entities"
	"user-service/application/dtos"
)

type mockAuthService struct {
	token  string
	userID uint
}

func NewMockAuthService(
	token  string,
	userID uint,
) AuthService {
	return &mockAuthService{
		token: token,
		userID: userID,
	}
}

func (m *mockAuthService) Register(creds dtos.Credentials) (string, error) {
	return m.token, nil
}

func (m *mockAuthService) Login(creds dtos.Credentials) (string, error) {
	return m.token, nil
}

func (m *mockAuthService) Check(tok string) error {
	if tok != m.token {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func (m *mockAuthService) Delete(tok string) error {
	if tok != m.token {
		return fmt.Errorf("invalid token")
	}

	return nil
}

func (m *mockAuthService) IsAdmin(tok string) (bool, error) {
	return false, nil
}

func (m *mockAuthService) Get(tok string) (entities.User, error) {
	if tok != m.token {
		return entities.User{}, fmt.Errorf("invalid token")
	}

	return entities.User{}, nil
}

func (m *mockAuthService) Refresh(reftok string) (string, error) {
	if reftok != m.token {
		return "", fmt.Errorf("invalid token")
	}

	return "", nil
}

func (m *mockAuthService) GetID(token string) (uint, error) {
	if token != m.token {
		return 0, fmt.Errorf("invalid token")
	}

	return m.userID, nil
}
