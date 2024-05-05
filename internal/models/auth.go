package models

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/vladimirbannikov/BIP-backend/internal/models/structs"
	"github.com/vladimirbannikov/BIP-backend/internal/utils/logger"

	"github.com/gofrs/uuid"
	"github.com/golang-jwt/jwt/v5"
)

type AuthStorager interface {
	GetSecretByUserID(ctx context.Context, input structs.GetUserSecretInput) (structs.UserSecretDTO, error)
	CreateUserSecret(ctx context.Context, secretDTO structs.UserSecretDTO) error
	DeleteUserSecret(ctx context.Context, input structs.DeleteUserSecretInput) error
}

func (m *ModelAuth) RegisterUser(ctx context.Context, input structs.RegisterUserInput) error {
	if len(input.Login) == 0 || len(input.Password) == 0 {
		logger.Log(logger.ErrPrefix, fmt.Sprintf("ModelAuth: RegisterUser: logger or password is nil"))
		return ErrInvalidInput
	}
	passHash := getHash(input.Password)
	userDTO := structs.UserDTO{
		Login:        input.Login,
		PasswordHash: passHash,
	}

	_, err := m.us.CreateUser(ctx, userDTO)
	if err != nil {
		logger.Log(logger.ErrPrefix, fmt.Sprintf("ModelAuth: RegisterUser: CreateUser error: %s", err.Error()))
		return err
	}

	return nil
}

const validHoursNum = 24

func (m *ModelAuth) LoginUser(ctx context.Context, input structs.LoginUserInput) (string, error) {
	userDTO, err := m.us.GetUserByLogin(ctx, input.Login)
	if err != nil {
		logger.Log(logger.ErrPrefix, fmt.Sprintf("ModelAuth: LoginUser: GetUserByLogin error: %s", err.Error()))
		if errors.Is(err, ErrNotFound) || errors.Is(err, ErrConflict) || errors.Is(err, ErrInvalidInput) {
			return "", ErrBadCredentials
		}
		return "", err
	}
	if userDTO.PasswordHash != getHash(input.Password) {
		logger.Log(logger.ErrPrefix, fmt.Sprintf("ModelAuth: LoginUser: password is incorrect"))
		return "", ErrBadCredentials
	}

	key, err := genKey(64)
	if err != nil {
		logger.Log(logger.ErrPrefix, fmt.Sprintf("ModelAuth: LoginUser: genKey error: %s", err.Error()))
		return "", err
	}
	sessionID, err := uuid.NewV4()
	if err != nil {
		logger.Log(logger.ErrPrefix, fmt.Sprintf("ModelAuth: LoginUser: uuid.NewV4 error: %s", err.Error()))
		return "", err
	}

	payload := jwt.MapClaims{
		"sub": input.Login,
		"sID": sessionID.String(),
		"exp": time.Now().Add(time.Hour * validHoursNum).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	tokStr, err := jwtToken.SignedString([]byte(key))
	if err != nil {
		logger.Log(logger.ErrPrefix, fmt.Sprintf("ModelAuth: LoginUser: jwtToken.SignedString error: %s", err.Error()))
		return "", err
	}

	userSecret := structs.UserSecretDTO{
		Login:     userDTO.Login,
		Secret:    key,
		SessionID: sessionID.String(),
	}
	err = m.as.CreateUserSecret(ctx, userSecret)
	if err != nil {
		logger.Log(logger.ErrPrefix, fmt.Sprintf("ModelAuth: LoginUser: CreateUserSecret error: %s", err.Error()))
		return "", err
	}

	return tokStr, nil
}

func (m *ModelAuth) LogoutUser(ctx context.Context, tokenStr string) error {
	claims, err := getUnverifiedTokenClaims(tokenStr)
	if err != nil {
		logger.Log(logger.ErrPrefix, fmt.Sprintf("ModelAuth: LogoutUser: invalid token"))
		return ErrInvalidToken
	}
	login := claims["sub"].(string)
	sessionID := claims["sID"].(string)

	err = m.as.DeleteUserSecret(ctx, structs.DeleteUserSecretInput{
		Login:     login,
		SessionID: sessionID,
	})
	if err != nil {
		logger.Log(logger.ErrPrefix, fmt.Sprintf("ModelAuth: LogoutUser: DeleteUserSecret: %s", err.Error()))
		return err
	}
	return nil
}

func (m *ModelAuth) ValidateUserToken(ctx context.Context, tokenStr string) error {
	claims, err := getUnverifiedTokenClaims(tokenStr)
	if err != nil {
		logger.Log(logger.ErrPrefix, fmt.Sprintf("ModelAuth: ValidateUserToken: invalid token"))
		return ErrInvalidToken
	}
	login := claims["sub"].(string)
	sessionID := claims["sID"].(string)

	userSecret, err := m.as.GetSecretByUserID(ctx, structs.GetUserSecretInput{
		Login:     login,
		SessionID: sessionID,
	})
	if err != nil {
		logger.Log(logger.ErrPrefix, fmt.Sprintf("ModelAuth: ValidateUserToken: GetSecretByUserID error: %s", err.Error()))
		return err
	}

	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return []byte(userSecret.Secret), nil
	})

	switch {
	case token.Valid:
		return nil
	case errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet):
		logger.Log(logger.ErrPrefix, fmt.Sprintf("ModelAuth: ValidateUserToken: token expired"))
		return ErrTokenExpired
	default:
		logger.Log(logger.ErrPrefix, fmt.Sprintf("ModelAuth: ValidateUserToken: token invalid"))
		return ErrInvalidToken
	}
}

func getHash(pass string) string {
	hash := sha256.New()
	hash.Write([]byte(pass))
	return hex.EncodeToString(hash.Sum(nil))
}

func getUnverifiedTokenClaims(tokenStr string) (jwt.MapClaims, error) {
	parser := jwt.Parser{}
	unverToken, _, err := parser.ParseUnverified(tokenStr,
		jwt.MapClaims{
			"sub": "",
			"sID": "",
			"exp": "",
		})
	if err != nil {
		return nil, ErrInvalidToken
	}

	claims, ok := unverToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, ErrInvalidToken
	}
	return claims, nil
}

func genKey(length int) (string, error) {
	result := ""
	for {
		if len(result) >= length {
			return result, nil
		}
		num, err := rand.Int(rand.Reader, big.NewInt(int64(127)))
		if err != nil {
			return "", err
		}
		n := num.Int64()
		if n > 32 && n < 127 {
			result += string(n)
		}
	}
}
