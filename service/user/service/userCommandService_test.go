package service

import (
	"context"
	"testing"

	"github.com/MamangRust/monolith-ecommerce-grpc-user/repository"
	db "github.com/MamangRust/monolith-ecommerce-pkg/database/schema"
	"github.com/MamangRust/monolith-ecommerce-pkg/hash"
	"github.com/MamangRust/monolith-ecommerce-pkg/logger"
	"github.com/MamangRust/monolith-ecommerce-shared/domain/requests"
	"github.com/MamangRust/monolith-ecommerce-shared/observability"
	"go.uber.org/zap"
)

type passwordUpdateRepositoryStub struct {
	password string
	userID   int
}

func (s *passwordUpdateRepositoryStub) Create(context.Context, *requests.CreateUserRequest) (*db.CreateUserRow, error) {
	return nil, nil
}

func (s *passwordUpdateRepositoryStub) Update(context.Context, *requests.UpdateUserRequest) (*db.User, error) {
	return nil, nil
}

func (s *passwordUpdateRepositoryStub) UpdateIsVerified(context.Context, int, bool) (*db.UpdateUserIsVerifiedRow, error) {
	return nil, nil
}

func (s *passwordUpdateRepositoryStub) UpdatePassword(_ context.Context, userID int, password string) (*db.UpdateUserPasswordRow, error) {
	s.userID = userID
	s.password = password
	return &db.UpdateUserPasswordRow{
		UserID: int32(userID),
		Email:  "user@example.com",
	}, nil
}

func (s *passwordUpdateRepositoryStub) Trash(context.Context, int) (*db.TrashUserRow, error) {
	return nil, nil
}

func (s *passwordUpdateRepositoryStub) Restore(context.Context, int) (*db.RestoreUserRow, error) {
	return nil, nil
}

func (s *passwordUpdateRepositoryStub) DeletePermanent(context.Context, int) (bool, error) {
	return false, nil
}

func (s *passwordUpdateRepositoryStub) RestoreAll(context.Context) (bool, error) {
	return false, nil
}

func (s *passwordUpdateRepositoryStub) DeleteAll(context.Context) (bool, error) {
	return false, nil
}

func TestUserCommandServiceUpdatePasswordHashesPassword(t *testing.T) {
	log := &logger.Logger{Log: zap.NewNop()}
	obs, err := observability.NewObservability("user-service-test", log)
	if err != nil {
		t.Fatalf("NewObservability() error = %v", err)
	}

	repo := &passwordUpdateRepositoryStub{}
	hasher := hash.NewHashingPassword()
	svc := &userCommandService{
		observability:         obs,
		userCommandRepository: repository.UserCommandRepository(repo),
		logger:                log,
		hashing:               hasher,
	}

	const (
		userID   = 42
		password = "new-reset-password"
	)

	got, err := svc.UpdatePassword(context.Background(), userID, password)
	if err != nil {
		t.Fatalf("UpdatePassword() error = %v", err)
	}
	if got == nil || got.UserID != userID {
		t.Fatalf("UpdatePassword() result = %#v, want user %d", got, userID)
	}
	if repo.userID != userID {
		t.Fatalf("repository received user ID %d, want %d", repo.userID, userID)
	}
	if repo.password == password {
		t.Fatal("repository received the plaintext password")
	}
	if err := hasher.ComparePassword(repo.password, password); err != nil {
		t.Fatalf("repository password is not a valid hash of the plaintext: %v", err)
	}
}
