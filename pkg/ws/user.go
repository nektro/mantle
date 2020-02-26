package ws

import (
	"time"

	"github.com/nektro/mantle/pkg/db"

	"github.com/gorilla/websocket"
)

type User struct {
	Conn  *websocket.Conn
	User  *db.User
	Perms *UserPerms
}

func (u *User) Disconnect() {
	if u.IsConnected() {
		delete(UserCache, u.User.UUID)
		listRemove(connected, u.User.UUID)
		BroadcastMessage(map[string]interface{}{
			"type": "user-disconnect",
			"user": u.User.UUID,
		})
	}
}

func (u *User) IsConnected() bool {
	return listHas(connected, u.User.UUID)
}

func (u *User) SendMessageRaw(msg map[string]interface{}) {
	u.Conn.WriteJSON(msg)
}

func (u *User) SendMessage(in *db.Channel, msg string) {
	if len(msg) == 0 {
		return
	}
	m := db.CreateMessage(u.User, in, msg)
	t, _ := time.Parse("2006-01-02 15:04:05", m.At)
	BroadcastMessage(map[string]interface{}{
		"type":    "message",
		"in":      in.UUID,
		"from":    m.By,
		"message": m.Body,
		"at":      t.Format("2 Jan 2006 15:04:05 MST"),
	})
}
