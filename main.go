package main

import (
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"leanote-sync/api"
	. "leanote-sync/config"
)

// 这些账本总能获取到，即使已经被删除了
var ConstNotebook = map[string]bool{
	"life":  true,
	"work":  true,
	"study": true,
}

var config = GetConfig()

func main() {
	now := time.Now().Unix()

	// 登录操作
	loginResp, err := api.Login(config.Address, config.Email, config.Passwd)
	if err != nil {
		fmt.Printf("api.Login(%v, %v, %v) failed: %v\n", config.Address, config.Email, config.Passwd, err)
		return
	}
	if !loginResp.Ok {
		fmt.Printf("api.Login(%v, %v, %v) failed: %v\n", config.Address, config.Email, config.Passwd, loginResp.Msg)
		return
	}
	config.Token = loginResp.Token

	// 读取笔记目录，构造本地目录树
	list, err := ioutil.ReadDir(config.RootDir)
	if err != nil {
		fmt.Println(err)
		return
	}
	localTree := LocalDirNode{
		Name:  "Note",
		Dir:   config.RootDir,
		Dirs:  make([]*LocalDirNode, 0),
		Files: make([]*LocalFileNode, 0),
	}
	BuildLocalTree(&localTree, list)

	// 读取云端笔记，并构造目录树
	nodes, err := api.GetNotebooks(config.Address, config.Token)
	if err != nil {
		fmt.Println(err)
		return
	}
	remoteTree, err := BuildRemoteTree(nodes)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 开始同步
	Sync(&localTree, remoteTree)

	// 同步结束后更新一下config文件，记录本次更新时间
	fd, err := os.OpenFile("/etc/leanote-sync.json", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Printf("os.OpenFile(%v) failed: %v\n", "/etc/leanote-sync.json", err)
		return
	}
	defer fd.Close()
	// 记录此次的更新时间
	config.Updated = now
	bs, err := json.Marshal(*config)
	if err != nil {
		fmt.Printf("json.Marshal(%+v) failed: %v\n", *config, err)
	}

	if _, err = fd.Write(bs); err != nil {
		fmt.Printf("fd.Write(%s) failed: %v\n", bs, err)
	}
	fmt.Printf("本次同步耗时：%v秒\n", time.Now().Unix()-now)
}

func Sync(local *LocalDirNode, remote *RemoteDirNode) error {
	// step1: 同步文件
	for _, file := range local.Files {
		arr := strings.Split(file.Name, ".")
		if len(arr) != 2 || arr[1] != "md" {
			continue
		}
		title := arr[0]
		filePath := local.Dir + "/" + file.Name
		bs, err := ioutil.ReadFile(filePath)
		if err != nil {
			return err
		}
		if note, ok := remote.Files[title]; ok { // 存在
			// 从记录里删除，最后剩下的要被删除掉
			delete(remote.Files, title)
			if file.Updated < config.Updated {
				continue
			}
			if _, err = api.UpdateNote(config.Address, config.Token, note.NotebookID, title, string(bs), "", note.Usn); err != nil {
				return err
			}
			fmt.Printf("UpdateNote: %v\n", filePath)
		} else { // 不存在新增
			if _, err = api.AddNote(config.Address, config.Token, remote.NotebookId, title, string(bs), ""); err != nil {
				return err
			}
			fmt.Printf("AddNote: %v\n", filePath)
		}
	}
	// step2: 从远端删除本地已经删除的文件
	for _, item := range remote.Files {
		if ret, err := api.DeleteNote(config.Address, config.Token, item.NoteID, item.Usn); err != nil {
			return err
		} else if !ret.Ok {
			fmt.Printf("api.DeleteNote(%v, %v, %v, %v)删除笔记失败\n\n", config.Address, config.Token, item.NoteID, item.Usn)
		}
		fmt.Printf("DeleteNote: %v\n", item.Title)
	}
	// step3: 同步目录(笔记本)
	for _, locaLdir := range local.Dirs {
		name := locaLdir.Name
		remoteDir, ok := remote.Dirs[name]
		if !ok { // 笔记本不存在，新增
			notebook, err := api.AddNotebook(config.Address, config.Token, remote.NotebookId, name, 0)
			if err != nil {
				return err
			}
			remoteDir = &RemoteDirNode{
				NotebookId:       notebook.NotebookID,
				ParentNotebookId: remote.NotebookId,
				Title:            name,
				Usn:              notebook.Usn,
				Dirs:             make(map[string]*RemoteDirNode),
				Files:            make(map[string]*api.GetNotesResp),
			}
		} else {
			delete(remote.Dirs, name) // 最后剩余的笔记本要被删除掉
		}
		Sync(locaLdir, remoteDir) // 递归
	}
	// step4: 从远端删除本地已经删除的目录
	for _, notebook := range remote.Dirs {
		if ConstNotebook[notebook.Title] { // 删不掉的，放弃吧
			continue
		}
		if ok, err := api.DeleteNotebook(config.Address, config.Token, notebook.NotebookId, notebook.Usn); err != nil {
			return err
		} else if !ok {
			fmt.Printf("api.DeleteNotebook(%v, %v, %v, %v)删除笔记失败\n\n", config.Address, config.Token, notebook.NotebookId, notebook.Usn)
		}
		fmt.Printf("DeleteNotebook: %v\n", notebook.Title)
	}
	return nil
}

