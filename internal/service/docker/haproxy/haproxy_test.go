package haproxy_test

import (
	"fmt"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/pygmystack/pygmy/internal/runtime/docker"
	"github.com/pygmystack/pygmy/internal/service/docker/haproxy"
)

func Example() {
	haproxy.New(&docker.Params{})
	haproxy.NewDefaultPorts()
}

func Test(t *testing.T) {
	Convey("HAProxy: Field equality tests...", t, func() {
		obj := haproxy.New(&docker.Params{Domain: "docker.amazee.io", TLSCertPath: "/path/to/ssl/cert.pem"})
		objPorts := haproxy.NewDefaultPorts()
		So(obj.Config.Image, ShouldContainSubstring, "pygmystack/haproxy")
		So(obj.Config.Labels["pygmy.defaults"], ShouldEqual, "true")
		So(obj.Config.Labels["pygmy.enable"], ShouldEqual, "true")
		So(obj.Config.Labels["pygmy.name"], ShouldEqual, "amazeeio-haproxy")
		So(obj.Config.Labels["pygmy.network"], ShouldEqual, "amazeeio-network")
		So(obj.Config.Labels["pygmy.url"], ShouldEqual, "docker.amazee.io/stats")
		So(obj.Config.Labels["pygmy.weight"], ShouldEqual, "14")
		So(obj.HostConfig.AutoRemove, ShouldBeFalse)
		So(fmt.Sprint(obj.HostConfig.Binds), ShouldEqual, fmt.Sprint([]string{"/var/run/docker.sock:/tmp/docker.sock", "/path/to/ssl/cert.pem:/app/server.pem:ro"}))
		So(obj.HostConfig.PortBindings, ShouldEqual, nat.PortMap(nil))
		So(obj.HostConfig.RestartPolicy.Name, ShouldEqual, container.RestartPolicyMode("unless-stopped"))
		So(obj.HostConfig.RestartPolicy.MaximumRetryCount, ShouldEqual, 0)
		So(fmt.Sprint(objPorts.HostConfig.PortBindings), ShouldEqual, fmt.Sprint(nat.PortMap{"80/tcp": []nat.PortBinding{{HostIP: "", HostPort: "80"}}, "443/tcp": []nat.PortBinding{{HostIP: "", HostPort: "443"}}}))
	})
}
