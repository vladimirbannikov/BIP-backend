package structs

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
	Password string
}

type LoginUserInput struct {
	Login    string
	Password string
}
