package main

import (
	"database/sql"

	"github.com/nektro/mantle/pkg/itypes"

	etc "github.com/nektro/go.etc"

	. "github.com/nektro/go-util/alias"
)

// // //

func scanChannel(rows *sql.Rows) itypes.RowChannel {
	var v itypes.RowChannel
	rows.Scan(&v.ID, &v.UUID, &v.Position, &v.Name, &v.Description)
	return v
}

func scanUser(rows *sql.Rows) itypes.RowUser {
	var v itypes.RowUser
	rows.Scan(&v.ID, &v.Provider, &v.Snowflake, &v.UUID, &v.IsMember, &v.IsBanned, &v.Name, &v.Nickname, &v.JoindedOn, &v.LastActive, &v.Roles)
	return v
}

func scanRole(rows *sql.Rows) itypes.RowRole {
	var v itypes.RowRole
	rows.Scan(v.ID, v.UUID, v.Position, v.Name, v.Color, v.PermManageChannels, v.PermManageRoles)
	return v
}

// // //

func queryAllChannels() []itypes.RowChannel {
	result := []itypes.RowChannel{}
	rows := etc.Database.Build().Se("*").Fr(cTableChannels).Exe()
	for rows.Next() {
		rch := scanChannel(rows)
		result = append(result, rch)
	}
	rows.Close()
	return result
}

func queryUserByUUID(uid string) (itypes.RowUser, bool) {
	rows := etc.Database.Build().Se("*").Fr(cTableUsers).Wh("uuid", uid).Exe()
	if !rows.Next() {
		return itypes.RowUser{}, false
	}
	ru := scanUser(rows)
	rows.Close()
	return ru, true
}

func queryUserBySnowflake(provider string, flake string, name string) itypes.RowUser {
	rows := etc.Database.Build().Se("*").Fr(cTableUsers).Wh("provider", provider).An("snowflake", flake).Exe()
	if rows.Next() {
		ru := scanUser(rows)
		rows.Close()
		return ru
	}
	// else
	id := etc.Database.QueryNextID(cTableUsers)
	uid := newUUID()
	now := T()
	roles := ""
	if id == 1 {
		roles += "o"
	}
	etc.Database.QueryPrepared(true, F("insert into %s values ('%d', '%s', '%s', '%s', '0', '0', ?, '', '%s', '%s', '%s')", cTableUsers, id, provider, flake, uid, now, now, roles), name)
	return queryUserBySnowflake(provider, flake, name)
}

func queryAssertUserName(uid string, name string) {
	etc.Database.Build().Up(cTableUsers, "name", name).Wh("uuid", uid).Exe()
}

func queryAllRoles() []itypes.RowRole {
	result := []itypes.RowRole{}
	rows := etc.Database.Build().Se("*").Fr(cTableRoles).Or("position", "asc").Exe()
	for rows.Next() {
		result = append(result, scanRole(rows))
	}
	rows.Close()
	return result
}
