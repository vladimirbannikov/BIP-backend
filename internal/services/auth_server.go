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
	RegisterUser(ctx context.Context, input structs.RegisterUserInput) (qr2fa string, err error)
	LoginUser(ctx context.Context, input structs.LoginUserInput) (string, error)
	LogoutUser(ctx context.Context, tokenStr string) error
	ValidateUserToken(ctx context.Context, tokenStr string) error
	CheckUser2FA(ctx context.Context, tokenStr string) error
	Validate2FACode(ctx context.Context, code string, tokenStr string) (valid bool, err error)
	// если просрочен 2fa код или его вообще нет, то возвращать, чтобы фронт дергал ручку логина
	// проверки правильности введенного кода (который ввел клиент), потом кидать запрос в гугл аутентификатор
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
		Email:    unm.Email,
	}
	data, status := s.register(req.Context(), userInput)
	logger.Log(logger.InfoPrefix, fmt.Sprintf("Response: %v %s", status))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "multipart/form-data")
	w.WriteHeader(status)
	_, err = w.Write(data)
	if err != nil {
		return
	}
}

func (s *authServer) register(ctx context.Context, input structs.RegisterUserInput) ([]byte, int) {
	qr, err := s.m.RegisterUser(ctx, input)
	if err != nil {
		if errors.Is(err, models.ErrInvalidInput) || errors.Is(err, models.ErrConflict) {
			return nil, http.StatusBadRequest
		}
		return nil, http.StatusInternalServerError
	}
	return []byte(qr), http.StatusOK
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
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
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
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

func (s *authServer) User2FA(w http.ResponseWriter, req *http.Request) {
	tokenStr := req.Header.Get(authHeaderStr)
	if tokenStr == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var unm User2FAInput
	if err := json.Unmarshal(body, &unm); err != nil {
		log.Print(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	data, status := s.user2fa(req.Context(), unm.Code, tokenStr)

	logger.Log(logger.InfoPrefix, fmt.Sprintf("Response: %v %s", status, data))
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, err = w.Write(data)
	if err != nil {
		return
	}
}

func (s *authServer) user2fa(ctx context.Context, code string, tokenStr string) ([]byte, int) {
	valid, err := s.m.Validate2FACode(ctx, code, tokenStr)
	if err != nil {
		if errors.Is(err, models.ErrInvalidToken) {
			return nil, http.StatusUnauthorized
		}
		return nil, http.StatusInternalServerError
	}
	if !valid {
		return nil, http.StatusUnauthorized
	}
	return nil, http.StatusOK
}
