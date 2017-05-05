package bede

import (
	"encoding/xml"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDict(t *testing.T) {
	Convey("Given a folder to parse", t, func() {
		versions := [3]string{"version1", "version3", "version5"}
		names := [3]string{"foo", "bar", "baz"}
		lengths := [3]int{574, 712, 613}
		paths := [3]string{"testdata/f01.xml", "testdata/f02.xml", "testdata/f03.xml"}

		Convey("client.xml and server.xml should be generated", func() {
			err := GenDict("testdata", "client.xml", "server.xml")
			So(err, ShouldBeNil)
		})

		Convey("client.xml should have the proper data", func() {
			f, err := os.Open("client.xml")
			So(err, ShouldBeNil)

			dict := dictionary{}
			err = xml.NewDecoder(f).Decode(&dict)

			So(err, ShouldBeNil)

			So(len(dict.Documents), ShouldEqual, 3)

			for _, d := range dict.Documents {
				So(names, ShouldContain, d.DocumentName)
				So(lengths, ShouldContain, d.Length)
				So(versions, ShouldContain, d.Version)
				So(d.PathToFile, ShouldBeEmpty)
			}
			os.Remove("client.xml")
		})

		Convey("server.xml should have the proper data", func() {
			f, err := os.Open("server.xml")
			So(err, ShouldBeNil)

			dict := dictionary{}
			err = xml.NewDecoder(f).Decode(&dict)

			So(err, ShouldBeNil)

			So(len(dict.Documents), ShouldEqual, 3)

			for _, d := range dict.Documents {
				So(names, ShouldContain, d.DocumentName)
				So(lengths, ShouldContain, d.Length)
				So(versions, ShouldContain, d.Version)
				So(paths, ShouldContain, d.PathToFile)
			}
			os.Remove("server.xml")
		})

	})
}
