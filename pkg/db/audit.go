package db

import (
	"database/sql"
	"strconv"

	"github.com/nektro/mantle/pkg/store"

	dbstorage "github.com/nektro/go.dbstorage"
)

type Audit struct {
	ID        int64  `json:"id"`
	UUID      string `json:"uuid" sqlite:"text"`
	CreatedOn Time   `json:"created_on" sqlite:"text"`
	Action    Action `json:"action" sqlite:"smallint"`
	Agent     string `json:"agent" sqlite:"text"`
	Affected  string `json:"affected" sqlite:"text"`
	Key       string `json:"a_key" sqlite:"text"`
	Value     string `json:"a_value" sqlite:"text"`
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
