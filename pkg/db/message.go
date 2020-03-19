package db

import (
	"database/sql"
	"strings"

	"github.com/nektro/go-util/alias"
	dbstorage "github.com/nektro/go.dbstorage"
)

type Message struct {
	ID   int64  `json:"id"`
	UUID string `json:"uuid" sqlite:"text"`
	At   string `json:"time" sqlite:"text"`
	By   string `json:"author" sqlite:"text"`
	Body string `json:"body" sqlite:"text"`
}

//
//

func CreateMessage(user *User, channel *Channel, body string) *Message {
	dbstorage.InsertsLock.Lock()
	defer dbstorage.InsertsLock.Unlock()
	m := &Message{
		db.QueryNextID(cTableMessagesPrefix + channel.UUID),
		newUUID(),
		alias.T(),
		user.UUID,
		body,
	}
	if channel.HistoryOff {
		m.At = sUTCto3339(m.At)
		return m
	}
	db.Build().Ins(cTableMessagesPrefix+channel.UUID, m.ID, m.UUID, m.At, m.By, m.Body).Exe()
	db.Build().Up(cTableChannels, "latest_message", m.UUID).Wh("uuid", channel.UUID).Exe()
	m.At = sUTCto3339(m.At)
	return m
}

//
//

// Scan implements dbstorage.Scannable
func (v Message) Scan(rows *sql.Rows) dbstorage.Scannable {
	rows.Scan(&v.ID, &v.UUID, &v.At, &v.By, &v.Body)
	v.At = sUTCto3339(v.At)
	return &v
}
