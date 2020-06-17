package db

import (
	"database/sql"
	"strconv"

	"github.com/nektro/mantle/pkg/store"

	"github.com/nektro/go-util/util"
	dbstorage "github.com/nektro/go.dbstorage"
)

type Invite struct {
	ID         int64  `json:"id"`
	UUID       string `json:"uuid" sqlite:"text"`
	CreatedOn  Time   `json:"created_on" sqlite:"text"`
	Code       string `json:"name" sqlite:"text"`
	Uses       int64  `json:"uses" sqlite:"int"`
	MaxUses    int64  `json:"max_uses" sqlite:"int"`
	Mode       int    `json:"mode" sqlite:"int"`
	ExpiresIn  string `json:"expires_in" sqlite:"text"`
	ExpiresOn  Time   `json:"expires_on" sqlite:"text"`
	IsFrozen   bool   `json:"is_frozen" sqlite:"tinyint(1)"`
	GivenRoles Array  `json:"given_roles" sqlite:"string"`
}

//
//

// CreateInvite creates a new permanent invite and returns it
func CreateInvite() *Invite {
	store.This.Lock()
	defer store.This.Unlock()
	//
	id := db.QueryNextID(cTableInvites)
	uid := newUUID()
	co := now()
	code := util.RandomString(8)
	n := &Invite{id, uid, co, code, 0, 0, 0, "", NewTime(timeZero), false, Array{}}
	db.Build().InsI(cTableInvites, n).Exe()
	Props.Increment("count_" + cTableInvites)
	return n
}

// QueryInviteByCode does exactly that
func QueryInviteByCode(c string) (*Invite, bool) {
	ch, ok := dbstorage.ScanFirst(db.Build().Se("*").Fr(cTableInvites).Wh("name", c), Invite{}).(*Invite)
	return ch, ok
}

// QueryInviteByUID does exactly that
func QueryInviteByUID(uid string) (*Invite, bool) {
	ch, ok := dbstorage.ScanFirst(db.Build().Se("*").Fr(cTableInvites).Wh("uuid", uid), Invite{}).(*Invite)
	return ch, ok
}

//
//

// Scan implements dbstorage.Scannable
func (v Invite) Scan(rows *sql.Rows) dbstorage.Scannable {
	rows.Scan(&v.ID, &v.UUID, &v.CreatedOn, &v.Code, &v.Uses, &v.MaxUses, &v.Mode, &v.ExpiresIn, &v.ExpiresOn, &v.IsFrozen, &v.GivenRoles)
	return &v
}

// All queries database for all currently existing Invites
func (v Invite) All() []*Invite {
	arr := dbstorage.ScanAll(db.Build().Se("*").Fr(cTableInvites), Invite{})
	res := []*Invite{}
	for _, item := range arr {
		res = append(res, item.(*Invite))
	}
	return res
}

//
//

// Use increments Uses by 1
func (v *Invite) Use(u *User) {
	v.Uses++
	db.Build().Up(cTableInvites, "uses", strconv.FormatInt(v.Uses, 10)).Wh("uuid", v.UUID).Exe()
	CreateAudit(ActionInviteUse, u, v.UUID, v.Code, "")
	u.SetAsMember(true)
}

// SetMaxUses sets
func (v *Invite) SetMaxUses(p int64) {
	db.Build().Up(cTableInvites, "max_uses", strconv.FormatInt(p, 10)).Wh("uuid", v.UUID).Exe()
	v.MaxUses = p
}

// Delete removes this item from the database
func (v *Invite) Delete() {
	db.Build().Del(cTableInvites).Wh("uuid", v.UUID).Exe()
	Props.Decrement("count_" + cTableInvites)
}
