package agent_test

import (
	"testing"

	model "github.com/fubarhouse/pygmy-go/service/interface"
	"github.com/fubarhouse/pygmy-go/service/ssh/agent"
	. "github.com/smartystreets/goconvey/convey"
)

func Example() {
	agent.New()
}

func ExampleList() {
	agent.List(model.Service{})
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
		So(obj.Config.Healthcheck.Test, ShouldEqual, []string{"CMD-SHELL", "if [ ! \"$SSH_AUTH_SOCK\" = \"\" ]; then exit 0; else exit 1; fi;"})
		So(obj.Config.Healthcheck.Interval, ShouldEqual, 30000000000)
		So(obj.Config.Healthcheck.Timeout, ShouldEqual, 5000000000)
		So(obj.Config.Healthcheck.StartPeriod, ShouldEqual, 5000000000)
		So(obj.HostConfig.AutoRemove, ShouldBeFalse)
		So(obj.HostConfig.IpcMode, ShouldEqual, "private")
		So(obj.HostConfig.PortBindings, ShouldEqual, nil)
		So(obj.HostConfig.RestartPolicy.Name, ShouldEqual, "unless-stopped")
		So(obj.HostConfig.RestartPolicy.MaximumRetryCount, ShouldEqual, 0)
	})
}
