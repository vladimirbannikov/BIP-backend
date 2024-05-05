package storage

import (
	"context"
	"errors"
	"github.com/vladimirbannikov/BIP-backend/internal/models"
	"github.com/vladimirbannikov/BIP-backend/internal/models/structs"
	"github.com/vladimirbannikov/BIP-backend/internal/storage/repository"
)

type UsersRepo interface {
	CreateUser(ctx context.Context, user *structs.UserDTO) (int, error)
	GetUserByLogin(ctx context.Context, login string) (*structs.UserDTO, error)
	DeleteUserByLogin(ctx context.Context, login string) error
}

type UsersStorage struct {
	usersRepo UsersRepo
}

func NewUsersStorage(usersRepo UsersRepo) UsersStorage {
	return UsersStorage{usersRepo: usersRepo}
}

// CreateUser user
func (s *UsersStorage) CreateUser(ctx context.Context, user structs.UserDTO) (int, error) {
	id, err := s.usersRepo.CreateUser(ctx, &user)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateKey) {
			return 0, models.ErrConflict
		}
		return 0, err
	}
	return id, nil
}

// GetUserByLogin user
func (s *UsersStorage) GetUserByLogin(ctx context.Context, login string) (structs.UserDTO, error) {
	user, err := s.usersRepo.GetUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return structs.UserDTO{}, models.ErrNotFound
		}
		return structs.UserDTO{}, err
	}
	return *user, nil
}

// DeleteUserByLogin user
func (s *UsersStorage) DeleteUserByLogin(ctx context.Context, login string) error {
	err := s.usersRepo.DeleteUserByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return models.ErrNotFound
		}
		return err
	}
	return nil
}
