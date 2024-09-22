package structs

type Rating struct {
	Rating   []RatingUnit
	Category int
}

type RatingUnit struct {
	Login string `json:"login"`
	Place int    `json:"place"`
	Score int    `json:"score"`
}

type RatingUnitDTO struct {
	Login string `db:"user_login"`
	Sum   int    `db:"sum"`
}
