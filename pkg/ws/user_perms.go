package ws

import (
	"github.com/nektro/mantle/pkg/db"
)

type UserPerms struct {
	ManageChannels bool `json:"manage_channels"`
	ManageRoles    bool `json:"manage_roles"`
	ManageServer   bool `json:"manage_server"`
	ManageInvites  bool `json:"manage_invites"`
	ViewAudits     bool `json:"view_audits"`
}

// From calculates a user's permissions based on the roles they have
func (v UserPerms) From(user *db.User) *UserPerms {
	rls := user.GetRolesSorted()
	for i := len(rls) - 1; i >= 0; i-- {
		role := rls[i]

		switch role.PermManageChannels {
		case db.PermDeny, db.PermAllow:
			v.ManageChannels = role.PermManageChannels.ToBool()
		}
		switch role.PermManageRoles {
		case db.PermDeny, db.PermAllow:
			v.ManageRoles = role.PermManageRoles.ToBool()
		}
		switch role.PermManageServer {
		case db.PermDeny, db.PermAllow:
			v.ManageServer = role.PermManageServer.ToBool()
		}
		switch role.PermManageInvites {
		case db.PermDeny, db.PermAllow:
			v.ManageInvites = role.PermManageInvites.ToBool()
		}
		switch role.PermViewAudits {
		case db.PermDeny, db.PermAllow:
			v.ViewAudits = role.PermViewAudits.ToBool()
		}
	}
	return &v
}
