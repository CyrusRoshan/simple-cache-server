package main

import (
	"testing"
	"time"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client

func TestRedisSetup(t *testing.T) {
	conf := getConfig()

	var pong string
	redisClient, pong = connectRedis(conf.RedisAddress)

	if pong != "PONG" {
		t.Fail()
	}
}

func TestRedisKeys(t *testing.T) {
	redisClient.Set("KEY1", "VAL1", time.Hour)
	redisClient.Set("KEY2", "VAL2", time.Hour)
	redisClient.Set("KEY3", "VAL3", time.Hour)

	if redisClient.Get("KEY1").Val() != "VAL1" ||
		redisClient.Get("KEY2").Val() != "VAL2" ||
		redisClient.Get("KEY3").Val() != "VAL3" {
		t.Fail()
	}

	redisClient.Set("KEY1", "VAL4", time.Hour)
	redisClient.Set("KEY2", "VAL5", time.Hour)

	if redisClient.Get("KEY1").Val() != "VAL4" ||
		redisClient.Get("KEY2").Val() != "VAL5" ||
		redisClient.Get("KEY3").Val() != "VAL3" {
		t.Fail()
	}
}

func TestProxy(t *testing.T) {

}

func TestProxyKeys(t *testing.T) {

}

func TestProxyCache(t *testing.T) {

}

func TestConcurrentClients(t *testing.T) {

}

func TestAsyncClients(t *testing.T) {

}
