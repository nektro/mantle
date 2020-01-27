package main

import (
	"sync"

	"github.com/nektro/mantle/pkg/itypes"

	"github.com/nektro/go-util/util"
	etc "github.com/nektro/go.etc"

	. "github.com/nektro/go-util/alias"
)

type Properties struct {
	cache sync.Map
}

var (
	props = Properties{}
)

func (p *Properties) SetDefault(key string, value string) {
	id := etc.Database.QueryNextID(cTableSettings)
	etc.Database.QueryPrepared(true, F("insert into %s(id,key,value) select %d,'%s',? where not exists(select 1 from %s where key = '%s' and value = ?)", cTableSettings, id, key, cTableSettings, key), value, value)
	id2 := etc.Database.QueryNextID(cTableSettings)
	if id2 > id {
		util.Log(F("Added missing property '%s' with default value '%s'", key, value))
	}
}

func (p *Properties) Init() {
	p.cache = sync.Map{}
	rows := etc.Database.Build().Se("*").Fr(cTableSettings).Exe()
	for rows.Next() {
		sr := itypes.Setting{}
		rows.Scan(&sr.ID, &sr.Key, &sr.Value)
		p.cache.Store(sr.Key, sr.Value)
	}
	rows.Close()
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
	if ok {
		return val.(string)
	}
	return ""
}

func (p *Properties) Set(key string, val string) {
	etc.Database.Build().Up(cTableSettings, key, val)
	p.cache.Store(key, val)
}
