package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/vladimirbannikov/BIP-backend/internal/utils/logger"
	"io"
	"log"
	"net/http"

	"github.com/vladimirbannikov/BIP-backend/internal/models"
	"github.com/vladimirbannikov/BIP-backend/internal/models/structs"
)

type AuthModelManager interface {
	RegisterUser(ctx context.Context, input structs.RegisterUserInput) error
	LoginUser(ctx context.Context, input structs.LoginUserInput) (string, error)
	LogoutUser(ctx context.Context, tokenStr string) error
	ValidateUserToken(ctx context.Context, tokenStr string) error
}

type authServer struct {
	m AuthModelManager
}

func (s *authServer) Register(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		logger.Log(logger.ErrPrefix, fmt.Sprintf("Service: authServer: Register: ReadAll error: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var unm registerUserReq
	if err = json.Unmarshal(body, &unm); err != nil {
		logger.Log(logger.ErrPrefix, fmt.Sprintf("Service: authServer: Register: json.Unmarshal error: %s", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userInput := structs.RegisterUserInput{
		Login:    unm.Login,
		Password: unm.Password,
	}
	data, status := s.register(req.Context(), userInput)
	logger.Log(logger.InfoPrefix, fmt.Sprintf("Response: %v %s", status, data))
	w.WriteHeader(status)
	_, err = w.Write(data)
	if err != nil {
		return
	}
}

func (s *authServer) register(ctx context.Context, input structs.RegisterUserInput) ([]byte, int) {
	err := s.m.RegisterUser(ctx, input)
	if err != nil {
		if errors.Is(err, models.ErrInvalidInput) || errors.Is(err, models.ErrConflict) {
			return nil, http.StatusBadRequest
		}
		return nil, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}

func (s *authServer) Login(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var unm loginUserReq
	if err = json.Unmarshal(body, &unm); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	loginInput := structs.LoginUserInput{
		Login:    unm.Login,
		Password: unm.Password,
	}
	data, status := s.login(req.Context(), loginInput)
	logger.Log(logger.InfoPrefix, fmt.Sprintf("Response: %v %s", status, data))
	w.WriteHeader(status)
	_, err = w.Write(data)
	if err != nil {
		return
	}
}

func (s *authServer) login(ctx context.Context, input structs.LoginUserInput) ([]byte, int) {
	token, err := s.m.LoginUser(ctx, input)
	if err != nil {
		if errors.Is(err, models.ErrBadCredentials) {
			return nil, http.StatusUnauthorized
		}
		return nil, http.StatusInternalServerError
	}

	articleJSON, _ := json.Marshal(loginUserResp{token})
	return articleJSON, http.StatusOK
}

func (s *authServer) Logout(w http.ResponseWriter, req *http.Request) {
	tokenStr := req.Header.Get(authHeaderStr)
	if tokenStr == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	data, status := s.logout(req.Context(), tokenStr)
	logger.Log(logger.InfoPrefix, fmt.Sprintf("Response: %v %s", status, data))
	w.WriteHeader(status)
	_, err := w.Write(data)
	if err != nil {
		return
	}
}

func (s *authServer) logout(ctx context.Context, tokenStr string) ([]byte, int) {
	err := s.m.LogoutUser(ctx, tokenStr)
	if err != nil {
		if errors.Is(err, models.ErrInvalidToken) {
			return nil, http.StatusUnauthorized
		}
		if errors.Is(err, models.ErrNotFound) {
			return nil, http.StatusNotFound
		}
		if errors.Is(err, models.ErrConflict) || errors.Is(err, models.ErrInvalidInput) {
			return nil, http.StatusBadRequest
		}
		return nil, http.StatusInternalServerError
	}
	return nil, http.StatusOK
}
