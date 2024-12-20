package commands

import (
	"fmt"
	"os"

	"github.com/ghodss/yaml"

	"github.com/pygmystack/pygmy/external/docker/setup"
	"github.com/pygmystack/pygmy/internal/runtime/docker/internals"
)

// Export will export validated configuration to a given path, or it will
// export by default to $HOME/.pygmy.yml
func Export(c setup.Config, output string) error {
	cli, ctx, err := internals.NewClient()
	if err != nil {
		return err
	}

	// Set up the configuration.
	setup.Setup(ctx, cli, &c)

	// Marshal to Yaml.
	x, err := yaml.Marshal(c)
	if err != nil {
		fmt.Println(err)
		return err
	}

	// Provide output for state.
	fmt.Println("Data has been marshalled to YAML")

	// Does the file exist?
	if _, e := os.Stat(output); !os.IsNotExist(e) {
		// Remove the existing file.
		if err := os.Remove(output); err != nil {
			fmt.Println(err)
			return err
		}

		// Provide output for state.
		fmt.Printf("Path %v has been removed\n", output)

	}

	if _, e := os.Stat(output); os.IsNotExist(e) {

		// Create the new file.
		file, err := os.Create(output)
		if err != nil {
			fmt.Println(err)
			return err
		}

		// Provide output for state.
		fmt.Printf("Path %v has been created\n", output)

		// Housekeeping.
		defer file.Close()

		_, err = file.WriteString(string(x))
		if err != nil {
			fmt.Println(err)
			return err
		}

		// Provide output for state.
		fmt.Printf("Data has been written to file %v\n", output)

		err = file.Sync()
		if err != nil {
			fmt.Println(err)
			return err
		}

		// Provide output for state.
		fmt.Println("Operation completed successfully")
	}

	return nil
}
