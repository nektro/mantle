package ws

import (
	"net/http"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/store"

	"github.com/gorilla/websocket"
)

const (
	keyOnline = "online_users"
)

var (
	reqUpgrader = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
)

var (
	// UserCache is the list of users currently with ws connections to this instance
	UserCache = map[string]*User{}
)

// Connect takes a db.User and upgrades it to a ws.User
func Connect(user *db.User, w http.ResponseWriter, r *http.Request) (*User, error) {
	conn, err := reqUpgrader.Upgrade(w, r, nil)
	if err != nil {
		return nil, err
	}
	u := &User{
		conn,
		user,
	}
	UserCache[u.User.UUID] = u

	if !u.IsConnected() {
		store.This.ListAdd(keyOnline, u.User.UUID)
		BroadcastMessage(map[string]interface{}{
			"type": "user-connect",
			"user": u.User.UUID,
		})
	}
	return u, nil
}

// Close disconnect all remaining users
func Close() {
	for _, item := range UserCache {
		item.Conn.Close()
	}
}

// BroadcastMessage sends message to all users
func BroadcastMessage(message map[string]interface{}) {
	for _, item := range UserCache {
		item.SendMessageRaw(message)
	}
}

// AllOnlineIDs returns ULID of every online user
func AllOnlineIDs() []string {
	return store.This.ListGet(keyOnline)
}

// OnlineUserCount is the total number of active users
func OnlineUserCount() int64 {
	return int64(store.This.ListLen(keyOnline))
}
