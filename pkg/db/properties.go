package db

import (
	"strconv"

	"github.com/nektro/mantle/pkg/store"
)

// Properties is a KV-store abstraction around a cached db table
type Properties struct {
	s *store.Store
}

// SetDefault sets a key's value only if that key has never been set
func (p *Properties) SetDefault(key string, value string) {
	_, ok := QuerySettingByKey(key)
	if ok {
		return
	}
	id := db.QueryNextID(cTableSettings)
	db.Build().Ins(cTableSettings, id, key, value).Exe()
}

// Init is called at all SetDefaults and loads cache from db table
func (p *Properties) Init() {
	p.s = store.This
	if !p.s.Has("loaded") {
		for _, item := range (Setting{}.All()) {
			p.s.Set(item.Key, item.Value)
		}
		p.s.Set("loaded", "1")
	}
}

// GetAll returns entire store in a map structure
func (p *Properties) GetAll() map[string]string {
	return p.GetSome(p.s.Keys()...)
}

// Get retrieves the value of a single key
func (p *Properties) Get(key string) string {
	return p.s.Get(key)
}

// Set sets the value of a single key
func (p *Properties) Set(key string, val string) {
	p.SetDefault(key, "")
	db.Build().Up(cTableSettings, "value", val).Wh("key", key).Exe()
	p.s.Set(key, val)
}

// Has tests whether this Properties contains a certain key
func (p *Properties) Has(key string) bool {
	return p.s.Has(key)
}

// GetSome returns a subset of the store in a map structure
func (p *Properties) GetSome(ks ...string) map[string]string {
	res := map[string]string{}
	for _, k := range ks {
		res[k] = p.Get(k)
	}
	return res
}

// SetDefaultInt64 sets key's value to value if never set before
func (p *Properties) SetDefaultInt64(key string, value int64) {
	p.SetDefault(key, strconv.FormatInt(value, 10))
}

// SetInt64 sets key's value to value
func (p *Properties) SetInt64(key string, value int64) {
	p.Set(key, strconv.FormatInt(value, 10))
}

// GetInt64 returns key's value as an int64
func (p *Properties) GetInt64(key string) int64 {
	i, _ := strconv.ParseInt(p.Get(key), 10, 64)
	return i
}

// Increment adds 1 to key's value if it is an integer
func (p *Properties) Increment(key string) {
	i, err := strconv.ParseInt(p.Get(key), 10, 64)
	if err != nil {
		return
	}
	p.SetInt64(key, i+1)
}

// Decrement subtracts 1 from key's value if it is an integer
func (p *Properties) Decrement(key string) {
	i, err := strconv.ParseInt(p.Get(key), 10, 64)
	if err != nil {
		return
	}
	p.SetInt64(key, i-1)
}
