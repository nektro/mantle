package db

import (
	"database/sql"
	"sort"
	"strconv"

	"github.com/nektro/mantle/pkg/store"

	"github.com/nektro/go-util/arrays/stringsu"
	"github.com/nektro/go-util/util"
	dbstorage "github.com/nektro/go.dbstorage"
)

type User struct {
	ID         int64  `json:"id"`
	Provider   string `json:"provider" sqlite:"text"`
	Snowflake  string `json:"snowflake" sqlite:"text"`
	UUID       string `json:"uuid" sqlite:"text"`
	IsMember   bool   `json:"is_member" sqlite:"tinyint(1)"`
	IsBanned   bool   `json:"is_banned" sqlite:"tinyint(1)"`
	Name       string `json:"name" sqlite:"text"`
	Nickname   string `json:"nickname" sqlite:"text"`
	JoindedOn  Time   `json:"joined_on" sqlite:"text"`
	LastActive Time   `json:"last_active" sqlite:"text"`
	Roles      Array  `json:"roles" sqlite:"text"`
}

//
//

func QueryUserByUUID(uid string) (*User, bool) {
	ch, ok := dbstorage.ScanFirst(db.Build().Se("*").Fr(cTableUsers).Wh("uuid", uid), User{}).(*User)
	return ch, ok
}

func QueryUserBySnowflake(provider string, flake string, name string) *User {
	us, ok := dbstorage.ScanFirst(db.Build().Se("*").Fr(cTableUsers).Wh("provider", provider).Wh("snowflake", flake), User{}).(*User)
	if ok {
		return us
	}
	store.This.Lock()
	defer store.This.Unlock()
	//
	id := db.QueryNextID(cTableUsers)
	uid := newUUID()
	co := now()
	roles := Array{}
	if id == 1 {
		roles = append(roles, "o")
		Props.Set("owner", uid)
	}
	u := &User{id, provider, flake, uid, false, false, name, "", co, co, roles}
	db.Build().InsI(cTableUsers, u).Exe()
	Props.Increment("count_" + cTableUsers)
	return u
}

//
//

// Scan implements dbstorage.Scannable
func (v User) Scan(rows *sql.Rows) dbstorage.Scannable {
	rows.Scan(&v.ID, &v.Provider, &v.Snowflake, &v.UUID, &v.IsMember, &v.IsBanned, &v.Name, &v.Nickname, &v.JoindedOn, &v.LastActive, &v.Roles)
	return &v
}

// Count returns the total number of Users
func (v User) Count() int64 {
	return db.QueryRowCount(cTableUsers)
}

func (v User) MemberCount() int64 {
	rows := db.Build().Se("count(*)").Fr(cTableUsers).Wh("is_member", "1").Exe()
	defer rows.Close()
	c := int64(0)
	rows.Next()
	rows.Scan(&c)
	return c
}

//
//

func (u *User) SetAsMember(b bool) {
	m := u.IsMember
	db.Build().Up(cTableUsers, "is_member", strconv.Itoa(util.Btoi(b))).Wh("uuid", u.UUID).Exe()
	if !m && b {
		Props.Increment("count_users_members")
	}
	if m && !b {
		Props.Decrement("count_users_members")
	}
	u.IsMember = b
}

func (u *User) SetName(s string) {
	db.Build().Up(cTableUsers, "name", s).Wh("uuid", u.UUID).Exe()
	u.Name = s
}

// DeleteMessage attempts to delete a UID from this Channel's associated message
// table. If the UID is not a message in this Channel, nothing happens.
func (u *User) DeleteMessage(c *Channel, uid string) {
	db.Build().Del(cTableMessagesPrefix+c.UUID).Wh("uuid", uid).Wh("author", u.UUID).Exe()
}

func (u *User) HasRole(role string) bool {
	return stringsu.Contains(u.Roles, role)
}

func (u *User) AddRole(role string) {
	if u.HasRole(role) {
		return
	}
	u.Roles = append(u.Roles, role)
	db.Build().Up(cTableUsers, "roles", u.Roles.String()).Wh("uuid", u.UUID).Exe()
}

func (u *User) RemoveRole(role string) {
	if !u.HasRole(role) {
		return
	}
	u.Roles = stringsu.Remove(u.Roles, role)
	db.Build().Up(cTableUsers, "roles", u.Roles.String()).Wh("uuid", u.UUID).Exe()
}

func (u *User) GetRoles() []*Role {
	res := []*Role{}
	for _, item := range u.Roles {
		r, ok := QueryRoleByUID(item)
		if !ok {
			continue
		}
		res = append(res, r)
	}
	return res
}

func (u *User) GetRolesSorted() []*Role {
	res := u.GetRoles()
	sort.Slice(res, func(i, j int) bool {
		return res[i].Position < res[j].Position
	})
	return res
}

func (u *User) SetUID(uid string) {
	oid := u.UUID
	db.Build().Up(cTableUsers, "uuid", uid).Wh("uuid", u.UUID).Exe()
	for _, item := range (Channel{}.All()) {
		db.Build().Up(cTableMessagesPrefix+item.UUID, "author", uid).Wh("author", u.UUID).Exe()
	}
	u.UUID = uid
	util.Log("user-update:", "updated", u.Name+"#"+strconv.FormatInt(u.ID, 10), "from", oid, "to", u.UUID)
}

func (u *User) ResetUID() {
	u.SetUID(newUUID())
	if u.HasRole("o") {
		Props.Set("owner", u.UUID)
	}
}

func (u *User) SetNickname(s string) {
	db.Build().Up(cTableUsers, "nickname", s).Wh("uuid", u.UUID).Exe()
	u.Nickname = s
}
