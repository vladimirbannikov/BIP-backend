package tests

import (
	"context"
	"database/sql"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/pkg/errors"
	"github.com/vladimirbannikov/BIP-backend/internal/models"
	"github.com/vladimirbannikov/BIP-backend/internal/models/structs"
	"github.com/vladimirbannikov/BIP-backend/internal/storage/db"
	"github.com/vladimirbannikov/BIP-backend/internal/storage/repository"
)

type Repo struct {
	db db.DBops
}

func New(db db.DBops) *Repo {
	return &Repo{db: db}
}

func (r *Repo) GetTests(ctx context.Context, limit int, offset int) ([]structs.TestSimple, error) {
	var tests []*structs.TestSimple
	err := r.db.Select(ctx, &tests,
		`SELECT id, name, description, category, diff_level, pictureFile  
		FROM test 
		ORDER by id 
		limit $1 offset $2;
		`, limit, offset)
	if err != nil {
		return nil, err
	}
	testsOut := make([]structs.TestSimple, len(tests))
	for i, test := range tests {
		testsOut[i] = *test
	}
	return testsOut, nil
}

func (r *Repo) GetFullTestByID(ctx context.Context, id int) (structs.TestFull, error) {
	var test structs.TestSimple
	err := r.db.Get(ctx, &test,
		`Select id, name, description, category, diff_level FROM test WHERE id = $1;`, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return structs.TestFull{}, repository.ErrObjectNotFound
		}
		return structs.TestFull{}, err
	}

	var questionsDTOsp []*structs.TestQuestionDTO
	err = r.db.Select(ctx, &questionsDTOsp,
		`Select id, test_id, question, is_song, song_file FROM test_questions WHERE test_id = $1;`, id)
	if err != nil {
		return structs.TestFull{}, err
	}

	questionsDTOs := make([]structs.TestQuestionDTO, len(questionsDTOsp))
	for i, q := range questionsDTOsp {
		questionsDTOs[i] = *q
	}

	questions := make([]structs.TestQuestionFull, 0, len(questionsDTOs))
	for _, question := range questionsDTOs {
		var answersp []*structs.QuestionAnswer
		err = r.db.Select(ctx, &answersp,
			`select id, question_id, answer, is_correct 
					from question_variants where question_id = $1`, question.ID)
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return structs.TestFull{}, repository.ErrObjectNotFound
		}

		answers := make([]structs.QuestionAnswer, len(answersp))
		for i, a := range answersp {
			answers[i] = *a
		}

		questions = append(questions, structs.TestQuestionFull{
			ID:       question.ID,
			TestID:   question.TestID,
			Question: question.Question,
			IsSong:   question.IsSong,
			Song:     []byte(question.SongFile),
			Answers:  answers,
		})
	}

	return structs.TestFull{
		ID:          test.ID,
		Name:        test.Name,
		Description: test.Description,
		Category:    test.Category,
		DiffLevel:   test.DiffLevel,
		Questions:   questions,
	}, nil
}

func (r *Repo) SaveScore(ctx context.Context, dto structs.UserScoreDTO) error {
	id := 0
	err := r.db.ExecQueryRow(ctx,
		`INSERT INTO user_test_score(test_id, user_login, score)
				VALUES($1,$2,$3) returning 1;`, dto.TestID, dto.UserLogin, dto.Score).Scan(&id)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)
		if pgErr.Code == "23505" {
			return repository.ErrDuplicateKey
		}
		return err
	}

	err = r.db.ExecQueryRow(ctx,
		`UPDATE users_schema.user_profile 
			SET total_score = total_score + $1,
			    tests_count = tests_count + 1 
			    WHERE login = $2 returning 1;`,
		dto.Score, dto.UserLogin).Scan(&id)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	return nil
}

func (r *Repo) GetTotalScores(ctx context.Context, category string, limit int, offset int) (structs.Rating, error) {
	var unitsp []*structs.RatingUnitDTO
	if category == models.TotalCategory {
		err := r.db.Select(ctx, &unitsp,
			`SELECT user_login, SUM(score) as sum
			FROM user_test_score
			GROUP BY user_login
			order by SUM(score) DESC
			limit $1 OFFSET $2;`, limit, offset)
		if err != nil {
			return structs.Rating{}, err
		}
	} else {
		err := r.db.Select(ctx, &unitsp,
			`SELECT user_login, SUM(score) as sum
			FROM user_test_score
			where category = $1
			GROUP BY user_login
			order by SUM(score) DESC 
			limit $2 OFFSET $3;`, category, limit, offset)
		if err != nil {
			return structs.Rating{}, err
		}
	}

	units := make([]structs.RatingUnitDTO, len(unitsp))
	for i, u := range unitsp {
		units[i] = *u
	}

	res := structs.Rating{}
	for place, unit := range units {
		res.Rating = append(res.Rating, structs.RatingUnit{
			Login: unit.Login,
			Place: place + 1,
			Score: unit.Sum,
		})
	}
	return res, nil
}
