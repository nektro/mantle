package db

import (
	"database/sql"
	"strconv"

	dbstorage "github.com/nektro/go.dbstorage"
	"github.com/nektro/go.etc/store"

	. "github.com/nektro/go.etc/dbt"
)

type Role struct {
	ID                 int64  `json:"id"`
	UUID               UUID   `json:"uuid" dbsorm:"1"`
	Position           int    `json:"position" dbsorm:"1"`
	Name               string `json:"name" dbsorm:"1"`
	Color              string `json:"color" dbsorm:"1"`
	PermManageChannels Perm   `json:"perm_manage_channels" dbsorm:"1"`
	PermManageRoles    Perm   `json:"perm_manage_roles" dbsorm:"1"`
	Distinguish        bool   `json:"distinguish" dbsorm:"1"`
	PermManageServer   Perm   `json:"perm_manage_server" dbsorm:"1"`
	PermManageInvites  Perm   `json:"perm_manage_invites" dbsorm:"1"`
	CreatedOn          Time   `json:"created_on" dbsorm:"1"`
	PermViewAudits     Perm   `json:"perm_view_audits" dbsorm:"1"`
}

//
//

func CreateRole(name string) *Role {
	store.This.Lock()
	defer store.This.Unlock()
	//
	id := db.QueryNextID(cTableRoles)
	uid := NewUUID()
	p := PermIgnore
	co := now()
	r := &Role{id, uid, int(id), name, "", p, p, false, p, p, co, p}
	db.Build().InsI(cTableRoles, r).Exe()
	Props.Increment("count_" + cTableRoles)
	return r
}

// QueryRoleByUID finds a Role with the specified uid
func QueryRoleByUID(uid UUID) (*Role, bool) {
	rl, ok := BuiltInRoles[uid]
	if ok {
		return rl, true
	}
	ch, ok := dbstorage.ScanFirst(db.Build().Se("*").Fr(cTableRoles).Wh("uuid", uid.String()), Role{}).(*Role)
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

func (v *Role) i() string {
	return v.UUID.String()
}

func (v Role) t() string {
	return cTableRoles
}

func (v Role) b() dbstorage.QueryBuilder {
	return db.Build().Se("*").Fr(v.t())
}


//
// searchers
//

//
// modifiers
//

// SetName sets name
func (v *Role) SetName(s string) {
	doUp(v, "name", s)
	v.Name = s
}

// SetColor sets color
func (v *Role) SetColor(s string) {
	doUp(v, "color", s)
	v.Color = s
}

// SetPosition sets position
func (v *Role) SetPosition(n int) {
	doUp(v, "position", strconv.Itoa(n))
	v.Position = n
}

// SetDistinguish sets position
func (v *Role) SetDistinguish(b bool) {
	doUp(v, "distinguish", strconv.FormatBool(b))
	v.Distinguish = b
}

// SetPermMngServer sets
func (v *Role) SetPermMngServer(p Perm) {
	doUp(v, "perm_manage_server", strconv.Itoa(int(p)))
	v.PermManageServer = p
}

// SetPermMngChannels sets
func (v *Role) SetPermMngChannels(p Perm) {
	doUp(v, "perm_manage_channels", strconv.Itoa(int(p)))
	v.PermManageChannels = p
}

// SetPermMngRoles sets
func (v *Role) SetPermMngRoles(p Perm) {
	doUp(v, "perm_manage_roles", strconv.Itoa(int(p)))
	v.PermManageRoles = p
}

// SetPermMngInvites sets
func (v *Role) SetPermMngInvites(p Perm) {
	doUp(v, "perm_manage_invites", strconv.Itoa(int(p)))
	v.PermManageInvites = p
}

// Delete removes this item from the database
// and returns the list of users that got the role removed
func (v *Role) Delete() []*User {
	v.MoveTo(len(Role{}.All()))
	arr := dbstorage.ScanAll(User{}.b().WR("roles", "like", "'%'||?||'%'", true, v.UUID), User{})
	aru := make([]*User, len(arr))
	for i, item := range arr {
		u := item.(*User)
		u.RemoveRole(v.UUID)
		aru[i] = u
	}
	doDel(v)
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
func (v *Role) SetPermViewAudits(p Perm) {
	doUp(v, "perm_view_audits", strconv.Itoa(int(p)))
	v.PermManageInvites = p
}
