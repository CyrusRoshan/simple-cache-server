# simple-cache-server
a simple cache server in Go, for caching Redis key values in an LRU cache

## Running:
- `docker run -d -p 6379:6379 redis`
- `go run main.go`

This also supports the following environment variables:
- `CONFIGFILE`: Config file location. See `config.toml` for an example. Config settings are overridden by env var settings
- `REDISADDRESS`: Redis server address, including port
- `PROXYPORT`: Port to bind the proxy server to
- `CACHEEXPIRY`: How long to let items remain in the cache until they're invalidated
- `CACHECAPACITY`: Maximum number of items to keep in the cache at a given time

## Testing:
- `docker run -d -p 6379:6379 redis`
- `go test ./...`

(the Redis port can be changed, as long as you also add the `REDISADDRESS` env var, or add the `CONFIGFILE` env var, with the new redisAddress)

## High level architecture overview

Client <-> Proxy (with LRU cache) <-> Redis

Requests go from the client, to the proxy, and if cached, are returned from the cache to the client. If uncached, the proxy fetches the value from the Redis instance, and adds it to the LRU cache.

## What the code does

Works.

## What the code actually does

After main.go is run, the proxy checks the config file location, reads the given file (or default file if none is given), and parses all config options from it. Config options can also be given by environment variables, in which case they will override options from the file.

The proxy then connects to the redis instance, using config options, and pings it to make sure it's working.

The proxy then creates and initializes an empty LRU cache, using config options.

The proxy then runs on the configured port, handling requests in the form `${BASEURL}:${PORT}/${KEY}`. Manual testing should be simple with curl.

When recieving a request, the proxy first checks for the key value in the LRU. If it doesn't exist, or is out of date (the LRU uses lazy expiration), the proxy fetches the new value from the Redis instance, and updates the LRU with the new value, then serves it back to the client.

If the client's request is in the LRU, it's of course served back, and the key's position in the cache is moved to the start.

# What's done:

## Required: 
- [x] HTTP web service
- [x] Single backing instance
- [x] Cached GET
- [x] Global expiry
- [x] LRU eviction
- [x] Fixed key size
- [x] Sequential concurrent processing
- [x] Configuration
- [x] System tests
- [x] Single-click build and test
- [x] Documentation

## Optional
- [ ] Platform (Dockerizing proxy) (In progress)
- [x] Parallel concurrent processing 
- [ ] Redis client protocol (not planned, but there's always [this](https://github.com/quorzz/redis-protocol), I guess)

