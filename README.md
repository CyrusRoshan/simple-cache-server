# simple-cache-server
a simple cache server in Go, for caching Redis key values

## Testing:
- `docker run -d -p 6379:6379 redis`
- `go test ./...`

(the Redis port can be changed, as long as you also add the `REDISADDRESS` env var, or add the `CONFIGFILE` env var, with the new redisAddress)

