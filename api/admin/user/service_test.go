package user

import (
	"errors"
	"testing"

	"boilerplate-api/api/user/user"
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
		CreateFn: func(u user.CUser) error {
			called = true
			return nil
		},
	}

	svc := NewService(repo)
	if err := svc.CreateUser(user.CUser{}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !called {
		t.Fatal("expected repository.Create to be invoked")
	}
}
