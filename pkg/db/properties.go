package db

import (
	"sync"
)

type Properties struct {
	cache sync.Map
}

func (p *Properties) SetDefault(key string, value string) {
	_, ok := QuerySettingByKey(key)
	if ok {
		return
	}
	id := db.QueryNextID(cTableSettings)
	db.QueryPrepared(true, "insert into "+cTableSettings+" values (?,?,?)", id, key, value)
}

func (p *Properties) Init() {
	p.cache = sync.Map{}
	for _, item := range (Setting{}.All()) {
		p.cache.Store(item.Key, item.Value)
	}
}

func (p *Properties) GetAll() map[string]string {
	res := map[string]string{}
	p.cache.Range(func(key interface{}, val interface{}) bool {
		res[key.(string)] = val.(string)
		return true
	})
	return res
}

func (p *Properties) Get(key string) string {
	val, ok := p.cache.Load(key)
	if !ok {
		return ""
	}
	return val.(string)
}

func (p *Properties) Set(key string, val string) {
	p.SetDefault(key, "")
	db.Build().Up(cTableSettings, "value", val).Wh("key", key).Exe()
	p.cache.Store(key, val)
}
