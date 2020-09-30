package db

import (
	"database/sql"

	dbstorage "github.com/nektro/go.dbstorage"
	"github.com/nektro/go.etc/store"

	. "github.com/nektro/go.etc/dbt"
)

type Message struct {
	ID   int64  `json:"id"`
	UUID UUID   `json:"uuid" dbsorm:"1"`
	At   Time   `json:"time" dbsorm:"1"`
	By   UUID   `json:"author" dbsorm:"1"`
	Body string `json:"body" dbsorm:"1"`
}

//
//

func CreateMessage(user *User, channel *Channel, body string) *Message {
	store.This.Lock()
	defer store.This.Unlock()
	//
	tn := cTableMessagesPrefix + channel.i()
	m := &Message{db.QueryNextID(tn), NewUUID(), now(), user.UUID, body}
	if channel.HistoryOff {
		return m
	}
	db.Build().InsI(tn, m).Exe()
	db.Build().Up(cTableChannels, "latest_message", m.i()).Wh("uuid", channel.i()).Exe()
	Props.Increment("count_" + tn)
	return m
}

// Scan implements dbstorage.Scannable
func (v Message) Scan(rows *sql.Rows) dbstorage.Scannable {
	rows.Scan(&v.ID, &v.UUID, &v.At, &v.By, &v.Body)
	return &v
}

func (v *Message) i() string {
	return v.UUID.String()
}

//
// searchers
//

//
// modifiers
//
