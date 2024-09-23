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
	tp        TestPicProvider
}

type TestPicProvider interface {
	GetPictureByte(pic string) ([]byte, error)
}

func NewTestsStorage(testsRepo TestsRepo, tp TestPicProvider) TestsStorage {
	return TestsStorage{testsRepo: testsRepo, tp: tp}
}

func (t *TestsStorage) GetTests(ctx context.Context, limit int, offset int) ([]structs.TestSimple, error) {
	tests, err := t.testsRepo.GetTests(ctx, limit, offset)
	if err != nil {
		return nil, err
	}
	testsOut := make([]structs.TestSimple, 0)
	for _, test := range tests {
		picByte, err := t.tp.GetPictureByte(string(test.Picture))
		if err != nil {
			return nil, err
		}
		testsOut = append(testsOut, structs.TestSimple{
			ID:          test.ID,
			Name:        test.Name,
			Description: test.Description,
			Category:    test.Category,
			DiffLevel:   test.DiffLevel,
			Picture:     picByte,
		})
	}
	return testsOut, nil
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
