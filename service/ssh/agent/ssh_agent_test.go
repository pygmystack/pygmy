package agent_test

import (
	"github.com/pygmystack/pygmy/internal/runtime"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"

	"github.com/pygmystack/pygmy/service/ssh/agent"
	. "github.com/smartystreets/goconvey/convey"
)

//func TestExampleList(t *testing.T) {
//	m := &model.Service{}
//	c, e := agent.List(*m)
//	if c != nil && e != nil {
//		t.Fail()
//	}
//}

func TestExampleSearch(t *testing.T) {
	_, err := agent.Search(runtime.Service{}, "id_rsa.pub")
	if err != nil {
		t.Fail()
	}
}

func Test(t *testing.T) {
	Convey("SSH Agent: Field equality tests...", t, func() {
		obj := agent.New()
		So(obj.Config.Image, ShouldContainSubstring, "pygmystack/ssh-agent")
		So(obj.Config.Labels["pygmy.defaults"], ShouldEqual, "true")
		So(obj.Config.Labels["pygmy.enable"], ShouldEqual, "true")
		So(obj.Config.Labels["pygmy.output"], ShouldEqual, "false")
		So(obj.Config.Labels["pygmy.name"], ShouldEqual, "amazeeio-ssh-agent")
		So(obj.Config.Labels["pygmy.network"], ShouldEqual, "amazeeio-network")
		So(obj.Config.Labels["pygmy.purpose"], ShouldEqual, "sshagent")
		So(obj.Config.Labels["pygmy.weight"], ShouldEqual, "10")
		So(obj.HostConfig.AutoRemove, ShouldBeFalse)
		So(obj.HostConfig.IpcMode.IsPrivate(), ShouldBeTrue)
		So(obj.HostConfig.PortBindings, ShouldEqual, nat.PortMap(nil))
		So(obj.HostConfig.RestartPolicy.Name, ShouldEqual, container.RestartPolicyMode("unless-stopped"))
		So(obj.HostConfig.RestartPolicy.MaximumRetryCount, ShouldEqual, 0)
	})
}
