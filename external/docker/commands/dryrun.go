package commands

import (
	"context"
	"fmt"
	"github.com/docker/docker/client"
	"net"
	"strings"
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

	return messages, nil
}
