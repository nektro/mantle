package db

import (
	"database/sql"
	"strconv"

	"github.com/nektro/mantle/pkg/store"

	dbstorage "github.com/nektro/go.dbstorage"
)

type Role struct {
	ID                 int64  `json:"id"`
	UUID               string `json:"uuid" sqlite:"text"`
	Position           int    `json:"position" sqlite:"int"`
	Name               string `json:"name" sqlite:"text"`
	Color              string `json:"color" sqlite:"text"`
	PermManageChannels Perm   `json:"perm_manage_channels" sqlite:"tinyint(1)"`
	PermManageRoles    Perm   `json:"perm_manage_roles" sqlite:"tinyint(1)"`
	Distinguish        bool   `json:"distinguish" sqlite:"tinyint(1)"`
	PermManageServer   Perm   `json:"perm_manage_server" sqlite:"tinyint(1)"`
	PermManageInvites  Perm   `json:"perm_manage_invites" sqlite:"tinyint(1)"`
	CreatedOn          Time   `json:"created_on" sqlite:"text"`
	PermViewAudits     Perm   `json:"perm_view_audits" sqlite:"tinyint(1)"`
}

//
//

func CreateRole(name string) *Role {
	store.This.Lock()
	defer store.This.Unlock()
	//
	id := db.QueryNextID(cTableRoles)
	uid := newUUID()
	p := PermIgnore
	co := now()
	r := &Role{id, uid, int(id), name, "", p, p, false, p, p, co, p}
	db.Build().InsI(cTableRoles, r).Exe()
	Props.Increment("count_" + cTableRoles)
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
	rows.Scan(&v.ID, &v.UUID, &v.Position, &v.Name, &v.Color, &v.PermManageChannels, &v.PermManageRoles, &v.Distinguish, &v.PermManageServer, &v.PermManageInvites, &v.CreatedOn,
		&v.PermViewAudits)
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

//
//

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

// SetPermMngServer sets
func (v *Role) SetPermMngServer(p int) {
	db.Build().Up(cTableRoles, "perm_manage_server", strconv.Itoa(p)).Wh("uuid", v.UUID).Exe()
	v.PermManageServer = Perm(p)
}

// SetPermMngChannels sets
func (v *Role) SetPermMngChannels(p int) {
	db.Build().Up(cTableRoles, "perm_manage_channels", strconv.Itoa(p)).Wh("uuid", v.UUID).Exe()
	v.PermManageChannels = Perm(p)
}

// SetPermMngRoles sets
func (v *Role) SetPermMngRoles(p int) {
	db.Build().Up(cTableRoles, "perm_manage_roles", strconv.Itoa(p)).Wh("uuid", v.UUID).Exe()
	v.PermManageRoles = Perm(p)
}

// SetPermMngInvites sets
func (v *Role) SetPermMngInvites(p int) {
	db.Build().Up(cTableRoles, "perm_manage_invites", strconv.Itoa(p)).Wh("uuid", v.UUID).Exe()
	v.PermManageInvites = Perm(p)
}

// Delete removes this item from the database
// and returns the list of users that got the role removed
func (v *Role) Delete() []*User {
	v.MoveTo(len(Role{}.All()))
	arr := dbstorage.ScanAll(db.Build().Se("*").Fr(cTableUsers).WR("roles", "like", "'%'||?||'%'", true, v.UUID), User{})
	aru := make([]*User, len(arr))
	for i, item := range arr {
		u := item.(*User)
		u.RemoveRole(v.UUID)
		aru[i] = u
	}
	db.Build().Del(cTableRoles).Wh("uuid", v.UUID).Exe()
	Props.Decrement("count_" + cTableRoles)
	return aru
}

// MoveTo sets position cleanly
func (v *Role) MoveTo(n int) {
	pH, pL := uHighLow(v.Position, n)
	allR := Role{}.AllSorted()
	for i, item := range allR {
		o := i + 1
		if o < pL {
			continue
		}
		if o > pH {
			continue
		}
		// role moving down
		if pL == v.Position {
			if o == pL {
				continue
			}
			if o == pH {
				v.SetPosition(n)
				continue
			}
			item.SetPosition(o - 1)
		}
		// role moving up
		if pL == n {
			if o == pH {
				v.SetPosition(n)
				continue
			}
			item.SetPosition(o + 1)
		}
	}
}

// SetPermViewAudits sets
func (v *Role) SetPermViewAudits(p int) {
	db.Build().Up(cTableRoles, "perm_view_audits", strconv.Itoa(p)).Wh("uuid", v.UUID).Exe()
	v.PermManageInvites = Perm(p)
}
