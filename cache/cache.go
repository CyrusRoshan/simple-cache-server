package cache

import (
	"container/list"
	"errors"
	"sync"
	"time"

	"github.com/go-redis/redis"
)

var ErrNegativeValues = errors.New("negative cache expiry and capacity values unsupported")

type element struct {
	Timestamp time.Time
	Key       string
	Response  *redis.StringCmd
}

type LRU struct {
	capacity int
	expiry   time.Duration
	list     *list.List
	lookup   map[string]*list.Element
	mutex    *sync.Mutex
}

func NewLRU(expiry int, capacity int) (lru *LRU, err error) {
	if expiry < 0 || capacity < 0 {
		return nil, ErrNegativeValues
	}

	lru = &LRU{
		capacity: capacity,
		expiry:   time.Duration(expiry) * time.Millisecond,
		lookup:   make(map[string]*list.Element, capacity),
		list:     list.New(),
		mutex:    &sync.Mutex{},
	}

	for i := 0; i < capacity; i++ {
		lru.list.PushFront(nil)
	}

	return lru, nil
}

func (lru *LRU) Set(key string, value *redis.StringCmd) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	cacheElement := element{
		Timestamp: time.Now(),
		Key:       key,
		Response:  value,
	}

	listElement, exists := lru.lookup[key]
	if exists {
		lru.list.MoveToFront(listElement)
		listElement.Value = &cacheElement

		return
	}

	lru.deleteElem(lru.list.Back())
	listElement = lru.list.PushFront(&cacheElement)
	lru.lookup[key] = listElement

	return
}

func (lru *LRU) Clear() {
	lru.mutex.Lock()
	lru.lookup = make(map[string]*list.Element, lru.capacity)
	lru.mutex.Unlock()
}

func (lru *LRU) Get(key string) (response *redis.StringCmd) {
	lru.mutex.Lock()
	defer lru.mutex.Unlock()

	listElement, exists := lru.lookup[key]
	if !exists {
		return nil
	}

	cacheElement := listElement.Value.(*element)
	expiryTime := cacheElement.Timestamp.Add(lru.expiry)

	if time.Now().Before(expiryTime) {
		lru.list.MoveToFront(listElement)

		return cacheElement.Response
	}

	delete(lru.lookup, cacheElement.Key)
	// We don't also delete listElement from list: we're preserving list size

	return nil
}

func (lru *LRU) deleteElem(listElement *list.Element) {
	if listElement.Value == nil { // the initial empty elements
		lru.list.Remove(listElement)
		return
	}

	cacheElement := listElement.Value.(*element)

	delete(lru.lookup, cacheElement.Key)
	lru.list.Remove(listElement)
}
