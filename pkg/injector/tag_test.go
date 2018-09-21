package injector

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTag(t *testing.T) {
	Convey("Test RouteName", t, func() {
		tests := []struct {
			input string
			name  string
		}{
			{"hello", "hello"},
			{"ehm, optional", "ehm"},
		}

		for _, item := range tests {
			So(Tag(item.input).Name(), ShouldEqual, item.name)
		}
	})

	Convey("Test Optional", t, func() {
		tests := []struct {
			input    string
			optional bool
		}{
			{"hello", false},
			{"ehm, optional", true},
		}

		for _, item := range tests {
			So(Tag(item.input).Optional(), ShouldEqual, item.optional)
		}
	})

}
