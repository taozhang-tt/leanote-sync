package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"leanote-sync/api"
)

type Config struct {
	Address string `json:"address"` // 服务地址
	Email   string `json:"email"`   // 登录邮箱
	Passwd  string `json:"passwd"`  // 登录密码
	RootDir string `json:"dir"`     // 本地笔记根目录
	Updated int64  `json:"updated"` // 上次更新时间
	Token   string `json:"token"`   // 不用配置，自动获取
}

var cnf = new(Config)

func init() {
	bs, err := ioutil.ReadFile("/etc/leanote-sync.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(bs, cnf)
	if err != nil {
		panic(err)
	}
	loginResp, err := api.Login(cnf.Address, cnf.Email, cnf.Passwd)
	if err != nil {
		panic(fmt.Sprintf("api.Login(%v, %v, %v) failed: %v\n", cnf.Address, cnf.Email, cnf.Passwd, err))
	}
	if !loginResp.Ok {
		panic(fmt.Sprintf("api.Login(%v, %v, %v) failed: %v\n", cnf.Address, cnf.Email, cnf.Passwd, loginResp.Msg))
	}
	cnf.Token = loginResp.Token
}

func GetConfig() *Config {
	return cnf
}
