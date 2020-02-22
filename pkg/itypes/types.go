package itypes

import (
	"github.com/nektro/mantle/pkg/db"

	"github.com/gorilla/websocket"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message interface{} `json:"message"`
}

type ConnCacheValue struct {
	Conn  *websocket.Conn
	User  *db.User
	Perms *UserPerms
}

type UserPerms struct {
	ManageChannels bool `json:"manage_channels"`
	ManageRoles    bool `json:"manage_roles"`
}
