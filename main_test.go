package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"testing"
	"time"

	"github.com/CyrusRoshan/simple-cache-server/cache"
	"github.com/CyrusRoshan/simple-cache-server/config"
	"github.com/CyrusRoshan/simple-cache-server/proxy"
	"github.com/go-redis/redis"
)

var redisClient *redis.Client
var conf *config.Config
var basePath string

func requestBody(path string) (body string, err error) {
	resp, err := http.Get(basePath + path)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(bodyByte), nil
}

func TestSetup(t *testing.T) {
	os.Setenv("CACHEEXPIRY", "100")
	os.Setenv("CACHECAPACITY", "5")

	conf = getConfig()
	basePath = fmt.Sprintf("http://localhost:%v/", conf.ProxyPort)

	var pong string
	redisClient, pong = connectRedis(conf.RedisAddress)
	if pong != "PONG" {
		t.Error("No pong")
	}

	lru, err := cache.NewLRU(conf.CacheExpiry, conf.CacheCapacity)
	if err != nil {
		t.Error(err)
	}

	http.HandleFunc("/", proxy.RedisProxyHandler(redisClient, lru))
	portString := fmt.Sprintf(":%v", conf.ProxyPort)
	go http.ListenAndServe(portString, nil)
}

func TestRedisKeys(t *testing.T) {
	redisClient.Set("KEY1", "VAL1", time.Hour)
	redisClient.Set("KEY2", "VAL2", time.Hour)
	redisClient.Set("KEY3", "VAL3", time.Hour)

	if redisClient.Get("KEY1").Val() != "VAL1" ||
		redisClient.Get("KEY2").Val() != "VAL2" ||
		redisClient.Get("KEY3").Val() != "VAL3" {
		t.Error("Initial value mismatch")
	}

	redisClient.Set("KEY1", "VAL4", time.Hour)
	redisClient.Set("KEY2", "VAL5", time.Hour)

	if redisClient.Get("KEY1").Val() != "VAL4" ||
		redisClient.Get("KEY2").Val() != "VAL5" ||
		redisClient.Get("KEY3").Val() != "VAL3" {
		t.Error("Changed value mismatch")
	}
}

func TestProxy404(t *testing.T) {
	cmd := redisClient.FlushAll()
	err := cmd.Err()
	if err != nil {
		t.Error(err)
	}

	resp, err := http.Get(basePath + "undefined")
	if err != nil {
		t.Error(err)
	}

	if resp.StatusCode != 404 {
		fmt.Println(resp.StatusCode)
		t.Error("Wrong status code")
	}
}

func TestProxyIncorrectURLEncode(t *testing.T) {
	url := basePath + "t%2%-%%%est"

	// https://superuser.com/a/442395
	// just get curl status code
	curl := exec.Command("curl", "-s", "-o", "/dev/null", "-w", `"%{http_code}"`, url)

	output, err := curl.Output()
	if err != nil {
		t.Error(err)
	}

	if string(output) != `"400"` {
		t.Error("Non-400 status code response")
	}
}

func TestProxyCache(t *testing.T) {
	cmd := redisClient.FlushAll()
	err := cmd.Err()
	if err != nil {
		t.Error(err)
	}

	redisClient.Set("KEY1", "VAL1", time.Hour)
	redisClient.Set("KEY2", "VAL2", time.Hour)
	redisClient.Set("KEY3", "VAL3", time.Hour)

	body, err := requestBody("KEY1")
	if err != nil {
		t.Error(err)
	}
	if body != "VAL1" {
		t.Error("Initial value mismatch")
	}

	redisClient.Set("KEY1", "VAL4", time.Hour)

	body, err = requestBody("KEY1")
	if err != nil {
		t.Error(err)
	}
	if string(body) != "VAL1" {
		t.Error("Key value not cached")
	}
}

func TestProxyCacheUpdate(t *testing.T) {
	cmd := redisClient.FlushAll()
	err := cmd.Err()
	if err != nil {
		t.Error(err)
	}

	redisClient.Set("KEY1", "VAL1", time.Hour)
	redisClient.Set("KEY2", "VAL2", time.Hour)
	redisClient.Set("KEY3", "VAL3", time.Hour)

	body, err := requestBody("KEY1")
	if err != nil {
		t.Error(err)
	}
	if body != "VAL1" {
		t.Error("Initial value mismatch")
	}

	redisClient.Set("KEY1", "VAL4", time.Hour)
	time.Sleep(time.Duration(100) * time.Millisecond) // time out cache

	body, err = requestBody("KEY1")
	if err != nil {
		t.Error(err)
	}
	if string(body) != "VAL4" {
		t.Error("Cache not updated")
	}
}

func TestConcurrentClients(t *testing.T) {

}

func TestAsyncClients(t *testing.T) {

}
