package setup_test

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/pygmystack/pygmy/external/docker/setup"
	"github.com/pygmystack/pygmy/internal/runtime/docker"
	"github.com/pygmystack/pygmy/internal/runtime/docker/internals"
)

// Tests the setup process.
func TestSetup(t *testing.T) {
	// Get our configuration object
	c := &setup.Config{
		Services: map[string]docker.Service{
			"pygmy-dns": {
				// Set an override config value so it can be tested.
				Image: "example-pygmy-dns",
			},
			"pygmy-mail": {
				// Set an override config value so it can be tested.
				Image: "example-pygmy-mail",
			},
		},
	}

	cli, ctx, err := internals.NewClient()
	if err != nil {
		t.Fatal(err)
	}

	setup.Setup(ctx, cli, c)
	c.SortedServices = setup.GetServicesSorted(ctx, cli, c)

	Convey("Setup Tests", t, func() {
		// SSH Agent must be 5 items long by default.
		So(c.SortedServices, ShouldHaveLength, 5)
		// SSH Agent must be the first item in the sorted list.
		So(c.SortedServices[0], ShouldEqual, "pygmy-ssh")
		// Test sorting result (ordered by pygmy.weight: 13, 14, 15, 31).
		So(c.SortedServices[1], ShouldEqual, "pygmy-dns")
		So(c.SortedServices[2], ShouldEqual, "pygmy-proxy")
		So(c.SortedServices[3], ShouldEqual, "pygmy-mail")
		So(c.SortedServices[4], ShouldEqual, "pygmy-ssh-add-key")
		// Test Image Override configuration.
		So(c.Services["pygmy-dns"].Image, ShouldEqual, "example-pygmy-dns")
		So(c.Services["pygmy-dns"].Config.Image, ShouldEqual, "example-pygmy-dns")
		So(c.Services["pygmy-proxy"].Image, ShouldEqual, "pygmystack/haproxy")
		So(c.Services["pygmy-proxy"].Config.Image, ShouldEqual, "pygmystack/haproxy")
		So(c.Services["pygmy-mail"].Image, ShouldEqual, "example-pygmy-mail")
		So(c.Services["pygmy-mail"].Config.Image, ShouldEqual, "example-pygmy-mail")
		So(c.Services["pygmy-ssh"].Image, ShouldEqual, "pygmystack/ssh-agent")
		So(c.Services["pygmy-ssh"].Config.Image, ShouldEqual, "pygmystack/ssh-agent")
	})
}
