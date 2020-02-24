package db

import (
	"database/sql"

	"github.com/nektro/go-util/util"
	dbstorage "github.com/nektro/go.dbstorage"
)

type Channel struct {
	ID          int64  `json:"id"`
	UUID        string `json:"uuid" sqlite:"text"`
	Position    int    `json:"position" sqlite:"int"`
	Name        string `json:"name" sqlite:"text"`
	Description string `json:"description" sqlite:"text"`
}

//
//

func CreateChannel(name string) *Channel {
	id := db.QueryNextID(cTableChannels)
	uid := newUUID()
	util.Log("[channel-create]", uid, "#"+name)
	ch := &Channel{id, uid, int(id), name, ""}
	db.QueryPrepared(true, "insert into "+cTableChannels+" values (?, ?, ?, ?, ?)", id, uid, id, name, "")
	ch.AssertMessageTableExists()
	return ch
}

//
//

func (v Channel) Scan(rows *sql.Rows) dbstorage.Scannable {
	rows.Scan(&v.ID, &v.UUID, &v.Position, &v.Name, &v.Description)
	return &v
}

func (v Channel) All() []*Channel {
	arr := dbstorage.ScanAll(db.Build().Se("*").Fr(cTableChannels), Channel{})
	res := []*Channel{}
	for _, item := range arr {
		res = append(res, item.(*Channel))
	}
	return res
}

func (c *Channel) AssertMessageTableExists() {
	db.CreateTableStruct(cTableMessagesPrefix+c.UUID, Message{})
}
