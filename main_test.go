package main

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/client"
	model "github.com/fubarhouse/pygmy-go/service/interface"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	dindContainerName = "exampleTestContainer"
)

var (
	dindID  string
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

	var cleanCmd = "/builds/pygmy-go-linux-x86 clean"
	var statusCmd = "/builds/pygmy-go-linux-x86 status"
	var upCmd = "/builds/pygmy-go-linux-x86 up"
	var downCmd = "/builds/pygmy-go-linux-x86 down"

	if config.configpath != "" {
		cleanCmd = fmt.Sprintf("/builds/pygmy-go-linux-x86 clean --config %v", config.configpath)
		statusCmd = fmt.Sprintf("/builds/pygmy-go-linux-x86 status --config %v", config.configpath)
		upCmd = fmt.Sprintf("/builds/pygmy-go-linux-x86 up --config %v", config.configpath)
		downCmd = fmt.Sprintf("/builds/pygmy-go-linux-x86 dow --config %v", config.configpath)
	}

	Convey("Pygmy Application Test: "+config.name, t, func() {

		ctx := context.Background()
		cli, err := client.NewEnvClient()

		Convey("Provision environment", func() {
			Convey("Connection to Docker Client", func() {
				So(err, ShouldEqual, nil)
			})

			Convey("Image pulled", func() {
				_, e := model.DockerPull("library/docker:dind")
				So(e, ShouldBeNil)
			})

			Convey("Container created", func() {
				currentWorkingDirectory, err := os.Getwd()
				So(err, ShouldBeNil)
				x, _ := cli.ContainerCreate(ctx, &container.Config{
					Image: "docker:dind",
				}, &container.HostConfig{
					AutoRemove: false,
					Binds: []string{
						fmt.Sprintf("%v%vbuilds%v:/builds", currentWorkingDirectory, string(os.PathSeparator), string(os.PathSeparator)),
						fmt.Sprintf("%v%vexamples%v:/examples", currentWorkingDirectory, string(os.PathSeparator), string(os.PathSeparator)),
					},
					Privileged: true,
				}, &network.NetworkingConfig{}, dindContainerName)

				dindID = x.ID
				So(dindID, ShouldNotEqual, "")
			})

			Convey("Container started", func() {
				err = cli.ContainerStart(ctx, dindContainerName, types.ContainerStartOptions{})
				So(err, ShouldEqual, nil)
			})
		})

		Convey("Populating Daemon", func() {

			Convey("Container has started the daemon", func() {
				_, e := model.DockerExec(dindContainerName, "dockerd")
				So(e, ShouldEqual, nil)
				time.Sleep(time.Second * 2)
			})

			e := cli.ContainerStart(ctx, dindContainerName, types.ContainerStartOptions{})
			if e != nil {
				fmt.Println(e)
			}

			for _, image := range config.images {
				Convey("Pulling "+image, func() {
					_, e := model.DockerExec(dindContainerName, "docker pull "+image)
					time.Sleep(time.Second * 2)
					So(e, ShouldBeNil)
				})
			}
		})

		Convey("Application Tests", func() {

			Convey("Container has configuration file ("+config.configpath+")", func() {
				d, _ := model.DockerExec(dindContainerName, "stat "+config.configpath)
				if config.configpath == "" {
					SkipSo(string(d), ShouldContainSubstring, config.configpath)
				} else {
					So(string(d), ShouldContainSubstring, config.configpath)
				}
			})

			Convey("Container has compiled binary from host", func() {
				d, _ := model.DockerExec(dindContainerName, "stat /builds/pygmy-go-linux-x86")
				So(string(d), ShouldContainSubstring, "/builds/pygmy-go-linux-x86")
			})

			d, _ := model.DockerExec(dindContainerName, "/builds/pygmy-go-linux-x86")
			Convey("Container can run pygmy", func() {
				So(string(d), ShouldContainSubstring, "local containers for local development")
			})

			// While it's safe, we should clean the environment.
			_, e := model.DockerExec(dindContainerName, cleanCmd)
			if e != nil {
				fmt.Println(e)
			}

			Convey("Default ports are not allocated", func() {
				g, _ := model.DockerExec(dindContainerName, statusCmd)
				for _, service := range config.servicewithports {
					So(string(g), ShouldContainSubstring, service+" is able to start")
				}
			})

			Convey("Pygmy started", func() {
				d, _ = model.DockerExec(dindContainerName, upCmd)
				if config.configpath != "" {
					So(string(d), ShouldContainSubstring, "Using config file: "+config.configpath)
				}
				for _, service := range config.services {
					So(string(d), ShouldContainSubstring, "Successfully started "+service)
				}
			})

			Convey("Endpoints are serving", func() {
				d, _ = model.DockerExec(dindContainerName, statusCmd)
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

				_, e := model.DockerExec(dindContainerName, downCmd)
				So(e, ShouldBeNil)
				_, e = model.DockerExec(dindContainerName, cleanCmd)
				So(e, ShouldBeNil)
				d, _ := model.DockerExec(dindContainerName, statusCmd)
				for _, service := range config.services {
					So(string(d), ShouldContainSubstring, service+" is not running")
				}
				So(e, ShouldBeNil)
			})
			// System prune container...
			Convey("Removing DinD Container", func() {
				ctx := context.Background()
				cli, err := client.NewEnvClient()
				So(err, ShouldBeNil)
				err = cli.ContainerKill(ctx, "exampleTestContainer", "")
				So(err, ShouldBeNil)
				err = cli.ContainerRemove(ctx, "exampleTestContainer", types.ContainerRemoveOptions{Force: true})
				So(err, ShouldBeNil)
			})
		})
	})
}

// TestDefault will test an environment with no additional configuration.
func TestDefault(t *testing.T) {
	configuration := &config{
		name:               "default",
		configpath:         "/examples/pygmy.travis.yml",
		endpoints:          []string{"http://docker.amazee.io/stats", "http://mailhog.docker.amazee.io"},
		images:             []string{"amazeeio/haproxy", "andyshinn/dnsmasq:2.78", "mailhog/mailhog"},
		services:           []string{"amazeeio-haproxy", "amazeeio-dnsmasq", "amazeeio-mailhog"},
		servicewithports:   []string{"amazeeio-haproxy", "amazeeio-mailhog"},
		skipendpointchecks: true,
	}
	setup(t, configuration)
}

// TestCustom will test a highly customised environment.
func TestCustom(t *testing.T) {
	configuration := &config{
		name:               "custom",
		configpath:         "/examples/pygmy.complex.yml",
		endpoints:          []string{"http://traefik.docker.amazee.io", "http://mailhog.docker.amazee.io", "http://portainer.docker.amazee.io", "http://phpmyadmin.docker.amazee.io"},
		images:             []string{"amazeeio/ssh-agent", "mailhog/mailhog", "phpmyadmin/phpmyadmin", "portainer/portainer", "library/traefik:v2.1.3"},
		services:           []string{"unofficial-portainer", "unofficial-traefik-2", "unofficial-phpmyadmin", "amazeeio-mailhog"},
		servicewithports:   []string{"amazeeio-mailhog", "unofficial-portainer", "unofficial-phpmyadmin", "unofficial-traefik-2"},
		skipendpointchecks: false,
	}
	setup(t, configuration)
}
