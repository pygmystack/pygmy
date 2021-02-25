package library

// Down will bring pygmy down safely
func Down(c Config) {

	Setup(&c)
	for _, Service := range c.Services {
		enabled, _ := Service.GetFieldBool("enable")
		if enabled {
			Service.Stop()
		}
	}
}
