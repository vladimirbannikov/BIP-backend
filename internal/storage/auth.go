package storage

import (
	"context"
	"errors"
	"github.com/vladimirbannikov/BIP-backend/internal/models"
	"github.com/vladimirbannikov/BIP-backend/internal/models/structs"
	"github.com/vladimirbannikov/BIP-backend/internal/storage/repository"
)

type AuthRepo interface {
	CreateUserSecret(ctx context.Context, secretDTO *structs.UserSecretDTO) error
	DeleteUserSecret(ctx context.Context, input structs.DeleteUserSecretInput) error
	GetSecretByUserID(ctx context.Context, input structs.GetUserSecretInput) (*structs.UserSecretDTO, error)
}

type AuthStorage struct {
	authRepo AuthRepo
}

func NewAuthStorage(authRepo AuthRepo) AuthStorage {
	return AuthStorage{authRepo: authRepo}
}

// GetSecretByUserID secret
func (s *AuthStorage) GetSecretByUserID(ctx context.Context, input structs.GetUserSecretInput) (structs.UserSecretDTO, error) {
	secretDTO, err := s.authRepo.GetSecretByUserID(ctx, input)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return structs.UserSecretDTO{}, models.ErrNotFound
		}
		return structs.UserSecretDTO{}, err
	}
	return *secretDTO, err
}

// CreateUserSecret secret
func (s *AuthStorage) CreateUserSecret(ctx context.Context, secretDTO structs.UserSecretDTO) error {
	err := s.authRepo.CreateUserSecret(ctx, &secretDTO)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateKey) {
			return models.ErrConflict
		}
		return err
	}
	return nil
}

// DeleteUserSecret secret
func (s *AuthStorage) DeleteUserSecret(ctx context.Context, input structs.DeleteUserSecretInput) error {
	err := s.authRepo.DeleteUserSecret(ctx, input)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return models.ErrNotFound
		}
		return err
	}
	return nil
}
