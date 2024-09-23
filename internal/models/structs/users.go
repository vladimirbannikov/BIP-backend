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
	Login             string `json:"login"`
	Email             string `json:"email"`
	TotalScore        int    `json:"total_score"`
	TestCount         int    `json:"tests_count"`
	GlobalRatingPlace int    `json:"global_rating"`
	Avatar            []byte `json:"avatar"`
}

type UserProfileDTO struct {
	Login      string `db:"login"`
	Email      string `db:"email"`
	TotalScore int    `db:"total_score"`
	TestCount  int    `db:"tests_count"`
	AvatarFile string `db:"avatarfile"`
}
