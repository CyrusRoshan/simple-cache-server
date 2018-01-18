package proxy

import (
	"fmt"
	"net/http"

	"github.com/CyrusRoshan/simple-cache-server/cache"
)

func RedisProxyHandler(lru *cache.LRU) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, r.URL.Path)
	}
}
