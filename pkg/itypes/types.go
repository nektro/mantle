package itypes

type APIResponse struct {
	Success bool        `json:"success"`
	Message interface{} `json:"message"`
}

type UserPerms struct {
	ManageChannels bool `json:"manage_channels"`
	ManageRoles    bool `json:"manage_roles"`
}
