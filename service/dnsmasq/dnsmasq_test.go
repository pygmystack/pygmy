package dnsmasq_test

import (
	"fmt"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/fubarhouse/pygmy-go/service/dnsmasq"
	model "github.com/fubarhouse/pygmy-go/service/interface"
	. "github.com/smartystreets/goconvey/convey"
)

func Example() {
	dnsmasq.New(&model.Params{})
}

func Test(t *testing.T) {
	Convey("DNSMasq: Field equality tests...", t, func() {
		obj := dnsmasq.New(&model.Params{})

		So(obj.Config.Image, ShouldEqual, "andyshinn/dnsmasq:2.83")
		So(fmt.Sprint(obj.Config.Cmd), ShouldEqual, fmt.Sprint([]string{"-A", "/docker.amazee.io/127.0.0.1"}))
		So(obj.Config.Labels["pygmy.defaults"], ShouldEqual, "true")
		So(obj.Config.Labels["pygmy.enable"], ShouldEqual, "true")
		So(obj.Config.Labels["pygmy.name"], ShouldEqual, "amazeeio-dnsmasq")
		So(obj.Config.Labels["pygmy.weight"], ShouldEqual, "13")
		So(obj.HostConfig.AutoRemove, ShouldBeFalse)
		So(fmt.Sprint(obj.HostConfig.CapAdd), ShouldEqual, fmt.Sprint([]string{"NET_ADMIN"}))
		So(obj.HostConfig.IpcMode, ShouldEqual, "private")
		So(fmt.Sprint(obj.HostConfig.PortBindings), ShouldEqual, fmt.Sprint(nat.PortMap{"53/tcp": []nat.PortBinding{{HostIP: "", HostPort: "6053"}}, "53/udp": []nat.PortBinding{{HostIP: "", HostPort: "6053"}}}))
		So(obj.HostConfig.RestartPolicy.Name, ShouldEqual, "unless-stopped")
		So(obj.HostConfig.RestartPolicy.MaximumRetryCount, ShouldBeZeroValue)
	})
}
