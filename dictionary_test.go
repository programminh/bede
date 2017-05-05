package bede

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDict(t *testing.T) {
	Convey("Given a folder to parse, make sure client.xml and server.xml are generated", t, func() {
		err := GenDict("testdata", "testdata/client.xml", "testdata/server.xml")

		So(err, ShouldBeNil)
	})
}
