package ws

type UserPerms struct {
	ManageChannels bool `json:"manage_channels"`
	ManageRoles    bool `json:"manage_roles"`
}
