package network_test

import (
	"fmt"
	"github.com/docker/docker/api/types/network"
	"testing"

	n "github.com/fubarhouse/pygmy-go/service/network"
	. "github.com/smartystreets/goconvey/convey"
)

func Example() {
	n.New()
}

func Test(t *testing.T) {
	Convey("Network: Field equality tests...", t, func() {
		obj := n.New()
		So(obj.Name, ShouldEqual, "amazeeio-network")
		So(obj.IPAM.Driver, ShouldEqual, "")
		So(obj.IPAM.Options, ShouldEqual, nil)
		So(fmt.Sprint(obj.IPAM.Config), ShouldEqual, fmt.Sprint([]network.IPAMConfig{{Subnet:  "10.99.99.0/24", Gateway: "10.99.99.1"}}))
		So(fmt.Sprint(obj.Labels), ShouldEqual, fmt.Sprint(map[string]string{"pygmy.name": "amazeeio-network"}))
	})
}
