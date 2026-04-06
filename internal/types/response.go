package types

type UserRstRsp struct {
	UID     uint   `json:"uid"`
	Message string `json:"message"`
}

type UserLoginRsp struct {
	Message string `json:"message"`
}
