package library_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	model "github.com/pygmystack/pygmy/service/interface"
	"github.com/pygmystack/pygmy/service/library"
)

// Tests the setup process.
func TestSetup(t *testing.T) {
	// Get our configuration object
	c := &library.Config{
		Services: map[string]model.Service{
			"amazeeio-dnsmasq": {
				// Set an override config value so it can be tested.
				Image: "example-amazeeio-dnsmasq:main",
			},
			"amazeeio-mailhog": {
				// Set an override config value so it can be tested.
				Image: "example-amazeeio-mailhog:main",
			},
		},
	}
	library.Setup(c)
	c.SortedServices = library.GetServicesSorted(c)

	Convey("Setup Tests", t, func() {
		// SSH Agent must be 5 items long by default.
		So(c.SortedServices, ShouldHaveLength, 5)
		// SSH Agent must be the first item in the sorted list.
		So(c.SortedServices[0], ShouldEqual, "amazeeio-ssh-agent")
		// Test sorting result.
		So(c.SortedServices[1], ShouldEqual, "amazeeio-dnsmasq")
		So(c.SortedServices[2], ShouldEqual, "amazeeio-haproxy")
		So(c.SortedServices[3], ShouldEqual, "amazeeio-mailhog")
		So(c.SortedServices[4], ShouldEqual, "amazeeio-ssh-agent-add-key")
		// Test Image Override configuration.
		So(c.Services["amazeeio-dnsmasq"].Image, ShouldEqual, "example-amazeeio-dnsmasq:main")
		So(c.Services["amazeeio-dnsmasq"].Config.Image, ShouldEqual, "example-amazeeio-dnsmasq:main")
		So(c.Services["amazeeio-haproxy"].Image, ShouldEqual, "pygmystack/haproxy:main")
		So(c.Services["amazeeio-haproxy"].Config.Image, ShouldEqual, "pygmystack/haproxy:main")
		So(c.Services["amazeeio-mailhog"].Image, ShouldEqual, "example-amazeeio-mailhog:main")
		So(c.Services["amazeeio-mailhog"].Config.Image, ShouldEqual, "example-amazeeio-mailhog:main")
		So(c.Services["amazeeio-ssh-agent"].Image, ShouldEqual, "pygmystack/ssh-agent:main")
		So(c.Services["amazeeio-ssh-agent"].Config.Image, ShouldEqual, "pygmystack/ssh-agent:main")
	})
}
