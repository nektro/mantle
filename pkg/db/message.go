package db

import (
	"database/sql"

	dbstorage "github.com/nektro/go.dbstorage"
	"github.com/nektro/go.etc/store"

	. "github.com/nektro/go.etc/dbt"
)

type Message struct {
	ID   int64  `json:"id"`
	UUID string `json:"uuid" dbsorm:"1"`
	At   Time   `json:"time" dbsorm:"1"`
	By   string `json:"author" dbsorm:"1"`
	Body string `json:"body" dbsorm:"1"`
}

//
//

func CreateMessage(user *User, channel *Channel, body string) *Message {
	store.This.Lock()
	defer store.This.Unlock()
	//
	m := &Message{db.QueryNextID(cTableMessagesPrefix + channel.UUID), newUUID(), now(), user.UUID, body}
	if channel.HistoryOff {
		return m
	}
	db.Build().InsI(cTableMessagesPrefix+channel.UUID, m).Exe()
	db.Build().Up(cTableChannels, "latest_message", m.UUID).Wh("uuid", channel.UUID).Exe()
	return m
}

//
//

// Scan implements dbstorage.Scannable
func (v Message) Scan(rows *sql.Rows) dbstorage.Scannable {
	rows.Scan(&v.ID, &v.UUID, &v.At, &v.By, &v.Body)
	return &v
}
