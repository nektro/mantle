package ws

import (
	"container/list"

	"github.com/nektro/mantle/pkg/db"

	"github.com/gorilla/websocket"
)

type User struct {
	Conn  *websocket.Conn
	User  *db.User
	Perms *UserPerms
}

var (
	ReqUpgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	ConnCache   = map[string]User{}
	RoleCache   = map[string]db.Role{}
	Connected   = list.New() // user UUIDs
)

func BroadcastMessage(message map[string]string) {
	for _, item := range ConnCache {
		item.Conn.WriteJSON(message)
	}
}
