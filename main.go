package main

import (
	"fmt"
	"os"

	"github.com/go-redis/redis"

	"github.com/CyrusRoshan/simple-cache-server/config"
)

const CONFIGFILE = "config.toml"

func main() {
	conf := getConfig()
	fmt.Println("Config read", *conf)
	fmt.Println()

	// log success
	redisClient, pong := connectRedis(conf.RedisAddress)
	fmt.Println("Successfully connected to redis, with a ping for a", pong, "| Client:", redisClient)
	fmt.Println()

	// start proxy with config info
	// log proxy running on info
}

func getConfig() *config.Config {
	configFile := os.Getenv("CONFIGFILE")
	if configFile == "" {
		configFile = CONFIGFILE
	}

	conf, err := config.LoadConfig(configFile)
	if err != nil {
		panic(err)
	}

	return conf
}

func connectRedis(address string) (*redis.Client, string) {
	client := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: "",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		panic(err)
	}

	return client, pong
}
