package db

import (
	"strings"

	"github.com/nektro/go-util/util"

	. "github.com/nektro/go-util/alias"
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
	now := T()
	roles := ""
	if id == 1 {
		roles += "o"
		Props.Set("owner", uid)
	}
	DB.QueryPrepared(true, F("insert into %s values ('%d', '%s', '%s', '%s', '0', '0', ?, '', '%s', '%s', '%s')", cTableUsers, id, provider, flake, uid, now, now, roles), name)
	return QueryUserBySnowflake(provider, flake, name)
}

func QueryAssertUserName(uid string, name string) {
	DB.Build().Up(cTableUsers, "name", name).Wh("uuid", uid).Exe()
}

func CreateRole(name string) string {
	id := DB.QueryNextID(cTableRoles)
	uid := newUUID()
	util.Log("[role-create]", uid, name)
	DB.QueryPrepared(true, F("insert into %s values ('%d', '%s', '%d', ?, '', 1, 1)", cTableRoles, id, uid, id), name)
	return uid
}

func CreateChannel(name string) string {
	id := DB.QueryNextID(cTableChannels)
	uid := newUUID()
	util.Log("[channel-create]", uid, "#"+name)
	DB.QueryPrepared(true, F("insert into %s values ('%d', '%s', '%d', ?, '')", cTableChannels, id, uid, id), name)
	AssertChannelMessagesTableExists(uid)
	return uid
}

func AssertChannelMessagesTableExists(uid string) {
	DB.CreateTableStruct(F("%s%s", cTableMessagesPrefix, strings.Replace(uid, "-", "_", -1)), Message{})
}
