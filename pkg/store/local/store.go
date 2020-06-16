package local

import (
	"sync"

	"github.com/nektro/go-util/util"
)

// Store is the struct for this store.Inner implementation
type Store struct {
	m map[string]string
	l map[string][]string
	*sync.Mutex
}

// Get returns a Store
func Get() *Store {
	return &Store{
		map[string]string{},
		map[string][]string{},
		new(sync.Mutex),
	}
}

// Type returns the loader for this Store
func (p *Store) Type() string {
	return "local"
}

// Ping ensures we have a valid connection to this store
func (p *Store) Ping() error {
	return nil
}

// Has tests whether this Store contains a certain key
func (p *Store) Has(key string) bool {
	_, ok := p.m[key]
	return ok
}

// Set sets the value of a single key
func (p *Store) Set(key string, val string) {
	p.m[key] = val
}

// Get retrieves the value of a single key
func (p *Store) Get(key string) string {
	val, ok := p.m[key]
	if !ok {
		return ""
	}
	return val
}

// Range loops over all values in this Store
func (p *Store) Range(f func(key string, val string) bool) {
	for k, v := range p.m {
		if !f(k, v) {
			break
		}
	}
}

// HasList returns true if there is a list associated with this key
func (p *Store) HasList(key string) bool {
	_, ok := p.l[key]
	return ok
}

func (p *Store) initListK(key string) {
	if !p.HasList(key) {
		p.l[key] = []string{}
	}
}

// ListHas returns true if the list at this key contains value
func (p *Store) ListHas(key, value string) bool {
	return util.Contains(p.ListGet(key), value)
}

// ListAdd appends value to the list at key
func (p *Store) ListAdd(key, value string) {
	p.initListK(key)
	p.l[key] = append(p.l[key], value)
}

// ListRemove removes value from the list at key
func (p *Store) ListRemove(key, value string) {
	r := []string{}
	for _, item := range p.ListGet(key) {
		if item != value {
			r = append(r, item)
		}
	}
	p.l[key] = r
}

// ListLen returns the number of items in the list at key
func (p *Store) ListLen(key string) int {
	return len(p.ListGet(key))
}

// ListGet returns all items for the list at key
func (p *Store) ListGet(key string) []string {
	p.initListK(key)
	return p.l[key]
}
