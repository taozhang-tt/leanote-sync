package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"
	"leanote-sync/api"
	"leanote-sync/config"
	"strings"
)

func main() {
	list, err := ioutil.ReadDir("/Users/tt/Documents/Note")
	if err != nil {
		fmt.Println(err)
		return
	}
	localTree := LocalDirNode{
		Name:  "Note",
		Dir:   "/Users/tt/Documents/Note",
		Dirs:  make([]*LocalDirNode, 0),
		Files: make([]string, 0),
	}
	BuildLocalTree(&localTree, list)

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

	Sync(&localTree, remoteTree)
}

func Sync(local *LocalDirNode, remote *RemoteDirNode) error {
	// 同步文件
	for _, file := range local.Files {
		arr := strings.Split(file, ".")
		if len(arr) != 2 || arr[1] != "md" {
			continue
		}
		title := arr[0]
		// 判断远端是否存在该笔记
		// TODO: md5校验一下是否需要更新
		bs, err := ioutil.ReadFile(local.Dir + "/" + file)
		if err != nil {
			return err
		}
		if note, ok := remote.Files[title]; ok { // 存在就更新
			if _, err = api.UpdateNote(config.Address, config.Token, note.NotebookID, title, string(bs), "", note.Usn); err != nil {
				return err
			}
			// 从记录里删除，最后剩下的要被删除掉
			delete(remote.Files, title)
		} else { // 不存在则新增
			if _, err = api.AddNote(config.Address, config.Token, remote.NotebookId, title, string(bs), ""); err != nil {
				return err
			}
		}
	}
	// 从远端删除本地已经删除的文件
	for _, item := range remote.Files {
		if ret, err := api.DeleteNote(config.Address, config.Token, item.NoteID, item.Usn); err != nil {
			return err
		} else if !ret.Ok {
			fmt.Printf("api.DeleteNote(%v, %v, %v, %v)删除笔记失败\n\n", config.Address, config.Token, item.NoteID, item.Usn)
		}
	}
	// 同步目录
	for _, locaLdir := range local.Dirs {
		name := locaLdir.Name
		remoteDir, ok := remote.Dirs[name]
		if !ok {
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
			delete(remote.Dirs, name)
		}
		Sync(locaLdir, remoteDir)
	}
	for _, notebook := range remote.Dirs {
		if ok, err := api.DeleteNotebook(config.Address, config.Token, notebook.NotebookId, notebook.Usn); err != nil {
			return err
		} else if !ok {
			fmt.Printf("api.DeleteNotebook(%v, %v, %v, %v)删除笔记失败\n\n", config.Address, config.Token, notebook.NotebookId, notebook.Usn)
		}
	}
	return nil
}

type LocalDirNode struct {
	Name  string          // 当前节点名
	Dir   string          // 当前节点路径
	Dirs  []*LocalDirNode // 子目录
	Files []string        // 子文件名
}

type RemoteDirNode struct {
	NotebookId       string
	ParentNotebookId string
	Title            string
	Usn              int
	Dirs             map[string]*RemoteDirNode    // 名字->笔记目录
	Files            map[string]*api.GetNotesResp // 名字->笔记
}

// 构造本地目录树
func BuildLocalTree(root *LocalDirNode, list []fs.FileInfo) error {
	for _, item := range list {
		name := item.Name()
		if item.IsDir() {
			node := LocalDirNode{
				Name:  name,
				Dir:   root.Dir + "/" + name,
				Dirs:  make([]*LocalDirNode, 0),
				Files: make([]string, 0),
			}
			root.Dirs = append(root.Dirs, &node)
			dirs, err := ioutil.ReadDir(node.Dir)
			if err != nil {
				return err
			}
			BuildLocalTree(&node, dirs)
		} else {
			root.Files = append(root.Files, name)
		}
	}
	return nil
}

func PrintLocalTree(root *LocalDirNode, seg string) {
	fmt.Printf("%s %s\n", seg, root.Name)
	for _, file := range root.Files {
		fmt.Printf("-----%s %s\n", seg, file)
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
