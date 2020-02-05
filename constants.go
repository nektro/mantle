package main

import (
	"github.com/nektro/mantle/pkg/itypes"
)

const (
	Name    = "Mantle"
	Version = "vMASTER"

	cTableSettings         = "server_settings"
	cTableUsers            = "users"
	cTableChannels         = "channels"
	cTableRoles            = "roles"
	cTableChannelRolePerms = "channel_role_perms"
	cTableMessagesPrefix   = "channel_messages__"
)

const (
	PermDeny itypes.Perm = iota
	PermIgnore
	PermAllow
)

func GetPermColumnRealVal(p itypes.Perm) bool {
	if p == PermAllow {
		return true
	}
	return false
}
