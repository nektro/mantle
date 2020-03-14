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
)

func Connect(user *db.User, w http.ResponseWriter, r *http.Request) (*User, error) {
	conn, err := reqUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	u := &User{
		conn,
		user,
		UserPerms{}.From(user),
	}
	UserCache[u.User.UUID] = u

	if !u.IsConnected() {
		connected.PushBack(u.User.UUID)
		BroadcastMessage(map[string]interface{}{
			"type": "user-connect",
			"user": u.User.UUID,
		})
	}
	return u, nil
}

func Close() {
	// disconnect all remaining users
	for _, item := range UserCache {
		item.Conn.Close()
	}
}

func BroadcastMessage(message map[string]interface{}) {
	for _, item := range UserCache {
		item.SendMessageRaw(message)
	}
}

func AllOnlineIDs() []string {
	return listToArray(connected)
}

func OnlineUserCount() int64 {
	return int64(len(UserCache))
}
