package ws

import (
	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/store"

	"github.com/gorilla/websocket"
)

type User struct {
	Conn *websocket.Conn
	User *db.User
}

func (u *User) Disconnect() {
	if u.IsConnected() {
		delete(UserCache, u.User.UUID)
		store.This.ListRemove(keyOnline, u.User.UUID)
		BroadcastMessage(map[string]interface{}{
			"type": "user-disconnect",
			"user": u.User.UUID,
		})
	}
}

func (u *User) IsConnected() bool {
	return store.This.ListHas(keyOnline, u.User.UUID)
}

func (u *User) SendMessageRaw(msg map[string]interface{}) {
	u.Conn.WriteJSON(msg)
}

func (u *User) SendMessage(in *db.Channel, msg string) {
	if len(msg) == 0 {
		return
	}
	m := db.CreateMessage(u.User, in, msg)
	BroadcastMessage(map[string]interface{}{
		"type":    "message-new",
		"in":      in.UUID,
		"from":    u,
		"message": m,
		"at":      m.At.T().Format("2 Jan 2006 15:04:05 MST"),
	})
}
