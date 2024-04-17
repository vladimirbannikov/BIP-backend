package structs

type User struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// UserDTO dto
type UserDTO struct {
	Login        string `db:"login"`
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
