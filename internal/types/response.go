package types

type UserRstRsp struct {
	UID     uint   `json:"uid"`
	Message string `json:"message"`
}

type UserLoginRsp struct {
	Message      string `json:"message"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}

type TokenRefRsp struct {
	Message      string `json:"message"`
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
}
