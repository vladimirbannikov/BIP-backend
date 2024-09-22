package models

import (
	"context"
	"fmt"
	"github.com/vladimirbannikov/BIP-backend/internal/models/structs"
	"github.com/vladimirbannikov/BIP-backend/internal/utils/logger"
)

type UsersStorager interface {
	CreateUser(ctx context.Context, user structs.UserDTO) (int, error)
	GetUserByLogin(ctx context.Context, login string) (structs.UserDTO, error)
	DeleteUserByLogin(ctx context.Context, login string) error
	GetUserProfileByLogin(ctx context.Context, login string) (structs.UserProfile, error)
	UpdateUser(ctx context.Context, user structs.UserDTO) error
}

func (u *ModelUsers) GetUserProfileOwn(ctx context.Context, tokenStr string) (structs.UserProfile, error) {
	login, err := getUserLoginFromToken(tokenStr)
	if err != nil {
		return structs.UserProfile{}, err
	}
	return u.GetUserProfile(ctx, login)
}

func (u *ModelUsers) GetUserProfile(ctx context.Context, login string) (structs.UserProfile, error) {
	user, err := u.us.GetUserProfileByLogin(ctx, login)
	if err != nil {
		return structs.UserProfile{}, err
	}
	return structs.UserProfile{
		Login:             user.Login,
		Email:             user.Email,
		TotalScore:        user.TotalScore,
		TestCount:         user.TestCount,
		GlobalRatingPlace: user.GlobalRatingPlace,
	}, nil
}

func (u *ModelUsers) UpdateUser(ctx context.Context, tokenStr string, newUser structs.User) error {
	login, err := getUserLoginFromToken(tokenStr)
	if err != nil {
		return err
	}

	err = u.us.UpdateUser(ctx, structs.UserDTO{
		Login:        login,
		Email:        newUser.Email,
		PasswordHash: getHash(newUser.Password),
	})
	if err != nil {
		return err
	}
	return nil
}

func getUserLoginFromToken(tokenStr string) (string, error) {
	claims, err := getUnverifiedTokenClaims(tokenStr)
	if err != nil {
		logger.Log(logger.ErrPrefix, fmt.Sprintf("ModelAuth: ValidateUserToken: invalid token"))
		return "", ErrInvalidToken
	}
	login := claims["sub"].(string)
	return login, nil
}
