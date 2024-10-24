package commands

import "fmt"

// Restart will stop and start Pygmy in its entirety.
func Restart(c Config) {
	err := Down(c)
	if err != nil {
		fmt.Println(err)
	}

	err = Up(c)
	if err != nil {
		fmt.Println(err)
	}
}
