package endpoint_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/pygmystack/pygmy/internal/utils/endpoint"
)

func Example() {
	endpoint.Validate("http://127.0.0.1:8080")
}

func Test(t *testing.T) {
	Convey("URL Endpoint tests...", t, func() {
		valid := endpoint.Validate("https://www.golang.org/")
		So(valid, ShouldBeTrue)
	})
}
