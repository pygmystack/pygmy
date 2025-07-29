package setup

import (
	"context"
	"fmt"
	"net"
	"strings"

	"github.com/docker/docker/client"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/go-connections/nat"
	ShirouNet "github.com/shirou/gopsutil/net"
    "github.com/shirou/gopsutil/process"
)

// CompatibilityCheck is a struct of fields associated to reporting of
// a result state.
type CompatibilityCheck struct {
	State   bool   `yaml:"value"`
	Message string `yaml:"string"`
}

// DryRun will check for. It is here to check for port compatibility before
// Pygmy attempts to start any containers and provide the user with a report.
func DryRun(ctx context.Context, cli *client.Client, c *Config) ([]CompatibilityCheck, error) {

	messages := []CompatibilityCheck{}

	for _, Service := range c.Services {
		name, _ := Service.GetFieldString(ctx, cli, "name")
		enabled, _ := Service.GetFieldBool(ctx, cli, "enable")
		if enabled {
			if s, _ := Service.Status(ctx, cli); !s {
				for PortBinding, Ports := range Service.HostConfig.PortBindings {
					if strings.Contains(string(PortBinding), "tcp") {
						for _, Port := range Ports {
							p := fmt.Sprint(Port.HostPort)
							conn, err := net.Dial("tcp", "localhost:"+p)
							if conn != nil {
								if e := conn.Close(); e != nil {
									fmt.Println(e)
								}
							}
							if err != nil {
								messages = append(messages, CompatibilityCheck{
									State:   true,
									Message: fmt.Sprintf("%v is able to start on port %v", name, p),
								})
							} else {
								conn, err := net.Listen("tcp", ":"+p)
								if conn != nil {
									conn.Close()
								}
								if err != nil {
									blockingProcId, procName, err := getBlockingProcess(p, ctx, cli)
									if err == nil {
										messages = append(messages, CompatibilityCheck{
											State:   false,
											Message: fmt.Sprintf("%v is not able to start on port %v as process %d (%v) is already using this port", name, p, blockingProcId, procName),
										})
									} else {
										messages = append(messages, CompatibilityCheck{
											State:   false,
											Message: fmt.Sprintf("%v is not able to start on port %v: %v", name, p, err),
										})
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return messages, nil
}

func getBlockingProcess(rawPort string, ctx context.Context, cli *client.Client) (int, string, error) {
	p, err := nat.ParsePort(rawPort)
	if err != nil {
		return 0, "", err
	}
	port := uint32(p)

	conns, err := ShirouNet.Connections("inet")
    if err != nil {
        panic(err)
    }

    for _, conn := range conns {
        if conn.Laddr.Port == port && conn.Status == "LISTEN" {
            if conn.Pid != 0 {
                proc, err := process.NewProcess(conn.Pid)
                if err == nil {
                    name, _ := proc.Name()
					if strings.Contains(name, "docker") {
						containerName, _ := getContainerNameFromPort(port, ctx, cli)
						name = fmt.Sprintf("docker container %v", containerName)
					}
					return int(conn.Pid), name, err
                } else {
                    return 0, "", fmt.Errorf("could not get process info for PID %d\n", conn.Pid)
                }
            } else {
                return 0, "", fmt.Errorf("no PID found\n")
            }
        }
    }

	return 0, "", fmt.Errorf("no process found listening on port %d\n", port)
}

func getContainerNameFromPort(port uint32, ctx context.Context, cli *client.Client) (string, error) {
	containers, err := cli.ContainerList(ctx, container.ListOptions{})
    if err != nil {
        return "", err
    }

    for _, c := range containers {
        for _, p := range c.Ports {
            if p.PublicPort == uint16(port) {
                return c.Names[0][1:], nil 
            }
        }
    }

	return "", fmt.Errorf("no container found bound to host port %d", port)
}
