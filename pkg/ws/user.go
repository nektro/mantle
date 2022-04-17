package ws

import (
	"github.com/nektro/mantle/pkg/db"

	"github.com/gorilla/websocket"
	"github.com/nektro/go.etc/store"
)

type User struct {
	conn *websocket.Conn
	User *db.User
}

func (u *User) Disconnect() {
	if u.IsConnected() {
		delete(UserCache, u.User.UUID)
		store.This.ListRemove(keyOnline, u.User.UUID.String())
		db.Props.Decrement("count_users_online")
		BroadcastMessage(map[string]interface{}{
			"type": "user-disconnect",
			"user": u.User.UUID,
		})
		u.conn.Close()
	}
}

func (u *User) ReadMessage() (int, []byte, error) {
	return u.conn.ReadMessage()
}

func (u *User) IsConnected() bool {
	return store.This.ListHas(keyOnline, u.User.UUID.String())
}

func (u *User) SendWsMessage(msg map[string]interface{}) {
	u.conn.WriteJSON(msg)
}

func (u *User) SendWsMessageRaw(msg []byte) {
	u.conn.WriteMessage(websocket.TextMessage, msg)
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
		"at":      m.At.V().Format("2 Jan 2006 15:04:05 MST"),
	})
}
