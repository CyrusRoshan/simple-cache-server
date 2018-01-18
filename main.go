package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-redis/redis"

	"github.com/CyrusRoshan/simple-cache-server/cache"
	"github.com/CyrusRoshan/simple-cache-server/config"
	"github.com/CyrusRoshan/simple-cache-server/proxy"
)

const CONFIGFILE = "config.toml"

func main() {
	conf := getConfig()
	fmt.Println("Config read", *conf)
	fmt.Println()

	redisClient, pong := connectRedis(conf.RedisAddress)
	fmt.Println("Successfully connected to redis, with a ping for a", pong, "| Client:", redisClient)
	fmt.Println()

	lru, err := cache.NewLRU(conf.CacheExpiry, conf.CacheCapacity)
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/", proxy.RedisProxyHandler(lru))
	portString := fmt.Sprintf(":%v", conf.ProxyPort)
	fmt.Println("Server running on port", conf.ProxyPort)
	log.Fatal(http.ListenAndServe(portString, nil))
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
