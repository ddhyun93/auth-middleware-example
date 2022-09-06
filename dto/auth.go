package dto

import "net/mail"

type CreateUserReq struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type LoginReq struct {
	CreateUserReq
}

type ActiveUserReq struct {
	Code string `json:"code"`
}

type TokenRes struct {
	Token        string `json:"token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshReq struct {
	RefreshToken string `json:"refresh_token"`
}

func (req *CreateUserReq) Validate() (bool, string) {
	var errMsg = ""
	var ok = true
	if req.Email == "" {
		ok = false
		errMsg += "email field is required\n"
	}
	_, err := mail.ParseAddress(req.Email)
	if err != nil {
		ok = false
		errMsg += "entered value is not valid email\n"
	}
	if req.Password == "" {
		ok = false
		errMsg += "password field is required\n"
	}
	return ok, errMsg
}
