package num2words

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestConvert(t *testing.T) {
	Convey("Should convert correctly", t, func() {
		Convey("Small numbers should convert correctly", func() {
			So(Convert(0), ShouldEqual, "zero")
			So(Convert(1), ShouldEqual, "one")
			So(Convert(5), ShouldEqual, "five")
			So(Convert(10), ShouldEqual, "ten")
			So(Convert(11), ShouldEqual, "eleven")
			So(Convert(12), ShouldEqual, "twelve")
			So(Convert(17), ShouldEqual, "seventeen")
		})
		Convey("Tens should convert correctly", func() {
			So(Convert(20), ShouldEqual, "twenty")
			So(Convert(30), ShouldEqual, "thirty")
			So(Convert(40), ShouldEqual, "forty")
			So(Convert(50), ShouldEqual, "fifty")
			So(Convert(60), ShouldEqual, "sixty")
			So(Convert(90), ShouldEqual, "ninety")
		})
		Convey("Combined numbers should convert correctly", func() {
			So(Convert(21), ShouldEqual, "twenty one")
			So(Convert(34), ShouldEqual, "thirty four")
			So(Convert(49), ShouldEqual, "forty nine")
			So(Convert(53), ShouldEqual, "fifty three")
			So(Convert(68), ShouldEqual, "sixty eight")
			So(Convert(99), ShouldEqual, "ninety nine")
		})
		Convey("Big numbers should convert correctly", func() {
			So(Convert(100), ShouldEqual, "one hundred")
			So(Convert(200), ShouldEqual, "two hundred")
			So(Convert(500), ShouldEqual, "five hundred")
			So(Convert(123), ShouldEqual, "one hundred twenty three")
			So(Convert(666), ShouldEqual, "six hundred sixty six")
			So(Convert(1024), ShouldEqual, "one thousand twenty four")
		})
		Convey("Negative numbers should convert correclty", func() {
			So(Convert(-123), ShouldEqual, "minus one hundred twenty three")
		})
	})
}

func ExampleConvert() {
	var str string
	str = Convert(17)   // outputs "seventeen"
	str = Convert(1024) // outputs "one thousand twenty four"
	str = Convert(-123) // outputs "minus one hundred twenty three"
	_ = str
}
