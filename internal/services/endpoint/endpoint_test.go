package endpoint_test

import (
	"github.com/pygmystack/pygmy/internal/services/endpoint"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
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
