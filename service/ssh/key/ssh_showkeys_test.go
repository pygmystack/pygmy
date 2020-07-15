package key_test

import (
	"fmt"
	"testing"

	"github.com/fubarhouse/pygmy-go/service/ssh/key"
	. "github.com/smartystreets/goconvey/convey"
)

func ExampleShow() {
	key.NewShower()
}

func TestShow(t *testing.T) {
	Convey("SSH Key Shower: Field equality tests...", t, func() {
		obj := key.NewShower()
		So(obj.Config.Image, ShouldEqual, "amazeeio/ssh-agent")
		So(fmt.Sprint(obj.Config.Cmd), ShouldEqual, fmt.Sprint([]string{"ssh-add", "-L"}))
		So(obj.Config.Labels["pygmy.defaults"], ShouldEqual, "true")
		So(obj.Config.Labels["pygmy.enable"], ShouldEqual, "true")
		So(obj.Config.Labels["pygmy.output"], ShouldEqual, "false")
		So(obj.Config.Labels["pygmy.discrete"], ShouldEqual, "true")
		So(obj.Config.Labels["pygmy.name"], ShouldEqual, "amazeeio-ssh-agent-show-keys")
		So(obj.Config.Labels["pygmy.network"], ShouldEqual, "amazeeio-network")
		So(obj.Config.Labels["pygmy.purpose"], ShouldEqual, "showkeys")
		So(obj.Config.Labels["pygmy.weight"], ShouldEqual, "32")
		So(obj.HostConfig.AutoRemove, ShouldBeTrue)
		So(obj.HostConfig.IpcMode, ShouldEqual, "private")
		So(fmt.Sprint(obj.HostConfig.VolumesFrom), ShouldEqual, fmt.Sprint([]string{"amazeeio-ssh-agent"}))
	})
}
