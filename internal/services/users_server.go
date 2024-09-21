package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/vladimirbannikov/BIP-backend/internal/models"
	"github.com/vladimirbannikov/BIP-backend/internal/models/structs"
	"github.com/vladimirbannikov/BIP-backend/internal/utils/logger"
	"io"
	"net/http"
)

type UsersModelManager interface {
	GetUserProfile(ctx context.Context, login string) (structs.UserProfile, error)
	GetUserProfileOwn(ctx context.Context, tokenStr string) (structs.UserProfile, error)
	UpdateUser(ctx context.Context, tokenStr string, newUser structs.User) error
}

type usersServer struct {
	m UsersModelManager
}

func (s *usersServer) GetUserProfileOwn(w http.ResponseWriter, req *http.Request) {
	tokenStr := req.Header.Get(authHeaderStr)
	if tokenStr == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	data, status := s.getUserProfileOwn(req.Context(), tokenStr)
	logger.Log(logger.InfoPrefix, fmt.Sprintf("Response: %v %s", status, data))
	w.WriteHeader(status)
	_, err := w.Write(data)
	if err != nil {
		return
	}
}

func (s *usersServer) getUserProfileOwn(ctx context.Context, tokenStr string) ([]byte, int) {
	p, err := s.m.GetUserProfileOwn(ctx, tokenStr)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, http.StatusNotFound
		}
		if errors.Is(err, models.ErrBadCredentials) || errors.Is(err, models.ErrInvalidToken) {
			return nil, http.StatusUnauthorized
		}
		logger.Log(logger.ErrPrefix, fmt.Sprintf(err.Error()))
		return nil, http.StatusInternalServerError
	}

	articleJSON, _ := json.Marshal(GetUserProfileOwnResp{
		Login:             p.Login,
		Email:             p.Email,
		TotalScore:        p.TotalScore,
		TestCount:         p.TestCount,
		GlobalRatingPlace: p.GlobalRatingPlace,
	})
	return articleJSON, http.StatusOK
}

func (s *usersServer) UpdateUserProfile(w http.ResponseWriter, req *http.Request) {
	tokenStr := req.Header.Get(authHeaderStr)
	if tokenStr == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	body, err := io.ReadAll(req.Body)
	if err != nil {
		logger.Log(logger.ErrPrefix, fmt.Sprintf("Service: usersServer: UpdateProfile: ReadAll error: %s", err.Error()))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	var unm UpdateUserProfileReq
	if err = json.Unmarshal(body, &unm); err != nil {
		logger.Log(logger.ErrPrefix, fmt.Sprintf("Service: usersServer: UpdateProfile: json.Unmarshal error: %s", err.Error()))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	userUpdate := structs.User{
		Login:    "",
		Email:    unm.Email,
		Password: unm.Password,
	}
	data, status := s.updateUserProfile(req.Context(), tokenStr, userUpdate)
	logger.Log(logger.InfoPrefix, fmt.Sprintf("Response: %v %s", status))
	w.WriteHeader(status)
	_, err = w.Write(data)
	if err != nil {
		return
	}
}

func (s *usersServer) updateUserProfile(ctx context.Context, tokenStr string, u structs.User) ([]byte, int) {
	err := s.m.UpdateUser(ctx, tokenStr, u)
	if err != nil {
		if errors.Is(err, models.ErrNotFound) {
			return nil, http.StatusNotFound
		}
		if errors.Is(err, models.ErrInvalidInput) || errors.Is(err, models.ErrConflict) {
			return nil, http.StatusBadRequest
		}
		if errors.Is(err, models.ErrInvalidToken) {
			return nil, http.StatusUnauthorized
		}
		return nil, http.StatusInternalServerError
	}
	return []byte(""), http.StatusOK
}

func (s *usersServer) GetUserProfile(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	login := vars[loginKey]

	data, status := s.getUserProfile(req.Context(), login)
	logger.Log(logger.InfoPrefix, fmt.Sprintf("Response: %v %s", status, data))
	w.WriteHeader(status)
	_, err := w.Write(data)
	if err != nil {
		return
	}
}

func (s *usersServer) getUserProfile(ctx context.Context, login string) ([]byte, int) {
	p, err := s.m.GetUserProfile(ctx, login)
	if err != nil {
		if errors.Is(err, models.ErrBadCredentials) || errors.Is(err, models.ErrInvalidToken) {
			return nil, http.StatusUnauthorized
		}
		if errors.Is(err, models.ErrNotFound) {
			return nil, http.StatusNotFound
		}
		return nil, http.StatusInternalServerError
	}

	articleJSON, _ := json.Marshal(GetUserProfileResp{
		Login:             p.Login,
		TotalScore:        p.TotalScore,
		TestCount:         p.TestCount,
		GlobalRatingPlace: p.GlobalRatingPlace,
	})
	return articleJSON, http.StatusOK
}
