package cache

import (
	"testing"
	"time"

	"github.com/go-redis/redis"
)

func TestCacheCreation(t *testing.T) {
	_, err := NewLRU(1000, 5)
	if err != nil {
		t.Error(err)
	}
}

func TestMaxCapacity(t *testing.T) {
	lru, err := NewLRU(1000, 5)
	if err != nil {
		panic(err)
	}

	a := redis.StringCmd{}
	b := redis.StringCmd{}
	c := redis.StringCmd{}
	d := redis.StringCmd{}
	e := redis.StringCmd{}

	lru.Set("a", &a)
	lru.Set("b", &b)
	lru.Set("c", &c)
	lru.Set("d", &d)
	lru.Set("e", &e)

	if lru.Get("a") != &a ||
		lru.Get("b") != &b ||
		lru.Get("c") != &c ||
		lru.Get("d") != &d ||
		lru.Get("e") != &e {

		t.Fail()
	}
}

func TestOverflowCapacity(t *testing.T) {
	lru, err := NewLRU(1000, 5)
	if err != nil {
		panic(err)
	}

	a := redis.StringCmd{}
	b := redis.StringCmd{}
	c := redis.StringCmd{}
	d := redis.StringCmd{}
	e := redis.StringCmd{}
	f := redis.StringCmd{}

	lru.Set("a", &a)
	lru.Set("b", &b)
	lru.Set("c", &c)
	lru.Set("d", &d)
	lru.Set("e", &e)
	lru.Set("f", &f)

	if lru.Get("a") != nil ||
		lru.Get("b") != &b ||
		lru.Get("c") != &c ||
		lru.Get("d") != &d ||
		lru.Get("e") != &e ||
		lru.Get("f") != &f {

		t.Fail()
	}
}

func TestCacheCycling(t *testing.T) {
	lru, err := NewLRU(1000, 3)
	if err != nil {
		panic(err)
	}

	a := redis.StringCmd{}
	b := redis.StringCmd{}
	c := redis.StringCmd{}
	d := redis.StringCmd{}
	e := redis.StringCmd{}
	f := redis.StringCmd{}

	lru.Set("a", &a)
	lru.Set("b", &b)
	lru.Set("c", &c)
	lru.Set("d", &d)
	lru.Set("e", &e)
	lru.Set("f", &f)

	if lru.Get("a") != nil ||
		lru.Get("b") != nil ||
		lru.Get("c") != nil ||
		lru.Get("d") != &d ||
		lru.Get("e") != &e ||
		lru.Get("f") != &f {

		t.Fail()
	}
}

func TestReassignment(t *testing.T) {
	lru, err := NewLRU(1000, 5)
	if err != nil {
		panic(err)
	}

	a := redis.StringCmd{}
	b := redis.StringCmd{}
	c := redis.StringCmd{}
	d := redis.StringCmd{}
	e := redis.StringCmd{}

	lru.Set("a", &a)
	lru.Set("b", &b)
	lru.Set("c->a", &c)
	lru.Set("d", &d)
	lru.Set("e", &e)
	lru.Set("c->a", &a)

	if lru.Get("a") != &a ||
		lru.Get("b") != &b ||
		lru.Get("c->a") != &a ||
		lru.Get("d") != &d ||
		lru.Get("e") != &e {

		t.Fail()
	}
}

func TestCacheTimeout(t *testing.T) {
	lru, err := NewLRU(100, 5)
	if err != nil {
		panic(err)
	}

	a := redis.StringCmd{}
	b := redis.StringCmd{}
	c := redis.StringCmd{}
	d := redis.StringCmd{}
	e := redis.StringCmd{}
	f := redis.StringCmd{}

	lru.Set("a", &a)
	time.Sleep(9 * time.Millisecond)
	lru.Set("b", &b)
	time.Sleep(9 * time.Millisecond)
	lru.Set("c", &c)
	time.Sleep(9 * time.Millisecond)
	lru.Set("d", &d)
	time.Sleep(9 * time.Millisecond)
	lru.Set("e", &e)
	time.Sleep(9 * time.Millisecond)
	lru.Set("f", &f)
	time.Sleep(9 * time.Millisecond)

	if lru.Get("a") != nil {
		t.Fail()
	}

	// b should have expired, c shouldn't have
	time.Sleep(50 * time.Millisecond)
	if lru.Get("b") != nil ||
		lru.Get("c") != &c {

		t.Fail()
	}

	// c should have expired
	time.Sleep(10 * time.Millisecond)
	if lru.Get("c") != nil {

		t.Fail()
	}

	// d should have expired
	time.Sleep(10 * time.Millisecond)
	if lru.Get("d") != nil {

		t.Fail()
	}

	// all should have expired
	time.Sleep(100 * time.Millisecond)
	if lru.Get("a") != nil ||
		lru.Get("b") != nil ||
		lru.Get("c") != nil ||
		lru.Get("d") != nil ||
		lru.Get("e") != nil ||
		lru.Get("f") != nil {

		t.Fail()
	}
}
