package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
)

var ErrMissingConfigField = errors.New("missing config field")

// Config info for Proxy and Redis
type Config struct {
	RedisAddress  string
	ProxyPort     int
	CacheExpiry   int
	CacheCapacity int
}

// LoadConfig loads config from file and ENV, ENV taking precedence
func LoadConfig(configFile string) (conf *Config, err error) {
	var config Config

	if _, fileErr := os.Stat(configFile); fileErr == nil {
		_, err = toml.DecodeFile(configFile, &config)
		if err != nil {
			return
		}
	}

	err = prioritizeEnvConfig(&config)
	if err != nil {
		return
	}

	if config.RedisAddress == "" ||
		config.ProxyPort == 0 ||
		config.CacheExpiry == 0 ||
		config.CacheCapacity == 0 {

		err = ErrMissingConfigField
		fmt.Println("test")
		return
	}

	return &config, nil
}

func prioritizeEnvConfig(config *Config) (err error) {
	if redisAddress := os.Getenv("REDISADDRESS"); redisAddress != "" {
		config.RedisAddress = redisAddress
	}

	if proxyPort := os.Getenv("PROXYPORT"); proxyPort != "" {
		var proxyPortInt int
		proxyPortInt, err = strconv.Atoi(proxyPort)
		if err != nil {
			return
		}

		config.ProxyPort = proxyPortInt
	}

	if cacheExpiry := os.Getenv("CACHEEXPIRY"); cacheExpiry != "" {
		var cacheExpiryInt int
		cacheExpiryInt, err = strconv.Atoi(cacheExpiry)
		if err != nil {
			return
		}

		config.CacheExpiry = cacheExpiryInt
	}

	if cacheCapacity := os.Getenv("CACHECAPACITY"); cacheCapacity != "" {
		var cacheCapacityInt int
		cacheCapacityInt, err = strconv.Atoi(cacheCapacity)
		if err != nil {
			return
		}

		config.CacheCapacity = cacheCapacityInt
	}

	return nil
}
