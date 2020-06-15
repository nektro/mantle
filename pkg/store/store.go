package store

import (
	"sync"

	"github.com/nektro/mantle/pkg/store/local"
)

// global singleton
var (
	This *Store
)

// PreInit registers flags
func PreInit() {
}

// Init takes flag values and initializes datastore
	This = &Store{local.Get()}
}

// Inner holds KV values
type Inner interface {
	Type() string
	Has(key string) bool
	Set(key string, val string)
	Get(key string) string
	Range(f func(key string, val string) bool)
	HasList(key string) bool
	ListHas(key, value string) bool
	ListAdd(key, value string)
	ListRemove(key, value string)
	ListLen(key string) int
	ListGet(key string) []string
	sync.Locker
}

// Store is
type Store struct {
	Inner
}

// Keys returns array of all Keys
func (p Store) Keys() []string {
	res := []string{}
	p.Range(func(k, _ string) bool {
		res = append(res, k)
		return true
	})
	return res
}
