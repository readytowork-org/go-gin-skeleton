package user

import (
	"errors"
	"testing"

	"boilerplate-api/api/user/user"

	"gorm.io/gorm"
)

func TestService_GetOneUser_PassesThroughError(t *testing.T) {
	wantErr := errors.New("not found")
	repo := &MockRepository{
		GetOneUserFn: func(id int64) (GetUserResponse, error) {
			if id != 42 {
				t.Fatalf("unexpected id: %d", id)
			}
			return GetUserResponse{}, wantErr
		},
	}

	svc := NewService(repo)
	_, err := svc.GetOneUser(42)
	if !errors.Is(err, wantErr) {
		t.Fatalf("expected wrapped error, got %v", err)
	}
}

func TestService_CreateUser_CallsRepository(t *testing.T) {
	var called bool
	repo := &MockRepository{
		GetOneUserWithEmailFn: func(email string) (user.CUser, error) {
			return user.CUser{}, gorm.ErrRecordNotFound
		},
		GetOneUserWithPhoneFn: func(phone string) (user.CUser, error) {
			return user.CUser{}, gorm.ErrRecordNotFound
		},
		CreateFn: func(u user.CUser) error {
			called = true
			return nil
		},
	}

	svc := NewService(repo)
	req := CreateUserRequestData{ConfirmPassword: "secret"}
	req.Password = "secret"
	if err := svc.CreateUser(req); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected repository.Create to be invoked")
	}
}

func TestService_CreateUser_PasswordMismatch(t *testing.T) {
	svc := NewService(&MockRepository{})
	req := CreateUserRequestData{ConfirmPassword: "other"}
	req.Password = "secret"

	if err := svc.CreateUser(req); !errors.Is(err, ErrPasswordMismatch) {
		t.Fatalf("expected ErrPasswordMismatch, got %v", err)
	}
}

func TestService_CreateUser_DuplicateEmail(t *testing.T) {
	repo := &MockRepository{
		GetOneUserWithEmailFn: func(email string) (user.CUser, error) {
			return user.CUser{}, nil
		},
	}
	svc := NewService(repo)
	req := CreateUserRequestData{ConfirmPassword: "secret"}
	req.Password = "secret"

	if err := svc.CreateUser(req); !errors.Is(err, ErrEmailAlreadyExists) {
		t.Fatalf("expected ErrEmailAlreadyExists, got %v", err)
	}
}

func TestService_CreateUser_DuplicatePhone(t *testing.T) {
	repo := &MockRepository{
		GetOneUserWithEmailFn: func(email string) (user.CUser, error) {
			return user.CUser{}, gorm.ErrRecordNotFound
		},
		GetOneUserWithPhoneFn: func(phone string) (user.CUser, error) {
			return user.CUser{}, nil
		},
	}
	svc := NewService(repo)
	req := CreateUserRequestData{ConfirmPassword: "secret"}
	req.Password = "secret"

	if err := svc.CreateUser(req); !errors.Is(err, ErrPhoneAlreadyExists) {
		t.Fatalf("expected ErrPhoneAlreadyExists, got %v", err)
	}
}
