package api

import (
	"encoding/json"
	"strconv"
	"time"

	"leanote-sync/util"
)

// 获取所有笔记本
type GetNotebooksResp struct {
	NotebookID       string    `json:"NotebookId"`
	UserID           string    `json:"UserId"`
	ParentNotebookID string    `json:"ParentNotebookId"`
	Seq              int       `json:"Seq"`
	Title            string    `json:"Title"`
	URLTitle         string    `json:"UrlTitle"`
	IsBlog           bool      `json:"IsBlog"`
	CreatedTime      time.Time `json:"CreatedTime"`
	UpdatedTime      time.Time `json:"UpdatedTime"`
	Usn              int       `json:"Usn"`
	IsDeleted        bool      `json:"IsDeleted"`
}

func GetNotebooks(address, token string) ([]*GetNotebooksResp, error) {
	ret := make([]*GetNotebooksResp, 0)
	param := map[string]string{
		"token": token,
	}
	resp, err := util.HttpRequest(address+"/api/notebook/getNotebooks", param)
	if err != nil {
		return ret, err
	}
	err = json.Unmarshal(resp, &ret)
	return ret, err
}

// 添加笔记本
type AddNotebookResp struct {
	NotebookID       string    `json:"NotebookId"`
	UserID           string    `json:"UserId"`
	ParentNotebookID string    `json:"ParentNotebookId"`
	Seq              int       `json:"Seq"`
	Title            string    `json:"Title"`
	URLTitle         string    `json:"UrlTitle"`
	IsBlog           bool      `json:"IsBlog"`
	CreatedTime      time.Time `json:"CreatedTime"`
	UpdatedTime      time.Time `json:"UpdatedTime"`
	Usn              int       `json:"Usn"`
	IsDeleted        bool      `json:"IsDeleted"`
}

func AddNotebook(address, token, parentNotebookId, title string, seq int) (*AddNotebookResp, error) {
	var (
		resp = []byte{}
		ret  = new(AddNotebookResp)
		err  error
	)
	param := map[string]string{
		"token":            token,
		"title":            title,
		"parentNotebookId": parentNotebookId,
		"seq":              strconv.Itoa(seq),
	}
	resp, err = util.HttpRequest(address+"/api/notebook/addNotebook", param)
	if err != nil {
		return ret, err
	}
	err = json.Unmarshal(resp, &ret)
	return ret, err
}

// 删除笔记本
type DeleteNotebookResp struct {
	Ok  bool   `json:"Ok"`
	Msg string `json:"Msg"`
}

func DeleteNotebook(address, token, notebookId string, usn int) (bool, error) {
	var (
		resp = []byte{}
		ret  = new(DeleteNotebookResp)
		err  error
	)
	param := map[string]string{
		"token":      token,
		"notebookId": notebookId,
		"usn":        strconv.Itoa(usn),
	}
	resp, err = util.HttpRequest(address+"/api/notebook/deleteNotebook", param)
	if err != nil {
		return false, err
	}
	err = json.Unmarshal(resp, &ret)
	return ret.Ok, err
}
