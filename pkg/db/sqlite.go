package db

import (
	"strings"

	"github.com/nektro/go-util/alias"
	"github.com/nektro/go-util/util"
)

func QueryUserByUUID(uid string) (*User, bool) {
	rows := DB.Build().Se("*").Fr(cTableUsers).Wh("uuid", uid).Exe()
	if !rows.Next() {
		return &User{}, false
	}
	ru := User{}.Scan(rows).(*User)
	rows.Close()
	return ru, true
}

func QueryUserBySnowflake(provider string, flake string, name string) *User {
	rows := DB.Build().Se("*").Fr(cTableUsers).Wh("provider", provider).Wh("snowflake", flake).Exe()
	if rows.Next() {
		ru := User{}.Scan(rows).(*User)
		rows.Close()
		return ru
	}
	// else
	id := DB.QueryNextID(cTableUsers)
	uid := newUUID()
	now := alias.T()
	roles := ""
	if id == 1 {
		roles += "o"
		Props.Set("owner", uid)
	}
	DB.QueryPrepared(true, "insert into "+cTableUsers+" values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)", id, provider, flake, uid, 0, 0, name, "", now, now, roles)
	return QueryUserBySnowflake(provider, flake, name)
}

func CreateRole(name string) string {
	id := DB.QueryNextID(cTableRoles)
	uid := newUUID()
	util.Log("[role-create]", uid, name)
	DB.QueryPrepared(true, "insert into "+cTableRoles+" values (?, ?, ?, ?, '', 1, 1)", id, uid, id, name)
	return uid
}

func CreateChannel(name string) string {
	id := DB.QueryNextID(cTableChannels)
	uid := newUUID()
	util.Log("[channel-create]", uid, "#"+name)
	DB.QueryPrepared(true, "insert into "+cTableChannels+" values (?, ?, ?, ?, '')", id, uid, id, name)
	AssertChannelMessagesTableExists(uid)
	return uid
}

func AssertChannelMessagesTableExists(uid string) {
	DB.CreateTableStruct(cTableMessagesPrefix+strings.Replace(uid, "-", "_", -1), Message{})
}
