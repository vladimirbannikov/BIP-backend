package models

import (
	"context"
	"github.com/pkg/errors"
	"github.com/vladimirbannikov/BIP-backend/internal/models/structs"
)

type TestStorager interface {
	GetTests(ctx context.Context, limit int, offset int) ([]structs.TestSimple, error)
	GetFullTestByID(ctx context.Context, id int) (structs.TestFull, error)
	SaveScore(ctx context.Context, score structs.UserScore) error
	GetTotalRating(ctx context.Context, category string, limit int, offset int) (structs.Rating, error)
}

const defaultLimit = 100
const defaultOffset = 0

const defaultRatingLimit = 200
const defaultRatingOffset = 0

const TotalCategory = "all"

func (m *ModelTests) GetTests(ctx context.Context, limit int, offset int) ([]structs.TestSimple, error) {
	if limit == 0 {
		limit = defaultLimit
	}
	if offset == 0 {
		offset = defaultOffset
	}

	tests, err := m.ts.GetTests(ctx, limit, offset)
	if err != nil {
		return nil, err
	}

	return tests, nil
}

func (m *ModelTests) GetFullTestByID(ctx context.Context, id int) (structs.TestFull, error) {
	test, err := m.ts.GetFullTestByID(ctx, id)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return structs.TestFull{}, ErrNotFound
		}
		return structs.TestFull{}, err
	}
	return test, nil
}

func (m *ModelTests) CheckUserAnswers(ctx context.Context, tokenStr string, testID int, userAnswers []structs.UserAnswer) (structs.TestResult, error) {
	userLogin, err := getUserLoginFromToken(tokenStr)
	if err != nil {
		return structs.TestResult{}, err
	}

	user, err := m.us.GetUserByLogin(ctx, userLogin)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return structs.TestResult{}, ErrNotFound
		}
		return structs.TestResult{}, err
	}

	test, err := m.ts.GetFullTestByID(ctx, testID)
	if err != nil {
		if errors.Is(err, ErrNotFound) {
			return structs.TestResult{}, ErrNotFound
		}
		return structs.TestResult{}, err
	}

	correctAnswers := make(map[int]int)
	for _, question := range test.Questions {
		for _, answer := range question.Answers {
			if answer.IsCorrect {
				correctAnswers[question.ID] = answer.ID
			}
		}
	}

	score := 0
	for _, answer := range userAnswers {
		if correct, ok := correctAnswers[answer.QuestionID]; ok {
			if answer.AnswerID == correct {
				score += 1 * test.DiffLevel
			}
		} else {
			return structs.TestResult{}, ErrInvalidQuestionOrAnswer
		}
	}

	// если пользователь уже проходил тест, скор не запишется в рейтинг
	err = m.ts.SaveScore(ctx, structs.UserScore{
		TestID:    testID,
		UserLogin: userLogin,
		Score:     score,
	})
	if err != nil {
		if !errors.Is(err, ErrConflict) {
			return structs.TestResult{}, err
		}
	}

	return structs.TestResult{
		User:      user.Login,
		TestID:    test.ID,
		Total:     len(correctAnswers),
		UserScore: score,
	}, nil
}

func (m *ModelTests) GetRating(ctx context.Context, category string, limit int, offset int) (structs.Rating, error) {
	if category == "" {
		category = TotalCategory
	}
	if limit == 0 {
		limit = defaultRatingLimit
	}
	if offset == 0 {
		offset = defaultRatingOffset
	}

	rating, err := m.ts.GetTotalRating(ctx, category, limit, offset)
	if err != nil {
		return structs.Rating{}, err
	}

	return rating, nil
}
