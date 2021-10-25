package haproxy_test

import (
	"fmt"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/fubarhouse/pygmy-go/service/haproxy"
	model "github.com/fubarhouse/pygmy-go/service/interface"
	. "github.com/smartystreets/goconvey/convey"
)

func Example() {
	haproxy.New(&model.Params{})
	haproxy.NewDefaultPorts()
}

func Test(t *testing.T) {
	Convey("HAProxy: Field equality tests...", t, func() {
		obj := haproxy.New(&model.Params{Domain: "docker.amazee.io"})
		objPorts := haproxy.NewDefaultPorts()
		So(obj.Config.Image, ShouldEqual, "pygmystack/haproxy")
		So(obj.Config.Labels["pygmy.defaults"], ShouldEqual, "true")
		So(obj.Config.Labels["pygmy.enable"], ShouldEqual, "true")
		So(obj.Config.Labels["pygmy.name"], ShouldEqual, "amazeeio-haproxy")
		So(obj.Config.Labels["pygmy.network"], ShouldEqual, "amazeeio-network")
		So(obj.Config.Labels["pygmy.url"], ShouldEqual, "http://docker.amazee.io/stats")
		So(obj.Config.Labels["pygmy.weight"], ShouldEqual, "14")
		So(obj.HostConfig.AutoRemove, ShouldBeFalse)
		So(fmt.Sprint(obj.HostConfig.Binds), ShouldEqual, fmt.Sprint([]string{"/var/run/docker.sock:/tmp/docker.sock"}))
		So(obj.HostConfig.PortBindings, ShouldEqual, nil)
		So(obj.HostConfig.RestartPolicy.Name, ShouldEqual, "unless-stopped")
		So(obj.HostConfig.RestartPolicy.MaximumRetryCount, ShouldEqual, 0)
		So(fmt.Sprint(objPorts.HostConfig.PortBindings), ShouldEqual, fmt.Sprint(nat.PortMap{"80/tcp": []nat.PortBinding{{HostIP: "", HostPort: "80"}}, "443/tcp": []nat.PortBinding{{HostIP: "", HostPort: "443"}}}))
	})
}
