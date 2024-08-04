package models

type ModelAuth struct {
	as AuthStorager
	us UsersStorager
}

type ModelUsers struct {
	us UsersStorager
}

func NewModelAuth(as AuthStorager, us UsersStorager) ModelAuth {
	return ModelAuth{as, us}
}
func NewModelUsers(us UsersStorager) ModelUsers {
	return ModelUsers{us}
}

const ServiceName = "BipServiceWithTests"
