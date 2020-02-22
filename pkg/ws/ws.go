package ws

import (
	"container/list"

	"github.com/nektro/mantle/pkg/db"

	"github.com/gorilla/websocket"
)

var (
	reqUpgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	connected   = list.New() // user UUIDs
)

var (
	ConnCache = map[string]*User{}
	RoleCache = map[string]db.Role{}
)

func BroadcastMessage(message map[string]string) {
	for _, item := range ConnCache {
		item.Conn.WriteJSON(message)
	}
}

func AllOnlineIDs() []string {
	return listToArray(connected)
}
