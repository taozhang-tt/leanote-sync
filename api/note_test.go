package api

import (
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"

	"leanote-sync/config"
)

func TestGetNoteContent(t *testing.T) {
	Convey("GetNoteConteng", t, func() {
		ret, err := GetNoteContent(config.GetConfig().Address, config.GetConfig().Token, "557eaa9905fcd14d95000001")
		So(err, ShouldBeNil)
		if err != nil {
			t.Errorf("GetNoteContent Err: %v", err)
		} else {
			t.Logf("GetNoteContent Ret: %+v", ret)
		}
	})
}

func TestGetNotes(t *testing.T) {
	Convey("GetNotes", t, func() {
		ret, err := GetNotes(config.GetConfig().Address, config.GetConfig().Token, "61779b2af76761037a000009")
		So(err, ShouldBeNil)
		if err != nil {
			t.Errorf("GetNotesErr: %v", err)
		} else {
			t.Logf("GetNotesRet: %+v", ret)
		}
	})
}

func TestUpdateNote(t *testing.T) {
	Convey("AddNote", t, func() {
		// 新增一个笔记
		newNote, err := AddNote(config.GetConfig().Address, config.GetConfig().Token, "557eab5705fcd14d95000002", "测试笔记", "这是笔记内容", "这是笔记摘要")
		So(err, ShouldBeNil)
		if err != nil {
			t.Errorf("AddNote Err: %v", err)
		}

		// 获取新增的笔记和内容
		Convey("GetNote", func() {
			noteAndContent, err := GetNoteAndContent(config.GetConfig().Address, config.GetConfig().Token, newNote.NoteId)
			So(err, ShouldBeNil)
			if err != nil {
				t.Errorf("GetNoteAndContent Err: %v", err)
			}

			// 更新新增的笔记
			Convey("UpdateNote", func() {
				title := "测试笔记" + time.Now().String()
				_, err = UpdateNote(config.GetConfig().Address, config.GetConfig().Token, noteAndContent.NoteID, title, "这是更新以后的笔记内容", "", noteAndContent.Usn)
				So(err, ShouldBeNil)
				if err != nil {
					t.Errorf("UpdateNote Err: %v", err)
				}

				// 获取更新后的笔记内容
				Convey("GetNote after update", func() {
					note, err := GetNoteContent(config.GetConfig().Address, config.GetConfig().Token, noteAndContent.NoteID)
					So(err, ShouldBeNil)
					if err != nil {
						t.Errorf("GetNote Err: %v", err)
					}
					So(note.Content, ShouldEqual, "这是更新以后的笔记内容")
				})
			})

		})
	})
}
