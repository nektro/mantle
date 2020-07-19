package db

import (
	"database/sql"
	"strconv"

	dbstorage "github.com/nektro/go.dbstorage"
	"github.com/nektro/go.etc/store"

	. "github.com/nektro/go.etc/dbt"
)

type Channel struct {
	ID          int64  `json:"id"`
	UUID        UUID   `json:"uuid" dbsorm:"1"`
	Position    int    `json:"position" dbsorm:"1"`
	Name        string `json:"name" dbsorm:"1"`
	Description string `json:"description" dbsorm:"1"`
	HistoryOff  bool   `json:"history_off" dbsorm:"1"`
	LatestMsg   UUID   `json:"latest_message" dbsorm:"1"`
	CreatedOn   Time   `json:"created_on" dbsorm:"1"`
}

//
//

func CreateChannel(name string) *Channel {
	store.This.Lock()
	defer store.This.Unlock()
	//
	id := db.QueryNextID(cTableChannels)
	uid := NewUUID()
	co := now()
	ch := &Channel{id, uid, int(id), name, "", false, "", co}
	db.Build().InsI(cTableChannels, ch).Exe()
	ch.AssertMessageTableExists()
	Props.Increment("count_" + cTableChannels)
	return ch
}

func QueryChannelByUUID(uid UUID) (*Channel, bool) {
	ch, ok := dbstorage.ScanFirst(db.Build().Se("*").Fr(cTableChannels).Wh("uuid", uid.String()), Channel{}).(*Channel)
	return ch, ok
}

//
//

// Scan implements dbstorage.Scannable
func (v Channel) Scan(rows *sql.Rows) dbstorage.Scannable {
	rows.Scan(&v.ID, &v.UUID, &v.Position, &v.Name, &v.Description, &v.HistoryOff, &v.LatestMsg, &v.CreatedOn)
	return &v
}

// All returns an array of all channels sorted by their position
func (v Channel) All() []*Channel {
	arr := dbstorage.ScanAll(db.Build().Se("*").Fr(cTableChannels).Or("position", "asc"), Channel{})
	res := []*Channel{}
	for _, item := range arr {
		res = append(res, item.(*Channel))
	}
	return res
}

//
//

func (v *Channel) i() string {
	return v.UUID.String()
}

func (v Channel) t() string {
	return cTableChannels
}

func (v Channel) m() string {
	return cTableMessagesPrefix + v.i()
}

func (v Channel) b() dbstorage.QueryBuilder {
	return db.Build().Se("*").Fr(v.t())
}

func (c *Channel) AssertMessageTableExists() {
	db.CreateTableStruct(c.m(), Message{})
}

// QueryMsgAfterUID runs 'select * from messages where uuid < ? order by uuid desc limit 50'
func (c *Channel) QueryMsgAfterUID(uid UUID, limit int) []*Message {
	res := []*Message{}
	qb := db.Build()
	qb.Se("*")
	qb.Fr(c.m())
	if len(uid) > 0 {
		if IsUUID(uid) {
			qb.Wr("uuid", "<=", uid.String())
		}
	}
	qb.Or("uuid", "desc")
	qb.Lm(int64(limit))
	arr := dbstorage.ScanAll(qb, Message{})
	for _, item := range arr {
		res = append(res, item.(*Message))
	}
	return res
}

// SetName sets name
func (v *Channel) SetName(s string) {
	doUp(v, "name", s)
	v.Name = s
}

// SetPosition sets position
func (v *Channel) SetPosition(n int) {
	doUp(v, "position", strconv.Itoa(n))
	v.Position = n
}

// SetDescription sets description
func (v *Channel) SetDescription(s string) {
	doUp(v, "description", s)
	v.Description = s
}

// EnableHistory sets position
func (v *Channel) EnableHistory(b bool) {
	doUp(v, "history_off", strconv.FormatBool(!b))
	v.HistoryOff = !b
}

// Delete removes this item from the database
func (v *Channel) Delete() {
	doDel(v)
	db.DropTable(v.m())
	Props.Decrement("count_" + cTableChannels)
}

// MoveTo sets position cleanly
func (v *Channel) MoveTo(n int) {
	pH, pL := uHighLow(v.Position, n)
	allC := Channel{}.All()
	for i, item := range allC {
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

func (v *Channel) MessageCount() int64 {
	return db.QueryRowCount(v.m())
}
