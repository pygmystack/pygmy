package resolv

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

type Resolv struct {
	Name string `yaml:"name"`
	Path string `yaml:"path"`
	Data string `yaml:"contents"`
	File string `yaml:"file"`
}

type resolv interface {
	Clean()
	Configure()
	New() Resolv
	Status() bool
}

func New(d Resolv) Resolv {
	p := strings.Split(d.Path, string(os.PathSeparator))
	filename := p[len(p)-1]
	return Resolv{
		Name: d.Name,
		File: filename,
		Data: d.Data,
		Path: strings.Replace(d.Path, string(os.PathSeparator)+filename, "", -1),
	}
}

func run(args []string) error {
	commandArgs := strings.Join(args, " ")
	command := exec.Command("sh", "-c", commandArgs)
	return command.Run()
}

func (resolv Resolv) Configure() {

	if resolv.Status() {
		fmt.Printf("Already configured resolvr %v\n", resolv.Name)
	} else {
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
				_, error = tmpFile.WriteString(resolv.Data)
				if error != nil {
					fmt.Println(error)
				}
				err := run([]string{"sudo", "cp", tmpFile.Name(), fullPath})
				if err != nil {
					fmt.Println(err)
				}
			}
		}
		if _, err := os.Stat(fullPath); err == nil {

			cmd := exec.Command("/bin/sh", "-c", "cat "+fullPath)
			cmdOut, cmdErr := cmd.Output()
			if cmdErr != nil {
				fmt.Println(cmdErr.Error())
			}

			tmpFile, error := ioutil.TempFile("", "pygmy-")
			if error != nil {
				fmt.Println(error)
			}
			error = os.Chmod(tmpFile.Name(), 0777)
			if error != nil {
				fmt.Println(error)
			}
			_, error = tmpFile.WriteString(string(cmdOut))
			_, error = tmpFile.WriteString(resolv.Data)
			if error != nil {
				fmt.Println(error)
			}
			err := run([]string{"sudo", "cp", tmpFile.Name(), fullPath})
			if err != nil {
				fmt.Println(err)
			}
		}

		if runtime.GOOS == "darwin" {
			ifConfig := exec.Command("/bin/sh", "-c", "sudo ifconfig lo0 alias 172.16.172.16")
			err := ifConfig.Run()
			if err != nil {
				fmt.Println("error creating loopback UP alias")
			}
			killAll := exec.Command("/bin/sh", "-c", "sudo killall mDNSResponder")
			err = killAll.Run()
			if err != nil {
				fmt.Println("error restarting mDNSResponder")
			}
		}

		if resolv.Status() {
			fmt.Printf("Successfully configured resolvr %v\n", resolv.Name)
		}
	}

}

func (resolv Resolv) Clean() {

	fullPath := fmt.Sprintf("%v%v%v", resolv.Path, string(os.PathSeparator), resolv.File)
	if _, err := os.Stat(fullPath); err == nil {

		cmd := exec.Command("/bin/sh", "-c", "cat "+fullPath)
		cmdOut, cmdErr := cmd.Output()
		if cmdErr != nil {
			fmt.Println(cmdErr.Error())
		}
		if strings.Contains(string(cmdOut), resolv.Data) {
			newFile := strings.Replace(string(cmdOut), resolv.Data, "", -1)
			tmpFile, error := ioutil.TempFile("", "pygmy-")
			if error != nil {
				fmt.Println(error)
			}
			error = os.Chmod(tmpFile.Name(), 0777)
			if error != nil {
				fmt.Println(error)
			}
			_, error = tmpFile.WriteString(newFile)
			if error != nil {
				fmt.Println(error)
			}
			err := run([]string{"sudo", "cp", tmpFile.Name(), fullPath})
			if err != nil {
				fmt.Println(err)
			}
		}

		if runtime.GOOS == "darwin" {

			if fullPath == "/etc/resolver/docker.amazee.io" {
				if err = os.Remove(fullPath); err != nil {
					fmt.Println(err)
				} else {
					fmt.Println("Successfully removed resolver file")
				}
			}

			fmt.Println("Removing loopback alias IP (may require sudo)")
			ifConfig := exec.Command("/bin/sh", "-c", "sudo ifconfig lo0 -alias 172.16.172.16")
			err = ifConfig.Run()
			if err != nil {
				fmt.Println("error removing loopback UP alias", err)
			} else {
				if !resolv.statusNet() {
					fmt.Println("Successfully removed loopback alias IP.")
				}
			}

			killAll := exec.Command("/bin/sh", "-c", "sudo killall mDNSResponder")
			err = killAll.Run()
			if err != nil {
				fmt.Println("error restarting mDNSResponder")
			} else {
				fmt.Println("Successfully restarted mDNSResponder")
			}
		}
	}

}

func (resolv Resolv) Status() bool {

	if runtime.GOOS == "darwin" {
		return resolv.statusFile() && resolv.statusNet()
	}
	fullPath := fmt.Sprintf("%v%v%v", resolv.Path, string(os.PathSeparator), resolv.File)
	if _, err := os.Stat(fullPath); err == nil {

		cmd := exec.Command("/bin/sh", "-c", "cat "+fullPath)
		cmdOut, cmdErr := cmd.Output()
		if cmdErr != nil {
			fmt.Println(cmdErr.Error())
		}
		if strings.Contains(string(cmdOut), resolv.Data) {
			return true
		}
	}

	return false

}

func (resolv Resolv) statusFile() bool {
	fullPath := fmt.Sprintf("%v%v%v", resolv.Path, string(os.PathSeparator), resolv.File)
	if _, err := os.Stat(fullPath); !os.IsExist(err) {
		return true
	}
	return false
}

func (resolv Resolv) statusNet() bool {
	ifConfigCmd := exec.Command("/bin/sh", "-c", "ifconfig lo0")
	out, ifConfigErr := ifConfigCmd.Output()
	if ifConfigErr != nil {
		fmt.Println(ifConfigErr.Error())
		return false
	}
	return strings.Contains(string(out), "172.16.172.16")
}
