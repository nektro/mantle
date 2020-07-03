package idata

import (
	"github.com/nektro/go.etc/store"
	"github.com/nektro/mantle/pkg/itypes"
)

var (
	Name = "Mantle"

	Config = new(itypes.Config)
)

func InitStore() {
	c := Config

	if len(c.RedisURL) > 0 {
		store.Init("redis", c.RedisURL)
		return
	}
	store.Init("local", "")
}
