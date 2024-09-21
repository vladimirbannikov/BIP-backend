package models

type ModelAuth struct {
	as AuthStorager
	us UsersStorager
}

type ModelUsers struct {
	us UsersStorager
}

type ModelTests struct {
	us UsersStorager
	ts TestStorager
}

func NewModelAuth(as AuthStorager, us UsersStorager) ModelAuth {
	return ModelAuth{as, us}
}
func NewModelUsers(us UsersStorager) ModelUsers {
	return ModelUsers{us}
}
func NewModelTests(us UsersStorager, ts TestStorager) ModelTests {
	return ModelTests{us, ts}
}

const ServiceName = "BipServiceWithTests"
