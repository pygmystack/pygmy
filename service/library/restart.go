package library

// Restart provides business logic for the `restart` command.
func Restart(c Config) {
	Down(c)
	Up(c)
}
