package db

import (
	"database/sql"

	"github.com/nektro/mantle/pkg/store"

	dbstorage "github.com/nektro/go.dbstorage"
)

type Message struct {
	ID   int64  `json:"id"`
	UUID string `json:"uuid" sqlite:"text"`
	At   Time   `json:"time" sqlite:"text"`
	By   string `json:"author" sqlite:"text"`
	Body string `json:"body" sqlite:"text"`
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
