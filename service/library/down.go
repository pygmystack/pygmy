package library

// Down will bring pygmy down safely
func Down(c Config) {

	Setup(&c)

	for _, Service := range c.Services {
		disabled, _ := Service.GetFieldBool("disabled")
		if !disabled {
			Service.Stop()
		}
	}

	for _, resolver := range c.Resolvers {
		resolver.Clean()
	}
}
