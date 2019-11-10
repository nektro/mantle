package main

import (
	"database/sql"
)

// // //

func scanChannel(rows *sql.Rows) RowChannel {
	var v RowChannel
	rows.Scan(&v.ID, &v.UUID, &v.Position, &v.Name, &v.Description)
	return v
}

func scanUser(rows *sql.Rows) RowUser {
	var v RowUser
	rows.Scan(&v.ID, &v.Provider, &v.Snowflake, &v.UUID, &v.IsMember, &v.IsBanned, &v.Name, &v.Nickname, &v.JoindedOn, &v.LastActive, &v.Roles)
	return v
}

func scanRole(rows *sql.Rows) RowRole {
	var v RowRole
	rows.Scan(v.ID, v.UUID, v.Position, v.Name, v.Color, v.PermManageChannels, v.PermManageRoles)
	return v
}

// // //
