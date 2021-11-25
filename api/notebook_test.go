package api

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"leanote-sync/config"
)

func TestGetNotebooks(t *testing.T) {
	Convey("GetNotebooks", t, func() {
		ret, err := GetNotebooks(config.GetConfig().Address, config.GetConfig().Token)
		So(err, ShouldBeNil)
		if err != nil {
			t.Errorf("GetNotebooks Err: %v", err)
		}
		t.Logf("%+v", ret)
	})
}
