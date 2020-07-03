package db

import (
	"database/sql"
	"strconv"

	dbstorage "github.com/nektro/go.dbstorage"
	"github.com/nektro/go.etc/store"
)

type Audit struct {
	ID        int64  `json:"id"`
	UUID      string `json:"uuid" dbsorm:"1"`
	CreatedOn Time   `json:"created_on" dbsorm:"1"`
	Action    Action `json:"action" dbsorm:"1"`
	Agent     string `json:"agent" dbsorm:"1"`
	Affected  string `json:"affected" dbsorm:"1"`
	Key       string `json:"a_key" dbsorm:"1"`
	Value     string `json:"a_value" dbsorm:"1"`
}

//
//

func CreateAudit(ac Action, agent *User, aff string, key, val string) *Audit {
	store.This.Lock()
	defer store.This.Unlock()
	//
	id := db.QueryNextID(cTableAudits)
	uid := newUUID()
	co := now()
	a := &Audit{id, uid, co, ac, agent.UUID, aff, key, val}
	db.Build().InsI(cTableAudits, a).Exe()
	Props.Increment("count_" + cTableAudits)
	Props.Increment("count_" + cTableAudits + "_action_" + strconv.Itoa(int(ac)))
	return a
}

// Scan implements dbstorage.Scannable
func (v Audit) Scan(rows *sql.Rows) dbstorage.Scannable {
	rows.Scan(&v.ID, &v.UUID, &v.CreatedOn, &v.Action, &v.Agent, &v.Affected, &v.Key, &v.Value)
	return &v
}

// All returns an array of all channels sorted by their position
func (v Audit) All() []*Audit {
	arr := dbstorage.ScanAll(db.Build().Se("*").Fr(cTableAudits), Audit{})
	res := []*Audit{}
	for _, item := range arr {
		res = append(res, item.(*Audit))
	}
	return res
}

//
//
