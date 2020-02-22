package main

import (
	"github.com/nektro/mantle/pkg/db"
	"github.com/nektro/mantle/pkg/iconst"

	dbstorage "github.com/nektro/go.dbstorage"

	. "github.com/nektro/go-util/alias"
)

func queryAllChannels() []db.Channel {
	arr := dbstorage.ScanAll(db.DB.Build().Se("*").Fr(iconst.TableChannels), db.Channel{})
	res := []db.Channel{}
	for _, item := range arr {
		res = append(res, *item.(*db.Channel))
	}
	return res
}

func queryUserByUUID(uid string) (*db.User, bool) {
	rows := db.DB.Build().Se("*").Fr(iconst.TableUsers).Wh("uuid", uid).Exe()
	if !rows.Next() {
		return &db.User{}, false
	}
	ru := db.User{}.Scan(rows).(*db.User)
	rows.Close()
	return ru, true
}

func queryUserBySnowflake(provider string, flake string, name string) *db.User {
	rows := db.DB.Build().Se("*").Fr(iconst.TableUsers).Wh("provider", provider).Wh("snowflake", flake).Exe()
	if rows.Next() {
		ru := db.User{}.Scan(rows).(*db.User)
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
		db.Props.Set("owner", uid)
	}
	db.DB.QueryPrepared(true, F("insert into %s values ('%d', '%s', '%s', '%s', '0', '0', ?, '', '%s', '%s', '%s')", iconst.TableUsers, id, provider, flake, uid, now, now, roles), name)
	return queryUserBySnowflake(provider, flake, name)
}

func queryAssertUserName(uid string, name string) {
	db.DB.Build().Up(iconst.TableUsers, "name", name).Wh("uuid", uid).Exe()
}

func queryAllRoles() []db.Role {
	arr := dbstorage.ScanAll(db.DB.Build().Se("*").Fr(iconst.TableRoles).Or("position", "asc"), db.Role{})
	res := []db.Role{}
	for _, item := range arr {
		res = append(res, *item.(*db.Role))
	}
	return res
}
