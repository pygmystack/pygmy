package library

func Restart(c Config) {
	Down(c)
	Up(c)
}
