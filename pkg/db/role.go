package db

import (
	"database/sql"
	"strconv"

	"github.com/nektro/go-util/util"
	dbstorage "github.com/nektro/go.dbstorage"
)

type Role struct {
	ID                 int64  `json:"id"`
	UUID               string `json:"uuid" sqlite:"text"`
	Position           int    `json:"position" sqlite:"int"`
	Name               string `json:"name" sqlite:"text"`
	Color              string `json:"color" sqlite:"text"`
	PermManageChannels uint8  `json:"perm_manage_channels" sqlite:"tinyint(1)"`
	PermManageRoles    uint8  `json:"perm_manage_roles" sqlite:"tinyint(1)"`
	Distinguish        bool   `json:"distinguish" sqlite:"tinyint(1)"`
	PermManageServer   uint8  `json:"perm_manage_server" sqlite:"tinyint(1)"`
}

//
//

func CreateRole(name string) *Role {
	id := db.QueryNextID(cTableRoles)
	uid := newUUID()
	util.Log("[role-create]", uid, name)
	p := PermIgnore
	db.Build().Ins(cTableRoles, id, uid, id, name, "", p, p, false, p).Exe()
	p8 := uint8(p)
	r := &Role{id, uid, int(id), name, "", p8, p8, false, p8}
	return r
}

// QueryRoleByUID finds a Role with the specified uid
func QueryRoleByUID(uid string) (*Role, bool) {
	rl, ok := BuiltInRoles[uid]
	if ok {
		return rl, true
	}
	ch, ok := dbstorage.ScanFirst(db.Build().Se("*").Fr(cTableRoles).Wh("uuid", uid), Role{}).(*Role)
	return ch, ok
}

//
//

// Scan implements dbstorage.Scannable
func (v Role) Scan(rows *sql.Rows) dbstorage.Scannable {
	rows.Scan(&v.ID, &v.UUID, &v.Position, &v.Name, &v.Color, &v.PermManageChannels, &v.PermManageRoles, &v.Distinguish, &v.PermManageServer)
	return &v
}

// All queries database for all currently existing Roles
func (v Role) All() []*Role {
	arr := dbstorage.ScanAll(db.Build().Se("*").Fr(cTableRoles), Role{})
	res := []*Role{}
	for _, item := range arr {
		res = append(res, item.(*Role))
	}
	return res
}

// AllSorted is the same as All but ordered by position
func (v Role) AllSorted() []*Role {
	arr := dbstorage.ScanAll(db.Build().Se("*").Fr(cTableRoles).Or("position", "asc"), Role{})
	res := []*Role{}
	for _, item := range arr {
		res = append(res, item.(*Role))
	}
	return res
}

// SetName sets name
func (v *Role) SetName(s string) {
	db.Build().Up(cTableRoles, "name", s).Wh("uuid", v.UUID).Exe()
	v.Name = s
}

// SetColor sets color
func (v *Role) SetColor(s string) {
	db.Build().Up(cTableRoles, "color", s).Wh("uuid", v.UUID).Exe()
	v.Color = s
}

// SetPosition sets position
func (v *Role) SetPosition(n int) {
	db.Build().Up(cTableRoles, "position", strconv.Itoa(n)).Wh("uuid", v.UUID).Exe()
	v.Position = n
}

// SetDistinguish sets position
func (v *Role) SetDistinguish(b bool) {
	db.Build().Up(cTableRoles, "distinguish", strconv.FormatBool(b)).Wh("uuid", v.UUID).Exe()
	v.Distinguish = b
}
