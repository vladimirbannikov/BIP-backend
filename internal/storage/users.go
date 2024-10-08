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
	GetUserProfileByLogin(ctx context.Context, login string) (*structs.UserProfile, error)
	UpdateUser(ctx context.Context, user *structs.UserDTO) error
	CreateUserProfile(ctx context.Context, user *structs.UserProfileDTO) (int, error)
}

type AvatarsRepo interface {
	GetAvatarByte(avatar string) ([]byte, error)
	GetRandomAvatarFileName() (string, error)
}

type UsersStorage struct {
	usersRepo   UsersRepo
	avatarsRepo AvatarsRepo
}

func NewUsersStorage(usersRepo UsersRepo, avatarsRepo AvatarsRepo) UsersStorage {
	return UsersStorage{usersRepo: usersRepo, avatarsRepo: avatarsRepo}
}

// CreateUser user
func (s *UsersStorage) CreateUser(ctx context.Context, user structs.UserDTO) (int, error) {
	avatar, err := s.avatarsRepo.GetRandomAvatarFileName()
	if err != nil {
		return 0, err
	}
	id, err := s.usersRepo.CreateUser(ctx, &user)
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateKey) {
			return 0, models.ErrConflict
		}
		return 0, err
	}
	id, err = s.usersRepo.CreateUserProfile(ctx, &structs.UserProfileDTO{
		Login:      user.Login,
		Email:      "",
		TotalScore: 0,
		TestCount:  0,
		AvatarFile: avatar,
	})
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

func (s *UsersStorage) GetUserProfileByLogin(ctx context.Context, login string) (structs.UserProfile, error) {
	profile, err := s.usersRepo.GetUserProfileByLogin(ctx, login)
	if err != nil {
		if errors.Is(err, repository.ErrObjectNotFound) {
			return structs.UserProfile{}, models.ErrNotFound
		}
		return structs.UserProfile{}, err
	}
	avatarByte, err := s.avatarsRepo.GetAvatarByte(string(profile.Avatar))
	if err != nil {
		return structs.UserProfile{}, err
	}
	profile.Avatar = avatarByte
	return *profile, nil
}

func (s *UsersStorage) UpdateUser(ctx context.Context, user structs.UserDTO) error {
	err := s.usersRepo.UpdateUser(ctx, &structs.UserDTO{
		Login:        user.Login,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
	})
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateKey) {
			return models.ErrConflict
		}
		return err
	}
	return nil
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
