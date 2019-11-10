package main

const (
	Name = "Mantle"

	cTableSettings         = "server_settings"
	cTableUsers            = "users"
	cTableChannels         = "channels"
	cTableRoles            = "roles"
	cTableChannelRolePerms = "channel_role_perms"
	cTableMessagesPrefix   = "channel_messages__"
)

const (
	PermDeny uint8 = iota
	PermIgnore
	PermAllow
)

func GetPermColumnRealVal(p uint8) bool {
	if p == PermAllow {
		return true
	}
	return false
}
