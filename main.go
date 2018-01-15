package main

import (
	"fmt"
	"os"

	"github.com/CyrusRoshan/simple-cache-server/config"
)

const CONFIGFILE = "config.toml"

func main() {
	// read config
	configFile := os.Getenv("CONFIGFILE")
	if configFile == "" {
		configFile = CONFIGFILE
	}

	conf, err := config.LoadConfig(configFile)
	if err != nil {
		panic(err)
	}

	fmt.Println(conf)

	// start redis with config info
	// log redis running on info

	// start proxy with config info
	// log proxy running on info
}
