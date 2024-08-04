package structs

type User struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// UserDTO dto
type UserDTO struct {
	Login        string `db:"login"`
	Email        string `db:"email"`
	PasswordHash string `db:"password_hash"`
}

type UserProfile struct {
	Login string `json:"login"`
	Info  string `json:"info"`
}

type UserProfileDTO struct {
	Login string `json:"login"`
	Info  string `json:"info"`
}
