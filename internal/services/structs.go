package services

type registerUserReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginUserReq struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

type loginUserResp struct {
	AccessToken string `json:"access_token"`
}
