package library

import "fmt"

// Version provides business logic for the `version` command.
func Version(c Config) {
	fmt.Println("version called")
}
