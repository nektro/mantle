package db

import (
	"database/sql"
	"strconv"
	"strings"

	"github.com/nektro/go-util/alias"
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
	JoindedOn  string `json:"joined_on" sqlite:"text"`
	LastActive string `json:"last_active" sqlite:"text"`
	Roles      string `json:"roles" sqlite:"text"`
	RolesA     []string
}

//
//

func QueryUserByUUID(uid string) (*User, bool) {
	ch, ok := dbstorage.ScanFirst(db.Build().Se("*").Fr(cTableUsers).Wh("uuid", uid), User{}).(*User)
	return ch, ok
}

func QueryUserBySnowflake(provider string, flake string, name string) *User {
	us := dbstorage.ScanFirst(db.Build().Se("*").Fr(cTableUsers).Wh("provider", provider).Wh("snowflake", flake), User{})
	if us != nil {
		return us.(*User)
	}
	// else
	id := db.QueryNextID(cTableUsers)
	uid := newUUID()
	now := alias.T()
	roles := ""
	if id == 1 {
		roles += "o"
		Props.Set("owner", uid)
	}
	db.Build().Ins(cTableUsers, id, provider, flake, uid, 0, 0, name, "", now, now, roles).Exe()
	return QueryUserBySnowflake(provider, flake, name)
}

//
//

// Scan implements dbstorage.Scannable
func (v User) Scan(rows *sql.Rows) dbstorage.Scannable {
	rows.Scan(&v.ID, &v.Provider, &v.Snowflake, &v.UUID, &v.IsMember, &v.IsBanned, &v.Name, &v.Nickname, &v.JoindedOn, &v.LastActive, &v.Roles)
	v.JoindedOn = strings.Replace(v.JoindedOn, " ", "T", 1) + "Z"
	v.LastActive = strings.Replace(v.LastActive, " ", "T", 1) + "Z"
	v.RolesA = strings.Split(v.Roles, ",")
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

func (u *User) SetAsMember(b bool) {
	db.Build().Up(cTableUsers, "is_member", strconv.Itoa(util.Btoi(b))).Wh("uuid", u.UUID).Exe()
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
	return util.Contains(u.RolesA, role)
}

func (u *User) AddRole(role string) {
	if u.HasRole(role) {
		return
	}
	u.RolesA = append(u.RolesA, role)
	u.Roles = strings.Join(u.RolesA, ",")
	db.Build().Up(cTableUsers, "roles", u.Roles).Wh("uuid", u.UUID).Exe()
}
