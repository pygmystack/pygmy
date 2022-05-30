package main_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/pygmystack/pygmy/service/interface/docker"
	. "github.com/smartystreets/goconvey/convey"
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
	domain             string
	prefix             string
	envs               []string
}

// setup is a configurable pipeline which allows different configurations to
// run to keep the consistency for as many tests as are required.
func setup(t *testing.T, config *config) {

	cmdPrefix := ""
	for _, v := range config.envs {
		cmdPrefix += fmt.Sprintf(" %v", v)
	}

	var cleanCmd = fmt.Sprintf("/builds/%v clean", binaryReference)
	var statusCmd = fmt.Sprintf("/builds/%v status", binaryReference)
	var upCmd = fmt.Sprintf("%v /builds/%v up", cmdPrefix, binaryReference)
	var downCmd = fmt.Sprintf("/builds/%v down", binaryReference)

	if config.configpath != "" {
		cleanCmd = fmt.Sprintf("/builds/%v clean --config %v", binaryReference, config.configpath)
		statusCmd = fmt.Sprintf("/builds/%v status --config %v", binaryReference, config.configpath)
		upCmd = fmt.Sprintf("/builds/%v up --config %v", binaryReference, config.configpath)
		downCmd = fmt.Sprintf("/builds/%v down --config %v", binaryReference, config.configpath)
	}

	time.Sleep(5)

	Convey("Pygmy Application Test: "+config.name, t, func() {

		Convey("Provision environment", func() {
			Convey("Image pulled", func() {
				_, e := docker.DockerPull("library/docker:dind")
				So(e, ShouldBeNil)
			})

			Convey("Container created", func() {
				currentWorkingDirectory, err := os.Getwd()
				So(err, ShouldBeNil)
				x, _ := docker.DockerContainerCreate(dindContainerName, container.Config{
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
				err := docker.DockerContainerStart(dindContainerName, types.ContainerStartOptions{})
				So(err, ShouldEqual, nil)
			})
		})

		Convey("Populating Daemon", func() {

			Convey("Container has started the daemon", func() {
				_, e := docker.DockerExec(dindContainerName, "dockerd")
				So(e, ShouldEqual, nil)
				time.Sleep(time.Second * 2)
			})

			e := docker.DockerContainerStart(dindContainerName, types.ContainerStartOptions{})
			if e != nil {
				fmt.Println(e)
			}

			for _, image := range config.images {
				Convey("Pulling "+image, func() {
					_, e := docker.DockerExec(dindContainerName, "docker pull "+image)
					time.Sleep(time.Second * 2)
					So(e, ShouldBeNil)
				})
			}
		})

		Convey("Application Tests", func() {

			Convey("Container has configuration file ("+config.configpath+")", func() {
				d, _ := docker.DockerExec(dindContainerName, "stat "+config.configpath)
				if config.configpath == "" {
					SkipSo(string(d), ShouldContainSubstring, config.configpath)
				} else {
					So(string(d), ShouldContainSubstring, config.configpath)
				}
			})

			Convey("Container has compiled binary from host", func() {
				d, _ := docker.DockerExec(dindContainerName, fmt.Sprintf("stat /builds/%v", binaryReference))
				So(string(d), ShouldContainSubstring, fmt.Sprintf("/builds/%v", binaryReference))
			})

			d, _ := docker.DockerExec(dindContainerName, fmt.Sprintf("/builds/%v", binaryReference))
			Convey("Container can run pygmy", func() {
				So(string(d), ShouldContainSubstring, "local containers for local development")
			})

			// While it's safe, we should clean the environment.
			_, e := docker.DockerExec(dindContainerName, cleanCmd)
			if e != nil {
				fmt.Println(e)
			}

			Convey("Default ports are not allocated", func() {
				g, _ := docker.DockerExec(dindContainerName, statusCmd)
				for _, service := range config.servicewithports {
					So(string(g), ShouldContainSubstring, service+" is able to start")
				}
			})

			Convey("Pygmy started", func() {
				d, _ = docker.DockerExec(dindContainerName, upCmd)
				if config.configpath != "" {
					So(string(d), ShouldContainSubstring, "Using config file: "+config.configpath)
				}
				for _, service := range config.services {
					So(string(d), ShouldContainSubstring, "Successfully started "+service)
				}
			})

			Convey("Endpoints are serving", func() {
				d, _ = docker.DockerExec(dindContainerName, statusCmd)
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

				_, e := docker.DockerExec(dindContainerName, downCmd)
				So(e, ShouldBeNil)
				_, e = docker.DockerExec(dindContainerName, cleanCmd)
				So(e, ShouldBeNil)
				d, _ := docker.DockerExec(dindContainerName, statusCmd)
				for _, service := range config.services {
					So(string(d), ShouldContainSubstring, service+" is not running")
				}
				So(e, ShouldBeNil)
			})
			// System prune container...
			Convey("Removing DinD Container", func() {
				err := docker.DockerKill("exampleTestContainer")
				So(err, ShouldBeNil)
				err = docker.DockerRemove("exampleTestContainer")
				So(err, ShouldBeNil)
			})
		})
	})
}

// TestDefault will test an environment with no additional configuration.
func TestDefault(t *testing.T) {
	domain := "docker.amazee.io"
	prefix := "custom"
	configuration := &config{
		name:               "default",
		configpath:         "/examples/pygmy.basic.yml",
		domain:             domain,
		envs:               []string{"PYGMY_DOMAIN=docker.amazee.io"},
		prefix:             prefix,
		endpoints:          []string{fmt.Sprintf("http://%s/stats", domain), fmt.Sprintf("http://mailhog.%s/", domain)},
		images:             []string{"pygmystack/haproxy", "pygmystack/dnsmasq", "pygmystack/mailhog"},
		services:           []string{fmt.Sprintf("%s-haproxy", prefix), fmt.Sprintf("%s-dnsmasq", prefix), fmt.Sprintf("%s-mailhog", prefix)},
		servicewithports:   []string{fmt.Sprintf("%s-haproxy", prefix), fmt.Sprintf("%s-mailhog", prefix)},
		skipendpointchecks: false,
	}
	setup(t, configuration)
}

// TestCustom will test a highly customised environment.
func TestCustom(t *testing.T) {
	domain := "pygmy.site"
	prefix := "unofficial"
	configuration := &config{
		name:               "custom",
		configpath:         "/examples/pygmy.complex.yml",
		domain:             domain,
		envs:               []string{"PYGMY_DOMAIN=docker.amazee.io"},
		prefix:             prefix,
		endpoints:          []string{fmt.Sprintf("http://traefik.%s", domain), fmt.Sprintf("http://mailhog.%s", domain), fmt.Sprintf("http://portainer.%s", domain), "http://phpmyadmin.pygmy.site"},
		images:             []string{"pygmystack/ssh-agent", "pygmystack/mailhog", "phpmyadmin/phpmyadmin", "portainer/portainer", "library/traefik:v2.1.3"},
		services:           []string{fmt.Sprintf("%s-portainer", prefix), fmt.Sprintf("%s-traefik-2", prefix), fmt.Sprintf("%s-phpmyadmin", prefix), fmt.Sprintf("%s-mailhog", prefix)},
		servicewithports:   []string{fmt.Sprintf("%s-mailhog", prefix), fmt.Sprintf("%s-portainer", prefix), fmt.Sprintf("%s-phpmyadmin", prefix), fmt.Sprintf("%s-traefik-2", prefix)},
		skipendpointchecks: false,
	}
	setup(t, configuration)
}
