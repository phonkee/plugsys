package plugin

import (
	"testing"

	"github.com/phonkee/plugsys/api"
	"github.com/pkg/errors"
	. "github.com/smartystreets/goconvey/convey"
)

type TestPlugin string

func (t TestPlugin) ID() string { return string(t) }

type TestPlugin2 TestPlugin
func (t TestPlugin2) ID() string { return string(t) }
func (t TestPlugin2) Additional() { }


func TestStorage(t *testing.T) {

	Convey("Test New", t, func() {
		s := NewStorage(nil)
		So(s, ShouldNotBeNil)
	})

	Convey("Test Add Plugin", t, func() {
		s := NewStorage(nil)

		So(s.Add(TestPlugin("test")), ShouldBeNil)
		So(errors.Cause(s.Add(TestPlugin("test"))), ShouldEqual, api.ErrPluginAlreadyRegistered)
		So(s.Add(TestPlugin("test2")), ShouldBeNil)
		So(errors.Cause(s.Add(TestPlugin("test2"))), ShouldEqual, api.ErrPluginAlreadyRegistered)
	})

	Convey("Test Each Plugin", t, func() {
		s := NewStorage(nil)

		So(s.Add(TestPlugin("test")), ShouldBeNil)
		So(errors.Cause(s.Add(TestPlugin("test"))), ShouldEqual, api.ErrPluginAlreadyRegistered)
		So(s.Add(TestPlugin("test2")), ShouldBeNil)
		So(errors.Cause(s.Add(TestPlugin("test2"))), ShouldEqual, api.ErrPluginAlreadyRegistered)

		found := make([]string, 0)

		s.Each(func(plugin api.Plugin) (err error) {
			found = append(found, plugin.ID())
			return
		})

		So(len(found), ShouldEqual, 2)
	})

	Convey("Test Filter Plugin", t, func() {
		s := NewStorage(nil)

		So(s.Add(TestPlugin("test")), ShouldBeNil)
		So(s.Add(TestPlugin2("test2")), ShouldBeNil)

		found := make([]string, 0)

		s.Filter(func(plugin api.Plugin) (err error) {
			found = append(found, plugin.ID())
			return
		})

		So(len(found), ShouldEqual, 2)

		foundFiltered := 0

		s.Filter(func(plugin interface{Additional()}) (err error) {
			foundFiltered += 1
			return
		})

		So(foundFiltered, ShouldEqual, 1)
	})


}
