package config

import (
	"encoding/json"
	"io/ioutil"
)

type Config struct {
	Address    string          `json:"address"`     // 服务地址
	Email      string          `json:"email"`       // 登录邮箱
	Passwd     string          `json:"passwd"`      // 登录密码
	RootDir    string          `json:"dir"`         // 本地笔记根目录
	Updated    int64           `json:"updated"`     // 上次更新时间
	Token      string          `json:"token"`       // 不用配置，自动获取
	IgnoreDirs []string        `json:"ignore_dirs"` // 忽略本地目录
	IgnoreMap  map[string]bool `json:"-"`
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

	cnf.IgnoreMap = make(map[string]bool)
	for _, dir := range cnf.IgnoreDirs {
		cnf.IgnoreMap[dir] = true
	}
}

func GetConfig() *Config {
	return cnf
}
