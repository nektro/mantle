package db

import (
	"database/sql"
	"strconv"

	"github.com/nektro/mantle/pkg/store"

	dbstorage "github.com/nektro/go.dbstorage"
)

type Channel struct {
	ID          int64  `json:"id"`
	UUID        string `json:"uuid" sqlite:"text"`
	Position    int    `json:"position" sqlite:"int"`
	Name        string `json:"name" sqlite:"text"`
	Description string `json:"description" sqlite:"text"`
	HistoryOff  bool   `json:"history_off" sqlite:"tinyint(1)"`
	LatestMsg   string `json:"latest_message" sqlite:"text"`
	CreatedOn   Time   `json:"created_on" sqlite:"text"`
}

//
//

func CreateChannel(name string) *Channel {
	store.This.Lock()
	defer store.This.Unlock()
	//
	id := db.QueryNextID(cTableChannels)
	uid := newUUID()
	co := now()
	ch := &Channel{id, uid, int(id), name, "", false, "", co}
	db.Build().InsI(cTableChannels, ch).Exe()
	ch.AssertMessageTableExists()
	Props.Increment("count_" + cTableChannels)
	return ch
}

func QueryChannelByUUID(uid string) (*Channel, bool) {
	ch, ok := dbstorage.ScanFirst(db.Build().Se("*").Fr(cTableChannels).Wh("uuid", uid), Channel{}).(*Channel)
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

func (c *Channel) AssertMessageTableExists() {
	db.CreateTableStruct(cTableMessagesPrefix+c.UUID, Message{})
}

// QueryMsgAfterUID runs 'select * from messages where uuid < ? order by uuid desc limit 50'
func (c *Channel) QueryMsgAfterUID(uid string, limit int) []*Message {
	res := []*Message{}
	qb := db.Build()
	qb.Se("*")
	qb.Fr(cTableMessagesPrefix + c.UUID)
	if len(uid) > 0 {
		if IsUID(uid) {
			qb.Wr("uuid", "<=", uid)
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
	db.Build().Up(cTableChannels, "name", s).Wh("uuid", v.UUID).Exe()
	v.Name = s
}

// SetPosition sets position
func (v *Channel) SetPosition(n int) {
	db.Build().Up(cTableChannels, "position", strconv.Itoa(n)).Wh("uuid", v.UUID).Exe()
	v.Position = n
}

// SetDescription sets description
func (v *Channel) SetDescription(s string) {
	db.Build().Up(cTableChannels, "description", s).Wh("uuid", v.UUID).Exe()
	v.Description = s
}

// EnableHistory sets position
func (v *Channel) EnableHistory(b bool) {
	db.Build().Up(cTableChannels, "history_off", strconv.FormatBool(!b)).Wh("uuid", v.UUID).Exe()
	v.HistoryOff = !b
}

// Delete removes this item from the database
func (v *Channel) Delete() {
	db.Build().Del(cTableChannels).Wh("uuid", v.UUID).Exe()
	db.DropTable(cTableMessagesPrefix + v.UUID)
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
	return db.QueryRowCount(cTableMessagesPrefix + v.UUID)
}
