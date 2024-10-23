package main_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	. "github.com/smartystreets/goconvey/convey"

	"github.com/pygmystack/pygmy/internal/runtime/docker/internals/containers"
	"github.com/pygmystack/pygmy/internal/runtime/docker/internals/images"
)

const (
	dindContainerName = "exampleTestContainer"
	binaryReference   = "pygmy-linux-amd64"
)

var (
	dindID string
)

type config struct {
	name               string
	configpath         string
	endpoints          []string
	images             []string
	services           []string
	servicewithports   []string
	skipendpointchecks bool
}

// setup is a configurable pipeline which allows different configurations to
// run to keep the consistency for as many tests as are required.
func setup(t *testing.T, config *config) {

	var cleanCmd = fmt.Sprintf("/builds/%v clean", binaryReference)
	var statusCmd = fmt.Sprintf("/builds/%v status", binaryReference)
	var upCmd = fmt.Sprintf("/builds/%v up", binaryReference)
	var downCmd = fmt.Sprintf("/builds/%v down", binaryReference)

	if config.configpath != "" {
		cleanCmd = fmt.Sprintf("/builds/%v clean --config %v", binaryReference, config.configpath)
		statusCmd = fmt.Sprintf("/builds/%v status --config %v", binaryReference, config.configpath)
		upCmd = fmt.Sprintf("/builds/%v up --config %v", binaryReference, config.configpath)
		downCmd = fmt.Sprintf("/builds/%v dow --config %v", binaryReference, config.configpath)
	}

	Convey("Pygmy Application Test: "+config.name, t, func() {

		Convey("Provision environment", func() {
			Convey("Image pulled", func() {
				_, e := images.Pull("library/docker:dind")
				So(e, ShouldBeNil)
			})

			Convey("Container created", func() {
				currentWorkingDirectory, err := os.Getwd()
				So(err, ShouldBeNil)
				x, _ := containers.Create(dindContainerName, container.Config{
					Image: "docker:dind",
				}, container.HostConfig{
					AutoRemove: false,
					Binds: []string{
						fmt.Sprintf("%v%vbuilds%v:/builds", currentWorkingDirectory, string(os.PathSeparator), string(os.PathSeparator)),
						fmt.Sprintf("%v%vexamples%v:/examples", currentWorkingDirectory, string(os.PathSeparator), string(os.PathSeparator)),
					},
					Privileged: true,
				}, network.NetworkingConfig{})

				dindID = x.ID
				So(dindID, ShouldNotEqual, "")
			})

			Convey("Container started", func() {
				err := containers.Start(dindContainerName, container.StartOptions{})
				So(err, ShouldEqual, nil)
			})
		})

		Convey("Populating Daemon", func() {

			Convey("Container has started the daemon", func() {
				_, e := containers.Exec(dindContainerName, "dockerd")
				So(e, ShouldEqual, nil)
				time.Sleep(time.Second * 2)
			})

			e := containers.Start(dindContainerName, container.StartOptions{})
			if e != nil {
				fmt.Println(e)
			}

			for _, image := range config.images {
				Convey("Pulling "+image, func() {
					_, e := containers.Exec(dindContainerName, "docker pull "+image)
					time.Sleep(time.Second * 2)
					So(e, ShouldBeNil)
				})
			}
		})

		Convey("Application Tests", func() {

			Convey("Container has configuration file ("+config.configpath+")", func() {
				d, _ := containers.Exec(dindContainerName, "stat "+config.configpath)
				if config.configpath == "" {
					SkipSo(string(d), ShouldContainSubstring, config.configpath)
				} else {
					So(string(d), ShouldContainSubstring, config.configpath)
				}
			})

			Convey("Container has compiled binary from host", func() {
				d, _ := containers.Exec(dindContainerName, fmt.Sprintf("stat /builds/%v", binaryReference))
				So(string(d), ShouldContainSubstring, fmt.Sprintf("/builds/%v", binaryReference))
			})

			d, _ := containers.Exec(dindContainerName, fmt.Sprintf("/builds/%v", binaryReference))
			Convey("Container can run pygmy", func() {
				So(string(d), ShouldContainSubstring, "local containers for local development")
			})

			// While it's safe, we should clean the environment.
			_, e := containers.Exec(dindContainerName, cleanCmd)
			if e != nil {
				fmt.Println(e)
			}

			Convey("Default ports are not allocated", func() {
				g, _ := containers.Exec(dindContainerName, statusCmd)
				for _, service := range config.servicewithports {
					So(string(g), ShouldContainSubstring, service+" is able to start")
				}
			})

			Convey("Pygmy started", func() {
				d, _ = containers.Exec(dindContainerName, upCmd)
				if config.configpath != "" {
					So(string(d), ShouldContainSubstring, "Using config file: "+config.configpath)
				}
				for _, service := range config.services {
					So(string(d), ShouldContainSubstring, "Successfully started "+service)
				}
			})

			Convey("Endpoints are serving", func() {
				d, _ = containers.Exec(dindContainerName, statusCmd)
				for _, endpoint := range config.endpoints {
					if config.skipendpointchecks {
						SkipSo(string(d), ShouldNotContainSubstring, "! "+endpoint)
					} else {
						So(string(d), ShouldNotContainSubstring, "! "+endpoint)
					}
				}
			})
		})

		Convey("Environment Cleanup", func() {
			Convey("Pygmy has cleaned the environment", func() {

				_, e := containers.Exec(dindContainerName, downCmd)
				So(e, ShouldBeNil)
				_, e = containers.Exec(dindContainerName, cleanCmd)
				So(e, ShouldBeNil)
				d, _ := containers.Exec(dindContainerName, statusCmd)
				for _, service := range config.services {
					So(string(d), ShouldContainSubstring, service+" is not running")
				}
				So(e, ShouldBeNil)
			})
			// System prune container...
			Convey("Removing DinD Container", func() {
				err := containers.Kill("exampleTestContainer")
				So(err, ShouldBeNil)
				err = containers.Remove("exampleTestContainer")
				So(err, ShouldBeNil)
			})
		})
	})
}

// TestDefault will test an environment with no additional configuration.
func TestDefault(t *testing.T) {
	configuration := &config{
		name:               "default",
		configpath:         "/examples/pygmy.basic.yml",
		endpoints:          []string{"http://docker.amazee.io/stats", "http://mailhog.docker.amazee.io"},
		images:             []string{"pygmystack/haproxy", "pygmystack/dnsmasq", "pygmystack/mailhog"},
		services:           []string{"amazeeio-haproxy", "amazeeio-dnsmasq", "amazeeio-mailhog"},
		servicewithports:   []string{"amazeeio-haproxy", "amazeeio-mailhog"},
		skipendpointchecks: false,
	}
	setup(t, configuration)
}

// TestCustom will test a highly customised environment.
func TestCustom(t *testing.T) {
	configuration := &config{
		name:               "custom",
		configpath:         "/examples/pygmy.complex.yml",
		endpoints:          []string{"http://traefik.docker.amazee.io", "http://mailhog.docker.amazee.io", "http://phpmyadmin.docker.amazee.io"},
		images:             []string{"pygmystack/ssh-agent", "pygmystack/mailhog", "phpmyadmin/phpmyadmin", "library/traefik:v2.1.3"},
		services:           []string{"unofficial-traefik-2", "unofficial-phpmyadmin", "amazeeio-mailhog"},
		servicewithports:   []string{"amazeeio-mailhog", "unofficial-phpmyadmin", "unofficial-traefik-2"},
		skipendpointchecks: false,
	}
	setup(t, configuration)
}
