// +build !windows

package resolv

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func New(resolv Resolv) Resolv {
	return resolv
}

func run(args []string) error {
	commandArgs := strings.Join(args, " ")
	command := exec.Command("sh", "-c", commandArgs)
	return command.Run()
}

func (resolv Resolv) Configure() {

	if resolv.Disabled {
		return
	}
	if resolv.Status() {
		fmt.Printf("Already configured resolvr %v\n", resolv.Name)
	} else {
		fullPath := fmt.Sprintf("%v%v%v", resolv.Folder, string(os.PathSeparator), resolv.File)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			// Create the directory if it doesn't exist.
			if _, err := os.Stat(resolv.Folder); os.IsNotExist(err) {
				err := run([]string{"sudo", "mkdir", "-p", resolv.Folder})
				if err != nil {
					fmt.Println(err)
				}
				err = run([]string{"sudo", "chmod", "777", resolv.Folder})
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
		} else {

			// If the bytes haven't already been written to the file:
			if !resolv.statusFileData() {

				if _, err := os.Stat(fullPath); err == nil {

					cmd := exec.Command("/bin/sh", "-c", "cat "+fullPath)
					cmdOut, cmdErr := cmd.Output()
					if cmdErr != nil {
						fmt.Println(cmdErr.Error())
						fmt.Println("/bin/sh", "-c", "cat "+fullPath)
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

	fullPath := fmt.Sprintf("%v%v%v", resolv.Folder, string(os.PathSeparator), resolv.File)
	if runtime.GOOS == "linux" {
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
		}
	}

	if runtime.GOOS == "darwin" {

		if strings.HasPrefix(fullPath, "/etc/resolver/") {
			if _, err := os.Stat(fullPath); err == nil {
				err := run([]string{"sudo", "rm", fullPath})
				if err != nil {
					fmt.Println(err)
				}
				if !resolv.statusFile() {
					fmt.Println("Successfully removed resolver file")
				}
			}
		}
	}

	if runtime.GOOS == "darwin" {

		if resolv.statusNet() {
			fmt.Println("Removing loopback alias IP (may require sudo)")
			ifConfig := exec.Command("/bin/sh", "-c", "sudo ifconfig lo0 -alias 172.16.172.16")
			err := ifConfig.Run()
			if err != nil {
				fmt.Println("error removing loopback UP alias", err)
			} else {
				if !resolv.statusNet() {
					fmt.Println("Successfully removed loopback alias IP.")
				}
			}
		}

		killAll := exec.Command("/bin/sh", "-c", "sudo killall mDNSResponder")
		err := killAll.Run()
		if err != nil {
			fmt.Println("error restarting mDNSResponder")
		} else {
			fmt.Println("Successfully restarted mDNSResponder")
		}
	}

}

func (resolv Resolv) Status() bool {

	if runtime.GOOS == "darwin" {
		return resolv.statusFile() && resolv.statusNet() && resolv.statusFileData()
	}
	fullPath := fmt.Sprintf("%v%v%v", resolv.Folder, string(os.PathSeparator), resolv.File)
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

func (resolv Resolv) statusFileData() bool {
	fullPath := fmt.Sprintf("%v%v%v", resolv.Folder, string(os.PathSeparator), resolv.File)
	cmd := exec.Command("/bin/sh", "-c", "cat "+fullPath)
	cmdOut, cmdErr := cmd.Output()
	if cmdErr != nil {
		fmt.Println(cmdErr.Error())
	}
	return strings.Contains(string(cmdOut), resolv.Data)
}

func (resolv Resolv) statusFile() bool {
	fullPath := fmt.Sprintf("%v%v%v", resolv.Folder, string(os.PathSeparator), resolv.File)
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
