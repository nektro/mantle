package redis

import (
	"time"

	"github.com/bsm/redislock"
	"github.com/go-redis/redis/v7"
	"github.com/nektro/go-util/util"
)

// Store is the struct for this store.Inner implementation
type Store struct {
	c *redis.Client
	l *redislock.Client
	k *redislock.Lock
}

// Get returns a Store
func Get(addr string) *Store {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	return &Store{
		client,
		redislock.New(client), nil,
	}
}

// Type returns the loader for this Store
func (p *Store) Type() string {
	return "redis"
}

// Ping ensures we have a valid connection to this store
func (p *Store) Ping() error {
	_, err := p.c.Ping().Result()
	return err
}

// Has tests whether this Store contains a certain key
func (p *Store) Has(key string) bool {
	n, _ := p.c.Exists(key).Result()
	return n == 1
}

// Set sets the value of a single key
func (p *Store) Set(key string, val string) {
	p.c.Set(key, val, 0).Err()
}

// Get retrieves the value of a single key
func (p *Store) Get(key string) string {
	s, _ := p.c.Get(key).Result()
	return s
}

// Range loops over all values in this Store
func (p *Store) Range(f func(key string, val string) bool) {
	arr, _ := p.c.Keys("*").Result()
	for _, item := range arr {
		if !f(item, p.Get(item)) {
			break
		}
	}
}

// HasList returns true if there is a list associated with this key
func (p *Store) HasList(key string) bool {
	return p.Has(key)
}

// ListHas returns true if the list at this key contains value
func (p *Store) ListHas(key, value string) bool {
	n, _ := p.c.LRem(key, 0, value).Result()
	for i := 0; i < int(n); i++ {
		p.ListAdd(key, value)
	}
	return n > 0
}

// ListAdd appends value to the list at key
func (p *Store) ListAdd(key, value string) {
	p.c.LPush(key, value).Result()
}

// ListRemove removes value from the list at key
func (p *Store) ListRemove(key, value string) {
	p.c.LRem(key, 0, value)
}

// ListLen returns the number of items in the list at key
func (p *Store) ListLen(key string) int {
	i, _ := p.c.LLen(key).Result()
	return int(i)
}

// ListGet returns all items for the list at key
func (p *Store) ListGet(key string) []string {
	r, _ := p.c.LRange(key, 0, -1).Result()
	return r
}

// Lock locks this mutex
func (p *Store) Lock() {
	lock, err := p.l.Obtain("my_lock", time.Hour*24*30*12, nil)
	if err != nil {
		if err == redislock.ErrNotObtained {
			// block until obtained
			p.Lock()
			return
		}
		util.LogError("store/redis:", "lock error:", err)
		return
	}
	p.k = lock
}

// Unlock unlocks this mutex
func (p *Store) Unlock() {
	if p.k == nil {
		return
	}
	p.k.Release()
	p.k = nil
}
