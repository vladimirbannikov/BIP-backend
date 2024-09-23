package app

import (
	"context"
	"fmt"
	"github.com/vladimirbannikov/BIP-backend/internal/models"
	"github.com/vladimirbannikov/BIP-backend/internal/services"
	"github.com/vladimirbannikov/BIP-backend/internal/storage"
	"github.com/vladimirbannikov/BIP-backend/internal/storage/file-storage/avatars"
	"github.com/vladimirbannikov/BIP-backend/internal/storage/file-storage/file_provider"
	"github.com/vladimirbannikov/BIP-backend/internal/storage/file-storage/music"
	"github.com/vladimirbannikov/BIP-backend/internal/storage/file-storage/test_pictures"
	"github.com/vladimirbannikov/BIP-backend/internal/storage/repository/postgresql/auth"
	"github.com/vladimirbannikov/BIP-backend/internal/storage/repository/postgresql/tests"
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
	testsRepo := tests.New(dbStor.DB)

	f := file_provider.NewFileProvider()
	ar := avatars.NewRepository(f)
	tp := test_pictures.NewRepository(f)
	mp := music.NewRepository(f)

	usersStorage := storage.NewUsersStorage(usersRepo, ar)
	authStorage := storage.NewAuthStorage(authRepo)
	testsStorage := storage.NewTestsStorage(testsRepo, tp, mp)

	amdl := models.NewModelAuth(&authStorage, &usersStorage)
	umdl := models.NewModelUsers(&usersStorage)
	tmdl := models.NewModelTests(&usersStorage, &testsStorage)

	serv := services.NewService(&amdl, &umdl, &tmdl)

	serv.Launch()

	return nil
}
