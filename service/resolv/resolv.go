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

// run will run a shell command and is not exported.
// Shell functionality is exclusive to this package.
func run(args []string) error {
	commandArgs := strings.Join(args, " ")
	command := exec.Command("sh", "-c", commandArgs)
	return command.Run()
}

// Configure will ensure the given Resolv type a method that can setup a file
// with the contents of Data at File in Folder. This file will route traffic
// on a configured namespace to the localhost and dnsmasq will accept this
// traffic and route it to the docker container. It will remove the file and/or
// rewrite the contents for both MacOS and Linux - Linux will however result in
// removing the string from the file, where MacOS will contain a file with only
// the contents of Data. MacOS will also run the following upon completion of
// this function:
// * sudo ifconfig lo0 alias 172.16.172.16
// * sudo killall mDNSResponder
func (resolv Resolv) Configure() {

	var cmdOut []byte
	var tmpFile *os.File

	if !resolv.Enabled {
		return
	}
	if resolv.Status() {
		fmt.Printf("Already configured resolvr %v\n", resolv.Name)
	} else {
		fullPath := fmt.Sprintf("%v%v%v", resolv.Folder, string(os.PathSeparator), resolv.File)
		if _, err := os.Stat(fullPath); os.IsNotExist(err) {
			// Create the directory if it doesn't exist.
			if _, err := os.Stat(resolv.Folder); os.IsNotExist(err) {
				if err := run([]string{"sudo", "mkdir", "-p", resolv.Folder}); err != nil {
					fmt.Println(err)
				}
				if err = run([]string{"sudo", "chmod", "777", resolv.Folder}); err != nil {
					fmt.Println(err)
				}
			}

			// Create the file if it doesn't exist.
			if _, err := os.Stat(fullPath); os.IsNotExist(err) {
				if tmpFile, err = ioutil.TempFile("", "pygmy-"); err != nil {
					fmt.Println(err)
				}
				if err = os.Chmod(tmpFile.Name(), 0777); err != nil {
					fmt.Println(err)
				}
				if _, err = tmpFile.WriteString(resolv.Data); err != nil {
					fmt.Println(err)
				}
				if err = run([]string{"sudo", "cp", tmpFile.Name(), fullPath}); err != nil {
					fmt.Println(err)
				}
			}
		} else {

			// If the bytes haven't already been written to the file:
			if !resolv.statusFileData() {

				if _, err := os.Stat(fullPath); err == nil {

					cmd := exec.Command("/bin/sh", "-c", "cat "+fullPath)

					if cmdOut, err = cmd.Output(); err != nil {
						fmt.Println(err.Error())
						fmt.Println("/bin/sh", "-c", "cat "+fullPath)
					}

					if tmpFile, err = ioutil.TempFile("", "pygmy-"); err != nil {
						fmt.Println(err)
					} else {
						if err = os.Chmod(tmpFile.Name(), 0777); err != nil {
							fmt.Println(err)
						}
						if _, err = tmpFile.WriteString(string(cmdOut)); err != nil {
							fmt.Println(err)
						}
						if _, err = tmpFile.WriteString(resolv.Data); err != nil {
							fmt.Println(err)
						}
						if err = tmpFile.Close(); err != nil {
							fmt.Println(err)
						}
						if err = run([]string{"sudo", "cp", tmpFile.Name(), fullPath}); err != nil {
							fmt.Println(err)
						}
					}
				}
			}
		}

		ifConfig := exec.Command("/bin/sh", "-c", "sudo ifconfig lo0 alias 172.16.172.16")
		if err := ifConfig.Run(); err != nil {
			fmt.Println("error creating loopback UP alias")
		}
		killAll := exec.Command("/bin/sh", "-c", "sudo killall mDNSResponder")
		if err := killAll.Run(); err != nil {
			fmt.Println("error restarting mDNSResponder")
		}

		if resolv.Status() {
			fmt.Printf("Successfully configured resolvr %v\n", resolv.Name)
		}
	}
}

// Clean will cleanup the resolv file configured to the system and run some
// cleanup commands which were ran at the end of resolv.Configure on MacOS.
func (resolv Resolv) Clean() {

	fullPath := fmt.Sprintf("%v%v%v", resolv.Folder, string(os.PathSeparator), resolv.File)
	if runtime.GOOS == "linux" {
		if _, err := os.Stat(fullPath); err == nil {

			cmd := exec.Command("/bin/sh", "-c", "cat "+fullPath)
			if cmdOut, cmdErr := cmd.Output(); cmdErr != nil {
				fmt.Println(cmdErr.Error())
			} else {
				if strings.Contains(string(cmdOut), resolv.Data) {
					newFile := strings.Replace(string(cmdOut), resolv.Data, "", -1)
					if tmpFile, err := ioutil.TempFile("", "pygmy-"); err != nil {
						fmt.Println(err)
					} else {
						if err = os.Chmod(tmpFile.Name(), 0777); err != nil {
							fmt.Println(err)
						}
						if _, err = tmpFile.WriteString(newFile); err != nil {
							fmt.Println(err)
						}
						if err = tmpFile.Close(); err != nil {
							fmt.Println(err)
						}
						if err = run([]string{"sudo", "cp", tmpFile.Name(), fullPath}); err != nil {
							fmt.Println(err)
						}
					}
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

// Status is an exported state function which will check the file contents
// matches Data on Linux, or return the result of three independent checks
// on MacOS including the file, network and data checks.
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

// statusFileData will check the resolv file contents matches what is expected
func (resolv Resolv) statusFileData() bool {
	fullPath := fmt.Sprintf("%v%v%v", resolv.Folder, string(os.PathSeparator), resolv.File)
	cmd := exec.Command("/bin/sh", "-c", "cat "+fullPath)
	if cmdOut, cmdErr := cmd.Output(); cmdErr != nil {
		fmt.Println(cmdErr.Error())
	} else {
		return strings.Contains(string(cmdOut), resolv.Data)
	}
	return false
}

// statusFile will check the expected file exists
func (resolv Resolv) statusFile() bool {
	fullPath := fmt.Sprintf("%v%v%v", resolv.Folder, string(os.PathSeparator), resolv.File)
	if _, err := os.Stat(fullPath); !os.IsExist(err) {
		return true
	}
	return false
}

// statusNet will check the network has the required config.
// One of the original Pygmy's more regular issues was that the network had  no
// checks, so the command to make that change was ran as much as logic provided
// and as a result there were some very unusual and unfixable issues.
// This has completely ruled that situation out.
func (resolv Resolv) statusNet() bool {
	ifConfigCmd := exec.Command("/bin/sh", "-c", "ifconfig lo0")
	if out, ifConfigErr := ifConfigCmd.Output(); ifConfigErr != nil {
		fmt.Println(ifConfigErr.Error())
		return false
	} else {
		return strings.Contains(string(out), "172.16.172.16")
	}
}
