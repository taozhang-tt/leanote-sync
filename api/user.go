package api

import (
	"encoding/json"
	"leanote-sync/util"
)

type LoginResp struct {
	Ok       bool   `json:"Ok"`
	Token    string `json:"Token"`
	UserID   string `json:"UserId"`
	Email    string `json:"Email"`
	Username string `json:"Username"`
	Msg      string `json:"Msg"`
}

// 登录
func Login(address, email, pwd string) (*LoginResp, error) {
	var (
		ret  = new(LoginResp)
		resp = []byte{}
		err  error
	)
	param := map[string]string{
		"email": email,
		"pwd":   pwd,
	}
	resp, err = util.HttpRequest(address+"/api/auth/login", param)
	if err != nil {
		return ret, err
	}
	err = json.Unmarshal(resp, ret)
	return ret, err
}
