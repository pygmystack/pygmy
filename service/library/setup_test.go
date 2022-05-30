package library_test

import (
	"strings"
	"testing"

	"github.com/pygmystack/pygmy/service/library"
)

// Tests that the SSH Agent is first in SortedServices.
func TestSortedServices(t *testing.T) {
	c := &library.Config{}
	library.Setup(c)
	c.SortedServices = library.GetServicesSorted(c)
	if !strings.HasSuffix(c.SortedServices[0], "-ssh-agent") {
		t.Fail()
	}
	for _, v := range c.SortedServices {
		t.Log(v)
	}
}
