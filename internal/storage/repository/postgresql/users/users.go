package users

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
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

// CreateUser create user
func (r *Repo) CreateUser(ctx context.Context, user *structs.UserDTO) (int, error) {
	id := 0
	err := r.db.ExecQueryRow(ctx,
		`INSERT INTO users_schema.users(login, email, password_hash)
				VALUES($1,$2,$3) returning 1;`, user.Login, user.Email, user.PasswordHash).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)
		if pgErr.Code == "23505" {
			return 0, repository.ErrDuplicateKey
		}
		return 0, err
	}
	return id, nil
}

func (r *Repo) CreateUserProfile(ctx context.Context, user *structs.UserProfileDTO) (int, error) {
	id := 0
	err := r.db.ExecQueryRow(ctx,
		`INSERT INTO users_schema.user_profile(login, tests_count, total_score, avatarFile)
				VALUES($1,$2,$3,$4) returning 1;`, user.Login, 0, 0, user.AvatarFile).Scan(&id)

	if err != nil {
		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)
		if pgErr.Code == "23505" {
			return 0, repository.ErrDuplicateKey
		}
		return 0, err
	}

	return id, nil
}

// GetUserByLogin get user
func (r *Repo) GetUserByLogin(ctx context.Context, login string) (*structs.UserDTO, error) {
	var info structs.UserDTO
	err := r.db.Get(ctx, &info,
		`SELECT login, password_hash, email FROM users_schema.users WHERE login=$1;`, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, err
	}
	return &info, nil
}

// DeleteUserByLogin delete user
func (r *Repo) DeleteUserByLogin(ctx context.Context, login string) error {
	tag, err := r.db.Exec(ctx,
		`DELETE FROM users_schema.users WHERE login = $1;`, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return repository.ErrObjectNotFound
		}
		return err
	}
	if tag.RowsAffected() == 0 {
		return repository.ErrObjectNotFound
	}
	return nil
}

func (r *Repo) GetUserProfileByLogin(ctx context.Context, login string) (*structs.UserProfile, error) {
	var info structs.UserProfileDTO
	err := r.db.Get(ctx, &info,
		`SELECT u.login as login, u.email as email, up.total_score as total_score, 
       up.tests_count as tests_count, up.avatarFile as avatarfile  FROM users_schema.users u
                left join users_schema.user_profile up on u.login = up.login 
                                                WHERE u.login=$1;`, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, err
	}

	rating, err := r.getTotalScores(ctx)
	if err != nil {
		return nil, err
	}

	place := 0
	for _, unit := range rating.Rating {
		if login == unit.Login {
			place = unit.Place
			break
		}
	}

	return &structs.UserProfile{
		Login:             login,
		Email:             info.Email,
		TotalScore:        info.TotalScore,
		TestCount:         info.TestCount,
		GlobalRatingPlace: place,
		Avatar:            []byte(info.AvatarFile),
	}, nil
}

// todo: лежит в tests
func (r *Repo) getTotalScores(ctx context.Context) (structs.Rating, error) {
	var unitsp []*structs.RatingUnitDTO
	err := r.db.Select(ctx, &unitsp,
		`SELECT user_login, SUM(score) as sum
			FROM user_test_score
			GROUP BY user_login
			order by SUM(score) DESC;`)
	if err != nil {
		return structs.Rating{}, err
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

func (r *Repo) UpdateUser(ctx context.Context, user *structs.UserDTO) error {
	id := 0
	err := r.db.ExecQueryRow(ctx,
		`UPDATE users_schema.users set 
				email = $1, password_hash = $2 WHERE login=$1 returning 1;`, user.Email, user.PasswordHash).Scan(&id)
	if err != nil {
		var pgErr *pgconn.PgError
		errors.As(err, &pgErr)
		if pgErr.Code == "23505" {
			return repository.ErrDuplicateKey
		}
		return err
	}
	return nil
}
