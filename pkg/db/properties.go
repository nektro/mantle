package db

import (
	"sync"

	"github.com/nektro/mantle/pkg/iconst"

	"github.com/nektro/go-util/util"

	. "github.com/nektro/go-util/alias"
)

type Properties struct {
	cache sync.Map
}

func (p *Properties) SetDefault(key string, value string) {
	id := DB.QueryNextID(iconst.TableSettings)
	DB.QueryPrepared(true, F("insert into %s(id,key,value) select %d,'%s',? where not exists(select 1 from %s where key = '%s' and value = ?)", iconst.TableSettings, id, key, iconst.TableSettings, key), value, value)
	id2 := DB.QueryNextID(iconst.TableSettings)
	if id2 > id {
		util.Log(F("Added missing property '%s' with default value '%s'", key, value))
	}
}

func (p *Properties) Init() {
	p.cache = sync.Map{}
	rows := DB.Build().Se("*").Fr(iconst.TableSettings).Exe()
	for rows.Next() {
		sr := Setting{}
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
	DB.Build().Up(iconst.TableSettings, key, val)
	p.cache.Store(key, val)
}
