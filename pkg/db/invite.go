package db

import (
	"database/sql"
	"strconv"
	"time"

	"github.com/nektro/go-util/util"
	dbstorage "github.com/nektro/go.dbstorage"
	"github.com/nektro/go.etc/store"

	. "github.com/nektro/go.etc/dbt"
)

type Invite struct {
	ID         int64    `json:"id"`
	UUID       UUID     `json:"uuid" dbsorm:"1"`
	CreatedOn  Time     `json:"created_on" dbsorm:"1"`
	Code       string   `json:"name" dbsorm:"1"`
	Uses       int64    `json:"uses" dbsorm:"1"`
	MaxUses    int64    `json:"max_uses" dbsorm:"1"`
	Mode       int      `json:"mode" dbsorm:"1"`
	ExpiresIn  Duration `json:"expires_in" dbsorm:"1"`
	ExpiresOn  Time     `json:"expires_on" dbsorm:"1"`
	IsFrozen   bool     `json:"is_frozen" dbsorm:"1"`
	GivenRoles List     `json:"given_roles" dbsorm:"1"`
}

//
//

// CreateInvite creates a new permanent invite and returns it
func CreateInvite() *Invite {
	store.This.Lock()
	defer store.This.Unlock()
	//
	id := db.QueryNextID(cTableInvites)
	uid := NewUUID()
	co := now()
	code := util.RandomString(8)
	n := &Invite{id, uid, co, code, 0, 0, 0, DurationZero, NewTime(TimeZero), false, List{}}
	db.Build().InsI(cTableInvites, n).Exe()
	Props.Increment("count_" + cTableInvites)
	return n
}

// QueryInviteByCode does exactly that
func QueryInviteByCode(c string) (*Invite, bool) {
	ch, ok := dbstorage.ScanFirst(Invite{}.b().Wh("name", c), Invite{}).(*Invite)
	return ch, ok
}

// QueryInviteByUID does exactly that
func QueryInviteByUID(uid UUID) (*Invite, bool) {
	ch, ok := dbstorage.ScanFirst(Invite{}.b().Wh("uuid", uid.String()), Invite{}).(*Invite)
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
	arr := dbstorage.ScanAll(v.b(), Invite{})
	res := []*Invite{}
	for _, item := range arr {
		res = append(res, item.(*Invite))
	}
	return res
}

//
//

func (v *Invite) i() string {
	return v.UUID.String()
}

func (v Invite) t() string {
	return cTableInvites
}

func (v Invite) b() dbstorage.QueryBuilder {
	return db.Build().Se("*").Fr(v.t())
}


//
// searchers
//

//
// modifiers
//

// Use increments Uses by 1
func (v *Invite) Use(u *User) {
	v.Uses++
	doUp(v, "uses", strconv.FormatInt(v.Uses, 10))
	CreateAudit(ActionInviteUse, u, v.UUID, v.Code, "")
	u.SetAsMember(true)
}

// SetMaxUses sets
func (v *Invite) SetMaxUses(p int64) {
	doUp(v, "max_uses", strconv.FormatInt(p, 10))
	v.MaxUses = p
}

// Delete removes this item from the database
func (v *Invite) Delete() {
	doDel(v)
	Props.Decrement("count_" + cTableInvites)
}

// SetMode does
func (v *Invite) SetMode(x int) {
	doUp(v, "mode", strconv.Itoa(x))
	v.Mode = x
}

// SetExpIn does
func (v *Invite) SetExpIn(x [2]int) {
	doUp(v, "expires_in", strconv.Itoa(x[0])+":"+strconv.Itoa(x[1]))
	v.ExpiresIn = Duration(x)
}

// SetExpOn does
func (v *Invite) SetExpOn(t time.Time) {
	doUp(v, "expires_on", t.Format(TimeFormat))
	v.ExpiresOn = Time(t)
}
