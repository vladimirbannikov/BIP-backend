package storage

import (
	"context"
	"errors"
	"github.com/vladimirbannikov/BIP-backend/internal/models"
	"github.com/vladimirbannikov/BIP-backend/internal/models/structs"
	"github.com/vladimirbannikov/BIP-backend/internal/storage/repository"
)

type TestsRepo interface {
	GetTests(ctx context.Context, limit int, offset int) ([]structs.TestSimple, error)
	GetFullTestByID(ctx context.Context, id int) (structs.TestFull, error)
	SaveScore(ctx context.Context, dto structs.UserScoreDTO) error
	GetTotalScores(ctx context.Context, category string, limit int, offset int) (structs.Rating, error)
}

type TestsStorage struct {
	testsRepo TestsRepo
}

func NewTestsStorage(testsRepo TestsRepo) TestsStorage {
	return TestsStorage{testsRepo: testsRepo}
}

func (t *TestsStorage) GetTests(ctx context.Context, limit int, offset int) ([]structs.TestSimple, error) {
	return t.testsRepo.GetTests(ctx, limit, offset)
}

func (t *TestsStorage) GetFullTestByID(ctx context.Context, id int) (structs.TestFull, error) {
	return t.testsRepo.GetFullTestByID(ctx, id)
}

func (t *TestsStorage) SaveScore(ctx context.Context, score structs.UserScore) error {
	err := t.testsRepo.SaveScore(ctx, structs.UserScoreDTO{
		TestID:    score.TestID,
		UserLogin: score.UserLogin,
		Score:     score.Score,
	})
	if err != nil {
		if errors.Is(err, repository.ErrDuplicateKey) {
			return models.ErrConflict
		}
		return err
	}
	return nil
}

func (t *TestsStorage) GetTotalRating(ctx context.Context, category string, limit int, offset int) (structs.Rating, error) {
	rating, err := t.testsRepo.GetTotalScores(ctx, category, limit, offset)
	if err != nil {
		return structs.Rating{}, err
	}
	return rating, nil
}