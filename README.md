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

## Algorithmic complexity for LRU operations

- Set: O(1) for map access, O(1) for accessing the list element the map object points to, O(1) for either:
    - (if list element exists) moving the list element to front
    - (if list element doesn't exist) deleting the last element from the (doubly linked) list, deleting the corresponding map k/v pair, and pushing the new element to the front, then assigning the new elem's key to the new elem, in the map). Probably could have been written better, but assuming pointer access and map access is O(1), this is an O(1) operation.
- Get: O(1) for map access, O(1) for accessing the list element the map object points to, O(1) for checking expiry time, and either:
    - O(1) for moving element to front of list and returning it
    - O(1) for deleting the element's corresponding k/v pair (if element is lazy expired) and returning nil
- Clear: lru.list is not modified, lru.lookup points to a new map, and the old map is garbage collected. Therefore, the Big O is O(1).

## How long I spent on each part:
In total, around 6-7hrs as a liberal estimate, not including this documentation commit.

I spent a while selecting the right redis library, reading a bit about how it works, looking into a good YAML alternative for simple config files (I went with TOML). Maybe an hour on config-related stuff, then a couple hours on the cache. Then 30-45m on the server part of the proxy. Then the rest of the time writing tests, adding concurrency, noticing issues with my implementation, cleaning up code, cleaning up tests, writing documentation, and recovering from cache misses inbetween putting down and picking up work.

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
- [ ] Platform (Dockerizing proxy) (Planned)
- [x] Parallel concurrent processing 
- [ ] Redis client protocol (not planned, but there's always [this](https://github.com/quorzz/redis-protocol), I guess)

