// +build windows

package resolv

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"sync"
)

func New(resolv Resolv) Resolv {
	return resolv
}

func runCommand(args []string) ([]byte, error) {

	powershell, err := exec.LookPath("powershell")
	if err != nil {
		fmt.Println(err)
	}

	// Generate the command, based on input.
	cmd := exec.Cmd{}
	cmd.Path = powershell
	cmd.Args = []string{powershell}

	// Add our arguments to the command.
	cmd.Args = append(cmd.Args, args...)

	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output

	// Check the errors, return as needed.
	var wg sync.WaitGroup
	wg.Add(1)
	err = cmd.Run()

	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}
	wg.Done()

	return output.Bytes(), nil

}

func (resolv Resolv) Clean() {
	_, error := runCommand([]string{"Clear-ItemProperty -Path HKLM:\\SYSTEM\\CurrentControlSet\\Services\\Tcpip\\Parameters -Name Domain"})
	if error != nil {
		fmt.Println(error.Error())
	}
}
func (resolv Resolv) Configure() {
	if !resolv.Disabled {
		_, error := runCommand([]string{"Set-ItemProperty -Path HKLM:\\SYSTEM\\CurrentControlSet\\Services\\Tcpip\\Parameters -Name Domain -Value docker.amazee.io"})
		if error != nil {
			fmt.Println(error.Error())
		}
	}
}

func (resolv Resolv) Status() bool {
	data, error := runCommand([]string{"Get-ItemProperty -Path HKLM:\\SYSTEM\\CurrentControlSet\\Services\\Tcpip\\Parameters"})
	if error != nil {
		return false
	}
	for _, v := range strings.Split(string(data), "\n") {
		if strings.HasPrefix(v, "Domain") && strings.Contains(v, "docker.amazee.io") {
			return true
		}
	}
	return false
}
