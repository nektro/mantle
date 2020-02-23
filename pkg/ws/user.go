package ws

import (
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
		BroadcastMessage(map[string]string{
			"type": "user-disconnect",
			"user": u.User.UUID,
		})
	}
}

func (u *User) IsConnected() bool {
	return listHas(connected, u.User.UUID)
}

func (u *User) SendMessage(msg map[string]string) {
	u.Conn.WriteJSON(msg)
}
