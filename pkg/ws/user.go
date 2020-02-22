package ws

import (
	"net/http"

	"github.com/nektro/mantle/pkg/db"

	"github.com/gorilla/websocket"
)

type User struct {
	Conn  *websocket.Conn
	User  *db.User
	Perms *UserPerms
}

func (v User) From(r *http.Request, w http.ResponseWriter, user *db.User) *User {
	conn, _ := reqUpgrader.Upgrade(w, r, nil)
	return &User{
		conn,
		user,
		UserPerms{}.From(user),
	}
}

func (u *User) Connect() {
	connCache[u.User.UUID] = u

	if !u.IsConnected() {
		connected.PushBack(u.User.UUID)
		BroadcastMessage(map[string]string{
			"type": "user-connect",
			"user": u.User.UUID,
		})
	}
}

func (u *User) Disconnect() {
	if u.IsConnected() {
		delete(connCache, u.User.UUID)
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
