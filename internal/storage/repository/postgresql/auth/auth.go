package auth

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

func (r *Repo) CreateUserSecret(ctx context.Context, secretDTO *structs.UserSecretDTO) error {
	id := 0
	err := r.db.ExecQueryRow(ctx,
		`INSERT INTO auth_schema.users_secrets(login, secret, session_id)
				VALUES($1,$2, $3) returning 1;`, secretDTO.Login, secretDTO.Secret, secretDTO.SessionID).Scan(&id)
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

// DeleteUserSecret delete user secret
func (r *Repo) DeleteUserSecret(ctx context.Context, input structs.DeleteUserSecretInput) error {
	tag, err := r.db.Exec(ctx,
		`DELETE FROM auth_schema.users_secrets 
       			WHERE login = $1 and session_id = $2;`, input.Login, input.SessionID)
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

// GetSecretByUserID get secret
func (r *Repo) GetSecretByUserID(ctx context.Context, input structs.GetUserSecretInput) (*structs.UserSecretDTO, error) {
	secretDTO := structs.UserSecretDTO{}
	err := r.db.Get(ctx, &secretDTO,
		`SELECT login, secret, session_id FROM auth_schema.users_secrets 
				WHERE login=$1 and session_id=$2;`, input.Login, input.SessionID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, err
	}
	return &secretDTO, nil
}

// GetUser2FAInfo get 2fa info
func (r *Repo) GetUser2FAInfo(ctx context.Context, login string) (*structs.User2FAInfoDTO, error) {
	infoDTO := structs.User2FAInfoDTO{}
	err := r.db.Get(ctx, &infoDTO,
		`SELECT login, valid_until, secret FROM auth_schema.users_2fa 
				WHERE login=$1;`, login)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) || errors.Is(err, pgx.ErrNoRows) {
			return nil, repository.ErrObjectNotFound
		}
		return nil, err
	}
	return &infoDTO, nil
}

func (r *Repo) CreateUser2FASecret(ctx context.Context, secretDTO *structs.User2FAInfoDTO) error {
	id := 0
	err := r.db.ExecQueryRow(ctx,
		`INSERT INTO auth_schema.users_2fa(login, valid_until, secret)
				VALUES($1,$2, $3) returning 1;`, secretDTO.Login, secretDTO.ValidUntil, secretDTO.Secret).Scan(&id)
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

func (r *Repo) UpdateUser2FASecret(ctx context.Context, secretDTO *structs.User2FAInfoDTO) error {
	id := 0
	err := r.db.ExecQueryRow(ctx,
		`UPDATE auth_schema.users_2fa set 
				valid_until = $1 returning 1;`, secretDTO.ValidUntil).Scan(&id)
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
