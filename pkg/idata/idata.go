package idata

import (
	"container/list"

	"github.com/nektro/mantle/pkg/itypes"

	"github.com/gorilla/websocket"
)

var (
	Config      = new(itypes.Config)
	WsUpgrader  = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	WsConnCache = map[string]itypes.ConnCacheValue{}
	RoleCache   = map[string]itypes.Role{}
	Connected   = list.New()
)
