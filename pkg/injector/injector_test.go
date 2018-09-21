package injector

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"go.uber.org/zap"
)

func TestInjector(t *testing.T) {

	Convey("Test New", t, func() {
		So(New(), ShouldNotBeNil)
		So(New(), ShouldImplement, (*Injector)(nil))
	})

	Convey("Test Options", t, func() {
		l := zap.NewExample()
		t := "test_tag"
		i := New(
			WithLogger(l),
			WithTag(t),
		)

		So(i.(*injector).logger, ShouldEqual, l)
		So(i.(*injector).tag, ShouldEqual, t)
	})

	Convey("Test Provide", t, func() {
		i := New()
		So(i.Provide(10, "my"), ShouldBeNil)
		So(len(i.(*injector).deps), ShouldEqual, 1)
		So(i.Provide(10, "my"), ShouldBeNil)
	})

	Convey("Test Inject type", t, func() {
		i := New(WithTag("mytag"))
		So(i.Provide(10, "my", "ns"), ShouldBeNil)

		x := struct {
			Field int `mytag:"ns:my"`
		}{
			0,
		}

		So(i.Inject(x, false), ShouldEqual, ErrTargetCannotBeSet)
		So(i.Inject(&x, false), ShouldBeNil)

	})

	Convey("Test Inject interface", t, func() {

	})

}
