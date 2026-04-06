package types

type UserRstReq struct {
	Username string
	Password string
	Email    string
}

type UserLoginReq struct {
	Username string
	Password string
}