// 本地目录树文件夹节点
type LocalDirNode struct {
	Name  string           // 当前节点名
	Dir   string           // 当前节点路径
	Dirs  []*LocalDirNode  // 子目录
	Files []*LocalFileNode // 子文件名
}

// 本地目录树文件节点
type LocalFileNode struct {
	Name    string
	Updated int64
}

// 云端目录树笔记本节点
type RemoteDirNode struct {
	NotebookId       string
	ParentNotebookId string
	Title            string
	Usn              int
	Dirs             map[string]*RemoteDirNode    // 名字->笔记目录
	Files            map[string]*api.GetNotesResp // 名字->笔记
}

// 云端目录树笔记节点直接使用 api.GetNotesResp即可

// 构造本地目录树
func BuildLocalTree(root *LocalDirNode, list []fs.FileInfo) error {
	for _, item := range list {
		name := item.Name()
		if item.IsDir() {
			node := LocalDirNode{
				Name:  name,
				Dir:   root.Dir + "/" + name,
				Dirs:  make([]*LocalDirNode, 0),
				Files: make([]*LocalFileNode, 0),
			}
			root.Dirs = append(root.Dirs, &node)
			dirs, err := ioutil.ReadDir(node.Dir)
			if err != nil {
				return err
			}
			BuildLocalTree(&node, dirs)
		} else {
			dir := root.Dir + "/" + name
			fileInfo, err := os.Stat(dir)
			if err != nil {
				return err
			}
			file := &LocalFileNode{
				Name:    name,
				Updated: fileInfo.ModTime().Unix(),
			}
			root.Files = append(root.Files, file)
		}
	}
	return nil
}

func PrintLocalTree(root *LocalDirNode, seg string) {
	fmt.Printf("%s %s\n", seg, root.Name)
	for _, file := range root.Files {
		fmt.Printf("-----%s %s\n", seg, file.Name)
	}
	for _, dir := range root.Dirs {
		PrintLocalTree(dir, seg+"-----")
	}
}

func BuildRemoteTree(nodes []*api.GetNotebooksResp) (*RemoteDirNode, error) {
	set := make(map[string][]*api.GetNotebooksResp) // parentNotebookId->[]notebook
	for _, node := range nodes {
		if node.ParentNotebookID == "" {
			node.ParentNotebookID = "root"
		}
		set[node.ParentNotebookID] = append(set[node.ParentNotebookID], node)
	}
	root := &RemoteDirNode{
		NotebookId: "root",
		Title:      "Note",
		Dirs:       make(map[string]*RemoteDirNode),
		Files:      make(map[string]*api.GetNotesResp),
	}
	err := buildRemoteTree(root, set)
	return root, err
}

func buildRemoteTree(root *RemoteDirNode, set map[string][]*api.GetNotebooksResp) error {
	// 获取本层的笔记
	notebookId := root.NotebookId
	notes, err := api.GetNotes(config.Address, config.Token, notebookId)
	if err != nil {
		return err
	}
	for _, note := range notes {
		root.Files[note.Title] = note
	}
	// 整理本层的目录
	for _, dir := range set[root.NotebookId] {
		node := &RemoteDirNode{
			NotebookId:       dir.NotebookID,
			ParentNotebookId: dir.ParentNotebookID,
			Title:            dir.Title,
			Dirs:             make(map[string]*RemoteDirNode),
			Files:            make(map[string]*api.GetNotesResp),
		}
		root.Dirs[dir.Title] = node
		buildRemoteTree(node, set)
	}
	return nil
}

func PrintRemoteTree(root *RemoteDirNode, seg string) {
	fmt.Printf("%s %s\n", seg, root.Title)
	for _, file := range root.Files {
		fmt.Printf("-----%s %s\n", seg, file.Title)
	}
	for _, dir := range root.Dirs {
		PrintRemoteTree(dir, seg+"-----")
	}
}
