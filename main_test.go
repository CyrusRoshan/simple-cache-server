package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"os/exec"
	"sync"
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
var lru *cache.LRU

func TestProxyStartup(t *testing.T) {
	os.Setenv("CACHEEXPIRY", "50")
	os.Setenv("CACHECAPACITY", "5")

	conf = getConfig()
	basePath = fmt.Sprintf("http://localhost:%v/", conf.ProxyPort)

	var pong string
	redisClient, pong = connectRedis(conf.RedisAddress)
	if pong != "PONG" {
		t.Error("No pong")
	}

	var err error
	lru, err = cache.NewLRU(conf.CacheExpiry, conf.CacheCapacity)
	if err != nil {
		t.Error(err)
	}

	http.HandleFunc("/", proxy.RedisProxyHandler(redisClient, lru))
	portString := fmt.Sprintf(":%v", conf.ProxyPort)
	go http.ListenAndServe(portString, nil)
}

func TestRedisKeys(t *testing.T) {
	testSetup(t)

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
	testSetup(t)

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
	testSetup(t)

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
	testSetup(t)

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
	testSetup(t)

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
	time.Sleep(time.Duration(50) * time.Millisecond) // time out cache

	body, err = requestBody("KEY1")
	if err != nil {
		t.Error(err)
	}
	if string(body) != "VAL4" {
		t.Error("Cache not updated")
	}
}

func TestConcurrentClients(t *testing.T) {
	testSetup(t)

	pairs := map[string]string{
		"KEY1": "VAL1",
		"KEY2": "VAL2",
		"KEY3": "VAL3",
		"KEY4": "VAL4",
		"KEY5": "VAL5",
	}

	for key, val := range pairs {
		redisClient.Set(key, val, time.Hour)
	}

	var wg sync.WaitGroup
	testRequest := func(key string, val string, sleepTime int) {
		defer wg.Done()
		time.Sleep(time.Duration(sleepTime) * time.Millisecond)

		resp, err := http.Get(basePath + key)
		if err != nil {
			t.Error(err)
		}

		bodyByte, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Error(err)
		}

		err = resp.Body.Close()
		if err != nil {
			t.Error(err)
		}

		if string(bodyByte) != val {
			t.Error("Value mismatch", bodyByte, val)
		}
	}

	for i := 0; i < 1000; i++ {
		for key, val := range pairs {
			wg.Add(1)
			go testRequest(key, val, rand.Intn(10))
		}
	}

	wg.Wait()
}

func TestAsyncClients(t *testing.T) {

}

// utility functions
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

func testSetup(t *testing.T) {
	lru.Clear()

	cmd := redisClient.FlushAll()
	err := cmd.Err()
	if err != nil {
		t.Error(err)
	}
}
