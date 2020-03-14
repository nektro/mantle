package ws

import (
	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/itypes"
)

type UserPerms struct {
	ManageChannels bool `json:"manage_channels"`
	ManageRoles    bool `json:"manage_roles"`
	ManageServer   bool `json:"manage_server"`
}

// From calculates a user's permissions based on the roles they have
func (v UserPerms) From(user *db.User) *UserPerms {
	rls := user.GetRolesSorted()
	for i := len(rls) - 1; i >= 0; i-- {
		role := rls[i]

		switch itypes.Perm(role.PermManageChannels) {
		case itypes.PermDeny, itypes.PermAllow:
			v.ManageChannels = itypes.Perm(role.PermManageChannels).ToBool()
		}
		switch itypes.Perm(role.PermManageRoles) {
		case itypes.PermDeny, itypes.PermAllow:
			v.ManageRoles = itypes.Perm(role.PermManageRoles).ToBool()
		}
		switch itypes.Perm(role.PermManageServer) {
		case itypes.PermDeny, itypes.PermAllow:
			v.ManageServer = itypes.Perm(role.PermManageServer).ToBool()
		}
	}
	return &v
}
