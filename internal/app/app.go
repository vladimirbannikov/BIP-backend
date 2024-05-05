package app

import (
	"context"
	"fmt"
	"github.com/vladimirbannikov/BIP-backend/internal/models"
	"github.com/vladimirbannikov/BIP-backend/internal/services"
	"github.com/vladimirbannikov/BIP-backend/internal/storage"
	"github.com/vladimirbannikov/BIP-backend/internal/storage/repository/postgresql/auth"
	"github.com/vladimirbannikov/BIP-backend/internal/storage/repository/postgresql/users"
	"github.com/vladimirbannikov/BIP-backend/internal/utils/logger"
)

func Start() error {
	ctx := context.Background()
	dbStor, err := storage.NewDbStorage(ctx)
	if err != nil {
		logger.Log(logger.ErrPrefix, fmt.Sprintf("App: Start: NewDbStorage: %s", err.Error()))
		return err
	}
	defer dbStor.Close(ctx)

	usersRepo := users.New(dbStor.DB)
	authRepo := auth.New(dbStor.DB)

	usersStorage := storage.NewUsersStorage(usersRepo)
	authStorage := storage.NewAuthStorage(authRepo)

	amdl := models.NewModelAuth(&authStorage, &usersStorage)
	umdl := models.NewModelUsers(&usersStorage)

	serv := services.NewService(&amdl, &umdl)

	serv.Launch()

	return nil
}
