package ws

import (
	"container/list"
	"net/http"

	"github.com/nektro/mantle/pkg/db"

	"github.com/gorilla/websocket"
)

var (
	reqUpgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	connected   = list.New() // user UUIDs
)

var (
	UserCache = map[string]*User{}
	RoleCache = map[string]db.Role{}
)

func Connect(user *db.User, w http.ResponseWriter, r *http.Request) *User {
	conn, _ := reqUpgrader.Upgrade(w, r, nil)
	u := &User{
		conn,
		user,
		UserPerms{}.From(user),
	}
	UserCache[u.User.UUID] = u

	if !u.IsConnected() {
		connected.PushBack(u.User.UUID)
		BroadcastMessage(map[string]string{
			"type": "user-connect",
			"user": u.User.UUID,
		})
	}
	return u
}

func BroadcastMessage(message map[string]string) {
	for _, item := range UserCache {
		item.Conn.WriteJSON(message)
	}
}

func AllOnlineIDs() []string {
	return listToArray(connected)
}
