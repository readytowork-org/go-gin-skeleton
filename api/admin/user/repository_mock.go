package user

import (
	"boilerplate-api/api/user/user"

	"gorm.io/gorm"
)

// MockRepository is a hand-rolled mock that satisfies UserRepository. Keep it
// alongside the production code so it stays in sync with the interface. Use
// it from any package that needs to construct a Service without a real DB.
type MockRepository struct {
	CreateFn              func(u user.CUser) error
	GetAllUsersFn         func(p Pagination) ([]GetUserResponse, int64, error)
	GetOneUserFn          func(id int64) (GetUserResponse, error)
	GetOneUserWithEmailFn func(email string) (user.CUser, error)
	GetOneUserWithPhoneFn func(phone string) (user.CUser, error)
}

func (m *MockRepository) WithTrx(*gorm.DB) UserRepository { return m }

func (m *MockRepository) Create(u user.CUser) error {
	if m.CreateFn != nil {
		return m.CreateFn(u)
	}
	return nil
}

func (m *MockRepository) GetAllUsers(p Pagination) ([]GetUserResponse, int64, error) {
	if m.GetAllUsersFn != nil {
		return m.GetAllUsersFn(p)
	}
	return nil, 0, nil
}

func (m *MockRepository) GetOneUser(id int64) (GetUserResponse, error) {
	if m.GetOneUserFn != nil {
		return m.GetOneUserFn(id)
	}
	return GetUserResponse{}, nil
}

func (m *MockRepository) GetOneUserWithEmail(email string) (user.CUser, error) {
	if m.GetOneUserWithEmailFn != nil {
		return m.GetOneUserWithEmailFn(email)
	}
	return user.CUser{}, nil
}

func (m *MockRepository) GetOneUserWithPhone(phone string) (user.CUser, error) {
	if m.GetOneUserWithPhoneFn != nil {
		return m.GetOneUserWithPhoneFn(phone)
	}
	return user.CUser{}, nil
}
