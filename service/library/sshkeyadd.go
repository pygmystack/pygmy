package library

import (
	"fmt"
	"os"

	"github.com/fubarhouse/pygmy/service/ssh/agent"
)

func SshKeyAdd(c Config, key string) {

	if c.SkipKey {
		return
	}

	Setup(&c)

	if key != "" {
		if _, err := os.Stat(key); err != nil {
			fmt.Printf("The file path %v does not exist, or is not readable.\n%v\n", key, err)
			return
		}
	}

	if !agent.Search(key) {

		if key != "" {
			c.Key = key
		}

		//data, _ := model.Start(Service)
		//fmt.Println(string(data))

	} else {
		fmt.Printf("Already added key file %v.\n", c.Key)
	}
}