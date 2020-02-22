package ws

import (
	"strings"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/itypes"
)

type UserPerms struct {
	ManageChannels bool `json:"manage_channels"`
	ManageRoles    bool `json:"manage_roles"`
}

func (v UserPerms) From(user *db.User) *UserPerms {
	for _, item := range strings.Split(user.Roles, ",") {
		if item == "" {
			continue
		}
		role := roleCache[item]

		switch itypes.Perm(role.PermManageChannels) {
		case itypes.PermDeny, itypes.PermAllow:
			v.ManageChannels = itypes.Perm(role.PermManageChannels).ToBool()
		}
		switch itypes.Perm(role.PermManageRoles) {
		case itypes.PermDeny, itypes.PermAllow:
			v.ManageRoles = itypes.Perm(role.PermManageRoles).ToBool()
		}
	}
	return &v
}
