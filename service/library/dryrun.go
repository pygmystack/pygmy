package library

import (
	"fmt"
	"net"
	"strings"
)

type CompatibilityCheck struct {
	State   bool   `yaml:"value"`
	Message string `yaml:"string"`
}

func DryRun(c *Config) []CompatibilityCheck {

	messages := []CompatibilityCheck{}

	for _, Service := range c.Services {
		if !Service.Disabled {
			if s, _ := Service.Status(); !s {
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
									Message: fmt.Sprintf("[*] %v is able to start on port %v", Service.Name, p),
								})
							} else {
								conn, err := net.Listen("tcp", ":"+p)
								if conn != nil {
									if e := conn.Close(); e != nil {
										fmt.Println(e)
									}
								}
								if err != nil {
									messages = append(messages, CompatibilityCheck{
										State:   false,
										Message: fmt.Sprintf("[ ] %v is not able to start on port %v: %v", Service.Name, p, err),
									})
								} else {
									messages = append(messages, CompatibilityCheck{
										State:   true,
										Message: fmt.Sprintf("[*] %v is able to start on port %v", Service.Name, p),
									})
								}
							}
						}
					}
				}
			}
		}
	}

	return messages
}
