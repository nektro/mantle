package itypes

type RowChannelRolePerms struct {
	ID      int    `json:"id"`
	Channel string `json:"channel" sqlite:"text"`
	Role    string `json:"role" sqlite:"text"`
}
