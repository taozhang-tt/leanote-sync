package api

import (
	"leanote-sync/config"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetNotebooks(t *testing.T) {
	Convey("GetNotebooks", t, func() {
		ret, err := GetNotebooks(config.Token, config.Address)
		So(err, ShouldBeNil)
		if err != nil {
			t.Errorf("GetNotebooks Err: %v", err)
		}
		t.Logf("%+v", ret)
	})
}
