package library

import (
	"fmt"
)

// Stop will bring pygmy down safely
func Stop(c Config) {

	Setup(&c)
	NetworksToClean := []string{}

	for _, Service := range c.Services {
		enabled, _ := Service.GetFieldBool("enable")
		if enabled {
			e := Service.Stop()
			if e != nil {
				fmt.Println(e)
			}
			if s, _ := Service.GetFieldString("network"); s != "" {
				NetworksToClean = append(NetworksToClean, s)
			}
		}
	}
}
