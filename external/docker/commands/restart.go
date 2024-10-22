package commands

// Restart will stop and start Pygmy in its entirety.
func Restart(c Config) {
	Down(c)
	Up(c)
}
