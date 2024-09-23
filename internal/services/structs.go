package services

type registerUserReq struct {
	Login    string `json:"login"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type loginUserReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginUserResp struct {
	AccessToken string `json:"access_token"`
}

type User2FAInput struct {
	Code string `json:"code"`
}

type GetUserProfileOwnResp struct {
	Login             string `json:"login"`
	Email             string `json:"email"`
	TotalScore        int    `json:"total_score"`
	TestCount         int    `json:"tests_count"`
	GlobalRatingPlace int    `json:"global_rating"`
	Avatar            string `json:"avatar"`
}

type GetUserProfileResp struct {
	Login             string `json:"login"`
	TotalScore        int    `json:"total_score"`
	TestCount         int    `json:"tests_count"`
	GlobalRatingPlace int    `json:"global_rating"`
	Avatar            string `json:"avatar"`
}

type UpdateUserProfileReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type GetAllTestsResp struct {
	Tests []GetAllTestsRespTest `json:"tests"`
}

type GetAllTestsRespTest struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	DiffLevel   int    `json:"diff_level"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Picture     string `json:"picture"`
}

type GetFullTestResp struct {
	Id          int                       `json:"id"`
	Name        string                    `json:"name"`
	DiffLevel   int                       `json:"diff_level"`
	Description string                    `json:"description"`
	Category    string                    `json:"category"`
	Questions   []GetFullTestRespQuestion `json:"questions"`
}

type GetFullTestRespQuestion struct {
	Id       int                     `json:"id"`
	Question string                  `json:"question"`
	IsSong   bool                    `json:"isSong"`
	Song     string                  `json:"song"`
	Answers  []GetFullTestRespAnswer `json:"answers"`
}

type GetFullTestRespAnswer struct {
	Id     int    `json:"id"`
	Answer string `json:"answer"`
}

type GetUserScoreReq struct {
	Answers []GetUserScoreReqAnswer `json:"user_answers"`
}

type GetUserScoreReqAnswer struct {
	QuestionId int `json:"question_id"`
	AnswerId   int `json:"answer_id"`
}

type GetUserScoreResp struct {
	UserScore int `json:"user_score"`
	Total     int `json:"total"`
}

type GetRatingResp struct {
	Rating []GetRatingRespUnit `json:"rating"`
}

type GetRatingRespUnit struct {
	Login string `json:"login"`
	Place int    `json:"place"`
	Score int    `json:"score"`
}

type QrCodeResp struct {
	Qr string `json:"qr"`
}
