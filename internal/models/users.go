package models

import (
	"context"
	"github.com/vladimirbannikov/BIP-backend/internal/models/structs"
)

type UsersStorager interface {
	CreateUser(ctx context.Context, user structs.UserDTO) (int, error)
	GetUserByLogin(ctx context.Context, login string) (structs.UserDTO, error)
	DeleteUserByLogin(ctx context.Context, login string) error
}
