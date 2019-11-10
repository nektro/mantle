package main

type RowSetting struct {
	ID    int    `json:"id"`
	Key   string `json:"key" sqlite:"text"`
	Value string `json:"value" sqlite:"text"`
}

type RowUser struct {
	ID         int    `json:"id"`
	Provider   string `json:"provider" sqlite:"text"`
	Snowflake  string `json:"snowflake" sqlite:"text"`
	UUID       string `json:"uuid" sqlite:"text"`
	IsMember   bool   `json:"is_member" sqlite:"tinyint(1)"`
	IsBanned   bool   `json:"is_banned" sqlite:"tinyint(1)"`
	Name       string `json:"name" sqlite:"text"`
	Nickname   string `json:"nickname" sqlite:"text"`
	JoindedOn  string `json:"joined_on" sqlite:"text"`
	LastActive string `json:"last_active" sqlite:"text"`
	Roles      string `json:"roles" sqlite:"text"`
}

type RowChannel struct {
	ID          int    `json:"id"`
	UUID        string `json:"uuid" sqlite:"text"`
	Position    int    `json:"position" sqlite:"int"`
	Name        string `json:"name" sqlite:"text"`
	Description string `json:"description" sqlite:"text"`
}

type RowRole struct {
	ID                 int    `json:"id"`
	UUID               string `json:"uuid" sqlite:"text"`
	Position           int    `json:"position" sqlite:"int"`
	Name               string `json:"name" sqlite:"text"`
	Color              string `json:"color" sqlite:"text"`
	PermManageChannels uint8  `json:"perm_manage_channels" sqlite:"tinyint(1)"`
	PermManageRoles    uint8  `json:"perm_manage_roles" sqlite:"tinyint(1)"`
}

type RowChannelRolePerms struct {
	ID      int    `json:"id"`
	Channel string `json:"channel" sqlite:"text"`
	Role    string `json:"role" sqlite:"text"`
}
