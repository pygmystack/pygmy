// +build darwin

package resolv

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	model "github.com/fubarhouse/pygmy/service/interface"
)

func New() Resolv {
	return ResolvGeneric
}

func run(args []string) error {
	commandArgs := strings.Join(args, " ")
	command := exec.Command("sh", "-c", commandArgs)
	return command.Run()
}

func (resolv Resolv) Configure() {

	fullPath := fmt.Sprintf("%v%v%v", resolv.Path, string(os.PathSeparator), resolv.File)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// Create the directory if it doesn't exist.
		if _, err := os.Stat(resolv.Path); os.IsNotExist(err) {
			err := run([]string{"sudo", "mkdir", "-p", resolv.Path})
			if err != nil {
				fmt.Println(err)
			}
			err = run([]string{"sudo", "chmod", "777", resolv.Path})
			if err != nil {
				fmt.Println(err)
			}
		}
		// Create the file if it doesn't exist.
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			tmpFile, error := ioutil.TempFile("", "pygmy-")
			if error != nil {
				fmt.Println(error)
			}
			error = os.Chmod(tmpFile.Name(), 0777)
			if error != nil {
				fmt.Println(error)
			}
			_, error = tmpFile.WriteString(resolv.Contents)
			if error != nil {
				fmt.Println(error)
			}
			err := run([]string{"sudo", "cp", tmpFile.Name(), fullPath})
			if err != nil {
				fmt.Println(err)
			}
		}
	} else {
		// Open and write to the file if it exists.
		file, error := os.Open(fullPath)
		if error != nil {
			fmt.Println(error)
		}
		_, error = file.WriteString(resolv.Contents)
		if error != nil {
			fmt.Println(error)
		}
		error = file.Close()
		if error != nil {
			fmt.Println(error)
		}
	}

	if resolv.Status() {
		model.Green(fmt.Sprintf("Successfully configured local resolver"))
	} else {
		model.Red(fmt.Sprintf("Could not configur local resolver"))
	}

}

func (resolv Resolv) Clean() {

	//fullPath := fmt.Sprintf("%v%v%v", resolv.Path, string(os.PathSeparator), resolv.File)

}

func (resolv Resolv) Status() bool {

	fullPath := fmt.Sprintf("%v%v%v", resolv.Path, string(os.PathSeparator), resolv.File)
	if _, err := os.Stat(fullPath); os.IsExist(err) {

		bytes, error := ioutil.ReadFile(fullPath)
		if error != nil {
			fmt.Println(error)
		}

		if strings.Contains(string(bytes), resolv.Contents) {
			return true
		}

	}

	return false

}
