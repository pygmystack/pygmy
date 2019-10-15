// +build linux

package resolv

import (
	"fmt"
	model "github.com/fubarhouse/pygmy/service/interface"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
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

	if resolv.Status() {
		model.Green(fmt.Sprintf("Successfully configured local resolver"))
	} else {
		model.Red(fmt.Sprintf("Could not configur local resolver"))
	}

	model.Green(fmt.Sprintf("setting up DNS resolution and loopback alias IP, this may require sudo"))
	ifConfig := exec.Command("sh", "-c", "ifconfig", "lo0 alias 172.16.172.16")
	err := ifConfig.Run()
	if err != nil {
		model.Red(fmt.Sprintf("error creating loopback UP alias"))
	}
	killAll := exec.Command("sh", "-c", "sudo killall mDNSResponder")
	err = killAll.Run()
	if err != nil {
		model.Red(fmt.Sprintf("error restarting mDNSResponder"))
	}
}

func (resolv Resolv) Clean() {

	fullPath := fmt.Sprintf("%v%v%v", resolv.Path, string(os.PathSeparator), resolv.File)

	if resolv.Status() {
		err := run([]string{"sudo", "rm", fullPath})
		if err == nil {
			model.Green(fmt.Sprintf("Resolver removed"))
		} else {
			model.Red(fmt.Sprintf("Error while removing the resolver"))
		}
	}

	model.Green(fmt.Sprintln("Removing resolver file and loopback alias IP, this may require sudo"))
	ifConfig := exec.Command("/bin/sh", "-c", "ifconfig", "lo0 -alias 172.16.172.16")
	err := ifConfig.Run()
	if err != nil {
		model.Red(fmt.Sprintf("error creating loopback UP alias"))
	}
	killAll := exec.Command("/bin/sh", "-c", "sudo killall mDNSResponder")
	err = killAll.Run()
	if err != nil {
		model.Green(fmt.Sprintf("error restarting mDNSResponder"))
	}

}

func (resolv Resolv) Status() bool {

	fullPath := fmt.Sprintf("%v%v%v", resolv.Path, string(os.PathSeparator), resolv.File)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		return false
	}

	return true

}
