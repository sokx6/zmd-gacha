package types

type UserRstReq struct {
	Username string
	Nickname string
	Profile  string
	Password string
	Email    string
}

type UserLoginReq struct {
	Username string
	UID      uint
	Email    string
	Password string
}
