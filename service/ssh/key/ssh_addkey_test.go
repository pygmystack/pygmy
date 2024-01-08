package key_test

import (
	"fmt"
	"testing"

	"github.com/pygmystack/pygmy/service/ssh/key"
	. "github.com/smartystreets/goconvey/convey"
)

//func ExampleAdd() {
//	key.NewAdder()
//}

func TestAdd(t *testing.T) {
	Convey("SSH Key Adder: Field equality tests...", t, func() {
		obj := key.NewAdder()
		So(obj.Config.Image, ShouldContainSubstring, "pygmystack/ssh-agent")
		So(obj.Config.Labels["pygmy.defaults"], ShouldEqual, "true")
		So(obj.Config.Labels["pygmy.enable"], ShouldEqual, "true")
		So(obj.Config.Labels["pygmy.output"], ShouldEqual, "false")
		So(obj.Config.Labels["pygmy.discrete"], ShouldEqual, "true")
		So(obj.Config.Labels["pygmy.interactive"], ShouldEqual, "true")
		So(obj.Config.Labels["pygmy.name"], ShouldEqual, "amazeeio-ssh-agent-add-key")
		So(obj.Config.Labels["pygmy.network"], ShouldEqual, "amazeeio-network")
		So(obj.Config.Labels["pygmy.purpose"], ShouldEqual, "addkeys")
		So(obj.Config.Labels["pygmy.weight"], ShouldEqual, "31")
		So(obj.Config.Tty, ShouldEqual, true)
		So(obj.Config.OpenStdin, ShouldEqual, true)
		So(obj.HostConfig.AutoRemove, ShouldBeFalse)
		So(obj.HostConfig.IpcMode.IsPrivate(), ShouldBeTrue)
		So(fmt.Sprint(obj.HostConfig.VolumesFrom), ShouldEqual, fmt.Sprint([]string{"amazeeio-ssh-agent"}))
	})
}
