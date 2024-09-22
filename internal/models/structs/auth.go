package structs

import "time"

type UserSecret struct {
	Login     string
	Secret    string
	SessionID string
}

type UserSecretDTO struct {
	Login     string `db:"login"`
	Secret    string `db:"secret"`
	SessionID string `db:"session_id"`
}

type User2FAInfo struct {
	Login      string
	ValidUntil time.Time
	Secret     string
}

type User2FAInfoDTO struct {
	Login      string    `db:"login"`
	ValidUntil time.Time `db:"valid_until"`
	Secret     string    `db:"secret"`
}

type GetUserSecretInput struct {
	Login     string
	SessionID string
}

type DeleteUserSecretInput struct {
	Login     string
	SessionID string
}

type RegisterUserInput struct {
	Login    string
	Email    string
	Password string
}

type LoginUserInput struct {
	Login    string
	Password string
}
