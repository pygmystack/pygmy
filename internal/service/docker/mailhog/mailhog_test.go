package mailhog_test

import (
	"fmt"
	"net"
	"testing"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/pygmystack/pygmy/internal/runtime/docker"
	"github.com/pygmystack/pygmy/internal/service/docker/mailhog"
)

func Example() {
	mailhog.New(&docker.Params{})
	mailhog.NewDefaultPorts()
}

func Test(t *testing.T) {
	Convey("MailHog: Field equality tests...", t, func() {
		obj := mailhog.New(&docker.Params{Domain: "docker.amazee.io", TLSCertPath: ""})
		objPorts := mailhog.NewDefaultPorts()
		So(obj.Config.User, ShouldEqual, "0")
		So(obj.Config.Image, ShouldContainSubstring, "pygmystack/mailhog")
		So(fmt.Sprint(obj.Config.ExposedPorts), ShouldEqual, fmt.Sprint(nat.PortSet{"80/tcp": struct{}{}, "1025/tcp": struct{}{}, "8025/tcp": struct{}{}}))
		So(fmt.Sprint(obj.Config.Env), ShouldEqual, fmt.Sprint([]string{"MH_UI_BIND_ADDR=0.0.0.0:80", "MH_API_BIND_ADDR=0.0.0.0:80", "AMAZEEIO=AMAZEEIO", "AMAZEEIO_URL=mailhog.docker.amazee.io", "LAGOON_ROUTE=http://mailhog.docker.amazee.io"}))
		So(obj.Config.Labels["pygmy.defaults"], ShouldEqual, "true")
		So(obj.Config.Labels["pygmy.enable"], ShouldEqual, "true")
		So(obj.Config.Labels["pygmy.name"], ShouldEqual, "amazeeio-mailhog")
		So(obj.Config.Labels["pygmy.network"], ShouldEqual, "amazeeio-network")
		So(obj.Config.Labels["pygmy.url"], ShouldEqual, "http://mailhog.docker.amazee.io")
		So(obj.Config.Labels["pygmy.weight"], ShouldEqual, "15")
		So(obj.HostConfig.AutoRemove, ShouldBeFalse)
		So(obj.HostConfig.PortBindings, ShouldEqual, nat.PortMap(nil))
		So(obj.HostConfig.RestartPolicy.Name, ShouldEqual, container.RestartPolicyMode("unless-stopped"))
		So(obj.HostConfig.RestartPolicy.MaximumRetryCount, ShouldEqual, 0)
		So(fmt.Sprint(objPorts.HostConfig.PortBindings), ShouldEqual, fmt.Sprint(nat.PortMap{"1025/tcp": []nat.PortBinding{{HostIP: "", HostPort: "1025"}}}))
	})
}

func TestGetFreePort(t *testing.T) {
	port, err := mailhog.GetRandomUnusedPort()
	if err != nil {
		t.Fatalf("GetRandomUnusedPort returned error: %v", err)
	}
	if port <= 0 || port > 65535 {
		t.Fatalf("GetRandomUnusedPort returned invalid port: %d", port)
	}

	// Test that the port can be bound again (means it actually is free now)
	ln, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", port))
	if err != nil {
		t.Fatalf("Failed to bind to port %d: %v", port, err)
	}
	ln.Close()
}
