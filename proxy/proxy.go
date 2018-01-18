package proxy

import (
	"net/http"
	"net/url"

	"github.com/CyrusRoshan/simple-cache-server/cache"
	"github.com/go-redis/redis"
)

const IMPROPERLY_ENCODED_PATH = "Error - improperly urlencoded url path"
const KEY_NOT_FOUND = "Error - key not found"

func RedisProxyHandler(redisClient *redis.Client, lru *cache.LRU) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		key, err := url.QueryUnescape(r.URL.Path)
		if err != nil {
			w.WriteHeader(400)
			w.Write([]byte(IMPROPERLY_ENCODED_PATH))
			return
		}

		cachedVal := lru.Get(key)
		if cachedVal != nil {
			result, err := cachedVal.Result()

			if errIf(err, &w, r) {
				return
			}

			w.WriteHeader(200)
			w.Write([]byte(result))
			return
		}

		updatedValue := redisClient.Get(key)

		result, err := updatedValue.Result()
		if err == redis.Nil {
			w.WriteHeader(404)
			w.Write([]byte(KEY_NOT_FOUND))
			return
		} else if errIf(err, &w, r) {
			return
		}

		lru.Set(key, updatedValue)

		w.WriteHeader(200)
		w.Write([]byte(result))
		return
	}
}

func errIf(err error, w *http.ResponseWriter, r *http.Request) bool {
	if err != nil {
		(*w).WriteHeader(500)
		(*w).Write([]byte(err.Error()))

		return true
	}

	return false
}
