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
	conn, _ := ReqUpgrader.Upgrade(w, r, nil)
	return &User{
		conn,
		user,
		UserPerms{}.From(user),
	}
}
