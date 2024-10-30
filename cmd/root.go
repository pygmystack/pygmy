// Copyright © 2019 Karl Hepworth <Karl.Hepworth@gmail.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to
// deal in the Software without restriction, including without limitation the
// rights to use, copy, modify, merge, publish, distribute, sublicense, and/or
// sell copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING
// FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS
// IN THE SOFTWARE.

package cmd

import (
	"fmt"
	"github.com/pygmystack/pygmy/external/docker/setup"
	"os"
	"runtime"
	"strings"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile   string
	c         setup.Config
	validArgs = []string{"addkey", "clean", "down", "export", "pull", "restart", "status", "up", "update", "version"}
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:       "pygmy",
	ValidArgs: validArgs,
	Short:     "amazeeio's local development tool",
	Long: `amazeeio's local development tool,
	
Runs DNSMasq, HAProxy, MailHog and an SSH Agent in local containers for local development.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", findConfig(), "")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// findConfig will find the first available configuration and return a
// sensible default any of the expected paths are not found. The default
// is assigned to the default flag. If the result which is returned does
// not exist, it will not be loaded into memory and it will not be reported.
func findConfig() string {

	// Find home directory.
	home, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Find a config for non-windows.
	if runtime.GOOS != "windows" {

		// Define a list of files we need to search for.
		// The file needs to have an extension supported
		// by Viper and to be included in the strings
		// declared below.
		searchFor := []string{
			home + "/.config/pygmy/config.yaml",
			home + "/.config/pygmy/config.yml",
			home + "/.config/pygmy/pygmy.yaml",
			home + "/.config/pygmy/pygmy.yml",
			home + "/.pygmy.yaml",
			home + "/.pygmy.yml",
			"/etc/pygmy/config.yaml",
			"/etc/pygmy/config.yml",
			"/etc/pygmy/pygmy.yaml",
			"/etc/pygmy/pygmy.yml",
		}

		// Look for each of the files listed above.
		for n := range searchFor {
			if _, err := os.Stat(searchFor[n]); err == nil {
				if !os.IsNotExist(err) {
					return searchFor[n]
				}
			}
		}
	}

	// Provide a default.
	if runtime.GOOS == "linux" {
		return strings.Join([]string{"etc", "pygmy", "config.yml"}, string(os.PathSeparator))
	}
	return strings.Join([]string{home, ".pygmy.yml"}, string(os.PathSeparator))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {

	if cfgFile == "" {
		viper.SetConfigFile(findConfig())
	} else {
		viper.SetConfigFile(cfgFile)
	}

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if os.Args[1] != "completion" && !jsonOutput {
			fmt.Println("Using config file:", viper.ConfigFileUsed())
		}
	}
}
