package library

import "fmt"

// Down will bring pygmy down safely
func Down(c Config) {

	Setup(&c)
	for _, Service := range c.Services {
		enabled, _ := Service.GetFieldBool("enable")
		if enabled {
			e := Service.Stop()
			if e != nil {
				name, _ := Service.GetFieldString("name")
				fmt.Printf("Failed to stop %s\n", name)
			}
		}
	}
}
