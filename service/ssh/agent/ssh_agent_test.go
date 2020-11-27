package agent_test

import (
	"fmt"
	"testing"

	model "github.com/fubarhouse/pygmy-go/service/interface"
	"github.com/fubarhouse/pygmy-go/service/ssh/agent"
	. "github.com/smartystreets/goconvey/convey"
)

func Example() {
	agent.New()
}

func ExampleList() {
	_, e := agent.List(model.Service{})
	if e != nil {
		fmt.Println(e)
	}
}

func ExampleSearch() {
	agent.Search(model.Service{}, "id_rsa.pub")
}

func Test(t *testing.T) {
	Convey("SSH Agent: Field equality tests...", t, func() {
		obj := agent.New()
		So(obj.Config.Image, ShouldEqual, "amazeeio/ssh-agent")
		So(obj.Config.Labels["pygmy.defaults"], ShouldEqual, "true")
		So(obj.Config.Labels["pygmy.enable"], ShouldEqual, "true")
		So(obj.Config.Labels["pygmy.output"], ShouldEqual, "false")
		So(obj.Config.Labels["pygmy.name"], ShouldEqual, "amazeeio-ssh-agent")
		So(obj.Config.Labels["pygmy.network"], ShouldEqual, "amazeeio-network")
		So(obj.Config.Labels["pygmy.purpose"], ShouldEqual, "sshagent")
		So(obj.Config.Labels["pygmy.weight"], ShouldEqual, "30")
		So(obj.HostConfig.AutoRemove, ShouldBeFalse)
		So(obj.HostConfig.IpcMode, ShouldEqual, "private")
		So(obj.HostConfig.PortBindings, ShouldEqual, nil)
		So(obj.HostConfig.RestartPolicy.Name, ShouldEqual, "unless-stopped")
		So(obj.HostConfig.RestartPolicy.MaximumRetryCount, ShouldEqual, 0)
	})
}
