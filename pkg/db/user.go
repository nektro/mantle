package db

import (
	"database/sql"
)

type User struct {
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

func ScanUser(rows *sql.Rows) *User {
	var v User
	rows.Scan(&v.ID, &v.Provider, &v.Snowflake, &v.UUID, &v.IsMember, &v.IsBanned, &v.Name, &v.Nickname, &v.JoindedOn, &v.LastActive, &v.Roles)
	return &v
}
