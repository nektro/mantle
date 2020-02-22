package main

import (
	"database/sql"

	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/iconst"
	"github.com/nektro/mantle/pkg/itypes"

	. "github.com/nektro/go-util/alias"
)

// // //

func scanChannel(rows *sql.Rows) itypes.Channel {
	var v itypes.Channel
	rows.Scan(&v.ID, &v.UUID, &v.Position, &v.Name, &v.Description)
	return v
}

func scanUser(rows *sql.Rows) itypes.User {
	var v itypes.User
	rows.Scan(&v.ID, &v.Provider, &v.Snowflake, &v.UUID, &v.IsMember, &v.IsBanned, &v.Name, &v.Nickname, &v.JoindedOn, &v.LastActive, &v.Roles)
	return v
}

func scanRole(rows *sql.Rows) itypes.Role {
	var v itypes.Role
	rows.Scan(v.ID, v.UUID, v.Position, v.Name, v.Color, v.PermManageChannels, v.PermManageRoles)
	return v
}

// // //

func queryAllChannels() []itypes.Channel {
	result := []itypes.Channel{}
	rows := db.DB.Build().Se("*").Fr(iconst.TableChannels).Exe()
	for rows.Next() {
		rch := scanChannel(rows)
		result = append(result, rch)
	}
	rows.Close()
	return result
}

func queryUserByUUID(uid string) (itypes.User, bool) {
	rows := db.DB.Build().Se("*").Fr(iconst.TableUsers).Wh("uuid", uid).Exe()
	if !rows.Next() {
		return itypes.User{}, false
	}
	ru := scanUser(rows)
	rows.Close()
	return ru, true
}

func queryUserBySnowflake(provider string, flake string, name string) itypes.User {
	rows := db.DB.Build().Se("*").Fr(iconst.TableUsers).Wh("provider", provider).Wh("snowflake", flake).Exe()
	if rows.Next() {
		ru := scanUser(rows)
		rows.Close()
		return ru
	}
	// else
	id := db.DB.QueryNextID(iconst.TableUsers)
	uid := newUUID()
	now := T()
	roles := ""
	if id == 1 {
		roles += "o"
		props.Set("owner", uid)
	}
	db.DB.QueryPrepared(true, F("insert into %s values ('%d', '%s', '%s', '%s', '0', '0', ?, '', '%s', '%s', '%s')", iconst.TableUsers, id, provider, flake, uid, now, now, roles), name)
	return queryUserBySnowflake(provider, flake, name)
}

func queryAssertUserName(uid string, name string) {
	db.DB.Build().Up(iconst.TableUsers, "name", name).Wh("uuid", uid).Exe()
}

func queryAllRoles() []itypes.Role {
	result := []itypes.Role{}
	rows := db.DB.Build().Se("*").Fr(iconst.TableRoles).Or("position", "asc").Exe()
	for rows.Next() {
		result = append(result, scanRole(rows))
	}
	rows.Close()
	return result
}
