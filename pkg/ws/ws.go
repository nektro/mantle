package ws

import (
	"container/list"

	"github.com/nektro/mantle/pkg/db"

	"github.com/gorilla/websocket"
)

var (
	reqUpgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	connCache   = map[string]*User{}
	roleCache   = map[string]db.Role{}
	connected   = list.New() // user UUIDs
)

func BroadcastMessage(message map[string]string) {
	for _, item := range connCache {
		item.Conn.WriteJSON(message)
	}
}

func AllOnlineIDs() []string {
	return listToArray(connected)
}
