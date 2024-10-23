package dnsmasq_test

import (
	"fmt"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/pygmystack/pygmy/internal/runtime/docker"
	"github.com/pygmystack/pygmy/internal/service/docker/dnsmasq"
)

func Example() {
	dnsmasq.New(&docker.Params{})
}

func Test(t *testing.T) {
	Convey("DNSMasq: Field equality tests...", t, func() {
		obj := dnsmasq.New(&docker.Params{Domain: "docker.amazee.io"})

		So(obj.Config.Image, ShouldContainSubstring, "pygmystack/dnsmasq")
		So(fmt.Sprint(obj.Config.Cmd), ShouldEqual, fmt.Sprint([]string{"--log-facility=-", "-A", "/docker.amazee.io/127.0.0.1"}))
		So(obj.Config.Labels["pygmy.defaults"], ShouldEqual, "true")
		So(obj.Config.Labels["pygmy.enable"], ShouldEqual, "true")
		So(obj.Config.Labels["pygmy.name"], ShouldEqual, "amazeeio-dnsmasq")
		So(obj.Config.Labels["pygmy.weight"], ShouldEqual, "13")
		So(obj.HostConfig.AutoRemove, ShouldBeFalse)
		So(fmt.Sprint(obj.HostConfig.CapAdd), ShouldEqual, fmt.Sprint([]string{"NET_ADMIN"}))
		So(obj.HostConfig.IpcMode.IsPrivate(), ShouldBeTrue)
		So(fmt.Sprint(obj.HostConfig.PortBindings), ShouldEqual, fmt.Sprint(nat.PortMap{"53/tcp": []nat.PortBinding{{HostIP: "", HostPort: "6053"}}, "53/udp": []nat.PortBinding{{HostIP: "", HostPort: "6053"}}}))
		So(obj.HostConfig.RestartPolicy.Name, ShouldEqual, container.RestartPolicyMode("unless-stopped"))
		So(obj.HostConfig.RestartPolicy.MaximumRetryCount, ShouldBeZeroValue)
	})
}
