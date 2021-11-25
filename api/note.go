package api

import (
	"encoding/json"
	"leanote-sync/util"
	"strconv"
	"time"
)

// 新增笔记，默认 markdown
type AddNoteResp struct {
	NoteId      string `json:"NoteId"`
	NotebookId  string `json:"NoteBookId"`
	UserId      string `json:"UserId"`
	Title       string `json:"Title"`
	Desc        string `json:"Desc"`
	Abstract    string `json:"Abstract"`
	Content     string `json:"Content"`
	IsMarkdown  bool   `json:"IsMarkdown"`
	IsBlog      bool   `json:"IsBlog"`
	IsTrash     bool   `json:"IsTrash"`
	IsDeleted   bool   `json:"IsDeleted"`
	Usn         int    `json:"Usn"`
	CreatedTime string `json:"CreatedTime"`
	UpdatedTime string `json:"UpdatedTime"`
	PublicTime  string `json:"PublicTime"`
}

func AddNote(address, token, notebookId, title, content, abstract string) (*AddNoteResp, error) {
	var (
		ret  = new(AddNoteResp)
		resp = []byte{}
		err  error
	)
	param := map[string]string{
		"token":      token,
		"NotebookId": notebookId,
		"Title":      title,
		"Content":    content,
		"Abstract":   abstract,
		"IsMarkdown": "true",
	}

	resp, err = util.HttpRequest(address+"/api/note/addNote", param)
	if err != nil {
		return ret, err
	}
	err = json.Unmarshal(resp, ret)
	return ret, err
}

// 获得笔记内容
type GetNoteContentResp struct {
	NoteID  string `json:"NoteId"`
	UserID  string `json:"UserId"`
	Content string `json:"Content"`
}

func GetNoteContent(address, token, noteId string) (GetNoteContentResp, error) {
	var (
		ret  = GetNoteContentResp{}
		resp = []byte{}
		err  error
	)
	param := map[string]string{
		"token":  token,
		"noteId": noteId,
		"format": "json",
	}
	resp, err = util.HttpRequest(address+"/api/note/getNoteContent", param)
	if err != nil {
		return ret, err
	}
	err = json.Unmarshal(resp, &ret)
	return ret, err
}

// 获得笔记与内容
type GetNoteAndContentResp struct {
	NoteID      string        `json:"NoteId"`
	NotebookID  string        `json:"NotebookId"`
	UserID      string        `json:"UserId"`
	Title       string        `json:"Title"`
	Desc        string        `json:"Desc"`
	Tags        interface{}   `json:"Tags"`
	Abstract    string        `json:"Abstract"`
	Content     string        `json:"Content"`
	IsMarkdown  bool          `json:"IsMarkdown"`
	IsBlog      bool          `json:"IsBlog"`
	IsTrash     bool          `json:"IsTrash"`
	IsDeleted   bool          `json:"IsDeleted"`
	Usn         int           `json:"Usn"`
	Files       []interface{} `json:"Files"`
	CreatedTime time.Time     `json:"CreatedTime"`
	UpdatedTime time.Time     `json:"UpdatedTime"`
	PublicTime  time.Time     `json:"PublicTime"`
}

func GetNoteAndContent(address, token, noteId string) (GetNoteAndContentResp, error) {
	var (
		ret  = GetNoteAndContentResp{}
		resp = []byte{}
		err  error
	)
	param := map[string]string{
		"token":  token,
		"noteId": noteId,
		"format": "json",
	}
	resp, err = util.HttpRequest(address+"/api/note/getNoteAndContent", param)
	if err != nil {
		return ret, err
	}
	err = json.Unmarshal(resp, &ret)
	return ret, err
}

// 获得笔记本下的笔记
type GetNotesResp struct {
	NoteID      string        `json:"NoteId"`
	NotebookID  string        `json:"NotebookId"`
	UserID      string        `json:"UserId"`
	Title       string        `json:"Title"`
	Desc        string        `json:"Desc"`
	Tags        interface{}   `json:"Tags"`
	Abstract    string        `json:"Abstract"`
	Content     string        `json:"Content"` // 这个Content一直为空，有点坑
	IsMarkdown  bool          `json:"IsMarkdown"`
	IsBlog      bool          `json:"IsBlog"`
	IsTrash     bool          `json:"IsTrash"`
	IsDeleted   bool          `json:"IsDeleted"`
	Usn         int           `json:"Usn"`
	Files       []interface{} `json:"Files"`
	CreatedTime time.Time     `json:"CreatedTime"`
	UpdatedTime time.Time     `json:"UpdatedTime"`
	PublicTime  time.Time     `json:"PublicTime"`
}

func GetNotes(address, token, notebookId string) ([]*GetNotesResp, error) {
	ret := make([]*GetNotesResp, 0)
	param := map[string]string{
		"token":      token,
		"notebookId": notebookId,
	}
	resp, err := util.HttpRequest(address+"/api/note/getNotes", param)
	if err != nil {
		return ret, err
	}
	json.Unmarshal(resp, &ret)
	return ret, err

}

// 更新笔记
type UpdateNoteResp struct {
	NoteID      string      `json:"NoteId"`
	NotebookID  string      `json:"NotebookId"`
	UserID      string      `json:"UserId"`
	Title       string      `json:"Title"`
	Desc        string      `json:"Desc"`
	Tags        interface{} `json:"Tags"`
	Abstract    string      `json:"Abstract"`
	Content     string      `json:"Content"`
	IsMarkdown  bool        `json:"IsMarkdown"`
	IsBlog      bool        `json:"IsBlog"`
	IsTrash     bool        `json:"IsTrash"`
	IsDeleted   bool        `json:"IsDeleted"`
	Usn         int         `json:"Usn"`
	Files       interface{} `json:"Files"`
	CreatedTime time.Time   `json:"CreatedTime"`
	UpdatedTime time.Time   `json:"UpdatedTime"`
	PublicTime  time.Time   `json:"PublicTime"`
}

func UpdateNote(address, token, noteId, title, content, abstract string, usn int) (UpdateNoteResp, error) {
	var (
		ret  = UpdateNoteResp{}
		resp = []byte{}
		err  error
	)
	param := map[string]string{
		"token":      token,
		"NoteId":     noteId,
		"Title":      title,
		"Content":    content,
		"Abstract":   abstract,
		"IsMarkdown": "true",
		"usn":        strconv.Itoa(usn),
	}

	resp, err = util.HttpRequest(address+"/api/note/updateNote", param)
	if err != nil {
		return ret, err
	}
	err = json.Unmarshal(resp, &ret)
	return ret, err
}

type DeleteNoteResp struct {
	Ok  bool   `json:"Ok"`
	Msg string `json:"Msg"`
	Usn int    `json:"Usn"`
}

func DeleteNote(address, token, noteId string, usn int) (DeleteNoteResp, error) {
	var (
		ret  = DeleteNoteResp{}
		resp = []byte{}
		err  error
	)
	param := map[string]string{
		"token":  token,
		"noteId": noteId,
		"usn":    strconv.Itoa(usn),
	}

	resp, err = util.HttpRequest(address+"/api/note/deleteTrash", param)
	if err != nil {
		return ret, err
	}
	err = json.Unmarshal(resp, &ret)
	return ret, err
}
