package ws

import (
	"container/list"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/itypes"

	"github.com/gorilla/websocket"
)

var (
	ReqUpgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	ConnCache   = map[string]itypes.ConnCacheValue{}
	RoleCache   = map[string]db.Role{}
	Connected   = list.New() // user UUIDs
)
