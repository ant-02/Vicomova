package mock

import (
	"context"

	userEntity "vicomova/internal/user/domain/entity"
)

type MockUserRepository struct {
	Users      map[string]*userEntity.User
	CreateErr  error
	GetByUNErr error
	GetByIDErr error
	UpdateErr  error
	DeleteErr  error
}

func NewMockUserRepository() *MockUserRepository {
	return &MockUserRepository{
		Users: make(map[string]*userEntity.User),
	}
}

// GetUser returns a user by username for testing
func (m *MockUserRepository) GetUser(username string) (*userEntity.User, bool) {
	u, ok := m.Users[username]
	return u, ok
}

func (m *MockUserRepository) Create(ctx context.Context, u *userEntity.User) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	u.ID = int64(len(m.Users) + 1)
	m.Users[u.Username] = u
	return nil
}

func (m *MockUserRepository) GetByUsername(ctx context.Context, username string) (*userEntity.User, error) {
	if m.GetByUNErr != nil {
		return nil, m.GetByUNErr
	}
	return m.Users[username], nil
}

func (m *MockUserRepository) GetByID(ctx context.Context, id int64) (*userEntity.User, error) {
	if m.GetByIDErr != nil {
		return nil, m.GetByIDErr
	}
	for _, u := range m.Users {
		if u.ID == id {
			return u, nil
		}
	}
	return nil, nil
}

func (m *MockUserRepository) Update(ctx context.Context, u *userEntity.User) error {
	if m.UpdateErr != nil {
		return m.UpdateErr
	}
	m.Users[u.Username] = u
	return nil
}

func (m *MockUserRepository) Delete(ctx context.Context, id int64) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	for username, u := range m.Users {
		if u.ID == id {
			delete(m.Users, username)
			break
		}
	}
	return nil
}
