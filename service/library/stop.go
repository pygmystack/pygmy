package library

import (
	"fmt"
)

// Stop will bring pygmy down safely
func Stop(c Config) {

	Setup(&c)

	for _, Service := range c.Services {
		enabled, _ := Service.GetFieldBool("enable")
		if enabled {
			e := Service.Stop()
			if e != nil {
				fmt.Println(e)
			}
		}
	}
}
