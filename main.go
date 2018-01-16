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

	fmt.Println("Config read", *conf)

	// connect to redis
	// log success

	// start proxy with config info
	// log proxy running on info
}
