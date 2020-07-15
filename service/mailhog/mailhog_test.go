package mailhog_test

import (
	"fmt"
	"testing"

	"github.com/docker/go-connections/nat"
	"github.com/fubarhouse/pygmy-go/service/mailhog"
	. "github.com/smartystreets/goconvey/convey"
)

// This example creates a Service with the default values for mailhog.
func Example() {
	mailhog.New()
	mailhog.NewDefaultPorts()
}

func Test(t *testing.T) {
	Convey("MailHog: Field equality tests...", t, func() {
		obj := mailhog.New()
		objPorts := mailhog.NewDefaultPorts()
		So(obj.Config.User, ShouldEqual, "0")
		So(obj.Config.Image, ShouldEqual, "mailhog/mailhog")
		So(fmt.Sprint(obj.Config.ExposedPorts), ShouldEqual, fmt.Sprint(nat.PortSet{"80/tcp": struct{}{}, "1025/tcp": struct{}{}, "8025/tcp": struct{}{}}))
		So(fmt.Sprint(obj.Config.Env), ShouldEqual, fmt.Sprint([]string{"MH_UI_BIND_ADDR=0.0.0.0:80", "MH_API_BIND_ADDR=0.0.0.0:80", "AMAZEEIO=AMAZEEIO", "AMAZEEIO_URL=mailhog.docker.amazee.io"}))
		So(obj.Config.Labels["pygmy.defaults"], ShouldEqual, "true")
		So(obj.Config.Labels["pygmy.enable"], ShouldEqual, "true")
		So(obj.Config.Labels["pygmy.name"], ShouldEqual, "amazeeio-mailhog")
		So(obj.Config.Labels["pygmy.network"], ShouldEqual, "amazeeio-network")
		So(obj.Config.Labels["pygmy.url"], ShouldEqual, "http://mailhog.docker.amazee.io")
		So(obj.Config.Labels["pygmy.weight"], ShouldEqual, "15")
		So(obj.HostConfig.AutoRemove, ShouldBeFalse)
		So(obj.HostConfig.PortBindings, ShouldEqual, nil)
		So(obj.HostConfig.RestartPolicy.Name, ShouldEqual, "on-failure")
		So(obj.HostConfig.RestartPolicy.MaximumRetryCount, ShouldEqual, 0)
		So(fmt.Sprint(objPorts.HostConfig.PortBindings), ShouldEqual, fmt.Sprint(nat.PortMap{"1025/tcp": []nat.PortBinding{{HostIP: "", HostPort: "1025"}}}))
	})
}
